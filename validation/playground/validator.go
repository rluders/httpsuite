package playground

import (
	"errors"
	"net/http"

	playgroundvalidator "github.com/go-playground/validator/v10"
	"github.com/rluders/httpsuite/v3"
)

// Validator adapts go-playground/validator to the httpsuite.Validator interface.
type Validator struct {
	validate *playgroundvalidator.Validate
	problems httpsuite.ProblemConfig
}

// New returns a validator with the default go-playground configuration.
func New() *Validator {
	return NewWithValidator(playgroundvalidator.New(), nil)
}

// RegisterDefault installs a playground validator as the package-level default in httpsuite.
func RegisterDefault() *Validator {
	validator := New()
	httpsuite.SetValidator(validator)
	return validator
}

// NewWithValidator returns a validator using a custom go-playground validator.
func NewWithValidator(validate *playgroundvalidator.Validate, problems *httpsuite.ProblemConfig) *Validator {
	if validate == nil {
		validate = playgroundvalidator.New()
	}

	return &Validator{
		validate: validate,
		problems: mergeProblems(problems),
	}
}

// Validate validates the request and converts errors into ProblemDetails.
func (v *Validator) Validate(request any) *httpsuite.ProblemDetails {
	if err := v.validate.Struct(request); err != nil {
		return v.problemDetails(err)
	}
	return nil
}

func (v *Validator) problemDetails(err error) *httpsuite.ProblemDetails {
	var validationErrors playgroundvalidator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return httpsuite.NewProblemDetails(
			http.StatusBadRequest,
			v.problems.TypeURL("bad_request_error"),
			"Invalid Request",
			"Invalid data format or structure",
		)
	}

	errorDetails := make([]httpsuite.ValidationErrorDetail, len(validationErrors))
	for i, validationErr := range validationErrors {
		errorDetails[i] = httpsuite.ValidationErrorDetail{
			Field:   validationErr.Field(),
			Message: validationErr.Field() + " failed " + validationErr.Tag() + " validation",
		}
	}

	return &httpsuite.ProblemDetails{
		Type:   v.problems.TypeURL("validation_error"),
		Title:  "Validation Error",
		Status: http.StatusBadRequest,
		Detail: "One or more fields failed validation.",
		Extensions: map[string]interface{}{
			"errors": errorDetails,
		},
	}
}

func mergeProblems(problems *httpsuite.ProblemConfig) httpsuite.ProblemConfig {
	config := httpsuite.DefaultProblemConfig()
	if problems == nil {
		return config
	}

	if problems.BaseURL != "" {
		config.BaseURL = problems.BaseURL
	}
	for key, value := range problems.ErrorTypePaths {
		config.ErrorTypePaths[key] = value
	}
	return config
}
