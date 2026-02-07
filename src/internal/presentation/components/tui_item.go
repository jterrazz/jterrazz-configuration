package components

import (
	"fmt"
	"strings"

	"github.com/jterrazz/jterrazz-cli/src/internal/presentation/theme"
)

const (
	// IndentSpaces is the number of spaces per indent level
	IndentSpaces = 4
)

// ItemKind defines the type of list item
type ItemKind int

const (
	KindHeader     ItemKind = iota // Non-selectable section header
	KindNavigation                 // Navigates to another view
	KindAction                     // Runs a one-shot action
	KindToggle                     // Has on/off state (checkbox)
	KindExpandable                 // Can expand/collapse (tree node)
)

// ItemState defines the visual state of an item
type ItemState int

const (
	StateNone      ItemState = iota // No state indicator
	StateChecked                    // Checked/enabled (✓)
	StateUnchecked                  // Unchecked/disabled (○)
	StateLoading                    // Loading state
)

// Item represents a generic list item
type Item struct {
	Kind        ItemKind
	Label       string
	Description string
	State       ItemState
	Expanded    bool
	Indent      int // Nesting level (0 = root, 1 = nested, etc.)
	DescWidth   int // Width for description alignment (0 = no padding)
}

// Selectable returns true if this item can be selected/focused
func (i Item) Selectable() bool {
	return i.Kind != KindHeader
}

// Render renders the item as a styled string
func (i Item) Render(selected bool, labelWidth int, width int, spinnerFrame string) string {
	switch i.Kind {
	case KindHeader:
		return renderSection(i.Label, width)

	case KindNavigation:
		return i.renderNavigation(selected, labelWidth)

	case KindAction:
		return i.renderAction(selected)

	case KindToggle:
		return i.renderToggle(selected, labelWidth, spinnerFrame)

	case KindExpandable:
		return i.renderExpandable(selected, labelWidth)
	}

	return ""
}

func renderSection(title string, width int) string {
	// Use SectionHeader style (thick bordered box)
	return SectionHeader(title, width)
}

func (i Item) renderNavigation(selected bool, labelWidth int) string {
	prefix := i.buildPrefix(selected)

	paddedLabel := i.Label
	if labelWidth > 0 {
		paddedLabel = fmt.Sprintf("%-*s", labelWidth, i.Label)
	}

	if selected {
		return theme.Selected.Render(fmt.Sprintf("%s%s", prefix, paddedLabel)) + RenderDescription(i.Description)
	}
	return prefix + theme.Action.Render(paddedLabel) + RenderDescription(i.Description)
}

func (i Item) renderAction(selected bool) string {
	prefix := i.buildPrefix(selected)
	style := theme.Normal
	if selected {
		style = theme.Selected
	}
	return style.Render(prefix + i.Label)
}

func (i Item) renderToggle(selected bool, labelWidth int, spinnerFrame string) string {
	var statusBadge string
	switch i.State {
	case StateChecked:
		statusBadge = BadgeOK()
	case StateUnchecked:
		statusBadge = BadgeError()
	case StateLoading:
		if spinnerFrame != "" {
			statusBadge = theme.SpinnerStyle.Render(spinnerFrame)
		} else {
			statusBadge = theme.BadgeLoading.Render(theme.IconLoading)
		}
	default:
		statusBadge = " "
	}

	prefix := i.buildPrefix(selected)
	style := theme.Normal
	if selected {
		style = theme.Selected
	}

	paddedLabel := i.Label
	if labelWidth > 0 {
		paddedLabel = fmt.Sprintf("%-*s", labelWidth, i.Label)
	}

	// Status on the right: label  description  status
	if i.Description != "" {
		desc := i.Description
		if i.DescWidth > 0 {
			desc = fmt.Sprintf("%-*s", i.DescWidth, i.Description)
		}
		return style.Render(fmt.Sprintf("%s%s", prefix, paddedLabel)) +
			RenderDescription(desc) + ColumnSeparator + statusBadge
	}

	return style.Render(fmt.Sprintf("%s%s", prefix, paddedLabel)) + ColumnSeparator + statusBadge
}

func (i Item) renderExpandable(selected bool, labelWidth int) string {
	prefix := i.buildPrefix(selected)

	paddedLabel := i.Label
	if labelWidth > 0 {
		paddedLabel = fmt.Sprintf("%-*s", labelWidth, i.Label)
	}

	style := theme.Normal
	if selected {
		style = theme.Selected
	}

	base := style.Render(fmt.Sprintf("%s%s", prefix, paddedLabel))
	return base + RenderDescription(i.Description)
}

func (i Item) indentPrefix() string {
	if i.Indent <= 0 {
		return ""
	}
	return strings.Repeat(" ", i.Indent*IndentSpaces)
}

// buildPrefix returns the prefix string for an item (indent + selection indicator)
func (i Item) buildPrefix(selected bool) string {
	indent := i.indentPrefix()
	if selected {
		return indent + " " + theme.IconSelected + " "
	}
	return indent + "   "
}
