package commands

import (
	"fmt"

	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
	"github.com/spf13/cobra"
)

var installAll bool

var installCmd = &cobra.Command{
	Use:   "install [tool...]",
	Short: "Install development tools",
	Long: `Install development tools.

Examples:
  j install --all           Install all tools
  j install homebrew        Install Homebrew
  j install nvm             Install NVM
  j install go python node  Install specific tools
  j install                 List available tools`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var all []string
		for _, tool := range config.Tools {
			all = append(all, tool.Name)
		}
		return filterUsedArgs(all, args), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if installAll {
			fmt.Println(ui.Cyan("🚀 Installing all development tools..."))
			installAllTools()
			return
		}

		if len(args) == 0 {
			listAvailableTools()
			return
		}

		fmt.Println(ui.Cyan("📦 Installing selected tools..."))
		for _, name := range args {
			installToolByName(name)
		}
		fmt.Println(ui.Green("✅ Done"))
	},
}

func init() {
	installCmd.Flags().BoolVarP(&installAll, "all", "a", false, "Install all tools")
	rootCmd.AddCommand(installCmd)
}

func listAvailableTools() {
	fmt.Println(ui.Cyan("Available tools:"))
	fmt.Println()

	currentCategory := config.ToolCategory("")
	for _, tool := range config.Tools {
		if tool.Category != currentCategory {
			currentCategory = tool.Category
			fmt.Printf("%s\n", ui.Dim(string(currentCategory)))
		}

		result := tool.Check()
		status := ui.Red("✗")
		if result.Installed {
			status = ui.Green("✓")
		}

		fmt.Printf("  %s %-14s %s\n", status, tool.Name, ui.Dim(tool.Method.String()))
	}

	fmt.Println()
	fmt.Println(ui.Dim("Usage: j install <tool> [tool...]"))
	fmt.Println(ui.Dim("       j install --all"))
}

func installToolByName(name string) {
	// Handle "brew" as alias for "homebrew"
	if name == "brew" {
		name = "homebrew"
	}

	tool := config.GetToolByName(name)
	if tool == nil {
		ui.PrintError(fmt.Sprintf("Unknown tool: %s", name))
		return
	}

	result := tool.Check()
	if result.Installed {
		fmt.Printf("  %s %s already installed\n", ui.Green("✓"), tool.Name)
		return
	}

	// Check dependencies
	for _, depName := range tool.Dependencies {
		depTool := config.GetToolByName(depName)
		if depTool == nil {
			continue
		}
		depResult := depTool.Check()
		if !depResult.Installed {
			ui.PrintError(fmt.Sprintf("%s required for %s. Run: j install %s", depName, tool.Name, depName))
			return
		}
	}

	fmt.Printf("  📥 Installing %s...\n", tool.Name)
	if err := tool.Install(); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to install %s: %v", tool.Name, err))
	} else {
		fmt.Printf("  %s %s installed\n", ui.Green("✓"), tool.Name)
		// Run post-install scripts
		for _, scriptName := range tool.Scripts {
			runSetupItem(scriptName)
		}
	}
}

func installAllTools() {
	// Get tools in dependency order (topological sort)
	tools := config.GetToolsInDependencyOrder()

	for _, tool := range tools {
		result := tool.Check()
		if result.Installed {
			fmt.Printf("  %s %s already installed\n", ui.Green("✓"), tool.Name)
			continue
		}

		// Check if all dependencies are satisfied
		depsMissing := false
		for _, depName := range tool.Dependencies {
			depTool := config.GetToolByName(depName)
			if depTool == nil {
				continue
			}
			depResult := depTool.Check()
			if !depResult.Installed {
				ui.PrintWarning(fmt.Sprintf("Skipping %s (%s not installed)", tool.Name, depName))
				depsMissing = true
				break
			}
		}
		if depsMissing {
			continue
		}

		// Special handling for nvm tools (need manual installation)
		if tool.Method == config.InstallNvm {
			fmt.Printf("  📥 Installing %s (via nvm)...\n", tool.Name)
			ui.PrintWarning(fmt.Sprintf("Run 'nvm install stable' to install %s", tool.Name))
			continue
		}

		fmt.Printf("  📥 Installing %s...\n", tool.Name)
		if err := tool.Install(); err != nil {
			ui.PrintError(fmt.Sprintf("Failed to install %s: %v", tool.Name, err))
		} else {
			// Run post-install scripts
			for _, scriptName := range tool.Scripts {
				runSetupItem(scriptName)
			}
		}
	}

	fmt.Println(ui.Green("✅ All tools installed"))
}
