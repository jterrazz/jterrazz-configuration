package commands

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCommandStructure(t *testing.T) {
	// Verify root command exists
	if rootCmd == nil {
		t.Fatal("rootCmd is nil")
	}

	if rootCmd.Use != "j" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "j")
	}

	// Verify expected subcommands are registered
	expectedCommands := []string{
		"status",
		"update",
		"clean",
		"setup",
		"git",
		"docker",
	}

	commands := rootCmd.Commands()
	commandNames := make(map[string]bool)
	for _, cmd := range commands {
		commandNames[cmd.Use] = true
	}

	for _, expected := range expectedCommands {
		if !commandNames[expected] {
			t.Errorf("missing subcommand: %s", expected)
		}
	}
}

func TestSetupSubcommands(t *testing.T) {
	var setupCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "setup" {
			setupCmd = cmd
			break
		}
	}

	if setupCmd == nil {
		t.Fatal("setup command not found")
	}

	expectedSubcommands := []string{
		"all",
		"brew",
		"ohmyzsh",
		"nvm",
		"git-ssh",
		"dock-spacer",
		"dock-reset",
	}

	subcommands := setupCmd.Commands()
	subcommandNames := make(map[string]bool)
	for _, cmd := range subcommands {
		subcommandNames[cmd.Use] = true
	}

	for _, expected := range expectedSubcommands {
		if !subcommandNames[expected] {
			t.Errorf("setup missing subcommand: %s", expected)
		}
	}
}

func TestGitSubcommands(t *testing.T) {
	var gitCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "git" {
			gitCmd = cmd
			break
		}
	}

	if gitCmd == nil {
		t.Fatal("git command not found")
	}

	expectedSubcommands := []string{
		"feat",
		"fix",
		"chore",
		"push",
		"sync",
		"wip",
		"unwip",
	}

	subcommands := gitCmd.Commands()
	subcommandNames := make(map[string]bool)
	for _, cmd := range subcommands {
		subcommandNames[cmd.Name()] = true
	}

	for _, expected := range expectedSubcommands {
		if !subcommandNames[expected] {
			t.Errorf("git missing subcommand: %s", expected)
		}
	}
}

func TestDockerSubcommands(t *testing.T) {
	var dockerCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "docker" {
			dockerCmd = cmd
			break
		}
	}

	if dockerCmd == nil {
		t.Fatal("docker command not found")
	}

	expectedSubcommands := []string{
		"ps",
		"images",
		"clean",
		"reset",
	}

	subcommands := dockerCmd.Commands()
	subcommandNames := make(map[string]bool)
	for _, cmd := range subcommands {
		subcommandNames[cmd.Use] = true
	}

	for _, expected := range expectedSubcommands {
		if !subcommandNames[expected] {
			t.Errorf("docker missing subcommand: %s", expected)
		}
	}
}
