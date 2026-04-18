package httpsuite

import (
	"net/http"
	"sync"
	"testing"
)

func TestValidateRequest(t *testing.T) {
	t.Parallel()

	problem := &ProblemDetails{
		Type:   GetProblemTypeURL("validation_error"),
		Title:  "Validation Error",
		Status: http.StatusBadRequest,
		Detail: "One or more fields failed validation.",
	}

	if got := ValidateRequest(&testRequest{}, nil); got != nil {
		t.Fatalf("expected nil validation problem, got %#v", got)
	}
	if got := ValidateRequest(&testRequest{}, stubValidator{problem: problem}); got != problem {
		t.Fatalf("expected validation problem to be returned")
	}
}

func TestValidationProblemStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		problem *ProblemDetails
		want    int
	}{
		{
			name: "nil problem",
			want: http.StatusBadRequest,
		},
		{
			name: "valid status",
			problem: &ProblemDetails{
				Status: http.StatusUnprocessableEntity,
			},
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid status",
			problem: &ProblemDetails{
				Status: 0,
			},
			want: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validationProblemStatus(tt.problem); got != tt.want {
				t.Fatalf("expected status %d, got %d", tt.want, got)
			}
		})
	}
}

func TestDefaultValidatorLifecycle(t *testing.T) {
	ClearValidator()
	t.Cleanup(ClearValidator)

	if DefaultValidator() != nil {
		t.Fatalf("expected nil default validator")
	}

	validator := stubValidator{}
	SetValidator(validator)
	if DefaultValidator() == nil {
		t.Fatalf("expected default validator to be set")
	}

	ClearValidator()
	if DefaultValidator() != nil {
		t.Fatalf("expected default validator to be cleared")
	}
}

func TestDefaultValidatorConcurrentAccess(t *testing.T) {
	ClearValidator()
	t.Cleanup(ClearValidator)

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				SetValidator(stubValidator{})
			} else {
				_ = DefaultValidator()
				ClearValidator()
			}
		}(i)
	}
	wg.Wait()
}
