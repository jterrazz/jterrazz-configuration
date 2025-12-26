package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show comprehensive system status",
	Run: func(cmd *cobra.Command, args []string) {
		showStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

// Styles
var (
	subtle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	highlight  = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	special    = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	success    = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	warning    = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			MarginTop(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			PaddingLeft(1)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

func showStatus() {
	fmt.Println(titleStyle.Render("j status"))
	fmt.Println()

	printSystemSection()
	printDevelopmentSection()
	printInfraSection()
	printAISection()
	printSystemToolsSection()
	printConfigSection()
	printPackageStats()
	printDiskUsage()
}

func printSystemSection() {
	hostname, _ := os.Hostname()
	osInfo := getCommandOutput("uname", "-sr")
	arch := getCommandOutput("uname", "-m")
	user := os.Getenv("USER")
	shell := filepath.Base(os.Getenv("SHELL"))

	fmt.Println(headerStyle.Render("System"))
	fmt.Printf(" %s • %s\n", special.Render(osInfo), dimStyle.Render(arch))
	fmt.Printf(" %s • %s • %s\n\n", dimStyle.Render(hostname), dimStyle.Render(user), dimStyle.Render(shell))
}

func printDevelopmentSection() {
	tools := []toolInfo{
		getToolInfo("bun", "bun", []string{"--version"}, "brew", trimVersion),
		getDockerInfo(),
		getToolInfo("git", "git", []string{"--version"}, "xcode", parseGitVersion),
		getToolInfo("go", "go", []string{"version"}, "brew", parseGoVersion),
		getToolInfo("homebrew", "brew", []string{"--version"}, "-", parseBrewVersion),
		getToolInfo("node", "node", []string{"--version"}, "nvm", trimVersion),
		getToolInfo("npm", "npm", []string{"--version"}, "node", trimVersion),
		getToolInfo("python", "python3", []string{"--version"}, "brew", parsePythonVersion),
	}

	printToolTable("Development", tools)
}

func printInfraSection() {
	tools := []toolInfo{
		getToolInfo("ansible", "ansible", []string{"--version"}, "brew", parseAnsibleVersion),
		getToolInfo("kubectl", "kubectl", []string{"version", "--client", "-o", "yaml"}, "brew", parseKubectlVersion),
		getToolInfo("multipass", "multipass", []string{"--version"}, "brew", parseMultipassVersion),
		getToolInfo("terraform", "terraform", []string{"--version"}, "brew", parseTerraformVersion),
	}

	printToolTable("Infrastructure", tools)
}

func printAISection() {
	tools := []toolInfo{
		getToolInfo("claude", "claude", []string{"--version"}, "brew", parseClaudeVersion),
		getToolInfo("codex", "codex", []string{"--version"}, "brew", parseCodexVersion),
		getToolInfo("gemini", "gemini", []string{"--version"}, "brew", trimVersion),
	}

	printToolTable("AI Tools", tools)
}

func printSystemToolsSection() {
	tools := []toolInfo{
		getToolInfoWithHint("mole", "mo", []string{"--version"}, "brew", parseMoleVersion, "brew install tw93/tap/mole"),
	}

	printToolTable("System Tools", tools)
}

func printConfigSection() {
	fmt.Println(headerStyle.Render("Configuration"))

	rows := [][]string{}

	// completions
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		rows = append(rows, []string{"completions", "zsh", success.Render("✓")})
	} else {
		rows = append(rows, []string{"completions", "", errorStyle.Render("✗")})
	}

	// oh-my-zsh
	omzPath := os.Getenv("HOME") + "/.oh-my-zsh"
	if _, err := os.Stat(omzPath); err == nil {
		rows = append(rows, []string{"oh-my-zsh", "~/.oh-my-zsh", success.Render("✓")})
	} else {
		rows = append(rows, []string{"oh-my-zsh", "", errorStyle.Render("✗")})
	}

	// ssh key
	sshKey := os.Getenv("HOME") + "/.ssh/id_github"
	if _, err := os.Stat(sshKey); err == nil {
		rows = append(rows, []string{"ssh-key", "~/.ssh/id_github", success.Render("✓")})
	} else {
		rows = append(rows, []string{"ssh-key", "", errorStyle.Render("✗")})
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if col == 0 {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).PaddingLeft(1).PaddingRight(1)
			}
			return lipgloss.NewStyle().Foreground(lipgloss.Color("241")).PaddingLeft(1).PaddingRight(1)
		}).
		Rows(rows...)

	fmt.Println("  " + strings.ReplaceAll(t.Render(), "\n", "\n  "))
	fmt.Println()
}

func printPackageStats() {
	fmt.Println(headerStyle.Render("Packages"))

	rows := [][]string{}

	// Homebrew stats
	if commandExists("brew") {
		formulaeOut := getCommandOutput("brew", "list", "--formula", "-1")
		caskOut := getCommandOutput("brew", "list", "--cask", "-1")
		formulaeCount := len(strings.Split(strings.TrimSpace(formulaeOut), "\n"))
		caskCount := len(strings.Split(strings.TrimSpace(caskOut), "\n"))
		if formulaeOut == "" {
			formulaeCount = 0
		}
		if caskOut == "" {
			caskCount = 0
		}
		rows = append(rows, []string{"brew", fmt.Sprintf("%d formulae, %d casks", formulaeCount, caskCount)})
	}

	// Docker stats
	if commandExists("docker") {
		containersOut, _ := exec.Command("docker", "ps", "-aq").Output()
		imagesOut, _ := exec.Command("docker", "images", "-q").Output()
		containerCount := 0
		imageCount := 0
		if len(strings.TrimSpace(string(containersOut))) > 0 {
			containerCount = len(strings.Split(strings.TrimSpace(string(containersOut)), "\n"))
		}
		if len(strings.TrimSpace(string(imagesOut))) > 0 {
			imageCount = len(strings.Split(strings.TrimSpace(string(imagesOut)), "\n"))
		}
		rows = append(rows, []string{"docker", fmt.Sprintf("%d containers, %d images", containerCount, imageCount)})
	}

	// npm global stats
	if commandExists("npm") {
		npmOut := getCommandOutput("npm", "list", "-g", "--depth=0", "--parseable")
		npmLines := strings.Split(strings.TrimSpace(npmOut), "\n")
		count := len(npmLines) - 1
		if count < 0 {
			count = 0
		}
		rows = append(rows, []string{"npm", fmt.Sprintf("%d global packages", count)})
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if col == 0 {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("241")).PaddingLeft(1).PaddingRight(1).Width(10)
			}
			return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).PaddingLeft(1).PaddingRight(1)
		}).
		Rows(rows...)

	fmt.Println("  " + strings.ReplaceAll(t.Render(), "\n", "\n  "))
	fmt.Println()
}

func printDiskUsage() {
	fmt.Println(headerStyle.Render("Disk Usage"))

	rows := [][]string{}
	var hasData bool
	home := os.Getenv("HOME")

	// Docker
	if commandExists("docker") {
		out, _ := exec.Command("docker", "system", "df", "--format", "{{.Size}}").Output()
		dockerLines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(dockerLines) > 0 && dockerLines[0] != "" {
			rows = append(rows, []string{"docker", warning.Render(strings.Join(dockerLines, " + "))})
			hasData = true
		}
	}

	// Homebrew cache
	brewCache := home + "/Library/Caches/Homebrew"
	if size := getDirSize(brewCache); size > 0 {
		rows = append(rows, []string{"homebrew cache", warning.Render(formatBytes(size))})
		hasData = true
	}

	// Multipass
	if commandExists("multipass") {
		multipassData := home + "/Library/Application Support/multipassd"
		if size := getDirSize(multipassData); size > 0 {
			rows = append(rows, []string{"multipass", warning.Render(formatBytes(size))})
			hasData = true
		}
	}

	// npm cache
	npmCache := home + "/.npm"
	if size := getDirSize(npmCache); size > 0 {
		rows = append(rows, []string{"npm cache", warning.Render(formatBytes(size))})
		hasData = true
	}

	// pnpm cache
	pnpmCache := home + "/Library/pnpm"
	if size := getDirSize(pnpmCache); size > 0 {
		rows = append(rows, []string{"pnpm cache", warning.Render(formatBytes(size))})
		hasData = true
	}

	// Trash
	trashPath := home + "/.Trash"
	if size := getDirSize(trashPath); size > 0 {
		rows = append(rows, []string{"trash", warning.Render(formatBytes(size))})
		hasData = true
	}

	if hasData {
		t := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
			StyleFunc(func(row, col int) lipgloss.Style {
				if col == 0 {
					return lipgloss.NewStyle().Foreground(lipgloss.Color("241")).PaddingLeft(1).PaddingRight(1).Width(16)
				}
				return lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)
			}).
			Rows(rows...)

		fmt.Println("  " + strings.ReplaceAll(t.Render(), "\n", "\n  "))
		fmt.Printf("\n %s %s", dimStyle.Render("run"), special.Render("j clean"))
		if commandExists("mo") {
			fmt.Printf(" %s %s", dimStyle.Render("or"), special.Render("mo clean"))
		}
		fmt.Println()
	}
	fmt.Println()
}

// Types and helpers

type toolInfo struct {
	name      string
	version   string
	source    string
	installed bool
	extra     string
	hint      string
}

func getToolInfo(name, cmd string, args []string, source string, parser func(string) string) toolInfo {
	return getToolInfoWithHint(name, cmd, args, source, parser, fmt.Sprintf("brew install %s", cmd))
}

func getToolInfoWithHint(name, cmd string, args []string, source string, parser func(string) string, hint string) toolInfo {
	info := toolInfo{name: name, source: source}

	if _, err := exec.LookPath(cmd); err != nil {
		info.installed = false
		info.hint = hint
		return info
	}

	info.installed = true
	out, err := exec.Command(cmd, args...).Output()
	if err == nil {
		info.version = parser(string(out))
	}

	return info
}

func getDockerInfo() toolInfo {
	info := toolInfo{name: "docker", source: "cask"}

	if _, err := exec.LookPath("docker"); err != nil {
		info.installed = false
		info.hint = "brew install --cask docker"
		return info
	}

	info.installed = true
	out, _ := exec.Command("docker", "--version").Output()
	parts := strings.Split(string(out), " ")
	if len(parts) >= 3 {
		info.version = strings.TrimSuffix(parts[2], ",")
	}

	if err := exec.Command("docker", "info").Run(); err == nil {
		info.extra = "running"
	} else {
		info.extra = "stopped"
	}

	return info
}

func printToolTable(title string, tools []toolInfo) {
	fmt.Println(headerStyle.Render(title))

	rows := [][]string{}
	for _, t := range tools {
		var status string
		if t.installed {
			status = success.Render("✓")
			if t.extra != "" {
				if t.extra == "running" {
					status += " " + success.Render(t.extra)
				} else {
					status += " " + warning.Render(t.extra)
				}
			}
		} else {
			status = errorStyle.Render("✗")
		}

		rows = append(rows, []string{t.name, t.version, t.source, status})
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch col {
			case 0:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).PaddingLeft(1).PaddingRight(1).Width(14)
			case 1:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("241")).PaddingLeft(1).PaddingRight(1).Width(14)
			case 2:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("241")).PaddingLeft(1).PaddingRight(1).Width(8)
			default:
				return lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)
			}
		}).
		Rows(rows...)

	fmt.Println("  " + strings.ReplaceAll(t.Render(), "\n", "\n  "))
	fmt.Println()
}

// Version parsers

func trimVersion(s string) string {
	return strings.TrimSpace(strings.TrimPrefix(s, "v"))
}

func parseBrewVersion(s string) string {
	lines := strings.Split(s, "\n")
	if len(lines) > 0 {
		parts := strings.Split(lines[0], " ")
		if len(parts) >= 2 {
			return parts[1]
		}
	}
	return ""
}

func parseGitVersion(s string) string {
	v := strings.TrimPrefix(strings.TrimSpace(s), "git version ")
	// Truncate Apple Git suffix for cleaner display
	if idx := strings.Index(v, " ("); idx != -1 {
		v = v[:idx]
	}
	return v
}

func parseGoVersion(s string) string {
	// "go version go1.23.4 darwin/arm64" -> "1.23.4"
	parts := strings.Fields(s)
	if len(parts) >= 3 {
		return strings.TrimPrefix(parts[2], "go")
	}
	return ""
}

func parsePythonVersion(s string) string {
	return strings.TrimPrefix(strings.TrimSpace(s), "Python ")
}

func parseTerraformVersion(s string) string {
	lines := strings.Split(s, "\n")
	if len(lines) > 0 {
		parts := strings.Split(lines[0], " ")
		if len(parts) >= 2 {
			return strings.TrimPrefix(parts[1], "v")
		}
	}
	return ""
}

func parseAnsibleVersion(s string) string {
	lines := strings.Split(s, "\n")
	if len(lines) > 0 {
		line := lines[0]
		if strings.Contains(line, "[core") {
			start := strings.Index(line, "[core ")
			end := strings.Index(line, "]")
			if start != -1 && end != -1 {
				return line[start+6 : end]
			}
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			return parts[1]
		}
	}
	return ""
}

func parseKubectlVersion(s string) string {
	for _, line := range strings.Split(s, "\n") {
		if strings.Contains(line, "gitVersion:") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				return strings.TrimSpace(strings.TrimPrefix(parts[1], " v"))
			}
		}
	}
	return ""
}

func parseMultipassVersion(s string) string {
	lines := strings.Split(s, "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 2 {
			return parts[1]
		}
	}
	return ""
}

func parseCodexVersion(s string) string {
	parts := strings.Fields(strings.TrimSpace(s))
	if len(parts) >= 2 {
		return parts[1]
	}
	return strings.TrimSpace(s)
}

func parseMoleVersion(s string) string {
	// "\nMole version 1.14.5\nmacOS: ..." -> "1.14.5"
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(line, "Mole version") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				return parts[2]
			}
		}
	}
	return ""
}

func parseClaudeVersion(s string) string {
	// "2.0.76 (Claude Code)" -> "2.0.76"
	parts := strings.Fields(strings.TrimSpace(s))
	if len(parts) >= 1 {
		return parts[0]
	}
	return ""
}

// Utility functions

func getDirSize(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func getCommandOutput(name string, args ...string) string {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
