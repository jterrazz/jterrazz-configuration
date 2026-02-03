package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jterrazz/jterrazz-cli/internal/domain/tool"
)

// CheckResult is the unified result type for all check operations
type CheckResult struct {
	Installed bool   // Whether the item is installed/configured
	Version   string // Version string (if applicable)
	Status    string // Additional status: "running", "stopped", "3 versions", etc.
	Detail    string // Extra detail: path, config location, etc.
}

// CheckResult constructors for common patterns

// Installed creates a CheckResult for an installed item
func Installed() CheckResult {
	return CheckResult{Installed: true}
}

// InstalledWithVersion creates a CheckResult for an installed item with version
func InstalledWithVersion(version string) CheckResult {
	return CheckResult{Installed: true, Version: version}
}

// InstalledWithDetail creates a CheckResult for an installed item with detail
func InstalledWithDetail(detail string) CheckResult {
	return CheckResult{Installed: true, Detail: detail}
}

// InstalledWithStatus creates a CheckResult for an installed item with status
func InstalledWithStatus(version, status string) CheckResult {
	return CheckResult{Installed: true, Version: version, Status: status}
}

// NotInstalled creates a CheckResult for a not installed item
func NotInstalled() CheckResult {
	return CheckResult{}
}

// ToolCategory groups tools by their purpose
type ToolCategory string

const (
	CategoryPackageManager ToolCategory = "Package Managers"
	CategoryLanguages      ToolCategory = "Languages"
	CategoryInfrastructure ToolCategory = "Infrastructure"
	CategoryAI             ToolCategory = "AI"
	CategoryApps           ToolCategory = "Apps"
	CategorySystemTools    ToolCategory = "System Tools"
)

// InstallMethod defines how a tool is installed
type InstallMethod string

const (
	InstallBrewFormula InstallMethod = "brew"
	InstallBrewCask    InstallMethod = "cask"
	InstallNpm         InstallMethod = "npm"
	InstallNvm         InstallMethod = "nvm"
	InstallXcode       InstallMethod = "xcode"
	InstallManual      InstallMethod = "manual"
)

// String returns a display string for the install method
func (m InstallMethod) String() string {
	switch m {
	case InstallBrewFormula, InstallBrewCask:
		return "brew"
	case InstallNpm:
		return "npm"
	case InstallNvm:
		return "nvm"
	case InstallXcode:
		return "xcode"
	case InstallManual:
		return "sh"
	default:
		return "-"
	}
}

// Tool represents an installable piece of software
type Tool struct {
	Name        string
	Description string
	Category    ToolCategory

	// Check - how to verify if installed
	Command string             // CLI command to check existence
	CheckFn func() CheckResult // Custom check (overrides Command)

	// Install - how to install
	Method       InstallMethod // brew, npm, manual, etc.
	Formula      string        // Brew formula or npm package name
	InstallFn    func() error  // Custom install (overrides Method)
	Dependencies []string      // Tool names this depends on

	// Version - how to get version info
	VersionFn func() string // Returns version string

	// Scripts - post-install or related scripts
	Scripts []string // Script names to run after install
}

// Tools is the single source of truth for all installable software
var Tools = []Tool{
	// ==========================================================================
	// Package Managers
	// ==========================================================================
	{
		Name:         "bun",
		Command:      "bun",
		Formula:      "bun",
		Method:       InstallBrewFormula,
		Category:     CategoryPackageManager,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("bun", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:         "cocoapods",
		Command:      "pod",
		Formula:      "cocoapods",
		Method:       InstallBrewFormula,
		Category:     CategoryPackageManager,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("pod", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:     "homebrew",
		Command:  "brew",
		Method:   InstallManual,
		Category: CategoryPackageManager,
		CheckFn: func() CheckResult {
			if _, err := exec.LookPath("brew"); err != nil {
				return CheckResult{}
			}
			out, _ := exec.Command("brew", "--version").Output()
			version := tool.ParseBrewVersion(string(out))
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
			return CheckResult{
				Installed: true,
				Version:   version,
				Status:    fmt.Sprintf("%d formulae, %d casks", formulaeCount, caskCount),
			}
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
		Name:         "npm",
		Command:      "npm",
		Method:       InstallNvm,
		Category:     CategoryPackageManager,
		Dependencies: []string{"node"},
		CheckFn: func() CheckResult {
			if _, err := exec.LookPath("npm"); err != nil {
				return CheckResult{}
			}
			out, _ := exec.Command("npm", "--version").Output()
			version := tool.TrimVersion(string(out))
			npmOut, _ := exec.Command("npm", "list", "-g", "--depth=0", "--parseable").Output()
			npmLines := strings.Split(strings.TrimSpace(string(npmOut)), "\n")
			count := len(npmLines) - 1
			if count < 0 {
				count = 0
			}
			return CheckResult{
				Installed: true,
				Version:   version,
				Status:    fmt.Sprintf("%d global", count),
			}
		},
	},
	{
		Name:         "nvm",
		Command:      "",
		Formula:      "nvm",
		Method:       InstallBrewFormula,
		Category:     CategoryPackageManager,
		Dependencies: []string{"homebrew"},
		CheckFn: func() CheckResult {
			nvmDir := os.Getenv("HOME") + "/.nvm"
			if _, err := os.Stat(nvmDir); err != nil {
				return CheckResult{}
			}
			versionsDir := nvmDir + "/versions/node"
			entries, err := os.ReadDir(versionsDir)
			status := ""
			if err == nil {
				count := 0
				for _, e := range entries {
					if e.IsDir() && strings.HasPrefix(e.Name(), "v") {
						count++
					}
				}
				if count > 0 {
					status = fmt.Sprintf("%d versions", count)
				}
			}
			version := tool.VersionFromBrewFormula("nvm")()
			return CheckResult{Installed: true, Version: version, Status: status}
		},
	},
	{
		Name:         "pnpm",
		Command:      "pnpm",
		Formula:      "pnpm",
		Method:       InstallBrewFormula,
		Category:     CategoryPackageManager,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("pnpm", []string{"--version"}, tool.TrimVersion),
	},

	// ==========================================================================
	// Languages
	// ==========================================================================
	{
		Name:         "go",
		Command:      "go",
		Formula:      "go",
		Method:       InstallBrewFormula,
		Category:     CategoryLanguages,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("go", []string{"version"}, tool.ParseGoVersion),
	},
	{
		Name:         "node",
		Command:      "node",
		Method:       InstallNvm,
		Category:     CategoryLanguages,
		Dependencies: []string{"nvm"},
		VersionFn:    tool.VersionFromCmd("node", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:         "openjdk",
		Command:      "java",
		Formula:      "openjdk",
		Method:       InstallBrewFormula,
		Category:     CategoryLanguages,
		Dependencies: []string{"homebrew"},
		Scripts:      []string{"java"},
		CheckFn: func() CheckResult {
			brewJava := "/opt/homebrew/opt/openjdk/bin/java"
			if _, err := os.Stat(brewJava); err == nil {
				out, _ := exec.Command(brewJava, "-version").CombinedOutput()
				return CheckResult{Installed: true, Version: tool.ParseJavaVersion(string(out))}
			}
			cmd := exec.Command("/usr/libexec/java_home")
			if err := cmd.Run(); err != nil {
				return CheckResult{}
			}
			out, _ := exec.Command("java", "-version").CombinedOutput()
			return CheckResult{Installed: true, Version: tool.ParseJavaVersion(string(out))}
		},
	},
	{
		Name:         "python",
		Command:      "python3",
		Formula:      "python",
		Method:       InstallBrewFormula,
		Category:     CategoryLanguages,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("python3", []string{"--version"}, tool.ParsePythonVersion),
	},

	// ==========================================================================
	// Infrastructure
	// ==========================================================================
	{
		Name:         "ansible",
		Command:      "ansible",
		Formula:      "ansible",
		Method:       InstallBrewFormula,
		Category:     CategoryInfrastructure,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("ansible", []string{"--version"}, tool.ParseAnsibleVersion),
	},

	{
		Name:         "multipass",
		Command:      "multipass",
		Formula:      "multipass",
		Method:       InstallBrewFormula,
		Category:     CategoryInfrastructure,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("multipass", []string{"--version"}, tool.ParseMultipassVersion),
	},
	{
		Name:         "pulumi",
		Command:      "pulumi",
		Formula:      "pulumi/tap/pulumi",
		Method:       InstallBrewFormula,
		Category:     CategoryInfrastructure,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("pulumi", []string{"version"}, tool.ParsePulumiVersion),
	},
	{
		Name:         "terraform",
		Command:      "terraform",
		Formula:      "terraform",
		Method:       InstallBrewFormula,
		Category:     CategoryInfrastructure,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("terraform", []string{"--version"}, tool.ParseTerraformVersion),
	},

	// ==========================================================================
	// AI
	// ==========================================================================
	{
		Name:         "claude",
		Command:      "claude",
		Formula:      "claude-code",
		Method:       InstallBrewCask,
		Category:     CategoryAI,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromBrewCask("claude-code"),
	},
	{
		Name:         "codex",
		Command:      "codex",
		Formula:      "codex",
		Method:       InstallBrewCask,
		Category:     CategoryAI,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("codex", []string{"--version"}, tool.ParseBrewVersion),
	},
	{
		Name:         "gemini",
		Command:      "gemini",
		Formula:      "gemini-cli",
		Method:       InstallBrewFormula,
		Category:     CategoryAI,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("gemini", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:         "ollama",
		Command:      "ollama",
		Formula:      "ollama-app",
		Method:       InstallBrewCask,
		Category:     CategoryAI,
		Dependencies: []string{"homebrew"},
		CheckFn: func() CheckResult {
			_, appErr := os.Stat("/Applications/Ollama.app")
			if appErr != nil {
				return CheckResult{}
			}
			version := tool.VersionFromBrewCask("ollama-app")()
			status := "stopped"
			if err := exec.Command("pgrep", "-x", "ollama").Run(); err == nil {
				status = "running"
			}
			return CheckResult{Installed: true, Version: version, Status: status}
		},
	},
	{
		Name:         "happy-coder",
		Command:      "happy",
		Formula:      "happy-coder",
		Method:       InstallNpm,
		Category:     CategoryAI,
		Dependencies: []string{"npm"},
		VersionFn:    tool.VersionFromCmd("happy", []string{"--version"}, tool.ParseHappyCoderVersion),
	},
	{
		Name:         "opencode",
		Command:      "opencode",
		Formula:      "opencode",
		Method:       InstallBrewFormula,
		Category:     CategoryAI,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("opencode", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:         "qmd",
		Command:      "qmd",
		Method:       InstallManual,
		Category:     CategoryAI,
		Dependencies: []string{"bun"},
		VersionFn:    tool.VersionFromCmd("qmd", []string{"--version"}, tool.TrimVersion),
		CheckFn: func() CheckResult {
			if _, err := exec.LookPath("qmd"); err != nil {
				return CheckResult{}
			}
			version := tool.VersionFromCmd("qmd", []string{"--version"}, tool.TrimVersion)()
			return CheckResult{Installed: true, Version: version}
		},
		InstallFn: func() error {
			cmd := exec.Command("bun", "install", "-g", "https://github.com/tobi/qmd")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			return cmd.Run()
		},
	},
	{
		Name:         "skills",
		Command:      "skills",
		Formula:      "skills",
		Method:       InstallNpm,
		Category:     CategoryAI,
		Dependencies: []string{"npm"},
		VersionFn:    tool.VersionFromCmd("skills", []string{"--version"}, tool.TrimVersion),
	},

	// ==========================================================================
	// Apps
	// ==========================================================================
	{
		Name:         "docker",
		Command:      "docker",
		Formula:      "docker",
		Method:       InstallBrewCask,
		Category:     CategoryApps,
		Dependencies: []string{"homebrew"},
		CheckFn: func() CheckResult {
			_, appErr := os.Stat("/Applications/Docker.app")
			if appErr != nil {
				return CheckResult{}
			}
			version := tool.VersionFromBrewCask("docker")()
			status := "stopped"
			if err := exec.Command("docker", "info").Run(); err == nil {
				status = "running"
			}
			return CheckResult{Installed: true, Version: version, Status: status}
		},
	},
	{
		Name:         "ghostty",
		Formula:      "ghostty",
		Method:       InstallBrewCask,
		Category:     CategoryApps,
		Dependencies: []string{"homebrew"},
		Scripts:      []string{"ghostty"},
		CheckFn: func() CheckResult {
			if _, err := os.Stat("/Applications/Ghostty.app"); err != nil {
				return CheckResult{}
			}
			version := tool.VersionFromBrewCask("ghostty")()
			return CheckResult{Installed: true, Version: version}
		},
	},
	{
		Name:         "gpg",
		Description:  "GNU Privacy Guard for encryption and signing",
		Command:      "gpg",
		Formula:      "gnupg",
		Method:       InstallBrewFormula,
		Category:     CategorySystemTools,
		Dependencies: []string{"homebrew"},
		Scripts:      []string{"gpg"},
		VersionFn:    tool.VersionFromBrewFormula("gnupg"),
	},
	{
		Name:        "ohmyzsh",
		Description: "Oh My Zsh shell framework",
		Command:     "",
		Method:      InstallManual,
		Category:    CategorySystemTools,
		CheckFn: func() CheckResult {
			omzPath := os.Getenv("HOME") + "/.oh-my-zsh"
			if _, err := os.Stat(omzPath); err != nil {
				return CheckResult{}
			}
			// Get git commit hash as version
			cmd := exec.Command("git", "-C", omzPath, "rev-parse", "--short", "HEAD")
			out, err := cmd.Output()
			version := ""
			if err == nil {
				version = strings.TrimSpace(string(out))
			}
			return CheckResult{Installed: true, Version: version}
		},
		InstallFn: func() error {
			cmd := exec.Command("sh", "-c", "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			return cmd.Run()
		},
	},
	{
		Name:         "zed",
		Description:  "Zed code editor",
		Formula:      "zed",
		Method:       InstallBrewCask,
		Category:     CategoryApps,
		Dependencies: []string{"homebrew"},
		Scripts:      []string{"zed"},
		CheckFn: func() CheckResult {
			if _, err := os.Stat("/Applications/Zed.app"); err != nil {
				return CheckResult{}
			}
			version := tool.VersionFromAppPlist("Zed")()
			return CheckResult{Installed: true, Version: version}
		},
	},

	// ==========================================================================
	// System Tools
	// ==========================================================================
	{
		Name:      "git",
		Command:   "git",
		Method:    InstallXcode,
		Category:  CategorySystemTools,
		VersionFn: tool.VersionFromCmd("git", []string{"--version"}, tool.ParseGitVersion),
	},
	{
		Name:         "mole",
		Command:      "mo",
		Formula:      "tw93/tap/mole",
		Method:       InstallBrewFormula,
		Category:     CategorySystemTools,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("mo", []string{"--version"}, tool.ParseMoleVersion),
	},
}

// =============================================================================
// Tool Functions
// =============================================================================

// GetAllTools returns all tools
func GetAllTools() []Tool {
	return Tools
}

// GetToolsByCategory returns tools filtered by category
func GetToolsByCategory(category ToolCategory) []Tool {
	var result []Tool
	for _, tool := range Tools {
		if tool.Category == category {
			result = append(result, tool)
		}
	}
	return result
}

// GetInstallableTools returns tools that can be installed
func GetInstallableTools() []Tool {
	var result []Tool
	for _, tool := range Tools {
		if tool.Method == InstallBrewFormula || tool.Method == InstallBrewCask || tool.Method == InstallNpm || tool.InstallFn != nil {
			result = append(result, tool)
		}
	}
	return result
}

// GetToolByName returns a tool by name
func GetToolByName(name string) *Tool {
	for i := range Tools {
		if Tools[i].Name == name {
			return &Tools[i]
		}
	}
	return nil
}

// GetToolsInDependencyOrder returns all installable tools sorted by dependencies
func GetToolsInDependencyOrder() []Tool {
	installable := GetInstallableTools()

	toolMap := make(map[string]*Tool)
	for i := range installable {
		toolMap[installable[i].Name] = &installable[i]
	}

	visited := make(map[string]bool)
	var result []Tool

	var visit func(name string)
	visit = func(name string) {
		if visited[name] {
			return
		}

		tool := toolMap[name]
		if tool == nil {
			tool = GetToolByName(name)
		}
		if tool == nil {
			return
		}

		for _, dep := range tool.Dependencies {
			visit(dep)
		}

		visited[name] = true

		if toolMap[name] != nil {
			result = append(result, *tool)
		}
	}

	for _, tool := range installable {
		visit(tool.Name)
	}

	return result
}

// Check checks if a tool is installed and returns its status
func (t Tool) Check() CheckResult {
	if t.CheckFn != nil {
		return t.CheckFn()
	}

	if t.Command == "" {
		return CheckResult{}
	}

	if _, err := exec.LookPath(t.Command); err != nil {
		return CheckResult{}
	}

	result := CheckResult{Installed: true}

	if t.VersionFn != nil {
		result.Version = t.VersionFn()
	}

	return result
}

// Install installs the tool
func (t Tool) Install() error {
	if t.InstallFn != nil {
		return t.InstallFn()
	}

	switch t.Method {
	case InstallBrewFormula:
		return RunBrewCommand("install", t.Formula)
	case InstallBrewCask:
		return RunBrewCommand("install", "--cask", t.Formula)
	case InstallNpm:
		return ExecCommand("npm", "install", "-g", t.Formula)
	default:
		return fmt.Errorf("cannot auto-install %s (method: %s)", t.Name, t.Method)
	}
}

// RunBrewCommand runs a brew command with ARM architecture forced
func RunBrewCommand(args ...string) error {
	cmd := exec.Command("arch", append([]string{"-arm64", "brew"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
