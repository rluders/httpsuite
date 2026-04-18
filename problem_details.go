package httpsuite

import (
	"encoding/json"
	"net/http"
)

const BlankURL = "about:blank"

// ProblemDetails conforms to RFC 9457, providing a standard format for describing errors in HTTP APIs.
type ProblemDetails struct {
	Type       string                 `json:"type"`
	Title      string                 `json:"title"`
	Status     int                    `json:"status"`
	Detail     string                 `json:"detail,omitempty"`
	Instance   string                 `json:"instance,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// ValidationErrorDetail provides structured details about a single validation error.
type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// MarshalJSON serializes RFC 9457 extension members at the top level.
func (p ProblemDetails) MarshalJSON() ([]byte, error) {
	payload := map[string]any{
		"type":   p.Type,
		"title":  p.Title,
		"status": p.Status,
	}
	if p.Detail != "" {
		payload["detail"] = p.Detail
	}
	if p.Instance != "" {
		payload["instance"] = p.Instance
	}
	for key, value := range p.Extensions {
		switch key {
		case "", "type", "title", "status", "detail", "instance", "extensions":
			continue
		default:
			payload[key] = value
		}
	}
	return json.Marshal(payload)
}

// NewProblemDetails creates a ProblemDetails instance with standard fields.
func NewProblemDetails(status int, problemType, title, detail string) *ProblemDetails {
	if status < 100 || status > 599 {
		status = http.StatusInternalServerError
	}
	if problemType == "" {
		problemType = BlankURL
	}
	if title == "" {
		title = http.StatusText(status)
		if title == "" {
			title = "Unknown error"
		}
	}
	return &ProblemDetails{
		Type:   problemType,
		Title:  title,
		Status: status,
		Detail: detail,
	}
}
