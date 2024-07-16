package upscale

// file: stability/upscale/models.go

type UpscaleParams struct {
	Image          string  `json:"-"` // Not sent in JSON, used for file upload
	Prompt         string  `json:"prompt"`
	NegativePrompt string  `json:"negative_prompt,omitempty"`
	Width          int     `json:"width,omitempty"`
	Height         int     `json:"height,omitempty"`
	Seed           int     `json:"seed,omitempty"`
	OutputFormat   string  `json:"output_format,omitempty"`
	Creativity     float64 `json:"creativity,omitempty"`
}

type UpscaleOption func(*UpscaleParams)

func WithImage(image string) UpscaleOption {
	return func(p *UpscaleParams) {
		p.Image = image
	}
}

func WithPrompt(prompt string) UpscaleOption {
	return func(p *UpscaleParams) {
		p.Prompt = prompt
	}
}

// Add more WithX functions for other fields...

func NewUpscaleParams(options ...UpscaleOption) *UpscaleParams {
	params := &UpscaleParams{}
	for _, option := range options {
		option(params)
	}
	return params
}

type UpscaleResponse struct {
	ID string `json:"id"`
}

type UpscaleResult struct {
	Image        string `json:"image,omitempty"`
	FinishReason string `json:"finish_reason"`
	Seed         int    `json:"seed"`
}
