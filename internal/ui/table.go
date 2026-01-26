package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// RenderTable creates a styled table with the WSO2 color scheme
func RenderTable(headers []string, rows [][]string) string {
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(BorderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return HeaderStyle
			case row%2 == 0:
				return EvenRowStyle
			default:
				return OddRowStyle
			}
		}).
		Headers(headers...).
		Rows(rows...)

	return t.String()
}

// RenderTableWithTitle creates a styled table with a title header
func RenderTableWithTitle(title string, headers []string, rows [][]string) string {
	output := TitleStyle.Render(title) + "\n\n"
	output += RenderTable(headers, rows)
	return output
}

// StatusCell returns a styled status string based on the status value
func StatusCell(status string) string {
	if status == "" {
		return MutedStyle.Render("—")
	}
	switch status {
	case "active", "running", "success", "healthy":
		return SuccessStyle.Render(status)
	case "inactive", "stopped", "failed", "error":
		return ErrorStyle.Render(status)
	case "pending", "building", "deploying":
		return WarningStyle.Render(status)
	default:
		return MutedStyle.Render(status)
	}
}

// RenderPaginationInfo returns a styled pagination footer
func RenderPaginationInfo(offset, limit, total int) string {
	// Handle edge cases
	if total == 0 {
		return ""
	}
	if limit <= 0 {
		return ""
	}
	if offset >= total {
		return MutedStyle.Render("No results in this range")
	}

	start := offset + 1
	end := offset + limit
	if end > total {
		end = total
	}

	info := MutedStyle.Render(fmt.Sprintf("Showing %d-%d of %d", start, end, total))

	// Show hint for next page if there are more results
	if end < total {
		nextOffset := offset + limit
		info += "\n" + MutedStyle.Render(fmt.Sprintf("Use --offset %d to see next page", nextOffset))
	}

	return info
}
