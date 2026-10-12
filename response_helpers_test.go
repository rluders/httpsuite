package httpsuite

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseHelpers(t *testing.T) {
	t.Parallel()

	t.Run("ok helper", func(t *testing.T) {
		w := httptest.NewRecorder()
		OK(w, testResponse{Key: "value"})
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("ok with meta helper", func(t *testing.T) {
		w := httptest.NewRecorder()
		OKWithMeta(w, []string{"a"}, NewPageMeta(1, 10, 15))
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("created helper", func(t *testing.T) {
		w := httptest.NewRecorder()
		Created(w, testResponse{Key: "value"}, "/users/1")
		if w.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
		}
		if got := w.Header().Get("Location"); got != "/users/1" {
			t.Fatalf("expected location header, got %q", got)
		}
	})

	t.Run("accepted helper", func(t *testing.T) {
		w := httptest.NewRecorder()
		Accepted(w, testResponse{Key: "value"})
		if w.Code != http.StatusAccepted {
			t.Fatalf("expected status %d, got %d", http.StatusAccepted, w.Code)
		}
		if w.Body.Len() == 0 {
			t.Fatal("expected response body")
		}
	})

	t.Run("no content helper", func(t *testing.T) {
		w := httptest.NewRecorder()
		NoContent(w)
		if w.Code != http.StatusNoContent {
			t.Fatalf("expected status %d, got %d", http.StatusNoContent, w.Code)
		}
		if w.Body.Len() != 0 {
			t.Fatalf("expected empty body, got %q", w.Body.String())
		}
	})

	t.Run("reply builder accepted and no content", func(t *testing.T) {
		w := httptest.NewRecorder()
		Reply().Header("X-Job", "1").Meta(map[string]int{"n": 1}).Accepted(w, testResponse{Key: "value"})
		if w.Code != http.StatusAccepted || w.Header().Get("X-Job") != "1" {
			t.Fatalf("unexpected accepted response: %d %v", w.Code, w.Header())
		}

		w = httptest.NewRecorder()
		Reply().Header("X-Job", "2").NoContent(w)
		if w.Code != http.StatusNoContent || w.Header().Get("X-Job") != "2" {
			t.Fatalf("unexpected no content response: %d %v", w.Code, w.Header())
		}
		if w.Body.Len() != 0 {
			t.Fatalf("expected empty body, got %q", w.Body.String())
		}
		if ct := w.Header().Get("Content-Type"); ct != "" {
			t.Fatalf("expected no content type, got %q", ct)
		}
	})

	t.Run("problem helper", func(t *testing.T) {
		w := httptest.NewRecorder()
		ProblemResponse(w, NewNotFoundProblem("user missing"))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}
