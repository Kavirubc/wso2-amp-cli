package ui

import "github.com/charmbracelet/lipgloss"

// Title and Header styles
var (
	// Main title style - used for command headers
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Orange500).
			MarginBottom(1)

	// Subtitle style
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Gray600).
			MarginBottom(1)
)

// Table styles
var (
	// Table header row style
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Orange500).
			Padding(0, 1)

	// Regular cell style
	CellStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Alternating row styles for better readability
	OddRowStyle = lipgloss.NewStyle().
			Foreground(Gray700).
			Padding(0, 1)

	EvenRowStyle = lipgloss.NewStyle().
			Foreground(Gray900).
			Padding(0, 1)

	// Border style for tables
	BorderStyle = lipgloss.NewStyle().
			Foreground(Orange500)
)

// Status indicator styles
var (
	SuccessStyle = lipgloss.NewStyle().
			Foreground(Green500).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Red500).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Yellow500).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(Blue500)

	MutedStyle = lipgloss.NewStyle().
			Foreground(Gray500)
)

// Config display styles
var (
	// Key style for config keys
	KeyStyle = lipgloss.NewStyle().
			Foreground(Teal500).
			Bold(true).
			Width(22)

	// Value style for config values
	ValueStyle = lipgloss.NewStyle().
			Foreground(Gray700)

	// Masked value style (for secrets)
	MaskedStyle = lipgloss.NewStyle().
			Foreground(Gray500).
			Italic(true)
)

// Box and container styles
var (
	// Box style for containing content
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Orange500).
			Padding(1, 2)

	// Section header style
	SectionStyle = lipgloss.NewStyle().
			Foreground(Orange600).
			Bold(true).
			MarginTop(1).
			MarginBottom(1)
)

// Icon constants for consistent usage
const (
	IconSuccess = "✓"
	IconError   = "✗"
	IconWarning = "⚠"
	IconInfo    = "ℹ"
	IconList    = "📋"
	IconConfig  = "⚙️"
	IconAgent   = "🤖"
	IconBuild   = "🔨"
	IconDeploy  = "🚀"
)

// RenderSuccess renders a success message
func RenderSuccess(msg string) string {
	return SuccessStyle.Render(IconSuccess + " " + msg)
}

// RenderError renders an error message
func RenderError(msg string) string {
	return ErrorStyle.Render(IconError + " " + msg)
}

// RenderWarning renders a warning message
func RenderWarning(msg string) string {
	return WarningStyle.Render(IconWarning + " " + msg)
}

// RenderInfo renders an info message
func RenderInfo(msg string) string {
	return InfoStyle.Render(IconInfo + " " + msg)
}
