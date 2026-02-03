package setup

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/presentation/components"
)

// Action represents navigation/action items in setup
type Action string

const (
	ActionSkills Action = "skills"
)

// itemNames maps list indices to script/action names
var itemNames []string
var loadingScript string // Script currently being run

// BuildItems builds the setup menu items
func BuildItems() []components.Item {
	var items []components.Item
	itemNames = []string{}

	// Navigation section
	items = append(items, components.Item{Kind: components.KindHeader, Label: "Navigation"})
	itemNames = append(itemNames, "")

	items = append(items, components.Item{Kind: components.KindNavigation, Label: "skills", Description: "Manage AI agent skills"})
	itemNames = append(itemNames, string(ActionSkills))

	// Configuration section - from config.Scripts with CheckFn
	items = append(items, components.Item{Kind: components.KindHeader, Label: "Setup"})
	itemNames = append(itemNames, "")

	type scriptEntry struct {
		item components.Item
		name string
	}
	var configuredItems, notConfiguredItems []scriptEntry

	for _, script := range config.Scripts {
		if script.CheckFn == nil {
			continue
		}

		result := script.CheckFn()
		state := components.StateUnchecked
		if loadingScript == script.Name {
			state = components.StateLoading
		} else if result.Installed {
			state = components.StateChecked
		}

		entry := scriptEntry{
			item: components.Item{
				Kind:        components.KindToggle,
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
	items = append(items, components.Item{Kind: components.KindHeader, Label: "Utilities"})
	itemNames = append(itemNames, "")

	for _, script := range config.Scripts {
		if script.CheckFn != nil {
			continue
		}

		items = append(items, components.Item{
			Kind:        components.KindAction,
			Label:       script.Name,
			Description: script.Description,
		})
		itemNames = append(itemNames, script.Name)
	}

	return items
}

// HandleSelect handles item selection in the setup menu
func HandleSelect(index int, item components.Item, runScript func(string)) tea.Cmd {
	if index >= len(itemNames) {
		return nil
	}
	name := itemNames[index]

	switch Action(name) {
	case ActionSkills:
		return func() tea.Msg {
			return components.NavigateMsg{
				InitFunc: InitSkillsState,
				Config:   SkillsConfig(),
			}
		}

	default:
		loadingScript = name
		return func() tea.Msg {
			runScript(name)
			loadingScript = ""
			return components.ActionDoneMsg{Message: "Completed " + name}
		}
	}
}

// RunOrExit runs the setup TUI
func RunOrExit(runScript func(string)) {
	components.RunOrExit(components.AppConfig{
		Title:      "Setup",
		BuildItems: BuildItems,
		OnSelect: func(index int, item components.Item) tea.Cmd {
			return HandleSelect(index, item, runScript)
		},
	})
}
