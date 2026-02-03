package setup

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/presentation/components/tui"
)

// Action represents navigation/action items in setup
type Action string

const (
	ActionSkills Action = "skills"
)

// itemNames maps list indices to script/action names
var itemNames []string

// BuildItems builds the setup menu items
func BuildItems() []tui.Item {
	var items []tui.Item
	itemNames = []string{}

	// Navigation section
	items = append(items, tui.Item{Kind: tui.KindHeader, Label: "Skills"})
	itemNames = append(itemNames, "")

	items = append(items, tui.Item{Kind: tui.KindAction, Label: "skills", Description: "Manage AI agent skills"})
	itemNames = append(itemNames, string(ActionSkills))

	// Configuration section - from config.Scripts with CheckFn
	items = append(items, tui.Item{Kind: tui.KindHeader, Label: "Setup"})
	itemNames = append(itemNames, "")

	type scriptEntry struct {
		item tui.Item
		name string
	}
	var configuredItems, notConfiguredItems []scriptEntry

	for _, script := range config.Scripts {
		if script.CheckFn == nil {
			continue
		}

		result := script.CheckFn()
		state := tui.StateUnchecked
		if result.Installed {
			state = tui.StateChecked
		}

		entry := scriptEntry{
			item: tui.Item{
				Kind:        tui.KindToggle,
				Label:       script.Name,
				Description: script.Description,
				State:       state,
			},
			name: script.Name,
		}

		if result.Installed {
			configuredItems = append(configuredItems, entry)
		} else {
			notConfiguredItems = append(notConfiguredItems, entry)
		}
	}

	for _, entry := range notConfiguredItems {
		items = append(items, entry.item)
		itemNames = append(itemNames, entry.name)
	}
	for _, entry := range configuredItems {
		items = append(items, entry.item)
		itemNames = append(itemNames, entry.name)
	}

	// Scripts section - scripts without CheckFn (utilities)
	items = append(items, tui.Item{Kind: tui.KindHeader, Label: "Utilities"})
	itemNames = append(itemNames, "")

	for _, script := range config.Scripts {
		if script.CheckFn != nil {
			continue
		}

		items = append(items, tui.Item{
			Kind:        tui.KindAction,
			Label:       script.Name,
			Description: script.Description,
		})
		itemNames = append(itemNames, script.Name)
	}

	return items
}

// HandleSelect handles item selection in the setup menu
func HandleSelect(index int, item tui.Item, runScript func(string)) tea.Cmd {
	if index >= len(itemNames) {
		return nil
	}
	name := itemNames[index]

	switch Action(name) {
	case ActionSkills:
		return func() tea.Msg {
			return tui.NavigateMsg{
				InitFunc: InitSkillsState,
				Config:   SkillsConfig(),
			}
		}

	default:
		return func() tea.Msg {
			runScript(name)
			return tui.ActionDoneMsg{Message: "Completed " + name}
		}
	}
}

// RunOrExit runs the setup TUI
func RunOrExit(runScript func(string)) {
	tui.RunOrExit(tui.AppConfig{
		Title:      "Setup",
		BuildItems: BuildItems,
		OnSelect: func(index int, item tui.Item) tea.Cmd {
			return HandleSelect(index, item, runScript)
		},
	})
}
