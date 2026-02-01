package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// AppConfig defines the configuration for a TUI application
type AppConfig struct {
	Title       string
	BuildItems  func() []Item
	OnSelect    func(index int, item Item) tea.Cmd
	OnMessage   func(msg tea.Msg) tea.Cmd // Custom message handler
	OnBack      func() bool               // Return true to handle back, false to quit
	Breadcrumbs []string
}

// ActionDoneMsg signals that an async action completed
type ActionDoneMsg struct {
	Message string
	Err     error
}

// RefreshMsg signals that items should be rebuilt
type RefreshMsg struct{}

// App is a generic TUI application model
type App struct {
	config     AppConfig
	list       *List
	page       *Page
	processing bool
	quitting   bool
}

// NewApp creates a new TUI application
func NewApp(config AppConfig) *App {
	items := config.BuildItems()
	list := NewList(items)
	list.CalculateLabelWidth()

	page := NewPage(config.Title)
	page.Breadcrumbs = config.Breadcrumbs

	return &App{
		config: config,
		list:   list,
		page:   page,
	}
}

// Init implements tea.Model
func (a *App) Init() tea.Cmd {
	return tea.WindowSize()
}

// Update implements tea.Model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if a.processing {
			return a, nil
		}

		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			if a.config.OnBack != nil && a.config.OnBack() {
				return a, nil
			}
			a.quitting = true
			return a, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))):
			a.quitting = true
			return a, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			a.list.Up()

		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			a.list.Down()

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			return a.handleSelect()
		}

	case ActionDoneMsg:
		a.processing = false
		a.page.Message = msg.Message
		a.page.Processing = false
		a.rebuildItems()
		return a, nil

	case RefreshMsg:
		a.rebuildItems()
		return a, nil

	case tea.WindowSizeMsg:
		a.list.SetSize(msg.Width, msg.Height)
		a.page.SetSize(msg.Width, msg.Height)

	default:
		// Custom message handler
		if a.config.OnMessage != nil {
			cmd := a.config.OnMessage(msg)
			a.rebuildItems()
			return a, cmd
		}
	}

	return a, nil
}

// View implements tea.Model
func (a *App) View() string {
	if a.quitting {
		return ""
	}

	if len(a.config.Breadcrumbs) > 0 {
		a.page.Help = DefaultHelpWithBack()
	} else {
		a.page.Help = DefaultHelp()
	}

	a.page.Content = a.list.Render(a.page.ContentHeight())
	return a.page.Render()
}

func (a *App) handleSelect() (*App, tea.Cmd) {
	idx := a.list.SelectedIndex()
	if idx < 0 || idx >= len(a.list.Items) {
		return a, nil
	}

	item := a.list.Items[idx]
	if item.Kind == KindHeader {
		return a, nil
	}

	if a.config.OnSelect != nil {
		cmd := a.config.OnSelect(idx, item)
		if cmd != nil {
			a.processing = true
			a.page.Processing = true
			return a, cmd
		}
	}

	return a, nil
}

func (a *App) rebuildItems() {
	cursor := a.list.Cursor
	items := a.config.BuildItems()
	a.list = NewList(items)
	a.list.CalculateLabelWidth()

	if cursor >= len(items) {
		cursor = len(items) - 1
	}
	a.list.SetCursor(cursor)

	for a.list.Cursor > 0 && !a.list.Items[a.list.Cursor].Selectable() {
		a.list.Cursor--
	}
}

// SetMessage sets a status message
func (a *App) SetMessage(msg string) {
	a.page.Message = msg
}

// SetProcessing sets the processing state
func (a *App) SetProcessing(processing bool, msg string) {
	a.processing = processing
	a.page.Processing = processing
	a.page.Message = msg
}

// GetList returns the list for direct manipulation
func (a *App) GetList() *List {
	return a.list
}

// GetPage returns the page for direct manipulation
func (a *App) GetPage() *Page {
	return a.page
}

// Run starts the TUI application
func Run(config AppConfig) error {
	app := NewApp(config)
	p := tea.NewProgram(app, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// RunOrExit runs the TUI and exits on error
func RunOrExit(config AppConfig) {
	if err := Run(config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
