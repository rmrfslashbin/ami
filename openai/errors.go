package openai

// file: openai/errors.go

import "fmt"

type ErrMissingClient struct {
	Err error
	Msg string
}

func (e *ErrMissingClient) Error() string {
	if e.Msg == "" {
		e.Msg = "missing OpenAI client - use WithAPIKey to set it"
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	}
	return e.Msg
}

// Add more error types as needed
