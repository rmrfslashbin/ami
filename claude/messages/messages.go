package messages

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/rmrfslashbin/ami/claude"
	"github.com/tmaxmax/go-sse"
)

// URL is the URL for the Messages API.
const URL = claude.URL + "/v1/messages"

// Slice of supported mime types.
var SUPPORTED_MIME_TYPES = []string{"image/jpeg", "image/png", "image/gif", "image/webp"}

// StopReasons is a map of stop reasons.
var StopReasons = map[string]string{
	"end_turn":      "the model reached a natural stopping point",
	"max_tokens":    "we exceeded the requested max_tokens or the model's maximum",
	"stop_sequence": "one of your provided custom stop_sequences was generated",
	"tool_use":      "the model requests use of a tool",
}

// Option is a configuration option.
type Option func(config *Messages)

// Messages is the messages configuration.
type Messages struct {
	claud            *claude.Claude
	conversation     *Conversation
	conversationFqpn *string
	url              string

	request Request
}

// New creates a new Messages configuration.
func New(opts ...func(*Messages)) (*Messages, error) {
	config := &Messages{}
	config.request = Request{}

	config.conversation = &Conversation{}
	now := time.Now()
	config.conversation.Created = now
	config.conversation.Updated = now
	config.url = URL

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}

	if config.claud == nil {
		return nil, &ErrMissingClaude{}
	}

	if config.request.Model == "" {
		return nil, &ErrMissingModel{}
	}

	if config.conversationFqpn != nil {
		err := config.Load()
		if err != nil {
			return nil, err
		}
	}

	if config.conversation.Model == nil {
		config.conversation.Model = &config.request.Model
	}

	return config, nil
}

// WithClaude sets the Claude configuration.
func WithClaude(claud *claude.Claude) Option {
	return func(config *Messages) {
		config.claud = claud
	}
}

func WithConversationFile(fpqn *string) Option {
	if fpqn != nil {
		return func(config *Messages) {
			config.conversationFqpn = fpqn
		}
	}
	return func(config *Messages) {
		config.conversation = &Conversation{}
	}
}

func WithOpus() Option {
	return withModel("opus")
}

func WithSonnet() Option {
	return withModel("sonnet")
}

func WithHaiku() Option {
	return withModel("haiku")
}

func withModel(model string) Option {
	return func(config *Messages) {
		model := claude.ModelsList[model]
		config.request.Model = model.Name
		config.request.MaxTokens = model.MaxOutputTokens
		config.request.modelMaxTokens = model.MaxOutputTokens
	}
}

func WithMaxTokens(n int) Option {
	return func(config *Messages) {
		config.request.MaxTokens = n
	}
}

func (messages *Messages) SetStreaming(stream bool) {
	messages.request.Stream = stream
}

// UserId sets the user id.
func (messages *Messages) SetUserId(id string) {
	messages.request.Metadata.UserId = id
}

func (messages *Messages) SetMaxTokens(n int) error {
	if n > messages.request.MaxTokens {
		return &ErrMaxTokensExceeded{Model: messages.request.Model, MaxTokens: messages.request.MaxTokens}
	}
	messages.request.MaxTokens = n
	return nil
}

func (messages *Messages) SetSystemPrompt(p string) {
	messages.request.System = p
}

func (messages *Messages) AddRoleAssistant(content string) {
	messages.conversation.Messages = append(
		messages.conversation.Messages,
		&Message{
			Role: "assistant",
			MessageContent: []*Content{
				{
					Type: "text",
					Text: content,
				},
			},
		},
	)
}

func (messages *Messages) AddRoleUser(content string) {
	/*
		newMessage :=
		messages.Messages = append(messages.Messages, newMessage)
	*/

	messages.conversation.Messages = append(
		messages.conversation.Messages,
		&Message{
			Role: "user",
			MessageContent: []*Content{
				{
					Type: "text",
					Text: content,
				},
			},
		},
	)
}

func (messages *Messages) AddRoleUserMedia(fqpn string, prompt string) error {
	mtype, err := mimetype.DetectFile(fqpn)
	if err != nil {
		return &ErrFetchingMimeType{Err: err}
	}
	if !slices.Contains(SUPPORTED_MIME_TYPES, mtype.String()) {
		return &ErrUnsupportedMimeType{MimeType: mtype.String()}
	}

	// Read the file content
	content, err := os.ReadFile(fqpn)
	if err != nil {
		return &ErrReadingFile{Err: err}
	}

	// Convert the file content to a base64-encoded string
	base64Content := base64.StdEncoding.EncodeToString(content)

	messages.conversation.Messages = append(
		messages.conversation.Messages,
		&Message{
			Role: "user",
			MessageContent: []*Content{
				{
					Type: "image",
					Source: &MediaSource{
						Type:      "base64",
						MediaType: mtype.String(),
						Data:      base64Content,
					},
				},
				{Type: "text", Text: prompt},
			},
		},
	)
	return nil
}

func (messages *Messages) AddTool(tool *Tool) {
	messages.request.Tools = append(messages.request.Tools, tool)
}

// SetToolChoiceAuto sets the tool choice to auto.
func (messages *Messages) SetToolChoiceAuto() {
	messages.request.ToolChoice = &ToolChoice{Type: "auto"}
}

// SetToolChoiceAny sets the tool choice to any.
func (messages *Messages) SetToolChoiceAny() {
	messages.request.ToolChoice = &ToolChoice{Type: "any"}
}

// SetToolChoiceTool sets the tool choice to a specific tool.
func (messages *Messages) SetToolChoiceTool(tool string) {
	messages.request.ToolChoice = &ToolChoice{Type: "tool", Name: tool}
}

func (messages *Messages) GetMessageRequest() *Request {
	return &messages.request

}

func (messages *Messages) GetConversation() *Conversation {
	return messages.conversation
}

type StreamResults struct {
	Response <-chan StreamingMessageResponse
	Error    <-chan error
}

func (messages *Messages) Stream(ctx context.Context) StreamResults {
	responseCh := make(chan StreamingMessageResponse)
	errCh := make(chan error)

	if len(messages.request.Tools) > 0 {
		errCh <- &ErrToolUseNotSupported{}
		close(responseCh)
		return StreamResults{Response: responseCh, Error: errCh}
	}

	if messages.request.MaxTokens > messages.request.modelMaxTokens {
		errCh <- &ErrMaxTokensExceeded{Model: messages.request.Model, MaxTokens: messages.request.MaxTokens}
		close(responseCh)
		return StreamResults{Response: responseCh, Error: errCh}
	}

	if messages.request.TopP != nil && messages.request.Temperature != nil {
		errCh <- &ErrConflictingOptions{Err: errors.New("top_p and temperature")}
		close(responseCh)
		return StreamResults{Response: responseCh, Error: errCh}
	}

	// Load the conversation
	messages.request.Messages = messages.conversation.Messages

	jsonData, err := json.Marshal(messages.request)
	if err != nil {
		errCh <- &ErrMarshalingInput{Err: err}
		close(responseCh)
		return StreamResults{Response: responseCh, Error: errCh}
	}

	go func() {
		defer close(responseCh)

		//ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		//defer cancel()

		req, err := http.NewRequestWithContext(ctx, "POST", messages.url, bytes.NewBuffer(jsonData))
		if err != nil {
			errCh <- err
			return
		}

		// Set headers
		for key, value := range messages.claud.GetHeaders() {
			req.Header.Set(key, value)
		}

		conn := sse.DefaultClient.NewConnection(req)

		conn.SubscribeEvent("message_start", func(event sse.Event) {
			var response StreamingMessageStart
			err := json.Unmarshal([]byte(event.Data), &response)
			if err != nil {
				errCh <- &ErrMarshalingReply{Err: err}
				return
			}
			responseCh <- StreamingMessageResponse{MessageStart: &response}
		})
		conn.SubscribeEvent("content_block_delta", func(event sse.Event) {
			var response StreamingMessageContentBlockDelta
			err := json.Unmarshal([]byte(event.Data), &response)
			if err != nil {
				errCh <- &ErrMarshalingReply{Err: err}
				return
			}
			responseCh <- StreamingMessageResponse{ContentBlock: &response}
		})

		conn.SubscribeEvent("message_delta", func(event sse.Event) {
			var response StreamingMessageStop
			err := json.Unmarshal([]byte(event.Data), &response)
			if err != nil {
				errCh <- &ErrMarshalingReply{Err: err}
				return
			}
			responseCh <- StreamingMessageResponse{MessageStop: &response}
			close(responseCh)
		})

		conn.SubscribeEvent("error", func(event sse.Event) {
			var response StreamingMessageError
			err := json.Unmarshal([]byte(event.Data), &response)
			if err != nil {
				errCh <- &ErrMarshalingReply{Err: err}
				return
			}
			responseCh <- StreamingMessageResponse{StreamingError: &response}
			errCh <- &ErrStreamingMessage{}
		})

		// noops for now
		conn.SubscribeEvent("ping", func(event sse.Event) {})
		conn.SubscribeEvent("content_block_start", func(event sse.Event) {})
		conn.SubscribeEvent("content_block_stop", func(event sse.Event) {})
		conn.SubscribeEvent("message_stop", func(event sse.Event) {})

		if err := conn.Connect(); err != nil {
			/*
				log.LogAttrs(context.TODO(), slog.LevelError,
					"error connecting to streaming service",
					slog.String("error", err.Error()))
			*/
			errCh <- err
			return
		}

	}()

	return StreamResults{Response: responseCh, Error: errCh}
}

func (messages *Messages) Send() (*Response, error) {
	if err := messages.request.Validate(); err != nil {
		return nil, err
	}

	// Load the conversation
	messages.request.Messages = messages.conversation.Messages

	jsonData, err := json.Marshal(messages.request)
	if err != nil {
		return nil, &ErrMarshalingInput{Err: err}
	}

	resp, err := messages.claud.Do(messages.url, jsonData)
	if err != nil {
		return nil, err
	}

	var reply Response
	err = json.Unmarshal(*resp, &reply)
	if err != nil {
		return nil, &ErrMarshalingReply{Err: err}
	}

	messages.conversation.Messages = append(
		messages.conversation.Messages,
		&Message{Role: reply.Role, MessageContent: reply.Content},
	)

	// Reset the messages
	messages.request.Messages = nil

	return &reply, nil
}

func (messages *Messages) Load() error {
	var err error
	var fqpn string
	fqpn, err = filepath.Abs(*messages.conversationFqpn)
	if err != nil {
		return err
	}
	messages.conversationFqpn = &fqpn

	// Check if the file exists
	_, err = os.Stat(fqpn)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		// Load the file
		file, err := os.Open(fqpn)
		if err != nil {
			return &ErrOpeningFile{Err: err}
		}
		defer file.Close()

		// Decode the GOB data
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(messages.conversation)
		if err != nil {
			return &ErrLoadingGOB{Err: err}
		}
	}

	return nil
}

// Save the conversation
func (messages *Messages) Save() error {
	// if no file path is set, return
	if messages.conversationFqpn == nil {
		return nil
	}

	// Create the file
	file, err := os.Create(*messages.conversationFqpn)
	if err != nil {
		return &ErrOpeningFile{Err: err}
	}
	defer file.Close()

	now := time.Now()
	messages.conversation.Updated = now

	// Encode the GOB data
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(messages.conversation)
	if err != nil {
		return &ErrSavingGOB{Err: err}
	}

	return nil
}

// Reset the conversation
func (messages *Messages) Reset(modelName *string) {
	// Change the model if needed
	if modelName != nil {
		model := claude.ModelsList[*modelName]
		messages.request.Model = model.Name
		messages.request.MaxTokens = model.MaxOutputTokens
		messages.request.modelMaxTokens = model.MaxOutputTokens
	}

	// Reset the messages
	messages.conversation.Messages = nil
}

// GetModel returns the model name.
func (messages *Messages) GetModel() string {
	return messages.request.Model
}
