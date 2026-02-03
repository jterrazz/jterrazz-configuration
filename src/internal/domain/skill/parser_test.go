package skill

import (
	"testing"
)

func TestParseSkillsListOutput(t *testing.T) {
	tests := []struct {
		name     string
		given    string
		expected []string
	}{
		{
			"standard output",
			`
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
`,
			[]string{"building-native-ui", "expo-api-routes", "expo-dev-client"},
		},
		{
			"empty output",
			`
◇  Available Skills
│
└  Use --skill <name> to install specific skills
`,
			[]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skills := ParseSkillsListOutput(tt.given)

			if len(skills) != len(tt.expected) {
				t.Errorf("got %d skills %v, expected %d %v", len(skills), skills, len(tt.expected), tt.expected)
				return
			}
			for i, skill := range tt.expected {
				if skills[i] != skill {
					t.Errorf("skills[%d] = %q, expected %q", i, skills[i], skill)
				}
			}
		})
	}
}

func TestIsValidName(t *testing.T) {
	tests := []struct {
		name     string
		given    string
		expected bool
	}{
		{"valid hyphenated", "frontend-design", true},
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
			result := IsValidName(tt.given)
			if result != tt.expected {
				t.Errorf("IsValidName(%q) = %v, expected %v", tt.given, result, tt.expected)
			}
		})
	}
}
