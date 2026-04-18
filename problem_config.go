package httpsuite

import "strings"

var defaultProblemConfig = NewProblemConfig()

// ProblemConfig controls how problem type URLs are generated.
type ProblemConfig struct {
	BaseURL        string
	ErrorTypePaths map[string]string
}

// NewProblemConfig returns a config preloaded with the default problem type paths.
func NewProblemConfig() ProblemConfig {
	return ProblemConfig{
		ErrorTypePaths: map[string]string{
			"validation_error":  "/errors/validation-error",
			"not_found_error":   "/errors/not-found",
			"server_error":      "/errors/server-error",
			"bad_request_error": "/errors/bad-request",
		},
	}
}

// DefaultProblemConfig returns a copy of the package default config.
func DefaultProblemConfig() ProblemConfig {
	return defaultProblemConfig.Clone()
}

func mergeProblemConfig(config *ProblemConfig) ProblemConfig {
	merged := DefaultProblemConfig()
	if config == nil {
		return merged
	}

	if config.BaseURL != "" {
		merged.BaseURL = config.BaseURL
	}
	for key, value := range config.ErrorTypePaths {
		merged.ErrorTypePaths[key] = value
	}
	return merged.Clone()
}

// Clone returns a deep copy of the config.
func (c ProblemConfig) Clone() ProblemConfig {
	clone := ProblemConfig{
		BaseURL:        strings.TrimRight(c.BaseURL, "/"),
		ErrorTypePaths: make(map[string]string, len(c.ErrorTypePaths)),
	}
	for key, value := range c.ErrorTypePaths {
		clone.ErrorTypePaths[key] = normalizeProblemPath(value)
	}
	return clone
}

// TypeURL builds the full type URL for a known error type.
func (c ProblemConfig) TypeURL(errorType string) string {
	path, exists := c.ErrorTypePaths[errorType]
	if !exists {
		return BlankURL
	}

	baseURL := strings.TrimRight(c.BaseURL, "/")
	if baseURL == "" || baseURL == BlankURL {
		return normalizeProblemPath(path)
	}

	return baseURL + normalizeProblemPath(path)
}

// GetProblemTypeURL returns the default problem type URL for a known error type.
func GetProblemTypeURL(errorType string) string {
	return defaultProblemConfig.TypeURL(errorType)
}

func normalizeProblemPath(path string) string {
	if path == "" {
		return BlankURL
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") || path == BlankURL {
		return path
	}
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}
