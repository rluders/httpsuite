package httpsuite

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func writeResponse[T any](w http.ResponseWriter, code int, data T, problem *ProblemDetails, meta any, headers http.Header) {
	if code >= 400 && problem != nil {
		writeProblemDetail(w, code, problem, headers)
		return
	}

	response := &Response[T]{
		Data: data,
		Meta: meta,
	}

	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(response); err != nil {
		log.Printf("Error writing response: %v", err)

		internalError := NewProblemDetails(
			http.StatusInternalServerError,
			GetProblemTypeURL("server_error"),
			"Internal Server Error",
			err.Error(),
		)
		writeProblemDetail(w, http.StatusInternalServerError, internalError, headers)
		return
	}

	applyHeaders(w, headers)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Printf("Failed to write response body (status=%d): %v", code, err)
	}
}

func writeProblemDetail(w http.ResponseWriter, code int, problem *ProblemDetails, headers http.Header) {
	effectiveStatus := code
	if effectiveStatus < 400 || effectiveStatus > 599 {
		effectiveStatus = problem.Status
	}
	if effectiveStatus < 400 || effectiveStatus > 599 {
		effectiveStatus = http.StatusInternalServerError
	}

	normalized := *problem
	normalized.Status = effectiveStatus

	applyHeaders(w, headers)
	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.WriteHeader(effectiveStatus)
	if err := json.NewEncoder(w).Encode(normalized); err != nil {
		log.Printf("Failed to encode problem details: %v", err)
	}
}

func applyHeaders(w http.ResponseWriter, headers http.Header) {
	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
}
