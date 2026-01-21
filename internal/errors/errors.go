package errors

import (
	"fmt"
	"strings"
)

// CLIError represents a user-friendly error with suggestions
type CLIError struct {
	Message    string
	Cause      error
	Suggestion string
	Context    map[string]string
}

func (e *CLIError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error
func (e *CLIError) Unwrap() error {
	return e.Cause
}

// New creates a new CLIError with just a message
func New(message string) *CLIError {
	return &CLIError{Message: message}
}

// Wrap wraps an existing error with a message
func Wrap(err error, message string) *CLIError {
	return &CLIError{
		Message: message,
		Cause:   err,
	}
}

// WithSuggestion adds a suggestion to the error
func (e *CLIError) WithSuggestion(suggestion string) *CLIError {
	e.Suggestion = suggestion
	return e
}

// WithContext adds context to the error
func (e *CLIError) WithContext(key, value string) *CLIError {
	if e.Context == nil {
		e.Context = make(map[string]string)
	}
	e.Context[key] = value
	return e
}

// Common error constructors with helpful suggestions

// AuthError creates an authentication error with suggestions
func AuthError(cause error) *CLIError {
	return &CLIError{
		Message:    "Authentication failed",
		Cause:      cause,
		Suggestion: "Check your credentials with 'amp config list' or run 'amp login' to reconfigure",
	}
}

// NotFoundError creates a not found error with suggestions
func NotFoundError(resourceType, resourceName string) *CLIError {
	plural := pluralizeResource(resourceType)
	listCmd := fmt.Sprintf("amp %s list", plural)
	return &CLIError{
		Message:    fmt.Sprintf("%s '%s' not found", resourceType, resourceName),
		Suggestion: fmt.Sprintf("List available %s with '%s'", plural, listCmd),
	}
}

// pluralizeResource returns a lowercase plural form of a resource type.
// Handles words ending in 's' to avoid incorrect forms like "statuss".
func pluralizeResource(resourceType string) string {
	lower := strings.ToLower(resourceType)
	if strings.HasSuffix(lower, "s") {
		return lower
	}
	return lower + "s"
}

// ConnectionError creates a connection error with suggestions
func ConnectionError(url string, cause error) *CLIError {
	return &CLIError{
		Message:    "Cannot connect to API server",
		Cause:      cause,
		Suggestion: "Check if the server is running and verify api_url with 'amp config list'",
		Context:    map[string]string{"url": url},
	}
}

// TimeoutError creates a timeout error
func TimeoutError() *CLIError {
	return &CLIError{
		Message:    "Request timed out",
		Suggestion: "The server may be slow or unreachable. Try again later.",
	}
}

// MissingFlagError creates an error for missing required flags
func MissingFlagError(flagName, command string) *CLIError {
	return &CLIError{
		Message:    fmt.Sprintf("Required flag '--%s' not provided", flagName),
		Suggestion: fmt.Sprintf("See 'amp %s --help' for usage", command),
	}
}

// MissingConfigError creates an error for missing configuration
func MissingConfigError(configKey, setCommand string) *CLIError {
	return &CLIError{
		Message:    fmt.Sprintf("Configuration '%s' is not set", configKey),
		Suggestion: fmt.Sprintf("Set it with 'amp config set %s' or run 'amp login'", setCommand),
	}
}

// APIError creates an error from an API response
func APIError(statusCode int, body string) *CLIError {
	message := fmt.Sprintf("API error (status %d)", statusCode)
	var suggestion string

	switch statusCode {
	case 401:
		message = "Authentication failed (401 Unauthorized)"
		suggestion = "Your credentials may be invalid or expired. Run 'amp login' to reconfigure"
	case 403:
		message = "Access denied (403 Forbidden)"
		suggestion = "You don't have permission for this action. Check your account permissions"
	case 404:
		message = "Resource not found (404)"
		suggestion = "The requested resource doesn't exist. Check the name and try again"
	case 500:
		message = "Server error (500)"
		suggestion = "The server encountered an error. Try again later or contact support"
	case 502, 503, 504:
		message = fmt.Sprintf("Service unavailable (%d)", statusCode)
		suggestion = "The server is temporarily unavailable. Try again later"
	}

	return &CLIError{
		Message:    message,
		Suggestion: suggestion,
		Context:    map[string]string{"response": truncateBody(body)},
	}
}

// truncateBody limits the error body length for display (UTF-8 safe)
func truncateBody(body string) string {
	const maxRunes = 200
	if len(body) <= maxRunes {
		return body
	}
	runes := []rune(body)
	if len(runes) <= maxRunes {
		return body
	}
	return string(runes[:maxRunes]) + "..."
}
