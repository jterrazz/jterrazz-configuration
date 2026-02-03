package components

import (
	"fmt"
	"strings"

	"github.com/jterrazz/jterrazz-cli/internal/presentation/theme"
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
	indent := i.indentPrefix()

	style := theme.Normal

	prefix := indent + "   "
	if selected {
		prefix = indent + " " + theme.IconSelected + " "
		style = theme.Selected
	}

	paddedLabel := i.Label
	if labelWidth > 0 {
		paddedLabel = fmt.Sprintf("%-*s", labelWidth, i.Label)
	}

	return style.Render(fmt.Sprintf("%s%s", prefix, paddedLabel)) +
		theme.Muted.Render("  "+i.Description)
}

func (i Item) renderAction(selected bool) string {
	indent := i.indentPrefix()

	prefix := indent + "   "

	if selected {
		prefix = indent + " " + theme.IconSelected + " "
		return theme.Selected.Render(prefix + i.Label)
	}
	return prefix + theme.Action.Render(i.Label)
}

func (i Item) renderToggle(selected bool, labelWidth int, spinnerFrame string) string {
	indent := i.indentPrefix()

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

	prefix := indent + "   "
	style := theme.Normal
	if selected {
		prefix = indent + " " + theme.IconSelected + " "
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
			theme.Muted.Render("  "+desc) + "  " + statusBadge
	}

	return style.Render(fmt.Sprintf("%s%s", prefix, paddedLabel)) + "  " + statusBadge
}

func (i Item) renderExpandable(selected bool, labelWidth int) string {
	indent := i.indentPrefix()

	prefix := indent + "   "
	if selected {
		prefix = indent + " " + theme.IconSelected + " "
	}

	paddedLabel := i.Label
	if labelWidth > 0 {
		paddedLabel = fmt.Sprintf("%-*s", labelWidth, i.Label)
	}

	if selected {
		base := theme.Selected.Render(fmt.Sprintf("%s%s", prefix, paddedLabel))
		if i.Description != "" {
			return base + theme.Muted.Render("  "+i.Description)
		}
		return base
	}

	base := theme.Normal.Render(fmt.Sprintf("%s%s", prefix, paddedLabel))
	if i.Description != "" {
		return base + theme.Muted.Render("  "+i.Description)
	}
	return base
}

func (i Item) indentPrefix() string {
	if i.Indent <= 0 {
		return ""
	}
	return strings.Repeat(" ", i.Indent*IndentSpaces)
}
