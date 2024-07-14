package messages

// file: claude/messages/messages.go

import (
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/rmrfslashbin/ami/claude"
	"github.com/rs/xid"
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
	conversationFile string
	url              string
	mu               sync.Mutex

	request MessageCreateParams
}

// New creates a new Messages configuration.
func New(opts ...func(*Messages)) (*Messages, error) {
	config := &Messages{
		conversation: &Conversation{
			ID:      xid.New().String(),
			Created: time.Now(),
			Updated: time.Now(),
		},
		url: URL,
	}

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

	config.conversation.Model = config.request.Model

	if config.conversationFile != "" {
		err := config.LoadConversation()
		if err != nil && !os.IsNotExist(err) {
			return nil, err
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

// WithConversationFile sets the file for storing the conversation.
func WithConversationFile(file string) Option {
	return func(config *Messages) {
		config.conversationFile = file
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

func WithSonnet35() Option {
	return withModel("sonnet35")
}

func withModel(model string) Option {
	return func(config *Messages) {
		model := claude.ModelsList[model]
		config.request.Model = model.Name
		config.request.MaxTokens = model.MaxOutputTokens
	}
}

func WithMaxTokens(n int) Option {
	return func(config *Messages) {
		config.request.MaxTokens = n
	}
}

func (m *Messages) SetStreaming(stream bool) {
	m.request.Stream = stream
}

// SetUserId sets the user id.
func (m *Messages) SetUserId(id string) {
	if m.request.Metadata == nil {
		m.request.Metadata = &Metadata{}
	}
	m.request.Metadata.UserID = &id
}

func (m *Messages) SetMaxTokens(n int) error {
	if n > m.request.MaxTokens {
		return &ErrMaxTokensExceeded{Model: m.request.Model, MaxTokens: m.request.MaxTokens}
	}
	m.request.MaxTokens = n
	return nil
}

func (m *Messages) SetSystemPrompt(p string) {
	m.request.System = &p
}

func (m *Messages) AddMessage(role string, content string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	message := &Message{
		Role: role,
		Content: []ContentBlock{
			{
				Type: "text",
				Text: &content,
			},
		},
	}

	// Validate the message before adding it
	if err := m.validateMessage(message); err != nil {
		return err
	}

	m.conversation.Messages = append(m.conversation.Messages, message)
	m.conversation.Updated = time.Now()

	return nil
}

func (m *Messages) AddRoleAssistant(content string) error {
	return m.AddMessage("assistant", content)
}

func (m *Messages) AddRoleUser(content string) error {
	return m.AddMessage("user", content)
}

func (m *Messages) AddRoleUserMedia(fqpn string, prompt string) error {
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

	message := &Message{
		Role: "user",
		Content: []ContentBlock{
			{
				Type: "image",
				Source: &ImageSource{
					Type:      "base64",
					MediaType: mtype.String(),
					Data:      base64Content,
				},
			},
			{
				Type: "text",
				Text: &prompt,
			},
		},
	}

	// Validate the message before adding it
	if err := m.validateMessage(message); err != nil {
		return err
	}

	m.conversation.Messages = append(m.conversation.Messages, message)
	m.conversation.Updated = time.Now()
	return nil
}

func (m *Messages) AddRoleUserToolResult(toolUseId string, content string) error {
	message := &Message{
		Role: "user",
		Content: []ContentBlock{
			{
				Type:      "tool_result",
				ToolUseID: &toolUseId,
				Content: []ContentBlockContent{
					{
						Type: "text",
						Text: &content,
					},
				},
			},
		},
	}
	m.conversation.Messages = append(m.conversation.Messages, message)
	m.conversation.Updated = time.Now()
	return nil
}

func (m *Messages) AddTool(tool *ToolParam) {
	m.request.Tools = append(m.request.Tools, tool)
}

// SetToolChoiceAuto sets the tool choice to auto.
func (m *Messages) SetToolChoiceAuto() {
	m.request.ToolChoice = &ToolChoice{Type: "auto"}
}

// SetToolChoiceAny sets the tool choice to any.
func (m *Messages) SetToolChoiceAny() {
	m.request.ToolChoice = &ToolChoice{Type: "any"}
}

// SetToolChoiceTool sets the tool choice to a specific tool.
func (m *Messages) SetToolChoiceTool(tool string) {
	m.request.ToolChoice = &ToolChoice{Type: "tool", Name: &tool}
}

func (m *Messages) GetMessageRequest() *MessageCreateParams {
	return &m.request
}

func (m *Messages) GetConversation() *Conversation {
	return m.conversation
}

func (m *Messages) Send() (*Message, error) {
	// Validate all messages in the conversation
	for _, msg := range m.conversation.Messages {
		if err := m.validateMessage(msg); err != nil {
			return nil, err
		}
	}

	// Convert conversation messages to MessageParams
	m.request.Messages = convertToMessageParams(m.conversation.Messages)

	jsonData, err := json.Marshal(m.request)
	if err != nil {
		return nil, &ErrMarshalingInput{Err: err}
	}

	resp, err := m.claud.Do(m.url, jsonData)
	if err != nil {
		return nil, err
	}

	var reply Message
	err = json.Unmarshal(*resp, &reply)
	if err != nil {
		return nil, &ErrMarshalingReply{Err: err}
	}

	m.conversation.Messages = append(m.conversation.Messages, &reply)
	m.conversation.Updated = time.Now()

	if m.conversationFile != "" {
		err = m.SaveConversation()
		if err != nil {
			return nil, err
		}
	}

	return &reply, nil
}

// validateMessage is a helper function to validate a single message
func (m *Messages) validateMessage(msg *Message) error {
	for _, block := range msg.Content {
		if err := block.Validate(); err != nil {
			return fmt.Errorf("invalid content block: %w", err)
		}
	}
	return nil
}

func (m *Messages) LoadConversation() error {
	file, err := os.Open(m.conversationFile)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(m.conversation)
}

func (m *Messages) SaveConversation() error {
	file, err := os.Create(m.conversationFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(m.conversation)
}

func (m *Messages) ResetConversation() {
	m.conversation = &Conversation{
		ID:      xid.New().String(),
		Model:   m.request.Model,
		Created: time.Now(),
		Updated: time.Now(),
	}
}

// GetModel returns the model name.
func (m *Messages) GetModel() string {
	return m.request.Model
}

func (c *ContentBlock) Validate() error {
	switch c.Type {
	case "text":
		if c.Text == nil || *c.Text == "" {
			return errors.New("text is required for text type content")
		}
	case "image":
		if c.Source == nil {
			return errors.New("source is required for image type content")
		}
		if err := c.Source.Validate(); err != nil {
			return err
		}
	case "tool_use":
		if c.ID == nil || c.Name == nil || c.Input == nil {
			return errors.New("id, name, and input are required for tool_use type content")
		}
	case "tool_result":
		if c.ToolUseID == nil || len(c.Content) == 0 {
			return errors.New("tool_use_id and content are required for tool_result type content")
		}
		for _, content := range c.Content {
			if err := content.Validate(); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("invalid content type: %s", c.Type)
	}
	return nil
}

func (is *ImageSource) Validate() error {
	if is == nil {
		return errors.New("image source cannot be nil")
	}
	if is.Type != "base64" {
		return errors.New("only base64 type is supported for image source")
	}
	validMediaTypes := []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
	if !slices.Contains(validMediaTypes, is.MediaType) {
		return fmt.Errorf("invalid media type: %s", is.MediaType)
	}
	if is.Data == "" {
		return errors.New("data is required for image source")
	}
	return nil
}

func (trc *ContentBlockContent) Validate() error {
	if trc == nil {
		return errors.New("content block content cannot be nil")
	}
	switch trc.Type {
	case "text":
		if trc.Text == nil || *trc.Text == "" {
			return errors.New("text is required for text type tool result content")
		}
	case "image":
		if trc.Source == nil {
			return errors.New("source is required for image type tool result content")
		}
		if err := trc.Source.Validate(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid tool result content type: %s", trc.Type)
	}
	return nil
}

// convertToMessageParams converts []*Message to []MessageParam
func convertToMessageParams(messages []*Message) []MessageParam {
	params := make([]MessageParam, len(messages))
	for i, msg := range messages {
		params[i] = MessageParam{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return params
}
