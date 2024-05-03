package stability

type StabilityResponse struct {
	Body    []byte              `json:"body,omitempty"`
	Headers map[string][]string `json:"headers,omitempty"`
	Errors  *[]byte             `json:"errors,omitempty"`
}
