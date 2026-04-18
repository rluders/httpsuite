package httpsuite

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBindPathParams(t *testing.T) {
	t.Parallel()

	t.Run("nil request pointer target", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/123", nil)
		got, err := BindPathParams[*testRequest](nil, req, testParamExtractor, "id")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil || got.ID != 123 {
			t.Fatalf("expected initialized request target, got %#v", got)
		}
	})

	t.Run("unsupported request target", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/123", nil)
		_, err := BindPathParams[bodyOnlyRequest](bodyOnlyRequest{}, req, testParamExtractor, "id")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if errors.Is(err, errNilHTTPRequest) || errors.Is(err, errNilParamExtractor) {
			t.Fatalf("expected unsupported target error, got %v", err)
		}
	})

	t.Run("nil extractor", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/123", nil)
		_, err := BindPathParams[*testRequest](&testRequest{}, req, nil, "id")
		if !errors.Is(err, errNilParamExtractor) {
			t.Fatalf("expected nil extractor error, got %v", err)
		}
	})

	t.Run("valid parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/123", nil)
		got, err := BindPathParams[*testRequest](&testRequest{Name: "ok"}, req, testParamExtractor, "id")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != 123 {
			t.Fatalf("expected id 123, got %d", got.ID)
		}
	})

	t.Run("missing parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		_, err := BindPathParams[*testRequest](&testRequest{}, req, testParamExtractor, "id")
		var pathErr *PathParamError
		if !errors.As(err, &pathErr) {
			t.Fatalf("expected PathParamError, got %v", err)
		}
		if !pathErr.Missing {
			t.Fatalf("expected missing parameter error")
		}
	})

	t.Run("invalid parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/nope", nil)
		_, err := BindPathParams[*testRequest](&testRequest{}, req, testParamExtractor, "id")
		var pathErr *PathParamError
		if !errors.As(err, &pathErr) {
			t.Fatalf("expected PathParamError, got %v", err)
		}
		if pathErr.Missing {
			t.Fatalf("expected invalid parameter, got missing")
		}
	})
}
