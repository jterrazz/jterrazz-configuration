package commands

import "testing"

func TestRootCommand(t *testing.T) {
	// Given: root command is initialized
	cmd := rootCmd

	// When: checking command properties
	// Then: it should not be nil
	if cmd == nil {
		t.Fatal("rootCmd is nil")
	}

	// Then: it should have correct use name
	if cmd.Use != "j" {
		t.Errorf("got %q, want %q", cmd.Use, "j")
	}
}

func TestRootSubcommands(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"status command", "status"},
		{"install command", "install"},
		{"update command", "update"},
		{"clean command", "clean"},
		{"setup command", "setup"},
		{"run command", "run"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: root command with subcommands
			commands := rootCmd.Commands()
			commandNames := make(map[string]bool)
			for _, cmd := range commands {
				commandNames[cmd.Name()] = true
			}

			// When: checking for subcommand
			found := commandNames[tt.want]

			// Then: it should be registered
			if !found {
				t.Errorf("missing subcommand: %s", tt.want)
			}
		})
	}
}
