package commands

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
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

// =============================================================================
// Helpers
// =============================================================================

// runSetupItem runs a setup item by name (used by install command for Tool.Scripts)
func runSetupItem(name string) {
	runScript(name)
}

// =============================================================================
// Setup TUI
// =============================================================================

type setupItemData struct {
	name       string
	configured *bool
}

type setupModel struct {
	list       *ui.List
	page       *ui.Page
	itemData   []setupItemData
	processing bool
	quitting   bool
	// Skills sub-view
	showSkills  bool
	skillsModel *skillsModel
}

func (m setupModel) buildItems() ([]ui.Item, []setupItemData) {
	var items []ui.Item
	var data []setupItemData

	// Navigation section
	items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Navigation"})
	data = append(data, setupItemData{})

	items = append(items, ui.Item{Kind: ui.KindAction, Label: "skills", Description: "Manage AI agent skills"})
	data = append(data, setupItemData{name: "skills"})

	// Actions section
	items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Actions"})
	data = append(data, setupItemData{})

	items = append(items, ui.Item{Kind: ui.KindAction, Label: "Setup all missing"})
	data = append(data, setupItemData{name: "setup-missing"})

	// Configuration section - from config.Scripts with CheckFn
	items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Configuration"})
	data = append(data, setupItemData{})

	var configuredItems []struct {
		item ui.Item
		data setupItemData
	}
	var notConfiguredItems []struct {
		item ui.Item
		data setupItemData
	}

	for _, script := range config.Scripts {
		if script.CheckFn == nil {
			continue // Skip scripts without check (utilities)
		}

		result := script.CheckFn()
		state := ui.StateUnchecked
		configured := result.Installed
		if configured {
			state = ui.StateChecked
		}

		entry := struct {
			item ui.Item
			data setupItemData
		}{
			item: ui.Item{
				Kind:        ui.KindToggle,
				Label:       script.Name,
				Description: script.Description,
				State:       state,
			},
			data: setupItemData{name: script.Name, configured: &configured},
		}

		if configured {
			configuredItems = append(configuredItems, entry)
		} else {
			notConfiguredItems = append(notConfiguredItems, entry)
		}
	}

	// Add pending items first, then configured
	for _, entry := range notConfiguredItems {
		items = append(items, entry.item)
		data = append(data, entry.data)
	}
	for _, entry := range configuredItems {
		items = append(items, entry.item)
		data = append(data, entry.data)
	}

	// Scripts section - scripts without CheckFn (utilities)
	items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Scripts"})
	data = append(data, setupItemData{})

	for _, script := range config.Scripts {
		if script.CheckFn != nil {
			continue // Skip scripts with check (configurations)
		}

		items = append(items, ui.Item{
			Kind:        ui.KindAction,
			Label:       script.Name,
			Description: script.Description,
		})
		data = append(data, setupItemData{name: script.Name})
	}

	return items, data
}

func (m setupModel) rebuildItems() setupModel {
	cursor := m.list.Cursor
	items, data := m.buildItems()
	m.list = ui.NewList(items)
	m.list.CalculateLabelWidth()
	m.itemData = data

	if cursor >= len(items) {
		cursor = len(items) - 1
	}
	m.list.SetCursor(cursor)

	for m.list.Cursor > 0 && !m.list.Items[m.list.Cursor].Selectable() {
		m.list.Cursor--
	}
	return m
}

func newSetupModel() setupModel {
	m := setupModel{
		page: ui.NewPage("Setup"),
	}

	items, data := m.buildItems()
	m.list = ui.NewList(items)
	m.list.CalculateLabelWidth()
	m.itemData = data

	return m
}

func (m setupModel) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m setupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If showing skills sub-view, delegate to it
	if m.showSkills && m.skillsModel != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, key.NewBinding(key.WithKeys("esc"))) {
				m.showSkills = false
				m.skillsModel = nil
				return m, nil
			}
			if key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))) {
				m.quitting = true
				return m, tea.Quit
			}
		case tea.WindowSizeMsg:
			m.list.SetSize(msg.Width, msg.Height)
			m.page.SetSize(msg.Width, msg.Height)
			m.skillsModel.list.SetSize(msg.Width, msg.Height)
			m.skillsModel.page.SetSize(msg.Width, msg.Height)
		}
		newModel, cmd := m.skillsModel.Update(msg)
		if sm, ok := newModel.(skillsModel); ok {
			m.skillsModel = &sm
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.processing {
			return m, nil
		}

		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"))):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			m.list.Up()

		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			m.list.Down()

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			return m.handleSelect()
		}

	case setupActionDoneMsg:
		m.processing = false
		m.page.Message = msg.message
		m.page.Processing = false
		m = m.rebuildItems()
		return m, nil

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		m.page.SetSize(msg.Width, msg.Height)
	}

	return m, nil
}

type setupActionDoneMsg struct {
	message string
	err     error
}

func (m setupModel) handleSelect() (setupModel, tea.Cmd) {
	idx := m.list.SelectedIndex()
	if idx < 0 || idx >= len(m.itemData) {
		return m, nil
	}

	item := m.list.Selected()
	data := m.itemData[idx]

	if item.Kind == ui.KindHeader {
		return m, nil
	}

	switch data.name {
	case "skills":
		sm := newSkillsModel()
		sm.list.SetSize(m.page.Width, m.page.Height)
		sm.page.SetSize(m.page.Width, m.page.Height)
		m.skillsModel = &sm
		m.showSkills = true
		return m, nil

	case "setup-missing":
		m.processing = true
		m.page.Processing = true
		m.page.Message = "Setting up all missing..."
		return m, m.runSetupMissing()

	default:
		m.processing = true
		m.page.Processing = true
		m.page.Message = fmt.Sprintf("Running %s...", data.name)
		return m, m.runSetup(data.name)
	}
}

func (m setupModel) runSetupMissing() tea.Cmd {
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
			return setupActionDoneMsg{message: "Everything already configured"}
		}
		return setupActionDoneMsg{message: fmt.Sprintf("Configured %d items", count)}
	}
}

func (m setupModel) runSetup(name string) tea.Cmd {
	return func() tea.Msg {
		runScript(name)
		return setupActionDoneMsg{message: fmt.Sprintf("Completed %s", name)}
	}
}

func (m setupModel) View() string {
	if m.quitting {
		return ""
	}

	if m.showSkills && m.skillsModel != nil {
		return m.skillsModel.viewWithBreadcrumb("Setup", "Skills")
	}

	m.page.Help = ui.DefaultHelp()
	m.page.Content = m.list.Render(m.page.ContentHeight())

	return m.page.Render()
}

func runSetupUI() {
	m := newSetupModel()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running setup UI: %v\n", err)
		os.Exit(1)
	}
}
