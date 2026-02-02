package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

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
	IconWarning    = "!"
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

// =============================================================================
// Status Indicators
// =============================================================================

// Status represents the state of an item (installed, missing, warning, etc.)
type Status int

const (
	StatusSuccess Status = iota // Installed, configured, running
	StatusDanger                // Not installed, error
	StatusWarning               // Needs attention
	StatusMuted                 // Neutral/informational
)

// StatusIcon returns the icon for a status
func StatusIcon(s Status) string {
	switch s {
	case StatusSuccess:
		return SuccessStyle.Render(IconCheck)
	case StatusDanger:
		return DangerStyle.Render(IconCross)
	case StatusWarning:
		return WarningStyle.Render(IconWarning)
	default:
		return MutedStyle.Render(IconBullet)
	}
}

// StatusFromBool returns StatusSuccess if true, StatusDanger if false
func StatusFromBool(ok bool) Status {
	if ok {
		return StatusSuccess
	}
	return StatusDanger
}

// =============================================================================
// Row Renderers for CLI output
// =============================================================================

// Row represents a styled output row
type Row struct {
	Status Status
	Label  string
	Detail string
	Extra  string
}

// Render renders a row with status icon
func (r Row) Render() string {
	icon := StatusIcon(r.Status)
	if r.Detail != "" {
		return fmt.Sprintf("  %s %-14s %s", icon, r.Label, MutedStyle.Render(r.Detail))
	}
	return fmt.Sprintf("  %s %s", icon, r.Label)
}

// RenderRow is a convenience function for quick row rendering
func RenderRow(ok bool, label string, detail string) string {
	return Row{Status: StatusFromBool(ok), Label: label, Detail: detail}.Render()
}

// =============================================================================
// Styled Value Rendering
// =============================================================================

// RenderStyledValue renders a value with a semantic style name.
// This is the central function for styling dynamic values from config/parsing.
// Supported styles: "success", "warning", "danger", "special", "muted", "normal"
func RenderStyledValue(value string, style string) string {
	switch style {
	case "success":
		return SuccessStyle.Render(value)
	case "warning":
		return WarningStyle.Render(value)
	case "danger":
		return DangerStyle.Render(value)
	case "special":
		return SpecialStyle.Render(value)
	case "normal":
		return NormalStyle.Render(value)
	case "muted":
		fallthrough
	default:
		return MutedStyle.Render(value)
	}
}
