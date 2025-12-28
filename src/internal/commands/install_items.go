package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// InstallMethod defines how a package is installed
type InstallMethod string

const (
	InstallBrewFormula InstallMethod = "brew"
	InstallBrewCask    InstallMethod = "cask"
	InstallNpm         InstallMethod = "npm"
	InstallNvm         InstallMethod = "nvm"
	InstallXcode       InstallMethod = "xcode"
	InstallManual      InstallMethod = "manual"
)

// PackageCategory groups packages by their purpose
type PackageCategory string

const (
	CategoryPackageManager PackageCategory = "Package Managers"
	CategoryDevelopment    PackageCategory = "Development"
	CategoryInfrastructure PackageCategory = "Infrastructure"
	CategoryAI             PackageCategory = "AI Tools"
	CategorySystemTools    PackageCategory = "System Tools"
)

// Package represents an installable package
type Package struct {
	Name          string
	Command       string          // Command to check if installed
	Formula       string          // Brew formula or npm package name
	Method        InstallMethod   // How to install
	Category      PackageCategory // Which category it belongs to
	Dependencies  []string        // Package names this depends on
	VersionArgs   []string        // Args to get version
	VersionParser func(string) string
	CheckFn       func() (installed bool, version string, extra string) // Custom check function
	InstallFn     func() error                                          // Custom install function
}

// Packages is the single source of truth for all installable packages
var Packages = []Package{
	// Package Managers
	{
		Name:     "homebrew",
		Command:  "brew",
		Method:   InstallManual,
		Category: CategoryPackageManager,
		CheckFn: func() (bool, string, string) {
			if _, err := exec.LookPath("brew"); err != nil {
				return false, "", ""
			}
			out, _ := exec.Command("brew", "--version").Output()
			version := parseBrewVersion(string(out))
			// Get package counts
			formulaeOut, _ := exec.Command("brew", "list", "--formula", "-1").Output()
			caskOut, _ := exec.Command("brew", "list", "--cask", "-1").Output()
			formulaeCount := 0
			caskCount := 0
			if len(strings.TrimSpace(string(formulaeOut))) > 0 {
				formulaeCount = len(strings.Split(strings.TrimSpace(string(formulaeOut)), "\n"))
			}
			if len(strings.TrimSpace(string(caskOut))) > 0 {
				caskCount = len(strings.Split(strings.TrimSpace(string(caskOut)), "\n"))
			}
			return true, version, fmt.Sprintf("%d formulae, %d casks", formulaeCount, caskCount)
		},
		InstallFn: func() error {
			cmd := exec.Command("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			return cmd.Run()
		},
	},
	{
		Name:          "bun",
		Command:       "bun",
		Formula:       "bun",
		Method:        InstallBrewFormula,
		Category:      CategoryPackageManager,
		Dependencies:  []string{"homebrew"},
		VersionArgs:   []string{"--version"},
		VersionParser: trimVersion,
	},
	{
		Name:         "npm",
		Command:      "npm",
		Method:       InstallNvm,
		Category:     CategoryPackageManager,
		Dependencies: []string{"node"},
		CheckFn: func() (bool, string, string) {
			if _, err := exec.LookPath("npm"); err != nil {
				return false, "", ""
			}
			out, _ := exec.Command("npm", "--version").Output()
			version := trimVersion(string(out))
			// Get global package count
			npmOut, _ := exec.Command("npm", "list", "-g", "--depth=0", "--parseable").Output()
			npmLines := strings.Split(strings.TrimSpace(string(npmOut)), "\n")
			count := len(npmLines) - 1
			if count < 0 {
				count = 0
			}
			return true, version, fmt.Sprintf("%d global", count)
		},
	},
	{
		Name:         "nvm",
		Command:      "",
		Formula:      "nvm",
		Method:       InstallBrewFormula,
		Category:     CategoryPackageManager,
		Dependencies: []string{"homebrew"},
		CheckFn: func() (bool, string, string) {
			nvmDir := os.Getenv("HOME") + "/.nvm"
			if _, err := os.Stat(nvmDir); err != nil {
				return false, "", ""
			}
			// Count installed node versions
			versionsDir := nvmDir + "/versions/node"
			entries, err := os.ReadDir(versionsDir)
			extra := ""
			if err == nil {
				count := 0
				for _, e := range entries {
					if e.IsDir() && strings.HasPrefix(e.Name(), "v") {
						count++
					}
				}
				if count > 0 {
					extra = fmt.Sprintf("%d versions", count)
				}
			}
			return true, "", extra
		},
	},
	{
		Name:          "pnpm",
		Command:       "pnpm",
		Formula:       "pnpm",
		Method:        InstallBrewFormula,
		Category:      CategoryPackageManager,
		Dependencies:  []string{"homebrew"},
		VersionArgs:   []string{"--version"},
		VersionParser: trimVersion,
	},

	// Development
	{
		Name:         "docker",
		Command:      "docker",
		Formula:      "docker",
		Method:       InstallBrewCask,
		Category:     CategoryDevelopment,
		Dependencies: []string{"homebrew"},
		CheckFn: func() (bool, string, string) {
			// Check if Docker Desktop app is installed
			_, appErr := os.Stat("/Applications/Docker.app")
			if appErr != nil {
				return false, "", ""
			}
			// Get version from CLI if available
			version := ""
			if _, err := exec.LookPath("docker"); err == nil {
				out, _ := exec.Command("docker", "--version").Output()
				parts := strings.Split(string(out), " ")
				if len(parts) >= 3 {
					version = strings.TrimSuffix(parts[2], ",")
				}
			}
			extra := "stopped"
			if err := exec.Command("docker", "info").Run(); err == nil {
				extra = "running"
			}
			return true, version, extra
		},
	},
	{
		Name:          "git",
		Command:       "git",
		Method:        InstallXcode,
		Category:      CategoryDevelopment,
		VersionArgs:   []string{"--version"},
		VersionParser: parseGitVersion,
	},
	{
		Name:          "go",
		Command:       "go",
		Formula:       "go",
		Method:        InstallBrewFormula,
		Category:      CategoryDevelopment,
		Dependencies:  []string{"homebrew"},
		VersionArgs:   []string{"version"},
		VersionParser: parseGoVersion,
	},
	{
		Name:          "node",
		Command:       "node",
		Method:        InstallNvm,
		Category:      CategoryDevelopment,
		Dependencies:  []string{"nvm"},
		VersionArgs:   []string{"--version"},
		VersionParser: trimVersion,
	},
	{
		Name:         "openjdk",
		Command:      "java",
		Formula:      "openjdk",
		Method:       InstallBrewFormula,
		Category:     CategoryDevelopment,
		Dependencies: []string{"homebrew"},
		CheckFn: func() (bool, string, string) {
			// Check for brew-installed openjdk first
			brewJava := "/opt/homebrew/opt/openjdk/bin/java"
			if _, err := os.Stat(brewJava); err == nil {
				out, _ := exec.Command(brewJava, "-version").CombinedOutput()
				return true, parseJavaVersion(string(out)), ""
			}
			// Fallback: check if system java_home finds a JDK
			cmd := exec.Command("/usr/libexec/java_home")
			if err := cmd.Run(); err != nil {
				return false, "", ""
			}
			out, _ := exec.Command("java", "-version").CombinedOutput()
			return true, parseJavaVersion(string(out)), ""
		},
	},
	{
		Name:          "python",
		Command:       "python3",
		Formula:       "python",
		Method:        InstallBrewFormula,
		Category:      CategoryDevelopment,
		Dependencies:  []string{"homebrew"},
		VersionArgs:   []string{"--version"},
		VersionParser: parsePythonVersion,
	},

	// Infrastructure
	{
		Name:          "ansible",
		Command:       "ansible",
		Formula:       "ansible",
		Method:        InstallBrewFormula,
		Category:      CategoryInfrastructure,
		Dependencies:  []string{"homebrew"},
		VersionArgs:   []string{"--version"},
		VersionParser: parseAnsibleVersion,
	},
	{
		Name:          "ansible-lint",
		Command:       "ansible-lint",
		Formula:       "ansible-lint",
		Method:        InstallBrewFormula,
		Category:      CategoryInfrastructure,
		Dependencies:  []string{"homebrew"},
		VersionArgs:   []string{"--version"},
		VersionParser: parseAnsibleLintVersion,
	},
	{
		Name:          "kubectl",
		Command:       "kubectl",
		Formula:       "kubectl",
		Method:        InstallBrewFormula,
		Category:      CategoryInfrastructure,
		Dependencies:  []string{"homebrew"},
		VersionArgs:   []string{"version", "--client", "-o", "yaml"},
		VersionParser: parseKubectlVersion,
	},
	{
		Name:          "multipass",
		Command:       "multipass",
		Formula:       "multipass",
		Method:        InstallBrewFormula,
		Category:      CategoryInfrastructure,
		Dependencies:  []string{"homebrew"},
		VersionArgs:   []string{"--version"},
		VersionParser: parseMultipassVersion,
	},
	{
		Name:          "terraform",
		Command:       "terraform",
		Formula:       "terraform",
		Method:        InstallBrewFormula,
		Category:      CategoryInfrastructure,
		Dependencies:  []string{"homebrew"},
		VersionArgs:   []string{"--version"},
		VersionParser: parseTerraformVersion,
	},

	// AI Tools
	{
		Name:          "claude",
		Command:       "claude",
		Formula:       "@anthropic-ai/claude-code",
		Method:        InstallNpm,
		Category:      CategoryAI,
		Dependencies:  []string{"npm"},
		VersionArgs:   []string{"--version"},
		VersionParser: parseClaudeVersion,
	},
	{
		Name:          "codex",
		Command:       "codex",
		Formula:       "@openai/codex",
		Method:        InstallNpm,
		Category:      CategoryAI,
		Dependencies:  []string{"npm"},
		VersionArgs:   []string{"--version"},
		VersionParser: parseCodexVersion,
	},
	{
		Name:          "gemini",
		Command:       "gemini",
		Formula:       "gemini-cli",
		Method:        InstallBrewFormula,
		Category:      CategoryAI,
		Dependencies:  []string{"homebrew"},
		VersionArgs:   []string{"--version"},
		VersionParser: trimVersion,
	},

	// System Tools
	{
		Name:          "mole",
		Command:       "mo",
		Formula:       "tw93/tap/mole",
		Method:        InstallBrewFormula,
		Category:      CategorySystemTools,
		Dependencies:  []string{"homebrew"},
		VersionArgs:   []string{"--version"},
		VersionParser: parseMoleVersion,
	},
	{
		Name:         "neohtop",
		Command:      "",
		Formula:      "neohtop",
		Method:       InstallBrewCask,
		Category:     CategorySystemTools,
		Dependencies: []string{"homebrew"},
	},
}

// GetPackagesByCategory returns packages filtered by category
func GetPackagesByCategory(category PackageCategory) []Package {
	var result []Package
	for _, pkg := range Packages {
		if pkg.Category == category {
			result = append(result, pkg)
		}
	}
	return result
}

// GetInstallablePackages returns packages that can be installed via brew or npm
func GetInstallablePackages() []Package {
	var result []Package
	for _, pkg := range Packages {
		if pkg.Method == InstallBrewFormula || pkg.Method == InstallBrewCask || pkg.Method == InstallNpm || pkg.InstallFn != nil {
			result = append(result, pkg)
		}
	}
	return result
}

// GetPackageByName returns a package by name
func GetPackageByName(name string) *Package {
	for i := range Packages {
		if Packages[i].Name == name {
			return &Packages[i]
		}
	}
	return nil
}

// GetPackagesInDependencyOrder returns all installable packages sorted by dependencies
func GetPackagesInDependencyOrder() []Package {
	installable := GetInstallablePackages()

	// Build a map for quick lookup
	pkgMap := make(map[string]*Package)
	for i := range installable {
		pkgMap[installable[i].Name] = &installable[i]
	}

	// Track visited and result
	visited := make(map[string]bool)
	var result []Package

	// Recursive function to add package and its dependencies
	var visit func(name string)
	visit = func(name string) {
		if visited[name] {
			return
		}

		pkg := pkgMap[name]
		if pkg == nil {
			// Check if it's in the full package list (might not be installable but needed as dep)
			pkg = GetPackageByName(name)
		}
		if pkg == nil {
			return
		}

		// Visit dependencies first
		for _, dep := range pkg.Dependencies {
			visit(dep)
		}

		visited[name] = true

		// Only add if it's installable
		if pkgMap[name] != nil {
			result = append(result, *pkg)
		}
	}

	// Visit all installable packages
	for _, pkg := range installable {
		visit(pkg.Name)
	}

	return result
}

// CheckPackage checks if a package is installed and returns its status
func CheckPackage(pkg Package) (installed bool, version string, extra string) {
	// Use custom check function if provided
	if pkg.CheckFn != nil {
		return pkg.CheckFn()
	}

	// Default check using command
	if pkg.Command == "" {
		return false, "", ""
	}

	if _, err := exec.LookPath(pkg.Command); err != nil {
		return false, "", ""
	}

	// Get version if version args provided
	if len(pkg.VersionArgs) > 0 && pkg.VersionParser != nil {
		out, err := exec.Command(pkg.Command, pkg.VersionArgs...).CombinedOutput()
		if err == nil {
			version = pkg.VersionParser(string(out))
		}
	}

	return true, version, ""
}

// InstallPackage installs a package
func InstallPackage(pkg Package) error {
	// Use custom install function if provided
	if pkg.InstallFn != nil {
		return pkg.InstallFn()
	}

	switch pkg.Method {
	case InstallBrewFormula:
		return runCommand("brew", "install", pkg.Formula)
	case InstallBrewCask:
		return runCommand("brew", "install", "--cask", pkg.Formula)
	case InstallNpm:
		return runCommand("npm", "install", "-g", pkg.Formula)
	default:
		return fmt.Errorf("cannot auto-install %s (method: %s)", pkg.Name, pkg.Method)
	}
}

// MethodString returns a display string for the install method
func (m InstallMethod) String() string {
	switch m {
	case InstallBrewFormula:
		return "brew"
	case InstallBrewCask:
		return "cask"
	case InstallNpm:
		return "npm"
	case InstallNvm:
		return "nvm"
	case InstallXcode:
		return "xcode"
	default:
		return "-"
	}
}
