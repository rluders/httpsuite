package httpsuite

import "testing"

func TestProblemConfigTypeURL(t *testing.T) {
	t.Parallel()

	config := ProblemConfig{
		BaseURL: "https://api.example.com/",
		ErrorTypePaths: map[string]string{
			"validation_error": "errors/validation-error",
		},
	}

	got := config.TypeURL("validation_error")
	want := "https://api.example.com/errors/validation-error"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestProblemConfigTypeURLUnknown(t *testing.T) {
	t.Parallel()

	config := NewProblemConfig()
	if got := config.TypeURL("missing"); got != BlankURL {
		t.Fatalf("expected %q, got %q", BlankURL, got)
	}
}

func TestDefaultProblemConfigReturnsCopy(t *testing.T) {
	t.Parallel()

	config := DefaultProblemConfig()
	config.ErrorTypePaths["validation_error"] = "/custom"

	if got := GetProblemTypeURL("validation_error"); got != "/errors/validation-error" {
		t.Fatalf("expected default config to stay unchanged, got %q", got)
	}
}
