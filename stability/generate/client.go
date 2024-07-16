package generate

// file: stability/generate/client.go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const (
	BaseURL = "https://api.stability.ai"
)

type Client struct {
	APIKey     string
	HTTPClient *http.Client
	BaseURL    string
}

type ClientOption func(*Client)

func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) {
		c.APIKey = apiKey
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

func New(options ...ClientOption) *Client {
	client := &Client{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		BaseURL:    BaseURL,
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func (c *Client) sendRequest(method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.HTTPClient.Do(req)
}

func (c *Client) generateImage(endpoint string, params GenerateParams) (*GenerateResponse, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add image file if it exists (for image-to-image mode)
	if params.Image != "" {
		file, err := os.Open(params.Image)
		if err != nil {
			return nil, fmt.Errorf("error opening image file: %w", err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile("image", params.Image)
		if err != nil {
			return nil, fmt.Errorf("error creating form file: %w", err)
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return nil, fmt.Errorf("error copying file to form: %w", err)
		}
	}

	// Add other fields
	for key, value := range params.toFormData() {
		_ = writer.WriteField(key, value)
	}

	err := writer.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing multipart writer: %w", err)
	}

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
		"Accept":       "application/json",
	}

	resp, err := c.sendRequest("POST", url, body, headers)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, response body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

func (c *Client) GenerateUltra(params GenerateParams) (*GenerateResponse, error) {
	return c.generateImage("/v2beta/stable-image/generate/ultra", params)
}

func (c *Client) GenerateCore(params GenerateParams) (*GenerateResponse, error) {
	return c.generateImage("/v2beta/stable-image/generate/core", params)
}

func (c *Client) GenerateSD3(params GenerateParams) (*GenerateResponse, error) {
	return c.generateImage("/v2beta/stable-image/generate/sd3", params)
}
