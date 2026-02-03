package setup

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/domain/skill"
	"github.com/jterrazz/jterrazz-cli/internal/presentation/components/tui"
)

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

// InitSkillsState initializes the skills state
func InitSkillsState() {
	skills = skillsState{
		expanded:   make(map[string]bool),
		installed:  skill.ListInstalled(),
		repoSkills: make(map[string][]string),
		itemData:   nil,
	}
}

// BuildSkillsItems builds the skills menu items
func BuildSkillsItems() []tui.Item {
	var items []tui.Item
	skills.itemData = nil

	// Favorites section
	favorites := config.GetFavoriteSkills()
	if len(favorites) > 0 {
		items = append(items, tui.Item{Kind: tui.KindHeader, Label: "Pinned"})
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
	items = append(items, tui.Item{Kind: tui.KindHeader, Label: "Repositories"})
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

// HandleSkillsSelect handles item selection in the skills menu
func HandleSkillsSelect(index int, item tui.Item) tea.Cmd {
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

// SkillsFetchedMsg is sent when skills are fetched from a repo
type SkillsFetchedMsg struct {
	Repo   string
	Skills []string
	Err    error
}

// HandleSkillsMessage handles messages for the skills view
func HandleSkillsMessage(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case SkillsFetchedMsg:
		skills.loadingRepo = ""
		if msg.Err == nil {
			skills.repoSkills[msg.Repo] = msg.Skills
		} else {
			skills.repoSkills[msg.Repo] = []string{}
		}
		return nil

	case tui.ActionDoneMsg:
		skills.installed = skill.ListInstalled()
		return nil
	}
	return nil
}

func fetchSkillsCmd(repo string) tea.Cmd {
	return func() tea.Msg {
		repoSkills, err := skill.ListFromRepo(repo)
		return SkillsFetchedMsg{Repo: repo, Skills: repoSkills, Err: err}
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

// SkillsConfig returns the TUI config for the skills view
func SkillsConfig() tui.AppConfig {
	return tui.AppConfig{
		Title:      "Skills",
		BuildItems: BuildSkillsItems,
		OnSelect:   HandleSkillsSelect,
		OnMessage:  HandleSkillsMessage,
	}
}
