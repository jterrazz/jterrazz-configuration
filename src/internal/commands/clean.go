package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean system caches, Docker, Multipass, and trash",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Println(cyan("🧹 Cleaning system..."))

		// Clean Homebrew
		if commandExists("brew") {
			fmt.Println(cyan("🍺 Cleaning Homebrew cache..."))
			runCommand("brew", "cleanup")
		}

		// Clean Docker
		if commandExists("docker") {
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
		} else {
			printWarning("Docker not found, skipping Docker cleanup")
		}

		// Clean Multipass
		if commandExists("multipass") {
			fmt.Println(cyan("🖥️  Cleaning Multipass..."))
			fmt.Println("  🗑️  Removing all instances...")
			exec.Command("multipass", "delete", "--all").Run()
			fmt.Println("  🗑️  Purging deleted instances...")
			runCommand("multipass", "purge")
			fmt.Println(green("  ✅ Multipass cleanup completed"))
		} else {
			printWarning("Multipass not found, skipping Multipass cleanup")
		}

		// Empty trash
		fmt.Println(cyan("🗑️  Emptying trash..."))
		os.RemoveAll(os.Getenv("HOME") + "/.Trash")
		os.MkdirAll(os.Getenv("HOME") + "/.Trash", 0755)

		fmt.Println(green("✅ System cleanup completed"))
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
