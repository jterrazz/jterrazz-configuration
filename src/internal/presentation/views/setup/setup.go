package setup

import (
	"os/exec"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jterrazz/jterrazz-cli/src/internal/config"
	"github.com/jterrazz/jterrazz-cli/src/internal/presentation/components"
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

	// Calculate max description width for alignment
	maxDescWidth := 0
	for _, script := range config.Scripts {
		if script.CheckFn == nil {
			continue
		}
		if len(script.Description) > maxDescWidth {
			maxDescWidth = len(script.Description)
		}
	}

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
				DescWidth:   maxDescWidth,
			},
			name: script.Name,
		}

		if result.Installed {
			configuredItems = append(configuredItems, entry)
		} else {
			notConfiguredItems = append(notConfiguredItems, entry)
		}
	}

	sort.Slice(notConfiguredItems, func(i, j int) bool {
		return notConfiguredItems[i].name < notConfiguredItems[j].name
	})
	sort.Slice(configuredItems, func(i, j int) bool {
		return configuredItems[i].name < configuredItems[j].name
	})

	for _, entry := range notConfiguredItems {
		items = append(items, entry.item)
		itemNames = append(itemNames, entry.name)
	}
	for _, entry := range configuredItems {
		items = append(items, entry.item)
		itemNames = append(itemNames, entry.name)
	}

	// Scripts section - scripts without CheckFn (run-once actions)
	items = append(items, components.Item{Kind: components.KindHeader, Label: "Scripts"})
	itemNames = append(itemNames, "")

	var runOnceItems []scriptEntry
	for _, script := range config.Scripts {
		if script.CheckFn != nil {
			continue
		}

		runOnceItems = append(runOnceItems, scriptEntry{
			item: components.Item{
				Kind:        components.KindAction,
				Label:       script.Name,
				Description: script.Description,
			},
			name: script.Name,
		})
	}
	sort.Slice(runOnceItems, func(i, j int) bool {
		return runOnceItems[i].name < runOnceItems[j].name
	})
	for _, entry := range runOnceItems {
		items = append(items, entry.item)
		itemNames = append(itemNames, entry.name)
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
		// Check if script uses ExecArgs (needs full terminal control)
		if script := config.GetScriptByName(name); script != nil && len(script.ExecArgs) > 0 {
			c := exec.Command(script.ExecArgs[0], script.ExecArgs[1:]...)
			return tea.ExecProcess(c, func(err error) tea.Msg {
				return components.ActionDoneMsg{Message: "Completed " + name}
			})
		}

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
