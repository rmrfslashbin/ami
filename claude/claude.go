package claude

// file: claude/claude.go

import (
	"bytes"
	"io"
	"net/http"

	"github.com/rmrfslashbin/ami"
)

const MODULE_NAME = "claude"
const ANTHROPIC_VERSION = ami.Anthropic_Version
const URL = "https://api.anthropic.com"

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
	"sonnet35": {
		Name:            "claude-3-5-sonnet-20240620",
		MaxOutputTokens: 4096,
	},
}

type Option func(config *Claude)

// Configuration structure.
type Claude struct {
	apikey  *string
	headers map[string]string
}

func New(opts ...func(*Claude)) (*Claude, error) {
	config := &Claude{}

	// init headers
	config.headers = make(map[string]string)

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}

	if config.apikey == nil {
		return nil, &ErrMissingAPIKey{}
	}

	config.headers["x-api-key"] = *config.apikey
	config.headers["content-type"] = "application/json"
	config.headers["anthropic-version"] = ANTHROPIC_VERSION
	return config, nil
}

func WithAPIKey(apikey string) Option {
	return func(config *Claude) {
		config.apikey = &apikey
	}
}

func (c *Claude) GetModelMaxOutputTokens(modelName string) int {
	if _, ok := ModelsList[modelName]; !ok {
		return -1
	} else {
		return ModelsList[modelName].MaxOutputTokens
	}

}

func (c *Claude) GetHeaders() map[string]string {
	return c.headers
}

func (c *Claude) Do(url string, jsonData []byte) (*[]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Create an HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &ErrHTTP{
			StatusCode: resp.StatusCode,
			URL:        url,
			Data:       &jsonData,
			Body:       &responseBody,
		}
	}

	return &responseBody, nil
}
