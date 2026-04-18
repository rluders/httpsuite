package httpsuite

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkSendResponse(b *testing.B) {
	payload := testResponse{Key: "value"}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		SendResponse(w, http.StatusOK, payload, nil, nil)
	}
}
