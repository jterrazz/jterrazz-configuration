package commands

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show comprehensive system status",
	Run: func(cmd *cobra.Command, args []string) {
		showStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

// =============================================================================
// Styles
// =============================================================================

var (
	// Section title embedded in border (main sections like SYSTEM, TOOLS)
	sectionBorderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")) // dark gray

	sectionTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ui.ColorPrimary)).
				Bold(true)

	// Box styles
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	subsectionTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ui.ColorText))

	// Table styles
	cellStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ui.ColorText))

	// Status badge styles
	badgeOk = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ui.ColorSuccess)).
		Bold(true)

	badgeErr = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ui.ColorDanger)).
			Bold(true)

	badgeLoading = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ui.ColorWarning))

	// Progress bar
	progressFull = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39"))

	progressEmpty = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ui.ColorBorder))

	// Service status
	serviceRunning = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ui.ColorSuccess)).
			Render("●")

	serviceStopped = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ui.ColorWarning)).
			Render("○")
)

// =============================================================================
// Status TUI Model
// =============================================================================

type statusModel struct {
	loader    *StatusLoader
	items     map[string]StatusItem
	itemOrder []StatusItem
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

func newStatusModel() statusModel {
	loader := NewStatusLoader()
	items := make(map[string]StatusItem)
	itemOrder := loader.GetItems()

	total := 0
	for _, item := range itemOrder {
		items[item.ID] = item
		if !item.Loaded && item.Kind != StatusItemHeader {
			total++
		}
	}

	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		FPS:    80 * time.Millisecond,
	}
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))

	return statusModel{
		loader:    loader,
		items:     items,
		itemOrder: itemOrder,
		spinner:   s,
		total:     total,
	}
}

func (m statusModel) Init() tea.Cmd {
	m.loader.Start()
	return tea.Batch(
		m.spinner.Tick,
		m.loader.WaitForUpdate(),
	)
}

func (m statusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case StatusUpdateMsg:
		if existing, ok := m.items[msg.ID]; ok {
			existing.Loaded = msg.Item.Loaded
			existing.Installed = msg.Item.Installed
			existing.Version = msg.Item.Version
			existing.Status = msg.Item.Status
			existing.Detail = msg.Item.Detail
			existing.Value = msg.Item.Value
			existing.Style = msg.Item.Style
			existing.Available = msg.Item.Available
			m.items[msg.ID] = existing
		} else {
			m.items[msg.ID] = msg.Item
		}
		if msg.Item.Loaded {
			m.loaded++
		}
		cmds = append(cmds, m.loader.WaitForUpdate())

	case AllLoadedMsg:
		m.allLoaded = true

	case spinner.TickMsg:
		if !m.allLoaded {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	if m.ready {
		m.viewport.SetContent(m.renderContent())
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m statusModel) View() string {
	if m.quitting {
		return ""
	}

	if !m.ready {
		return m.spinner.View() + " Initializing..."
	}

	var b strings.Builder

	// Header
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ui.ColorPrimary)).
		Bold(true).
		Render("STATUS")

	// System info
	sysInfo := ""
	if sysinfo, ok := m.items["sysinfo"]; ok && sysinfo.Loaded {
		sysInfo = ui.MutedStyle.Render(sysinfo.Detail)
	} else {
		sysInfo = m.spinner.View() + " Loading..."
	}

	header := lipgloss.JoinVertical(lipgloss.Left,
		"",
		"  "+title,
		"  "+sysInfo,
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
		footer := ui.HelpStyle.Render(help) + "  " + ui.SuccessStyle.Render("✓ All checks complete")
		b.WriteString(footer)
	} else {
		progressBar := m.renderProgressBar()
		footer := ui.HelpStyle.Render(help) + "  " + progressBar
		b.WriteString(footer)
	}

	return b.String()
}

func (m statusModel) renderProgressBar() string {
	if m.allLoaded {
		return ui.SuccessStyle.Render("✓ All checks complete")
	}

	width := 30
	filled := int(float64(m.loaded) / float64(m.total) * float64(width))
	if filled > width {
		filled = width
	}

	bar := progressFull.Render(strings.Repeat("█", filled)) +
		progressEmpty.Render(strings.Repeat("░", width-filled))

	return fmt.Sprintf("%s %s %d/%d",
		m.spinner.View(),
		bar,
		m.loaded,
		m.total,
	)
}

func (m statusModel) renderContent() string {
	var b strings.Builder

	sections := m.groupBySection()
	boxWidth := m.width - 4
	if boxWidth < 40 {
		boxWidth = 40
	}

	for _, section := range []string{"System", "Tools", "Resources"} {
		subsections, ok := sections[section]
		if !ok {
			continue
		}

		// Section header with decorative line
		b.WriteString("\n")
		sectionHeader := m.renderSectionHeader(section, boxWidth)
		b.WriteString(sectionHeader)
		b.WriteString("\n")

		// Collect all items in this section for column width calculation
		var allSectionItems []StatusItem
		for _, subsection := range m.getSubsectionOrder(section) {
			items, ok := subsections[subsection]
			if !ok {
				continue
			}
			for _, item := range items {
				if item.Kind == StatusItemNetwork || item.Kind == StatusItemDisk || item.Kind == StatusItemCache {
					if !item.Loaded || item.Available {
						allSectionItems = append(allSectionItems, item)
					}
				} else {
					allSectionItems = append(allSectionItems, item)
				}
			}
		}

		// Calculate column widths for the entire section
		colWidths := m.calculateColumnWidths(allSectionItems)

		// Render subsections as boxes
		for _, subsection := range m.getSubsectionOrder(section) {
			items, ok := subsections[subsection]
			if !ok {
				continue
			}

			// Filter out unavailable items
			var visibleItems []StatusItem
			for _, item := range items {
				if item.Kind == StatusItemNetwork || item.Kind == StatusItemDisk || item.Kind == StatusItemCache {
					if !item.Loaded || item.Available {
						visibleItems = append(visibleItems, item)
					}
				} else {
					visibleItems = append(visibleItems, item)
				}
			}

			if len(visibleItems) == 0 {
				continue
			}

			box := m.renderSubsectionBox(subsection, visibleItems, boxWidth, colWidths)
			b.WriteString(box)
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (m statusModel) renderSectionHeader(title string, width int) string {
	// Box header style:
	// ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
	// ┃  SYSTEM                                                                ┃
	// ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
	innerWidth := width - 2 // for the borders
	if innerWidth < 10 {
		innerWidth = 10
	}

	icon := m.getSectionIcon(title)
	displayTitle := icon + strings.ToUpper(title)
	padding := innerWidth - len(displayTitle) - 2 // -2 for "  " prefix
	if padding < 0 {
		padding = 0
	}

	top := sectionBorderStyle.Render("┏" + strings.Repeat("━", innerWidth) + "┓")
	middle := sectionBorderStyle.Render("┃") + "  " + sectionTitleStyle.Render(displayTitle) + strings.Repeat(" ", padding) + sectionBorderStyle.Render("┃")
	bottom := sectionBorderStyle.Render("┗" + strings.Repeat("━", innerWidth) + "┛")

	return top + "\n" + middle + "\n" + bottom
}

func (m statusModel) getSectionIcon(section string) string {
	return ""
}

func (m statusModel) groupBySection() map[string]map[string][]StatusItem {
	sections := make(map[string]map[string][]StatusItem)

	for _, baseItem := range m.itemOrder {
		item := m.items[baseItem.ID]
		if item.Kind == StatusItemHeader || item.Kind == StatusItemSystemInfo {
			continue
		}

		if sections[item.Section] == nil {
			sections[item.Section] = make(map[string][]StatusItem)
		}
		sections[item.Section][item.SubSection] = append(sections[item.Section][item.SubSection], item)
	}

	// Sort items A-Z by name within each subsection
	for _, subsections := range sections {
		for subsection, items := range subsections {
			sort.Slice(items, func(i, j int) bool {
				return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
			})
			subsections[subsection] = items
		}
	}

	return sections
}

func (m statusModel) getSubsectionOrder(section string) []string {
	switch section {
	case "System":
		return []string{"Setup", "MacOS Security", "Identity"}
	case "Tools":
		return []string{"Package Managers", "Languages", "Infrastructure", "AI", "Apps", "System Tools"}
	case "Resources":
		return []string{"Network", "Disk Usage", "Caches & Cleanable"}
	}
	return nil
}

func (m statusModel) getSubsectionIcon(subsection string) string {
	return ""
}

// ColumnWidths holds calculated column widths for alignment
type ColumnWidths struct {
	Name    int
	Desc    int
	Version int
	Method  int
	Detail  int
}

func (m statusModel) calculateColumnWidths(items []StatusItem) ColumnWidths {
	widths := ColumnWidths{}
	for _, item := range items {
		if len(item.Name) > widths.Name {
			widths.Name = len(item.Name)
		}
		if len(item.Description) > widths.Desc {
			widths.Desc = len(item.Description)
		}
		if len(item.Version) > widths.Version {
			widths.Version = len(item.Version)
		}
		if len(item.Method) > widths.Method {
			widths.Method = len(item.Method)
		}
		if len(item.Detail) > widths.Detail {
			widths.Detail = len(item.Detail)
		}
	}
	return widths
}

func (m statusModel) renderSubsectionBox(title string, items []StatusItem, width int, colWidths ColumnWidths) string {
	var content strings.Builder

	// Table content (no title inside, it goes in the border)
	for i, item := range items {
		row := m.renderTableRow(item, colWidths)
		content.WriteString(row)
		if i < len(items)-1 {
			content.WriteString("\n")
		}
	}

	// Custom border with title embedded in top
	innerWidth := width - 4 // account for border + padding
	if innerWidth < 20 {
		innerWidth = 20
	}

	// Build the box manually with title: ╭─ 🌐 Title ─────────────────╮
	icon := m.getSubsectionIcon(title)
	displayTitle := icon + title
	titleLen := len(title) + len(icon) + 4 // "─ " + icon + title + " ─"
	topBorderRight := strings.Repeat("─", innerWidth-titleLen)
	if topBorderRight == "" {
		topBorderRight = "─"
	}

	top := sectionBorderStyle.Render("╭─ ") + subsectionTitleStyle.Render(displayTitle) + sectionBorderStyle.Render(" ─"+topBorderRight+"╮")
	bottom := sectionBorderStyle.Render("╰" + strings.Repeat("─", innerWidth+2) + "╯")

	// Pad content lines
	lines := strings.Split(content.String(), "\n")
	var paddedLines []string
	for _, line := range lines {
		// Calculate visible length (approximate, ignoring ANSI codes for padding)
		padding := innerWidth - visibleLen(line)
		if padding < 0 {
			padding = 0
		}
		paddedLine := sectionBorderStyle.Render("│ ") + line + strings.Repeat(" ", padding) + sectionBorderStyle.Render(" │")
		paddedLines = append(paddedLines, paddedLine)
	}

	return top + "\n" + strings.Join(paddedLines, "\n") + "\n" + bottom
}

func (m statusModel) renderTableRow(item StatusItem, colWidths ColumnWidths) string {
	if !item.Loaded {
		switch item.Kind {
		case StatusItemSetup:
			return m.renderSetupRowLoading(item, colWidths)
		case StatusItemSecurity, StatusItemIdentity:
			return m.renderCheckRowLoading(item, colWidths)
		case StatusItemTool:
			return m.renderToolRowLoading(item, colWidths)
		case StatusItemNetwork, StatusItemDisk, StatusItemCache:
			return m.renderResourceRowLoading(item, colWidths)
		default:
			return fmt.Sprintf("  %s  %-*s", m.spinner.View(), colWidths.Name, item.Name)
		}
	}

	switch item.Kind {
	case StatusItemSetup:
		return m.renderSetupRow(item, colWidths)
	case StatusItemSecurity, StatusItemIdentity:
		return m.renderCheckRow(item, colWidths)
	case StatusItemTool:
		return m.renderToolRow(item, colWidths)
	case StatusItemNetwork, StatusItemDisk, StatusItemCache:
		return m.renderResourceRow(item, colWidths)
	}

	return ""
}

func (m statusModel) renderSetupRowLoading(item StatusItem, colWidths ColumnWidths) string {
	name := cellStyle.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))
	desc := ui.MutedStyle.Render(fmt.Sprintf("%-*s", colWidths.Desc, item.Description))
	return fmt.Sprintf("  %s  %s  %s", name, desc, m.spinner.View())
}

func (m statusModel) renderSetupRow(item StatusItem, colWidths ColumnWidths) string {
	name := cellStyle.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))
	desc := ui.MutedStyle.Render(fmt.Sprintf("%-*s", colWidths.Desc, item.Description))
	badge := m.statusBadge(item.Installed)
	detail := ""
	if item.Detail != "" {
		detail = ui.SpecialStyle.Render(item.Detail)
	}
	return fmt.Sprintf("  %s  %s  %s  %s", name, desc, badge, detail)
}

func (m statusModel) renderCheckRowLoading(item StatusItem, colWidths ColumnWidths) string {
	name := cellStyle.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))
	desc := ui.MutedStyle.Render(fmt.Sprintf("%-*s", colWidths.Desc, item.Description))
	return fmt.Sprintf("  %s  %s  %s", name, desc, m.spinner.View())
}

func (m statusModel) renderCheckRow(item StatusItem, colWidths ColumnWidths) string {
	name := cellStyle.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))
	desc := ui.MutedStyle.Render(fmt.Sprintf("%-*s", colWidths.Desc, item.Description))
	ok := item.Installed == item.GoodWhen
	badge := m.statusBadge(ok)
	detail := ""
	if item.Detail != "" {
		detail = ui.SpecialStyle.Render(item.Detail)
	}
	return fmt.Sprintf("  %s  %s  %s  %s", name, desc, badge, detail)
}

func (m statusModel) renderToolRowLoading(item StatusItem, colWidths ColumnWidths) string {
	name := cellStyle.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))
	method := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(fmt.Sprintf("%-*s", colWidths.Method, item.Method))
	return fmt.Sprintf("  %s  %s  %s", name, method, m.spinner.View())
}

func (m statusModel) renderToolRow(item StatusItem, colWidths ColumnWidths) string {
	name := cellStyle.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))

	// Method very dim (least important)
	method := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(fmt.Sprintf("%-*s", colWidths.Method, item.Method))

	badge := m.statusBadge(item.Installed)

	// Version in cyan if present, dimmed dash if not
	versionStr := fmt.Sprintf("%-*s", colWidths.Version, item.Version)
	var version string
	if item.Version != "" {
		version = ui.SpecialStyle.Render(versionStr)
	} else {
		version = ui.MutedStyle.Render(versionStr)
	}

	// Show service status for docker/ollama
	extra := ""
	if item.Status != "" {
		if item.Status == "running" {
			extra = "  " + serviceRunning + " " + ui.SuccessStyle.Render("running")
		} else if item.Status == "stopped" {
			extra = "  " + serviceStopped + " " + ui.WarningStyle.Render("stopped")
		} else {
			// Other status like "199 formulae, 6 casks" or "2 versions"
			extra = "  " + ui.MutedStyle.Render(item.Status)
		}
	}

	return fmt.Sprintf("  %s  %s  %s  %s%s", name, method, badge, version, extra)
}

func (m statusModel) renderResourceRowLoading(item StatusItem, colWidths ColumnWidths) string {
	name := ui.MutedStyle.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))
	return fmt.Sprintf("  %s  %s", name, m.spinner.View())
}

func (m statusModel) renderResourceRow(item StatusItem, colWidths ColumnWidths) string {
	name := ui.MutedStyle.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))
	value := ui.RenderStyledValue(item.Value, item.Style)
	return fmt.Sprintf("  %s  %s", name, value)
}

func (m statusModel) statusBadge(ok bool) string {
	if ok {
		return badgeOk.Render("✓")
	}
	return badgeErr.Render("✗")
}

// visibleLen returns the visible length of a string, stripping ANSI escape codes
func visibleLen(s string) int {
	// Strip ANSI escape sequences
	inEscape := false
	length := 0
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		length++
	}
	return length
}

// =============================================================================
// Entry Point
// =============================================================================

func showStatus() {
	p := tea.NewProgram(newStatusModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// =============================================================================
// Helpers
// =============================================================================

func shouldShowCleanHint() bool {
	for _, check := range config.CacheChecks {
		result := check.Check()
		if result.Available {
			return true
		}
	}
	return false
}
