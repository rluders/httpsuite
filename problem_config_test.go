package httpsuite

import (
	"sync"
	"testing"
)

func restoreProblemConfig(t *testing.T) {
	t.Helper()
	SetProblemConfig(NewProblemConfig())
	t.Cleanup(func() { SetProblemConfig(NewProblemConfig()) })
}

func TestSetProblemConfig(t *testing.T) {
	restoreProblemConfig(t)

	cfg := NewProblemConfig()
	cfg.BaseURL = "https://api.example.com/"
	cfg.ErrorTypePaths["unauthorized"] = "errors/unauthorized"
	SetProblemConfig(cfg)

	if got, want := GetProblemTypeURL("unauthorized"), "https://api.example.com/errors/unauthorized"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
	if got := DefaultProblemConfig().TypeURL("not_found_error"); got != "https://api.example.com/errors/not-found" {
		t.Fatalf("unexpected default type URL %q", got)
	}
	if got := NewBadRequestProblem("x").Type; got != "https://api.example.com/errors/bad-request" {
		t.Fatalf("helper did not use configured default, got %q", got)
	}
}

func TestSetProblemConfigCopiesInput(t *testing.T) {
	restoreProblemConfig(t)

	cfg := NewProblemConfig()
	SetProblemConfig(cfg)
	cfg.ErrorTypePaths["validation_error"] = "/mutated"

	if got := GetProblemTypeURL("validation_error"); got != "/errors/validation-error" {
		t.Fatalf("expected default to be isolated from caller map, got %q", got)
	}
}

func TestParseOptionsProblemsOverrideSetProblemConfig(t *testing.T) {
	restoreProblemConfig(t)

	cfg := NewProblemConfig()
	cfg.BaseURL = "https://global.example.com"
	SetProblemConfig(cfg)

	merged := mergeProblemConfig(&ProblemConfig{BaseURL: "https://local.example.com"})
	if got, want := merged.TypeURL("validation_error"), "https://local.example.com/errors/validation-error"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
	merged = mergeProblemConfig(nil)
	if got, want := merged.TypeURL("validation_error"), "https://global.example.com/errors/validation-error"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestSetProblemConfigRestore(t *testing.T) {
	restoreProblemConfig(t)

	SetProblemConfig(ProblemConfig{BaseURL: "https://api.example.com"})
	if got := GetProblemTypeURL("validation_error"); got != BlankURL {
		t.Fatalf("expected %q, got %q", BlankURL, got)
	}

	SetProblemConfig(NewProblemConfig())
	if got := GetProblemTypeURL("validation_error"); got != "/errors/validation-error" {
		t.Fatalf("expected built-in default, got %q", got)
	}
}

func TestSetProblemConfigConcurrent(t *testing.T) {
	restoreProblemConfig(t)

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			SetProblemConfig(NewProblemConfig())
		}()
		go func() {
			defer wg.Done()
			_ = GetProblemTypeURL("validation_error")
			_ = DefaultProblemConfig()
		}()
	}
	wg.Wait()
}

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
