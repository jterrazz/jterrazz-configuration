package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	updateAll  bool
	updateBrew bool
	updateNpm  bool
	updatePnpm bool
)

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
		// Suggest brew packages for completion
		var suggestions []string
		for _, pkg := range Packages {
			if pkg.Method == InstallBrewFormula || pkg.Method == InstallBrewCask {
				alreadyUsed := false
				for _, arg := range args {
					if arg == pkg.Name {
						alreadyUsed = true
						break
					}
				}
				if !alreadyUsed {
					suggestions = append(suggestions, pkg.Name)
				}
			}
		}
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		// If --all flag, update everything
		if updateAll {
			fmt.Println(cyan("🔄 Updating all packages..."))
			updateBrewPackages()
			updateNpmPackages()
			updatePnpmPackages()
			fmt.Println(green("✅ All updates completed"))
			return
		}

		// If specific manager flags
		if updateBrew || updateNpm || updatePnpm {
			if updateBrew {
				updateBrewPackages()
			}
			if updateNpm {
				updateNpmPackages()
			}
			if updatePnpm {
				updatePnpmPackages()
			}
			fmt.Println(green("✅ Updates completed"))
			return
		}

		// If specific package names provided
		if len(args) > 0 {
			fmt.Println(cyan("🔄 Updating selected packages..."))
			for _, name := range args {
				updatePackageByName(name)
			}
			fmt.Println(green("✅ Updates completed"))
			return
		}

		// No args, list options
		listUpdateOptions()
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&updateAll, "all", "a", false, "Update all package managers")
	updateCmd.Flags().BoolVar(&updateBrew, "brew", false, "Update Homebrew packages")
	updateCmd.Flags().BoolVar(&updateNpm, "npm", false, "Update npm global packages")
	updateCmd.Flags().BoolVar(&updatePnpm, "pnpm", false, "Update pnpm global packages")
	rootCmd.AddCommand(updateCmd)
}

func listUpdateOptions() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	dim := color.New(color.FgHiBlack).SprintFunc()

	fmt.Println(cyan("Available update targets:"))
	fmt.Println()

	managers := []struct {
		name      string
		flag      string
		available bool
	}{
		{"homebrew", "--brew", commandExists("brew")},
		{"npm", "--npm", commandExists("npm")},
		{"pnpm", "--pnpm", commandExists("pnpm")},
	}

	for _, m := range managers {
		status := red("✗")
		if m.available {
			status = green("✓")
		}
		fmt.Printf("  %s %-12s %s\n", status, m.name, dim(m.flag))
	}

	fmt.Println()
	fmt.Println(dim("Usage: j update <package> [package...]"))
	fmt.Println(dim("       j update --brew --npm"))
	fmt.Println(dim("       j update --all"))
}

func updateBrewPackages() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	if !commandExists("brew") {
		printWarning("Homebrew not found, skipping")
		return
	}
	fmt.Println(cyan("🍺 Updating Homebrew packages..."))
	runBrewCommand("update")
	runBrewCommand("upgrade")
	fmt.Println(green("  ✅ Homebrew update completed"))
}

func updateNpmPackages() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	if !commandExists("npm") {
		printWarning("npm not found, skipping")
		return
	}
	fmt.Println(cyan("📦 Updating npm global packages..."))
	runCommand("npm", "update", "-g")
	fmt.Println(green("  ✅ npm update completed"))
}

func updatePnpmPackages() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	if !commandExists("pnpm") {
		printWarning("pnpm not found, skipping")
		return
	}
	fmt.Println(cyan("📦 Updating pnpm global packages..."))
	runCommand("pnpm", "update", "-g")
	fmt.Println(green("  ✅ pnpm update completed"))
}

func updatePackageByName(name string) {
	green := color.New(color.FgGreen).SprintFunc()

	// Find package in our list
	pkg := GetPackageByName(name)
	if pkg != nil {
		switch pkg.Method {
		case InstallBrewFormula:
			if !commandExists("brew") {
				printError("Homebrew not found")
				return
			}
			fmt.Printf("  📥 Updating %s...\n", name)
			runBrewCommand("upgrade", pkg.Formula)
			fmt.Printf("  %s %s updated\n", green("✓"), name)
			return
		case InstallBrewCask:
			if !commandExists("brew") {
				printError("Homebrew not found")
				return
			}
			fmt.Printf("  📥 Updating %s...\n", name)
			runBrewCommand("upgrade", "--cask", pkg.Formula)
			fmt.Printf("  %s %s updated\n", green("✓"), name)
			return
		case InstallNpm:
			if !commandExists("npm") {
				printError("npm not found")
				return
			}
			fmt.Printf("  📥 Updating %s...\n", name)
			runCommand("npm", "update", "-g", pkg.Formula)
			fmt.Printf("  %s %s updated\n", green("✓"), name)
			return
		}
	}

	// Try as a direct brew package name
	if commandExists("brew") {
		fmt.Printf("  📥 Updating %s...\n", name)
		runBrewCommand("upgrade", name)
		fmt.Printf("  %s %s updated\n", green("✓"), name)
		return
	}

	printError(fmt.Sprintf("Unknown package: %s", name))
}
