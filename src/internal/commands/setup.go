package commands

import (
	"fmt"
	"os"
	"os/exec"

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
  j setup git-ssh            Setup Git SSH key
  j setup ohmyzsh git-ssh    Setup specific items
  j setup                    List available setup items`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		suggestions := []string{"ohmyzsh", "git-ssh", "dock-spacer", "dock-reset", "java"}
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
			runSetupItem("git-ssh")
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
		{"ohmyzsh", "Install and configure Oh My Zsh", func() *bool {
			_, err := os.Stat(os.Getenv("HOME") + "/.oh-my-zsh")
			result := err == nil
			return &result
		}},
		{"git-ssh", "Generate SSH key and configure Git", func() *bool {
			_, err := os.Stat(os.Getenv("HOME") + "/.ssh/id_github")
			result := err == nil
			return &result
		}},
		{"java", "Configure Java runtime symlink for macOS", func() *bool {
			_, err := os.Lstat("/Library/Java/JavaVirtualMachines/openjdk.jdk")
			result := err == nil
			return &result
		}},
		{"dock-spacer", "Add a small spacer tile to the dock", func() *bool { return nil }},
		{"dock-reset", "Reset dock to system defaults", func() *bool { return nil }},
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
	case "ohmyzsh":
		setupOhMyZshCmd.Run(nil, nil)
	case "git-ssh":
		setupGitSSHCmd.Run(nil, nil)
	case "java":
		setupJavaCmd.Run(nil, nil)
	case "dock-spacer":
		setupDockSpacerCmd.Run(nil, nil)
	case "dock-reset":
		setupDockResetCmd.Run(nil, nil)
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
