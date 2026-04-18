package httpsuite

import "net/http"

// Reply starts a fluent response helper configuration.
func Reply() *ReplyBuilder {
	return &ReplyBuilder{}
}

// SendResponse sends a JSON response to the client, supporting both success and error scenarios.
func SendResponse[T any](w http.ResponseWriter, code int, data T, problem *ProblemDetails, meta any) {
	writeResponse(w, code, data, problem, meta, nil)
}
