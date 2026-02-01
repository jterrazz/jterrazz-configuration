package commands

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/skill"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
)

// =============================================================================
// Skills State
// =============================================================================

// skillAction represents action items in skills UI
type skillAction string

const (
	skillActionInstallRepo skillAction = "install-repo"
)

type skillsState struct {
	expanded    map[string]bool
	installed   []string
	repoSkills  map[string][]string
	loadingRepo string
	itemData    []skillItemData
}

type skillItemData struct {
	repo   string
	skill  string
	action skillAction
}

var skills skillsState

func initSkillsState() {
	skills = skillsState{
		expanded:   make(map[string]bool),
		installed:  skill.ListInstalled(),
		repoSkills: make(map[string][]string),
		itemData:   nil,
	}
}

// =============================================================================
// Skills Items Builder
// =============================================================================

func buildSkillsItems() []ui.Item {
	var items []ui.Item
	skills.itemData = nil

	// Favorites section
	favorites := config.GetFavoriteSkills()
	if len(favorites) > 0 {
		items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Favorites"})
		skills.itemData = append(skills.itemData, skillItemData{})

		for _, s := range favorites {
			state := ui.StateUnchecked
			if isSkillInstalled(s.Skill) {
				state = ui.StateChecked
			}
			items = append(items, ui.Item{
				Kind:  ui.KindToggle,
				Label: s.Skill,
				State: state,
			})
			skills.itemData = append(skills.itemData, skillItemData{repo: s.Repo, skill: s.Skill})
		}
	}

	// Installed section (skills not in Favorites)
	var otherInstalled []string
	for _, s := range skills.installed {
		if !config.IsFavoriteSkill("", s) {
			otherInstalled = append(otherInstalled, s)
		}
	}
	if len(otherInstalled) > 0 {
		items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Installed"})
		skills.itemData = append(skills.itemData, skillItemData{})

		for _, s := range otherInstalled {
			repo := findRepoForSkill(s)
			items = append(items, ui.Item{
				Kind:        ui.KindToggle,
				Label:       s,
				Description: repo,
				State:       ui.StateChecked,
			})
			skills.itemData = append(skills.itemData, skillItemData{repo: repo, skill: s})
		}
	}

	// Browse section
	items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Browse"})
	skills.itemData = append(skills.itemData, skillItemData{})

	for _, repo := range config.GetAllSkillRepos() {
		expanded := skills.expanded[repo.Name]
		repoSkills := skills.repoSkills[repo.Name]
		isLoading := skills.loadingRepo == repo.Name

		desc := repo.Description
		if isLoading {
			desc = "Loading..."
		} else if len(repoSkills) > 0 {
			installedCount := 0
			for _, s := range repoSkills {
				if isSkillInstalled(s) {
					installedCount++
				}
			}
			desc = strconv.Itoa(len(repoSkills)) + " skills"
			if installedCount > 0 {
				desc = strconv.Itoa(installedCount) + "/" + strconv.Itoa(len(repoSkills)) + " installed"
			}
		}

		items = append(items, ui.Item{
			Kind:        ui.KindExpandable,
			Label:       repo.Name,
			Description: desc,
			Expanded:    expanded,
		})
		skills.itemData = append(skills.itemData, skillItemData{repo: repo.Name})

		if expanded && repoSkills != nil {
			items = append(items, ui.Item{
				Kind:   ui.KindAction,
				Label:  "Install all",
				Indent: 1,
			})
			skills.itemData = append(skills.itemData, skillItemData{repo: repo.Name, action: skillActionInstallRepo})

			for _, s := range repoSkills {
				state := ui.StateUnchecked
				if isSkillInstalled(s) {
					state = ui.StateChecked
				}
				items = append(items, ui.Item{
					Kind:   ui.KindToggle,
					Label:  s,
					State:  state,
					Indent: 1,
				})
				skills.itemData = append(skills.itemData, skillItemData{repo: repo.Name, skill: s})
			}
		}
	}

	return items
}

// =============================================================================
// Skills Helpers
// =============================================================================

func isSkillInstalled(name string) bool {
	for _, s := range skills.installed {
		if s == name {
			return true
		}
	}
	return false
}

func findRepoForSkill(name string) string {
	for repoName, repoSkills := range skills.repoSkills {
		for _, s := range repoSkills {
			if s == name {
				return repoName
			}
		}
	}
	for _, s := range config.FavoriteSkills {
		if s.Skill == name {
			return s.Repo
		}
	}
	return ""
}

// =============================================================================
// Skills Event Handlers
// =============================================================================

func handleSkillsSelect(index int, item ui.Item) tea.Cmd {
	if index >= len(skills.itemData) {
		return nil
	}
	data := skills.itemData[index]

	switch item.Kind {
	case ui.KindExpandable:
		if skills.expanded[data.repo] {
			skills.expanded[data.repo] = false
			return func() tea.Msg { return ui.RefreshMsg{} }
		} else {
			skills.expanded[data.repo] = true
			if skills.repoSkills[data.repo] == nil {
				skills.loadingRepo = data.repo
				return fetchSkillsCmd(data.repo)
			}
			return func() tea.Msg { return ui.RefreshMsg{} }
		}

	case ui.KindAction:
		if data.action == skillActionInstallRepo {
			return installRepoCmd(data.repo)
		}

	case ui.KindToggle:
		if item.State == ui.StateChecked {
			return removeSkillCmd(data.skill)
		}
		return installSkillCmd(data.repo, data.skill)
	}

	return nil
}

// =============================================================================
// Skills Messages
// =============================================================================

type skillsFetchedMsg struct {
	repo   string
	skills []string
	err    error
}

func handleSkillsMessage(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case skillsFetchedMsg:
		skills.loadingRepo = ""
		if msg.err == nil {
			skills.repoSkills[msg.repo] = msg.skills
		} else {
			skills.repoSkills[msg.repo] = []string{}
		}
		return nil

	case ui.ActionDoneMsg:
		skills.installed = skill.ListInstalled()
		return nil
	}
	return nil
}

// =============================================================================
// Skills Commands
// =============================================================================

func fetchSkillsCmd(repo string) tea.Cmd {
	return func() tea.Msg {
		repoSkills, err := skill.ListFromRepo(repo)
		return skillsFetchedMsg{repo: repo, skills: repoSkills, err: err}
	}
}

func installSkillCmd(repo, name string) tea.Cmd {
	return func() tea.Msg {
		if err := skill.Install(repo, name); err != nil {
			return ui.ActionDoneMsg{Message: "Error: " + err.Error(), Err: err}
		}
		return ui.ActionDoneMsg{Message: "Installed " + name}
	}
}

func installRepoCmd(repo string) tea.Cmd {
	return func() tea.Msg {
		if err := skill.InstallAll(repo); err != nil {
			return ui.ActionDoneMsg{Message: "Error: " + err.Error(), Err: err}
		}
		return ui.ActionDoneMsg{Message: "Installed all from " + repo}
	}
}

func removeSkillCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := skill.Remove(name); err != nil {
			return ui.ActionDoneMsg{Message: "Error: " + err.Error(), Err: err}
		}
		return ui.ActionDoneMsg{Message: "Removed " + name}
	}
}

// =============================================================================
// Skills Config
// =============================================================================

func skillsConfig() ui.AppConfig {
	initSkillsState()
	return ui.AppConfig{
		Title:      "Skills",
		BuildItems: buildSkillsItems,
		OnSelect:   handleSkillsSelect,
		OnMessage:  handleSkillsMessage,
	}
}

func runSkillsUI() {
	if !skill.IsInstalled() {
		ui.PrintError("skills CLI not installed. Run: npm install -g skills")
		return
	}

	ui.RunOrExit(skillsConfig())
}
