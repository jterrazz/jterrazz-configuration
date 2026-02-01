package ui

import (
	"fmt"
	"strings"
)

const (
	// IndentSpaces is the number of spaces per indent level
	IndentSpaces = 4
)

// ItemKind defines the type of list item
type ItemKind int

const (
	KindHeader     ItemKind = iota // Non-selectable section header
	KindAction                     // Clickable action
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
		return RenderSection(i.Label)

	case KindAction:
		return i.renderAction(selected)

	case KindToggle:
		return i.renderToggle(selected, labelWidth)

	case KindExpandable:
		return i.renderExpandable(selected, labelWidth)
	}

	return ""
}

func (i Item) renderAction(selected bool) string {
	indent := i.indentPrefix()
	prefix := indent + "  "

	if selected {
		prefix = indent + IconSelected + " "
		return SelectedStyle.Render(prefix + i.Label)
	}
	return ActionStyle.Render(prefix + i.Label)
}

func (i Item) renderToggle(selected bool, labelWidth int) string {
	indent := i.indentPrefix()

	var status string
	var style = MutedStyle

	switch i.State {
	case StateChecked:
		status = IconCheck
		style = SuccessStyle
	case StateUnchecked:
		status = "○"
		style = MutedStyle
	case StateLoading:
		status = "◌"
		style = ActionStyle
	default:
		status = " "
	}

	prefix := indent + "  "
	if selected {
		prefix = indent + IconSelected + " "
		style = SelectedStyle
	}

	if labelWidth > 0 && i.Description != "" {
		paddedLabel := fmt.Sprintf("%-*s", labelWidth, i.Label)
		return style.Render(fmt.Sprintf("%s%s %s", prefix, status, paddedLabel)) +
			MutedStyle.Render("  "+i.Description)
	}

	return style.Render(fmt.Sprintf("%s%s %s", prefix, status, i.Label))
}

func (i Item) renderExpandable(selected bool, labelWidth int) string {
	indent := i.indentPrefix()

	var arrow string
	if i.Expanded {
		arrow = IconArrowDown
	} else {
		arrow = IconArrowRight
	}

	prefix := indent + "  "
	if selected {
		prefix = indent + IconSelected + " "
	}

	paddedLabel := i.Label
	if labelWidth > 0 {
		paddedLabel = fmt.Sprintf("%-*s", labelWidth, i.Label)
	}

	if selected {
		base := SelectedStyle.Render(fmt.Sprintf("%s%s %s", prefix, arrow, paddedLabel))
		if i.Description != "" {
			return base + MutedStyle.Render("  "+i.Description)
		}
		return base
	}

	base := NormalStyle.Render(fmt.Sprintf("%s%s %s", prefix, arrow, paddedLabel))
	if i.Description != "" {
		return base + MutedStyle.Render("  "+i.Description)
	}
	return base
}

func (i Item) indentPrefix() string {
	if i.Indent <= 0 {
		return ""
	}
	return strings.Repeat(" ", i.Indent*IndentSpaces)
}
