package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const Version = "0.1.3"

// Banner styles
var (
	bannerBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Orange500).
			Padding(1, 2)

	logoStyle = lipgloss.NewStyle().
			Foreground(Orange500).
			Bold(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(Orange500).
			Bold(true)

	versionStyle = lipgloss.NewStyle().
			Foreground(Gray500)

	labelStyle = lipgloss.NewStyle().
			Foreground(Gray500)

	valueStyle = lipgloss.NewStyle().
			Foreground(Gray700)

	tipHeaderStyle = lipgloss.NewStyle().
			Foreground(Orange500).
			Bold(true)

	tipStyle = lipgloss.NewStyle().
			Foreground(Gray600)

	separatorStyle = lipgloss.NewStyle().
			Foreground(Orange500)
)

// RenderBanner renders the welcome banner
func RenderBanner(org, project string) string {
	// Left side content
	var leftLines []string
	leftLines = append(leftLines, logoStyle.Render("  ▄▀█ █▀▄▀█ █▀█"))
	leftLines = append(leftLines, logoStyle.Render("  █▀█ █░▀░█ █▀▀"))
	leftLines = append(leftLines, "")
	leftLines = append(leftLines, titleStyle.Render("WSO2 AI Agent Management Platform"))
	leftLines = append(leftLines, "")

	// Context info
	if org != "" {
		leftLines = append(leftLines, labelStyle.Render("org: ")+valueStyle.Render(org))
	}
	if project != "" {
		leftLines = append(leftLines, labelStyle.Render("project: ")+valueStyle.Render(project))
	}

	// Right side content - tips
	var rightLines []string
	rightLines = append(rightLines, tipHeaderStyle.Render("Quick Start"))
	rightLines = append(rightLines, tipStyle.Render("  amp orgs list"))
	rightLines = append(rightLines, tipStyle.Render("  amp projects list"))
	rightLines = append(rightLines, tipStyle.Render("  amp agents list"))
	rightLines = append(rightLines, "")
	rightLines = append(rightLines, tipHeaderStyle.Render("Help"))
	rightLines = append(rightLines, tipStyle.Render("  amp --help"))
	rightLines = append(rightLines, tipStyle.Render("  amp config list"))

	// Pad arrays to same length
	for len(leftLines) < len(rightLines) {
		leftLines = append(leftLines, "")
	}
	for len(rightLines) < len(leftLines) {
		rightLines = append(rightLines, "")
	}

	// Build content with separator
	var contentLines []string
	leftWidth := 40
	for i := range leftLines {
		left := lipgloss.NewStyle().Width(leftWidth).Render(leftLines[i])
		sep := separatorStyle.Render("│")
		right := "  " + rightLines[i]
		contentLines = append(contentLines, left+sep+right)
	}

	content := strings.Join(contentLines, "\n")

	// Header line
	header := fmt.Sprintf("─ %s %s ─",
		titleStyle.Render("amp-cli"),
		versionStyle.Render("v"+Version))

	return header + "\n" + bannerBorder.Render(content)
}
