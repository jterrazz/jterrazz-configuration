package commands

import (
	"testing"

	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
)

func TestSetupCommand(t *testing.T) {
	// Given: setup command is initialized
	cmd := setupCmd

	// When: checking command properties
	// Then: it should not be nil
	if cmd == nil {
		t.Fatal("setupCmd is nil")
	}

	// Then: it should have correct use name
	if cmd.Use != "setup" {
		t.Errorf("got %q, want %q", cmd.Use, "setup")
	}

	// Then: it should have a description
	if cmd.Short == "" {
		t.Error("setupCmd.Short should not be empty")
	}
}

func TestScriptsNotEmpty(t *testing.T) {
	// Given: scripts are defined in config
	scripts := config.Scripts

	// When: checking scripts list
	// Then: it should not be empty
	if len(scripts) == 0 {
		t.Error("Scripts returned empty list")
	}
}

func TestScriptsHaveRequiredFields(t *testing.T) {
	for _, script := range config.Scripts {
		t.Run(script.Name, func(t *testing.T) {
			// Given: a script from config
			// When: checking required fields
			// Then: name should not be empty
			if script.Name == "" {
				t.Error("Script has empty name")
			}

			// Then: description should not be empty
			if script.Description == "" {
				t.Errorf("Script %q has empty description", script.Name)
			}
		})
	}
}

func TestExpectedScriptsExist(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"dock-reset", "dock-reset"},
		{"dock-spacer", "dock-spacer"},
		{"ghostty-config", "ghostty-config"},
		{"gpg-setup", "gpg-setup"},
		{"hushlogin", "hushlogin"},
		{"java-symlink", "java-symlink"},
		{"ssh", "ssh"},
		{"zed-config", "zed-config"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: expected script name
			// When: looking up script
			script := config.GetScriptByName(tt.want)

			// Then: it should be found
			if script == nil {
				t.Errorf("missing expected script: %s", tt.want)
			}
		})
	}
}

func TestSetupModelInit(t *testing.T) {
	// Given: setup model is initialized
	m := newSetupModel()

	// When: checking model properties
	// Then: list should not be nil
	if m.list == nil {
		t.Error("newSetupModel() should have a list")
	}

	// Then: list should have items
	if len(m.list.Items) == 0 {
		t.Error("newSetupModel() list should have items")
	}

	// Then: cursor should be valid
	if m.list.Cursor < 0 || m.list.Cursor >= len(m.list.Items) {
		t.Errorf("cursor = %d, out of bounds", m.list.Cursor)
	}

	// Then: cursor should not be on a header
	if m.list.Items[m.list.Cursor].Kind == ui.KindHeader {
		t.Error("cursor should not be on a header")
	}
}

func TestScriptCheckFunctions(t *testing.T) {
	tests := []string{
		"ghostty-config",
		"gpg-setup",
		"hushlogin",
		"java-symlink",
		"ssh",
		"zed-config",
	}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			// Given: a script with check function
			script := config.GetScriptByName(name)
			if script == nil {
				t.Fatalf("Script %q not found", name)
			}

			// When: checking CheckFn
			// Then: it should not be nil
			if script.CheckFn == nil {
				t.Errorf("Script %q has nil CheckFn", name)
				return
			}

			// Then: calling it should not panic
			_ = script.CheckFn()
		})
	}
}

func TestUIStylesAccessible(t *testing.T) {
	// Given: UI package is imported
	// When: accessing styles
	// Then: they should be accessible without panic
	_ = ui.TitleStyle
	_ = ui.SelectedStyle
	_ = ui.SuccessStyle
}

func TestSkillsInSetupNavigation(t *testing.T) {
	// Given: setup model with items
	m := newSetupModel()

	// When: looking for skills in navigation
	foundSkills := false
	for i, item := range m.list.Items {
		if item.Kind == ui.KindAction && m.itemData[i].name == "skills" {
			foundSkills = true
			break
		}
	}

	// Then: skills should be found
	if !foundSkills {
		t.Error("skills should be in setup items (Navigation section)")
	}
}

func TestScriptCategories(t *testing.T) {
	// Given: scripts with categories
	categories := map[config.ScriptCategory]int{
		config.ScriptCategoryTerminal: 0,
		config.ScriptCategorySecurity: 0,
		config.ScriptCategoryEditor:   0,
		config.ScriptCategorySystem:   0,
	}

	for _, script := range config.Scripts {
		categories[script.Category]++
	}

	// When: checking category counts
	for cat, count := range categories {
		t.Run(string(cat), func(t *testing.T) {
			// Then: each category should have at least one script
			if count == 0 {
				t.Errorf("Category %v has no scripts", cat)
			}
		})
	}
}

func TestToolsWithScriptsReference(t *testing.T) {
	for _, tool := range config.Tools {
		for _, scriptName := range tool.Scripts {
			t.Run(tool.Name+"/"+scriptName, func(t *testing.T) {
				// Given: a tool that references a script
				// When: looking up the script
				script := config.GetScriptByName(scriptName)

				// Then: the script should exist
				if script == nil {
					t.Errorf("Tool %q references non-existent script %q", tool.Name, scriptName)
				}
			})
		}
	}
}
