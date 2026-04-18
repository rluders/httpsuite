package httpsuite

import "sync"

var (
	defaultValidatorMu sync.RWMutex
	defaultValidator   Validator
)

// ValidateRequest applies a validator without writing HTTP responses.
func ValidateRequest(request any, validator Validator) *ProblemDetails {
	if validator == nil {
		return nil
	}
	return validator.Validate(request)
}

// SetValidator configures the package-level default validator used by ParseRequest.
func SetValidator(v Validator) {
	defaultValidatorMu.Lock()
	defer defaultValidatorMu.Unlock()
	defaultValidator = v
}

// ClearValidator removes the package-level default validator.
func ClearValidator() {
	SetValidator(nil)
}

// DefaultValidator returns the current package-level default validator.
func DefaultValidator() Validator {
	defaultValidatorMu.RLock()
	defer defaultValidatorMu.RUnlock()
	return defaultValidator
}
