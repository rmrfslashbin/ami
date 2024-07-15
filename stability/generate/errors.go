package generate

// file: stability/generate/errors.go

import "fmt"

type APIError struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Stability API Error: [%s] %s - %s", e.ID, e.Name, e.Message)
}

// Add more specific error types if needed...
