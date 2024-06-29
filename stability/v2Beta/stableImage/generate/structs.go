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
	Errors       *ResponseErrors `json:"errors,omitempty"`
	Image        *string         `json:"-"`
	FinishReason *string         `json:"finish_reason,omitempty"`
	Seed         *int            `json:"seed,omitempty"`
	Filename     *string         `json:"filename,omitempty"`
}

type ResponseErrors struct {
	Errors []string `json:"errors,omitempty"`
	Id     string   `json:"id,omitempty"`
	Name   string   `json:"name,omitempty"`
}
