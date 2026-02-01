package commands

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
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
		ui.PrintError(fmt.Sprintf("Unknown script: %s", name))
		return
	}

	if script.RunFn == nil {
		ui.PrintError(fmt.Sprintf("No runner for script: %s", name))
		return
	}

	if err := script.RunFn(); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to run %s: %v", name, err))
	}
}

// runSetupItem runs a setup item by name (used by install command for Tool.Scripts)
func runSetupItem(name string) {
	runScript(name)
}

// =============================================================================
// Setup Data
// =============================================================================

// setupItemNames maps list indices to script/action names
var setupItemNames []string

func buildSetupItems() []ui.Item {
	var items []ui.Item
	setupItemNames = []string{}

	// Navigation section
	items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Navigation"})
	setupItemNames = append(setupItemNames, "")

	items = append(items, ui.Item{Kind: ui.KindAction, Label: "skills", Description: "Manage AI agent skills"})
	setupItemNames = append(setupItemNames, "skills")

	// Actions section
	items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Actions"})
	setupItemNames = append(setupItemNames, "")

	items = append(items, ui.Item{Kind: ui.KindAction, Label: "Setup all missing"})
	setupItemNames = append(setupItemNames, "setup-missing")

	// Configuration section - from config.Scripts with CheckFn
	items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Configuration"})
	setupItemNames = append(setupItemNames, "")

	type scriptEntry struct {
		item ui.Item
		name string
	}
	var configuredItems, notConfiguredItems []scriptEntry

	for _, script := range config.Scripts {
		if script.CheckFn == nil {
			continue
		}

		result := script.CheckFn()
		state := ui.StateUnchecked
		if result.Installed {
			state = ui.StateChecked
		}

		entry := scriptEntry{
			item: ui.Item{
				Kind:        ui.KindToggle,
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
	items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Scripts"})
	setupItemNames = append(setupItemNames, "")

	for _, script := range config.Scripts {
		if script.CheckFn != nil {
			continue
		}

		items = append(items, ui.Item{
			Kind:        ui.KindAction,
			Label:       script.Name,
			Description: script.Description,
		})
		setupItemNames = append(setupItemNames, script.Name)
	}

	return items
}

func handleSetupSelect(index int, item ui.Item) tea.Cmd {
	if index >= len(setupItemNames) {
		return nil
	}
	name := setupItemNames[index]

	switch name {
	case "skills":
		// Open skills sub-view - handled specially
		return nil

	case "setup-missing":
		return func() tea.Msg {
			count := 0
			for _, script := range config.Scripts {
				if script.CheckFn != nil {
					result := script.CheckFn()
					if !result.Installed {
						runScript(script.Name)
						count++
					}
				}
			}
			if count == 0 {
				return ui.ActionDoneMsg{Message: "Everything already configured"}
			}
			return ui.ActionDoneMsg{Message: fmt.Sprintf("Configured %d items", count)}
		}

	default:
		return func() tea.Msg {
			runScript(name)
			return ui.ActionDoneMsg{Message: fmt.Sprintf("Completed %s", name)}
		}
	}
}

func runSetupUI() {
	ui.RunOrExit(ui.AppConfig{
		Title:      "Setup",
		BuildItems: buildSetupItems,
		OnSelect:   handleSetupSelect,
	})
}
