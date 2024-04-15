package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type MessageContent struct {
	Type   string      `json:"type"`
	Text   string      `json:"text,omitempty"`
	Source *DataSource `json:"source,omitempty"`
}

type DataSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      []byte `json:"data"`
}

type Request struct {
	Model     string                 `json:"model"`
	Messages  []Message              `json:"messages"`
	MaxTokens int                    `json:"max_tokens"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Stop      []string               `json:"stop_sequences,omitempty"`
	Stream    bool                   `json:"stream,omitempty"`
	System    string                 `json:"system,omitempty"`
	Temp      float64                `json:"temperature,omitempty"`
	Tools     []Tool                 `json:"tools,omitempty"`
	TopK      int                    `json:"top_k,omitempty"`
	TopP      float64                `json:"top_p,omitempty"`
}

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"input_schema,omitempty"`
}

type Response struct {
	ID         string        `json:"id"`
	Content    []interface{} `json:"content"`
	Model      string        `json:"model"`
	StopReason string        `json:"stop_reason"`
	StopSeq    string        `json:"stop_sequence"`
	Usage      Usage         `json:"usage"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

func main() {
	// Set up the request payload
	request := Request{
		Model: "your_model_name",
		Messages: []Message{
			{
				Role:    "user",
				Content: "Hello, how are you?",
			},
		},
		MaxTokens: 100,
	}

	// Convert the request payload to JSON
	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshaling request:", err)
		return
	}

	// Set up the HTTP request
	url := "https://api.anthropic.com/v1/messages"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// Add any required authentication headers

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Parse the response JSON
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshaling response:", err)
		return
	}

	// Process the response data
	fmt.Println("Response ID:", response.ID)
	fmt.Println("Content:", response.Content)
	// ...
}
