package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup development tools",
}

var setupBrewCmd = &cobra.Command{
	Use:   "brew",
	Short: "Install Homebrew + development packages",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Println(cyan("🍺 Installing Homebrew..."))

		if !commandExists("brew") {
			fmt.Println("📥 Downloading and installing Homebrew...")
			runCommand("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)")
		} else {
			fmt.Println(green("✅ Homebrew already installed"))
		}

		fmt.Println(cyan("📦 Installing essential development packages..."))

		brewPackages := []struct {
			name    string
			cmd     string
			formula string
			cask    bool
		}{
			{"ansible-lint", "ansible-lint", "ansible-lint", false},
			{"ansible", "ansible", "ansible", false},
			{"terraform", "terraform", "terraform", false},
			{"kubectl", "kubectl", "kubectl", false},
			{"multipass", "multipass", "multipass", false},
			{"bun", "bun", "bun", false},
			{"python", "python3", "python", false},
			{"neohtop", "", "neohtop", true},
			{"codex", "codex", "codex", false},
			{"mole", "mo", "tw93/tap/mole", false},
			{"gemini-cli", "gemini", "gemini-cli", false},
		}

		for _, pkg := range brewPackages {
			if pkg.cmd != "" && commandExists(pkg.cmd) {
				fmt.Printf("  %s %s already installed\n", green("✅"), pkg.name)
			} else if pkg.cask {
				fmt.Printf("  📥 Installing %s...\n", pkg.name)
				runCommand("brew", "install", "--cask", pkg.formula)
			} else {
				fmt.Printf("  📥 Installing %s...\n", pkg.name)
				runCommand("brew", "install", pkg.formula)
			}
		}

		// Claude via npm
		if commandExists("claude") {
			fmt.Printf("  %s claude already installed\n", green("✅"))
		} else {
			fmt.Println("  📥 Installing claude...")
			if commandExists("npm") {
				runCommand("npm", "install", "-g", "@anthropic-ai/claude-code")
			} else {
				printError("npm not found, cannot install claude")
			}
		}

		fmt.Println(green("✅ Development packages check completed"))
	},
}

var setupOhMyZshCmd = &cobra.Command{
	Use:   "ohmyzsh",
	Short: "Install Oh My Zsh and configure shell",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Println(cyan("🐚 Installing Oh My Zsh..."))

		omzPath := os.Getenv("HOME") + "/.oh-my-zsh"
		if _, err := os.Stat(omzPath); err == nil {
			fmt.Println(green("✅ Oh My Zsh already installed"))
			return
		}

		fmt.Println("📥 Downloading and installing Oh My Zsh...")
		runCommand("sh", "-c", "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh) --unattended")

		fmt.Println(green("✅ Oh My Zsh installed successfully"))
	},
}

var setupNvmCmd = &cobra.Command{
	Use:   "nvm",
	Short: "Install NVM and Node.js stable",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Println(cyan("📦 Installing NVM (Node Version Manager)..."))

		if !commandExists("brew") {
			printError("Homebrew required for NVM installation")
			fmt.Println("💡 Run: j setup brew")
			return
		}

		fmt.Println("📥 Installing NVM via Homebrew...")
		runCommand("brew", "install", "nvm")

		fmt.Println("⚙️  Setting up NVM...")
		nvmDir := os.Getenv("HOME") + "/.nvm"
		os.MkdirAll(nvmDir, 0755)

		fmt.Println(green("✅ NVM installed - restart terminal and run 'nvm install stable'"))
	},
}

var setupGitSSHCmd = &cobra.Command{
	Use:   "git-ssh",
	Short: "Generate SSH key and configure Git",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Println(cyan("🔑 Setting up Git SSH..."))

		sshKey := os.Getenv("HOME") + "/.ssh/id_github"
		email := "contact@jterrazz.com"

		if _, err := os.Stat(sshKey); err == nil {
			fmt.Printf("%s SSH key already exists at %s\n", green("✅"), sshKey)
		} else {
			fmt.Println("🔐 Generating SSH key...")
			runCommand("ssh-keygen", "-t", "ed25519", "-C", email, "-f", sshKey, "-N", "")
			fmt.Println(green("✅ SSH key generated"))
		}

		// Configure SSH
		fmt.Println("⚙️  Configuring SSH...")
		sshConfig := os.Getenv("HOME") + "/.ssh/config"

		configContent := `
Host github.com
  AddKeysToAgent yes
  UseKeychain yes
  IdentityFile ~/.ssh/id_github
`
		f, err := os.OpenFile(sshConfig, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			f.WriteString(configContent)
			f.Close()
			fmt.Println(green("✅ SSH config updated"))
		}

		// Add key to SSH agent
		fmt.Println("🔗 Adding key to SSH agent...")
		runCommand("ssh-add", "--apple-use-keychain", sshKey)

		fmt.Println("📋 Your public key (add this to GitHub):")
		fmt.Println("----------------------------------------")
		pubKey, _ := os.ReadFile(sshKey + ".pub")
		fmt.Println(string(pubKey))
		fmt.Println("----------------------------------------")
		fmt.Println("💡 Copy the above key and add it to: https://github.com/settings/ssh/new")

		fmt.Println(green("✅ Git SSH setup completed"))
	},
}

var setupAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Setup everything (brew, ohmyzsh, nvm, git-ssh)",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Println(cyan("🚀 Setting up full development environment..."))
		setupBrewCmd.Run(cmd, args)
		setupOhMyZshCmd.Run(cmd, args)
		setupNvmCmd.Run(cmd, args)
		setupGitSSHCmd.Run(cmd, args)
	},
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

func init() {
	setupCmd.AddCommand(setupBrewCmd)
	setupCmd.AddCommand(setupOhMyZshCmd)
	setupCmd.AddCommand(setupNvmCmd)
	setupCmd.AddCommand(setupGitSSHCmd)
	setupCmd.AddCommand(setupAllCmd)
	setupCmd.AddCommand(setupDockSpacerCmd)
	setupCmd.AddCommand(setupDockResetCmd)
	rootCmd.AddCommand(setupCmd)
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
