package httpsuite

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"reflect"
)

// RequestParamSetter defines the interface used to set the parameters to the HTTP request object by the request parser.
// Implementing this interface allows custom handling of URL parameters.
type RequestParamSetter interface {
	// SetParam assigns a value to a specified field in the request struct.
	// The fieldName parameter is the name of the field, and value is the value to set.
	SetParam(fieldName, value string) error
}

// ParseRequest parses the incoming HTTP request into a specified struct type, handling JSON decoding and URL parameters.
// It validates the parsed request and returns it along with any potential errors.
// The pathParams variadic argument allows specifying URL parameters to be extracted.
// If an error occurs during parsing, validation, or parameter setting, it responds with an appropriate HTTP status.
func ParseRequest[T RequestParamSetter](w http.ResponseWriter, r *http.Request, pathParams ...string) (T, error) {
	var request T
	var empty T

	defer func() {
		_ = r.Body.Close()
	}()

	if r.Body != http.NoBody {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			SendResponse[any](w, "Invalid JSON format", http.StatusBadRequest, nil)
			return empty, err
		}
	}

	// If body wasn't parsed request may be nil and cause problems ahead
	if isRequestNil(request) {
		request = reflect.New(reflect.TypeOf(request).Elem()).Interface().(T)
	}

	// Parse URL parameters
	for _, key := range pathParams {
		value := chi.URLParam(r, key)
		if value == "" {
			SendResponse[any](w, "Parameter "+key+" not found in request", http.StatusBadRequest, nil)
			return empty, errors.New("missing parameter: " + key)
		}

		if err := request.SetParam(key, value); err != nil {
			SendResponse[any](w, "Failed to set field "+key, http.StatusInternalServerError, nil)
			return empty, err
		}
	}

	// Validate the combined request struct
	if validationErr := IsRequestValid(request); validationErr != nil {
		SendResponse[ValidationErrors](w, "Validation error", http.StatusBadRequest, validationErr)
		return empty, errors.New("validation error")
	}

	return request, nil
}

func isRequestNil(i interface{}) bool {
	return i == nil || (reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil())
}
