package httpsuite

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseRequestPrecognition(t *testing.T) {
	ClearValidator()
	t.Cleanup(ClearValidator)

	invalid := stubValidator{problem: NewProblemDetails(http.StatusBadRequest, "", "Validation Error", "name is required")}

	tests := []struct {
		name         string
		body         string
		path         string
		header       bool
		opts         *ParseOptions
		wantErr      error
		wantStatus   int
		wantResponse bool
		wantVary     bool
	}{
		{
			name:       "valid precognitive request",
			body:       `{"name":"Test"}`,
			path:       "/test/1",
			header:     true,
			opts:       &ParseOptions{Precognition: true},
			wantErr:    ErrPrecognitionHandled,
			wantStatus: http.StatusNoContent,
			wantVary:   true,
		},
		{
			name:       "validation failure returns 422",
			body:       `{"name":""}`,
			path:       "/test/1",
			header:     true,
			opts:       &ParseOptions{Precognition: true, Validator: invalid},
			wantErr:    errValidationFailed,
			wantStatus: http.StatusUnprocessableEntity,
			wantVary:   true,
		},
		{
			name:       "malformed json stays a 400 parse error",
			body:       `{bad`,
			path:       "/test/1",
			header:     true,
			opts:       &ParseOptions{Precognition: true},
			wantStatus: http.StatusBadRequest,
			wantVary:   true,
		},
		{
			name:       "oversized body stays 413",
			body:       `{"name":"` + strings.Repeat("a", 64) + `"}`,
			path:       "/test/1",
			header:     true,
			opts:       &ParseOptions{Precognition: true, MaxBodyBytes: 8},
			wantStatus: http.StatusRequestEntityTooLarge,
			wantVary:   true,
		},
		{
			name:       "invalid path param stays 400",
			body:       `{"name":"Test"}`,
			path:       "/test/nope",
			header:     true,
			opts:       &ParseOptions{Precognition: true},
			wantStatus: http.StatusBadRequest,
			wantVary:   true,
		},
		{
			name:         "ordinary request on opted-in endpoint",
			body:         `{"name":"Test"}`,
			path:         "/test/1",
			opts:         &ParseOptions{Precognition: true},
			wantResponse: true,
			wantVary:     true,
		},
		{
			name:         "header ignored when not opted in",
			body:         `{"name":"Test"}`,
			path:         "/test/1",
			header:       true,
			opts:         nil,
			wantResponse: true,
		},
		{
			name:       "validation status untouched when not opted in",
			body:       `{"name":""}`,
			path:       "/test/1",
			header:     true,
			opts:       &ParseOptions{Validator: invalid},
			wantErr:    errValidationFailed,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			if tt.header {
				r.Header.Set(PrecognitionHeader, "true")
			}
			w := httptest.NewRecorder()

			got, err := ParseRequest[*testRequest](w, r, testParamExtractor, tt.opts, "id")

			if tt.wantResponse {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got == nil || got.ID != 1 || got.Name != "Test" {
					t.Fatalf("unexpected request: %#v", got)
				}
				if w.Body.Len() != 0 {
					t.Fatalf("expected nothing written, got %q", w.Body.String())
				}
			} else {
				if err == nil {
					t.Fatal("expected error")
				}
				if got != nil {
					t.Fatalf("expected zero request, got %#v", got)
				}
				if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				if w.Code != tt.wantStatus {
					t.Fatalf("expected status %d, got %d", tt.wantStatus, w.Code)
				}
			}

			if IsPrecognitionHandled(err) != errors.Is(tt.wantErr, ErrPrecognitionHandled) {
				t.Fatalf("IsPrecognitionHandled mismatch for %v", err)
			}
			if vary := w.Header().Get("Vary") == PrecognitionHeader; vary != tt.wantVary {
				t.Fatalf("Vary header = %q, want precognition vary %v", w.Header().Get("Vary"), tt.wantVary)
			}
		})
	}
}

func TestParseRequestPrecognitionResponses(t *testing.T) {
	ClearValidator()
	t.Cleanup(ClearValidator)

	t.Run("success has empty body and success header", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/test/1", strings.NewReader(`{"name":"Test"}`))
		r.Header.Set(PrecognitionHeader, "true")
		w := httptest.NewRecorder()

		_, err := ParseRequest[*testRequest](w, r, testParamExtractor, &ParseOptions{Precognition: true}, "id")
		if !IsPrecognitionHandled(err) {
			t.Fatalf("expected ErrPrecognitionHandled, got %v", err)
		}
		if w.Body.Len() != 0 {
			t.Fatalf("expected empty body, got %q", w.Body.String())
		}
		if w.Header().Get(PrecognitionSuccessHeader) != "true" {
			t.Fatalf("expected %s header", PrecognitionSuccessHeader)
		}
	})

	t.Run("failure keeps problem body and extensions", func(t *testing.T) {
		problem := NewProblemDetails(http.StatusBadRequest, "", "Validation Error", "invalid")
		problem.Extensions = map[string]interface{}{
			"errors": []ValidationErrorDetail{{Field: "name", Message: "required"}},
		}
		r := httptest.NewRequest(http.MethodPost, "/test/1", strings.NewReader(`{"name":""}`))
		r.Header.Set(PrecognitionHeader, "true")
		w := httptest.NewRecorder()

		_, _ = ParseRequest[*testRequest](w, r, testParamExtractor, &ParseOptions{
			Precognition: true,
			Validator:    stubValidator{problem: problem},
		}, "id")

		var body map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body["status"] != float64(http.StatusUnprocessableEntity) {
			t.Fatalf("expected status 422 in body, got %v", body["status"])
		}
		if _, ok := body["errors"]; !ok {
			t.Fatalf("expected errors extension, got %v", body)
		}
		if problem.Status != http.StatusBadRequest {
			t.Fatal("validator problem must not be mutated")
		}
	})
}

func TestIsPrecognitionRequest(t *testing.T) {
	t.Parallel()

	for value, want := range map[string]bool{"true": true, "TRUE": true, " true ": true, "false": false, "": false, "1": false} {
		r := httptest.NewRequest(http.MethodPost, "/", nil)
		r.Header.Set(PrecognitionHeader, value)
		if got := IsPrecognitionRequest(r); got != want {
			t.Errorf("IsPrecognitionRequest(%q) = %v, want %v", value, got, want)
		}
	}
	if IsPrecognitionRequest(nil) {
		t.Error("nil request must not be precognitive")
	}
}
