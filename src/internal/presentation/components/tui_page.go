package components

import (
	"strings"

	"github.com/jterrazz/jterrazz-cli/internal/presentation/theme"
)

// Page represents a full-screen TUI page layout
type Page struct {
	Title       string   // Page title (used if no breadcrumbs)
	Breadcrumbs []string // Navigation breadcrumbs (overrides title)
	Content     string   // Main content area (pre-rendered)
	Help        string   // Help text at bottom
	Message     string   // Status message
	Processing  bool     // True if action in progress (affects message style)
	Width       int
	Height      int
}

// NewPage creates a new page with the given title
func NewPage(title string) *Page {
	return &Page{
		Title:  title,
		Width:  80,
		Height: 24,
	}
}

// SetSize updates the page dimensions
func (p *Page) SetSize(width, height int) {
	p.Width = width
	p.Height = height
}

// ContentHeight returns the available height for content
// Subtracts space for: title (1) + blank line (1) + help (1) + message (1) = 4 lines
func (p *Page) ContentHeight() int {
	h := p.Height - 4
	if h < 1 {
		h = 1
	}
	return h
}

// Render renders the complete page
func (p *Page) Render() string {
	var b strings.Builder

	// Header: simple title like status view
	if len(p.Breadcrumbs) > 0 {
		b.WriteString(PageIndent + theme.SectionTitle.Render(strings.ToUpper(p.Breadcrumbs[len(p.Breadcrumbs)-1])) + "\n\n")
	} else if p.Title != "" {
		b.WriteString(PageIndent + theme.SectionTitle.Render(strings.ToUpper(p.Title)) + "\n\n")
	}

	// Main content
	if p.Content != "" {
		b.WriteString(p.Content)
	}

	// Help text
	if p.Help != "" {
		b.WriteString(theme.Help.Render(p.Help))
	}

	// Status message
	if p.Message != "" {
		b.WriteString("\n")
		if p.Processing {
			b.WriteString(theme.Action.Render(p.Message))
		} else {
			b.WriteString(theme.Success.Render(p.Message))
		}
	}

	return b.String()
}

// DefaultHelp returns standard navigation help text
func DefaultHelp() string {
	return "↑/↓ navigate • enter select • q quit"
}

// DefaultHelpWithBack returns navigation help text with back option
func DefaultHelpWithBack() string {
	return "↑/↓ navigate • enter select/toggle • esc back • q quit"
}
