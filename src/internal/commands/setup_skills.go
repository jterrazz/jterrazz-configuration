package commands

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/skill"
	"github.com/jterrazz/jterrazz-cli/internal/ui/components/tui"
	"github.com/jterrazz/jterrazz-cli/internal/ui/print"
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

func buildSkillsItems() []tui.Item {
	var items []tui.Item
	skills.itemData = nil

	// Favorites section
	favorites := config.GetFavoriteSkills()
	if len(favorites) > 0 {
		items = append(items, tui.Item{Kind: tui.KindHeader, Label: "Favorites"})
		skills.itemData = append(skills.itemData, skillItemData{})

		for _, s := range favorites {
			state := tui.StateUnchecked
			if isSkillInstalled(s.Skill) {
				state = tui.StateChecked
			}
			items = append(items, tui.Item{
				Kind:  tui.KindToggle,
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
		items = append(items, tui.Item{Kind: tui.KindHeader, Label: "Installed"})
		skills.itemData = append(skills.itemData, skillItemData{})

		for _, s := range otherInstalled {
			repo := findRepoForSkill(s)
			items = append(items, tui.Item{
				Kind:        tui.KindToggle,
				Label:       s,
				Description: repo,
				State:       tui.StateChecked,
			})
			skills.itemData = append(skills.itemData, skillItemData{repo: repo, skill: s})
		}
	}

	// Browse section
	items = append(items, tui.Item{Kind: tui.KindHeader, Label: "Browse"})
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

		items = append(items, tui.Item{
			Kind:        tui.KindExpandable,
			Label:       repo.Name,
			Description: desc,
			Expanded:    expanded,
		})
		skills.itemData = append(skills.itemData, skillItemData{repo: repo.Name})

		if expanded && repoSkills != nil {
			items = append(items, tui.Item{
				Kind:   tui.KindAction,
				Label:  "Install all",
				Indent: 1,
			})
			skills.itemData = append(skills.itemData, skillItemData{repo: repo.Name, action: skillActionInstallRepo})

			for _, s := range repoSkills {
				state := tui.StateUnchecked
				if isSkillInstalled(s) {
					state = tui.StateChecked
				}
				items = append(items, tui.Item{
					Kind:   tui.KindToggle,
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

func handleSkillsSelect(index int, item tui.Item) tea.Cmd {
	if index >= len(skills.itemData) {
		return nil
	}
	data := skills.itemData[index]

	switch item.Kind {
	case tui.KindExpandable:
		if skills.expanded[data.repo] {
			skills.expanded[data.repo] = false
			return func() tea.Msg { return tui.RefreshMsg{} }
		} else {
			skills.expanded[data.repo] = true
			if skills.repoSkills[data.repo] == nil {
				skills.loadingRepo = data.repo
				return fetchSkillsCmd(data.repo)
			}
			return func() tea.Msg { return tui.RefreshMsg{} }
		}

	case tui.KindAction:
		if data.action == skillActionInstallRepo {
			return installRepoCmd(data.repo)
		}

	case tui.KindToggle:
		if item.State == tui.StateChecked {
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

	case tui.ActionDoneMsg:
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
			return tui.ActionDoneMsg{Message: "Error: " + err.Error(), Err: err}
		}
		return tui.ActionDoneMsg{Message: "Installed " + name}
	}
}

func installRepoCmd(repo string) tea.Cmd {
	return func() tea.Msg {
		if err := skill.InstallAll(repo); err != nil {
			return tui.ActionDoneMsg{Message: "Error: " + err.Error(), Err: err}
		}
		return tui.ActionDoneMsg{Message: "Installed all from " + repo}
	}
}

func removeSkillCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := skill.Remove(name); err != nil {
			return tui.ActionDoneMsg{Message: "Error: " + err.Error(), Err: err}
		}
		return tui.ActionDoneMsg{Message: "Removed " + name}
	}
}

// =============================================================================
// Skills Config
// =============================================================================

func skillsConfig() tui.AppConfig {
	initSkillsState()
	return tui.AppConfig{
		Title:      "Skills",
		BuildItems: buildSkillsItems,
		OnSelect:   handleSkillsSelect,
		OnMessage:  handleSkillsMessage,
	}
}

func runSkillsUI() {
	if !skill.IsInstalled() {
		print.Error("skills CLI not installed. Run: npm install -g skills")
		return
	}

	tui.RunOrExit(skillsConfig())
}
