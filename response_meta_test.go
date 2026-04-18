package httpsuite

import (
	"encoding/json"
	"strings"
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

func TestResponseSerializationPreservesMeaningfulZeroValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		response     any
		wantContains []string
	}{
		{
			name:         "zero data value",
			response:     Response[int]{Data: 0},
			wantContains: []string{`"data":0`},
		},
		{
			name:         "false bool data value",
			response:     Response[bool]{Data: false},
			wantContains: []string{`"data":false`},
		},
		{
			name:         "cursor false flags",
			response:     Response[[]string]{Data: []string{"a"}, Meta: NewCursorMeta("", "", false, false)},
			wantContains: []string{`"has_next":false`, `"has_prev":false`},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.response)
			if err != nil {
				t.Fatalf("marshal response: %v", err)
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(string(body), want) {
					t.Fatalf("expected %q in %q", want, string(body))
				}
			}
		})
	}
}
