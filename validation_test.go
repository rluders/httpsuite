package httpsuite

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestValidationRequest struct {
	Name string `validate:"required"`
	Age  int    `validate:"required,min=18"`
}

func TestNewValidationProblemDetails(t *testing.T) {
	validate := validator.New()
	request := TestValidationRequest{} // Missing required fields to trigger validation errors

	err := validate.Struct(request)
	if err == nil {
		t.Fatal("Expected validation errors, but got none")
	}

	validationProblem := NewValidationProblemDetails(err)

	expectedProblem := &ProblemDetails{
		Type:   "https://example.com/validation-error",
		Title:  "Validation Error",
		Status: 400,
		Detail: "One or more fields failed validation.",
		Extensions: map[string]interface{}{
			"errors": []ValidationErrorDetail{
				{Field: "Name", Message: "Name failed required validation"},
				{Field: "Age", Message: "Age failed required validation"},
			},
		},
	}

	assert.Equal(t, expectedProblem.Type, validationProblem.Type)
	assert.Equal(t, expectedProblem.Title, validationProblem.Title)
	assert.Equal(t, expectedProblem.Status, validationProblem.Status)
	assert.Equal(t, expectedProblem.Detail, validationProblem.Detail)
	assert.ElementsMatch(t, expectedProblem.Extensions["errors"], validationProblem.Extensions["errors"])
}

func TestIsRequestValid(t *testing.T) {
	tests := []struct {
		name            string
		request         TestValidationRequest
		expectedProblem *ProblemDetails
	}{
		{
			name:            "Valid request",
			request:         TestValidationRequest{Name: "Alice", Age: 25},
			expectedProblem: nil, // No errors expected for valid input
		},
		{
			name:    "Missing Name and Age below minimum",
			request: TestValidationRequest{Age: 17},
			expectedProblem: &ProblemDetails{
				Type:   "https://example.com/validation-error",
				Title:  "Validation Error",
				Status: 400,
				Detail: "One or more fields failed validation.",
				Extensions: map[string]interface{}{
					"errors": []ValidationErrorDetail{
						{Field: "Name", Message: "Name failed required validation"},
						{Field: "Age", Message: "Age failed min validation"},
					},
				},
			},
		},
		{
			name:    "Missing Age",
			request: TestValidationRequest{Name: "Alice"},
			expectedProblem: &ProblemDetails{
				Type:   "https://example.com/validation-error",
				Title:  "Validation Error",
				Status: 400,
				Detail: "One or more fields failed validation.",
				Extensions: map[string]interface{}{
					"errors": []ValidationErrorDetail{
						{Field: "Age", Message: "Age failed required validation"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			problem := IsRequestValid(tt.request)

			if tt.expectedProblem == nil {
				assert.Nil(t, problem)
			} else {
				assert.NotNil(t, problem)
				assert.Equal(t, tt.expectedProblem.Type, problem.Type)
				assert.Equal(t, tt.expectedProblem.Title, problem.Title)
				assert.Equal(t, tt.expectedProblem.Status, problem.Status)
				assert.Equal(t, tt.expectedProblem.Detail, problem.Detail)
				assert.ElementsMatch(t, tt.expectedProblem.Extensions["errors"], problem.Extensions["errors"])
			}
		})
	}
}
