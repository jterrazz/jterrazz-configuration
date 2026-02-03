package components

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jterrazz/jterrazz-cli/internal/presentation/theme"
)

// Spinner wraps the bubbles spinner with our theme
type Spinner struct {
	Model spinner.Model
}

// NewSpinner creates a new spinner with braille dots animation
func NewSpinner() Spinner {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: theme.BrailleSpinner,
		FPS:    120 * time.Millisecond,
	}
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorSpinner))
	return Spinner{Model: s}
}

// NewSpinnerWithStyle creates a spinner with custom color
func NewSpinnerWithStyle(color string) Spinner {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: theme.BrailleSpinner,
		FPS:    120 * time.Millisecond,
	}
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	return Spinner{Model: s}
}

// View returns the current spinner frame
func (s Spinner) View() string {
	return s.Model.View()
}

// Tick returns the tick command for animation
func (s Spinner) Tick() tea.Cmd {
	return s.Model.Tick
}

// Update updates the spinner state
func (s *Spinner) Update(msg spinner.TickMsg) {
	s.Model, _ = s.Model.Update(msg)
}

// SpinnerFrames returns the raw spinner frames for custom rendering
func SpinnerFrames() []string {
	return theme.BrailleSpinner
}
