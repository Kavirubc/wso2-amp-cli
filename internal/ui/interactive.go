package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	promptStyle = lipgloss.NewStyle().
			Foreground(Orange500).
			Bold(true)

	inputStyle = lipgloss.NewStyle().
			Foreground(Gray700)

	helpHintStyle = lipgloss.NewStyle().
			Foreground(Gray500).
			Italic(true)
)

// InteractiveModel is the Bubble Tea model for interactive mode
type InteractiveModel struct {
	input    textinput.Model
	output   string
	quitting bool
	handler  func(string) string
}

// NewInteractiveModel creates a new interactive model
func NewInteractiveModel(handler func(string) string) InteractiveModel {
	ti := textinput.New()
	ti.Placeholder = "type a command..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50
	ti.PromptStyle = promptStyle
	ti.TextStyle = inputStyle

	return InteractiveModel{
		input:   ti,
		handler: handler,
	}
}

func (m InteractiveModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InteractiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			input := strings.TrimSpace(m.input.Value())
			if input == "" {
				return m, nil
			}

			// Handle exit commands
			lower := strings.ToLower(input)
			if lower == "exit" || lower == "quit" || lower == "q" {
				m.quitting = true
				return m, tea.Quit
			}

			// Handle clear
			if lower == "clear" || lower == "cls" {
				m.output = ""
				m.input.SetValue("")
				return m, nil
			}

			// Execute command via handler
			if m.handler != nil {
				m.output = m.handler(input)
			}

			m.input.SetValue("")
			return m, nil
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m InteractiveModel) View() string {
	if m.quitting {
		return helpHintStyle.Render("Goodbye!") + "\n"
	}

	var b strings.Builder

	// Show output
	if m.output != "" {
		b.WriteString(m.output)
		b.WriteString("\n\n")
	}

	// Prompt
	b.WriteString(promptStyle.Render("amp > "))
	b.WriteString(m.input.View())
	b.WriteString("\n")
	b.WriteString(helpHintStyle.Render("  ctrl+c to exit"))
	b.WriteString("\n")

	return b.String()
}
