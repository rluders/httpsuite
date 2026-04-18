package httpsuite

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseBuilderSuccess(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	Respond([]string{"a", "b"}).
		Status(http.StatusCreated).
		Meta(NewPageMeta(1, 10, 15)).
		Header("X-Request-ID", "req-123").
		Write(w)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
	if got := w.Header().Get("X-Request-ID"); got != "req-123" {
		t.Fatalf("expected request id header, got %q", got)
	}
	if got := w.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("unexpected content type %q", got)
	}
	if !json.Valid(w.Body.Bytes()) {
		t.Fatalf("expected valid json response")
	}
}

func TestResponseBuilderProblem(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	RespondProblem(ProblemNotFound("user missing").Build()).
		Header("X-Trace-ID", "trace-123").
		Write(w)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
	if got := w.Header().Get("X-Trace-ID"); got != "trace-123" {
		t.Fatalf("expected trace id header, got %q", got)
	}

	var problem ProblemDetails
	if err := json.NewDecoder(w.Body).Decode(&problem); err != nil {
		t.Fatalf("decode problem details: %v", err)
	}
	if problem.Status != http.StatusNotFound {
		t.Fatalf("expected normalized problem status %d, got %d", http.StatusNotFound, problem.Status)
	}
}

func TestReplyBuilderHelpers(t *testing.T) {
	t.Parallel()

	t.Run("meta then ok", func(t *testing.T) {
		w := httptest.NewRecorder()
		Reply().
			Meta(NewPageMeta(1, 10, 15)).
			Header("X-Request-ID", "req-123").
			OK(w, []string{"a", "b"})

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}
		if got := w.Header().Get("X-Request-ID"); got != "req-123" {
			t.Fatalf("expected request id header, got %q", got)
		}
	})

	t.Run("headers then created", func(t *testing.T) {
		w := httptest.NewRecorder()
		Reply().
			Headers(http.Header{"X-Trace-ID": []string{"trace-123", "trace-456"}}).
			Created(w, testResponse{Key: "value"}, "/users/1")

		if w.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
		}
		if got := w.Header().Get("Location"); got != "/users/1" {
			t.Fatalf("expected location header, got %q", got)
		}
		if got := w.Header().Values("X-Trace-ID"); len(got) != 2 || got[0] != "trace-123" || got[1] != "trace-456" {
			t.Fatalf("expected repeated trace headers, got %#v", got)
		}
	})
}
