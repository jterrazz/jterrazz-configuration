package ui

import "github.com/charmbracelet/lipgloss"

// Shared color palette for TUI
const (
	ColorPrimary   = lipgloss.Color("212") // Pink/magenta for selection
	ColorSecondary = lipgloss.Color("99")  // Purple for headers
	ColorSuccess   = lipgloss.Color("42")  // Green for success/installed
	ColorWarning   = lipgloss.Color("214") // Orange for actions
	ColorDanger    = lipgloss.Color("196") // Red for errors/not configured
	ColorMuted     = lipgloss.Color("241") // Gray for dimmed text
	ColorText      = lipgloss.Color("252") // Light gray for normal text
)

// Shared TUI styles
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true)

	SectionStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true).
			MarginTop(1)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	ActionStyle = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true)

	DangerStyle = lipgloss.NewStyle().
			Foreground(ColorDanger)

	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	NormalStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			MarginTop(1)

	BreadcrumbStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	BreadcrumbActiveStyle = lipgloss.NewStyle().
				Foreground(ColorSecondary).
				Bold(true)
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
