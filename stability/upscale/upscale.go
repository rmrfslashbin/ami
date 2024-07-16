package upscale

// file: stability/upscale/client.go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/rmrfslashbin/ami/stability/generate"
)

func (c *Client) ConservativeUpscale(params *UpscaleParams) (*generate.GenerateResponse, error) {
	return c.upscale("/v2beta/stable-image/upscale/conservative", params)
}

func (c *Client) CreativeUpscale(params *UpscaleParams) (*UpscaleResponse, error) {
	return c.upscaleAsync("/v2beta/stable-image/upscale/creative", params)
}

func (c *Client) GetCreativeUpscaleResult(generationID string) (*UpscaleResult, error) {
	url := fmt.Sprintf("%s/v2beta/stable-image/upscale/creative/result/%s", c.BaseURL, generationID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, response body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result UpscaleResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

func (c *Client) upscale(endpoint string, params *UpscaleParams) (*generate.GenerateResponse, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

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
	addFieldIfNotEmpty(writer, "prompt", params.Prompt)
	addFieldIfNotEmpty(writer, "negative_prompt", params.NegativePrompt)
	addFieldIfNotZero(writer, "width", params.Width)
	addFieldIfNotZero(writer, "height", params.Height)
	addFieldIfNotZero(writer, "seed", params.Seed)
	addFieldIfNotEmpty(writer, "output_format", params.OutputFormat)
	addFieldIfNotZero(writer, "creativity", params.Creativity)

	err := writer.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, response body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result generate.GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

func (c *Client) upscaleAsync(endpoint string, params *UpscaleParams) (*UpscaleResponse, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

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
	addFieldIfNotEmpty(writer, "prompt", params.Prompt)
	addFieldIfNotEmpty(writer, "negative_prompt", params.NegativePrompt)
	addFieldIfNotZero(writer, "width", params.Width)
	addFieldIfNotZero(writer, "height", params.Height)
	addFieldIfNotZero(writer, "seed", params.Seed)
	addFieldIfNotEmpty(writer, "output_format", params.OutputFormat)
	addFieldIfNotZero(writer, "creativity", params.Creativity)

	err := writer.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, response body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result UpscaleResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

func addFieldIfNotEmpty(writer *multipart.Writer, fieldName, fieldValue string) {
	if fieldValue != "" {
		writer.WriteField(fieldName, fieldValue)
	}
}

func addFieldIfNotZero(writer *multipart.Writer, fieldName string, fieldValue interface{}) {
	switch v := fieldValue.(type) {
	case int:
		if v != 0 {
			writer.WriteField(fieldName, fmt.Sprintf("%d", v))
		}
	case float64:
		if v != 0 {
			writer.WriteField(fieldName, fmt.Sprintf("%f", v))
		}
	}
}
