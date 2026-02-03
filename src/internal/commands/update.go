package commands

import (
	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/domain/tool"
	"github.com/jterrazz/jterrazz-cli/internal/presentation/print"
	"github.com/spf13/cobra"
)

var updateFlags = make(map[string]*bool)

var updateCmd = &cobra.Command{
	Use:   "update [package...]",
	Short: "Update system packages (Homebrew + npm global)",
	Long: `Update system packages.

Examples:
  j update --all             Update all package managers
  j update --brew            Update Homebrew packages only
  j update --npm             Update npm global packages only
  j update --pnpm            Update pnpm global packages only
  j update node              Update specific brew package
  j update claude opencode   Update specific packages
  j update                   List available options`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var all []string
		for _, pkg := range config.Tools {
			if pkg.Method == config.InstallBrewFormula || pkg.Method == config.InstallBrewCask {
				all = append(all, pkg.Name)
			}
		}
		return tool.FilterStrings(all, args), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Check for --all flag
		allFlag, _ := cmd.Flags().GetBool("all")
		if allFlag {
			print.Action("🔄", "Updating all packages...")
			config.UpdateAll()
			print.Done("All updates completed")
			return
		}

		// Check for specific manager flags
		anyFlagSet := false
		for _, pm := range config.PackageManagers {
			if flagVal, ok := updateFlags[pm.Flag]; ok && *flagVal {
				anyFlagSet = true
				config.UpdatePackageManager(pm)
			}
		}
		if anyFlagSet {
			print.Done("Updates completed")
			return
		}

		// If specific package names provided
		if len(args) > 0 {
			print.Action("🔄", "Updating selected packages...")
			for _, name := range args {
				if err := config.UpdatePackageByName(name); err != nil {
					print.Error(err.Error())
				}
			}
			print.Done("Updates completed")
			return
		}

		// No args, list options
		listUpdateOptions()
	},
}

func init() {
	updateCmd.Flags().BoolP("all", "a", false, "Update all package managers")

	// Dynamically add flags for each package manager
	for _, pm := range config.PackageManagers {
		flagPtr := new(bool)
		updateFlags[pm.Flag] = flagPtr
		updateCmd.Flags().BoolVar(flagPtr, pm.Flag, false, "Update "+pm.Name+" packages")
	}

	rootCmd.AddCommand(updateCmd)
}

func listUpdateOptions() {
	print.Info("Available update targets:")
	print.Empty()

	for _, pm := range config.PackageManagers {
		available := config.CommandExists(pm.RequiresCmd)
		print.Row(available, pm.Name, "--"+pm.Flag)
	}

	print.Empty()
	print.Usage(
		"Usage: j update <package> [package...]",
		"       j update --brew --npm",
		"       j update --all",
	)
}
