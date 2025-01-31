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

// TestRequest includes custom type annotation for UUID type.
type TestRequest struct {
	ID   int    `json:"id" validate:"required,gt=0"`
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

// MyParamExtractor extracts parameters from the path, assuming the request URL follows a pattern like "/test/{id}".
func MyParamExtractor(r *http.Request, key string) string {
	pathSegments := strings.Split(r.URL.Path, "/")
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
		name       string
		args       args
		want       *TestRequest
		wantErr    assert.ErrorAssertionFunc
		wantDetail *ProblemDetails
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
			want:       &TestRequest{ID: 123, Name: "Test"},
			wantErr:    assert.NoError,
			wantDetail: nil,
		},
		{
			name: "Missing body",
			args: args{
				w:          httptest.NewRecorder(),
				r:          httptest.NewRequest("POST", "/test/123", nil),
				pathParams: []string{"ID"},
			},
			want:       nil,
			wantErr:    assert.Error,
			wantDetail: NewProblemDetails(http.StatusBadRequest, "Validation Error", "One or more fields failed validation."),
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
			want:       nil,
			wantErr:    assert.Error,
			wantDetail: NewProblemDetails(http.StatusBadRequest, "Invalid Request", "invalid character 'i' looking for beginning of object key string"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function under test.
			w := tt.args.w
			got, err := ParseRequest[*TestRequest](w, tt.args.r, MyParamExtractor, tt.args.pathParams...)

			// Validate the error response if applicable.
			if !tt.wantErr(t, err, fmt.Sprintf("parseRequest(%v, %v, %v)", tt.args.w, tt.args.r, tt.args.pathParams)) {
				return
			}

			// Check ProblemDetails if an error was expected.
			if tt.wantDetail != nil {
				rec := w.(*httptest.ResponseRecorder)
				var pd ProblemDetails
				decodeErr := json.NewDecoder(rec.Body).Decode(&pd)
				assert.NoError(t, decodeErr, "Failed to decode problem details response")
				assert.Equal(t, tt.wantDetail.Title, pd.Title, "Problem detail title mismatch")
				assert.Equal(t, tt.wantDetail.Status, pd.Status, "Problem detail status mismatch")
				assert.Contains(t, pd.Detail, tt.wantDetail.Detail, "Problem detail message mismatch")
			}

			// Validate successful response.
			assert.Equalf(t, tt.want, got, "parseRequest(%v, %v, %v)", w, tt.args.r, tt.args.pathParams)
		})
	}
}
