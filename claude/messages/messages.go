package messages

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"slices"

	"github.com/gabriel-vasile/mimetype"
	"github.com/rmrfslashbin/ami/claude"
	"github.com/rmrfslashbin/ami/claude/state"
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
}

// Option is a configuration option.
type Option func(config *Messages)

// Messages is the messages configuration.
type Messages struct {
	claud *claude.Claude
	state *state.State

	// Model is the model that will complete your prompt.
	// Required.
	// See models (https://docs.anthropic.com/claude/docs/models-overview) for additional details and options.
	Model *string `json:"model"`

	// modelMaxTokens is the maximum number of tokens for the model.
	modelMaxTokens int

	// Messages is the messages to send to the API.
	// Required.
	Messages []*Message `json:"messages"`

	// System is a system prompt is a way of providing context and instructions to Claude, such as specifying a particular goal or role.
	System string `json:"system,omitempty"`

	// MaxToken is the maximum number of tokens to generate before stopping.
	// Required.
	MaxTokens *int `json:"max_tokens"`

	// Metadata is an object describing metadata about the request.
	Metadata *Metadata `json:"metadata,omitempty"`

	// StopSequences is a list of strings that, if generated, will cause the model to stop generating tokens.
	StopSequences []string `json:"stop_sequences,omitempty"`

	// Stream is a boolean that indicates whether the model should generate a single response or a stream of responses.
	// Default is false.
	Stream bool `json:"stream"`

	// Temperature is a float that controls the randomness of the model's output. The higher the temperature, the more random the output.
	Temperature *float32 `json:"temperature,omitempty"`

	// TopP is an integer that controls nucleus sampling. The higher the top_p, the more diverse the output.
	// Recommended for advanced use cases only. You usually only need to use temperature.
	// You should either alter temperature or top_p, but not both.
	TopP *int `json:"top_p,omitempty"`

	// TopK is an integer that specifies sampling from the top K options for each subsequent token.
	// Recommended for advanced use cases only. You usually only need to use temperature.
	TopK *int `json:"top_k,omitempty"`
}

// New creates a new Messages configuration.
func New(opts ...func(*Messages)) (*Messages, error) {
	config := &Messages{}

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}

	if config.claud == nil {
		return nil, &ErrMissingClaude{}
	}

	if config.Model == nil {
		return nil, &ErrMissingModel{}
	}

	// Stream is not supported yet; always set to false.
	config.Stream = false

	if config.state != nil {
		for _, message := range config.state.Conversation.Messages {
			config.Messages = append(
				config.Messages,
				&Message{
					Role: message.Role,
					MessageContent: []*Content{
						{
							Type: "text",
							Text: message.Content,
						},
					},
				},
			)
		}
	}

	return config, nil
}

// WithClaude sets the Claude configuration.
func WithClaude(claud *claude.Claude) Option {
	return func(config *Messages) {
		config.claud = claud
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
		config.Model = &model.Name
		config.MaxTokens = &model.MaxOutputTokens
		config.modelMaxTokens = model.MaxOutputTokens
	}
}

func WithMaxTokens(n int) Option {
	return func(config *Messages) {
		config.MaxTokens = &n
	}
}

func WithState(state *state.State) Option {
	return func(config *Messages) {
		config.state = state
	}
}

// Streaming activates streaming mode.
// Streaming mode is not supported yet.
func (messages *Messages) Streaming() error {
	//request.Stream = true
	return &ErrUnsupportedOption{Err: errors.New("streaming not supported yet")}
}

// UserId sets the user id.
func (messages *Messages) SetUserId(id string) {
	messages.Metadata.UserId = id
}

func (messages *Messages) SetMaxTokens(n int) error {
	if n > *messages.MaxTokens {
		return &ErrMaxTokensExceeded{Model: *messages.Model, MaxTokens: *messages.MaxTokens}
	}
	messages.MaxTokens = &n
	return nil
}

func (messages *Messages) SetSystemPrompt(p string) {
	messages.System = p
}

func (messages *Messages) AddRoleAssistant(content string) {
	newMessage := &Message{
		Role: "assistant",
		MessageContent: []*Content{
			{
				Type: "text",
				Text: content,
			},
		},
	}
	messages.Messages = append(messages.Messages, newMessage)
	if messages.state != nil {
		messages.state.Conversation.Messages = append(messages.state.Conversation.Messages, &state.Message{Role: string(state.Assistant), Content: content})
	}
}

func (messages *Messages) AddRoleUser(content string) {
	newMessage := &Message{
		Role: "user",
		MessageContent: []*Content{
			{
				Type: "text",
				Text: content,
			},
		},
	}
	messages.Messages = append(messages.Messages, newMessage)
	if messages.state != nil {
		messages.state.Conversation.Messages = append(messages.state.Conversation.Messages, &state.Message{Role: string(state.User), Content: content})
	}
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

	messages.Messages = append(
		messages.Messages, &Message{
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

func (messages *Messages) Send() (*Response, error) {

	if *messages.MaxTokens > messages.modelMaxTokens {
		return nil, &ErrMaxTokensExceeded{Model: *messages.Model, MaxTokens: *messages.MaxTokens}
	}

	if messages.TopP != nil && messages.Temperature != nil {
		return nil, &ErrConflictingOptions{Err: errors.New("top_p and temperature")}
	}

	if messages.Stream {
		return nil, &ErrUnsupportedOption{Err: errors.New("streaming not supported yet")}
	}

	jsonData, err := json.Marshal(messages)
	if err != nil {
		return nil, &ErrMarshalingInput{Err: err}
	}

	resp, err := messages.claud.Do(URL, jsonData)
	if err != nil {
		return nil, err
	}

	var reply Response
	err = json.Unmarshal(*resp, &reply)
	if err != nil {
		return nil, &ErrMarshalingReply{Err: err}
	}

	if messages.state != nil {
		for _, content := range reply.Content {
			messages.state.Conversation.Messages = append(messages.state.Conversation.Messages, &state.Message{Role: string(state.Assistant), Content: content.Text})
		}
	}

	return &reply, nil
}
