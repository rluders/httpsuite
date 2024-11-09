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

func TestResponse_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		response Response[any]
		expected string
	}{
		{
			name:     "Basic Response",
			response: Response[any]{Code: 200, Message: "OK"},
			expected: `{"code":200,"message":"OK"}`,
		},
		{
			name:     "Response with Body",
			response: Response[any]{Code: 201, Message: "Created", Body: map[string]string{"id": "123"}},
			expected: `{"code":201,"message":"Created","body":{"id":"123"}}`,
		},
		{
			name:     "Response with Empty Body",
			response: Response[any]{Code: 204, Message: "No Content", Body: nil},
			expected: `{"code":204,"message":"No Content"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonResponse := tt.response.Marshal()
			assert.JSONEq(t, tt.expected, string(jsonResponse))
		})
	}
}

func Test_SendResponse(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		code           int
		body           any
		expectedCode   int
		expectedBody   string
		expectedHeader string
	}{
		{
			name:           "200 OK with TestResponse body",
			message:        "Success",
			code:           http.StatusOK,
			body:           &TestResponse{Key: "value"},
			expectedCode:   http.StatusOK,
			expectedBody:   `{"code":200,"message":"Success","body":{"key":"value"}}`,
			expectedHeader: "application/json",
		},
		{
			name:           "404 Not Found without body",
			message:        "Not Found",
			code:           http.StatusNotFound,
			body:           nil,
			expectedCode:   http.StatusNotFound,
			expectedBody:   `{"code":404,"message":"Not Found"}`,
			expectedHeader: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			switch body := tt.body.(type) {
			case *TestResponse:
				SendResponse[TestResponse](recorder, tt.message, tt.code, body)
			default:
				SendResponse(recorder, tt.message, tt.code, &tt.body)
			}

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedHeader, recorder.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.expectedBody, recorder.Body.String())
		})
	}
}

func TestWriteResponse(t *testing.T) {
	tests := []struct {
		name         string
		response     Response[any]
		expectedCode int
		expectedBody string
	}{
		{
			name:         "200 OK with Body",
			response:     Response[any]{Code: 200, Message: "OK", Body: map[string]string{"id": "123"}},
			expectedCode: 200,
			expectedBody: `{"code":200,"message":"OK","body":{"id":"123"}}`,
		},
		{
			name:         "500 Internal Server Error without Body",
			response:     Response[any]{Code: 500, Message: "Internal Server Error"},
			expectedCode: 500,
			expectedBody: `{"code":500,"message":"Internal Server Error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			writeResponse(recorder, &tt.response)

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.expectedBody, recorder.Body.String())
		})
	}
}
