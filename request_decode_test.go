package httpsuite

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecodeRequestBody(t *testing.T) {
	t.Parallel()

	t.Run("nil request", func(t *testing.T) {
		_, err := DecodeRequestBody[*testRequest](nil, defaultMaxBodyBytes)
		if !errors.Is(err, errNilHTTPRequest) {
			t.Fatalf("expected nil request error, got %v", err)
		}
	})

	t.Run("nil body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Body = nil

		_, err := DecodeRequestBody[*testRequest](req, defaultMaxBodyBytes)
		if !errors.Is(err, errNilRequestBody) {
			t.Fatalf("expected nil body error, got %v", err)
		}
	})

	t.Run("valid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{"id":42,"name":"OnlyBody"}`))
		got, err := DecodeRequestBody[*testRequest](req, defaultMaxBodyBytes)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil || got.ID != 42 || got.Name != "OnlyBody" {
			t.Fatalf("unexpected decoded request: %#v", got)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{invalid-json}`))
		_, err := DecodeRequestBody[*testRequest](req, defaultMaxBodyBytes)
		var decodeErr *BodyDecodeError
		if !errors.As(err, &decodeErr) {
			t.Fatalf("expected BodyDecodeError, got %v", err)
		}
		if decodeErr.Kind != BodyDecodeErrorInvalidJSON {
			t.Fatalf("expected invalid json error, got %s", decodeErr.Kind)
		}
	})

	t.Run("body too large", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{"name":"TooLarge"}`))
		_, err := DecodeRequestBody[*testRequest](req, 8)
		var decodeErr *BodyDecodeError
		if !errors.As(err, &decodeErr) {
			t.Fatalf("expected BodyDecodeError, got %v", err)
		}
		if decodeErr.Kind != BodyDecodeErrorBodyTooLarge {
			t.Fatalf("expected body too large error, got %s", decodeErr.Kind)
		}
	})

	t.Run("multiple json documents", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{"id":1}{"id":2}`))
		_, err := DecodeRequestBody[*testRequest](req, defaultMaxBodyBytes)
		var decodeErr *BodyDecodeError
		if !errors.As(err, &decodeErr) {
			t.Fatalf("expected BodyDecodeError, got %v", err)
		}
		if decodeErr.Kind != BodyDecodeErrorMultipleDocuments {
			t.Fatalf("expected multiple documents error, got %s", decodeErr.Kind)
		}
	})

	t.Run("trailing decode body too large", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{} {"name":"trailing"}`))
		_, err := DecodeRequestBody[*testRequest](req, 3)
		var decodeErr *BodyDecodeError
		if !errors.As(err, &decodeErr) {
			t.Fatalf("expected BodyDecodeError, got %v", err)
		}
		if decodeErr.Kind != BodyDecodeErrorBodyTooLarge {
			t.Fatalf("expected body too large error, got %s", decodeErr.Kind)
		}
	})
}

func BenchmarkParseRequestBody(b *testing.B) {
	body := []byte(`{"id":42,"name":"OnlyBody"}`)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		if _, err := DecodeRequestBody[*testRequest](req, defaultMaxBodyBytes); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
