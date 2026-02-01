package commands

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/skill"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
)

// skillItemData holds domain-specific data for each item
type skillItemData struct {
	repo       string
	skill      string
	actionType string
	isMySkill  bool
}

type skillsModel struct {
	list        *ui.List
	page        *ui.Page
	itemData    []skillItemData     // Domain data parallel to list.Items
	expanded    map[string]bool     // tracks which repos are expanded
	installed   []string            // list of installed skill names (ordered)
	repoSkills  map[string][]string // cache of fetched skills per repo
	loadingRepo string              // repo currently being loaded
	processing  bool
	quitting    bool
}

func newSkillsModel() skillsModel {
	m := skillsModel{
		expanded:   make(map[string]bool),
		installed:  skill.ListInstalled(),
		repoSkills: make(map[string][]string),
		page:       ui.NewPage("Skills"),
	}

	items, data := m.buildItems()
	m.list = ui.NewList(items)
	m.list.CalculateLabelWidth()
	m.itemData = data

	return m
}

// isInstalled checks if a skill is in the installed list
func (m skillsModel) isInstalled(skill string) bool {
	for _, s := range m.installed {
		if s == skill {
			return true
		}
	}
	return false
}

// getSkillsForRepo returns cached skills or nil if not loaded yet
func (m skillsModel) getSkillsForRepo(repoName string) []string {
	if s, ok := m.repoSkills[repoName]; ok {
		return s
	}
	return nil
}

// findRepoForSkill finds which repo a skill belongs to
func (m skillsModel) findRepoForSkill(skill string) string {
	for repoName, repoSkills := range m.repoSkills {
		for _, s := range repoSkills {
			if s == skill {
				return repoName
			}
		}
	}
	for _, s := range config.FavoriteSkills {
		if s.Skill == skill {
			return s.Repo
		}
	}
	return ""
}

func (m skillsModel) buildItems() ([]ui.Item, []skillItemData) {
	var items []ui.Item
	var data []skillItemData

	// Actions section
	items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Actions"})
	data = append(data, skillItemData{actionType: "install-my-skills"})

	items = append(items, ui.Item{Kind: ui.KindAction, Label: "Install favorites"})
	data = append(data, skillItemData{actionType: "install-my-skills"})

	items = append(items, ui.Item{Kind: ui.KindAction, Label: "Remove all skills"})
	data = append(data, skillItemData{actionType: "remove-all"})

	// Favorites section
	favorites := config.GetFavoriteSkills()
	if len(favorites) > 0 {
		items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Favorites"})
		data = append(data, skillItemData{})

		for _, s := range favorites {
			state := ui.StateUnchecked
			if m.isInstalled(s.Skill) {
				state = ui.StateChecked
			}
			items = append(items, ui.Item{
				Kind:  ui.KindToggle,
				Label: s.Skill,
				State: state,
			})
			data = append(data, skillItemData{repo: s.Repo, skill: s.Skill, isMySkill: true})
		}
	}

	// Installed section (skills not in Favorites)
	var otherInstalled []string
	for _, skill := range m.installed {
		if !config.IsFavoriteSkill("", skill) {
			otherInstalled = append(otherInstalled, skill)
		}
	}
	if len(otherInstalled) > 0 {
		items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Installed"})
		data = append(data, skillItemData{})

		for _, skill := range otherInstalled {
			repo := m.findRepoForSkill(skill)
			items = append(items, ui.Item{
				Kind:        ui.KindToggle,
				Label:       skill,
				Description: repo,
				State:       ui.StateChecked,
			})
			data = append(data, skillItemData{repo: repo, skill: skill})
		}
	}

	// Browse section
	items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Browse"})
	data = append(data, skillItemData{})

	for _, repo := range config.GetAllSkillRepos() {
		expanded := m.expanded[repo.Name]
		repoSkills := m.getSkillsForRepo(repo.Name)
		isLoading := m.loadingRepo == repo.Name

		// Build description
		desc := repo.Description
		if isLoading {
			desc = "Loading..."
		} else if repoSkills != nil && len(repoSkills) > 0 {
			installedCount := 0
			for _, skill := range repoSkills {
				if m.isInstalled(skill) {
					installedCount++
				}
			}
			desc = fmt.Sprintf("%d skills", len(repoSkills))
			if installedCount > 0 {
				desc = fmt.Sprintf("%d/%d installed", installedCount, len(repoSkills))
			}
		}

		items = append(items, ui.Item{
			Kind:        ui.KindExpandable,
			Label:       repo.Name,
			Description: desc,
			Expanded:    expanded,
		})
		data = append(data, skillItemData{repo: repo.Name})

		// If expanded, show repo actions and skills
		if expanded && repoSkills != nil {
			items = append(items, ui.Item{
				Kind:   ui.KindAction,
				Label:  "Install all",
				Indent: 1,
			})
			data = append(data, skillItemData{repo: repo.Name, actionType: "install-repo"})

			for _, skill := range repoSkills {
				state := ui.StateUnchecked
				if m.isInstalled(skill) {
					state = ui.StateChecked
				}
				items = append(items, ui.Item{
					Kind:   ui.KindToggle,
					Label:  skill,
					State:  state,
					Indent: 1,
				})
				data = append(data, skillItemData{repo: repo.Name, skill: skill})
			}
		}
	}

	return items, data
}

func (m skillsModel) rebuildItems() skillsModel {
	items, data := m.buildItems()
	cursor := m.list.Cursor
	m.list = ui.NewList(items)
	m.list.CalculateLabelWidth()
	m.itemData = data

	// Restore cursor position
	if cursor >= len(items) {
		cursor = len(items) - 1
	}
	m.list.SetCursor(cursor)

	// Skip headers if cursor landed on one
	for m.list.Cursor > 0 && !m.list.Items[m.list.Cursor].Selectable() {
		m.list.Cursor--
	}
	return m
}

func (m skillsModel) Init() tea.Cmd {
	return tea.WindowSize()
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
			m.list.Up()

		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			m.list.Down()

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			return m.handleSelect()
		}

	case skillActionDoneMsg:
		m.processing = false
		m.page.Message = msg.message
		m.page.Processing = false
		m.installed = skill.ListInstalled()
		m = m.rebuildItems()
		return m, nil

	case repoSkillsFetchedMsg:
		m.loadingRepo = ""
		if msg.err == nil {
			m.repoSkills[msg.repo] = msg.skills
		} else {
			m.page.Message = fmt.Sprintf("Failed to fetch skills for %s", msg.repo)
			m.repoSkills[msg.repo] = []string{}
		}
		m = m.rebuildItems()
		return m, nil

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		m.page.SetSize(msg.Width, msg.Height)
	}

	return m, nil
}

func (m skillsModel) handleSelect() (skillsModel, tea.Cmd) {
	idx := m.list.SelectedIndex()
	if idx < 0 || idx >= len(m.itemData) {
		return m, nil
	}

	item := m.list.Selected()
	data := m.itemData[idx]

	switch item.Kind {
	case ui.KindHeader:
		return m, nil

	case ui.KindExpandable:
		// Toggle expand/collapse
		if m.expanded[data.repo] {
			m.expanded[data.repo] = false
			m = m.rebuildItems()
		} else {
			m.expanded[data.repo] = true
			if m.getSkillsForRepo(data.repo) == nil {
				m.loadingRepo = data.repo
				m = m.rebuildItems()
				return m, m.fetchRepoSkills(data.repo)
			}
			m = m.rebuildItems()
		}
		return m, nil

	case ui.KindAction:
		m.processing = true
		m.page.Processing = true
		if data.actionType == "install-repo" {
			m.page.Message = fmt.Sprintf("Installing all from %s...", data.repo)
			return m, m.installRepo(data.repo)
		}
		m.page.Message = "Processing..."
		return m, m.runGlobalAction(data.actionType)

	case ui.KindToggle:
		m.processing = true
		m.page.Processing = true
		if item.State == ui.StateChecked {
			m.page.Message = fmt.Sprintf("Removing %s...", data.skill)
			return m, m.removeSkill(data.skill)
		} else {
			m.page.Message = fmt.Sprintf("Installing %s...", data.skill)
			return m, m.installSkill(data.repo, data.skill)
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
		repoSkills, err := skill.ListFromRepo(repo)
		if err != nil {
			return repoSkillsFetchedMsg{repo: repo, skills: []string{}, err: err}
		}
		return repoSkillsFetchedMsg{repo: repo, skills: repoSkills, err: nil}
	}
}

func (m skillsModel) installSkill(repo, name string) tea.Cmd {
	return func() tea.Msg {
		if err := skill.Install(repo, name); err != nil {
			return skillActionDoneMsg{message: fmt.Sprintf("Error: %s", err), err: err}
		}
		return skillActionDoneMsg{message: fmt.Sprintf("Installed %s", name)}
	}
}

func (m skillsModel) installRepo(repo string) tea.Cmd {
	return func() tea.Msg {
		if err := skill.InstallAll(repo); err != nil {
			return skillActionDoneMsg{message: fmt.Sprintf("Error: %s", err), err: err}
		}
		return skillActionDoneMsg{message: fmt.Sprintf("Installed all from %s", repo)}
	}
}

func (m skillsModel) removeSkill(name string) tea.Cmd {
	return func() tea.Msg {
		if err := skill.Remove(name); err != nil {
			return skillActionDoneMsg{message: fmt.Sprintf("Error: %v", err), err: err}
		}
		return skillActionDoneMsg{message: fmt.Sprintf("Removed %s", name)}
	}
}

func (m skillsModel) runGlobalAction(actionType string) tea.Cmd {
	return func() tea.Msg {
		switch actionType {
		case "install-my-skills":
			favorites := config.GetFavoriteSkills()
			installed := 0
			for _, s := range favorites {
				if err := skill.Install(s.Repo, s.Skill); err == nil {
					installed++
				}
			}
			if installed < len(favorites) {
				return skillActionDoneMsg{message: fmt.Sprintf("Installed %d/%d favorites", installed, len(favorites))}
			}
			return skillActionDoneMsg{message: fmt.Sprintf("Installed %d favorites", installed)}
		case "remove-all":
			if err := skill.RemoveAll(); err != nil {
				return skillActionDoneMsg{message: "Failed to remove skills", err: err}
			}
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

	m.page.Breadcrumbs = breadcrumbs
	if len(breadcrumbs) > 0 {
		m.page.Help = ui.DefaultHelpWithBack()
	} else {
		m.page.Help = ui.DefaultHelp()
	}

	m.page.Content = m.list.Render(m.page.ContentHeight())

	return m.page.Render()
}

func runSkillsUI() {
	if !skill.IsInstalled() {
		fmt.Printf("%s skills CLI not installed. Run: npm install -g skills\n", ui.Red("❌"))
		return
	}

	m := newSkillsModel()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running skills UI: %v\n", err)
		os.Exit(1)
	}
}
