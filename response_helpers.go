package httpsuite

import "net/http"

// OK writes a 200 JSON response without metadata.
func OK[T any](w http.ResponseWriter, data T) {
	Reply().OK(w, data)
}

// OKWithMeta writes a 200 JSON response with metadata.
func OKWithMeta[T any](w http.ResponseWriter, data T, meta any) {
	Reply().Meta(meta).OK(w, data)
}

// Created writes a 201 JSON response and optionally sets the Location header.
func Created[T any](w http.ResponseWriter, data T, location string) {
	Reply().Created(w, data, location)
}

// Accepted writes a 202 JSON response without metadata.
func Accepted[T any](w http.ResponseWriter, data T) {
	Reply().Accepted(w, data)
}

// NoContent writes a 204 response with no body.
func NoContent(w http.ResponseWriter) {
	Reply().NoContent(w)
}

// ProblemResponse writes a problem response using the problem's status.
func ProblemResponse(w http.ResponseWriter, problem *ProblemDetails) {
	if problem == nil {
		problem = NewProblemDetails(http.StatusInternalServerError, GetProblemTypeURL("server_error"), "Internal Server Error", "")
	}
	Reply().Problem(w, problem)
}
