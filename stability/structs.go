package stability

type StabilityResponse struct {
	Body    []byte
	Headers map[string][]string
}

type StabilityV3Request struct {
	Prompt         string `json:"prompt"`
	AspectRatio    string `json:"aspect_ratio,omitempty"`
	Mode           string `json:"mode,omitempty"`
	NegativePrompt string `json:"negative_prompt,omitempty"`
	Model          string `json:"model,omitempty"`
	Seed           int    `json:"seed,omitempty"`
	OutputFormat   string `json:"output_format,omitempty"`
}

type StabilityV3Response struct {
	Data *StabilityV3ImageData `json:"data"`
	Json *StabilityV3ImageJSON `json:"json"`
}

// StabilityV3ImageData is the response data for the V3 endpoint with image/* accept header
type StabilityV3ImageData struct {
	ContextType  string `json:"context_type"`
	FinishReason string `json:"finish_reason"`
	Seed         int    `json:"seed,omitempty"`
	ImageData    []byte `json:"image_data"`
}

// StabilityV3ImageJSON is the response data for the V3 endpoint with application/json accept header
type StabilityV3ImageJSON struct {
	Image        string `json:"image"`
	FinishReason string `json:"finish_reason"`
	Seed         int    `json:"seed"`
}
