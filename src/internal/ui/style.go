package ui

import "github.com/charmbracelet/lipgloss"

// Shared color palette for TUI
const (
	ColorPrimary   = "212" // Pink/magenta for selection
	ColorSecondary = "99"  // Purple for headers
	ColorSuccess   = "42"  // Green for success/installed
	ColorWarning   = "214" // Orange for actions
	ColorDanger    = "196" // Red for errors/not configured
	ColorMuted     = "241" // Gray for dimmed text
	ColorText      = "252" // Light gray for normal text
	ColorSpecial   = "86"  // Cyan for special highlights
	ColorBorder    = "238" // Dark gray for borders
)

// Shared TUI styles
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSecondary)).
			Bold(true)

	SectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSecondary)).
			Bold(true).
			MarginTop(1)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPrimary)).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSuccess))

	ActionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWarning)).
			Bold(true)

	DangerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorDanger))

	MutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted))

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorText))

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted)).
			MarginTop(1)

	BreadcrumbStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted))

	BreadcrumbActiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorSecondary)).
				Bold(true)

	SpecialStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSpecial))

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWarning))

	SubSectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted)).
			Italic(true)

	BorderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorBorder))
)

// Icons
const (
	IconSelected   = "›"
	IconCheck      = "✓"
	IconCross      = "✗"
	IconBullet     = "•"
	IconArrowRight = "▶"
	IconArrowDown  = "▼"
)

// RenderSection renders a header section
func RenderSection(title string) string {
	line := "───"
	return SectionStyle.Render(line + " " + title + " " + line)
}

// RenderBreadcrumb renders breadcrumb navigation
func RenderBreadcrumb(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return BreadcrumbActiveStyle.Render(parts[0])
	}

	result := ""
	for i, part := range parts {
		if i == len(parts)-1 {
			result += BreadcrumbActiveStyle.Render(part)
		} else {
			result += BreadcrumbStyle.Render(part + " → ")
		}
	}
	return result
}
