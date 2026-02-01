package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// ColumnConfig defines the style configuration for a table column
type ColumnConfig struct {
	Width int
	Color string // Use ColorPrimary, ColorMuted, ColorSpecial, or empty for default
}

// RenderTable creates a styled table with the given rows and column configuration
func RenderTable(rows [][]string, columns []ColumnConfig) string {
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(BorderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)

			if col < len(columns) {
				cfg := columns[col]
				if cfg.Width > 0 {
					style = style.Width(cfg.Width)
				}
				if cfg.Color != "" {
					style = style.Foreground(lipgloss.Color(cfg.Color))
				}
			}

			return style
		}).
		Rows(rows...)

	return t.Render()
}

// Common column configurations for reuse
var (
	// StatusTableColumns is for tables with name, detail, status
	StatusTableColumns = []ColumnConfig{
		{Width: 18, Color: ColorPrimary},
		{Width: 0, Color: ColorMuted},
		{Width: 0, Color: ""},
	}

	// CheckTableColumns is for tables with name, description, detail, status
	CheckTableColumns = []ColumnConfig{
		{Width: 16, Color: ColorPrimary},
		{Width: 30, Color: ColorMuted},
		{Width: 0, Color: ColorSpecial},
		{Width: 0, Color: ""},
	}

	// ToolTableColumns is for tool tables with name, version, method, status
	ToolTableColumns = []ColumnConfig{
		{Width: 14, Color: ColorPrimary},
		{Width: 14, Color: ColorMuted},
		{Width: 8, Color: ColorMuted},
		{Width: 0, Color: ""},
	}

	// ResourceTableColumns is for resource tables with name, value
	ResourceTableColumns = []ColumnConfig{
		{Width: 14, Color: ColorPrimary},
		{Width: 0, Color: ""},
	}

	// DiskTableColumns is for disk usage tables
	DiskTableColumns = []ColumnConfig{
		{Width: 18, Color: ColorPrimary},
		{Width: 0, Color: ""},
	}

	// CacheTableColumns is for cache tables
	CacheTableColumns = []ColumnConfig{
		{Width: 20, Color: ColorMuted},
		{Width: 0, Color: ""},
	}
)
