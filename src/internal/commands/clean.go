package commands

import (
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
			ui.PrintAction("🧹", "Cleaning everything...")
			for _, c := range config.Cleanables {
				runCleanable(c)
			}
			ui.PrintDone("System cleanup completed")
			return
		}

		if len(args) == 0 {
			listCleanItems()
			return
		}

		ui.PrintAction("🧹", "Cleaning selected items...")
		for _, name := range args {
			c := config.GetCleanableByName(name)
			if c == nil {
				ui.PrintError("Unknown clean item: " + name)
				continue
			}
			runCleanable(*c)
		}
		ui.PrintDone("Cleanup completed")
	},
}

func init() {
	cleanCmd.Flags().BoolVarP(&cleanAll, "all", "a", false, "Clean everything")
	rootCmd.AddCommand(cleanCmd)
}

func listCleanItems() {
	ui.PrintInfo("Available clean items:")
	ui.PrintEmpty()

	for _, c := range config.Cleanables {
		available := c.RequiresCmd == "" || config.CommandExists(c.RequiresCmd)
		ui.PrintRow(available, c.Name, c.Description)
	}

	ui.PrintEmpty()
	ui.PrintUsage(
		"Usage: j clean <item> [item...]",
		"       j clean --all",
	)
}

func runCleanable(c config.Cleanable) {
	if c.RequiresCmd != "" && !config.CommandExists(c.RequiresCmd) {
		ui.PrintWarning(c.RequiresCmd + " not found, skipping")
		return
	}

	ui.PrintAction("🧹", c.Description+"...")
	if c.CleanFn != nil {
		c.CleanFn()
	}
	ui.PrintRow(true, c.Name, "completed")
}
