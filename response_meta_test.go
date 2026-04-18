package httpsuite

import (
	"encoding/json"
	"testing"
)

func TestMetaSerialization(t *testing.T) {
	t.Parallel()

	pageResponse := Response[[]string]{
		Data: []string{"a"},
		Meta: NewPageMeta(1, 10, 15),
	}
	pageBody, err := json.Marshal(pageResponse)
	if err != nil {
		t.Fatalf("marshal page response: %v", err)
	}
	if !json.Valid(pageBody) {
		t.Fatalf("expected valid page response json")
	}

	cursorResponse := Response[[]string]{
		Data: []string{"a"},
		Meta: NewCursorMeta("next", "", true, false),
	}
	cursorBody, err := json.Marshal(cursorResponse)
	if err != nil {
		t.Fatalf("marshal cursor response: %v", err)
	}
	if !json.Valid(cursorBody) {
		t.Fatalf("expected valid cursor response json")
	}
}
