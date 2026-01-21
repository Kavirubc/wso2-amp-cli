package errors

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles for error rendering
var (
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Bold(true)

	causeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF"))

	suggestionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981"))

	contextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))
)

// Render returns a formatted string representation of the error
func (e *CLIError) Render() string {
	var sb strings.Builder

	// Main error message
	sb.WriteString(errorStyle.Render("✗ " + e.Message))

	// Underlying cause
	if e.Cause != nil {
		sb.WriteString("\n  ")
		sb.WriteString(causeStyle.Render(e.Cause.Error()))
	}

	// Context information (sorted for deterministic output)
	if len(e.Context) > 0 {
		sb.WriteString("\n")
		keys := make([]string, 0, len(e.Context))
		for k := range e.Context {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			value := e.Context[key]
			if value != "" {
				fmt.Fprintf(&sb, "\n  %s: %s",
					contextStyle.Render(key),
					causeStyle.Render(value))
			}
		}
	}

	// Suggestion
	if e.Suggestion != "" {
		sb.WriteString("\n\n")
		sb.WriteString(suggestionStyle.Render("💡 " + e.Suggestion))
	}

	return sb.String()
}

// RenderError formats any error for CLI output
func RenderError(err error) string {
	if err == nil {
		return ""
	}
	// Use errors.As to handle wrapped CLIError instances
	var cliErr *CLIError
	if errors.As(err, &cliErr) {
		return cliErr.Render()
	}
	// Fallback for regular errors
	return errorStyle.Render("✗ " + err.Error())
}
