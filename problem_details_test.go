package httpsuite

import (
	"encoding/json"
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
