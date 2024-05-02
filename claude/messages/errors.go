package messages

import "fmt"

type ErrMissingClaude struct {
	Err error
	Msg string
}

func (e *ErrMissingClaude) Error() string {
	if e.Msg != "" {
		e.Msg = "missing Claude- use WithClaude to set it"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrMaxTokensExceeded struct {
	Err       error
	Msg       string
	Model     string
	MaxTokens int
}

func (e *ErrMaxTokensExceeded) Error() string {
	if e.Msg != "" {
		e.Msg = "max tokens exceeded"
	}
	if e.Model != "" {
		e.Msg += " for model " + e.Model
	}
	if e.MaxTokens != 0 {
		e.Msg += fmt.Sprintf(" (max tokens: %d)", e.MaxTokens)
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrMarshalingInput struct {
	Err error
	Msg string
}

func (e *ErrMarshalingInput) Error() string {
	if e.Msg != "" {
		e.Msg = "marshaling input"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrMarshalingReply struct {
	Err error
	Msg string
}

func (e *ErrMarshalingReply) Error() string {
	if e.Msg != "" {
		e.Msg = "marshaling reply"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrConflictingOptions struct {
	Err error
	Msg string
}

func (e *ErrConflictingOptions) Error() string {
	if e.Msg != "" {
		e.Msg = "conflicting options"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrFetchingMimeType struct {
	Err error
	Msg string
}

func (e *ErrFetchingMimeType) Error() string {
	if e.Msg != "" {
		e.Msg = "fetching mime type"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrUnsupportedMimeType struct {
	Err      error
	Msg      string
	MimeType string
}

func (e *ErrUnsupportedMimeType) Error() string {
	if e.Msg != "" {
		e.Msg = "unsupported mime type"
	}
	if e.MimeType != "" {
		e.Msg += " " + e.MimeType
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrReadingFile struct {
	Err error
	Msg string
}

func (e *ErrReadingFile) Error() string {
	if e.Msg != "" {
		e.Msg = "reading file"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrUnsupportedOption struct {
	Err error
	Msg string
}

func (e *ErrUnsupportedOption) Error() string {
	if e.Msg != "" {
		e.Msg = "unsupported option"
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
	if e.Msg != "" {
		e.Msg = "invalid model"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrMissingModel struct {
	Err error
	Msg string
}

func (e *ErrMissingModel) Error() string {
	if e.Msg != "" {
		e.Msg = "missing model"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrOpeningFile struct {
	Err error
	Msg string
}

func (e *ErrOpeningFile) Error() string {
	if e.Msg != "" {
		e.Msg = "error opening file"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrLoadingGOB struct {
	Err error
	Msg string
}

func (e *ErrLoadingGOB) Error() string {
	if e.Msg != "" {
		e.Msg = "error loading GOB"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrSavingGOB struct {
	Err error
	Msg string
}

func (e *ErrSavingGOB) Error() string {
	if e.Msg != "" {
		e.Msg = "error saving GOB"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrStreamingMessage struct {
	Err error
	Msg string
}

func (e *ErrStreamingMessage) Error() string {
	if e.Msg != "" {
		e.Msg = "error streaming message"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}
