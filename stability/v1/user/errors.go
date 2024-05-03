package user

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
