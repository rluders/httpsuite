package httpsuite

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type testRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type bodyOnlyRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (r *testRequest) SetParam(fieldName, value string) error {
	switch strings.ToLower(fieldName) {
	case "id":
		id, err := strconv.Atoi(value)
		if err != nil {
			return errors.New("invalid id")
		}
		r.ID = id
	default:
		return fmt.Errorf("parameter %s cannot be set", fieldName)
	}
	return nil
}

func testParamExtractor(r *http.Request, key string) string {
	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) > 2 && strings.EqualFold(key, "id") {
		return pathSegments[2]
	}
	return ""
}

type stubValidator struct {
	problem *ProblemDetails
}

func (s stubValidator) Validate(any) *ProblemDetails {
	return s.problem
}
