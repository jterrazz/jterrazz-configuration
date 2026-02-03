package theme

// =============================================================================
// Unicode Icons & Symbols
// =============================================================================

const (
	// Navigation
	IconSelected = "›"

	// Status
	IconCheck   = "✓"
	IconCross   = "✗"
	IconWarning = "!"
	IconBullet  = "•"

	// Service status
	IconServiceOn  = "●"
	IconServiceOff = "○"

	// Progress
	IconProgressFull  = "█"
	IconProgressEmpty = "░"

	// Loading
	IconLoading = "◌"
)

// =============================================================================
// Spinner Frames
// =============================================================================

// BrailleSpinner is the braille dots spinner animation frames
var BrailleSpinner = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// =============================================================================
// Box Drawing Characters
// =============================================================================

// Rounded box (for subsections)
const (
	BoxRoundedTopLeft     = "╭"
	BoxRoundedTopRight    = "╮"
	BoxRoundedBottomLeft  = "╰"
	BoxRoundedBottomRight = "╯"
	BoxRoundedHorizontal  = "─"
	BoxRoundedVertical    = "│"
)

// Thick box (for section headers)
const (
	BoxThickTopLeft     = "┏"
	BoxThickTopRight    = "┓"
	BoxThickBottomLeft  = "┗"
	BoxThickBottomRight = "┛"
	BoxThickHorizontal  = "━"
	BoxThickVertical    = "┃"
)

// =============================================================================
// Rendered Status Icons (convenience)
// =============================================================================

// StatusIcon returns a styled status icon
func StatusIcon(ok bool) string {
	if ok {
		return BadgeOK.Render(IconCheck)
	}
	return BadgeError.Render(IconCross)
}

// ServiceIcon returns a styled service status icon
func ServiceIcon(running bool) string {
	if running {
		return ServiceRunning.Render(IconServiceOn)
	}
	return ServiceStopped.Render(IconServiceOff)
}
