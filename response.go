package httpsuite

import (
	"encoding/json"
	"log"
	"net/http"
)

// Response represents the structure of an HTTP response, including a status code, message, and optional body.
type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Body    T      `json:"body,omitempty"`
}

// Marshal serializes the Response struct into a JSON byte slice.
// It logs an error if marshalling fails.
func (r *Response[T]) Marshal() []byte {
	jsonResponse, err := json.Marshal(r)
	if err != nil {
		log.Printf("failed to marshal response: %v", err)
	}

	return jsonResponse
}

// SendResponse creates a Response struct, serializes it to JSON, and writes it to the provided http.ResponseWriter.
// If the body parameter is non-nil, it will be included in the response body.
func SendResponse[T any](w http.ResponseWriter, message string, code int, body *T) {
	response := &Response[T]{
		Code:    code,
		Message: message,
	}
	if body != nil {
		response.Body = *body
	}

	writeResponse[T](w, response)
}

// writeResponse serializes a Response and writes it to the http.ResponseWriter with appropriate headers.
// If an error occurs during the write, it logs the error and sends a 500 Internal Server Error response.
func writeResponse[T any](w http.ResponseWriter, r *Response[T]) {
	jsonResponse := r.Marshal()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)

	if _, err := w.Write(jsonResponse); err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
