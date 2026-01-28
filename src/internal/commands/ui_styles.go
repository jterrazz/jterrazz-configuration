package commands

import "github.com/charmbracelet/lipgloss"

// Shared color palette for TUI
const (
	uiColorPrimary   = lipgloss.Color("212") // Pink/magenta for selection
	uiColorSecondary = lipgloss.Color("99")  // Purple for headers
	uiColorSuccess   = lipgloss.Color("42")  // Green for success/installed
	uiColorWarning   = lipgloss.Color("214") // Orange for actions
	uiColorDanger    = lipgloss.Color("196") // Red for errors/not configured
	uiColorMuted     = lipgloss.Color("241") // Gray for dimmed text
	uiColorText      = lipgloss.Color("252") // Light gray for normal text
)

// Shared TUI styles
var (
	// Title style
	uiTitleStyle = lipgloss.NewStyle().
			Foreground(uiColorSecondary).
			Bold(true).
			MarginBottom(1)

	// Section header style
	uiSectionStyle = lipgloss.NewStyle().
			Foreground(uiColorSecondary).
			Bold(true).
			MarginTop(1)

	// Selected item
	uiSelectedStyle = lipgloss.NewStyle().
			Foreground(uiColorPrimary).
			Bold(true)

	// Success/installed items
	uiSuccessStyle = lipgloss.NewStyle().
			Foreground(uiColorSuccess)

	// Warning/action items
	uiActionStyle = lipgloss.NewStyle().
			Foreground(uiColorWarning).
			Bold(true)

	// Danger/error items
	uiDangerStyle = lipgloss.NewStyle().
			Foreground(uiColorDanger)

	// Muted/dimmed text
	uiMutedStyle = lipgloss.NewStyle().
			Foreground(uiColorMuted)

	// Normal text
	uiNormalStyle = lipgloss.NewStyle().
			Foreground(uiColorText)

	// Help bar at bottom
	uiHelpStyle = lipgloss.NewStyle().
			Foreground(uiColorMuted).
			MarginTop(1)

	// Breadcrumb navigation
	uiBreadcrumbStyle = lipgloss.NewStyle().
				Foreground(uiColorMuted)

	uiBreadcrumbActiveStyle = lipgloss.NewStyle().
				Foreground(uiColorText).
				Bold(true)
)

// Icons
const (
	iconSelected   = "›"
	iconCheck      = "✓"
	iconCross      = "✗"
	iconBullet     = "•"
	iconArrowRight = "▶"
	iconArrowDown  = "▼"
)

// Helper to render a header section
func uiRenderSection(title string) string {
	line := "───"
	return uiSectionStyle.Render(line + " " + title + " " + line)
}

// Helper to render breadcrumb
func uiRenderBreadcrumb(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return uiBreadcrumbActiveStyle.Render(parts[0])
	}

	result := ""
	for i, part := range parts {
		if i == len(parts)-1 {
			result += uiBreadcrumbActiveStyle.Render(part)
		} else {
			result += uiBreadcrumbStyle.Render(part + " → ")
		}
	}
	return result
}
