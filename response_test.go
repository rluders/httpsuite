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
		errs         []Error
		meta         *Meta
		expectedCode int
		expectedJSON string
	}{
		{
			name:         "200 OK with TestResponse body",
			code:         http.StatusOK,
			data:         &TestResponse{Key: "value"},
			errs:         nil,
			expectedCode: http.StatusOK,
			expectedJSON: `{"data":{"key":"value"}}`,
		},
		{
			name:         "404 Not Found without body",
			code:         http.StatusNotFound,
			data:         nil,
			errs:         []Error{{Code: http.StatusNotFound, Message: "Not Found"}},
			expectedCode: http.StatusNotFound,
			expectedJSON: `{"errors":[{"code":404,"message":"Not Found"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			switch data := tt.data.(type) {
			case TestResponse:
				SendResponse[TestResponse](w, tt.code, data, tt.errs, tt.meta)
			default:
				SendResponse[any](w, tt.code, tt.data, tt.errs, tt.meta)
			}

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.expectedJSON, w.Body.String())
		})
	}
}
