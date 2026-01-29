package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
)

// Item types
const (
	itemTypeAction     = "action"
	itemTypeHeader     = "header"
	itemTypeSkill      = "skill"
	itemTypeRepo       = "repo"
	itemTypeRepoAction = "repo-action"
)

// skillItem represents a selectable item in the TUI
type skillItem struct {
	itemType    string
	repo        string
	skill       string
	description string
	installed   bool
	expanded    bool
	actionType  string
	isMySkill   bool // true if this skill is from "My Skills" section
	isNested    bool // true if this skill is under an expanded repo
}

type skillsModel struct {
	items       []skillItem
	cursor      int
	expanded    map[string]bool     // tracks which repos are expanded
	installed   []string            // list of installed skill names (ordered)
	repoSkills  map[string][]string // cache of fetched skills per repo
	loadingRepo string              // repo currently being loaded
	width       int
	height      int
	message     string
	processing  bool
	quitting    bool
	maxSkillLen int // for aligning descriptions
}

func initialSkillsModel() skillsModel {
	installed := getInstalledSkills()

	// Calculate max skill name length for alignment
	maxLen := 0
	for _, repo := range SkillRepos {
		if len(repo.Name) > maxLen {
			maxLen = len(repo.Name)
		}
	}

	return skillsModel{
		expanded:    make(map[string]bool),
		installed:   installed,
		repoSkills:  make(map[string][]string),
		width:       80,
		height:      24,
		maxSkillLen: maxLen,
	}
}

// getSkillsForRepo returns cached skills or empty slice if not loaded yet
func (m *skillsModel) getSkillsForRepo(repoName string) []string {
	if skills, ok := m.repoSkills[repoName]; ok {
		return skills
	}
	return nil
}

// isInstalled checks if a skill is in the installed list
func (m *skillsModel) isInstalled(skill string) bool {
	for _, s := range m.installed {
		if s == skill {
			return true
		}
	}
	return false
}

// findRepoForSkill finds which repo a skill belongs to
func (m *skillsModel) findRepoForSkill(skill string) string {
	// Check in cached repo skills
	for repoName, skills := range m.repoSkills {
		for _, s := range skills {
			if s == skill {
				return repoName
			}
		}
	}
	// Check in MySkills as fallback
	for _, s := range MySkills {
		if s.Skill == skill {
			return s.Repo
		}
	}
	return ""
}

func (m *skillsModel) buildItems() []skillItem {
	var items []skillItem

	// Actions section
	items = append(items, skillItem{itemType: itemTypeHeader, description: "Actions"})
	items = append(items, skillItem{itemType: itemTypeAction, description: "Install favorites", actionType: "install-my-skills"})
	items = append(items, skillItem{itemType: itemTypeAction, description: "Remove all skills", actionType: "remove-all"})

	// Favorites section
	if len(MySkills) > 0 {
		items = append(items, skillItem{itemType: itemTypeHeader, description: "Favorites"})
		for _, s := range MySkills {
			items = append(items, skillItem{
				itemType:  itemTypeSkill,
				repo:      s.Repo,
				skill:     s.Skill,
				installed: m.isInstalled(s.Skill),
				isMySkill: true,
			})
		}
	}

	// Installed section (skills not in Favorites)
	var otherInstalled []string
	for _, skill := range m.installed {
		isFavorite := false
		for _, s := range MySkills {
			if s.Skill == skill {
				isFavorite = true
				break
			}
		}
		if !isFavorite {
			otherInstalled = append(otherInstalled, skill)
		}
	}
	if len(otherInstalled) > 0 {
		items = append(items, skillItem{itemType: itemTypeHeader, description: "Installed"})
		for _, skill := range otherInstalled {
			repo := m.findRepoForSkill(skill)
			items = append(items, skillItem{
				itemType:  itemTypeSkill,
				repo:      repo,
				skill:     skill,
				installed: true,
			})
		}
	}

	// Browse section
	items = append(items, skillItem{itemType: itemTypeHeader, description: "Browse"})
	for _, repo := range SkillRepos {
		expanded := m.expanded[repo.Name]
		repoSkills := m.getSkillsForRepo(repo.Name)
		isLoading := m.loadingRepo == repo.Name

		// Count installed vs total
		installedCount := 0
		for _, skill := range repoSkills {
			if m.isInstalled(skill) {
				installedCount++
			}
		}

		desc := repo.Description
		if isLoading {
			desc = "Loading..."
		} else if repoSkills != nil && len(repoSkills) > 0 {
			desc = fmt.Sprintf("%d skills", len(repoSkills))
			if installedCount > 0 {
				desc = fmt.Sprintf("%d/%d installed", installedCount, len(repoSkills))
			}
		}

		items = append(items, skillItem{
			itemType:    itemTypeRepo,
			repo:        repo.Name,
			description: desc,
			expanded:    expanded,
		})

		// If expanded, show repo actions and skills
		if expanded && repoSkills != nil {
			items = append(items, skillItem{
				itemType:    itemTypeRepoAction,
				repo:        repo.Name,
				description: "Install all",
				actionType:  "install-repo",
			})

			for _, skill := range repoSkills {
				isInstalled := m.isInstalled(skill)
				items = append(items, skillItem{
					itemType:  itemTypeSkill,
					repo:      repo.Name,
					skill:     skill,
					installed: isInstalled,
					isNested:  true,
				})
			}
		}
	}

	return items
}

func (m skillsModel) Init() tea.Cmd {
	return nil
}

func (m skillsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.moveCursor(-1)

		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			m.moveCursor(1)

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			newModel, cmd := m.handleSelect()
			return newModel, cmd
		}

	case skillActionDoneMsg:
		m.processing = false
		m.message = msg.message
		m.installed = getInstalledSkills()
		// Adjust cursor if it's now out of bounds or on a header
		items := m.buildItems()
		if m.cursor >= len(items) {
			m.cursor = len(items) - 1
		}
		for m.cursor > 0 && items[m.cursor].itemType == itemTypeHeader {
			m.cursor--
		}
		return m, nil

	case repoSkillsFetchedMsg:
		m.loadingRepo = ""
		if msg.err == nil {
			m.repoSkills[msg.repo] = msg.skills
			// Update maxSkillLen if needed
			for _, skill := range msg.skills {
				if len(skill) > m.maxSkillLen {
					m.maxSkillLen = len(skill)
				}
			}
		} else {
			m.message = fmt.Sprintf("Failed to fetch skills for %s", msg.repo)
			m.repoSkills[msg.repo] = []string{} // Cache empty to avoid retrying
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m *skillsModel) moveCursor(delta int) {
	items := m.buildItems()
	newCursor := m.cursor + delta

	// Skip headers
	for newCursor >= 0 && newCursor < len(items) && items[newCursor].itemType == itemTypeHeader {
		newCursor += delta
	}

	if newCursor >= 0 && newCursor < len(items) {
		m.cursor = newCursor
	}
}

func (m skillsModel) handleSelect() (skillsModel, tea.Cmd) {
	items := m.buildItems()
	if m.cursor < 0 || m.cursor >= len(items) {
		return m, nil
	}

	item := items[m.cursor]

	switch item.itemType {
	case itemTypeHeader:
		return m, nil

	case itemTypeRepo:
		// Toggle expand/collapse
		if m.expanded[item.repo] {
			// Collapse
			m.expanded[item.repo] = false
		} else {
			// Expand - fetch skills if not cached
			m.expanded[item.repo] = true
			if m.getSkillsForRepo(item.repo) == nil {
				m.loadingRepo = item.repo
				return m, m.fetchRepoSkills(item.repo)
			}
		}
		return m, nil

	case itemTypeAction:
		m.processing = true
		m.message = "Processing..."
		return m, m.runGlobalAction(item.actionType)

	case itemTypeRepoAction:
		m.processing = true
		m.message = fmt.Sprintf("Installing all from %s...", item.repo)
		return m, m.installRepo(item.repo)

	case itemTypeSkill:
		m.processing = true
		if item.installed {
			m.message = fmt.Sprintf("Removing %s...", item.skill)
			return m, m.removeSkill(item.skill)
		} else {
			m.message = fmt.Sprintf("Installing %s...", item.skill)
			return m, m.installSkill(item.repo, item.skill)
		}
	}

	return m, nil
}

type skillActionDoneMsg struct {
	message string
	err     error
}

type repoSkillsFetchedMsg struct {
	repo   string
	skills []string
	err    error
}

func (m skillsModel) fetchRepoSkills(repo string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("skills", "add", repo, "--list")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return repoSkillsFetchedMsg{repo: repo, skills: []string{}, err: err}
		}

		// Parse the output to extract skill names
		skills := parseSkillsListOutput(string(output))
		return repoSkillsFetchedMsg{repo: repo, skills: skills, err: nil}
	}
}

func parseSkillsListOutput(output string) []string {
	var skills []string
	cleanOutput := stripAnsi(output)
	lines := strings.Split(cleanOutput, "\n")

	// Skills appear after "Available Skills" line
	// Format (with box-drawing characters):
	// │    skill-name
	// │
	// │      Description text...
	inSkillsSection := false
	for _, line := range lines {
		if strings.Contains(line, "Available Skills") {
			inSkillsSection = true
			continue
		}

		if !inSkillsSection {
			continue
		}

		// Stop at "Use --skill" line
		if strings.Contains(line, "Use --skill") {
			break
		}

		// Remove box-drawing characters and trim
		// The line format is: "│    skill-name" or "│      description"
		cleaned := line
		cleaned = strings.ReplaceAll(cleaned, "│", "")
		cleaned = strings.ReplaceAll(cleaned, "├", "")
		cleaned = strings.ReplaceAll(cleaned, "└", "")
		cleaned = strings.ReplaceAll(cleaned, "┌", "")
		cleaned = strings.ReplaceAll(cleaned, "◇", "")

		// Count leading spaces before trimming
		leadingSpaces := len(cleaned) - len(strings.TrimLeft(cleaned, " "))
		trimmed := strings.TrimSpace(cleaned)

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		// Skill names have ~4 leading spaces, descriptions have more (~6)
		// Skill names are single words (no spaces)
		if leadingSpaces <= 5 && !strings.Contains(trimmed, " ") && len(trimmed) > 0 {
			if isValidSkillName(trimmed) {
				skills = append(skills, trimmed)
			}
		}
	}

	return skills
}

func isValidSkillName(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}
	return true
}

func (m skillsModel) installSkill(repo, skill string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("skills", "add", repo, "-g", "-y", "--skill", skill)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return skillActionDoneMsg{message: fmt.Sprintf("Error: %s", strings.TrimSpace(string(output))), err: err}
		}
		return skillActionDoneMsg{message: fmt.Sprintf("Installed %s", skill)}
	}
}

func (m skillsModel) installRepo(repo string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("skills", "add", repo, "-g", "-y", "--all")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return skillActionDoneMsg{message: fmt.Sprintf("Error: %s", strings.TrimSpace(string(output))), err: err}
		}
		return skillActionDoneMsg{message: fmt.Sprintf("Installed all from %s", repo)}
	}
}

func (m skillsModel) removeSkill(skill string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("skills", "remove", "-g", "-y", skill)
		if err := cmd.Run(); err != nil {
			return skillActionDoneMsg{message: fmt.Sprintf("Error removing %s: %v", skill, err), err: err}
		}
		return skillActionDoneMsg{message: fmt.Sprintf("Removed %s", skill)}
	}
}

func (m skillsModel) runGlobalAction(actionType string) tea.Cmd {
	return func() tea.Msg {
		switch actionType {
		case "install-my-skills":
			for _, s := range MySkills {
				cmd := exec.Command("skills", "add", s.Repo, "-g", "-y", "--skill", s.Skill)
				cmd.Run()
			}
			return skillActionDoneMsg{message: fmt.Sprintf("Installed %d favorites", len(MySkills))}
		case "remove-all":
			cmd := exec.Command("skills", "remove", "-g", "-y", "--all")
			cmd.Run()
			return skillActionDoneMsg{message: "Removed all skills"}
		}
		return skillActionDoneMsg{message: "Unknown action"}
	}
}

func (m skillsModel) View() string {
	return m.viewWithBreadcrumb()
}

func (m skillsModel) viewWithBreadcrumb(breadcrumbs ...string) string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Title with optional breadcrumb
	if len(breadcrumbs) > 0 {
		b.WriteString(ui.RenderBreadcrumb(breadcrumbs...) + "\n\n")
	} else {
		b.WriteString(ui.TitleStyle.Render("Skills") + "\n\n")
	}

	items := m.buildItems()

	// Calculate visible range
	visibleHeight := m.height - 6 // Account for title, help, message
	startIdx := 0
	if m.cursor > visibleHeight-3 {
		startIdx = m.cursor - visibleHeight + 3
	}
	endIdx := startIdx + visibleHeight
	if endIdx > len(items) {
		endIdx = len(items)
	}

	for i := startIdx; i < endIdx; i++ {
		item := items[i]
		selected := i == m.cursor

		line := m.renderItem(item, selected)
		b.WriteString(line + "\n")
	}

	// Help
	helpText := "↑/↓ navigate • enter select/toggle • q quit"
	if len(breadcrumbs) > 0 {
		helpText = "↑/↓ navigate • enter select/toggle • esc back • q quit"
	}
	b.WriteString(ui.HelpStyle.Render(helpText))

	// Message
	if m.message != "" {
		b.WriteString("\n")
		if m.processing {
			b.WriteString(ui.ActionStyle.Render(m.message))
		} else {
			b.WriteString(ui.SuccessStyle.Render(m.message))
		}
	}

	return b.String()
}

func (m skillsModel) renderItem(item skillItem, selected bool) string {
	switch item.itemType {
	case itemTypeHeader:
		return ui.RenderSection(item.description)

	case itemTypeAction:
		prefix := "  "
		if selected {
			prefix = ui.IconSelected + " "
			return ui.SelectedStyle.Render(prefix + item.description)
		}
		return ui.ActionStyle.Render(prefix + item.description)

	case itemTypeRepo:
		var arrow string
		if item.expanded {
			arrow = ui.IconArrowDown
		} else {
			arrow = ui.IconArrowRight
		}

		prefix := "  "
		paddedRepo := fmt.Sprintf("%-*s", m.maxSkillLen, item.repo)
		if selected {
			prefix = ui.IconSelected + " "
			return ui.SelectedStyle.Render(fmt.Sprintf("%s%s %s", prefix, arrow, paddedRepo)) +
				ui.MutedStyle.Render(fmt.Sprintf("  %s", item.description))
		}
		return ui.NormalStyle.Render(fmt.Sprintf("%s%s %s", prefix, arrow, paddedRepo)) +
			ui.MutedStyle.Render(fmt.Sprintf("  %s", item.description))

	case itemTypeRepoAction:
		prefix := "      "
		if selected {
			prefix = "    " + ui.IconSelected + " "
			return ui.SelectedStyle.Render(prefix + item.description)
		}
		return ui.ActionStyle.Render(prefix + item.description)

	case itemTypeSkill:
		var status string
		var style lipgloss.Style

		if item.installed {
			status = ui.IconCheck
			style = ui.SuccessStyle
		} else {
			status = "○"
			style = ui.MutedStyle
		}

		// Nested skill under expanded repo
		if item.isNested {
			prefix := "      "
			if selected {
				prefix = "    " + ui.IconSelected + " "
				return ui.SelectedStyle.Render(fmt.Sprintf("%s%s %s", prefix, status, item.skill))
			}
			return style.Render(fmt.Sprintf("%s%s %s", prefix, status, item.skill))
		}

		// Top-level skill (in My Skills or Installed section)
		prefix := "  "
		paddedSkill := fmt.Sprintf("%-*s", m.maxSkillLen, item.skill)
		if selected {
			prefix = ui.IconSelected + " "
		}

		// My Skills section: don't show repo name
		if item.isMySkill {
			if selected {
				return ui.SelectedStyle.Render(fmt.Sprintf("%s%s %s", prefix, status, item.skill))
			}
			return style.Render(fmt.Sprintf("%s%s %s", prefix, status, item.skill))
		}

		// Installed section: show repo name aligned
		if selected {
			return ui.SelectedStyle.Render(fmt.Sprintf("%s%s %s", prefix, status, paddedSkill)) +
				ui.MutedStyle.Render(fmt.Sprintf("  %s", item.repo))
		}

		repoInfo := ""
		if item.repo != "" {
			repoInfo = ui.MutedStyle.Render(fmt.Sprintf("  %s", item.repo))
		}

		return style.Render(fmt.Sprintf("%s%s %s", prefix, status, paddedSkill)) + repoInfo
	}

	return ""
}

func getInstalledSkills() []string {
	var installed []string

	cmd := exec.Command("skills", "list", "-g")
	output, err := cmd.Output()
	if err != nil {
		return installed
	}

	// Strip ANSI escape codes and parse output format:
	// Global Skills
	//
	// skill-name ~/.agents/skills/skill-name
	//   Agents: ...
	cleanOutput := stripAnsi(string(output))
	lines := strings.Split(cleanOutput, "\n")
	for _, line := range lines {
		// Skip lines that start with whitespace (agent info lines)
		if len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
			continue
		}

		line = strings.TrimSpace(line)

		// Skip empty lines, headers, and info messages
		if line == "" ||
			strings.Contains(line, "No global skills") ||
			strings.Contains(line, "Global") ||
			strings.Contains(line, "Skills") ||
			strings.HasPrefix(line, "Try ") {
			continue
		}

		// Format: "skill-name ~/.agents/skills/skill-name"
		// First word is the skill name
		parts := strings.Fields(line)
		if len(parts) >= 1 {
			skillName := parts[0]
			// Validate it looks like a skill name (not a path, not a header)
			if !strings.HasPrefix(skillName, "/") &&
				!strings.HasPrefix(skillName, "~") &&
				!strings.Contains(skillName, ":") &&
				len(skillName) > 0 {
				installed = append(installed, skillName)
			}
		}
	}

	return installed
}

func runSkillsUI() {
	// Check if skills CLI is installed
	if _, err := exec.LookPath("skills"); err != nil {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("%s skills CLI not installed. Run: npm install -g skills\n", red("❌"))
		return
	}

	m := initialSkillsModel()
	m.items = m.buildItems()

	// Skip first header
	m.cursor = 1

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running skills UI: %v\n", err)
		os.Exit(1)
	}
}
