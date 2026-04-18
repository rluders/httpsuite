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

	t.Run("problem helper", func(t *testing.T) {
		w := httptest.NewRecorder()
		ProblemResponse(w, NewNotFoundProblem("user missing"))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}
