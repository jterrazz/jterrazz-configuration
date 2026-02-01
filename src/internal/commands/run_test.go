package commands

import (
	"testing"

	"github.com/spf13/cobra"
)

func findCommand(parent *cobra.Command, name string) *cobra.Command {
	for _, cmd := range parent.Commands() {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}

func TestRunCommand(t *testing.T) {
	// Given: run command exists under root
	runCmd := findCommand(rootCmd, "run")

	// When: checking command
	// Then: it should not be nil
	if runCmd == nil {
		t.Fatal("run command not found")
	}
}

func TestRunSubcommands(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"git subcommand", "git"},
		{"docker subcommand", "docker"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: run command with subcommands
			runCmd := findCommand(rootCmd, "run")

			// When: checking for subcommand
			found := findCommand(runCmd, tt.want)

			// Then: it should exist
			if found == nil {
				t.Errorf("missing subcommand: %s", tt.want)
			}
		})
	}
}

func TestGitSubcommands(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"feat", "feat"},
		{"fix", "fix"},
		{"chore", "chore"},
		{"push", "push"},
		{"sync", "sync"},
		{"wip", "wip"},
		{"unwip", "unwip"},
		{"status", "status"},
		{"log", "log"},
		{"branches", "branches"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: git command under run
			runCmd := findCommand(rootCmd, "run")
			gitCmd := findCommand(runCmd, "git")
			if gitCmd == nil {
				t.Fatal("git command not found")
			}

			// When: checking for subcommand
			found := findCommand(gitCmd, tt.want)

			// Then: it should exist
			if found == nil {
				t.Errorf("missing subcommand: %s", tt.want)
			}
		})
	}
}

func TestDockerSubcommands(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"rm", "rm"},
		{"rmi", "rmi"},
		{"clean", "clean"},
		{"reset", "reset"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: docker command under run
			runCmd := findCommand(rootCmd, "run")
			dockerCmd := findCommand(runCmd, "docker")
			if dockerCmd == nil {
				t.Fatal("docker command not found")
			}

			// When: checking for subcommand
			found := findCommand(dockerCmd, tt.want)

			// Then: it should exist
			if found == nil {
				t.Errorf("missing subcommand: %s", tt.want)
			}
		})
	}
}
