package commands

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/skill"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
)

// =============================================================================
// Skills State
// =============================================================================

type skillsState struct {
	expanded    map[string]bool
	installed   []string
	repoSkills  map[string][]string
	loadingRepo string
	itemData    []skillItemData
}

type skillItemData struct {
	repo       string
	skill      string
	actionType string
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

	// Actions section
	items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Actions"})
	skills.itemData = append(skills.itemData, skillItemData{actionType: "install-favorites"})

	items = append(items, ui.Item{Kind: ui.KindAction, Label: "Install favorites"})
	skills.itemData = append(skills.itemData, skillItemData{actionType: "install-favorites"})

	items = append(items, ui.Item{Kind: ui.KindAction, Label: "Remove all skills"})
	skills.itemData = append(skills.itemData, skillItemData{actionType: "remove-all"})

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
		skills.itemData = append(skills.itemData, skillItemData{repo: repo.Name})

		if expanded && repoSkills != nil {
			items = append(items, ui.Item{
				Kind:   ui.KindAction,
				Label:  "Install all",
				Indent: 1,
			})
			skills.itemData = append(skills.itemData, skillItemData{repo: repo.Name, actionType: "install-repo"})

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
		} else {
			skills.expanded[data.repo] = true
			if skills.repoSkills[data.repo] == nil {
				skills.loadingRepo = data.repo
				return fetchSkillsCmd(data.repo)
			}
		}
		return nil

	case ui.KindAction:
		switch data.actionType {
		case "install-repo":
			return installRepoCmd(data.repo)
		case "install-favorites":
			return installFavoritesCmd()
		case "remove-all":
			return removeAllSkillsCmd()
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
			return ui.ActionDoneMsg{Message: fmt.Sprintf("Error: %s", err), Err: err}
		}
		return ui.ActionDoneMsg{Message: fmt.Sprintf("Installed %s", name)}
	}
}

func installRepoCmd(repo string) tea.Cmd {
	return func() tea.Msg {
		if err := skill.InstallAll(repo); err != nil {
			return ui.ActionDoneMsg{Message: fmt.Sprintf("Error: %s", err), Err: err}
		}
		return ui.ActionDoneMsg{Message: fmt.Sprintf("Installed all from %s", repo)}
	}
}

func installFavoritesCmd() tea.Cmd {
	return func() tea.Msg {
		favorites := config.GetFavoriteSkills()
		installed := 0
		for _, s := range favorites {
			if err := skill.Install(s.Repo, s.Skill); err == nil {
				installed++
			}
		}
		if installed < len(favorites) {
			return ui.ActionDoneMsg{Message: fmt.Sprintf("Installed %d/%d favorites", installed, len(favorites))}
		}
		return ui.ActionDoneMsg{Message: fmt.Sprintf("Installed %d favorites", installed)}
	}
}

func removeSkillCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := skill.Remove(name); err != nil {
			return ui.ActionDoneMsg{Message: fmt.Sprintf("Error: %v", err), Err: err}
		}
		return ui.ActionDoneMsg{Message: fmt.Sprintf("Removed %s", name)}
	}
}

func removeAllSkillsCmd() tea.Cmd {
	return func() tea.Msg {
		if err := skill.RemoveAll(); err != nil {
			return ui.ActionDoneMsg{Message: "Failed to remove skills", Err: err}
		}
		return ui.ActionDoneMsg{Message: "Removed all skills"}
	}
}

// =============================================================================
// Skills Runner
// =============================================================================

func runSkillsUI() {
	if !skill.IsInstalled() {
		fmt.Printf("%s skills CLI not installed. Run: npm install -g skills\n", ui.Red("❌"))
		return
	}

	initSkillsState()

	ui.RunOrExit(ui.AppConfig{
		Title:      "Skills",
		BuildItems: buildSkillsItems,
		OnSelect:   handleSkillsSelect,
		OnMessage:  handleSkillsMessage,
	})
}
