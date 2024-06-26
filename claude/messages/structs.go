package messages

import (
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/invopop/jsonschema"
)

// Conversation represents a conversation. This is not part of the API.
type Conversation struct {
	// Id is the unique object identifier.
	Id string `json:"id"`

	// Model is the model used in the conversation.
	Model *string `json:"model"`

	// Created is the time the conversation was created.
	Created time.Time `json:"created"`

	// Updated is the time the conversation was updated.
	Updated time.Time `json:"updated"`

	// Messages is a list of messages in the conversation.
	Messages []*Message `json:"messages"`
}

// https://docs.anthropic.com/claude/reference/messages_post
// Request is the request to send to the Messages API.
type Request struct {
	// Model is the model that will complete your prompt.
	// Required.
	// See models (https://docs.anthropic.com/claude/docs/models-overview) for additional details and options.
	Model string `json:"model" required:"true" enum:"haiku,sonnet,opus"`

	// Messages is the messages to send to the API.
	// Required.
	Messages []*Message `json:"messages" required:"true"`

	// MaxToken is the maximum number of tokens to generate before stopping.
	// Required.
	MaxTokens int `json:"max_tokens" required:"true"`

	// Metadata is an object describing metadata about the request.
	Metadata *Metadata `json:"metadata,omitempty"`

	// StopSequences is a list of strings that, if generated, will cause the model to stop generating tokens.
	StopSequences []string `json:"stop_sequences,omitempty"`

	// Stream is a boolean that indicates whether the model should generate a single response or a stream of responses.
	// Default is false.
	Stream bool `json:"stream"`

	// System is a system prompt is a way of providing context and instructions to Claude, such as specifying a particular goal or role.
	System string `json:"system,omitempty"`

	// Temperature is a float that controls the randomness of the model's output. The higher the temperature, the more random the output.
	Temperature *float32 `json:"temperature,omitempty"`

	// ToolChoice is a tool choice.
	ToolChoice *ToolChoice `json:"tool_choice,omitempty" enum:"auto,any,tool"`

	// Tools are definitions of tools that the model may use
	Tools []*Tool `json:"tools,omitempty"`

	// TopK is an integer that specifies sampling from the top K options for each subsequent token.
	// Recommended for advanced use cases only. You usually only need to use temperature.
	TopK *int `json:"top_k,omitempty"`

	// TopP is an integer that controls nucleus sampling. The higher the top_p, the more diverse the output.
	// Recommended for advanced use cases only. You usually only need to use temperature.
	// You should either alter temperature or top_p, but not both.
	TopP *int `json:"top_p,omitempty"`

	// modelMaxTokens is the maximum number of tokens for the model.
	modelMaxTokens int
}

// Metadata is an object describing metadata about the request.
type Metadata struct {
	// UserId is an external identifier for the user who is associated with the request.
	// This should be a uuid, hash value, or other opaque identifier.
	// Anthropic may use this id to help detect abuse.
	// Do not include any identifying information such as name, email address, or phone number.
	UserId string `json:"user_id,omitempty"`
}

// Message is a sinble nput message.
type Message struct {
	// Role is the conversational role of the message.
	// Specify a single user-role message, or you can include multiple "user" and "assistant" messages.
	// The first message must always use the "user" role.
	Role string `json:"role"`

	// MessageContent is the content of the message.
	MessageContent []*Content `json:"content"`
}

// Content is the content of the message.
type Content struct {
	// Type is the type of content.
	Type string `json:"type"`

	// Text is the text of the content.
	Text string `json:"text,omitempty"`

	// Content is the content of the content.
	Content string `json:"content,omitempty"`

	// Id is the unique object identifier for a tool_use block.
	Id string `json:"id,omitempty"`

	// ToolUseId is the unique object identifier for a tool_use block.
	ToolUseId string `json:"tool_use_id,omitempty"`

	// Input is the input request for a tool_use block.
	Input interface{} `json:"input,omitempty"`

	// Name is the name of the tool used in a tool_use block.
	Name string `json:"name,omitempty"`

	// Source is the source of the media.
	Source *MediaSource `json:"source,omitempty"`
}

// MediaSource is the source of the media.
type MediaSource struct {
	// Type is the type of media source.
	// The only type is "base64".
	Type string `json:"type"`

	// MediaType is the media type of the data.
	// Valid image types: image/jpeg, image/png, image/gif, and image/webp.
	MediaType string `json:"media_type"`

	// Data is the base64 encoded data.
	Data string `json:"data"`
}

// ToolChoice is a tool choice. The model can use a specific tool, any available tool, or decide by itself.
type ToolChoice struct {
	Type string `json:"type" enum:"auto,any,tool" required:"true"`
	Name string `json:"name" require_if:"type=tool"`
}

// https://docs.anthropic.com/claude/reference/messages_post
// Response is the response from the Messages API.
type Response struct {
	// Id is the unique object identifier.
	// Required.
	// The format and length of IDs may change over time.
	Id string `json:"id"`

	// Object type.
	// Required.
	// For Messages, this is always "message".
	Type string `json:"type"`

	// Error is the error message.
	Error Error `json:"error,omitempty"`

	// Conversational role of the generated message.
	// Required.
	// This will always be "assistant".
	Role string `json:"role"`

	// Content generated by the model.
	// Required.
	// This is an array of content blocks, each of which has a type that determines its shape. Currently, the only type in responses is "text".
	Content []*Content `json:"content"`

	// The model that handled the request.
	// Required.
	Model string `json:"model"`

	// StopReason is the reason that we stopped.
	// Required.
	StopReason string `json:"stop_reason"`

	// StopSequences indicates which custom stop sequence was generated, if any.
	// Required.
	StopSequences string `json:"stop_sequences"`

	// Usage is the usage of the API billing and rate-limit data.
	// Required.
	Usage Usage `json:"usage"`
}

type Usage struct {
	// InputTokens is the number of tokens used as input to the model.
	// Required.

	InputTokens int `json:"input_tokens"`
	// OutputTokens is the number of tokens generated by the model.
	// Required.
	OutputTokens int `json:"output_tokens"`
}

// Tool defines a tool that the model may use.
type Tool struct {
	// Name is the name of the tool.
	Name string `json:"name"`

	// Description is the optional description of the tool.
	Description string `json:"description"`

	// Input_schema specified the JSON schema for the tool input shape that the model will produce in tool_use output content blocks.
	InputSchema *jsonschema.Schema `json:"input_schema"`
}

// ToolReply is the reply from the tool.
type ToolReply struct {
	/*
			{
		      "type": "tool_use",
		      "id": "toolu_01A09q90qw90lq917835lq9",
		      "name": "get_weather",
		      "input": {"location": "San Francisco, CA", "unit": "celsius"}
		    }
	*/

	// Id is the unique object identifier for a tool_use block.
	Id string `json:"id"`

	// Name is the name of the tool used in a tool_use block.
	Name string `json:"name"`

	// Input is the input request for a tool_use block.
	Input interface{} `json:"input"`

	// Type is the type of content. Should ne "tool_use".
	Type string `json:"type"`
}

type StreamingMessageResponse struct {
	MessageStart   *StreamingMessageStart             `json:"message_start"`
	MessageDelta   *StreamingMessageDelta             `json:"message_delta"`
	ContentBlock   *StreamingMessageContentBlockDelta `json:"content_block_delta"`
	MessageStop    *StreamingMessageStop              `json:"message_stop"`
	StreamingError *StreamingMessageError             `json:"streaming_error"`
}

type StreamingMessageStart struct {
	Type    string   `json:"type"`
	Message Response `json:"message"`
}

type StreamingMessageDelta struct {
	Type  string `json:"type"`
	Delta struct {
		StopReason   string `json:"stop_reason"`
		StopSequence string `json:"stop_sequence"`
	} `json:"delta"`
	Usage struct {
		OutputTokens int `json:"output_tokens"`
	}
}

type StreamingMessageContentBlockDelta struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
	Delta struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta"`
}

type StreamingMessageStop struct {
	Type  string `json:"type"`
	Delta struct {
		StopReason   string `json:"stop_reason"`
		EndTurn      bool   `json:"end_turn"`
		StopSequence string `json:"stop_sequence"`
	} `json:"delta"`
	Usage struct {
		OutputTokens int `json:"output_tokens"`
	}
}

type StreamingMessageError struct {
	Type  string `json:"type"`
	Error Error  `json:"error"`
}

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Validater is an interface for validating structs.
type Validater interface {
	Validate() error
}

// Validate validates the struct.
func Validate(v Validater) error {
	// reflect the struct and parse the tags
	val := reflect.ValueOf(v)
	for i := 0; i < val.NumField(); i++ {
		requiredTag := val.Type().Field(i).Tag.Get("required")
		enumTag := val.Type().Field(i).Tag.Get("enum")
		requireIfTag := val.Type().Field(i).Tag.Get("require_if")

		// check if the field is required
		if requiredTag == "true" {
			if val.Field(i).IsZero() {
				return &ValidationError{Field: val.Type().Field(i).Name, Message: "is required"}
			}
		}

		// check if the field is in the enum; the format of the enum field is "enum1,enum2,enum3"
		if enumTag != "" {
			enum := strings.Split(enumTag, ",")
			if !slices.Contains(enum, val.Field(i).String()) {
				return &ValidationError{Field: val.Type().Field(i).Name, Message: "is not in the enum"}
			}
		}

		// check if the field is required if another field is set
		if requireIfTag != "" {
			// split the require_if tag

			requireIfSplit := strings.Split(requireIfTag, "=")
			requireIfField := requireIfSplit[0]
			requireIfValue := requireIfSplit[1]

			// if a value is set, the field is required
			if val.FieldByName(requireIfField).String() == requireIfValue {
				if val.Field(i).IsZero() {
					return &ValidationError{Field: val.Type().Field(i).Name, Message: "is required"}
				}
			}
		}
	}

	// call the struct's validate method
	return v.Validate()
}

// Validate validates the Request struct.
func (r *Request) Validate() error {
	if r.MaxTokens > r.modelMaxTokens {
		return &ValidationError{Field: "max_tokens", Message: "cannot be greater than the model's max tokens"}
	}

	if r.TopP != nil && r.Temperature != nil {
		return &ValidationError{Field: "top_p and temperature", Message: "cannot both be set"}
	}

	return nil
}
