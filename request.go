package httpsuite

import (
	"encoding/json"
	"errors"
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

// ParamExtractor is a function type that extracts a URL parameter from the incoming HTTP request.
// It takes the `http.Request` and a `key` as arguments, and returns the value of the URL parameter
// as a string. This function allows flexibility for extracting parameters from different routers,
// such as Chi, Echo, Gorilla Mux, or the default Go router.
//
// Example usage:
//
//	paramExtractor := func(r *http.Request, key string) string {
//	    return r.URL.Query().Get(key)
//	}
type ParamExtractor func(r *http.Request, key string) string

// ParseRequest parses the incoming HTTP request into a specified struct type,
// handling JSON decoding and extracting URL parameters using the provided `paramExtractor` function.
// The `paramExtractor` allows flexibility to integrate with various routers (e.g., Chi, Echo, Gorilla Mux).
// It extracts the specified parameters from the URL and sets them on the struct.
//
// The `pathParams` variadic argument is used to specify which URL parameters to extract and set on the struct.
//
// The function also validates the parsed request. If the request fails validation or if any error occurs during
// JSON parsing or parameter extraction, it responds with an appropriate HTTP status and error message.
//
// Parameters:
//   - `w`: The `http.ResponseWriter` used to send the response to the client.
//   - `r`: The incoming HTTP request to be parsed.
//   - `paramExtractor`: A function that extracts URL parameters from the request. This function allows custom handling
//     of parameters based on the router being used.
//   - `pathParams`: A variadic argument specifying which URL parameters to extract and set on the struct.
//
// Returns:
//   - A parsed struct of the specified type `T`, if successful.
//   - An error, if parsing, validation, or parameter extraction fails.
//
// Example usage:
//
//	request, err := ParseRequest[MyRequestType](w, r, MyParamExtractor, "id", "name")
//	if err != nil {
//	    // Handle error
//	}
//
//	// Continue processing the valid request...
func ParseRequest[T RequestParamSetter](w http.ResponseWriter, r *http.Request, paramExtractor ParamExtractor, pathParams ...string) (T, error) {
	var request T
	var empty T
	defer func() { _ = r.Body.Close() }()

	// Decode JSON body if present
	if r.Body != http.NoBody {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			problem := NewProblemDetails(http.StatusBadRequest, "Invalid Request", err.Error())
			SendResponse[any](w, http.StatusBadRequest, nil, problem, nil)
			return empty, err
		}
	}

	// Ensure request object is properly initialized
	if isRequestNil(request) {
		request = reflect.New(reflect.TypeOf(request).Elem()).Interface().(T)
	}

	// Extract and set URL parameters
	for _, key := range pathParams {
		value := paramExtractor(r, key)
		if value == "" {
			problem := NewProblemDetails(http.StatusBadRequest, "Missing Parameter", "Parameter "+key+" not found in request")
			SendResponse[any](w, http.StatusBadRequest, nil, problem, nil)
			return empty, errors.New("missing parameter: " + key)
		}
		if err := request.SetParam(key, value); err != nil {
			problem := NewProblemDetails(http.StatusInternalServerError, "Parameter Error", "Failed to set field "+key)
			problem.Extensions = map[string]interface{}{"error": err.Error()}
			SendResponse[any](w, http.StatusInternalServerError, nil, problem, nil)
			return empty, err
		}
	}

	// Validate the request
	if validationErr := IsRequestValid(request); validationErr != nil {
		SendResponse[any](w, http.StatusBadRequest, nil, validationErr, nil)
		return empty, errors.New("validation error")
	}

	return request, nil
}

// isRequestNil checks if a request object is nil or an uninitialized pointer.
func isRequestNil(i interface{}) bool {
	return i == nil || (reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil())
}
