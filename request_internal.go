package httpsuite

import (
	"errors"
	"net/http"
	"reflect"
)

var (
	errNilHTTPRequest     = errors.New("nil http request")
	errNilRequestBody     = errors.New("nil request body")
	errNilParamExtractor  = errors.New("nil param extractor")
	errInvalidRequestType = errors.New("invalid request target")
	errValidationFailed   = errors.New("validation error")
)

func normalizeParseOptions(opts *ParseOptions) ParseOptions {
	normalized := ParseOptions{
		MaxBodyBytes: defaultMaxBodyBytes,
		Problems:     nil,
		Validator:    DefaultValidator(),
	}
	if opts != nil {
		if opts.MaxBodyBytes > 0 {
			normalized.MaxBodyBytes = opts.MaxBodyBytes
		}
		if opts.Problems != nil {
			problems := mergeProblemConfig(opts.Problems)
			normalized.Problems = &problems
		}
		if opts.Validator != nil {
			normalized.Validator = opts.Validator
		}
		normalized.SkipValidation = opts.SkipValidation
	}
	if normalized.Problems == nil {
		problems := DefaultProblemConfig()
		normalized.Problems = &problems
	}
	return normalized
}

func problemFromDecodeError(err error, problems *ProblemConfig) (*ProblemDetails, int) {
	status := http.StatusBadRequest
	var decodeErr *BodyDecodeError
	if errors.As(err, &decodeErr) {
		switch decodeErr.Kind {
		case BodyDecodeErrorBodyTooLarge:
			return NewProblemDetails(
				status,
				problems.TypeURL("bad_request_error"),
				"Request Body Too Large",
				decodeErr.Error(),
			), status
		case BodyDecodeErrorMultipleDocuments:
			return NewProblemDetails(
				status,
				problems.TypeURL("bad_request_error"),
				"Invalid Request",
				"Request body must contain a single JSON document",
			), status
		default:
			return NewProblemDetails(
				status,
				problems.TypeURL("bad_request_error"),
				"Invalid Request",
				decodeErr.Error(),
			), status
		}
	}

	return NewProblemDetails(
		status,
		problems.TypeURL("bad_request_error"),
		"Invalid Request",
		err.Error(),
	), status
}

func problemFromPathParamError(err error, problems *ProblemConfig) (*ProblemDetails, int) {
	status := http.StatusBadRequest
	var pathErr *PathParamError
	if errors.As(err, &pathErr) {
		if pathErr.Missing {
			return NewProblemDetails(
				status,
				problems.TypeURL("bad_request_error"),
				"Missing Parameter",
				"Parameter "+pathErr.Param+" not found in request",
			), status
		}

		problem := NewProblemDetails(
			status,
			problems.TypeURL("bad_request_error"),
			"Invalid Parameter",
			"Failed to bind parameter "+pathErr.Param,
		)
		if pathErr.Err != nil {
			problem.Extensions = map[string]interface{}{"error": pathErr.Err.Error()}
		}
		return problem, status
	}

	return NewProblemDetails(
		status,
		problems.TypeURL("bad_request_error"),
		"Invalid Parameter",
		err.Error(),
	), status
}

func isRequestNil(i interface{}) bool {
	if i == nil {
		return true
	}

	value := reflect.ValueOf(i)
	return value.Kind() == reflect.Ptr && value.IsNil()
}

func ensureRequestInitialized[T any](request T) (T, error) {
	if !isRequestNil(request) {
		return request, nil
	}

	value := reflect.ValueOf(request)
	if !value.IsValid() || value.Kind() != reflect.Ptr {
		var empty T
		return empty, errInvalidRequestType
	}

	elem := value.Type().Elem()
	if elem == nil {
		var empty T
		return empty, errInvalidRequestType
	}

	request = reflect.New(elem).Interface().(T)
	return request, nil
}

func validationProblemStatus(problem *ProblemDetails) int {
	if problem == nil {
		return http.StatusBadRequest
	}
	if problem.Status >= 400 && problem.Status <= 599 {
		return problem.Status
	}
	return http.StatusBadRequest
}
