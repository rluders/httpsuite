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
	Data   T       `json:"data,omitempty"`
	Errors []Error `json:"errors,omitempty"`
	Meta   *Meta   `json:"meta,omitempty"`
}

// Error represents an error in the aPI response, with a structured format to describe issues in a consistent manner.
type Error struct {
	// Code unique error code or HTTP status code for categorizing the error
	Code int `json:"code"`
	// Message user-friendly message describing the error.
	Message string `json:"message"`
	// Details additional details about the error, often used for validation errors.
	Details interface{} `json:"details,omitempty"`
}

// Meta provides additional information about the response, such as pagination details.
// This is particularly useful for endpoints returning lists of data.
type Meta struct {
	// Page the current page number
	Page int `json:"page,omitempty"`
	// PageSize the number of items per page
	PageSize int `json:"page_size,omitempty"`
	// TotalPages the total number of pages available.
	TotalPages int `json:"total_pages,omitempty"`
	// TotalItems the total number of items across all pages.
	TotalItems int `json:"total_items,omitempty"`
}

// SendResponse sends a JSON response to the client, using a unified structure for both success and error responses.
// T represents the type of the `data` payload. This function automatically adapts the response structure
// based on whether `data` or `errors` is provided, promoting a consistent API format.
//
// Parameters:
//   - w: The http.ResponseWriter to send the response.
//   - code: HTTP status code to indicate success or failure.
//   - data: The main payload of the response. Use `nil` for error responses.
//   - errs: A slice of Error structs to describe issues. Use `nil` for successful responses.
//   - meta: Optional metadata, such as pagination information. Use `nil` if not needed.
func SendResponse[T any](w http.ResponseWriter, code int, data T, errs []Error, meta *Meta) {
	w.Header().Set("Content-Type", "application/json")

	response := &Response[T]{
		Data:   data,
		Errors: errs,
		Meta:   meta,
	}

	// Set the status code after encoding to ensure no issues with writing the response body
	w.WriteHeader(code)

	// Attempt to encode the response as JSON
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(response); err != nil {
		log.Printf("Error writing response: %v", err)

		errResponse := `{"errors":[{"code":500,"message":"Internal Server Error"}]}`
		http.Error(w, errResponse, http.StatusInternalServerError)
		return
	}

	// Set the status code after success encoding
	w.WriteHeader(code)

	// Write the encoded response to the ResponseWriter
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
