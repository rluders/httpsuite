package httpsuite

import (
	"github.com/go-playground/validator/v10"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestValidationRequest struct {
	Name string `validate:"required"`
	Age  int    `validate:"required,min=18"`
}

func TestNewValidationErrors(t *testing.T) {
	validate := validator.New()
	request := TestValidationRequest{} // Missing required fields to trigger validation errors

	err := validate.Struct(request)
	if err == nil {
		t.Fatal("Expected validation errors, but got none")
	}

	validationErrors := NewValidationErrors(err)

	expectedErrors := map[string][]string{
		"Name": {"Name required"},
		"Age":  {"Age required"},
	}

	assert.Equal(t, expectedErrors, validationErrors.Errors)
}

func TestIsRequestValid(t *testing.T) {
	tests := []struct {
		name           string
		request        TestValidationRequest
		expectedErrors *ValidationErrors
	}{
		{
			name:           "Valid request",
			request:        TestValidationRequest{Name: "Alice", Age: 25},
			expectedErrors: nil, // No errors expected for valid input
		},
		{
			name:    "Missing Name and Age below minimum",
			request: TestValidationRequest{Age: 17},
			expectedErrors: &ValidationErrors{
				Errors: map[string][]string{
					"Name": {"Name required"},
					"Age":  {"Age min"},
				},
			},
		},
		{
			name:    "Missing Age",
			request: TestValidationRequest{Name: "Alice"},
			expectedErrors: &ValidationErrors{
				Errors: map[string][]string{
					"Age": {"Age required"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := IsRequestValid(tt.request)
			if tt.expectedErrors == nil {
				assert.Nil(t, errs)
			} else {
				assert.NotNil(t, errs)
				assert.Equal(t, tt.expectedErrors.Errors, errs.Errors)
			}
		})
	}
}
