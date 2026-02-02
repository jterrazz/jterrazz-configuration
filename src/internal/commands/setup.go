package commands

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
	"github.com/jterrazz/jterrazz-cli/internal/ui/tui"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup system configurations (interactive)",
	Run: func(cmd *cobra.Command, args []string) {
		runSetupUI()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

// =============================================================================
// Script Runners
// =============================================================================

func runScript(name string) {
	script := config.GetScriptByName(name)
	if script == nil {
		ui.PrintError("Unknown script: " + name)
		return
	}

	if script.RunFn == nil {
		ui.PrintError("No runner for script: " + name)
		return
	}

	if err := script.RunFn(); err != nil {
		ui.PrintError("Failed to run " + name + ": " + err.Error())
	}
}

// runSetupItem runs a setup item by name (used by install command for Tool.Scripts)
func runSetupItem(name string) {
	runScript(name)
}

// =============================================================================
// Setup Data
// =============================================================================

// setupAction represents navigation/action items in setup
type setupAction string

const (
	setupActionSkills setupAction = "skills"
)

// setupItemNames maps list indices to script/action names
var setupItemNames []string

func buildSetupItems() []tui.Item {
	var items []tui.Item
	setupItemNames = []string{}

	// Navigation section
	items = append(items, tui.Item{Kind: tui.KindHeader, Label: "Navigation"})
	setupItemNames = append(setupItemNames, "")

	items = append(items, tui.Item{Kind: tui.KindAction, Label: "skills", Description: "Manage AI agent skills"})
	setupItemNames = append(setupItemNames, string(setupActionSkills))

	// Configuration section - from config.Scripts with CheckFn
	items = append(items, tui.Item{Kind: tui.KindHeader, Label: "Configuration"})
	setupItemNames = append(setupItemNames, "")

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
		setupItemNames = append(setupItemNames, entry.name)
	}
	for _, entry := range configuredItems {
		items = append(items, entry.item)
		setupItemNames = append(setupItemNames, entry.name)
	}

	// Scripts section - scripts without CheckFn (utilities)
	items = append(items, tui.Item{Kind: tui.KindHeader, Label: "Scripts"})
	setupItemNames = append(setupItemNames, "")

	for _, script := range config.Scripts {
		if script.CheckFn != nil {
			continue
		}

		items = append(items, tui.Item{
			Kind:        tui.KindAction,
			Label:       script.Name,
			Description: script.Description,
		})
		setupItemNames = append(setupItemNames, script.Name)
	}

	return items
}

func handleSetupSelect(index int, item tui.Item) tea.Cmd {
	if index >= len(setupItemNames) {
		return nil
	}
	name := setupItemNames[index]

	switch setupAction(name) {
	case setupActionSkills:
		return func() tea.Msg {
			return tui.NavigateMsg{
				InitFunc: initSkillsState,
				Config: tui.AppConfig{
					Title:      "Skills",
					BuildItems: buildSkillsItems,
					OnSelect:   handleSkillsSelect,
					OnMessage:  handleSkillsMessage,
				},
			}
		}

	default:
		return func() tea.Msg {
			runScript(name)
			return tui.ActionDoneMsg{Message: "Completed " + name}
		}
	}
}

func runSetupUI() {
	tui.RunOrExit(tui.AppConfig{
		Title:      "Setup",
		BuildItems: buildSetupItems,
		OnSelect:   handleSetupSelect,
	})
}
