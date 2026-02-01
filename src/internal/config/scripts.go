package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jterrazz/jterrazz-cli/internal/tool"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
)

// ScriptCategory groups scripts by their purpose
type ScriptCategory string

const (
	ScriptCategoryTerminal ScriptCategory = "Terminal"
	ScriptCategorySecurity ScriptCategory = "Security"
	ScriptCategoryEditor   ScriptCategory = "Editor"
	ScriptCategorySystem   ScriptCategory = "System"
)

// Script represents a setup/configuration task
// Scripts can be standalone or attached to a Tool via Tool.Scripts
type Script struct {
	Name        string
	Description string
	Category    ScriptCategory

	// Check - verify if already configured (optional)
	// If nil, script is "run-once" with no checkable state
	CheckFn func() CheckResult

	// Run - execute the script
	RunFn func() error

	// Dependencies
	RequiresTool string // Tool that must be installed first (e.g., "openjdk")
}

// Scripts is the single source of truth for all setup/configuration scripts
var Scripts = []Script{
	// ==========================================================================
	// Terminal Scripts
	// ==========================================================================
	{
		Name:        "hushlogin",
		Description: "Silence terminal login message",
		Category:    ScriptCategoryTerminal,
		CheckFn: func() CheckResult {
			hushPath := os.Getenv("HOME") + "/.hushlogin"
			if _, err := os.Stat(hushPath); err == nil {
				return CheckResult{Installed: true, Detail: "~/.hushlogin"}
			}
			return CheckResult{}
		},
		RunFn: runHushlogin,
	},
	{
		Name:         "ghostty-config",
		Description:  "Install Ghostty terminal config",
		Category:     ScriptCategoryTerminal,
		RequiresTool: "ghostty",
		CheckFn: func() CheckResult {
			configPath := os.Getenv("HOME") + "/.config/ghostty/config"
			if _, err := os.Stat(configPath); err == nil {
				return CheckResult{Installed: true, Detail: "~/.config/ghostty/config"}
			}
			return CheckResult{}
		},
		RunFn: runGhosttyConfig,
	},

	// ==========================================================================
	// Security Scripts
	// ==========================================================================
	{
		Name:         "gpg-setup",
		Description:  "Configure GPG for commit signing",
		Category:     ScriptCategorySecurity,
		RequiresTool: "gpg",
		CheckFn: func() CheckResult {
			out, _ := exec.Command("git", "config", "--global", "commit.gpgsign").Output()
			if strings.TrimSpace(string(out)) == "true" {
				return CheckResult{Installed: true, Detail: "commit.gpgsign=true"}
			}
			return CheckResult{}
		},
		RunFn: runGPGSetup,
	},
	{
		Name:        "ssh",
		Description: "Generate SSH key with Keychain integration",
		Category:    ScriptCategorySecurity,
		CheckFn: func() CheckResult {
			sshKey := os.Getenv("HOME") + "/.ssh/id_ed25519"
			if _, err := os.Stat(sshKey); err == nil {
				return CheckResult{Installed: true, Detail: "~/.ssh/id_ed25519"}
			}
			return CheckResult{}
		},
		RunFn: runSSHSetup,
	},

	// ==========================================================================
	// Editor Scripts
	// ==========================================================================
	{
		Name:         "zed-config",
		Description:  "Install Zed editor config",
		Category:     ScriptCategoryEditor,
		RequiresTool: "zed",
		CheckFn: func() CheckResult {
			configPath := os.Getenv("HOME") + "/.config/zed/settings.json"
			if _, err := os.Stat(configPath); err == nil {
				return CheckResult{Installed: true, Detail: "~/.config/zed/settings.json"}
			}
			return CheckResult{}
		},
		RunFn: runZedConfig,
	},

	// ==========================================================================
	// System Scripts
	// ==========================================================================
	{
		Name:         "java-symlink",
		Description:  "Configure Java runtime symlink for macOS",
		Category:     ScriptCategorySystem,
		RequiresTool: "openjdk",
		CheckFn: func() CheckResult {
			if _, err := os.Lstat("/Library/Java/JavaVirtualMachines/openjdk.jdk"); err == nil {
				return CheckResult{Installed: true, Detail: "/Library/Java/JavaVirtualMachines/openjdk.jdk"}
			}
			return CheckResult{}
		},
		RunFn: runJavaSymlink,
	},
	{
		Name:        "dock-reset",
		Description: "Reset dock to system defaults",
		Category:    ScriptCategorySystem,
		RunFn:       runDockReset,
		// No CheckFn - this is a run-once action with no state
	},
	{
		Name:        "dock-spacer",
		Description: "Add a small spacer tile to the dock",
		Category:    ScriptCategorySystem,
		RunFn:       runDockSpacer,
		// No CheckFn - can be run multiple times
	},
}

// =============================================================================
// Script Runners
// =============================================================================

func runHushlogin() error {
	fmt.Println(ui.Cyan("Setting up hushlogin..."))

	hushPath := os.Getenv("HOME") + "/.hushlogin"
	if _, err := os.Stat(hushPath); err == nil {
		fmt.Printf("%s .hushlogin already exists\n", ui.Green("Done"))
		return nil
	}

	f, err := os.Create(hushPath)
	if err != nil {
		return fmt.Errorf("failed to create .hushlogin: %w", err)
	}
	f.Close()

	fmt.Println(ui.Green("Done - terminal login message silenced"))
	return nil
}

func runGhosttyConfig() error {
	fmt.Println(ui.Cyan("Setting up Ghostty config..."))

	configDir := os.Getenv("HOME") + "/.config/ghostty"
	configPath := configDir + "/config"

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	repoConfig, err := GetRepoConfigPath("configuration/applications/ghostty/config")
	if err != nil {
		return fmt.Errorf("failed to find repo config: %w", err)
	}

	configContent, err := os.ReadFile(repoConfig)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", repoConfig, err)
	}

	if err := os.WriteFile(configPath, configContent, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", configPath, err)
	}

	fmt.Println(ui.Green("Done - Ghostty config installed"))
	return nil
}

func runGPGSetup() error {
	fmt.Println(ui.Cyan("Setting up GPG for commit signing..."))

	email := UserEmail
	name := UserName

	if !CommandExists("gpg") {
		return fmt.Errorf("GPG not installed. Run: brew install gnupg")
	}

	checkCmd := exec.Command("gpg", "--list-secret-keys", "--keyid-format", "long", email)
	if output, err := checkCmd.Output(); err == nil && len(output) > 0 {
		fmt.Println(ui.Green("GPG key already exists for " + email))
		configureGitGPG(email)
		return nil
	}

	fmt.Println("Generating GPG key...")
	fmt.Println(ui.Dim("Using ed25519 algorithm"))

	batchConfig := fmt.Sprintf(`%%no-protection
Key-Type: eddsa
Key-Curve: ed25519
Name-Real: %s
Name-Email: %s
Expire-Date: 0
%%commit
`, name, email)

	genCmd := exec.Command("gpg", "--batch", "--generate-key")
	genCmd.Stdin = strings.NewReader(batchConfig)
	genCmd.Stdout = os.Stdout
	genCmd.Stderr = os.Stderr
	if err := genCmd.Run(); err != nil {
		return fmt.Errorf("failed to generate GPG key: %w", err)
	}
	fmt.Println(ui.Green("GPG key generated"))

	configureGitGPG(email)
	return nil
}

func configureGitGPG(email string) {
	listCmd := exec.Command("gpg", "--list-secret-keys", "--keyid-format", "long", email)
	output, err := listCmd.Output()
	if err != nil {
		ui.PrintError("Failed to list GPG keys")
		return
	}

	lines := strings.Split(string(output), "\n")
	var keyID string
	for _, line := range lines {
		if strings.Contains(line, "ed25519/") || strings.Contains(line, "rsa") {
			parts := strings.Split(line, "/")
			if len(parts) >= 2 {
				keyID = strings.Fields(parts[1])[0]
				break
			}
		}
	}

	if keyID == "" {
		ui.PrintError("Could not find GPG key ID")
		return
	}

	fmt.Println("Configuring Git to use GPG key...")

	exec.Command("git", "config", "--global", "user.signingkey", keyID).Run()
	exec.Command("git", "config", "--global", "commit.gpgsign", "true").Run()
	exec.Command("git", "config", "--global", "gpg.program", "gpg").Run()

	fmt.Println(ui.Green("Git configured for commit signing"))

	fmt.Println()
	fmt.Println("Your GPG public key (add to GitHub):")
	fmt.Println("----------------------------------------")
	exportCmd := exec.Command("gpg", "--armor", "--export", email)
	exportCmd.Stdout = os.Stdout
	exportCmd.Run()
	fmt.Println("----------------------------------------")
	fmt.Println("Add at: https://github.com/settings/gpg/new")

	fmt.Println()
	fmt.Println(ui.Green("GPG setup completed"))
	fmt.Println(ui.Dim("All future commits will be signed automatically"))
}

func runSSHSetup() error {
	fmt.Println(ui.Cyan("Setting up SSH..."))

	sshDir := os.Getenv("HOME") + "/.ssh"
	sshKey := sshDir + "/id_ed25519"
	email := UserEmail

	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return fmt.Errorf("failed to create .ssh directory: %w", err)
	}

	if _, err := os.Stat(sshKey); err == nil {
		fmt.Printf("%s SSH key already exists at %s\n", ui.Green("Done"), sshKey)
	} else {
		fmt.Println("Generating SSH key with macOS Keychain integration...")
		fmt.Println(ui.Dim("You'll be prompted to create a passphrase"))
		fmt.Println()

		genCmd := exec.Command("ssh-keygen", "-t", "ed25519", "-C", email, "-f", sshKey)
		genCmd.Stdin = os.Stdin
		genCmd.Stdout = os.Stdout
		genCmd.Stderr = os.Stderr
		if err := genCmd.Run(); err != nil {
			return fmt.Errorf("failed to generate SSH key: %w", err)
		}
		fmt.Println(ui.Green("SSH key generated"))
	}

	fmt.Println("Configuring SSH...")
	sshConfig := sshDir + "/config"

	existingConfig, _ := os.ReadFile(sshConfig)
	if !strings.Contains(string(existingConfig), "AddKeysToAgent yes") {
		configContent := `
Host *
  AddKeysToAgent yes
  UseKeychain yes
  IdentityFile ~/.ssh/id_ed25519
`
		f, err := os.OpenFile(sshConfig, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err == nil {
			f.WriteString(configContent)
			f.Close()
			fmt.Println(ui.Green("SSH config updated"))
		}
	} else {
		fmt.Println(ui.Green("SSH config already configured"))
	}

	fmt.Println("Adding key to SSH agent with Keychain...")
	fmt.Println(ui.Dim("Passphrase will be stored in macOS Keychain"))
	fmt.Println()

	addCmd := exec.Command("ssh-add", "--apple-use-keychain", sshKey)
	addCmd.Stdin = os.Stdin
	addCmd.Stdout = os.Stdout
	addCmd.Stderr = os.Stderr
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed to add key to SSH agent: %w", err)
	}

	fmt.Println()
	fmt.Println("Your public key (add to GitHub):")
	fmt.Println("----------------------------------------")
	pubKey, _ := os.ReadFile(sshKey + ".pub")
	fmt.Println(string(pubKey))
	fmt.Println("----------------------------------------")
	fmt.Println("Add at: https://github.com/settings/ssh/new")

	fmt.Println(ui.Green("SSH setup completed"))
	return nil
}

func runZedConfig() error {
	fmt.Println(ui.Cyan("Setting up Zed config..."))

	configDir := os.Getenv("HOME") + "/.config/zed"
	configPath := configDir + "/settings.json"

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	repoConfig, err := GetRepoConfigPath("configuration/applications/zed/settings.json")
	if err != nil {
		return fmt.Errorf("failed to find repo config: %w", err)
	}

	configContent, err := os.ReadFile(repoConfig)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", repoConfig, err)
	}

	if err := os.WriteFile(configPath, configContent, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", configPath, err)
	}

	fmt.Println(ui.Green("Done - Zed config installed"))
	return nil
}

func runJavaSymlink() error {
	fmt.Println(ui.Cyan("Setting up Java runtime..."))

	brewJava := "/opt/homebrew/opt/openjdk/libexec/openjdk.jdk"
	if _, err := os.Stat(brewJava); err != nil {
		return fmt.Errorf("OpenJDK not installed. Run: j install openjdk")
	}

	symlinkPath := "/Library/Java/JavaVirtualMachines/openjdk.jdk"

	if _, err := os.Lstat(symlinkPath); err == nil {
		fmt.Printf("%s Java symlink already exists\n", ui.Green("Done"))
		return nil
	}

	fmt.Println("Creating symlink for macOS Java recognition...")

	sudoCmd := exec.Command("sudo", "ln", "-sfn", brewJava, symlinkPath)
	sudoCmd.Stdout = os.Stdout
	sudoCmd.Stderr = os.Stderr
	sudoCmd.Stdin = os.Stdin
	if err := sudoCmd.Run(); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	fmt.Println(ui.Green("Done - Java configured for macOS"))
	return nil
}

func runDockReset() error {
	fmt.Println(ui.Cyan("Resetting macOS Dock..."))
	ExecCommand("defaults", "delete", "com.apple.dock")
	ExecCommand("killall", "Dock")
	fmt.Println(ui.Green("Done - Dock reset to defaults"))
	return nil
}

func runDockSpacer() error {
	fmt.Println(ui.Cyan("Adding spacer to Dock..."))
	ExecCommand("defaults", "write", "com.apple.dock", "persistent-apps", "-array-add", `{"tile-type"="small-spacer-tile";}`)
	ExecCommand("killall", "Dock")
	fmt.Println(ui.Green("Done - Dock spacer added"))
	return nil
}

// =============================================================================
// Helper Functions
// =============================================================================

// GetRepoConfigPath returns the full path for a file in the repo
func GetRepoConfigPath(relativePath string) (string, error) {
	possibleRoots := []string{
		os.Getenv("HOME") + "/Developer/jterrazz-cli",
		"/usr/local/share/jterrazz-cli",
	}

	for _, root := range possibleRoots {
		fullPath := root + "/" + relativePath
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil
		}
	}

	return "", fmt.Errorf("config file not found: %s", relativePath)
}

// CommandExists checks if a command is available in PATH
// Re-exported from system package for convenience
var CommandExists = tool.CommandExists

// =============================================================================
// Script Functions
// =============================================================================

// GetAllScripts returns all scripts
func GetAllScripts() []Script {
	return Scripts
}

// GetScriptByName returns a script by name
func GetScriptByName(name string) *Script {
	for i := range Scripts {
		if Scripts[i].Name == name {
			return &Scripts[i]
		}
	}
	return nil
}

// GetScriptsByCategory returns scripts filtered by category
func GetScriptsByCategory(category ScriptCategory) []Script {
	var result []Script
	for _, script := range Scripts {
		if script.Category == category {
			result = append(result, script)
		}
	}
	return result
}

// GetScriptsForTool returns scripts that belong to a tool
func GetScriptsForTool(toolName string) []Script {
	tool := GetToolByName(toolName)
	if tool == nil || len(tool.Scripts) == 0 {
		return nil
	}

	var result []Script
	for _, scriptName := range tool.Scripts {
		if script := GetScriptByName(scriptName); script != nil {
			result = append(result, *script)
		}
	}
	return result
}

// GetStandaloneScripts returns scripts not attached to any tool
func GetStandaloneScripts() []Script {
	// Build set of tool-attached scripts
	attached := make(map[string]bool)
	for _, tool := range Tools {
		for _, scriptName := range tool.Scripts {
			attached[scriptName] = true
		}
	}

	var result []Script
	for _, script := range Scripts {
		if !attached[script.Name] {
			result = append(result, script)
		}
	}
	return result
}

// GetConfigurableScripts returns scripts that have a CheckFn (can be checked)
func GetConfigurableScripts() []Script {
	var result []Script
	for _, script := range Scripts {
		if script.CheckFn != nil {
			result = append(result, script)
		}
	}
	return result
}

// GetUnconfiguredScripts returns scripts that haven't been run yet
func GetUnconfiguredScripts() []Script {
	var result []Script
	for _, script := range Scripts {
		if script.CheckFn != nil {
			check := script.CheckFn()
			if !check.Installed {
				result = append(result, script)
			}
		}
	}
	return result
}

// CheckScript checks if a script has been configured
func CheckScript(script Script) CheckResult {
	if script.CheckFn != nil {
		return script.CheckFn()
	}
	return CheckResult{} // No check = unknown state
}

// RunScript runs a script (placeholder - actual implementation in commands)
func RunScript(script Script) error {
	if script.RunFn != nil {
		return script.RunFn()
	}
	// Scripts without RunFn are invoked via `j setup <RunCmd>`
	return nil
}
