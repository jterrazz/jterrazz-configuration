package status

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/domain/status"
	"github.com/jterrazz/jterrazz-cli/internal/presentation/components"
	"github.com/jterrazz/jterrazz-cli/internal/presentation/theme"
)

// ProcessRefreshMsg triggers a refresh of process data
type ProcessRefreshMsg struct{}

// ProcessDataMsg carries refreshed process data
type ProcessDataMsg struct {
	Data map[string][]config.ProcessInfo
}

// Model is the Bubble Tea model for the status view
type Model struct {
	loader    *status.Loader
	items     map[string]status.Item
	itemOrder []status.Item
	spinner   spinner.Model
	viewport  viewport.Model
	ready     bool
	width     int
	height    int
	loaded    int
	total     int
	quitting  bool
	allLoaded bool
}

// New creates a new status view model
func New() Model {
	loader := status.NewLoader()
	items := make(map[string]status.Item)
	itemOrder := loader.GetItems()

	total := 0
	for _, item := range itemOrder {
		items[item.ID] = item
		if !item.Loaded && item.Kind != status.KindHeader {
			total++
		}
	}

	return Model{
		loader:    loader,
		items:     items,
		itemOrder: itemOrder,
		spinner:   components.NewSpinnerModel(),
		total:     total,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	m.loader.Start()
	return tea.Batch(
		m.spinner.Tick,
		m.loader.WaitForUpdate(),
		scheduleProcessRefresh(),
	)
}

// scheduleProcessRefresh returns a command that triggers a process refresh after 1 second
func scheduleProcessRefresh() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return ProcessRefreshMsg{}
	})
}

// refreshProcesses runs process checks in background and returns the data
func refreshProcesses() tea.Cmd {
	return func() tea.Msg {
		data := make(map[string][]config.ProcessInfo)
		for _, check := range config.ProcessChecks {
			data["process-"+check.Name] = check.CheckFn()
		}
		return ProcessDataMsg{Data: data}
	}
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "home", "g":
			m.viewport.GotoTop()
		case "end", "G":
			m.viewport.GotoBottom()
		}

	case tea.WindowSizeMsg:
		headerHeight := 5 // blank + title + sysinfo + blank + newline
		footerHeight := 1

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-headerHeight-footerHeight)
			m.viewport.YPosition = headerHeight
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - headerHeight - footerHeight
		}
		m.width = msg.Width
		m.height = msg.Height

	case status.UpdateMsg:
		if existing, ok := m.items[msg.ID]; ok {
			existing.Loaded = msg.Item.Loaded
			existing.Installed = msg.Item.Installed
			existing.Version = msg.Item.Version
			existing.Status = msg.Item.Status
			existing.Detail = msg.Item.Detail
			existing.Value = msg.Item.Value
			existing.Style = msg.Item.Style
			existing.Available = msg.Item.Available
			existing.Processes = msg.Item.Processes
			m.items[msg.ID] = existing
		} else {
			m.items[msg.ID] = msg.Item
		}
		if msg.Item.Loaded {
			m.loaded++
		}
		cmds = append(cmds, m.loader.WaitForUpdate())

	case status.AllLoadedMsg:
		m.allLoaded = true

	case ProcessRefreshMsg:
		// Trigger async process data refresh
		cmds = append(cmds, refreshProcesses(), scheduleProcessRefresh())
		return m, tea.Batch(cmds...) // Don't re-render yet

	case ProcessDataMsg:
		// Apply refreshed process data
		for id, processes := range msg.Data {
			if existing, ok := m.items[id]; ok {
				existing.Processes = processes
				existing.Available = len(processes) > 0
				m.items[id] = existing
			}
		}

	case spinner.TickMsg:
		if !m.allLoaded {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
		// Only re-render content on spinner tick if still loading
		if m.ready && !m.allLoaded {
			m.viewport.SetContent(m.renderContent())
		}
		return m, tea.Batch(cmds...)
	}

	// Re-render content for data updates (UpdateMsg, ProcessDataMsg, etc.)
	if m.ready {
		m.viewport.SetContent(m.renderContent())
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View implements tea.Model
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	if !m.ready {
		return m.spinner.View() + " Initializing..."
	}

	var b strings.Builder

	// Header
	title := theme.SectionTitle.Render("STATUS")

	// System info
	sysInfo := ""
	if sysinfo, ok := m.items["sysinfo"]; ok && sysinfo.Loaded {
		sysInfo = theme.Muted.Render(sysinfo.Detail)
	} else {
		sysInfo = m.spinner.View() + " Loading..."
	}

	header := lipgloss.JoinVertical(lipgloss.Left,
		"",
		components.PageIndent+title,
		components.PageIndent+sysInfo,
		"",
	)
	b.WriteString(header)
	b.WriteString("\n")

	// Content
	b.WriteString(m.viewport.View())
	b.WriteString("\n")

	// Footer
	scrollPercent := int(m.viewport.ScrollPercent() * 100)
	help := fmt.Sprintf("↑/↓ scroll • g/G top/bottom • %d%% • q quit", scrollPercent)

	if m.allLoaded {
		footer := theme.Help.Render(help) + components.ColumnSeparator + theme.Success.Render(theme.IconCheck+" All checks complete")
		b.WriteString(footer)
	} else {
		progressBar := m.renderProgressBar()
		footer := theme.Help.Render(help) + components.ColumnSeparator + progressBar
		b.WriteString(footer)
	}

	return b.String()
}

func (m Model) renderProgressBar() string {
	if m.allLoaded {
		return theme.Success.Render(theme.IconCheck + " All checks complete")
	}

	width := 30
	filled := int(float64(m.loaded) / float64(m.total) * float64(width))
	if filled > width {
		filled = width
	}

	bar := theme.ProgressFilled.Render(strings.Repeat(theme.IconProgressFull, filled)) +
		theme.ProgressEmpty.Render(strings.Repeat(theme.IconProgressEmpty, width-filled))

	return fmt.Sprintf("%s %s %d/%d",
		m.spinner.View(),
		bar,
		m.loaded,
		m.total,
	)
}

// Run starts the status TUI
func Run() error {
	p := tea.NewProgram(New(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// RunOrExit runs the status TUI and exits on error
func RunOrExit() {
	if err := Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// Helper to get visible length (strip ANSI)
func visibleLen(s string) int {
	return components.VisibleLen(s)
}
