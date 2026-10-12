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

// ProblemUnauthorized returns an unauthorized problem builder.
func ProblemUnauthorized(detail string) *ProblemBuilder {
	return Problem(http.StatusUnauthorized).
		Type(GetProblemTypeURL("unauthorized_error")).
		Title("Unauthorized").
		Detail(detail)
}

// ProblemForbidden returns a forbidden problem builder.
func ProblemForbidden(detail string) *ProblemBuilder {
	return Problem(http.StatusForbidden).
		Type(GetProblemTypeURL("forbidden_error")).
		Title("Forbidden").
		Detail(detail)
}

// ProblemConflict returns a conflict problem builder.
func ProblemConflict(detail string) *ProblemBuilder {
	return Problem(http.StatusConflict).
		Type(GetProblemTypeURL("conflict_error")).
		Title("Conflict").
		Detail(detail)
}

// ProblemUnprocessable returns an unprocessable entity problem builder.
func ProblemUnprocessable(detail string) *ProblemBuilder {
	return Problem(http.StatusUnprocessableEntity).
		Type(GetProblemTypeURL("unprocessable_error")).
		Title("Unprocessable Entity").
		Detail(detail)
}

// ProblemInternal returns an internal server error problem builder.
func ProblemInternal(detail string) *ProblemBuilder {
	return Problem(http.StatusInternalServerError).
		Type(GetProblemTypeURL("server_error")).
		Title("Internal Server Error").
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

// NewUnauthorizedProblem returns a ready-to-use unauthorized problem.
func NewUnauthorizedProblem(detail string) *ProblemDetails {
	return ProblemUnauthorized(detail).Build()
}

// NewForbiddenProblem returns a ready-to-use forbidden problem.
func NewForbiddenProblem(detail string) *ProblemDetails {
	return ProblemForbidden(detail).Build()
}

// NewConflictProblem returns a ready-to-use conflict problem.
func NewConflictProblem(detail string) *ProblemDetails {
	return ProblemConflict(detail).Build()
}

// NewUnprocessableProblem returns a ready-to-use unprocessable entity problem.
func NewUnprocessableProblem(detail string) *ProblemDetails {
	return ProblemUnprocessable(detail).Build()
}

// NewInternalProblem returns a ready-to-use internal server error problem.
func NewInternalProblem(detail string) *ProblemDetails {
	return ProblemInternal(detail).Build()
}
