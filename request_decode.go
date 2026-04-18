package httpsuite

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// BodyDecodeErrorKind identifies the decode failure category.
type BodyDecodeErrorKind string

const (
	BodyDecodeErrorInvalidJSON       BodyDecodeErrorKind = "invalid_json"
	BodyDecodeErrorBodyTooLarge      BodyDecodeErrorKind = "body_too_large"
	BodyDecodeErrorMultipleDocuments BodyDecodeErrorKind = "multiple_documents"
)

// BodyDecodeError represents a request body parsing error.
type BodyDecodeError struct {
	Kind  BodyDecodeErrorKind
	Err   error
	Limit int64
}

func (e *BodyDecodeError) Error() string {
	switch e.Kind {
	case BodyDecodeErrorBodyTooLarge:
		return fmt.Sprintf("request body exceeds the limit of %d bytes", e.Limit)
	case BodyDecodeErrorMultipleDocuments:
		return "request body must contain a single JSON document"
	default:
		if e.Err != nil {
			return e.Err.Error()
		}
		return "invalid request body"
	}
}

func (e *BodyDecodeError) Unwrap() error {
	return e.Err
}

// DecodeRequestBody decodes a JSON request body into T without writing HTTP responses.
func DecodeRequestBody[T any](r *http.Request, maxBodyBytes int64) (T, error) {
	var request T
	if r == nil {
		return request, errNilHTTPRequest
	}
	if r.Body == nil {
		return request, errNilRequestBody
	}
	if r.Body == http.NoBody {
		return request, nil
	}

	limit := maxBodyBytes
	if limit <= 0 {
		limit = defaultMaxBodyBytes
	}

	body := http.MaxBytesReader(nilResponseWriter{}, r.Body, limit)
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&request); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return request, &BodyDecodeError{
				Kind:  BodyDecodeErrorBodyTooLarge,
				Err:   err,
				Limit: maxBytesErr.Limit,
			}
		}

		return request, &BodyDecodeError{
			Kind: BodyDecodeErrorInvalidJSON,
			Err:  err,
		}
	}

	var trailing json.RawMessage
	if err := decoder.Decode(&trailing); err != nil {
		if errors.Is(err, io.EOF) {
			return request, nil
		}

		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return request, &BodyDecodeError{
				Kind:  BodyDecodeErrorBodyTooLarge,
				Err:   err,
				Limit: maxBytesErr.Limit,
			}
		}

		return request, &BodyDecodeError{
			Kind: BodyDecodeErrorMultipleDocuments,
			Err:  err,
		}
	}
	return request, &BodyDecodeError{Kind: BodyDecodeErrorMultipleDocuments}
}

type nilResponseWriter struct{}

func (nilResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (nilResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (nilResponseWriter) WriteHeader(int) {}
