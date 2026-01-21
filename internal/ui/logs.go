package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Log level styles for color coding
var (
	LogErrorStyle = lipgloss.NewStyle().
			Foreground(Red500).
			Bold(true)

	LogWarnStyle = lipgloss.NewStyle().
			Foreground(Yellow500).
			Bold(true)

	LogInfoStyle = lipgloss.NewStyle().
			Foreground(Blue500)

	LogDebugStyle = lipgloss.NewStyle().
			Foreground(Gray500)

	LogTimestampStyle = lipgloss.NewStyle().
				Foreground(Gray600)
)

// LogLevelStyle returns the appropriate style for a log level
func LogLevelStyle(level string) lipgloss.Style {
	switch strings.ToUpper(level) {
	case "ERROR":
		return LogErrorStyle
	case "WARN", "WARNING":
		return LogWarnStyle
	case "INFO":
		return LogInfoStyle
	case "DEBUG":
		return LogDebugStyle
	default:
		return MutedStyle
	}
}

// FormatLogLevel returns a styled log level prefix
func FormatLogLevel(level string) string {
	upperLevel := strings.ToUpper(level)
	style := LogLevelStyle(level)

	switch upperLevel {
	case "ERROR":
		return style.Render("[ERROR]")
	case "WARN", "WARNING":
		return style.Render("[WARN ]")
	case "INFO":
		return style.Render("[INFO ]")
	case "DEBUG":
		return style.Render("[DEBUG]")
	default:
		if level != "" {
			return style.Render("[" + upperLevel + "]")
		}
		return ""
	}
}

// FormatLogTimestamp formats a timestamp for log display
func FormatLogTimestamp(timestamp string) string {
	// Try parsing as RFC3339
	if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
		return LogTimestampStyle.Render(t.Format("15:04:05"))
	}

	// Try parsing as RFC3339Nano
	if t, err := time.Parse(time.RFC3339Nano, timestamp); err == nil {
		return LogTimestampStyle.Render(t.Format("15:04:05"))
	}

	// Fallback: extract time from ISO-like format
	if strings.Contains(timestamp, "T") {
		parts := strings.SplitN(timestamp, "T", 2)
		if len(parts) == 2 {
			timePart := parts[1]
			// Strip fractional seconds first (before timezone)
			if idx := strings.Index(timePart, "."); idx != -1 {
				// Keep everything before the dot, but check for timezone after
				afterDot := timePart[idx+1:]
				timePart = timePart[:idx]
				// Check if timezone info exists after fractional seconds
				if tzIdx := strings.IndexAny(afterDot, "Z+-"); tzIdx != -1 {
					// Timezone found, timePart is already correct
				}
			}
			// Strip timezone markers
			if idx := strings.Index(timePart, "Z"); idx != -1 {
				timePart = timePart[:idx]
			}
			if idx := strings.Index(timePart, "+"); idx != -1 {
				timePart = timePart[:idx]
			}
			// Use LastIndex for negative timezone offset to avoid cutting time components
			if idx := strings.LastIndex(timePart, "-"); idx != -1 && idx > 5 {
				timePart = timePart[:idx]
			}
			return LogTimestampStyle.Render(timePart)
		}
	}

	// Return as-is if we can't parse
	return LogTimestampStyle.Render(timestamp)
}
