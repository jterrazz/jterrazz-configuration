package components

import (
	"strings"

	"github.com/jterrazz/jterrazz-cli/internal/ui/theme"
)

// BoxStyle defines the border style for a box
type BoxStyle int

const (
	BoxRounded BoxStyle = iota // ╭─╮ │ │ ╰─╯
	BoxThick                   // ┏━┓ ┃ ┃ ┗━┛
)

// =============================================================================
// Section Header Box (thick borders)
// =============================================================================

// SectionHeader renders a section header with thick borders
// ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
// ┃  SYSTEM                                                                ┃
// ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
func SectionHeader(title string, width int) string {
	innerWidth := width - 2 // for the borders
	if innerWidth < 10 {
		innerWidth = 10
	}

	displayTitle := strings.ToUpper(title)
	padding := innerWidth - len(displayTitle) - 2 // -2 for "  " prefix
	if padding < 0 {
		padding = 0
	}

	borderStyle := theme.SectionBorder

	top := borderStyle.Render(theme.BoxThickTopLeft + strings.Repeat(theme.BoxThickHorizontal, innerWidth) + theme.BoxThickTopRight)
	middle := borderStyle.Render(theme.BoxThickVertical) + "  " + theme.SectionTitle.Render(displayTitle) + strings.Repeat(" ", padding) + borderStyle.Render(theme.BoxThickVertical)
	bottom := borderStyle.Render(theme.BoxThickBottomLeft + strings.Repeat(theme.BoxThickHorizontal, innerWidth) + theme.BoxThickBottomRight)

	return top + "\n" + middle + "\n" + bottom
}

// =============================================================================
// Subsection Box (rounded borders with title)
// =============================================================================

// SubsectionBox renders a subsection with rounded borders and embedded title
// ╭─ Title ────────────────────────────────────────────────────────────────╮
// │ content line 1                                                         │
// │ content line 2                                                         │
// ╰────────────────────────────────────────────────────────────────────────╯
func SubsectionBox(title string, lines []string, width int) string {
	innerWidth := width - 4 // account for border + padding
	if innerWidth < 20 {
		innerWidth = 20
	}

	borderStyle := theme.SectionBorder

	// Build top border with title: ╭─ Title ─────────────────╮
	titleLen := len(title) + 4 // "─ " + title + " ─"
	topBorderRight := strings.Repeat(theme.BoxRoundedHorizontal, innerWidth-titleLen)
	if topBorderRight == "" {
		topBorderRight = theme.BoxRoundedHorizontal
	}

	top := borderStyle.Render(theme.BoxRoundedTopLeft+theme.BoxRoundedHorizontal+" ") +
		theme.SubSection.Render(title) +
		borderStyle.Render(" "+theme.BoxRoundedHorizontal+topBorderRight+theme.BoxRoundedTopRight)

	bottom := borderStyle.Render(theme.BoxRoundedBottomLeft + strings.Repeat(theme.BoxRoundedHorizontal, innerWidth+2) + theme.BoxRoundedBottomRight)

	// Pad content lines
	var paddedLines []string
	for _, line := range lines {
		padding := innerWidth - VisibleLen(line)
		if padding < 0 {
			padding = 0
		}
		paddedLine := borderStyle.Render(theme.BoxRoundedVertical+" ") + line + strings.Repeat(" ", padding) + borderStyle.Render(" "+theme.BoxRoundedVertical)
		paddedLines = append(paddedLines, paddedLine)
	}

	return top + "\n" + strings.Join(paddedLines, "\n") + "\n" + bottom
}

// =============================================================================
// Simple Box (no title)
// =============================================================================

// SimpleBox renders a simple box with rounded borders
func SimpleBox(content string, width int) string {
	lines := strings.Split(content, "\n")
	innerWidth := width - 4
	if innerWidth < 10 {
		innerWidth = 10
	}

	borderStyle := theme.SectionBorder

	top := borderStyle.Render(theme.BoxRoundedTopLeft + strings.Repeat(theme.BoxRoundedHorizontal, innerWidth+2) + theme.BoxRoundedTopRight)
	bottom := borderStyle.Render(theme.BoxRoundedBottomLeft + strings.Repeat(theme.BoxRoundedHorizontal, innerWidth+2) + theme.BoxRoundedBottomRight)

	var paddedLines []string
	for _, line := range lines {
		padding := innerWidth - VisibleLen(line)
		if padding < 0 {
			padding = 0
		}
		paddedLine := borderStyle.Render(theme.BoxRoundedVertical+" ") + line + strings.Repeat(" ", padding) + borderStyle.Render(" "+theme.BoxRoundedVertical)
		paddedLines = append(paddedLines, paddedLine)
	}

	return top + "\n" + strings.Join(paddedLines, "\n") + "\n" + bottom
}

// =============================================================================
// Helpers
// =============================================================================

// VisibleLen returns the visible length of a string, stripping ANSI escape codes
func VisibleLen(s string) int {
	inEscape := false
	length := 0
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		length++
	}
	return length
}
