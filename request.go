package httpsuite

import (
	"errors"
	"net/http"
)

// ParseRequest parses the incoming HTTP request into a specified struct type,
// handling JSON decoding, request body limits, path parameter binding, and
// optional validation. Invalid inputs return regular errors instead of panicking.
// With ParseOptions.Precognition enabled, a request carrying the Precognition
// header stops after validation and ParseRequest returns ErrPrecognitionHandled.
func ParseRequest[T any](w http.ResponseWriter, r *http.Request, paramExtractor ParamExtractor, opts *ParseOptions, pathParams ...string) (T, error) {
	var empty T
	if r == nil {
		return empty, errNilHTTPRequest
	}
	if r.Body != nil {
		defer func() { _ = r.Body.Close() }()
	}

	options := normalizeParseOptions(opts)

	precognitive := false
	if options.Precognition {
		// Responses differ by the Precognition header on opted-in endpoints.
		w.Header().Add("Vary", PrecognitionHeader)
		precognitive = IsPrecognitionRequest(r)
	}

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
			status := validationProblemStatus(problem)
			if precognitive {
				status = http.StatusUnprocessableEntity
			}
			SendResponse[any](w, status, nil, problem, nil)
			return empty, errValidationFailed
		}
	}

	if precognitive {
		writePrecognitionSuccess(w)
		return empty, ErrPrecognitionHandled
	}

	return request, nil
}
