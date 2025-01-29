package httpsuite

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// TestRequest includes custom type annotation for UUID type
type TestRequest struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

func (r *TestRequest) SetParam(fieldName, value string) error {
	switch strings.ToLower(fieldName) {
	case "id":
		id, err := strconv.Atoi(value)
		if err != nil {
			return errors.New("invalid id")
		}
		r.ID = id
	default:
		fmt.Printf("Parameter %s cannot be set", fieldName)
	}
	return nil
}

// This implementation extracts parameters from the path, assuming the request URL follows a pattern
// like "/test/{id}", where "id" is a path parameter.
func MyParamExtractor(r *http.Request, key string) string {
	// Here, we can extract parameters directly from the URL path for simplicity.
	// Example: for "/test/123", if key is "ID", we want to capture "123".
	pathSegments := strings.Split(r.URL.Path, "/")

	// You should know how the path is structured; in this case, we expect the ID to be the second segment.
	if len(pathSegments) > 2 && key == "ID" {
		return pathSegments[2]
	}
	return ""
}

func Test_ParseRequest(t *testing.T) {
	type args struct {
		w          http.ResponseWriter
		r          *http.Request
		pathParams []string
	}
	type testCase[T any] struct {
		name    string
		args    args
		want    *TestRequest
		wantErr assert.ErrorAssertionFunc
	}

	tests := []testCase[TestRequest]{
		{
			name: "Successful Request",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					body, _ := json.Marshal(TestRequest{Name: "Test"})
					req := httptest.NewRequest("POST", "/test/123", bytes.NewBuffer(body))
					req.URL.Path = "/test/123"
					return req
				}(),
				pathParams: []string{"ID"},
			},
			want:    &TestRequest{ID: 123, Name: "Test"},
			wantErr: assert.NoError,
		},
		{
			name: "Missing body",
			args: args{
				w:          httptest.NewRecorder(),
				r:          httptest.NewRequest("POST", "/test/123", nil),
				pathParams: []string{"ID"},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "Invalid JSON Body",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest("POST", "/test/123", bytes.NewBufferString("{invalid-json}"))
					req.URL.Path = "/test/123"
					return req
				}(),
				pathParams: []string{"ID"},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRequest[*TestRequest](tt.args.w, tt.args.r, MyParamExtractor, tt.args.pathParams...)
			if !tt.wantErr(t, err, fmt.Sprintf("parseRequest(%v, %v, %v)", tt.args.w, tt.args.r, tt.args.pathParams)) {
				return
			}
			assert.Equalf(t, tt.want, got, "parseRequest(%v, %v, %v)", tt.args.w, tt.args.r, tt.args.pathParams)
		})
	}
}
