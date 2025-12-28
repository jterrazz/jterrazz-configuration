package commands

import "testing"

func TestVersionParsers(t *testing.T) {
	tests := []struct {
		name     string
		parser   func(string) string
		input    string
		expected string
	}{
		// trimVersion
		{"trimVersion basic", trimVersion, "1.2.3", "1.2.3"},
		{"trimVersion with v", trimVersion, "v1.2.3", "1.2.3"},
		{"trimVersion with newline", trimVersion, "v1.2.3\n", "1.2.3"},

		// parseBrewVersion
		{"parseBrewVersion", parseBrewVersion, "Homebrew 4.2.0\nHomebrew/homebrew-core", "4.2.0"},
		{"parseBrewVersion single line", parseBrewVersion, "Homebrew 5.0.7", "5.0.7"},

		// parseGitVersion
		{"parseGitVersion basic", parseGitVersion, "git version 2.39.0", "2.39.0"},
		{"parseGitVersion apple", parseGitVersion, "git version 2.39.0 (Apple Git-145)", "2.39.0"},

		// parseGoVersion
		{"parseGoVersion", parseGoVersion, "go version go1.21.0 darwin/arm64", "1.21.0"},
		{"parseGoVersion linux", parseGoVersion, "go version go1.22.5 linux/amd64", "1.22.5"},

		// parsePythonVersion
		{"parsePythonVersion", parsePythonVersion, "Python 3.12.0", "3.12.0"},
		{"parsePythonVersion with newline", parsePythonVersion, "Python 3.11.5\n", "3.11.5"},

		// parseTerraformVersion
		{"parseTerraformVersion", parseTerraformVersion, "Terraform v1.5.7\non darwin_arm64", "1.5.7"},

		// parseAnsibleVersion
		{"parseAnsibleVersion core", parseAnsibleVersion, "ansible [core 2.15.0]", "2.15.0"},
		{"parseAnsibleVersion simple", parseAnsibleVersion, "ansible 2.14.0", "2.14.0"},

		// parseMultipassVersion
		{"parseMultipassVersion", parseMultipassVersion, "multipass 1.12.0+mac\nmultipassd 1.12.0+mac", "1.12.0+mac"},

		// parseCodexVersion
		{"parseCodexVersion", parseCodexVersion, "codex 0.1.0", "0.1.0"},
		{"parseCodexVersion single", parseCodexVersion, "0.1.0", "0.1.0"},

		// parseMoleVersion
		{"parseMoleVersion", parseMoleVersion, "\nMole version 1.14.5\nmacOS: 14.0", "1.14.5"},
		{"parseMoleVersion no leading newline", parseMoleVersion, "Mole version 1.13.0\nmacOS: 13.0", "1.13.0"},

		// parseClaudeVersion
		{"parseClaudeVersion", parseClaudeVersion, "2.0.76 (Claude Code)", "2.0.76"},
		{"parseClaudeVersion simple", parseClaudeVersion, "1.0.0", "1.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.parser(tt.input)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{1610612736, "1.5 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatBytes(%d) = %q, want %q", tt.bytes, result, tt.expected)
			}
		})
	}
}
