package commands

import (
	"testing"
)

func TestMySkillsNotEmpty(t *testing.T) {
	if len(MySkills) == 0 {
		t.Error("MySkills should not be empty")
	}
}

func TestMySkillsHaveValidRepos(t *testing.T) {
	for _, s := range MySkills {
		if s.Repo == "" {
			t.Errorf("MySkill %q has empty repo", s.Skill)
		}
		if s.Skill == "" {
			t.Errorf("MySkill in repo %q has empty skill name", s.Repo)
		}
		// Verify the skill exists in the corresponding repo
		repo := GetSkillRepoByName(s.Repo)
		if repo == nil {
			t.Errorf("MySkill %q references unknown repo %q", s.Skill, s.Repo)
			continue
		}
		found := false
		for _, skill := range repo.Skills {
			if skill == s.Skill {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("MySkill %q not found in repo %q skills list", s.Skill, s.Repo)
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
		if len(repo.Skills) == 0 {
			t.Errorf("SkillRepo %s should have at least one skill", repo.Name)
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

func TestFindRepoForSkill(t *testing.T) {
	tests := []struct {
		skill    string
		expected string
	}{
		{"frontend-design", "anthropics/skills"},
		{"vercel-react-best-practices", "vercel-labs/agent-skills"},
		{"remotion-best-practices", "remotion-dev/skills"},
		{"nonexistent-skill", ""},
	}

	for _, tt := range tests {
		repo := findRepoForSkill(tt.skill)
		if repo != tt.expected {
			t.Errorf("findRepoForSkill(%q) = %q, expected %q", tt.skill, repo, tt.expected)
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

func TestSkillsModelBuildItems(t *testing.T) {
	m := &skillsModel{
		expanded:  make(map[string]bool),
		installed: []string{},
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

	// Should have My Skills header if MySkills is not empty
	if len(MySkills) > 0 {
		hasMySkillsHeader := false
		for _, item := range items {
			if item.itemType == itemTypeHeader && item.description == "My Skills" {
				hasMySkillsHeader = true
				break
			}
		}
		if !hasMySkillsHeader {
			t.Error("Missing 'My Skills' header")
		}
	}

	// Should have Browse header
	hasReposHeader := false
	for _, item := range items {
		if item.itemType == itemTypeHeader && item.description == "Browse" {
			hasReposHeader = true
			break
		}
	}
	if !hasReposHeader {
		t.Error("Missing 'Repositories' header")
	}
}

func TestSkillsModelBuildItemsWithInstalled(t *testing.T) {
	// Test with a skill that's in MySkills - should show in My Skills section
	m := &skillsModel{
		expanded:  make(map[string]bool),
		installed: []string{"frontend-design"},
		
	}

	items := m.buildItems()

	// frontend-design is in MySkills, so it should show there as installed
	hasMySkillsHeader := false
	hasInstalledSkill := false
	for _, item := range items {
		if item.itemType == itemTypeHeader && item.description == "My Skills" {
			hasMySkillsHeader = true
		}
		if item.itemType == itemTypeSkill && item.skill == "frontend-design" && item.installed {
			hasInstalledSkill = true
		}
	}
	if !hasMySkillsHeader {
		t.Error("Missing 'My Skills' header")
	}
	if !hasInstalledSkill {
		t.Error("Installed skill 'frontend-design' not found in items")
	}
}

func TestSkillsModelBuildItemsWithOtherInstalled(t *testing.T) {
	// Test with a skill NOT in MySkills - should show in Other Installed section
	m := &skillsModel{
		expanded:  make(map[string]bool),
		installed: []string{"remotion-best-practices"}, // Not in MySkills
		
	}

	items := m.buildItems()

	hasOtherInstalledHeader := false
	hasInstalledSkill := false
	for _, item := range items {
		if item.itemType == itemTypeHeader && item.description == "Installed" {
			hasOtherInstalledHeader = true
		}
		if item.itemType == itemTypeSkill && item.skill == "remotion-best-practices" && item.installed {
			hasInstalledSkill = true
		}
	}
	if !hasOtherInstalledHeader {
		t.Error("Missing 'Other Installed' header when non-MySkills are installed")
	}
	if !hasInstalledSkill {
		t.Error("Installed skill 'remotion-best-practices' not found in items")
	}
}

func TestSkillsModelBuildItemsWithExpanded(t *testing.T) {
	m := &skillsModel{
		expanded:  map[string]bool{"anthropics/skills": true},
		installed: []string{},
		
	}

	items := m.buildItems()

	// Should have Install All action for expanded repo
	hasRepoAction := false
	hasSkillsFromRepo := false
	for _, item := range items {
		if item.itemType == itemTypeRepoAction && item.repo == "anthropics/skills" {
			hasRepoAction = true
		}
		if item.itemType == itemTypeSkill && item.repo == "anthropics/skills" {
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
