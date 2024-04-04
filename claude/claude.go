package claude

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
)

const MODULE_NAME = "claude"
const ANTHROPIC_VERSION = "2023-06-01"
const URL = "https://api.anthropic.com"

// headers:
// x-api-key: YOUR_API_KEY"
// content-type: application/json
// anthropic-version: 2023-06-01

/* https://docs.anthropic.com/claude/reference/errors
Our API follows a predictable HTTP error code format:

400 - invalid_request_error: There was an issue with the format or content of your request.
401 - authentication_error: There's an issue with your API key.
403 - permission_error: Your API key does not have permission to use the specified resource.
404 - not_found_error: The requested resource was not found.
429 - rate_limit_error: Your account has hit a rate limit.
500 - api_error: An unexpected error has occurred internal to Anthropic's systems.
529 - overloaded_error: Anthropic's API is temporarily overloaded.
*/

/* https://docs.anthropic.com/claude/reference/errors
{
  "type": "error",
  "error": {
    "type": "not_found_error",
    "message": "The requested resource could not be found."
  }
}
*/

/*
Model	        | API Name
Claude 3 Opus	| claude-3-opus-20240229
Claude 3 Sonnet	| claude-3-sonnet-20240229
Claude 3 Haiku	| claude-3-haiku-20240307
*/

type Model struct {
	Name            string `json:"name"`
	MaxOutputTokens int    `json:"max_output_tokens"`
}

var ModelsList = map[string]*Model{
	"opus": {
		Name:            "claude-3-opus-20240229",
		MaxOutputTokens: 4096,
	},
	"sonnet": {
		Name:            "claude-3-sonnet-20240229",
		MaxOutputTokens: 4096,
	},
	"haiku": {
		Name:            "claude-3-haiku-20240307",
		MaxOutputTokens: 4096,
	},
}

type Option func(config *Claude)

// Configuration structure.
type Claude struct {
	apikey *string
	log    *slog.Logger
}

func New(opts ...func(*Claude)) (*Claude, error) {
	config := &Claude{}

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}

	if config.apikey == nil {
		return nil, &ErrMissingAPIKey{}
	}

	return config, nil
}

func WithAPIKey(apikey string) Option {
	return func(config *Claude) {
		config.apikey = &apikey
	}
}

func WithLogger(log *slog.Logger) Option {
	return func(config *Claude) {
		moduleLogger := log.With(
			slog.Group("module_info",
				slog.String("module", MODULE_NAME),
			),
		)
		config.log = moduleLogger
	}
}

func (c *Claude) GetModelMaxOutputTokens(modelName string) int {
	if _, ok := ModelsList[modelName]; !ok {
		return -1
	} else {
		return ModelsList[modelName].MaxOutputTokens
	}

}

func (c *Claude) Do(url string, jsonData []byte) (*[]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", *c.apikey)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("anthropic-version", ANTHROPIC_VERSION)

	// Create an HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &ErrHTTP{
			StatusCode: resp.StatusCode,
			URL:        url,
			Data:       &jsonData,
		}
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}
