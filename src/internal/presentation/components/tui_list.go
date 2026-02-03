package components

import "strings"

const (
	// DefaultTerminalWidth is the default terminal width
	DefaultTerminalWidth = 80
	// DefaultTerminalHeight is the default terminal height
	DefaultTerminalHeight = 24
	// ScrollMargin is the number of lines to keep visible above/below cursor
	ScrollMargin = 3
)

// List represents a scrollable list of items with cursor navigation
type List struct {
	Items      []Item
	Cursor     int
	Width      int
	Height     int
	LabelWidth int // For aligning descriptions (0 = no alignment)
}

// NewList creates a new list with the given items
func NewList(items []Item) *List {
	l := &List{
		Items:  items,
		Cursor: 0,
		Width:  DefaultTerminalWidth,
		Height: DefaultTerminalHeight,
	}
	// Start cursor on first selectable item
	l.skipToSelectable(1)
	return l
}

// SetSize updates the list dimensions
func (l *List) SetSize(width, height int) {
	l.Width = width
	l.Height = height
}

// Up moves the cursor up to the previous selectable item
func (l *List) Up() {
	l.move(-1)
}

// Down moves the cursor down to the next selectable item
func (l *List) Down() {
	l.move(1)
}

// move moves the cursor by delta, skipping non-selectable items
func (l *List) move(delta int) {
	newCursor := l.Cursor + delta

	// Skip non-selectable items
	for newCursor >= 0 && newCursor < len(l.Items) && !l.Items[newCursor].Selectable() {
		newCursor += delta
	}

	if newCursor >= 0 && newCursor < len(l.Items) {
		l.Cursor = newCursor
	}
}

// skipToSelectable moves cursor to the first selectable item in direction
func (l *List) skipToSelectable(direction int) {
	for l.Cursor >= 0 && l.Cursor < len(l.Items) && !l.Items[l.Cursor].Selectable() {
		l.Cursor += direction
	}
}

// Selected returns the currently selected item, or nil if none
func (l *List) Selected() *Item {
	if l.Cursor >= 0 && l.Cursor < len(l.Items) {
		return &l.Items[l.Cursor]
	}
	return nil
}

// SelectedIndex returns the current cursor position
func (l *List) SelectedIndex() int {
	return l.Cursor
}

// SetCursor sets the cursor to a specific position
func (l *List) SetCursor(index int) {
	if index >= 0 && index < len(l.Items) {
		l.Cursor = index
	}
}

// UpdateItem updates an item at the given index
func (l *List) UpdateItem(index int, item Item) {
	if index >= 0 && index < len(l.Items) {
		l.Items[index] = item
	}
}

// Render renders the visible portion of the list
func (l *List) Render(visibleHeight int) string {
	if len(l.Items) == 0 {
		return ""
	}

	if visibleHeight <= 0 {
		visibleHeight = l.Height
	}

	// Calculate visible range with scrolling
	startIdx := 0
	if l.Cursor > visibleHeight-ScrollMargin {
		startIdx = l.Cursor - visibleHeight + ScrollMargin
	}
	endIdx := startIdx + visibleHeight
	if endIdx > len(l.Items) {
		endIdx = len(l.Items)
	}

	var b strings.Builder
	for i := startIdx; i < endIdx; i++ {
		item := l.Items[i]
		selected := i == l.Cursor
		line := item.Render(selected, l.LabelWidth, l.Width)
		b.WriteString(line + "\n")
	}

	return b.String()
}

// CalculateLabelWidth calculates the max label width for alignment
func (l *List) CalculateLabelWidth() int {
	maxLen := 0
	for _, item := range l.Items {
		if len(item.Label) > maxLen {
			maxLen = len(item.Label)
		}
	}
	l.LabelWidth = maxLen
	return maxLen
}
