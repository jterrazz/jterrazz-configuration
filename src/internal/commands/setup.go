package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup system configurations (interactive)",
	Run: func(cmd *cobra.Command, args []string) {
		runSetupUI()
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
	rootCmd.AddCommand(setupCmd)
}

// runSetupItem runs a setup item by name (used by install command for SetupCmd)
func runSetupItem(name string) {
	switch name {
	case "java":
		setupJavaCmd.Run(nil, nil)
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

// Setup TUI

// setupItemData holds domain-specific data for each setup item
type setupItemData struct {
	name       string
	configured *bool
}

type setupModel struct {
	list       *ui.List
	page       *ui.Page
	itemData   []setupItemData
	processing bool
	quitting   bool
	// Skills sub-view
	showSkills  bool
	skillsModel *skillsModel
}

type setupItemDef struct {
	name        string
	description string
	checkFn     func() *bool // nil for utilities
}

var setupConfigItems = []setupItemDef{
	{"ghostty", "Install Ghostty terminal config", checkGhostty},
	{"gpg", "Generate GPG key for commit signing", checkGPG},
	{"hushlogin", "Silence terminal login message", checkHushlogin},
	{"java", "Configure Java runtime symlink for macOS", checkJava},
	{"ohmyzsh", "Install and configure Oh My Zsh", checkOhMyZsh},
	{"ssh", "Generate SSH key with Keychain integration", checkSSH},
	{"zed", "Install Zed editor config", checkZed},
}

var setupUtilityItems = []setupItemDef{
	{"dock-reset", "Reset dock to system defaults", nil},
	{"dock-spacer", "Add a small spacer tile to the dock", nil},
}

func (m *setupModel) buildItems() ([]ui.Item, []setupItemData) {
	var items []ui.Item
	var data []setupItemData

	// Navigation section
	items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Navigation"})
	data = append(data, setupItemData{})

	items = append(items, ui.Item{Kind: ui.KindAction, Label: "skills", Description: "Manage AI agent skills"})
	data = append(data, setupItemData{name: "skills"})

	// Actions section
	items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Actions"})
	data = append(data, setupItemData{})

	items = append(items, ui.Item{Kind: ui.KindAction, Label: "Setup all missing"})
	data = append(data, setupItemData{name: "setup-missing"})

	// Configuration section - pending items first, then configured
	items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Configuration"})
	data = append(data, setupItemData{})

	var configuredItems []struct {
		item ui.Item
		data setupItemData
	}
	var notConfiguredItems []struct {
		item ui.Item
		data setupItemData
	}

	for _, def := range setupConfigItems {
		status := def.checkFn()
		state := ui.StateUnchecked
		if status != nil && *status {
			state = ui.StateChecked
		}
		entry := struct {
			item ui.Item
			data setupItemData
		}{
			item: ui.Item{
				Kind:        ui.KindToggle,
				Label:       def.name,
				Description: def.description,
				State:       state,
			},
			data: setupItemData{name: def.name, configured: status},
		}
		if status != nil && *status {
			configuredItems = append(configuredItems, entry)
		} else {
			notConfiguredItems = append(notConfiguredItems, entry)
		}
	}

	// Add pending items first, then configured
	for _, entry := range notConfiguredItems {
		items = append(items, entry.item)
		data = append(data, entry.data)
	}
	for _, entry := range configuredItems {
		items = append(items, entry.item)
		data = append(data, entry.data)
	}

	// Scripts section
	items = append(items, ui.Item{Kind: ui.KindHeader, Label: "Scripts"})
	data = append(data, setupItemData{})

	for _, def := range setupUtilityItems {
		items = append(items, ui.Item{
			Kind:        ui.KindAction,
			Label:       def.name,
			Description: def.description,
		})
		data = append(data, setupItemData{name: def.name})
	}

	return items, data
}

func (m *setupModel) rebuildItems() {
	cursor := m.list.Cursor
	items, data := m.buildItems()
	m.list = ui.NewList(items)
	m.list.CalculateLabelWidth()
	m.itemData = data

	// Restore cursor position
	if cursor >= len(items) {
		cursor = len(items) - 1
	}
	m.list.SetCursor(cursor)

	// Skip headers if cursor landed on one
	for m.list.Cursor > 0 && !m.list.Items[m.list.Cursor].Selectable() {
		m.list.Cursor--
	}
}

func checkGhostty() *bool {
	_, err := os.Stat(os.Getenv("HOME") + "/.config/ghostty/config")
	result := err == nil
	return &result
}

func checkGPG() *bool {
	out, _ := exec.Command("git", "config", "--global", "commit.gpgsign").Output()
	result := strings.TrimSpace(string(out)) == "true"
	return &result
}

func checkHushlogin() *bool {
	_, err := os.Stat(os.Getenv("HOME") + "/.hushlogin")
	result := err == nil
	return &result
}

func checkJava() *bool {
	_, err := os.Lstat("/Library/Java/JavaVirtualMachines/openjdk.jdk")
	result := err == nil
	return &result
}

func checkOhMyZsh() *bool {
	_, err := os.Stat(os.Getenv("HOME") + "/.oh-my-zsh")
	result := err == nil
	return &result
}

func checkSSH() *bool {
	_, err := os.Stat(os.Getenv("HOME") + "/.ssh/id_ed25519")
	result := err == nil
	return &result
}

func checkZed() *bool {
	_, err := os.Stat(os.Getenv("HOME") + "/.config/zed/settings.json")
	result := err == nil
	return &result
}

func initialSetupModel() setupModel {
	m := setupModel{
		page: ui.NewPage("Setup"),
	}

	items, data := m.buildItems()
	m.list = ui.NewList(items)
	m.list.CalculateLabelWidth()
	m.itemData = data

	return m
}

func (m setupModel) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m setupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If showing skills sub-view, delegate to it
	if m.showSkills && m.skillsModel != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			// Check for back/quit in skills view
			if key.Matches(msg, key.NewBinding(key.WithKeys("esc"))) {
				m.showSkills = false
				m.skillsModel = nil
				return m, nil
			}
			if key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))) {
				m.quitting = true
				return m, tea.Quit
			}
		case tea.WindowSizeMsg:
			m.list.SetSize(msg.Width, msg.Height)
			m.page.SetSize(msg.Width, msg.Height)
			m.skillsModel.list.SetSize(msg.Width, msg.Height)
			m.skillsModel.page.SetSize(msg.Width, msg.Height)
		}
		newModel, cmd := m.skillsModel.Update(msg)
		if sm, ok := newModel.(skillsModel); ok {
			m.skillsModel = &sm
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.processing {
			return m, nil
		}

		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"))):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			m.list.Up()

		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			m.list.Down()

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			return m.handleSelect()
		}

	case setupActionDoneMsg:
		m.processing = false
		m.page.Message = msg.message
		m.page.Processing = false
		m.rebuildItems()
		return m, nil

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		m.page.SetSize(msg.Width, msg.Height)
	}

	return m, nil
}

type setupActionDoneMsg struct {
	message string
	err     error
}

func (m setupModel) handleSelect() (setupModel, tea.Cmd) {
	idx := m.list.SelectedIndex()
	if idx < 0 || idx >= len(m.itemData) {
		return m, nil
	}

	item := m.list.Selected()
	data := m.itemData[idx]

	if item.Kind == ui.KindHeader {
		return m, nil
	}

	switch data.name {
	case "skills":
		sm := initialSkillsModel()
		sm.list.SetSize(m.page.Width, m.page.Height)
		sm.page.SetSize(m.page.Width, m.page.Height)
		m.skillsModel = &sm
		m.showSkills = true
		return m, nil

	case "setup-missing":
		m.processing = true
		m.page.Processing = true
		m.page.Message = "Setting up all missing..."
		return m, m.runSetupMissing()

	default:
		// Config item or utility
		m.processing = true
		m.page.Processing = true
		m.page.Message = fmt.Sprintf("Running %s...", data.name)
		return m, m.runSetup(data.name)
	}
}

func (m setupModel) runSetupMissing() tea.Cmd {
	return func() tea.Msg {
		count := 0
		for _, def := range setupConfigItems {
			status := def.checkFn()
			if status == nil || !*status {
				runSetupByName(def.name)
				count++
			}
		}
		if count == 0 {
			return setupActionDoneMsg{message: "Everything already configured"}
		}
		return setupActionDoneMsg{message: fmt.Sprintf("Configured %d items", count)}
	}
}

func runSetupByName(name string) {
	switch name {
	case "dock-reset":
		exec.Command("defaults", "delete", "com.apple.dock").Run()
		exec.Command("killall", "Dock").Run()
	case "dock-spacer":
		exec.Command("defaults", "write", "com.apple.dock", "persistent-apps", "-array-add", `{"tile-type"="small-spacer-tile";}`).Run()
		exec.Command("killall", "Dock").Run()
	case "ghostty":
		setupGhosttyCmd.Run(nil, nil)
	case "gpg":
		setupGPGCmd.Run(nil, nil)
	case "hushlogin":
		f, _ := os.Create(os.Getenv("HOME") + "/.hushlogin")
		if f != nil {
			f.Close()
		}
	case "java":
		setupJavaCmd.Run(nil, nil)
	case "ohmyzsh":
		setupOhMyZshCmd.Run(nil, nil)
	case "ssh":
		setupSSHCmd.Run(nil, nil)
	case "zed":
		setupZedCmd.Run(nil, nil)
	}
}

func (m setupModel) runSetup(name string) tea.Cmd {
	return func() tea.Msg {
		runSetupByName(name)
		return setupActionDoneMsg{message: fmt.Sprintf("Completed %s", name)}
	}
}

func (m setupModel) View() string {
	if m.quitting {
		return ""
	}

	// Show skills sub-view if active
	if m.showSkills && m.skillsModel != nil {
		return m.skillsModel.viewWithBreadcrumb("Setup", "Skills")
	}

	m.page.Help = ui.DefaultHelp()
	m.page.Content = m.list.Render(m.page.ContentHeight())

	return m.page.Render()
}

func runSetupUI() {
	m := initialSetupModel()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running setup UI: %v\n", err)
		os.Exit(1)
	}
}
