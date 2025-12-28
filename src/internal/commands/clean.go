package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var cleanAll bool

var cleanCmd = &cobra.Command{
	Use:   "clean [item...]",
	Short: "Clean system caches, Docker, Multipass, and trash",
	Long: `Clean system caches and resources.

Examples:
  j clean --all              Clean everything
  j clean brew               Clean Homebrew cache
  j clean docker             Clean Docker resources
  j clean multipass          Clean Multipass instances
  j clean trash              Empty trash
  j clean brew docker        Clean specific items
  j clean                    List available clean items`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		suggestions := []string{"brew", "docker", "multipass", "trash"}
		var filtered []string
		for _, s := range suggestions {
			alreadyUsed := false
			for _, arg := range args {
				if arg == s {
					alreadyUsed = true
					break
				}
			}
			if !alreadyUsed {
				filtered = append(filtered, s)
			}
		}
		return filtered, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		if cleanAll {
			fmt.Println(cyan("🧹 Cleaning everything..."))
			cleanBrew()
			cleanDocker()
			cleanMultipass()
			cleanTrash()
			fmt.Println(green("✅ System cleanup completed"))
			return
		}

		if len(args) == 0 {
			listCleanItems()
			return
		}

		fmt.Println(cyan("🧹 Cleaning selected items..."))
		for _, name := range args {
			runCleanItem(name)
		}
		fmt.Println(green("✅ Cleanup completed"))
	},
}

func init() {
	cleanCmd.Flags().BoolVarP(&cleanAll, "all", "a", false, "Clean everything")
	rootCmd.AddCommand(cleanCmd)
}

func listCleanItems() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	dim := color.New(color.FgHiBlack).SprintFunc()

	fmt.Println(cyan("Available clean items:"))
	fmt.Println()

	items := []struct {
		name        string
		description string
		available   bool
	}{
		{"brew", "Clean Homebrew cache", commandExists("brew")},
		{"docker", "Clean Docker containers, images, volumes", commandExists("docker")},
		{"multipass", "Remove all Multipass instances", commandExists("multipass")},
		{"trash", "Empty system trash", true},
	}

	for _, item := range items {
		status := red("✗")
		if item.available {
			status = green("✓")
		}
		fmt.Printf("  %s %-14s %s\n", status, item.name, dim(item.description))
	}

	fmt.Println()
	fmt.Println(dim("Usage: j clean <item> [item...]"))
	fmt.Println(dim("       j clean --all"))
}

func runCleanItem(name string) {
	switch name {
	case "brew":
		cleanBrew()
	case "docker":
		cleanDocker()
	case "multipass":
		cleanMultipass()
	case "trash":
		cleanTrash()
	default:
		printError(fmt.Sprintf("Unknown clean item: %s", name))
	}
}

func cleanBrew() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	if !commandExists("brew") {
		printWarning("Homebrew not found, skipping")
		return
	}
	fmt.Println(cyan("🍺 Cleaning Homebrew cache..."))
	runCommand("brew", "cleanup")
	fmt.Println(green("  ✅ Homebrew cleanup completed"))
}

func cleanDocker() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	if !commandExists("docker") {
		printWarning("Docker not found, skipping")
		return
	}
	fmt.Println(cyan("🐳 Cleaning Docker..."))
	fmt.Println("  🗑️  Removing stopped containers...")
	runCommand("docker", "container", "prune", "-f")
	fmt.Println("  🗑️  Removing unused images...")
	runCommand("docker", "image", "prune", "-f")
	fmt.Println("  🗑️  Removing unused volumes...")
	runCommand("docker", "volume", "prune", "-f")
	fmt.Println("  🗑️  Removing unused networks...")
	runCommand("docker", "network", "prune", "-f")
	fmt.Println("  🗑️  Cleaning build cache...")
	runCommand("docker", "builder", "prune", "-f")
	fmt.Println(green("  ✅ Docker cleanup completed"))
}

func cleanMultipass() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	if !commandExists("multipass") {
		printWarning("Multipass not found, skipping")
		return
	}
	fmt.Println(cyan("🖥️  Cleaning Multipass..."))
	fmt.Println("  🗑️  Removing all instances...")
	exec.Command("multipass", "delete", "--all").Run()
	fmt.Println("  🗑️  Purging deleted instances...")
	runCommand("multipass", "purge")
	fmt.Println(green("  ✅ Multipass cleanup completed"))
}

func cleanTrash() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	fmt.Println(cyan("🗑️  Emptying trash..."))
	os.RemoveAll(os.Getenv("HOME") + "/.Trash")
	os.MkdirAll(os.Getenv("HOME")+"/.Trash", 0755)
	fmt.Println(green("  ✅ Trash emptied"))
}
