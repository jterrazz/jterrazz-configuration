package commands

import (
	"testing"

	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
)

func TestFavoriteSkillsNotEmpty(t *testing.T) {
	// Given: favorite skills are defined
	skills := config.GetFavoriteSkills()

	// When: checking list
	// Then: it should not be empty
	if len(skills) == 0 {
		t.Error("FavoriteSkills should not be empty")
	}
}

func TestFavoriteSkillsHaveValidFields(t *testing.T) {
	for _, s := range config.GetFavoriteSkills() {
		t.Run(s.Skill, func(t *testing.T) {
			// Given: a favorite skill
			// When: checking fields
			// Then: repo should not be empty
			if s.Repo == "" {
				t.Errorf("FavoriteSkill %q has empty repo", s.Skill)
			}

			// Then: skill name should not be empty
			if s.Skill == "" {
				t.Errorf("FavoriteSkill in repo %q has empty skill name", s.Repo)
			}
		})
	}
}

func TestSkillReposNotEmpty(t *testing.T) {
	// Given: skill repos are defined
	repos := config.GetAllSkillRepos()

	// When: checking list
	// Then: it should not be empty
	if len(repos) == 0 {
		t.Error("SkillRepos should not be empty")
	}
}

func TestSkillReposHaveRequiredFields(t *testing.T) {
	for _, repo := range config.GetAllSkillRepos() {
		t.Run(repo.Name, func(t *testing.T) {
			// Given: a skill repo
			// When: checking fields
			// Then: name should not be empty
			if repo.Name == "" {
				t.Error("SkillRepo.Name should not be empty")
			}

			// Then: description should not be empty
			if repo.Description == "" {
				t.Errorf("SkillRepo %s should have a description", repo.Name)
			}

			// Then: name should be in owner/repo format
			if !containsSlash(repo.Name) {
				t.Errorf("SkillRepo.Name %q should be in owner/repo format", repo.Name)
			}
		})
	}
}

func containsSlash(s string) bool {
	for _, c := range s {
		if c == '/' {
			return true
		}
	}
	return false
}

func TestGetSkillRepoByName(t *testing.T) {
	tests := []struct {
		name     string
		given    string
		expected bool
	}{
		{"existing repo", "anthropics/skills", true},
		{"another existing repo", "vercel-labs/agent-skills", true},
		{"nonexistent repo", "nonexistent/repo", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: a repo name
			// When: getting repo by name
			repo := config.GetSkillRepoByName(tt.given)

			// Then: result should match expectation
			if tt.expected && repo == nil {
				t.Errorf("GetSkillRepoByName(%q) returned nil, expected repo", tt.given)
			}
			if !tt.expected && repo != nil {
				t.Errorf("GetSkillRepoByName(%q) returned repo, expected nil", tt.given)
			}
		})
	}
}

func TestSkillsModelIsInstalled(t *testing.T) {
	tests := []struct {
		name     string
		given    string
		expected bool
	}{
		{"installed skill-a", "skill-a", true},
		{"installed skill-b", "skill-b", true},
		{"installed skill-c", "skill-c", true},
		{"not installed skill-d", "skill-d", false},
		{"empty string", "", false},
	}

	// Given: skills model with installed skills
	m := &skillsModel{
		installed: []string{"skill-a", "skill-b", "skill-c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When: checking if installed
			result := m.isInstalled(tt.given)

			// Then: result should match
			if result != tt.expected {
				t.Errorf("isInstalled(%q) = %v, expected %v", tt.given, result, tt.expected)
			}
		})
	}
}

func TestSkillsModelFindRepoForSkill(t *testing.T) {
	tests := []struct {
		name     string
		given    string
		expected string
	}{
		{"skill in anthropics", "frontend-design", "anthropics/skills"},
		{"skill in vercel-labs", "vercel-react-best-practices", "vercel-labs/agent-skills"},
		{"nonexistent skill", "nonexistent-skill", ""},
	}

	// Given: skills model with repo skills cache
	m := &skillsModel{
		repoSkills: map[string][]string{
			"anthropics/skills":        {"frontend-design", "skill-creator"},
			"vercel-labs/agent-skills": {"vercel-react-best-practices"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When: finding repo for skill
			repo := m.findRepoForSkill(tt.given)

			// Then: correct repo should be returned
			if repo != tt.expected {
				t.Errorf("findRepoForSkill(%q) = %q, expected %q", tt.given, repo, tt.expected)
			}
		})
	}
}

func TestSkillsModelFindRepoFromFavorites(t *testing.T) {
	// Given: skills model with empty cache
	m := &skillsModel{
		repoSkills: map[string][]string{},
	}

	// When: finding skill from favorites
	repo := m.findRepoForSkill("frontend-design")

	// Then: it should find from FavoriteSkills
	if repo != "anthropics/skills" {
		t.Errorf("findRepoForSkill('frontend-design') = %q, expected 'anthropics/skills'", repo)
	}
}

func TestSkillsModelBuildItems(t *testing.T) {
	// Given: empty skills model
	m := &skillsModel{
		expanded:   make(map[string]bool),
		installed:  []string{},
		repoSkills: make(map[string][]string),
		page:       ui.NewPage("Skills"),
	}

	// When: building items
	items, _ := m.buildItems()

	// Then: should have minimum items
	if len(items) < 4 {
		t.Errorf("buildItems() returned %d items, expected at least 4", len(items))
	}

	// Then: first item should be Actions header
	if items[0].Kind != ui.KindHeader || items[0].Label != "Actions" {
		t.Error("First item should be Actions header")
	}

	// Then: should have Install favorites action
	hasInstallFavorites := false
	hasRemoveAll := false
	hasBrowseHeader := false
	for _, item := range items {
		if item.Kind == ui.KindAction && item.Label == "Install favorites" {
			hasInstallFavorites = true
		}
		if item.Kind == ui.KindAction && item.Label == "Remove all skills" {
			hasRemoveAll = true
		}
		if item.Kind == ui.KindHeader && item.Label == "Browse" {
			hasBrowseHeader = true
		}
	}
	if !hasInstallFavorites {
		t.Error("Missing 'Install favorites' action")
	}
	if !hasRemoveAll {
		t.Error("Missing 'Remove all skills' action")
	}
	if !hasBrowseHeader {
		t.Error("Missing 'Browse' header")
	}
}

func TestSkillsModelBuildItemsWithInstalled(t *testing.T) {
	// Given: skills model with installed favorite
	m := &skillsModel{
		expanded:   make(map[string]bool),
		installed:  []string{"frontend-design"},
		repoSkills: make(map[string][]string),
		page:       ui.NewPage("Skills"),
	}

	// When: building items
	items, data := m.buildItems()

	// Then: installed skill should be checked in Favorites
	hasInstalledSkill := false
	for i, item := range items {
		if item.Kind == ui.KindToggle && data[i].skill == "frontend-design" && item.State == ui.StateChecked {
			hasInstalledSkill = true
			break
		}
	}
	if !hasInstalledSkill {
		t.Error("Installed skill 'frontend-design' not found as checked")
	}
}

func TestSkillsModelBuildItemsWithOtherInstalled(t *testing.T) {
	// Given: skills model with installed non-favorite
	m := &skillsModel{
		expanded:   make(map[string]bool),
		installed:  []string{"some-other-skill"},
		repoSkills: make(map[string][]string),
		page:       ui.NewPage("Skills"),
	}

	// When: building items
	items, data := m.buildItems()

	// Then: Installed header should exist
	hasInstalledHeader := false
	hasInstalledSkill := false
	for i, item := range items {
		if item.Kind == ui.KindHeader && item.Label == "Installed" {
			hasInstalledHeader = true
		}
		if item.Kind == ui.KindToggle && data[i].skill == "some-other-skill" && item.State == ui.StateChecked {
			hasInstalledSkill = true
		}
	}
	if !hasInstalledHeader {
		t.Error("Missing 'Installed' header")
	}
	if !hasInstalledSkill {
		t.Error("Installed skill 'some-other-skill' not found")
	}
}

func TestSkillsModelBuildItemsWithExpanded(t *testing.T) {
	// Given: skills model with expanded repo
	m := &skillsModel{
		expanded:  map[string]bool{"anthropics/skills": true},
		installed: []string{},
		repoSkills: map[string][]string{
			"anthropics/skills": {"frontend-design", "skill-creator"},
		},
		page: ui.NewPage("Skills"),
	}

	// When: building items
	items, data := m.buildItems()

	// Then: Install All action should exist for repo
	hasRepoAction := false
	hasSkillsFromRepo := false
	for i, item := range items {
		if item.Kind == ui.KindAction && data[i].repo == "anthropics/skills" && data[i].actionType == "install-repo" {
			hasRepoAction = true
		}
		if item.Kind == ui.KindToggle && data[i].repo == "anthropics/skills" && item.Indent == 1 {
			hasSkillsFromRepo = true
		}
	}
	if !hasRepoAction {
		t.Error("Missing 'Install All' action for expanded repo")
	}
	if !hasSkillsFromRepo {
		t.Error("Missing skills from expanded repo")
	}
}

func TestParseSkillsListOutput(t *testing.T) {
	// Given: skills CLI output
	output := `
◇  Available Skills
│
│    building-native-ui
│
│      Complete guide for building beautiful apps
│
│    expo-api-routes
│
│      Guidelines for creating API routes
│
│    expo-dev-client
│
│      Build and distribute development clients
│
└  Use --skill <name> to install specific skills
`

	// When: parsing output
	skills := parseSkillsListOutput(output)

	// Then: correct skills should be extracted
	expected := []string{"building-native-ui", "expo-api-routes", "expo-dev-client"}
	if len(skills) != len(expected) {
		t.Errorf("got %d skills, expected %d: %v", len(skills), len(expected), skills)
		return
	}
	for i, skill := range expected {
		if skills[i] != skill {
			t.Errorf("skills[%d] = %q, expected %q", i, skills[i], skill)
		}
	}
}

func TestParseSkillsListOutputEmpty(t *testing.T) {
	// Given: empty skills output
	output := `
◇  Available Skills
│
└  Use --skill <name> to install specific skills
`

	// When: parsing output
	skills := parseSkillsListOutput(output)

	// Then: no skills should be returned
	if len(skills) != 0 {
		t.Errorf("got %d skills, expected 0", len(skills))
	}
}

func TestIsValidSkillName(t *testing.T) {
	tests := []struct {
		name     string
		given    string
		expected bool
	}{
		{"valid hyphenated", "frontend-design", true},
		{"valid hyphenated 2", "skill-creator", true},
		{"valid hyphenated 3", "expo-api-routes", true},
		{"valid underscore", "valid_skill", true},
		{"valid with numbers", "skill123", true},
		{"invalid uppercase", "Invalid", false},
		{"invalid space", "has space", false},
		{"invalid dot", "has.dot", false},
		{"invalid empty", "", false},
		{"invalid slash", "has/slash", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: a skill name
			// When: validating
			result := isValidSkillName(tt.given)

			// Then: result should match
			if result != tt.expected {
				t.Errorf("isValidSkillName(%q) = %v, expected %v", tt.given, result, tt.expected)
			}
		})
	}
}

func TestSkillsInToolsList(t *testing.T) {
	// Given: skills tool is defined
	pkg := config.GetToolByName("skills")

	// When: checking tool
	// Then: it should exist
	if pkg == nil {
		t.Fatal("'skills' tool not found in Tools list")
	}

	// Then: it should use npm install method
	if pkg.Method != config.InstallNpm {
		t.Errorf("method = %v, expected InstallNpm", pkg.Method)
	}

	// Then: it should be in AI category
	if pkg.Category != config.CategoryAI {
		t.Errorf("category = %v, expected CategoryAI", pkg.Category)
	}

	// Then: it should have correct command
	if pkg.Command != "skills" {
		t.Errorf("command = %q, expected 'skills'", pkg.Command)
	}
}

func TestIsFavoriteSkill(t *testing.T) {
	tests := []struct {
		name     string
		repo     string
		skill    string
		expected bool
	}{
		{"exact match", "anthropics/skills", "frontend-design", true},
		{"skill only match", "", "frontend-design", true},
		{"wrong repo", "wrong/repo", "frontend-design", false},
		{"nonexistent skill", "anthropics/skills", "nonexistent", false},
		{"empty both", "", "nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: repo and skill names
			// When: checking if favorite
			result := config.IsFavoriteSkill(tt.repo, tt.skill)

			// Then: result should match
			if result != tt.expected {
				t.Errorf("IsFavoriteSkill(%q, %q) = %v, expected %v", tt.repo, tt.skill, result, tt.expected)
			}
		})
	}
}
