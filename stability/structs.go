package stability

type StabilityResponse struct {
	Body    []byte
	Headers map[string][]string
	Errors  *[]byte
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
	Data   *StabilityV3ImageData `json:"data"`
	Json   *StabilityV3ImageJSON `json:"json"`
	Errors *StabilityV3Errors    `json:"errors"`
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

type StabilityV3Errors struct {
	Name   string   `json:"name"`
	Errors []string `json:"errors"`
}

type ResponseUser struct {
	Id             string             `json:"id"`
	Email          string             `json:"email"`
	ProfilePicture string             `json:"profile_picture"`
	Organizations  []ResponseUserOrg  `json:"organizations"`
	Error          *ResponseUserError `json:"error,omitempty"`
}

type ResponseUserOrg struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	IsDefault bool   `json:"is_default"`
}

type ResponseUserError struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

type ResponseUserBalance struct {
	Credits float64            `json:"credits"`
	Error   *ResponseUserError `json:"error,omitempty"`
}

type BalanceInput struct {
	Organization           *string
	StabilityClientID      *string
	StabilityClientVersion *string
}
