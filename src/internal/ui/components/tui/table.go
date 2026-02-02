package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jterrazz/jterrazz-cli/internal/ui/theme"
)

// TableColumn defines a column configuration
type TableColumn struct {
	Width int
	Color string // lipgloss color code (e.g., "212", "241")
}

// Table renders a bordered table with the given rows and column configurations
func Table(rows [][]string, columns []TableColumn) string {
	if len(rows) == 0 {
		return ""
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorBorder))).
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)

			if col < len(columns) {
				if columns[col].Width > 0 {
					style = style.Width(columns[col].Width)
				}
				if columns[col].Color != "" {
					style = style.Foreground(lipgloss.Color(columns[col].Color))
				}
			}

			return style
		}).
		Rows(rows...)

	return t.Render()
}

// SimpleTable renders a table with default styling
// First column is highlighted, rest are muted
func SimpleTable(rows [][]string, firstColWidth int) string {
	if len(rows) == 0 {
		return ""
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorBorder))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if col == 0 {
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color(theme.ColorPrimary)).
					PaddingLeft(1).PaddingRight(1).
					Width(firstColWidth)
			}
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.ColorMuted)).
				PaddingLeft(1).PaddingRight(1)
		}).
		Rows(rows...)

	return t.Render()
}

// StatusTable renders a table with name, detail, and status columns
// Commonly used for showing check results with ✓/✗ indicators
func StatusTable(rows [][]string, nameWidth int) string {
	if len(rows) == 0 {
		return ""
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorBorder))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch col {
			case 0:
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color(theme.ColorPrimary)).
					PaddingLeft(1).PaddingRight(1).
					Width(nameWidth)
			case 1:
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color(theme.ColorMuted)).
					PaddingLeft(1).PaddingRight(1).
					Width(30)
			case 2:
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color(theme.ColorSecondary)).
					PaddingLeft(1).PaddingRight(1)
			default:
				return lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)
			}
		}).
		Rows(rows...)

	return t.Render()
}
