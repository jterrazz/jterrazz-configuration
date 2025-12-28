package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var installAll bool

var installCmd = &cobra.Command{
	Use:   "install [package...]",
	Short: "Install development packages",
	Long: `Install development packages.

Examples:
  j install --all           Install all packages
  j install brew            Install Homebrew
  j install nvm             Install NVM
  j install go python node  Install specific packages
  j install                  List available packages`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var suggestions []string
		for _, pkg := range Packages {
			// Don't suggest already-specified packages
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
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		// If --all flag, install everything
		if installAll {
			fmt.Println(cyan("🚀 Installing all development tools..."))
			installHomebrew()
			installAllPackages()
			return
		}

		// If no args, list available packages
		if len(args) == 0 {
			listAvailablePackages()
			return
		}

		// Install specific packages
		fmt.Println(cyan("📦 Installing selected packages..."))
		for _, name := range args {
			installPackageByName(name)
		}
		fmt.Println(green("✅ Done"))
	},
}

func init() {
	installCmd.Flags().BoolVarP(&installAll, "all", "a", false, "Install all packages")
	rootCmd.AddCommand(installCmd)
}

func listAvailablePackages() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	dim := color.New(color.FgHiBlack).SprintFunc()

	fmt.Println(cyan("Available packages:"))
	fmt.Println()

	currentCategory := PackageCategory("")
	for _, pkg := range Packages {
		if pkg.Category != currentCategory {
			currentCategory = pkg.Category
			fmt.Printf("%s\n", dim(string(currentCategory)))
		}

		installed, _, _ := CheckPackage(pkg)
		status := red("✗")
		if installed {
			status = green("✓")
		}

		fmt.Printf("  %s %-14s %s\n", status, pkg.Name, dim(pkg.Method.String()))
	}

	fmt.Println()
	fmt.Println(dim("Usage: j install <package> [package...]"))
	fmt.Println(dim("       j install --all"))
}

func installPackageByName(name string) {
	green := color.New(color.FgGreen).SprintFunc()

	// Special case for brew
	if name == "brew" || name == "homebrew" {
		installHomebrew()
		return
	}

	// Find package
	pkg := GetPackageByName(name)
	if pkg == nil {
		printError(fmt.Sprintf("Unknown package: %s", name))
		return
	}

	installed, _, _ := CheckPackage(*pkg)
	if installed {
		fmt.Printf("  %s %s already installed\n", green("✓"), pkg.Name)
		return
	}

	// Check dependencies using the Dependencies field
	for _, depName := range pkg.Dependencies {
		depPkg := GetPackageByName(depName)
		if depPkg == nil {
			continue
		}
		depInstalled, _, _ := CheckPackage(*depPkg)
		if !depInstalled {
			printError(fmt.Sprintf("%s required for %s. Run: j install %s", depName, pkg.Name, depName))
			return
		}
	}

	fmt.Printf("  📥 Installing %s...\n", pkg.Name)
	if err := InstallPackage(*pkg); err != nil {
		printError(fmt.Sprintf("Failed to install %s: %v", pkg.Name, err))
	} else {
		fmt.Printf("  %s %s installed\n", green("✓"), pkg.Name)
	}
}

func installHomebrew() {
	green := color.New(color.FgGreen).SprintFunc()

	for _, pkg := range Packages {
		if pkg.Name == "homebrew" {
			installed, _, _ := CheckPackage(pkg)
			if installed {
				fmt.Printf("  %s Homebrew already installed\n", green("✓"))
				return
			}
			if pkg.InstallFn != nil {
				fmt.Println("  📥 Installing Homebrew...")
				if err := pkg.InstallFn(); err != nil {
					printError(fmt.Sprintf("Failed to install Homebrew: %v", err))
				}
			}
			return
		}
	}
}

func installAllPackages() {
	green := color.New(color.FgGreen).SprintFunc()

	// Get packages in dependency order (topological sort)
	packages := GetPackagesInDependencyOrder()

	for _, pkg := range packages {
		// Skip homebrew (handled separately before this function)
		if pkg.Name == "homebrew" {
			continue
		}

		installed, _, _ := CheckPackage(pkg)
		if installed {
			fmt.Printf("  %s %s already installed\n", green("✓"), pkg.Name)
			continue
		}

		// Check if all dependencies are satisfied
		depsMissing := false
		for _, depName := range pkg.Dependencies {
			depPkg := GetPackageByName(depName)
			if depPkg == nil {
				continue
			}
			depInstalled, _, _ := CheckPackage(*depPkg)
			if !depInstalled {
				printWarning(fmt.Sprintf("Skipping %s (%s not installed)", pkg.Name, depName))
				depsMissing = true
				break
			}
		}
		if depsMissing {
			continue
		}

		// Special handling for nvm packages (need manual installation)
		if pkg.Method == InstallNvm {
			fmt.Printf("  📥 Installing %s (via nvm)...\n", pkg.Name)
			printWarning(fmt.Sprintf("Run 'nvm install stable' to install %s", pkg.Name))
			continue
		}

		fmt.Printf("  📥 Installing %s...\n", pkg.Name)
		if err := InstallPackage(pkg); err != nil {
			printError(fmt.Sprintf("Failed to install %s: %v", pkg.Name, err))
		}
	}

	fmt.Println(green("✅ All packages installed"))
}
