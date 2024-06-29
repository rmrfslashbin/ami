package stability

type StabilityResponse struct {
	StatusCode int                 `json:"status_code,omitempty"`
	Body       []byte              `json:"body,omitempty"`
	Headers    map[string][]string `json:"headers,omitempty"`
}
