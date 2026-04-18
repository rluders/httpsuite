package httpsuite

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testResponse struct {
	Key string `json:"key"`
}

func TestSendResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		code         int
		data         any
		problem      *ProblemDetails
		meta         any
		expectedCode int
		contentType  string
	}{
		{
			name:         "success response",
			code:         http.StatusOK,
			data:         &testResponse{Key: "value"},
			expectedCode: http.StatusOK,
			contentType:  "application/json; charset=utf-8",
		},
		{
			name:         "problem response",
			code:         http.StatusNotFound,
			problem:      NewProblemDetails(http.StatusNotFound, "", "Not Found", "The requested resource was not found"),
			expectedCode: http.StatusNotFound,
			contentType:  "application/problem+json; charset=utf-8",
		},
		{
			name:         "success response with page meta",
			code:         http.StatusOK,
			data:         []string{"a", "b"},
			meta:         NewPageMeta(2, 10, 25),
			expectedCode: http.StatusOK,
			contentType:  "application/json; charset=utf-8",
		},
		{
			name:         "success response with cursor meta",
			code:         http.StatusOK,
			data:         []string{"a", "b"},
			meta:         NewCursorMeta("next-1", "prev-1", true, true),
			expectedCode: http.StatusOK,
			contentType:  "application/json; charset=utf-8",
		},
		{
			name:         "problem response normalizes status",
			code:         http.StatusBadRequest,
			problem:      NewProblemDetails(http.StatusUnprocessableEntity, "", "Invalid Request", "invalid"),
			expectedCode: http.StatusBadRequest,
			contentType:  "application/problem+json; charset=utf-8",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			SendResponse[any](w, tt.code, tt.data, tt.problem, tt.meta)

			if w.Code != tt.expectedCode {
				t.Fatalf("expected status %d, got %d", tt.expectedCode, w.Code)
			}
			if got := w.Header().Get("Content-Type"); got != tt.contentType {
				t.Fatalf("expected content type %q, got %q", tt.contentType, got)
			}
			if !json.Valid(w.Body.Bytes()) {
				t.Fatalf("expected valid json, got %q", w.Body.String())
			}

			if tt.problem != nil {
				var problem ProblemDetails
				if err := json.NewDecoder(w.Body).Decode(&problem); err != nil {
					t.Fatalf("decode problem details: %v", err)
				}
				if problem.Status != tt.expectedCode {
					t.Fatalf("expected problem status %d, got %d", tt.expectedCode, problem.Status)
				}
			}
		})
	}
}
