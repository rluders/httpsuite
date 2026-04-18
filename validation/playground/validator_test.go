package playground

import (
	"testing"

	"github.com/rluders/httpsuite/v3"
)

type request struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required,min=18"`
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
	errorsValue, ok := problem.Extensions["errors"].([]httpsuite.ValidationErrorDetail)
	if !ok {
		t.Fatalf("expected validation error details, got %#v", problem.Extensions["errors"])
	}
	if len(errorsValue) == 0 || errorsValue[0].Field != "name" {
		t.Fatalf("expected json field name in validation error, got %#v", errorsValue)
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
