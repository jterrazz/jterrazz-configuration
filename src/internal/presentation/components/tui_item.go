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
}

// Selectable returns true if this item can be selected/focused
func (i Item) Selectable() bool {
	return i.Kind != KindHeader
}

// Render renders the item as a styled string
func (i Item) Render(selected bool, labelWidth int) string {
	switch i.Kind {
	case KindHeader:
		return renderSection(i.Label)

	case KindNavigation:
		return i.renderNavigation(selected, labelWidth)

	case KindAction:
		return i.renderAction(selected)

	case KindToggle:
		return i.renderToggle(selected, labelWidth)

	case KindExpandable:
		return i.renderExpandable(selected, labelWidth)
	}

	return ""
}

func renderSection(title string) string {
	line := "───"
	return theme.Section.Render(line + " " + title + " " + line)
}

func (i Item) renderNavigation(selected bool, labelWidth int) string {
	indent := i.indentPrefix()

	icon := "→"
	prefix := indent + "  "

	paddedLabel := i.Label
	if labelWidth > 0 {
		paddedLabel = fmt.Sprintf("%-*s", labelWidth, i.Label)
	}

	if selected {
		prefix = indent + theme.IconSelected + " "
		return theme.Selected.Render(prefix+icon+" "+paddedLabel) + "  " + theme.Muted.Render(i.Description)
	}
	return prefix + theme.Special.Render(icon) + " " + theme.Normal.Render(paddedLabel) + "  " + theme.Muted.Render(i.Description)
}

func (i Item) renderAction(selected bool) string {
	indent := i.indentPrefix()

	icon := theme.IconArrowRight
	prefix := indent + "  "

	if selected {
		prefix = indent + theme.IconSelected + " "
		return theme.Selected.Render(prefix + icon + " " + i.Label)
	}
	return prefix + theme.Action.Render(icon) + " " + theme.Normal.Render(i.Label)
}

func (i Item) renderToggle(selected bool, labelWidth int) string {
	indent := i.indentPrefix()

	var status string
	var style = theme.Muted

	switch i.State {
	case StateChecked:
		status = theme.IconCheck
		style = theme.Success
	case StateUnchecked:
		status = theme.IconCheckboxUnchecked
		style = theme.Muted
	case StateLoading:
		status = theme.IconCheckboxLoading
		style = theme.Action
	default:
		status = " "
	}

	prefix := indent + "  "
	if selected {
		prefix = indent + theme.IconSelected + " "
		style = theme.Selected
	}

	if labelWidth > 0 && i.Description != "" {
		paddedLabel := fmt.Sprintf("%-*s", labelWidth, i.Label)
		return style.Render(fmt.Sprintf("%s%s %s", prefix, status, paddedLabel)) +
			theme.Muted.Render("  "+i.Description)
	}

	return style.Render(fmt.Sprintf("%s%s %s", prefix, status, i.Label))
}

func (i Item) renderExpandable(selected bool, labelWidth int) string {
	indent := i.indentPrefix()

	var arrow string
	if i.Expanded {
		arrow = theme.IconArrowDown
	} else {
		arrow = theme.IconArrowRight
	}

	prefix := indent + "  "
	if selected {
		prefix = indent + theme.IconSelected + " "
	}

	paddedLabel := i.Label
	if labelWidth > 0 {
		paddedLabel = fmt.Sprintf("%-*s", labelWidth, i.Label)
	}

	if selected {
		base := theme.Selected.Render(fmt.Sprintf("%s%s %s", prefix, arrow, paddedLabel))
		if i.Description != "" {
			return base + theme.Muted.Render("  "+i.Description)
		}
		return base
	}

	base := theme.Normal.Render(fmt.Sprintf("%s%s %s", prefix, arrow, paddedLabel))
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
