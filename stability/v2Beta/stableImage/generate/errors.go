package generate

import (
	"fmt"
	"strings"
)

type ErrMissingStability struct {
	Err error
	Msg string
}

func (e *ErrMissingStability) Error() string {
	if e.Msg != "" {
		e.Msg = "missing stability- use WithStability to set it"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrMissingLogger struct {
	Err error
	Msg string
}

func (e *ErrMissingLogger) Error() string {
	if e.Msg != "" {
		e.Msg = "missing logger- use WithLogger to set it"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrMissingPrompt struct {
	Err error
	Msg string
}

func (e *ErrMissingPrompt) Error() string {
	if e.Msg != "" {
		e.Msg = "missing prompt- use WithPrompt or SetPrompt to set it"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrInvalidAspectRatio struct {
	Err error
	Msg string
}

func (e *ErrInvalidAspectRatio) Error() string {
	validAspectRatios := strings.Join(ASPECT_RATIOS, ", ")
	if e.Msg != "" {
		e.Msg = "invalid aspect ratio- use WithAspectRatio or SetAspectRatio to set it. Must be one of " + validAspectRatios
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrInvalidModel struct {
	Err error
	Msg string
}

func (e *ErrInvalidModel) Error() string {
	validModels := strings.Join(MODELS, ", ")
	if e.Msg != "" {
		e.Msg = "invalid model- use WithModel or SetModel to set it. Must be one of " + validModels
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrInvalidSeed struct {
	Err error
	Msg string
}

func (e *ErrInvalidSeed) Error() string {
	if e.Msg != "" {
		e.Msg = fmt.Sprintf("invalid seed- use WithSeed or SetSeed to set it. Must be between 0 and %d", MAX_SEED)
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrInvalidOutputFormat struct {
	Err error
	Msg string
}

func (e *ErrInvalidOutputFormat) Error() string {
	validOutputFormats := strings.Join(OUTPUT_FORMATS, ", ")
	if e.Msg != "" {
		e.Msg = "invalid output format- use WithOutputFormat or SetOutputFormat to set it. Must be one of " + validOutputFormats
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrEmptyResponse struct {
	Err error
	Msg string
}

func (e *ErrEmptyResponse) Error() string {
	if e.Msg != "" {
		e.Msg = "empty response from Stability.Ai"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrJsonResponseNotSupported struct {
	Err error
	Msg string
}

func (e *ErrJsonResponseNotSupported) Error() string {
	if e.Msg != "" {
		e.Msg = "json response not supported"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrFetchingReturnHeader struct {
	Err    error
	Msg    string
	Header string
}

func (e *ErrFetchingReturnHeader) Error() string {
	if e.Msg != "" {
		e.Msg = "error fetching return header"
	}
	if e.Header != "" {
		e.Msg = e.Msg + ": " + e.Header
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrUnableToParseResponse struct {
	Err      error
	Msg      string
	Response []byte
}

func (e *ErrUnableToParseResponse) Error() string {
	if e.Msg != "" {
		e.Msg = "unable to parse response"
	}
	if len(e.Response) > 0 {
		e.Msg = e.Msg + ": " + string(e.Response)
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrInvalidPromptLength struct {
	Err error
	Msg string
}

func (e *ErrInvalidPromptLength) Error() string {
	if e.Msg != "" {
		e.Msg = fmt.Sprintf("invalid prompt- use WithPrompt or SetPrompt to set it. Must be between 1 and %d characters", MAX_PROMPT_LENGTH)
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrInvalidNegativePromptLength struct {
	Err error
	Msg string
}

func (e *ErrInvalidNegativePromptLength) Error() string {
	if e.Msg != "" {
		e.Msg = fmt.Sprintf("invalid negative prompt- use WithNegativePrompt or SetNegativePrompt to set it. If set, the prompt must be between 1 and %d characters", MAX_PROMPT_LENGTH)
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}
