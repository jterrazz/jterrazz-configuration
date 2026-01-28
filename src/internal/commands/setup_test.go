package commands

import (
	"testing"
)

func TestSetupCommandDefined(t *testing.T) {
	if setupCmd == nil {
		t.Fatal("setupCmd is nil")
	}

	if setupCmd.Use != "setup" {
		t.Errorf("setupCmd.Use = %q, expected 'setup'", setupCmd.Use)
	}

	if setupCmd.Short == "" {
		t.Error("setupCmd.Short should not be empty")
	}
}

func TestGetSetupItems(t *testing.T) {
	items := getSetupItems()

	if len(items) == 0 {
		t.Error("getSetupItems() returned empty list")
	}

	// Check that all items have required fields
	for _, item := range items {
		if item.name == "" {
			t.Error("Setup item has empty name")
		}
		if item.description == "" {
			t.Errorf("Setup item %q has empty description", item.name)
		}
	}
}

func TestGetSetupItemsContainsExpectedItems(t *testing.T) {
	items := getSetupItems()

	expectedItems := []string{
		"dock-reset",
		"dock-spacer",
		"ghostty",
		"gpg",
		"hushlogin",
		"java",
		"ohmyzsh",
		"skills",
		"ssh",
		"zed",
	}

	itemNames := make(map[string]bool)
	for _, item := range items {
		itemNames[item.name] = true
	}

	for _, expected := range expectedItems {
		if !itemNames[expected] {
			t.Errorf("Missing expected setup item: %s", expected)
		}
	}
}

func TestSetupModelInit(t *testing.T) {
	m := initialSetupModel()

	if len(m.items) == 0 {
		t.Error("initialSetupModel() should have items")
	}

	// Cursor should start on first non-header item (after Actions header)
	if m.cursor < 0 || m.cursor >= len(m.items) {
		t.Errorf("initialSetupModel() cursor = %d, out of bounds", m.cursor)
	}

	// Cursor should not be on a header
	if m.items[m.cursor].itemType == setupItemTypeHeader {
		t.Error("initialSetupModel() cursor should not be on a header")
	}
}

func TestCheckFunctions(t *testing.T) {
	// These functions should not panic and should return valid pointers
	tests := []struct {
		name  string
		check func() *bool
	}{
		{"checkGhostty", checkGhostty},
		{"checkGPG", checkGPG},
		{"checkHushlogin", checkHushlogin},
		{"checkJava", checkJava},
		{"checkOhMyZsh", checkOhMyZsh},
		{"checkSSH", checkSSH},
		{"checkZed", checkZed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.check()
			if result == nil {
				t.Errorf("%s() returned nil, expected *bool", tt.name)
			}
		})
	}
}

func TestUIStylesInitialized(t *testing.T) {
	// Just verify shared UI styles are accessible
	_ = uiTitleStyle
	_ = uiSelectedStyle
	_ = uiSuccessStyle
}

func TestSkillsItemIsAction(t *testing.T) {
	items := getSetupItems()

	for _, item := range items {
		if item.name == "skills" {
			if !item.isAction {
				t.Error("skills item should have isAction=true")
			}
			if item.configured != nil {
				t.Error("skills item should have configured=nil")
			}
			return
		}
	}
	t.Error("skills item not found in setup items")
}
