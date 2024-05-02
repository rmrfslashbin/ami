package stabilityV3

import (
	"fmt"
	"strings"

	"github.com/rmrfslashbin/ami/stability"
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
	validAspectRatios := strings.Join(stability.ASPECT_RATIOS, ", ")
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
	validModels := strings.Join(stability.MODELS_V3, ", ")
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
		e.Msg = fmt.Sprintf("invalid seed- use WithSeed or SetSeed to set it. Must be between 0 and %d", stability.MAX_SEED)
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
	validOutputFormats := strings.Join(stability.OUTPUT_FORMATS_V3, ", ")
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
