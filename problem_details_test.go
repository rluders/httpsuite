package httpsuite

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewProblemDetailsDefaults(t *testing.T) {
	t.Parallel()

	details := NewProblemDetails(700, "", "", "broken")
	if details.Status != 500 {
		t.Fatalf("expected status 500, got %d", details.Status)
	}
	if details.Type != BlankURL {
		t.Fatalf("expected type %q, got %q", BlankURL, details.Type)
	}
	if details.Title != "Internal Server Error" {
		t.Fatalf("expected fallback title, got %q", details.Title)
	}
}

func TestProblemDetailsJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := ProblemDetails{
		Type:     "https://example.com/errors/bad",
		Title:    "Bad Request",
		Status:   400,
		Detail:   "broken",
		Instance: "/things/1",
		Extensions: map[string]interface{}{
			"messageKey":    "errors.bad",
			"correlationId": "abc-123",
		},
	}

	body, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded ProblemDetails
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Fatalf("round trip mismatch:\nwant %#v\ngot  %#v", original, decoded)
	}
}

func TestProblemDetailsUnmarshalJSONNoExtensions(t *testing.T) {
	t.Parallel()

	decoded := ProblemDetails{Extensions: map[string]interface{}{"stale": true}}
	if err := json.Unmarshal([]byte(`{"type":"about:blank","title":"T","status":404}`), &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded.Extensions != nil {
		t.Fatalf("expected nil extensions, got %#v", decoded.Extensions)
	}
	if decoded.Status != 404 || decoded.Title != "T" {
		t.Fatalf("unexpected decode: %#v", decoded)
	}
}

func TestProblemDetailsUnmarshalJSONNull(t *testing.T) {
	t.Parallel()

	decoded := ProblemDetails{Title: "keep", Status: 400}
	if err := json.Unmarshal([]byte(`null`), &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded.Title != "keep" || decoded.Status != 400 {
		t.Fatalf("expected null to be a no-op, got %#v", decoded)
	}
}

func TestProblemDetailsUnmarshalJSONInvalid(t *testing.T) {
	t.Parallel()

	var decoded ProblemDetails
	if err := json.Unmarshal([]byte(`{"status":"nope"}`), &decoded); err == nil {
		t.Fatal("expected error for wrong status type")
	}
}

func TestProblemDetailsMarshalJSONFlattensExtensions(t *testing.T) {
	t.Parallel()

	problem := &ProblemDetails{
		Type:   BlankURL,
		Title:  "Bad Request",
		Status: 400,
		Detail: "broken",
		Extensions: map[string]interface{}{
			"trace_id":   "trace-123",
			"extensions": "ignored",
		},
	}

	body, err := json.Marshal(problem)
	if err != nil {
		t.Fatalf("marshal problem details: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("unmarshal problem details: %v", err)
	}

	if _, exists := payload["extensions"]; exists {
		t.Fatalf("expected flattened extensions, got nested payload %q", string(body))
	}
	if payload["trace_id"] != "trace-123" {
		t.Fatalf("expected trace_id extension, got %#v", payload["trace_id"])
	}
}
