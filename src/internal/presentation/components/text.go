package components

import (
	"fmt"
	"strings"

	"github.com/jterrazz/jterrazz-cli/internal/presentation/theme"
)

// =============================================================================
// Styled Text Helpers
// =============================================================================

// Muted renders text in muted style
func Muted(text string) string {
	return theme.Muted.Render(text)
}

// Success renders text in success style
func Success(text string) string {
	return theme.Success.Render(text)
}

// Warning renders text in warning style
func Warning(text string) string {
	return theme.Warning.Render(text)
}

// Danger renders text in danger style
func Danger(text string) string {
	return theme.Danger.Render(text)
}

// Special renders text in special/highlight style
func Special(text string) string {
	return theme.Special.Render(text)
}

// Bold renders text in bold
func Bold(text string) string {
	return theme.Selected.Render(text)
}

// =============================================================================
// Semantic Style Rendering
// =============================================================================

// StyledValue renders a value using a semantic style name
// Supported: "success", "warning", "danger", "special", "muted", "normal"
func StyledValue(value string, style string) string {
	return theme.Render(value, style)
}

// =============================================================================
// Padded Text
// =============================================================================

// PadRight pads text to the right to reach the given width
func PadRight(text string, width int) string {
	visibleLen := VisibleLen(text)
	if visibleLen >= width {
		return text
	}
	return text + strings.Repeat(" ", width-visibleLen)
}

// PadLeft pads text to the left to reach the given width
func PadLeft(text string, width int) string {
	visibleLen := VisibleLen(text)
	if visibleLen >= width {
		return text
	}
	return strings.Repeat(" ", width-visibleLen) + text
}

// =============================================================================
// Formatted Cell
// =============================================================================

// Cell renders a cell with specific width and style
func Cell(text string, width int, style string) string {
	padded := fmt.Sprintf("%-*s", width, text)
	return theme.Render(padded, style)
}

// CellNormal renders a normal styled cell
func CellNormal(text string, width int) string {
	return theme.Cell.Render(fmt.Sprintf("%-*s", width, text))
}

// CellMuted renders a muted styled cell
func CellMuted(text string, width int) string {
	return theme.CellMuted.Render(fmt.Sprintf("%-*s", width, text))
}

// CellMethod renders a method column cell (very dim)
func CellMethod(text string, width int) string {
	return theme.Method.Render(fmt.Sprintf("%-*s", width, text))
}

// CellSpecial renders a special/highlighted cell
func CellSpecial(text string, width int) string {
	if text == "" {
		return theme.Muted.Render(fmt.Sprintf("%-*s", width, text))
	}
	return theme.Special.Render(fmt.Sprintf("%-*s", width, text))
}

// =============================================================================
// Breadcrumb
// =============================================================================

// Breadcrumb renders breadcrumb navigation
func Breadcrumb(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return theme.BreadcrumbActive.Render(parts[0])
	}

	var result strings.Builder
	for i, part := range parts {
		if i == len(parts)-1 {
			result.WriteString(theme.BreadcrumbActive.Render(part))
		} else {
			result.WriteString(theme.Breadcrumb.Render(part + " > "))
		}
	}
	return result.String()
}

// =============================================================================
// Section Headers
// =============================================================================

// SectionLine renders a section header with lines
func SectionLine(title string) string {
	line := "───"
	return theme.Section.Render(line + " " + title + " " + line)
}
