package messages

// file: claude/messages/streaming.go

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

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

// StreamResults is a struct that contains the response and error channels for the streaming API
type StreamResults struct {
	Response <-chan StreamEvent
	Error    <-chan error
}

// errorBuffer is a helper struct to manage multiple errors
type errorBuffer struct {
	errors []error
	mu     sync.Mutex
}

func (eb *errorBuffer) add(err error) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.errors = append(eb.errors, err)
}

func (eb *errorBuffer) get() []error {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	return eb.errors
}

func (m *Messages) Stream(ctx context.Context) StreamResults {
	eventCh := make(chan StreamEvent)
	errCh := make(chan error)
	errBuf := &errorBuffer{}

	if len(m.request.Tools) > 0 {
		errBuf.add(&ErrToolUseNotSupported{})
		close(eventCh)
		go func() {
			for _, err := range errBuf.get() {
				errCh <- err
			}
			close(errCh)
		}()
		return StreamResults{Response: eventCh, Error: errCh}
	}

	if m.request.TopP != nil && m.request.Temperature != nil {
		errBuf.add(&ErrConflictingOptions{Err: errors.New("top_p and temperature cannot be used together")})
		close(eventCh)
		go func() {
			for _, err := range errBuf.get() {
				errCh <- err
			}
			close(errCh)
		}()
		return StreamResults{Response: eventCh, Error: errCh}
	}

	// Validate all messages in the conversation
	for _, msg := range m.conversation.Messages {
		if err := m.validateMessage(msg); err != nil {
			errBuf.add(fmt.Errorf("invalid message in conversation: %w", err))
			close(eventCh)
			go func() {
				for _, err := range errBuf.get() {
					errCh <- err
				}
				close(errCh)
			}()
			return StreamResults{Response: eventCh, Error: errCh}
		}
	}

	// Convert conversation messages to MessageParams
	m.request.Messages = convertToMessageParams(m.conversation.Messages)

	jsonData, err := json.Marshal(m.request)
	if err != nil {
		errBuf.add(&ErrMarshalingInput{Err: err})
		close(eventCh)
		go func() {
			for _, err := range errBuf.get() {
				errCh <- err
			}
			close(errCh)
		}()
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
			errBuf.add(err)
			for _, err := range errBuf.get() {
				errCh <- err
			}
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
				errBuf.add(&ErrMarshalingReply{Err: err})
				return
			}
			eventCh <- response
		})

		conn.SubscribeEvent("content_block_start", func(event sse.Event) {
			var response ContentBlockStartEvent
			err := json.Unmarshal([]byte(event.Data), &response)
			if err != nil {
				errBuf.add(&ErrMarshalingReply{Err: err})
				return
			}
			eventCh <- response
		})

		conn.SubscribeEvent("content_block_delta", func(event sse.Event) {
			var response ContentBlockDeltaEvent
			err := json.Unmarshal([]byte(event.Data), &response)
			if err != nil {
				errBuf.add(&ErrMarshalingReply{Err: err})
				return
			}
			eventCh <- response
		})

		conn.SubscribeEvent("message_delta", func(event sse.Event) {
			var response MessageDeltaEvent
			err := json.Unmarshal([]byte(event.Data), &response)
			if err != nil {
				errBuf.add(&ErrMarshalingReply{Err: err})
				return
			}
			eventCh <- response
		})

		conn.SubscribeEvent("message_stop", func(event sse.Event) {
			var response MessageStopEvent
			err := json.Unmarshal([]byte(event.Data), &response)
			if err != nil {
				errBuf.add(&ErrMarshalingReply{Err: err})
				return
			}
			eventCh <- response
		})

		conn.SubscribeEvent("error", func(event sse.Event) {
			var response StreamingErrorEvent
			err := json.Unmarshal([]byte(event.Data), &response)
			if err != nil {
				errBuf.add(&ErrMarshalingReply{Err: err})
				return
			}
			eventCh <- response
			errBuf.add(&ErrStreamingMessage{Err: errors.New(response.Error.Message)})
		})

		conn.SubscribeEvent("ping", func(event sse.Event) {
			// Handle ping event if needed
		})

		if err := conn.Connect(); err != nil {
			if streamCtx.Err() != nil {
				// The context was cancelled, so this isn't an unexpected error
				return
			}
			errBuf.add(err)
		}

		// Send all buffered errors
		for _, err := range errBuf.get() {
			errCh <- err
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
