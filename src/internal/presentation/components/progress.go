package components

import (
	"fmt"
	"strings"

	"github.com/jterrazz/jterrazz-cli/internal/presentation/theme"
)

// ProgressBar renders a progress bar with optional count and spinner
type ProgressBar struct {
	Current     int
	Total       int
	Width       int
	ShowCount   bool
	ShowSpinner bool
	Spinner     string // Current spinner frame (if ShowSpinner is true)
}

// NewProgressBar creates a new progress bar with defaults
func NewProgressBar(current, total int) *ProgressBar {
	return &ProgressBar{
		Current:   current,
		Total:     total,
		Width:     30,
		ShowCount: true,
	}
}

// Render returns the progress bar as a string
func (p *ProgressBar) Render() string {
	if p.Total == 0 {
		return ""
	}

	// Calculate fill percentage
	filled := int(float64(p.Current) / float64(p.Total) * float64(p.Width))
	if filled > p.Width {
		filled = p.Width
	}

	// Build bar
	bar := theme.ProgressFilled.Render(strings.Repeat(theme.IconProgressFull, filled)) +
		theme.ProgressEmpty.Render(strings.Repeat(theme.IconProgressEmpty, p.Width-filled))

	// Build result
	var parts []string

	if p.ShowSpinner && p.Spinner != "" {
		parts = append(parts, theme.SpinnerStyle.Render(p.Spinner))
	}

	parts = append(parts, bar)

	if p.ShowCount {
		parts = append(parts, fmt.Sprintf("%d/%d", p.Current, p.Total))
	}

	return strings.Join(parts, " ")
}

// SetSpinner updates the current spinner frame
func (p *ProgressBar) SetSpinner(frame string) {
	p.ShowSpinner = true
	p.Spinner = frame
}

// IsComplete returns true if progress is at 100%
func (p *ProgressBar) IsComplete() bool {
	return p.Current >= p.Total
}

// CompletionMessage returns a success message when complete
func (p *ProgressBar) CompletionMessage() string {
	return theme.Success.Render(theme.IconCheck + " All checks complete")
}
