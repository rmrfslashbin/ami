package upscale

// file: stability/upscale/client.go

import (
	"github.com/rmrfslashbin/ami/stability/generate"
)

type Client struct {
	*generate.Client
}

func New(options ...generate.ClientOption) *Client {
	return &Client{
		Client: generate.New(options...),
	}
}

// WithAPIKey returns a ClientOption that sets the API key for the client
func WithAPIKey(apiKey string) generate.ClientOption {
	return generate.WithAPIKey(apiKey)
}

// WithBaseURL returns a ClientOption that sets the base URL for the client
func WithBaseURL(baseURL string) generate.ClientOption {
	return generate.WithBaseURL(baseURL)
}
