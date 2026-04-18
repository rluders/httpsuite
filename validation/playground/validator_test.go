package playground

import (
	"testing"

	"github.com/rluders/httpsuite/v3"
)

type request struct {
	Name string `validate:"required"`
	Age  int    `validate:"required,min=18"`
}

func TestValidate(t *testing.T) {
	t.Parallel()

	validator := New()
	problem := validator.Validate(request{Age: 17})
	if problem == nil {
		t.Fatal("expected validation problem, got nil")
	}
	if problem.Title != "Validation Error" {
		t.Fatalf("expected title %q, got %q", "Validation Error", problem.Title)
	}
}

func TestRegisterDefault(t *testing.T) {
	httpsuite.ClearValidator()
	t.Cleanup(httpsuite.ClearValidator)

	validator := RegisterDefault()
	if validator == nil {
		t.Fatal("expected validator, got nil")
	}
	if httpsuite.DefaultValidator() == nil {
		t.Fatal("expected default validator to be registered")
	}
}
