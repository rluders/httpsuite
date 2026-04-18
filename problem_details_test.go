package httpsuite

import "testing"

func TestNewProblemDetailsDefaults(t *testing.T) {
	t.Parallel()

	details := NewProblemDetails(700, "", "", "broken")
	if details.Status != 500 {
		t.Fatalf("expected status 500, got %d", details.Status)
	}
	if details.Type != BlankURL {
		t.Fatalf("expected type %q, got %q", BlankURL, details.Type)
	}
	if details.Title != "Internal Server Error" {
		t.Fatalf("expected fallback title, got %q", details.Title)
	}
}
