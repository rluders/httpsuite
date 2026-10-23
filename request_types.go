package httpsuite

import "net/http"

// RequestParamSetter defines custom path parameter binding for request structs.
type RequestParamSetter interface {
	SetParam(fieldName, value string) error
}

// ParamExtractor extracts a path parameter from a request.
type ParamExtractor func(r *http.Request, key string) string

// Validator validates request payloads without coupling the core package to a validation library.
type Validator interface {
	Validate(any) *ProblemDetails
}

// ParseOptions configures request parsing behavior.
type ParseOptions struct {
	MaxBodyBytes   int64
	Problems       *ProblemConfig
	Validator      Validator
	SkipValidation bool
	// Precognition enables validation-only processing for requests that send
	// the Precognition: true header. Such requests are decoded, bound and
	// validated as usual; on success ParseRequest writes 204 No Content and
	// returns ErrPrecognitionHandled, on validation failure it writes the
	// validator's problem with status 422. When false, the header is ignored.
	Precognition bool
}

const defaultMaxBodyBytes int64 = 1 << 20
