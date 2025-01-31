package httpsuite

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// ValidationErrorDetail provides structured details about a single validation error.
type ValidationErrorDetail struct {
	Field   string `json:"field"`   // The name of the field that failed validation.
	Message string `json:"message"` // A human-readable message describing the error.
}

// NewValidationProblemDetails creates a ProblemDetails instance based on validation errors.
// It maps field-specific validation errors into structured details.
func NewValidationProblemDetails(err error) *ProblemDetails {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		// If the error is not of type ValidationErrors, return a generic problem response.
		return NewProblemDetails(http.StatusBadRequest, "Invalid Request", "Invalid data format or structure")
	}

	// Collect structured details about each validation error.
	errorDetails := make([]ValidationErrorDetail, len(validationErrors))
	for i, vErr := range validationErrors {
		errorDetails[i] = ValidationErrorDetail{
			Field:   vErr.Field(),
			Message: formatValidationMessage(vErr),
		}
	}

	return &ProblemDetails{
		Type:   "https://example.com/validation-error",
		Title:  "Validation Error",
		Status: http.StatusBadRequest,
		Detail: "One or more fields failed validation.",
		Extensions: map[string]interface{}{
			"errors": errorDetails,
		},
	}
}

// formatValidationMessage generates a descriptive message for a validation error.
func formatValidationMessage(vErr validator.FieldError) string {
	return vErr.Field() + " failed " + vErr.Tag() + " validation"
}

// IsRequestValid validates the provided request struct using the go-playground/validator package.
// It returns a ProblemDetails instance if validation fails, or nil if the request is valid.
func IsRequestValid(request any) *ProblemDetails {
	err := validate.Struct(request)
	if err != nil {
		return NewValidationProblemDetails(err)
	}
	return nil
}
