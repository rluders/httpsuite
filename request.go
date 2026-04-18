package httpsuite

import (
	"errors"
	"net/http"
)

// ParseRequest parses the incoming HTTP request into a specified struct type,
// handling JSON decoding, request body limits, path parameter binding, and
// optional validation. Invalid inputs return regular errors instead of panicking.
func ParseRequest[T any](w http.ResponseWriter, r *http.Request, paramExtractor ParamExtractor, opts *ParseOptions, pathParams ...string) (T, error) {
	var empty T
	if r == nil {
		return empty, errNilHTTPRequest
	}
	if r.Body != nil {
		defer func() { _ = r.Body.Close() }()
	}

	options := normalizeParseOptions(opts)

	request, err := DecodeRequestBody[T](r, options.MaxBodyBytes)
	if err != nil {
		var decodeErr *BodyDecodeError
		if !errors.As(err, &decodeErr) {
			return empty, err
		}
		problem, status := problemFromDecodeError(err, options.Problems)
		SendResponse[any](w, status, nil, problem, nil)
		return empty, err
	}

	request, err = BindPathParams(request, r, paramExtractor, pathParams...)
	if err != nil {
		var pathErr *PathParamError
		if !errors.As(err, &pathErr) {
			return empty, err
		}
		problem, status := problemFromPathParamError(err, options.Problems)
		SendResponse[any](w, status, nil, problem, nil)
		return empty, err
	}

	if !options.SkipValidation {
		if problem := ValidateRequest(request, options.Validator); problem != nil {
			SendResponse[any](w, validationProblemStatus(problem), nil, problem, nil)
			return empty, errValidationFailed
		}
	}

	return request, nil
}
