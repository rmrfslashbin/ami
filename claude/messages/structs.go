package messages

// file: claude/messages/structs.go

import (
	"time"

	"github.com/invopop/jsonschema"
)

// Conversation represents a conversation with Claude.
type Conversation struct {
	ID       string     `json:"id"`
	Model    string     `json:"model"`
	Messages []*Message `json:"messages"`
	Created  time.Time  `json:"created"`
	Updated  time.Time  `json:"updated"`
}

// Message represents a message in the conversation.
type Message struct {
	ID           string         `json:"id"`
	Role         string         `json:"role"`
	Content      []ContentBlock `json:"content"`
	Model        string         `json:"model"`
	StopReason   *string        `json:"stop_reason,omitempty"`
	StopSequence *string        `json:"stop_sequence,omitempty"`
	Type         string         `json:"type"`
	Usage        *Usage         `json:"usage"`
}

// ContentBlock represents a single content block within a message.
type ContentBlock struct {
	Type      string                `json:"type"`
	Text      *string               `json:"text,omitempty"`
	Source    *ImageSource          `json:"source,omitempty"`
	ID        *string               `json:"id,omitempty"`
	Name      *string               `json:"name,omitempty"`
	Input     interface{}           `json:"input,omitempty"`
	ToolUseID *string               `json:"tool_use_id,omitempty"`
	IsError   *bool                 `json:"is_error,omitempty"`
	Content   []ContentBlockContent `json:"content,omitempty"`
}

// ContentBlockContent represents the content of a tool result.
type ContentBlockContent struct {
	Type   string       `json:"type"`
	Text   *string      `json:"text,omitempty"`
	Source *ImageSource `json:"source,omitempty"`
}

// ImageSource represents the source of an image.
type ImageSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

// Usage represents the token usage information.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// MessageCreateParams represents the parameters for creating a new message.
type MessageCreateParams struct {
	Model         string         `json:"model"`
	Messages      []MessageParam `json:"messages"`
	MaxTokens     int            `json:"max_tokens"`
	Metadata      *Metadata      `json:"metadata,omitempty"`
	StopSequences []string       `json:"stop_sequences,omitempty"`
	Stream        bool           `json:"stream"`
	System        *string        `json:"system,omitempty"`
	Temperature   *float64       `json:"temperature,omitempty"`
	TopK          *int           `json:"top_k,omitempty"`
	TopP          *float64       `json:"top_p,omitempty"`
	ToolChoice    *ToolChoice    `json:"tool_choice,omitempty"`
	Tools         []*ToolParam   `json:"tools,omitempty"`
}

// MessageParam represents a message parameter for API requests.
type MessageParam struct {
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

// Metadata represents metadata about the request.
type Metadata struct {
	UserID *string `json:"user_id,omitempty"`
}

// ToolChoice represents how the model should use the provided tools.
type ToolChoice struct {
	Type string  `json:"type"`
	Name *string `json:"name,omitempty"`
}

// ToolParam defines a tool that the model may use.
type ToolParam struct {
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	InputSchema *jsonschema.Schema `json:"input_schema"`
}

// StreamingMessageResponse represents the structure of streaming message responses.
type StreamingMessageResponse struct {
	MessageStart   *StreamingMessageStart             `json:"message_start,omitempty"`
	MessageDelta   *StreamingMessageDelta             `json:"message_delta,omitempty"`
	ContentBlock   *StreamingMessageContentBlockDelta `json:"content_block_delta,omitempty"`
	MessageStop    *StreamingMessageStop              `json:"message_stop,omitempty"`
	StreamingError *StreamingMessageError             `json:"streaming_error,omitempty"`
}

// StreamingMessageStart represents the start of a streaming message.
type StreamingMessageStart struct {
	Type    string  `json:"type"`
	Message Message `json:"message"`
}

// StreamingMessageDelta represents a delta in a streaming message.
type StreamingMessageDelta struct {
	Type  string `json:"type"`
	Delta struct {
		StopReason   string `json:"stop_reason"`
		StopSequence string `json:"stop_sequence"`
	} `json:"delta"`
	Usage struct {
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// StreamingMessageContentBlockDelta represents a delta in a content block.
type StreamingMessageContentBlockDelta struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
	Delta struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta"`
}

// StreamingMessageStop represents the stop event in a streaming message.
type StreamingMessageStop struct {
	Type  string `json:"type"`
	Delta struct {
		StopReason   string `json:"stop_reason"`
		EndTurn      bool   `json:"end_turn"`
		StopSequence string `json:"stop_sequence"`
	} `json:"delta"`
	Usage struct {
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// StreamingMessageError represents an error in a streaming message.
type StreamingMessageError struct {
	Type  string `json:"type"`
	Error Error  `json:"error"`
}

// Error represents an error message.
type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
