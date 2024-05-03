package generate

type Request struct {
	Prompt         string `json:"prompt"`
	AspectRatio    string `json:"aspect_ratio"`
	Mode           string `json:"mode"`
	NegativePrompt string `json:"negative_prompt"`
	Model          string `json:"model"`
	Seed           int    `json:"seed"`
	OutputFormat   string `json:"output_format"`
}

type Response struct {
	Metadata *ResponseMetadata `json:"data,omitempty"`
	Errors   *ResponseErrors   `json:"errors,omitempty"`
	Image    *[]byte           `json:"image,omitempty"`
}

type ResponseMetadata struct {
	ContextType  string `json:"context_type"`
	FinishReason string `json:"finish_reason"`
	Seed         int    `json:"seed"`
}
type ResponseErrors struct {
	Name   string   `json:"name"`
	Errors []string `json:"errors"`
}
