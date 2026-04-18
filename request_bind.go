package httpsuite

import (
	"errors"
	"fmt"
	"net/http"
)

// PathParamError represents a path parameter binding error.
type PathParamError struct {
	Param   string
	Missing bool
	Err     error
}

func (e *PathParamError) Error() string {
	if e.Missing {
		return "missing parameter: " + e.Param
	}
	if e.Err != nil {
		return fmt.Sprintf("invalid parameter %s: %v", e.Param, e.Err)
	}
	return "invalid parameter: " + e.Param
}

func (e *PathParamError) Unwrap() error {
	return e.Err
}

// BindPathParams applies extracted path params to a request object without writing HTTP responses.
func BindPathParams[T any](request T, r *http.Request, paramExtractor ParamExtractor, pathParams ...string) (T, error) {
	if len(pathParams) == 0 {
		return request, nil
	}
	if r == nil {
		var empty T
		return empty, errNilHTTPRequest
	}
	if paramExtractor == nil {
		var empty T
		return empty, errNilParamExtractor
	}

	var err error
	request, err = ensureRequestInitialized(request)
	if err != nil {
		var empty T
		return empty, err
	}

	setter, ok := any(request).(RequestParamSetter)
	if !ok {
		var empty T
		return empty, errors.Join(errInvalidRequestType, errors.New("request type does not implement RequestParamSetter"))
	}

	for _, key := range pathParams {
		value := paramExtractor(r, key)
		if value == "" {
			var empty T
			return empty, &PathParamError{
				Param:   key,
				Missing: true,
			}
		}

		if err := setter.SetParam(key, value); err != nil {
			var empty T
			return empty, &PathParamError{
				Param: key,
				Err:   err,
			}
		}
	}

	return request, nil
}
