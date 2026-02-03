package components

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Viewport wraps bubbles/viewport with convenience methods
type Viewport struct {
	Model               viewport.Model
	ShowScrollIndicator bool
}

// NewViewport creates a new viewport with the given dimensions
func NewViewport(width, height int) *Viewport {
	return &Viewport{
		Model:               viewport.New(width, height),
		ShowScrollIndicator: true,
	}
}

// SetSize updates the viewport dimensions
func (v *Viewport) SetSize(width, height int) {
	v.Model.Width = width
	v.Model.Height = height
}

// SetContent sets the viewport content
func (v *Viewport) SetContent(content string) {
	v.Model.SetContent(content)
}

// Update processes messages and returns commands
func (v *Viewport) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	v.Model, cmd = v.Model.Update(msg)
	return cmd
}

// View renders the viewport content
func (v *Viewport) View() string {
	return v.Model.View()
}

// ScrollPercent returns the current scroll position as a percentage (0-100)
func (v *Viewport) ScrollPercent() int {
	return int(v.Model.ScrollPercent() * 100)
}

// GotoTop scrolls to the top of the content
func (v *Viewport) GotoTop() {
	v.Model.GotoTop()
}

// GotoBottom scrolls to the bottom of the content
func (v *Viewport) GotoBottom() {
	v.Model.GotoBottom()
}

// AtTop returns true if scrolled to the top
func (v *Viewport) AtTop() bool {
	return v.Model.AtTop()
}

// AtBottom returns true if scrolled to the bottom
func (v *Viewport) AtBottom() bool {
	return v.Model.AtBottom()
}

// YPosition returns/sets the Y position of the viewport
func (v *Viewport) YPosition() int {
	return v.Model.YPosition
}

// SetYPosition sets the Y position
func (v *Viewport) SetYPosition(pos int) {
	v.Model.YPosition = pos
}

// Width returns the viewport width
func (v *Viewport) Width() int {
	return v.Model.Width
}

// Height returns the viewport height
func (v *Viewport) Height() int {
	return v.Model.Height
}
