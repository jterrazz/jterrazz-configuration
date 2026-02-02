package commands

import (
	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/tool"
	"github.com/jterrazz/jterrazz-cli/internal/ui/print"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [tool...]",
	Short: "Install development tools",
	Long: `Install development tools.

Examples:
  j install homebrew        Install Homebrew
  j install nvm             Install NVM
  j install go python node  Install specific tools
  j install                 List available tools`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var all []string
		for _, t := range config.Tools {
			all = append(all, t.Name)
		}
		return tool.FilterStrings(all, args), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			listAvailableTools()
			return
		}

		print.Action("📦", "Installing selected tools...")
		for _, name := range args {
			installToolByName(name)
		}
		print.Done("Done")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func listAvailableTools() {
	print.Info("Available tools:")
	print.Empty()

	currentCategory := config.ToolCategory("")
	for _, t := range config.Tools {
		if t.Category != currentCategory {
			currentCategory = t.Category
			print.Category(string(currentCategory))
		}

		result := t.Check()
		print.Row(result.Installed, t.Name, t.Method.String())
	}

	print.Empty()
	print.Usage("Usage: j install <tool> [tool...]")
}

func installToolByName(name string) {
	// Handle "brew" as alias for "homebrew"
	if name == "brew" {
		name = "homebrew"
	}

	t := config.GetToolByName(name)
	if t == nil {
		print.Error("Unknown tool: " + name)
		return
	}

	result := t.Check()
	if result.Installed {
		print.Row(true, t.Name, "already installed")
		return
	}

	// Check dependencies
	for _, depName := range t.Dependencies {
		depTool := config.GetToolByName(depName)
		if depTool == nil {
			continue
		}
		depResult := depTool.Check()
		if !depResult.Installed {
			print.Error(depName + " required for " + t.Name + ". Run: j install " + depName)
			return
		}
	}

	print.Installing(t.Name)
	if err := t.Install(); err != nil {
		print.Error("Failed to install " + t.Name + ": " + err.Error())
	} else {
		print.Row(true, t.Name, "installed")
		// Run post-install scripts
		for _, scriptName := range t.Scripts {
			runSetupItem(scriptName)
		}
	}
}
