package claude

// file: claude/errors.go

import "fmt"

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

/*
	https://docs.anthropic.com/claude/reference/errors

400 - invalid_request_error: There was an issue with the format or content of your request.
401 - authentication_error: There's an issue with your API key.
403 - permission_error: Your API key does not have permission to use the specified resource.
404 - not_found_error: The requested resource was not found.
429 - rate_limit_error: Your account has hit a rate limit.
500 - api_error: An unexpected error has occurred internal to Anthropic's systems.
529 - overloaded_error: Anthropic's API is temporarily overloaded.
*/
var apiErrors = map[int]string{
	400: "invalid_request_error: There was an issue with the format or content of your request",
	401: "authentication_error: There's an issue with your API key",
	403: "permission_error: Your API key does not have permission to use the specified resource",
	404: "not_found_error: The requested resource was not found",
	413: "request_too_large_error: The request exceeds the maximum allowed number of bytes",
	429: "rate_limit_error: Your account has hit a rate limit",
	500: "api_error: An unexpected error has occurred internal to Anthropic's systems",
	529: "overloaded_error: Anthropic's API is temporarily overloaded",
}

type ErrHTTP struct {
	Err        error
	Msg        string
	URL        string
	Data       *[]byte
	Body       *[]byte
	StatusCode int
}

func (e *ErrHTTP) Error() string {
	if e.Msg != "" {
		e.Msg = "HTTP error"
	}
	if e.StatusCode != 0 {
		e.Msg += fmt.Sprintf(": %d", e.StatusCode)
	}
	if errorAPIText, ok := apiErrors[e.StatusCode]; ok {
		e.Msg += fmt.Sprintf(": %s", errorAPIText)
	}
	if e.URL != "" {
		e.Msg += fmt.Sprintf(" for %s", e.URL)
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrInvalidRequest struct {
	Err error
	Msg string
}

func (e *ErrInvalidRequest) Error() string {
	if e.Msg != "" {
		e.Msg = "invalid request"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrAuthentication struct {
	Err error
	Msg string
}

func (e *ErrAuthentication) Error() string {
	if e.Msg != "" {
		e.Msg = "authentication error"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}
