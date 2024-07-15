package stability

// file: ami/stability/models.go

import (
	"fmt"
	"strconv"
)

// ValidAspectRatios defines the valid aspect ratios for different endpoints
var ValidAspectRatios = map[string][]string{
	"ultra": {"21:9", "16:9", "3:2", "5:4", "1:1", "4:5", "2:3", "9:16", "9:21"},
	"core":  {"21:9", "16:9", "3:2", "5:4", "1:1", "4:5", "2:3", "9:16", "9:21"},
	"sd3":   {"1:1"}, // SD3 might have different aspect ratio requirements, adjust as needed
}

type GenerateParams struct {
	Model          string  `json:"model,omitempty"`
	Prompt         string  `json:"prompt"`
	NegativePrompt string  `json:"negative_prompt,omitempty"`
	Mode           string  `json:"mode,omitempty"`
	Image          string  `json:"-"` // Not sent in JSON, used for file upload
	Strength       float64 `json:"strength,omitempty"`
	AspectRatio    string  `json:"aspect_ratio,omitempty"`
	Seed           int     `json:"seed,omitempty"`
	OutputFormat   string  `json:"output_format,omitempty"`
}

type GenerateOption func(*GenerateParams)

func WithModel(model string) GenerateOption {
	return func(p *GenerateParams) {
		p.Model = model
	}
}

func WithPrompt(prompt string) GenerateOption {
	return func(p *GenerateParams) {
		p.Prompt = prompt
	}
}

func WithNegativePrompt(negativePrompt string) GenerateOption {
	return func(p *GenerateParams) {
		p.NegativePrompt = negativePrompt
	}
}

func WithMode(mode string) GenerateOption {
	return func(p *GenerateParams) {
		p.Mode = mode
	}
}

func WithImage(image string) GenerateOption {
	return func(p *GenerateParams) {
		p.Image = image
	}
}

func WithStrength(strength float64) GenerateOption {
	return func(p *GenerateParams) {
		p.Strength = strength
	}
}

func WithAspectRatio(aspectRatio string) GenerateOption {
	return func(p *GenerateParams) {
		p.AspectRatio = aspectRatio
	}
}

func WithSeed(seed int) GenerateOption {
	return func(p *GenerateParams) {
		p.Seed = seed
	}
}

func WithOutputFormat(outputFormat string) GenerateOption {
	return func(p *GenerateParams) {
		p.OutputFormat = outputFormat
	}
}

func NewGenerateParams(options ...GenerateOption) (*GenerateParams, error) {
	params := &GenerateParams{}
	for _, option := range options {
		option(params)
	}

	// Determine the endpoint based on the model
	var endpoint string
	switch params.Model {
	case "sd3-large":
		endpoint = "sd3"
	default:
		endpoint = "core" // Default to core if not specified
	}

	// Validate aspect ratio
	if params.AspectRatio != "" {
		validRatios, ok := ValidAspectRatios[endpoint]
		if !ok {
			return nil, fmt.Errorf("unknown endpoint: %s", endpoint)
		}
		isValid := false
		for _, ratio := range validRatios {
			if ratio == params.AspectRatio {
				isValid = true
				break
			}
		}
		if !isValid {
			return nil, fmt.Errorf("invalid aspect ratio %s for endpoint %s", params.AspectRatio, endpoint)
		}
	}

	return params, nil
}

func (p *GenerateParams) toFormData() map[string]string {
	data := make(map[string]string)
	if p.Model != "" {
		data["model"] = p.Model
	}
	data["prompt"] = p.Prompt
	if p.NegativePrompt != "" {
		data["negative_prompt"] = p.NegativePrompt
	}
	if p.Mode != "" {
		data["mode"] = p.Mode
	}
	if p.Strength != 0 {
		data["strength"] = fmt.Sprintf("%f", p.Strength)
	}
	if p.AspectRatio != "" {
		data["aspect_ratio"] = p.AspectRatio
	}
	if p.Seed != 0 {
		data["seed"] = strconv.Itoa(p.Seed)
	}
	if p.OutputFormat != "" {
		data["output_format"] = p.OutputFormat
	}
	return data
}

type GenerateResponse struct {
	Image        string `json:"image,omitempty"`
	FinishReason string `json:"finish_reason"`
	Seed         int    `json:"seed"`
}

// Existing types (TextPrompt, TextToImageParams, etc.) remain unchanged...
