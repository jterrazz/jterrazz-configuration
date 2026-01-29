package commands

import (
	"testing"
)

func TestMySkillsNotEmpty(t *testing.T) {
	if len(MySkills) == 0 {
		t.Error("MySkills should not be empty")
	}
}

func TestMySkillsHaveValidFields(t *testing.T) {
	for _, s := range MySkills {
		if s.Repo == "" {
			t.Errorf("MySkill %q has empty repo", s.Skill)
		}
		if s.Skill == "" {
			t.Errorf("MySkill in repo %q has empty skill name", s.Repo)
		}
	}
}

func TestSkillReposNotEmpty(t *testing.T) {
	repos := GetAllSkillRepos()
	if len(repos) == 0 {
		t.Error("SkillRepos should not be empty")
	}
}

func TestSkillReposHaveRequiredFields(t *testing.T) {
	for _, repo := range SkillRepos {
		if repo.Name == "" {
			t.Error("SkillRepo.Name should not be empty")
		}
		if repo.Description == "" {
			t.Errorf("SkillRepo %s should have a description", repo.Name)
		}
	}
}

func TestSkillReposHaveValidFormat(t *testing.T) {
	for _, repo := range SkillRepos {
		// Repo name should be in "owner/repo" format
		if !containsSlash(repo.Name) {
			t.Errorf("SkillRepo.Name %q should be in owner/repo format", repo.Name)
		}
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
		expected bool
	}{
		{"anthropics/skills", true},
		{"vercel-labs/agent-skills", true},
		{"nonexistent/repo", false},
	}

	for _, tt := range tests {
		repo := GetSkillRepoByName(tt.name)
		if tt.expected && repo == nil {
			t.Errorf("GetSkillRepoByName(%q) returned nil, expected repo", tt.name)
		}
		if !tt.expected && repo != nil {
			t.Errorf("GetSkillRepoByName(%q) returned repo, expected nil", tt.name)
		}
	}
}

func TestSkillsModelIsInstalled(t *testing.T) {
	m := &skillsModel{
		installed: []string{"skill-a", "skill-b", "skill-c"},
	}

	tests := []struct {
		skill    string
		expected bool
	}{
		{"skill-a", true},
		{"skill-b", true},
		{"skill-c", true},
		{"skill-d", false},
		{"", false},
	}

	for _, tt := range tests {
		result := m.isInstalled(tt.skill)
		if result != tt.expected {
			t.Errorf("isInstalled(%q) = %v, expected %v", tt.skill, result, tt.expected)
		}
	}
}

func TestSkillsModelFindRepoForSkill(t *testing.T) {
	m := &skillsModel{
		repoSkills: map[string][]string{
			"anthropics/skills":        {"frontend-design", "skill-creator"},
			"vercel-labs/agent-skills": {"vercel-react-best-practices"},
		},
	}

	tests := []struct {
		skill    string
		expected string
	}{
		{"frontend-design", "anthropics/skills"},
		{"vercel-react-best-practices", "vercel-labs/agent-skills"},
		{"nonexistent-skill", ""},
	}

	for _, tt := range tests {
		repo := m.findRepoForSkill(tt.skill)
		if repo != tt.expected {
			t.Errorf("findRepoForSkill(%q) = %q, expected %q", tt.skill, repo, tt.expected)
		}
	}
}

func TestSkillsModelFindRepoForSkillFromMySkills(t *testing.T) {
	// When skill is in MySkills but not in cached repoSkills
	m := &skillsModel{
		repoSkills: map[string][]string{},
	}

	// frontend-design is in MySkills with repo "anthropics/skills"
	repo := m.findRepoForSkill("frontend-design")
	if repo != "anthropics/skills" {
		t.Errorf("findRepoForSkill('frontend-design') = %q, expected 'anthropics/skills'", repo)
	}
}

func TestSkillsModelBuildItems(t *testing.T) {
	m := &skillsModel{
		expanded:   make(map[string]bool),
		installed:  []string{},
		repoSkills: make(map[string][]string),
	}

	items := m.buildItems()

	// Should have at least: Actions header + 2 actions + My Skills header + skills + Browse header + repos
	if len(items) < 4 {
		t.Errorf("buildItems() returned %d items, expected at least 4", len(items))
	}

	// First item should be Actions header
	if items[0].itemType != itemTypeHeader || items[0].description != "Actions" {
		t.Error("First item should be Actions header")
	}

	// Should have Install My Skills and Remove All actions
	hasInstallMySkills := false
	hasRemoveAll := false
	for _, item := range items {
		if item.itemType == itemTypeAction && item.actionType == "install-my-skills" {
			hasInstallMySkills = true
		}
		if item.itemType == itemTypeAction && item.actionType == "remove-all" {
			hasRemoveAll = true
		}
	}
	if !hasInstallMySkills {
		t.Error("Missing 'Install My Skills' action")
	}
	if !hasRemoveAll {
		t.Error("Missing 'Remove All' action")
	}

	// Should have Favorites header if MySkills is not empty
	if len(MySkills) > 0 {
		hasFavoritesHeader := false
		for _, item := range items {
			if item.itemType == itemTypeHeader && item.description == "Favorites" {
				hasFavoritesHeader = true
				break
			}
		}
		if !hasFavoritesHeader {
			t.Error("Missing 'Favorites' header")
		}
	}

	// Should have Browse header
	hasBrowseHeader := false
	for _, item := range items {
		if item.itemType == itemTypeHeader && item.description == "Browse" {
			hasBrowseHeader = true
			break
		}
	}
	if !hasBrowseHeader {
		t.Error("Missing 'Browse' header")
	}
}

func TestSkillsModelBuildItemsWithInstalled(t *testing.T) {
	// Test with a skill that's in MySkills - should show in Favorites section
	m := &skillsModel{
		expanded:   make(map[string]bool),
		installed:  []string{"frontend-design"},
		repoSkills: make(map[string][]string),
	}

	items := m.buildItems()

	// frontend-design is in MySkills, so it should show there as installed
	hasFavoritesHeader := false
	hasInstalledSkill := false
	for _, item := range items {
		if item.itemType == itemTypeHeader && item.description == "Favorites" {
			hasFavoritesHeader = true
		}
		if item.itemType == itemTypeSkill && item.skill == "frontend-design" && item.installed {
			hasInstalledSkill = true
		}
	}
	if !hasFavoritesHeader {
		t.Error("Missing 'Favorites' header")
	}
	if !hasInstalledSkill {
		t.Error("Installed skill 'frontend-design' not found in items")
	}
}

func TestSkillsModelBuildItemsWithOtherInstalled(t *testing.T) {
	// Test with a skill NOT in MySkills - should show in Installed section
	m := &skillsModel{
		expanded:   make(map[string]bool),
		installed:  []string{"some-other-skill"},
		repoSkills: make(map[string][]string),
	}

	items := m.buildItems()

	hasInstalledHeader := false
	hasInstalledSkill := false
	for _, item := range items {
		if item.itemType == itemTypeHeader && item.description == "Installed" {
			hasInstalledHeader = true
		}
		if item.itemType == itemTypeSkill && item.skill == "some-other-skill" && item.installed {
			hasInstalledSkill = true
		}
	}
	if !hasInstalledHeader {
		t.Error("Missing 'Installed' header when non-MySkills are installed")
	}
	if !hasInstalledSkill {
		t.Error("Installed skill 'some-other-skill' not found in items")
	}
}

func TestSkillsModelBuildItemsWithExpanded(t *testing.T) {
	m := &skillsModel{
		expanded:  map[string]bool{"anthropics/skills": true},
		installed: []string{},
		repoSkills: map[string][]string{
			"anthropics/skills": {"frontend-design", "skill-creator"},
		},
	}

	items := m.buildItems()

	// Should have Install All action for expanded repo
	hasRepoAction := false
	hasSkillsFromRepo := false
	for _, item := range items {
		if item.itemType == itemTypeRepoAction && item.repo == "anthropics/skills" {
			hasRepoAction = true
		}
		if item.itemType == itemTypeSkill && item.repo == "anthropics/skills" && item.isNested {
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

func TestSkillsModelBuildItemsExpandedWithoutCache(t *testing.T) {
	// When repo is expanded but skills not yet fetched (nil cache)
	m := &skillsModel{
		expanded:   map[string]bool{"anthropics/skills": true},
		installed:  []string{},
		repoSkills: make(map[string][]string), // Empty cache, nil for this repo
	}

	items := m.buildItems()

	// Should NOT have skills listed (they're not loaded yet)
	hasSkillsFromRepo := false
	for _, item := range items {
		if item.itemType == itemTypeSkill && item.repo == "anthropics/skills" && item.isNested {
			hasSkillsFromRepo = true
		}
	}
	if hasSkillsFromRepo {
		t.Error("Should not show skills from repo when not yet cached")
	}
}

func TestParseSkillsListOutput(t *testing.T) {
	// Sample output similar to what skills CLI produces
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

	skills := parseSkillsListOutput(output)

	expected := []string{"building-native-ui", "expo-api-routes", "expo-dev-client"}
	if len(skills) != len(expected) {
		t.Errorf("parseSkillsListOutput() returned %d skills, expected %d: %v", len(skills), len(expected), skills)
		return
	}

	for i, skill := range expected {
		if skills[i] != skill {
			t.Errorf("parseSkillsListOutput()[%d] = %q, expected %q", i, skills[i], skill)
		}
	}
}

func TestParseSkillsListOutputEmpty(t *testing.T) {
	output := `
◇  Available Skills
│
└  Use --skill <name> to install specific skills
`
	skills := parseSkillsListOutput(output)
	if len(skills) != 0 {
		t.Errorf("parseSkillsListOutput() returned %d skills for empty output, expected 0", len(skills))
	}
}

func TestIsValidSkillName(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"frontend-design", true},
		{"skill-creator", true},
		{"expo-api-routes", true},
		{"valid_skill", true},
		{"skill123", true},
		{"Invalid", false},   // uppercase
		{"has space", false}, // space
		{"has.dot", false},   // dot
		{"", false},          // empty
		{"has/slash", false}, // slash
	}

	for _, tt := range tests {
		result := isValidSkillName(tt.name)
		if result != tt.expected {
			t.Errorf("isValidSkillName(%q) = %v, expected %v", tt.name, result, tt.expected)
		}
	}
}

func TestSkillsInPackageList(t *testing.T) {
	// Verify 'skills' package is in the install list
	pkg := GetPackageByName("skills")
	if pkg == nil {
		t.Fatal("'skills' package not found in Packages list")
	}

	if pkg.Method != InstallNpm {
		t.Errorf("skills package method = %v, expected InstallNpm", pkg.Method)
	}

	if pkg.Category != CategoryAI {
		t.Errorf("skills package category = %v, expected CategoryAI", pkg.Category)
	}

	if pkg.Command != "skills" {
		t.Errorf("skills package command = %q, expected 'skills'", pkg.Command)
	}
}
