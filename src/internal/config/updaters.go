package config

import (
	"fmt"

	"github.com/jterrazz/jterrazz-cli/internal/ui"
)

// PackageManager represents an updatable package manager
type PackageManager struct {
	Name        string
	Flag        string // CLI flag name (e.g., "brew" for --brew)
	RequiresCmd string // Command that must exist
	UpdateFn    func() // Function to run updates
}

// PackageManagers is the list of all package managers that can be updated
var PackageManagers = []PackageManager{
	{
		Name:        "homebrew",
		Flag:        "brew",
		RequiresCmd: "brew",
		UpdateFn:    updateBrew,
	},
	{
		Name:        "npm",
		Flag:        "npm",
		RequiresCmd: "npm",
		UpdateFn:    updateNpm,
	},
	{
		Name:        "pnpm",
		Flag:        "pnpm",
		RequiresCmd: "pnpm",
		UpdateFn:    updatePnpm,
	},
}

// GetPackageManagerByFlag returns a package manager by its flag name
func GetPackageManagerByFlag(flag string) *PackageManager {
	for i := range PackageManagers {
		if PackageManagers[i].Flag == flag {
			return &PackageManagers[i]
		}
	}
	return nil
}

// UpdateAll updates all available package managers
func UpdateAll() {
	for _, pm := range PackageManagers {
		if CommandExists(pm.RequiresCmd) {
			pm.UpdateFn()
		}
	}
}

// UpdatePackageManager updates a specific package manager
func UpdatePackageManager(pm PackageManager) {
	if !CommandExists(pm.RequiresCmd) {
		fmt.Printf("%s %s not found, skipping\n", ui.Yellow("Warning:"), pm.RequiresCmd)
		return
	}
	pm.UpdateFn()
}

// UpdatePackageByName updates a specific package by name
func UpdatePackageByName(name string) error {
	// Find package in our tools list
	pkg := GetToolByName(name)
	if pkg != nil {
		switch pkg.Method {
		case InstallBrewFormula:
			if !CommandExists("brew") {
				return fmt.Errorf("Homebrew not found")
			}
			fmt.Printf("  📥 Updating %s...\n", name)
			ExecCommand("brew", "upgrade", pkg.Formula)
			fmt.Printf("  %s %s updated\n", ui.Green("✓"), name)
			return nil
		case InstallBrewCask:
			if !CommandExists("brew") {
				return fmt.Errorf("Homebrew not found")
			}
			fmt.Printf("  📥 Updating %s...\n", name)
			ExecCommand("brew", "upgrade", "--cask", pkg.Formula)
			fmt.Printf("  %s %s updated\n", ui.Green("✓"), name)
			return nil
		case InstallNpm:
			if !CommandExists("npm") {
				return fmt.Errorf("npm not found")
			}
			fmt.Printf("  📥 Updating %s...\n", name)
			ExecCommand("npm", "update", "-g", pkg.Formula)
			fmt.Printf("  %s %s updated\n", ui.Green("✓"), name)
			return nil
		}
	}

	// Try as a direct brew package name
	if CommandExists("brew") {
		fmt.Printf("  📥 Updating %s...\n", name)
		ExecCommand("brew", "upgrade", name)
		fmt.Printf("  %s %s updated\n", ui.Green("✓"), name)
		return nil
	}

	return fmt.Errorf("unknown package: %s", name)
}

// =============================================================================
// Update Functions
// =============================================================================

func updateBrew() {
	fmt.Println(ui.Cyan("🍺 Updating Homebrew packages..."))
	ExecCommand("brew", "update")
	ExecCommand("brew", "upgrade")
	fmt.Println(ui.Green("  ✅ Homebrew update completed"))
}

func updateNpm() {
	fmt.Println(ui.Cyan("📦 Updating npm global packages..."))
	ExecCommand("npm", "update", "-g")
	fmt.Println(ui.Green("  ✅ npm update completed"))
}

func updatePnpm() {
	fmt.Println(ui.Cyan("📦 Updating pnpm global packages..."))
	ExecCommand("pnpm", "update", "-g")
	fmt.Println(ui.Green("  ✅ pnpm update completed"))
}
