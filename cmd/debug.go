package cmd

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var debugStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))

// Debug prints a message only if verbose mode is enabled
func Debug(format string, args ...interface{}) {
	if Verbose {
		msg := fmt.Sprintf(format, args...)
		fmt.Println(debugStyle.Render("DEBUG: " + msg))
	}
}

// DebugRequest logs an API request in verbose mode
func DebugRequest(method, path string) {
	Debug("%s %s", method, path)
}

// DebugResponse logs an API response in verbose mode
func DebugResponse(status int, duration time.Duration) {
	Debug("Response: %d (%v)", status, duration.Round(time.Millisecond))
}
