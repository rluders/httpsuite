package httpsuite

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseRequest(t *testing.T) {
	ClearValidator()
	t.Cleanup(ClearValidator)

	tests := []struct {
		name               string
		body               string
		path               string
		pathParams         []string
		opts               *ParseOptions
		want               *testRequest
		wantErr            bool
		wantStatus         int
		wantTitle          string
		wantDetailContains string
	}{
		{
			name:       "successful request",
			body:       `{"name":"Test"}`,
			path:       "/test/123",
			pathParams: []string{"id"},
			want:       &testRequest{ID: 123, Name: "Test"},
		},
		{
			name:       "body only request",
			body:       `{"id":42,"name":"OnlyBody"}`,
			path:       "/test",
			pathParams: nil,
			want:       &testRequest{ID: 42, Name: "OnlyBody"},
		},
		{
			name:               "invalid json body",
			body:               `{invalid-json}`,
			path:               "/test/123",
			pathParams:         []string{"id"},
			wantErr:            true,
			wantStatus:         http.StatusBadRequest,
			wantTitle:          "Invalid Request",
			wantDetailContains: "invalid character",
		},
		{
			name:               "multiple json documents",
			body:               `{"name":"Test"}{"name":"Again"}`,
			path:               "/test/123",
			pathParams:         []string{"id"},
			wantErr:            true,
			wantStatus:         http.StatusBadRequest,
			wantTitle:          "Invalid Request",
			wantDetailContains: "single JSON document",
		},
		{
			name:               "missing parameter",
			body:               `{"name":"Test"}`,
			path:               "/test",
			pathParams:         []string{"id"},
			wantErr:            true,
			wantStatus:         http.StatusBadRequest,
			wantTitle:          "Missing Parameter",
			wantDetailContains: "Parameter id not found",
		},
		{
			name:               "invalid parameter",
			body:               `{"name":"Test"}`,
			path:               "/test/nope",
			pathParams:         []string{"id"},
			wantErr:            true,
			wantStatus:         http.StatusBadRequest,
			wantTitle:          "Invalid Parameter",
			wantDetailContains: "Failed to bind parameter id",
		},
		{
			name:               "body exceeds configured limit",
			body:               `{"name":"TooLarge"}`,
			path:               "/test/123",
			pathParams:         []string{"id"},
			opts:               &ParseOptions{MaxBodyBytes: 8},
			wantErr:            true,
			wantStatus:         http.StatusBadRequest,
			wantTitle:          "Request Body Too Large",
			wantDetailContains: "exceeds the limit",
		},
		{
			name:       "custom problem config",
			body:       `{"name":"Test"}`,
			path:       "/test/123",
			pathParams: []string{"id"},
			opts: &ParseOptions{
				Problems: &ProblemConfig{
					BaseURL: "https://api.example.com",
					ErrorTypePaths: map[string]string{
						"bad_request_error": "/errors/bad-request",
						"server_error":      "/errors/server-error",
					},
				},
			},
			want: &testRequest{ID: 123, Name: "Test"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var body *bytes.Buffer
			if tt.body != "" {
				body = bytes.NewBufferString(tt.body)
			} else {
				body = bytes.NewBuffer(nil)
			}

			req := httptest.NewRequest(http.MethodPost, tt.path, body)
			w := httptest.NewRecorder()

			got, err := ParseRequest[*testRequest](w, req, testParamExtractor, tt.opts, tt.pathParams...)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}

				if w.Code != tt.wantStatus {
					t.Fatalf("expected status %d, got %d", tt.wantStatus, w.Code)
				}

				var problem ProblemDetails
				if decodeErr := json.NewDecoder(w.Body).Decode(&problem); decodeErr != nil {
					t.Fatalf("decode problem details: %v", decodeErr)
				}
				if problem.Title != tt.wantTitle {
					t.Fatalf("expected title %q, got %q", tt.wantTitle, problem.Title)
				}
				if !strings.Contains(problem.Detail, tt.wantDetailContains) {
					t.Fatalf("expected detail %q to contain %q", problem.Detail, tt.wantDetailContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got == nil {
				t.Fatalf("expected request, got nil")
			}
			if *got != *tt.want {
				t.Fatalf("expected %+v, got %+v", *tt.want, *got)
			}
		})
	}
}

func TestParseRequestWithoutRequestParamSetter(t *testing.T) {
	ClearValidator()
	t.Cleanup(ClearValidator)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{"name":"Ada","age":36}`))
	w := httptest.NewRecorder()

	got, err := ParseRequest[*bodyOnlyRequest](w, req, testParamExtractor, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected parsed request, got nil")
	}
	if got.Name != "Ada" || got.Age != 36 {
		t.Fatalf("unexpected parsed request: %#v", got)
	}
}

func TestParseRequestInvalidInputs(t *testing.T) {
	ClearValidator()
	t.Cleanup(ClearValidator)

	tests := []struct {
		name       string
		makeReq    func() *http.Request
		extractor  ParamExtractor
		pathParams []string
		wantErr    error
	}{
		{
			name:      "nil request",
			makeReq:   func() *http.Request { return nil },
			extractor: testParamExtractor,
			wantErr:   errNilHTTPRequest,
		},
		{
			name: "nil body",
			makeReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/test/123", nil)
				req.Body = nil
				return req
			},
			extractor: testParamExtractor,
			wantErr:   errNilRequestBody,
		},
		{
			name: "nil extractor with params",
			makeReq: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/test/123", bytes.NewBufferString(`{"name":"Test"}`))
			},
			pathParams: []string{"id"},
			wantErr:    errNilParamExtractor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			got, err := ParseRequest[*testRequest](w, tt.makeReq(), tt.extractor, nil, tt.pathParams...)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if got != nil {
				t.Fatalf("expected nil request, got %#v", got)
			}
			if w.Body.Len() != 0 {
				t.Fatalf("expected no response body to be written, got %q", w.Body.String())
			}
		})
	}
}

func TestParseRequestUsesDefaultValidator(t *testing.T) {
	ClearValidator()
	t.Cleanup(ClearValidator)

	problem := &ProblemDetails{
		Type:   GetProblemTypeURL("validation_error"),
		Title:  "Validation Error",
		Status: http.StatusBadRequest,
		Detail: "One or more fields failed validation.",
	}
	SetValidator(stubValidator{problem: problem})

	req := httptest.NewRequest(http.MethodPost, "/test/123", bytes.NewBufferString(`{"name":""}`))
	w := httptest.NewRecorder()

	got, err := ParseRequest[*testRequest](w, req, testParamExtractor, nil, "id")
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}
	if got != nil {
		t.Fatalf("expected nil request, got %#v", got)
	}
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestParseRequestValidatorOverride(t *testing.T) {
	ClearValidator()
	t.Cleanup(ClearValidator)

	SetValidator(stubValidator{
		problem: &ProblemDetails{
			Type:   GetProblemTypeURL("validation_error"),
			Title:  "Validation Error",
			Status: http.StatusBadRequest,
			Detail: "global validator failed",
		},
	})

	override := stubValidator{
		problem: &ProblemDetails{
			Type:   GetProblemTypeURL("validation_error"),
			Title:  "Validation Error",
			Status: http.StatusBadRequest,
			Detail: "override validator failed",
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/test/123", bytes.NewBufferString(`{"name":"ok"}`))
	w := httptest.NewRecorder()

	_, err := ParseRequest[*testRequest](w, req, testParamExtractor, &ParseOptions{Validator: override}, "id")
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	var problem ProblemDetails
	if decodeErr := json.NewDecoder(w.Body).Decode(&problem); decodeErr != nil {
		t.Fatalf("decode problem details: %v", decodeErr)
	}
	if problem.Detail != "override validator failed" {
		t.Fatalf("expected override validator detail, got %q", problem.Detail)
	}
}

func TestParseRequestValidationStatus(t *testing.T) {
	ClearValidator()
	t.Cleanup(ClearValidator)

	tests := []struct {
		name       string
		problem    *ProblemDetails
		wantStatus int
	}{
		{
			name: "custom status preserved",
			problem: &ProblemDetails{
				Type:   GetProblemTypeURL("validation_error"),
				Title:  "Validation Error",
				Status: http.StatusUnprocessableEntity,
				Detail: "unprocessable payload",
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid status falls back to bad request",
			problem: &ProblemDetails{
				Type:   GetProblemTypeURL("validation_error"),
				Title:  "Validation Error",
				Status: 0,
				Detail: "bad payload",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetValidator(stubValidator{problem: tt.problem})

			req := httptest.NewRequest(http.MethodPost, "/test/123", bytes.NewBufferString(`{"name":"ok"}`))
			w := httptest.NewRecorder()

			_, err := ParseRequest[*testRequest](w, req, testParamExtractor, nil, "id")
			if !errors.Is(err, errValidationFailed) {
				t.Fatalf("expected validation error, got %v", err)
			}
			if w.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestParseRequestSkipValidation(t *testing.T) {
	ClearValidator()
	t.Cleanup(ClearValidator)

	SetValidator(stubValidator{
		problem: &ProblemDetails{
			Type:   GetProblemTypeURL("validation_error"),
			Title:  "Validation Error",
			Status: http.StatusBadRequest,
			Detail: "global validator failed",
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/test/123", bytes.NewBufferString(`{"name":"ok"}`))
	w := httptest.NewRecorder()

	got, err := ParseRequest[*testRequest](w, req, testParamExtractor, &ParseOptions{SkipValidation: true}, "id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != 123 {
		t.Fatalf("expected parsed request, got %#v", got)
	}
}
