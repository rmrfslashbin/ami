package stability

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
)

// Path: stability/stabilityV2.go

type Option func(config *Stability)

// Configuration structure.
type Stability struct {
	apikey    *string
	log       *slog.Logger
	headers   map[string]string
	formParts []map[string]interface{}
}

func New(opts ...func(*Stability)) (*Stability, error) {
	config := &Stability{}

	// init headers
	config.headers = make(map[string]string)

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}

	if config.apikey == nil {
		return nil, &ErrMissingAPIKey{}
	}

	config.headers["authorization"] = "Bearer " + *config.apikey
	//config.headers["content-type"] = "multipart/form-data"
	config.headers["accept"] = "image/png" // default to png
	return config, nil
}

func WithAPIKey(apikey string) Option {
	return func(config *Stability) {
		config.apikey = &apikey
	}
}

func WithLogger(log *slog.Logger) Option {
	return func(config *Stability) {
		moduleLogger := log.With(
			slog.Group("module_info",
				slog.String("module", MODULE_NAME),
			),
		)
		config.log = moduleLogger
	}
}

func (stability *Stability) AddFormPart(key string, value interface{}) {
	stability.formParts = append(stability.formParts, map[string]interface{}{
		key: value,
	})
}

func (stability *Stability) AddHeader(key string, value string) {
	stability.headers[key] = value
}

func (stability *Stability) MakeFormBody() (*bytes.Buffer, *string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, part := range stability.formParts {
		for key, value := range part {
			writer.WriteField(key, fmt.Sprint(value))
		}
	}
	contentType := writer.FormDataContentType()
	writer.Close()
	return body, &contentType
}

func (stability *Stability) Do(url *string, httpMethod HttpMethod) (*StabilityResponse, error) {
	method := string(httpMethod)
	body, contentType := stability.MakeFormBody()

	req, err := http.NewRequest(method, *url, body)
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range stability.headers {
		req.Header.Set(key, value)
	}

	req.Header.Set("Content-Type", *contentType)

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
		return &StabilityResponse{
			Errors: &responseBody,
		}, nil
	}

	responseHeaders := map[string][]string{}

	for key, value := range resp.Header {
		responseHeaders[key] = value
	}

	if err != nil {
		return nil, err
	}

	return &StabilityResponse{
		Body:    responseBody,
		Headers: responseHeaders,
	}, nil
}
