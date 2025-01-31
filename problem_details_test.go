package httpsuite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetProblemBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Set valid base URL",
			input:    "https://api.example.com",
			expected: "https://api.example.com",
		},
		{
			name:     "Set base URL to blank",
			input:    BlankUrl,
			expected: BlankUrl,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetProblemBaseURL(tt.input)
			assert.Equal(t, tt.expected, problemBaseURL)
		})
	}
}

func Test_SetProblemErrorTypePath(t *testing.T) {
	tests := []struct {
		name     string
		errorKey string
		path     string
		expected string
	}{
		{
			name:     "Set custom error path",
			errorKey: "custom_error",
			path:     "/errors/custom-error",
			expected: "/errors/custom-error",
		},
		{
			name:     "Override existing path",
			errorKey: "validation_error",
			path:     "/errors/new-validation-error",
			expected: "/errors/new-validation-error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetProblemErrorTypePath(tt.errorKey, tt.path)
			assert.Equal(t, tt.expected, errorTypePaths[tt.errorKey])
		})
	}
}

func Test_GetProblemTypeURL(t *testing.T) {
	// Setup initial state
	SetProblemBaseURL("https://api.example.com")
	SetProblemErrorTypePath("validation_error", "/errors/validation-error")

	tests := []struct {
		name        string
		errorType   string
		expectedURL string
	}{
		{
			name:        "Valid error type",
			errorType:   "validation_error",
			expectedURL: "https://api.example.com/errors/validation-error",
		},
		{
			name:        "Unknown error type",
			errorType:   "unknown_error",
			expectedURL: BlankUrl,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetProblemTypeURL(tt.errorType)
			assert.Equal(t, tt.expectedURL, result)
		})
	}
}

func Test_getProblemBaseURL(t *testing.T) {
	tests := []struct {
		name           string
		baseURL        string
		expectedResult string
	}{
		{
			name:           "Base URL is set",
			baseURL:        "https://api.example.com",
			expectedResult: "https://api.example.com",
		},
		{
			name:           "Base URL is about:blank",
			baseURL:        BlankUrl,
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			problemBaseURL = tt.baseURL
			assert.Equal(t, tt.expectedResult, getProblemBaseURL())
		})
	}
}

func Test_NewProblemDetails(t *testing.T) {
	tests := []struct {
		name         string
		status       int
		problemType  string
		title        string
		detail       string
		expectedType string
	}{
		{
			name:         "All fields provided",
			status:       400,
			problemType:  "https://api.example.com/errors/validation-error",
			title:        "Validation Error",
			detail:       "Invalid input",
			expectedType: "https://api.example.com/errors/validation-error",
		},
		{
			name:         "Empty problem type",
			status:       404,
			problemType:  "",
			title:        "Not Found",
			detail:       "The requested resource was not found",
			expectedType: BlankUrl,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			details := NewProblemDetails(tt.status, tt.problemType, tt.title, tt.detail)
			assert.Equal(t, tt.status, details.Status)
			assert.Equal(t, tt.title, details.Title)
			assert.Equal(t, tt.detail, details.Detail)
			assert.Equal(t, tt.expectedType, details.Type)
		})
	}
}
