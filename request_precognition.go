package httpsuite

import (
	"errors"
	"net/http"
	"strings"
)

// PrecognitionHeader is the request header a client sends to ask for
// validation-only processing on an endpoint that enables ParseOptions.Precognition.
const PrecognitionHeader = "Precognition"

// PrecognitionSuccessHeader is set on the 204 response of a successful precognitive request.
const PrecognitionSuccessHeader = "Precognition-Success"

// ErrPrecognitionHandled is returned by ParseRequest when a precognitive request
// was validated and the response has already been written. Handlers must stop
// processing; the usual `if err != nil { return }` does exactly that.
var ErrPrecognitionHandled = errors.New("precognition request handled")

// IsPrecognitionRequest reports whether the request asks for validation-only processing.
func IsPrecognitionRequest(r *http.Request) bool {
	if r == nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(r.Header.Get(PrecognitionHeader)), "true")
}

// IsPrecognitionHandled reports whether err signals a handled precognitive request.
func IsPrecognitionHandled(err error) bool {
	return errors.Is(err, ErrPrecognitionHandled)
}

func writePrecognitionSuccess(w http.ResponseWriter) {
	w.Header().Set(PrecognitionSuccessHeader, "true")
	w.WriteHeader(http.StatusNoContent)
}
