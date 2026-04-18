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
