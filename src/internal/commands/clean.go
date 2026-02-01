package commands

import (
	"fmt"

	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/tool"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
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
		var all []string
		for _, c := range config.Cleanables {
			all = append(all, c.Name)
		}
		return tool.FilterStrings(all, args), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if cleanAll {
			fmt.Println(ui.Cyan("🧹 Cleaning everything..."))
			for _, c := range config.Cleanables {
				runCleanable(c)
			}
			fmt.Println(ui.Green("✅ System cleanup completed"))
			return
		}

		if len(args) == 0 {
			listCleanItems()
			return
		}

		fmt.Println(ui.Cyan("🧹 Cleaning selected items..."))
		for _, name := range args {
			c := config.GetCleanableByName(name)
			if c == nil {
				ui.PrintError(fmt.Sprintf("Unknown clean item: %s", name))
				continue
			}
			runCleanable(*c)
		}
		fmt.Println(ui.Green("✅ Cleanup completed"))
	},
}

func init() {
	cleanCmd.Flags().BoolVarP(&cleanAll, "all", "a", false, "Clean everything")
	rootCmd.AddCommand(cleanCmd)
}

func listCleanItems() {
	fmt.Println(ui.Cyan("Available clean items:"))
	fmt.Println()

	for _, c := range config.Cleanables {
		available := c.RequiresCmd == "" || config.CommandExists(c.RequiresCmd)
		status := ui.Red("✗")
		if available {
			status = ui.Green("✓")
		}
		fmt.Printf("  %s %-14s %s\n", status, c.Name, ui.Dim(c.Description))
	}

	fmt.Println()
	fmt.Println(ui.Dim("Usage: j clean <item> [item...]"))
	fmt.Println(ui.Dim("       j clean --all"))
}

func runCleanable(c config.Cleanable) {
	if c.RequiresCmd != "" && !config.CommandExists(c.RequiresCmd) {
		ui.PrintWarning(fmt.Sprintf("%s not found, skipping", c.RequiresCmd))
		return
	}

	fmt.Println(ui.Cyan(fmt.Sprintf("🧹 %s...", c.Description)))
	if c.CleanFn != nil {
		c.CleanFn()
	}
	fmt.Println(ui.Green(fmt.Sprintf("  ✅ %s completed", c.Name)))
}
