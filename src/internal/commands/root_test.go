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
		"install",
		"update",
		"clean",
		"setup",
		"run",
	}

	commands := rootCmd.Commands()
	commandNames := make(map[string]bool)
	for _, cmd := range commands {
		commandNames[cmd.Name()] = true
	}

	for _, expected := range expectedCommands {
		if !commandNames[expected] {
			t.Errorf("missing subcommand: %s", expected)
		}
	}
}

func TestRunSubcommands(t *testing.T) {
	var runCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "run" {
			runCmd = cmd
			break
		}
	}

	if runCmd == nil {
		t.Fatal("run command not found")
	}

	expectedSubcommands := []string{
		"git",
		"docker",
	}

	subcommands := runCmd.Commands()
	subcommandNames := make(map[string]bool)
	for _, cmd := range subcommands {
		subcommandNames[cmd.Name()] = true
	}

	for _, expected := range expectedSubcommands {
		if !subcommandNames[expected] {
			t.Errorf("run missing subcommand: %s", expected)
		}
	}
}

func TestGitSubcommands(t *testing.T) {
	// Find run -> git
	var runCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "run" {
			runCmd = cmd
			break
		}
	}

	if runCmd == nil {
		t.Fatal("run command not found")
	}

	var gitCmd *cobra.Command
	for _, cmd := range runCmd.Commands() {
		if cmd.Name() == "git" {
			gitCmd = cmd
			break
		}
	}

	if gitCmd == nil {
		t.Fatal("git command not found under run")
	}

	expectedSubcommands := []string{
		"feat",
		"fix",
		"chore",
		"push",
		"sync",
		"wip",
		"unwip",
		"status",
		"log",
		"branches",
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
	// Find run -> docker
	var runCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "run" {
			runCmd = cmd
			break
		}
	}

	if runCmd == nil {
		t.Fatal("run command not found")
	}

	var dockerCmd *cobra.Command
	for _, cmd := range runCmd.Commands() {
		if cmd.Name() == "docker" {
			dockerCmd = cmd
			break
		}
	}

	if dockerCmd == nil {
		t.Fatal("docker command not found under run")
	}

	expectedSubcommands := []string{
		"rm",
		"rmi",
		"clean",
		"reset",
	}

	subcommands := dockerCmd.Commands()
	subcommandNames := make(map[string]bool)
	for _, cmd := range subcommands {
		subcommandNames[cmd.Name()] = true
	}

	for _, expected := range expectedSubcommands {
		if !subcommandNames[expected] {
			t.Errorf("docker missing subcommand: %s", expected)
		}
	}
}
