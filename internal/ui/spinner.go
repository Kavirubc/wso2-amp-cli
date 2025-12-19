package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// NewSpinner creates a new spinner with WSO2 theming
func NewSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(Orange500)
	return s
}

// NewMiniSpinner creates a smaller spinner for inline use
func NewMiniSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(Orange500)
	return s
}

// NewPulseSpinner creates a pulsing spinner
func NewPulseSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse
	s.Style = lipgloss.NewStyle().Foreground(Teal500)
	return s
}
