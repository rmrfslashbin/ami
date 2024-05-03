package stability

import "fmt"

// path: stability/stability.go

type ErrMissingAPIKey struct {
	Err error
	Msg string
}

func (e *ErrMissingAPIKey) Error() string {
	if e.Msg != "" {
		e.Msg = "missing API key- use WithAPIKey to set it"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrHTTP struct {
	Err        error
	Msg        string
	Url        string
	StatusCode int
}

func (e *ErrHTTP) Error() string {
	if e.Msg != "" {
		e.Msg = "HTTP error"
	}
	if e.StatusCode != 0 {
		e.Msg += fmt.Sprintf(": %d", e.StatusCode)
	}
	if errorAPIText, ok := STATUS_CODES[e.StatusCode]; ok {
		e.Msg += fmt.Sprintf(": %s", errorAPIText)
	}
	if e.Url != "" {
		e.Msg += fmt.Sprintf(" for %s", e.Url)
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}
