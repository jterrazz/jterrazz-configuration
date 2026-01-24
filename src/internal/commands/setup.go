package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var setupAll bool

var setupCmd = &cobra.Command{
	Use:   "setup [item...]",
	Short: "Setup system configurations",
	Long: `Setup system configurations.

Examples:
  j setup --all              Setup all configurations
  j setup ohmyzsh            Setup Oh My Zsh
  j setup ssh                Setup SSH key
  j setup ohmyzsh ssh        Setup specific items
  j setup                    List available setup items`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		suggestions := []string{"dock-reset", "dock-spacer", "ghostty", "gpg", "hushlogin", "java", "ohmyzsh", "ssh", "zed"}
		var filtered []string
		for _, s := range suggestions {
			alreadyUsed := false
			for _, arg := range args {
				if arg == s {
					alreadyUsed = true
					break
				}
			}
			if !alreadyUsed {
				filtered = append(filtered, s)
			}
		}
		return filtered, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()

		if setupAll {
			fmt.Println(cyan("🚀 Setting up all configurations..."))
			runSetupItem("ohmyzsh")
			runSetupItem("ssh")
			return
		}

		if len(args) == 0 {
			listSetupItems()
			return
		}

		fmt.Println(cyan("🚀 Setting up selected configurations..."))
		for _, name := range args {
			runSetupItem(name)
		}
	},
}

var setupOhMyZshCmd = &cobra.Command{
	Use:   "ohmyzsh",
	Short: "Install and configure Oh My Zsh",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Println(cyan("🐚 Setting up Oh My Zsh..."))

		omzPath := os.Getenv("HOME") + "/.oh-my-zsh"
		if _, err := os.Stat(omzPath); err == nil {
			fmt.Println(green("✅ Oh My Zsh already installed"))
			return
		}

		fmt.Println("📥 Downloading and installing Oh My Zsh...")
		runCommand("sh", "-c", "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh) --unattended")

		fmt.Println(green("✅ Oh My Zsh configured"))
	},
}

var setupSSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Generate SSH key with macOS Keychain integration",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		dim := color.New(color.FgHiBlack).SprintFunc()

		fmt.Println(cyan("🔑 Setting up Git SSH..."))

		sshDir := os.Getenv("HOME") + "/.ssh"
		sshKey := sshDir + "/id_ed25519"
		email := "admin@jterrazz.com"

		// Ensure .ssh directory exists with correct permissions
		if err := os.MkdirAll(sshDir, 0700); err != nil {
			printError(fmt.Sprintf("Failed to create .ssh directory: %v", err))
			return
		}

		if _, err := os.Stat(sshKey); err == nil {
			fmt.Printf("%s SSH key already exists at %s\n", green("✅"), sshKey)
		} else {
			fmt.Println("🔐 Generating SSH key with macOS Keychain integration...")
			fmt.Println(dim("   You'll be prompted to create a passphrase (stored securely in Keychain)"))
			fmt.Println()

			// Generate key with passphrase prompt (user enters it interactively)
			// Using ed25519 which is the current best practice for SSH keys
			genCmd := exec.Command("ssh-keygen",
				"-t", "ed25519",
				"-C", email,
				"-f", sshKey,
			)
			genCmd.Stdin = os.Stdin
			genCmd.Stdout = os.Stdout
			genCmd.Stderr = os.Stderr
			if err := genCmd.Run(); err != nil {
				printError(fmt.Sprintf("Failed to generate SSH key: %v", err))
				return
			}
			fmt.Println(green("✅ SSH key generated"))
		}

		// Configure SSH with macOS Keychain integration
		fmt.Println("⚙️  Configuring SSH...")
		sshConfig := sshDir + "/config"

		// Configure SSH to use macOS Keychain for all hosts
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
				fmt.Println(green("✅ SSH config updated"))
			}
		} else {
			fmt.Println(green("✅ SSH config already configured"))
		}

		// Add key to SSH agent with Keychain storage
		fmt.Println("🔗 Adding key to SSH agent with Keychain...")
		fmt.Println(dim("   Your passphrase will be stored in macOS Keychain"))
		fmt.Println(dim("   Future authentications will use Touch ID or auto-unlock"))
		fmt.Println()

		addCmd := exec.Command("ssh-add", "--apple-use-keychain", sshKey)
		addCmd.Stdin = os.Stdin
		addCmd.Stdout = os.Stdout
		addCmd.Stderr = os.Stderr
		if err := addCmd.Run(); err != nil {
			printError(fmt.Sprintf("Failed to add key to SSH agent: %v", err))
			return
		}

		fmt.Println()
		fmt.Println("📋 Your public key (add this to GitHub):")
		fmt.Println("----------------------------------------")
		pubKey, _ := os.ReadFile(sshKey + ".pub")
		fmt.Println(string(pubKey))
		fmt.Println("----------------------------------------")
		fmt.Println("💡 Copy the above key and add it to: https://github.com/settings/ssh/new")

		fmt.Println(green("✅ Git SSH setup completed"))
		fmt.Println(dim("   Passphrase stored in macOS Keychain - unlocks automatically"))
	},
}

var setupGPGCmd = &cobra.Command{
	Use:   "gpg",
	Short: "Generate GPG key and configure Git commit signing",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		dim := color.New(color.FgHiBlack).SprintFunc()

		fmt.Println(cyan("🔐 Setting up GPG for Git commit signing..."))

		email := "admin@jterrazz.com"
		name := "Jean-Baptiste Music"

		// Check if gpg is installed
		if !commandExists("gpg") {
			printError("GPG not installed. Run: brew install gnupg")
			return
		}

		// Check if a key already exists for this email
		checkCmd := exec.Command("gpg", "--list-secret-keys", "--keyid-format", "long", email)
		if output, err := checkCmd.Output(); err == nil && len(output) > 0 {
			fmt.Println(green("✅ GPG key already exists for " + email))
			// Extract key ID and configure git
			configureGitGPG(email, dim, green)
			return
		}

		fmt.Println("🔑 Generating GPG key...")
		fmt.Println(dim("   Using ed25519 algorithm (modern, secure)"))
		fmt.Println()

		// Generate key using batch mode with ed25519
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
			printError(fmt.Sprintf("Failed to generate GPG key: %v", err))
			return
		}
		fmt.Println(green("✅ GPG key generated"))

		configureGitGPG(email, dim, green)
	},
}

func configureGitGPG(email string, dim, green func(a ...interface{}) string) {
	// Get the key ID
	listCmd := exec.Command("gpg", "--list-secret-keys", "--keyid-format", "long", email)
	output, err := listCmd.Output()
	if err != nil {
		printError("Failed to list GPG keys")
		return
	}

	// Parse key ID from output (format: "ed25519/KEYID")
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
		printError("Could not find GPG key ID")
		return
	}

	fmt.Println("⚙️  Configuring Git to use GPG key...")

	// Configure git
	exec.Command("git", "config", "--global", "user.signingkey", keyID).Run()
	exec.Command("git", "config", "--global", "commit.gpgsign", "true").Run()
	exec.Command("git", "config", "--global", "gpg.program", "gpg").Run()

	fmt.Println(green("✅ Git configured for commit signing"))

	// Export public key
	fmt.Println()
	fmt.Println("📋 Your GPG public key (add this to GitHub):")
	fmt.Println("----------------------------------------")
	exportCmd := exec.Command("gpg", "--armor", "--export", email)
	exportCmd.Stdout = os.Stdout
	exportCmd.Run()
	fmt.Println("----------------------------------------")
	fmt.Println("💡 Copy the above key and add it to: https://github.com/settings/gpg/new")

	fmt.Println()
	fmt.Println(green("✅ GPG setup completed"))
	fmt.Println(dim("   All future commits will be signed automatically"))
}

var setupDockSpacerCmd = &cobra.Command{
	Use:   "dock-spacer",
	Short: "Add a small spacer tile to the dock",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Println(cyan("🔧 Adding spacer to macOS Dock..."))
		runCommand("defaults", "write", "com.apple.dock", "persistent-apps", "-array-add", `{"tile-type"="small-spacer-tile";}`)
		runCommand("killall", "Dock")
		fmt.Println(green("✅ Dock spacer added and restarted"))
	},
}

var setupDockResetCmd = &cobra.Command{
	Use:   "dock-reset",
	Short: "Reset dock to system defaults",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Println(cyan("🔧 Resetting macOS Dock to defaults..."))
		runCommand("defaults", "delete", "com.apple.dock")
		runCommand("killall", "Dock")
		fmt.Println(green("✅ Dock reset to defaults"))
	},
}

var setupGhosttyCmd = &cobra.Command{
	Use:   "ghostty",
	Short: "Install Ghostty terminal configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Println(cyan("👻 Setting up Ghostty config..."))

		configDir := os.Getenv("HOME") + "/.config/ghostty"
		configPath := configDir + "/config"

		// Create config directory if it doesn't exist
		if err := os.MkdirAll(configDir, 0755); err != nil {
			printError(fmt.Sprintf("Failed to create config directory: %v", err))
			return
		}

		// Get the source config from the repo
		repoConfig, err := getRepoConfigPath("configuration/applications/ghostty/config")
		if err != nil {
			printError(fmt.Sprintf("Failed to find repo config: %v", err))
			return
		}

		configContent, err := os.ReadFile(repoConfig)
		if err != nil {
			printError(fmt.Sprintf("Failed to read config file: %v", err))
			return
		}

		// Write config file
		if err := os.WriteFile(configPath, configContent, 0644); err != nil {
			printError(fmt.Sprintf("Failed to write config file: %v", err))
			return
		}

		fmt.Println(green("✅ Ghostty config installed at ~/.config/ghostty/config"))
	},
}

var setupZedCmd = &cobra.Command{
	Use:   "zed",
	Short: "Install Zed editor configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Println(cyan("⚡ Setting up Zed config..."))

		configDir := os.Getenv("HOME") + "/.config/zed"
		configPath := configDir + "/settings.json"

		// Create config directory if it doesn't exist
		if err := os.MkdirAll(configDir, 0755); err != nil {
			printError(fmt.Sprintf("Failed to create config directory: %v", err))
			return
		}

		// Get the source config from the repo
		repoConfig, err := getRepoConfigPath("configuration/applications/zed/settings.json")
		if err != nil {
			printError(fmt.Sprintf("Failed to find repo config: %v", err))
			return
		}

		configContent, err := os.ReadFile(repoConfig)
		if err != nil {
			printError(fmt.Sprintf("Failed to read config file: %v", err))
			return
		}

		// Write config file
		if err := os.WriteFile(configPath, configContent, 0644); err != nil {
			printError(fmt.Sprintf("Failed to write config file: %v", err))
			return
		}

		fmt.Println(green("✅ Zed config installed at ~/.config/zed/settings.json"))
	},
}

// getRepoConfigPath returns the absolute path to a config file in the repo
func getRepoConfigPath(relativePath string) (string, error) {
	// Try to find the repo root by looking for known paths
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

var setupHushloginCmd = &cobra.Command{
	Use:   "hushlogin",
	Short: "Create .hushlogin to silence terminal login message",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Println(cyan("🤫 Setting up hushlogin..."))

		hushPath := os.Getenv("HOME") + "/.hushlogin"
		if _, err := os.Stat(hushPath); err == nil {
			fmt.Printf("%s .hushlogin already exists at %s\n", green("✅"), hushPath)
			return
		}

		// Create empty .hushlogin file
		f, err := os.Create(hushPath)
		if err != nil {
			printError(fmt.Sprintf("Failed to create .hushlogin: %v", err))
			return
		}
		f.Close()

		fmt.Println(green("✅ .hushlogin created - terminal login message silenced"))
	},
}

var setupJavaCmd = &cobra.Command{
	Use:   "java",
	Short: "Configure Java runtime symlink for macOS",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Println(cyan("☕ Setting up Java runtime..."))

		// Check if openjdk is installed via Homebrew
		brewJava := "/opt/homebrew/opt/openjdk/libexec/openjdk.jdk"
		if _, err := os.Stat(brewJava); err != nil {
			printError("OpenJDK not installed. Run: j install openjdk")
			return
		}

		// Create symlink for macOS to recognize Java
		jvmDir := "/Library/Java/JavaVirtualMachines"
		symlinkPath := jvmDir + "/openjdk.jdk"

		// Check if symlink already exists
		if _, err := os.Lstat(symlinkPath); err == nil {
			fmt.Printf("%s Java symlink already exists at %s\n", green("✅"), symlinkPath)
			return
		}

		fmt.Println("🔗 Creating symlink for macOS Java recognition...")
		fmt.Printf("   %s -> %s\n", symlinkPath, brewJava)

		// Need sudo for /Library/Java/JavaVirtualMachines
		sudoCmd := exec.Command("sudo", "ln", "-sfn", brewJava, symlinkPath)
		sudoCmd.Stdout = os.Stdout
		sudoCmd.Stderr = os.Stderr
		sudoCmd.Stdin = os.Stdin
		if err := sudoCmd.Run(); err != nil {
			printError(fmt.Sprintf("Failed to create symlink: %v", err))
			return
		}

		fmt.Println(green("✅ Java configured - macOS will now recognize OpenJDK"))
	},
}

func init() {
	setupCmd.Flags().BoolVarP(&setupAll, "all", "a", false, "Setup all configurations")
	rootCmd.AddCommand(setupCmd)
}

func listSetupItems() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	dim := color.New(color.FgHiBlack).SprintFunc()

	fmt.Println(cyan("Available setup items:"))
	fmt.Println()

	items := []struct {
		name        string
		description string
		check       func() *bool // nil = action (no state), true = configured, false = not configured
	}{
		{"dock-reset", "Reset dock to system defaults", func() *bool { return nil }},
		{"dock-spacer", "Add a small spacer tile to the dock", func() *bool { return nil }},
		{"ghostty", "Install Ghostty terminal config", func() *bool {
			_, err := os.Stat(os.Getenv("HOME") + "/.config/ghostty/config")
			result := err == nil
			return &result
		}},
		{"gpg", "Generate GPG key for commit signing", func() *bool {
			out, _ := exec.Command("git", "config", "--global", "commit.gpgsign").Output()
			result := strings.TrimSpace(string(out)) == "true"
			return &result
		}},
		{"hushlogin", "Silence terminal login message", func() *bool {
			_, err := os.Stat(os.Getenv("HOME") + "/.hushlogin")
			result := err == nil
			return &result
		}},
		{"java", "Configure Java runtime symlink for macOS", func() *bool {
			_, err := os.Lstat("/Library/Java/JavaVirtualMachines/openjdk.jdk")
			result := err == nil
			return &result
		}},
		{"ohmyzsh", "Install and configure Oh My Zsh", func() *bool {
			_, err := os.Stat(os.Getenv("HOME") + "/.oh-my-zsh")
			result := err == nil
			return &result
		}},
		{"ssh", "Generate SSH key with Keychain integration", func() *bool {
			_, err := os.Stat(os.Getenv("HOME") + "/.ssh/id_ed25519")
			result := err == nil
			return &result
		}},
		{"zed", "Install Zed editor config", func() *bool {
			_, err := os.Stat(os.Getenv("HOME") + "/.config/zed/settings.json")
			result := err == nil
			return &result
		}},
	}

	for _, item := range items {
		checkResult := item.check()
		var status string
		if checkResult == nil {
			status = dim("•")
		} else if *checkResult {
			status = green("✓")
		} else {
			status = red("✗")
		}
		fmt.Printf("  %s %-14s %s\n", status, item.name, dim(item.description))
	}

	fmt.Println()
	fmt.Println(dim("Usage: j setup <item> [item...]"))
	fmt.Println(dim("       j setup --all"))
}

func runSetupItem(name string) {
	switch name {
	case "dock-reset":
		setupDockResetCmd.Run(nil, nil)
	case "dock-spacer":
		setupDockSpacerCmd.Run(nil, nil)
	case "ghostty":
		setupGhosttyCmd.Run(nil, nil)
	case "gpg":
		setupGPGCmd.Run(nil, nil)
	case "hushlogin":
		setupHushloginCmd.Run(nil, nil)
	case "java":
		setupJavaCmd.Run(nil, nil)
	case "ohmyzsh":
		setupOhMyZshCmd.Run(nil, nil)
	case "ssh":
		setupSSHCmd.Run(nil, nil)
	case "zed":
		setupZedCmd.Run(nil, nil)
	default:
		printError(fmt.Sprintf("Unknown setup item: %s", name))
	}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func printError(msg string) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("%s %s\n", red("❌"), msg)
}

func printWarning(msg string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("%s %s\n", yellow("⚠️ "), msg)
}
