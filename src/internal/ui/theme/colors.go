package theme

import "github.com/charmbracelet/lipgloss"

// Color palette - ANSI 256 color codes
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
	ColorSpinner   = "39"  // Blue for spinners/progress
)

// Color returns a lipgloss.Color for the given color code
func Color(code string) lipgloss.Color {
	return lipgloss.Color(code)
}
