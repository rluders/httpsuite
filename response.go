package httpsuite

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// Response represents the structure of an HTTP response, including a status code, message, and optional body.
// T represents the type of the `Data` field, allowing this structure to be used flexibly across different endpoints.
type Response[T any] struct {
	Data T     `json:"data,omitempty"`
	Meta *Meta `json:"meta,omitempty"`
}

// Meta provides additional information about the response, such as pagination details.
type Meta struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"page_size,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	TotalItems int `json:"total_items,omitempty"`
}

// ProblemDetails conforms to RFC 9457, providing a standard format for describing errors in HTTP APIs.
type ProblemDetails struct {
	Type       string                 `json:"type"`                 // A URI reference identifying the problem type.
	Title      string                 `json:"title"`                // A short, human-readable summary of the problem.
	Status     int                    `json:"status"`               // The HTTP status code.
	Detail     string                 `json:"detail,omitempty"`     // Detailed explanation of the problem.
	Instance   string                 `json:"instance,omitempty"`   // A URI reference identifying the specific instance of the problem.
	Extensions map[string]interface{} `json:"extensions,omitempty"` // Custom fields for additional details.
}

// NewProblemDetails creates a ProblemDetails instance with standard fields.
func NewProblemDetails(status int, title, detail string) *ProblemDetails {
	return &ProblemDetails{
		Type:   "about:blank", // Replace with a custom URI if desired.
		Title:  title,
		Status: status,
		Detail: detail,
	}
}

// SendResponse sends a JSON response to the client, supporting both success and error scenarios.
//
// Parameters:
//   - w: The http.ResponseWriter to send the response.
//   - code: HTTP status code to indicate success or failure.
//   - data: The main payload of the response (only for successful responses).
//   - problem: An optional ProblemDetails struct (used for error responses).
//   - meta: Optional metadata for successful responses (e.g., pagination details).
func SendResponse[T any](w http.ResponseWriter, code int, data T, problem *ProblemDetails, meta *Meta) {

	// Handle error responses
	if code >= 400 && problem != nil {
		writeProblemDetail(w, code, problem)
		return
	}

	// Construct and encode the success response
	response := &Response[T]{
		Data: data,
		Meta: meta,
	}

	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(response); err != nil {
		log.Printf("Error writing response: %v", err)

		// Internal server error fallback using ProblemDetails
		internalError := NewProblemDetails(http.StatusInternalServerError, "Internal Server Error", err.Error())
		writeProblemDetail(w, http.StatusInternalServerError, internalError)
		return
	}

	// Send the success response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Printf("Failed to write response body (status=%d): %v", code, err)
	}
}

func writeProblemDetail(w http.ResponseWriter, code int, problem *ProblemDetails) {
	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.WriteHeader(problem.Status)
	if err := json.NewEncoder(w).Encode(problem); err != nil {
		log.Printf("Failed to encode problem details: %v", err)
	}
}
