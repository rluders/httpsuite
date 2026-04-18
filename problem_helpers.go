package httpsuite

import "net/http"

// ProblemBadRequest returns a bad request problem builder.
func ProblemBadRequest(detail string) *ProblemBuilder {
	return Problem(http.StatusBadRequest).
		Type(GetProblemTypeURL("bad_request_error")).
		Title("Bad Request").
		Detail(detail)
}

// ProblemNotFound returns a not found problem builder.
func ProblemNotFound(detail string) *ProblemBuilder {
	return Problem(http.StatusNotFound).
		Type(GetProblemTypeURL("not_found_error")).
		Title("Not Found").
		Detail(detail)
}

// NewBadRequestProblem returns a ready-to-use bad request problem.
func NewBadRequestProblem(detail string) *ProblemDetails {
	return ProblemBadRequest(detail).Build()
}

// NewNotFoundProblem returns a ready-to-use not found problem.
func NewNotFoundProblem(detail string) *ProblemDetails {
	return ProblemNotFound(detail).Build()
}
