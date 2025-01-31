package httpsuite

import "sync"

const BlankUrl = "about:blank"

var (
	mu             sync.RWMutex
	problemBaseURL = BlankUrl
	errorTypePaths = map[string]string{
		"validation_error":  "/errors/validation-error",
		"not_found_error":   "/errors/not-found",
		"server_error":      "/errors/server-error",
		"bad_request_error": "/errors/bad-request",
	}
)

// ProblemDetails conforms to RFC 9457, providing a standard format for describing errors in HTTP APIs.
type ProblemDetails struct {
	Type       string                 `json:"type"`                 // A URI reference identifying the problem type.
	Title      string                 `json:"title"`                // A short, human-readable summary of the problem.
	Status     int                    `json:"status"`               // The HTTP status code.
	Detail     string                 `json:"detail,omitempty"`     // Detailed explanation of the problem.
	Instance   string                 `json:"instance,omitempty"`   // A URI reference identifying the specific instance of the problem.
	Extensions map[string]interface{} `json:"extensions,omitempty"` // Custom fields for additional details.
}

// NewProblemDetails creates a ProblemDetails instance with standard fields.
func NewProblemDetails(status int, problemType, title, detail string) *ProblemDetails {
	if problemType == "" {
		problemType = BlankUrl
	}
	return &ProblemDetails{
		Type:   problemType,
		Title:  title,
		Status: status,
		Detail: detail,
	}
}

// SetProblemBaseURL configures the base URL used in the "type" field for ProblemDetails.
//
// This function allows applications using httpsuite to provide a custom domain and structure
// for error documentation URLs. By setting this base URL, the library can generate meaningful
// and discoverable problem types.
//
// Parameters:
// - baseURL: The base URL where error documentation is hosted (e.g., "https://api.mycompany.com").
//
// Example usage:
//
//	httpsuite.SetProblemBaseURL("https://api.mycompany.com")
//
// Once configured, generated ProblemDetails will include a "type" such as:
//
//	"https://api.mycompany.com/errors/validation-error"
//
// If the base URL is not set, the default value for the "type" field will be "about:blank".
func SetProblemBaseURL(baseURL string) {
	mu.Lock()
	defer mu.Unlock()
	problemBaseURL = baseURL
}

// SetProblemErrorTypePath sets or updates the path for a specific error type.
//
// This allows applications to define custom paths for error documentation.
//
// Parameters:
// - errorType: The unique key identifying the error type (e.g., "validation_error").
// - path: The path under the base URL where the error documentation is located.
//
// Example usage:
//
//	httpsuite.SetProblemErrorTypePath("validation_error", "/errors/validation-error")
//
// After setting this path, the generated problem type for "validation_error" will be:
//
//	"https://api.mycompany.com/errors/validation-error"
func SetProblemErrorTypePath(errorType, path string) {
	mu.Lock()
	defer mu.Unlock()
	errorTypePaths[errorType] = path
}

// SetProblemErrorTypePaths sets or updates multiple paths for different error types.
//
// This allows applications to define multiple custom paths at once.
//
// Parameters:
// - paths: A map of error types to paths (e.g., {"validation_error": "/errors/validation-error"}).
//
// Example usage:
//
//	paths := map[string]string{
//	    "validation_error":  "/errors/validation-error",
//	    "not_found_error":   "/errors/not-found",
//	}
//	httpsuite.SetProblemErrorTypePaths(paths)
//
// This method overwrites any existing paths with the same keys.
func SetProblemErrorTypePaths(paths map[string]string) {
	mu.Lock()
	defer mu.Unlock()
	for errorType, path := range paths {
		errorTypePaths[errorType] = path
	}
}

// GetProblemTypeURL get the full problem type URL based on the error type.
//
// If the error type is not found in the predefined paths, it returns a default unknown error path.
//
// Parameters:
// - errorType: The unique key identifying the error type (e.g., "validation_error").
//
// Example usage:
//
//	problemTypeURL := GetProblemTypeURL("validation_error")
func GetProblemTypeURL(errorType string) string {
	if path, exists := errorTypePaths[errorType]; exists {
		return getProblemBaseURL() + path
	}

	return BlankUrl
}

// getProblemBaseURL just return the baseURL if it isn't "about:blank"
func getProblemBaseURL() string {
	if problemBaseURL == BlankUrl {
		return ""
	}
	return problemBaseURL
}
