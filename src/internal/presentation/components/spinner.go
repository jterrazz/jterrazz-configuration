package components

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jterrazz/jterrazz-cli/src/internal/presentation/theme"
)

// SpinnerFPS is the animation speed for all spinners
const SpinnerFPS = 80 * time.Millisecond

// Spinner wraps the bubbles spinner with our theme
type Spinner struct {
	Model spinner.Model
}

// NewSpinner creates a new spinner with braille dots animation
func NewSpinner() Spinner {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: theme.BrailleSpinner,
		FPS:    SpinnerFPS,
	}
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorSpinner))
	return Spinner{Model: s}
}

// NewSpinnerWithStyle creates a spinner with custom color
func NewSpinnerWithStyle(color string) Spinner {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: theme.BrailleSpinner,
		FPS:    SpinnerFPS,
	}
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	return Spinner{Model: s}
}

// NewSpinnerModel creates a raw spinner.Model for use in custom views
func NewSpinnerModel() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: theme.BrailleSpinner,
		FPS:    SpinnerFPS,
	}
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorSpinner))
	return s
}

// View returns the current spinner frame
func (s Spinner) View() string {
	return s.Model.View()
}

// Tick returns the tick command for animation
func (s Spinner) Tick() tea.Cmd {
	return s.Model.Tick
}

// Update updates the spinner state and returns the next tick command
func (s *Spinner) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	s.Model, cmd = s.Model.Update(msg)
	return cmd
}

// SpinnerFrames returns the raw spinner frames for custom rendering
func SpinnerFrames() []string {
	return theme.BrailleSpinner
}
