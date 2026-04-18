package httpsuite

import "net/http"

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
