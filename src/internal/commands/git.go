package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git workflow shortcuts",
}

var gitFeatCmd = &cobra.Command{
	Use:   "feat [message]",
	Short: "Add all and commit with 'feat:' prefix",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gitCommit("feat", strings.Join(args, " "))
	},
}

var gitFixCmd = &cobra.Command{
	Use:   "fix [message]",
	Short: "Add all and commit with 'fix:' prefix",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gitCommit("fix", strings.Join(args, " "))
	},
}

var gitChoreCmd = &cobra.Command{
	Use:   "chore [message]",
	Short: "Add all and commit with 'chore:' prefix",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gitCommit("chore", strings.Join(args, " "))
	},
}

var gitPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push current branch to origin",
	Run: func(cmd *cobra.Command, args []string) {
		runCommand("git", "push", "-u", "origin", "HEAD")
	},
}

var gitSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Fetch and pull from remote",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Println(cyan("🔄 Syncing with remote..."))
		runCommand("git", "fetch", "-p")
		runCommand("git", "pull")
	},
}

var gitWipCmd = &cobra.Command{
	Use:   "wip",
	Short: "Add all and commit as 'WIP'",
	Run: func(cmd *cobra.Command, args []string) {
		runCommand("git", "add", "--all")
		runCommand("git", "commit", "-m", "WIP")
	},
}

var gitUnwipCmd = &cobra.Command{
	Use:   "unwip",
	Short: "Undo last commit and unstage",
	Run: func(cmd *cobra.Command, args []string) {
		runCommand("git", "reset", "--soft", "HEAD~1")
		runCommand("git", "reset", "HEAD")
	},
}

var gitStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show git status",
	Run: func(cmd *cobra.Command, args []string) {
		runCommand("git", "status")
	},
}

var gitLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Show recent commits",
	Run: func(cmd *cobra.Command, args []string) {
		runCommand("git", "log", "--oneline", "-10")
	},
}

var gitBranchesCmd = &cobra.Command{
	Use:   "branches",
	Short: "List local branches",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Println(cyan("🌿 Local branches:"))
		runCommand("git", "branch")
	},
}

func init() {
	gitCmd.AddCommand(gitFeatCmd)
	gitCmd.AddCommand(gitFixCmd)
	gitCmd.AddCommand(gitChoreCmd)
	gitCmd.AddCommand(gitPushCmd)
	gitCmd.AddCommand(gitSyncCmd)
	gitCmd.AddCommand(gitWipCmd)
	gitCmd.AddCommand(gitUnwipCmd)
	gitCmd.AddCommand(gitStatusCmd)
	gitCmd.AddCommand(gitLogCmd)
	gitCmd.AddCommand(gitBranchesCmd)
	runCmd.AddCommand(gitCmd)
}

func gitCommit(prefix, message string) {
	runCommand("git", "add", ".")
	runCommand("git", "commit", "-m", fmt.Sprintf("%s: %s", prefix, message))
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
