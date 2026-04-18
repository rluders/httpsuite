package httpsuite

import (
	"net/http"
	"testing"
)

func TestProblemBuilder(t *testing.T) {
	t.Parallel()

	problem := Problem(http.StatusBadRequest).
		Type(GetProblemTypeURL("bad_request_error")).
		Title("Validation Error").
		Detail("invalid payload").
		Instance("/users").
		Extension("request_id", "req-123").
		Extensions(map[string]any{"trace_id": "trace-456"}).
		Build()

	if problem.Status != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, problem.Status)
	}
	if problem.Type != GetProblemTypeURL("bad_request_error") {
		t.Fatalf("unexpected type %q", problem.Type)
	}
	if problem.Instance != "/users" {
		t.Fatalf("unexpected instance %q", problem.Instance)
	}
	if problem.Extensions["request_id"] != "req-123" {
		t.Fatalf("expected request_id extension")
	}
	if problem.Extensions["trace_id"] != "trace-456" {
		t.Fatalf("expected trace_id extension")
	}
}
