package ollama

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// APIError represents an error response from the Ollama API.
type APIError struct {
	StatusCode int
	Status     string
	Message    string
	Body       []byte
}

func (e *APIError) Error() string {
	if e == nil {
		return "ollama: api error"
	}
	if e.Message != "" && e.Status != "" {
		return fmt.Sprintf("ollama: %s: %s", e.Status, e.Message)
	}
	if e.Message != "" {
		return "ollama: " + e.Message
	}
	if e.Status != "" {
		return "ollama: " + e.Status
	}
	return "ollama: api error"
}

// MissingFieldError is returned before a request when a required builder field is empty.
type MissingFieldError struct {
	Field string
}

func (e *MissingFieldError) Error() string {
	return fmt.Sprintf("ollama: missing required field %q", e.Field)
}

func missingField(field string) error {
	return &MissingFieldError{Field: field}
}

func decodeAPIError(resp *http.Response) error {
	data, readErr := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if readErr != nil {
		return fmt.Errorf("ollama: read error response: %w", readErr)
	}

	var parsed struct {
		Error string `json:"error"`
	}
	_ = json.Unmarshal(data, &parsed)

	return &APIError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Message:    parsed.Error,
		Body:       data,
	}
}
