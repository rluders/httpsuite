package httpsuite

import (
	"net/http"
	"testing"
)

func TestProblemBuilderHelpers(t *testing.T) {
	t.Parallel()

	badRequest := ProblemBadRequest("invalid page").Build()
	if badRequest.Status != http.StatusBadRequest {
		t.Fatalf("expected bad request status, got %d", badRequest.Status)
	}

	notFound := ProblemNotFound("user missing").Build()
	if notFound.Status != http.StatusNotFound {
		t.Fatalf("expected not found status, got %d", notFound.Status)
	}
}

func TestStatusProblemHelpers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		builder func(string) *ProblemBuilder
		direct  func(string) *ProblemDetails
		status  int
		title   string
		errType string
	}{
		{"unauthorized", ProblemUnauthorized, NewUnauthorizedProblem, http.StatusUnauthorized, "Unauthorized", "unauthorized_error"},
		{"forbidden", ProblemForbidden, NewForbiddenProblem, http.StatusForbidden, "Forbidden", "forbidden_error"},
		{"conflict", ProblemConflict, NewConflictProblem, http.StatusConflict, "Conflict", "conflict_error"},
		{"unprocessable", ProblemUnprocessable, NewUnprocessableProblem, http.StatusUnprocessableEntity, "Unprocessable Entity", "unprocessable_error"},
		{"internal", ProblemInternal, NewInternalProblem, http.StatusInternalServerError, "Internal Server Error", "server_error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for _, got := range []*ProblemDetails{tt.builder("some detail").Build(), tt.direct("some detail")} {
				if got.Status != tt.status {
					t.Errorf("status = %d, want %d", got.Status, tt.status)
				}
				if got.Title != tt.title {
					t.Errorf("title = %q, want %q", got.Title, tt.title)
				}
				if got.Detail != "some detail" {
					t.Errorf("detail = %q", got.Detail)
				}
				if want := GetProblemTypeURL(tt.errType); got.Type != want || want == BlankURL {
					t.Errorf("type = %q, want %q", got.Type, want)
				}
			}
		})
	}
}

func TestDirectProblemHelpers(t *testing.T) {
	t.Parallel()

	badRequest := NewBadRequestProblem("invalid payload")
	if badRequest.Status != http.StatusBadRequest {
		t.Fatalf("expected bad request status, got %d", badRequest.Status)
	}

	notFound := NewNotFoundProblem("user missing")
	if notFound.Status != http.StatusNotFound {
		t.Fatalf("expected not found status, got %d", notFound.Status)
	}
}
