package httpsuite

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestResponse struct {
	Key string `json:"key"`
}

func Test_SendResponse(t *testing.T) {
	tests := []struct {
		name         string
		code         int
		data         any
		problem      *ProblemDetails
		meta         *Meta
		expectedCode int
		expectedJSON string
	}{
		{
			name:         "200 OK with TestResponse body",
			code:         http.StatusOK,
			data:         &TestResponse{Key: "value"},
			expectedCode: http.StatusOK,
			expectedJSON: `{
				"data": {
					"key": "value"
				}
			}`,
		},
		{
			name:         "404 Not Found without body",
			code:         http.StatusNotFound,
			problem:      NewProblemDetails(http.StatusNotFound, "", "Not Found", "The requested resource was not found"),
			expectedCode: http.StatusNotFound,
			expectedJSON: `{
				"type": "about:blank",
				"title": "Not Found",
				"status": 404,
				"detail": "The requested resource was not found"
			}`,
		},
		{
			name:         "200 OK with pagination metadata",
			code:         http.StatusOK,
			data:         &TestResponse{Key: "value"},
			meta:         &Meta{TotalPages: 100, Page: 1, PageSize: 10},
			expectedCode: http.StatusOK,
			expectedJSON: `{
				"data": {
					"key": "value"
				},
				"meta": {
					"total_pages": 100,
					"page": 1,
					"page_size": 10
				}
			}`,
		},
		{
			name: "400 Bad Request with validation error",
			code: http.StatusBadRequest,
			problem: &ProblemDetails{
				Type:   "https://example.com/validation-error",
				Title:  "Validation Error",
				Status: http.StatusBadRequest,
				Detail: "One or more fields failed validation.",
				Extensions: map[string]interface{}{
					"errors": []ValidationErrorDetail{
						{Field: "email", Message: "Email is required"},
						{Field: "password", Message: "Password is required"},
					},
				},
			},
			expectedCode: http.StatusBadRequest,
			expectedJSON: `{
				"type": "https://example.com/validation-error",
				"title": "Validation Error",
				"status": 400,
				"detail": "One or more fields failed validation.",
				"extensions": {
					"errors": [
						{"field": "email", "message": "Email is required"},
						{"field": "password", "message": "Password is required"}
					]
				}
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			// Call SendResponse with the appropriate data or problem
			SendResponse[any](w, tt.code, tt.data, tt.problem, tt.meta)

			// Assert response status code and content type
			assert.Equal(t, tt.expectedCode, w.Code)
			if w.Code >= 400 {
				assert.Equal(t, "application/problem+json; charset=utf-8", w.Header().Get("Content-Type"))
			} else {
				assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
			}
			// Assert response body
			assert.JSONEq(t, tt.expectedJSON, w.Body.String())
		})
	}
}
