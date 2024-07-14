package messages

// file: claude/messages/streaming.go

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/tmaxmax/go-sse"
)

type StreamEvent interface {
	Type() string
}

type MessageStartEvent struct {
	Message Message `json:"message"`
}

func (e MessageStartEvent) Type() string {
	return "message_start"
}

type ContentBlockStartEvent struct {
	Index        int          `json:"index"`
	ContentBlock ContentBlock `json:"content_block"`
}

func (e ContentBlockStartEvent) Type() string {
	return "content_block_start"
}

type StreamResults struct {
	Response <-chan StreamEvent
	Error    <-chan error
}

func (m *Messages) Stream(ctx context.Context) StreamResults {
	eventCh := make(chan StreamEvent)
	errCh := make(chan error)

	if len(m.request.Tools) > 0 {
		errCh <- &ErrToolUseNotSupported{}
		close(eventCh)
		return StreamResults{Response: eventCh, Error: errCh}
	}

	if m.request.TopP != nil && m.request.Temperature != nil {
		errCh <- &ErrConflictingOptions{Err: errors.New("top_p and temperature cannot be used together")}
		close(eventCh)
		return StreamResults{Response: eventCh, Error: errCh}
	}

	// Validate all messages in the conversation
	for _, msg := range m.conversation.Messages {
		if err := m.validateMessage(msg); err != nil {
			errCh <- fmt.Errorf("invalid message in conversation: %w", err)
			close(eventCh)
			return StreamResults{Response: eventCh, Error: errCh}
		}
	}

	// Convert conversation messages to MessageParams
	m.request.Messages = convertToMessageParams(m.conversation.Messages)

	jsonData, err := json.Marshal(m.request)
	if err != nil {
		errCh <- &ErrMarshalingInput{Err: err}
		close(eventCh)
		return StreamResults{Response: eventCh, Error: errCh}
	}

	go func() {
		defer close(eventCh)
		defer close(errCh)

		// Create a new context that we can cancel
		streamCtx, cancel := context.WithCancel(ctx)
		defer cancel() // Ensure all resources are cleaned up

		req, err := http.NewRequestWithContext(streamCtx, "POST", m.url, bytes.NewBuffer(jsonData))
		if err != nil {
			errCh <- err
			return
		}

		// Set headers
		for key, value := range m.claud.GetHeaders() {
			req.Header.Set(key, value)
		}

		conn := sse.DefaultClient.NewConnection(req)

		// Handle context cancellation
		go func() {
			<-ctx.Done()
			cancel() // This will cause the SSE connection to stop
		}()

		conn.SubscribeEvent("message_start", func(event sse.Event) {
			var response MessageStartEvent
			err := json.Unmarshal([]byte(event.Data), &response)
			if err != nil {
				errCh <- &ErrMarshalingReply{Err: err}
				return
			}
			eventCh <- response
		})

		conn.SubscribeEvent("content_block_start", func(event sse.Event) {
			var response ContentBlockStartEvent
			err := json.Unmarshal([]byte(event.Data), &response)
			if err != nil {
				errCh <- &ErrMarshalingReply{Err: err}
				return
			}
			eventCh <- response
		})

		conn.SubscribeEvent("content_block_delta", func(event sse.Event) {
			var response ContentBlockDeltaEvent
			err := json.Unmarshal([]byte(event.Data), &response)
			if err != nil {
				errCh <- &ErrMarshalingReply{Err: err}
				return
			}
			eventCh <- response
		})

		conn.SubscribeEvent("message_delta", func(event sse.Event) {
			var response MessageDeltaEvent
			err := json.Unmarshal([]byte(event.Data), &response)
			if err != nil {
				errCh <- &ErrMarshalingReply{Err: err}
				return
			}
			eventCh <- response
		})

		conn.SubscribeEvent("message_stop", func(event sse.Event) {
			var response MessageStopEvent
			err := json.Unmarshal([]byte(event.Data), &response)
			if err != nil {
				errCh <- &ErrMarshalingReply{Err: err}
				return
			}
			eventCh <- response
		})

		conn.SubscribeEvent("error", func(event sse.Event) {
			var response StreamingErrorEvent
			err := json.Unmarshal([]byte(event.Data), &response)
			if err != nil {
				errCh <- &ErrMarshalingReply{Err: err}
				return
			}
			eventCh <- response
			errCh <- &ErrStreamingMessage{Err: errors.New(response.Error.Message)}
		})

		conn.SubscribeEvent("ping", func(event sse.Event) {
			// Handle ping event if needed
		})

		if err := conn.Connect(); err != nil {
			if streamCtx.Err() != nil {
				// The context was cancelled, so this isn't an unexpected error
				return
			}
			errCh <- err
			return
		}
	}()

	return StreamResults{Response: eventCh, Error: errCh}
}

type ContentBlockDeltaEvent struct {
	EventType string `json:"type"` // Changed from Type to EventType
	Index     int    `json:"index"`
	Delta     struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta"`
}

func (e ContentBlockDeltaEvent) Type() string {
	return e.EventType // Return the EventType field
}

type MessageDeltaEvent struct {
	EventType string `json:"type"` // Changed from Type to EventType
	Delta     struct {
		StopReason   string `json:"stop_reason"`
		StopSequence string `json:"stop_sequence"`
	} `json:"delta"`
	Usage struct {
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

func (e MessageDeltaEvent) Type() string {
	return e.EventType // Return the EventType field
}

type MessageStopEvent struct {
	EventType string `json:"type"` // Changed from Type to EventType
}

func (e MessageStopEvent) Type() string {
	return e.EventType // Return the EventType field
}

type StreamingErrorEvent struct {
	EventType string `json:"type"` // Changed from Type to EventType
	Error     struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e StreamingErrorEvent) Type() string {
	return e.EventType // Return the EventType field
}
