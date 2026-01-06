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

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99"))

	subSectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

func showStatus() {
	fmt.Println(titleStyle.Render("j status"))
	fmt.Println()

	printSystemInfo()
	printSystemSection()
	printToolsSection()
}

func printSystemInfo() {
	hostname, _ := os.Hostname()
	osInfo := getCommandOutput("uname", "-sr")
	arch := getCommandOutput("uname", "-m")
	user := os.Getenv("USER")
	shell := filepath.Base(os.Getenv("SHELL"))

	fmt.Printf("%s • %s\n", special.Render(osInfo), dimStyle.Render(arch))
	fmt.Printf("%s • %s • %s\n\n", dimStyle.Render(hostname), dimStyle.Render(user), dimStyle.Render(shell))
}

func printSystemSection() {
	fmt.Println(sectionStyle.Render("System"))
	fmt.Println()

	// Setup subsection
	fmt.Println(subSectionStyle.Render("Setup"))
	printSetupTable()

	// Security subsection
	fmt.Println(subSectionStyle.Render("Security"))
	printSecurityTable()

	// Network subsection
	fmt.Println(subSectionStyle.Render("Network"))
	printNetworkTable()

	// Disk Usage subsection
	fmt.Println(subSectionStyle.Render("Disk Usage"))
	printDiskTable()
}

func printSetupTable() {
	rows := [][]string{}

	for _, item := range SetupItems {
		installed, detail := item.CheckFn()
		if installed {
			rows = append(rows, []string{item.Name, detail, success.Render("✓")})
		} else {
			rows = append(rows, []string{item.Name, "", errorStyle.Render("✗")})
		}
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if col == 0 {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).PaddingLeft(1).PaddingRight(1).Width(14)
			}
			return lipgloss.NewStyle().Foreground(lipgloss.Color("241")).PaddingLeft(1).PaddingRight(1)
		}).
		Rows(rows...)

	fmt.Println(t.Render())
	fmt.Println()
}

func printSecurityTable() {
	type securityCheck struct {
		name        string
		description string
		checkFn     func() (ok bool, detail string)
		goodWhen    bool // true = check passes when enabled, false = check passes when disabled
	}

	checks := []securityCheck{
		// System Protection
		{
			name:        "filevault",
			description: "Full disk encryption",
			checkFn: func() (bool, string) {
				out, _ := exec.Command("fdesetup", "status").Output()
				return strings.Contains(string(out), "FileVault is On"), ""
			},
			goodWhen: true,
		},
		{
			name:        "firewall",
			description: "Block incoming connections",
			checkFn: func() (bool, string) {
				out, _ := exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--getglobalstate").Output()
				return strings.Contains(string(out), "enabled"), ""
			},
			goodWhen: true,
		},
		{
			name:        "sip",
			description: "System Integrity Protection",
			checkFn: func() (bool, string) {
				out, _ := exec.Command("csrutil", "status").Output()
				return strings.Contains(string(out), "enabled"), ""
			},
			goodWhen: true,
		},
		{
			name:        "gatekeeper",
			description: "App signature verification",
			checkFn: func() (bool, string) {
				out, _ := exec.Command("spctl", "--status").Output()
				return strings.Contains(string(out), "enabled"), ""
			},
			goodWhen: true,
		},
		// Network
		{
			name:        "remote-login",
			description: "SSH server disabled",
			checkFn: func() (bool, string) {
				out, _ := exec.Command("launchctl", "list").Output()
				sshRunning := strings.Contains(string(out), "com.openssh.sshd")
				return !sshRunning, ""
			},
			goodWhen: true,
		},
		// Keys
		{
			name:        "ssh",
			description: "SSH key for authentication",
			checkFn: func() (bool, string) {
				sshKey := os.Getenv("HOME") + "/.ssh/id_ed25519"
				if _, err := os.Stat(sshKey); err == nil {
					return true, "~/.ssh/id_ed25519"
				}
				return false, ""
			},
			goodWhen: true,
		},
		{
			name:        "gpg",
			description: "GPG key for signing",
			checkFn: func() (bool, string) {
				out, err := exec.Command("gpg", "--list-secret-keys", "--keyid-format", "long").Output()
				if err != nil || len(out) == 0 {
					return false, ""
				}
				return true, "~/.gnupg"
			},
			goodWhen: true,
		},
		// Git config
		{
			name:        "git-email",
			description: "Git commit email",
			checkFn: func() (bool, string) {
				out, _ := exec.Command("git", "config", "--global", "user.email").Output()
				email := strings.TrimSpace(string(out))
				return email == "admin@jterrazz.com", email
			},
			goodWhen: true,
		},
		{
			name:        "git-name",
			description: "Git commit author name",
			checkFn: func() (bool, string) {
				out, _ := exec.Command("git", "config", "--global", "user.name").Output()
				name := strings.TrimSpace(string(out))
				return name != "", name
			},
			goodWhen: true,
		},
		{
			name:        "commit-signing",
			description: "Sign git commits with GPG",
			checkFn: func() (bool, string) {
				out, _ := exec.Command("git", "config", "--global", "commit.gpgsign").Output()
				return strings.TrimSpace(string(out)) == "true", ""
			},
			goodWhen: true,
		},
	}

	rows := [][]string{}
	for _, check := range checks {
		ok, detail := check.checkFn()
		var status string
		if ok == check.goodWhen {
			status = success.Render("✓")
		} else {
			status = warning.Render("!")
		}
		rows = append(rows, []string{check.name, check.description, detail, status})
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch col {
			case 0:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).PaddingLeft(1).PaddingRight(1).Width(16)
			case 1:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("241")).PaddingLeft(1).PaddingRight(1).Width(28)
			case 2:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("86")).PaddingLeft(1).PaddingRight(1)
			default:
				return lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)
			}
		}).
		Rows(rows...)

	fmt.Println(t.Render())
	fmt.Println()
}

func printNetworkTable() {
	rows := [][]string{}

	// WiFi network name
	wifiOut, _ := exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport", "-I").Output()
	for _, line := range strings.Split(string(wifiOut), "\n") {
		if strings.Contains(line, " SSID:") {
			ssid := strings.TrimSpace(strings.TrimPrefix(line, " SSID:"))
			rows = append(rows, []string{"wifi", special.Render(ssid)})
			break
		}
	}

	// VPN status
	vpnOut, _ := exec.Command("scutil", "--nc", "list").Output()
	vpnConnected := strings.Contains(string(vpnOut), "(Connected)")
	if vpnConnected {
		rows = append(rows, []string{"vpn", success.Render("connected")})
	} else {
		rows = append(rows, []string{"vpn", dimStyle.Render("disconnected")})
	}

	// Local IP
	ifconfigOut, _ := exec.Command("ipconfig", "getifaddr", "en0").Output()
	localIP := strings.TrimSpace(string(ifconfigOut))
	if localIP != "" {
		rows = append(rows, []string{"local ip", dimStyle.Render(localIP)})
	}

	// Public IP (with timeout, prefer IPv4)
	publicIP := ""
	cmd := exec.Command("curl", "-s", "--max-time", "2", "-4", "ifconfig.me")
	if out, err := cmd.Output(); err == nil {
		publicIP = strings.TrimSpace(string(out))
	}
	if publicIP != "" {
		rows = append(rows, []string{"public ip", dimStyle.Render(publicIP)})
	}

	// DNS servers (only show valid IPs)
	dnsOut, _ := exec.Command("scutil", "--dns").Output()
	var dnsServers []string
	for _, line := range strings.Split(string(dnsOut), "\n") {
		if strings.Contains(line, "nameserver[") {
			// Format: "  nameserver[0] : 192.168.1.254" or "  nameserver[1] : fd0f:ee:b0::1"
			idx := strings.Index(line, "] : ")
			if idx == -1 {
				continue
			}
			server := strings.TrimSpace(line[idx+4:])
			if server == "" {
				continue
			}
			// Skip localhost
			if server == "127.0.0.1" || server == "::1" {
				continue
			}
			// Avoid duplicates
			found := false
			for _, s := range dnsServers {
				if s == server {
					found = true
					break
				}
			}
			if !found && len(dnsServers) < 2 {
				dnsServers = append(dnsServers, server)
			}
		}
	}
	if len(dnsServers) > 0 {
		rows = append(rows, []string{"dns", dimStyle.Render(strings.Join(dnsServers, ", "))})
	}

	// Listening ports count
	lsofOut, _ := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-P", "-n").Output()
	listenLines := strings.Split(strings.TrimSpace(string(lsofOut)), "\n")
	listenCount := 0
	if len(listenLines) > 1 {
		listenCount = len(listenLines) - 1 // subtract header
	}
	if listenCount > 0 {
		rows = append(rows, []string{"listening", warning.Render(fmt.Sprintf("%d ports", listenCount))})
	} else {
		rows = append(rows, []string{"listening", dimStyle.Render("0 ports")})
	}

	// Active connections count
	netstatOut, _ := exec.Command("netstat", "-an").Output()
	establishedCount := 0
	for _, line := range strings.Split(string(netstatOut), "\n") {
		if strings.Contains(line, "ESTABLISHED") {
			establishedCount++
		}
	}
	rows = append(rows, []string{"connections", dimStyle.Render(fmt.Sprintf("%d active", establishedCount))})

	if len(rows) > 0 {
		t := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
			StyleFunc(func(row, col int) lipgloss.Style {
				if col == 0 {
					return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).PaddingLeft(1).PaddingRight(1).Width(14)
				}
				return lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)
			}).
			Rows(rows...)

		fmt.Println(t.Render())
	}
	fmt.Println()
}

func printDiskTable() {
	home := os.Getenv("HOME")

	// Main directories section
	mainRows := [][]string{}

	// Developer folder
	developerPath := home + "/Developer"
	if size := getDirSize(developerPath); size > 0 {
		mainRows = append(mainRows, []string{"~/Developer", special.Render(formatBytes(size))})
	}

	// Applications
	appsPath := "/Applications"
	if size := getDirSize(appsPath); size > 0 {
		mainRows = append(mainRows, []string{"/Applications", dimStyle.Render(formatBytes(size))})
	}

	// Documents
	documentsPath := home + "/Documents"
	if size := getDirSize(documentsPath); size > 0 {
		mainRows = append(mainRows, []string{"~/Documents", dimStyle.Render(formatBytes(size))})
	}

	// Downloads
	downloadsPath := home + "/Downloads"
	if size := getDirSize(downloadsPath); size > 0 {
		mainRows = append(mainRows, []string{"~/Downloads", warning.Render(formatBytes(size))})
	}

	if len(mainRows) > 0 {
		t := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
			StyleFunc(func(row, col int) lipgloss.Style {
				if col == 0 {
					return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).PaddingLeft(1).PaddingRight(1).Width(18)
				}
				return lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)
			}).
			Rows(mainRows...)
		fmt.Println(t.Render())
		fmt.Println()
	}

	// Caches & cleanable section
	fmt.Println(subSectionStyle.Render("Caches & Cleanable"))
	cacheRows := [][]string{}

	// Docker
	if commandExists("docker") {
		out, _ := exec.Command("docker", "system", "df", "--format", "{{.Size}}").Output()
		dockerLines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(dockerLines) > 0 && dockerLines[0] != "" {
			cacheRows = append(cacheRows, []string{"docker", warning.Render(strings.Join(dockerLines, " + "))})
		}
	}

	// Xcode derived data
	xcodeDerivedData := home + "/Library/Developer/Xcode/DerivedData"
	if size := getDirSize(xcodeDerivedData); size > 0 {
		cacheRows = append(cacheRows, []string{"xcode derived", warning.Render(formatBytes(size))})
	}

	// Xcode archives
	xcodeArchives := home + "/Library/Developer/Xcode/Archives"
	if size := getDirSize(xcodeArchives); size > 0 {
		cacheRows = append(cacheRows, []string{"xcode archives", warning.Render(formatBytes(size))})
	}

	// iOS Device Support
	iosDeviceSupport := home + "/Library/Developer/Xcode/iOS DeviceSupport"
	if size := getDirSize(iosDeviceSupport); size > 0 {
		cacheRows = append(cacheRows, []string{"ios device support", warning.Render(formatBytes(size))})
	}

	// CocoaPods cache
	cocoapodsCache := home + "/Library/Caches/CocoaPods"
	if size := getDirSize(cocoapodsCache); size > 0 {
		cacheRows = append(cacheRows, []string{"cocoapods cache", warning.Render(formatBytes(size))})
	}

	// Homebrew cache
	brewCache := home + "/Library/Caches/Homebrew"
	if size := getDirSize(brewCache); size > 0 {
		cacheRows = append(cacheRows, []string{"homebrew cache", warning.Render(formatBytes(size))})
	}

	// Multipass
	if commandExists("multipass") {
		multipassData := home + "/Library/Application Support/multipassd"
		if size := getDirSize(multipassData); size > 0 {
			cacheRows = append(cacheRows, []string{"multipass", warning.Render(formatBytes(size))})
		}
	}

	// npm cache
	npmCache := home + "/.npm"
	if size := getDirSize(npmCache); size > 0 {
		cacheRows = append(cacheRows, []string{"npm cache", warning.Render(formatBytes(size))})
	}

	// pnpm cache
	pnpmCache := home + "/Library/pnpm"
	if size := getDirSize(pnpmCache); size > 0 {
		cacheRows = append(cacheRows, []string{"pnpm cache", warning.Render(formatBytes(size))})
	}

	// Yarn cache
	yarnCache := home + "/Library/Caches/Yarn"
	if size := getDirSize(yarnCache); size > 0 {
		cacheRows = append(cacheRows, []string{"yarn cache", warning.Render(formatBytes(size))})
	}

	// Go module cache
	goCache := home + "/go/pkg/mod"
	if size := getDirSize(goCache); size > 0 {
		cacheRows = append(cacheRows, []string{"go modules", warning.Render(formatBytes(size))})
	}

	// Gradle cache
	gradleCache := home + "/.gradle/caches"
	if size := getDirSize(gradleCache); size > 0 {
		cacheRows = append(cacheRows, []string{"gradle cache", warning.Render(formatBytes(size))})
	}

	// System logs
	systemLogs := "/var/log"
	if size := getDirSize(systemLogs); size > 0 {
		cacheRows = append(cacheRows, []string{"system logs", dimStyle.Render(formatBytes(size))})
	}

	// User logs
	userLogs := home + "/Library/Logs"
	if size := getDirSize(userLogs); size > 0 {
		cacheRows = append(cacheRows, []string{"user logs", warning.Render(formatBytes(size))})
	}

	// Trash
	trashPath := home + "/.Trash"
	if size := getDirSize(trashPath); size > 0 {
		cacheRows = append(cacheRows, []string{"trash", warning.Render(formatBytes(size))})
	}

	if len(cacheRows) > 0 {
		t := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
			StyleFunc(func(row, col int) lipgloss.Style {
				if col == 0 {
					return lipgloss.NewStyle().Foreground(lipgloss.Color("241")).PaddingLeft(1).PaddingRight(1).Width(20)
				}
				return lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)
			}).
			Rows(cacheRows...)

		fmt.Println(t.Render())
		fmt.Printf("%s %s", dimStyle.Render("run"), special.Render("j clean"))
		if commandExists("mo") {
			fmt.Printf(" %s %s", dimStyle.Render("or"), special.Render("mo clean"))
		}
		fmt.Println()
	}
	fmt.Println()
}

func printToolsSection() {
	fmt.Println(sectionStyle.Render("Tools"))
	fmt.Println()

	categories := []PackageCategory{
		CategoryPackageManager,
		CategoryDevelopment,
		CategoryInfrastructure,
		CategoryAI,
		CategorySystemTools,
	}

	for _, category := range categories {
		packages := GetPackagesByCategory(category)
		if len(packages) == 0 {
			continue
		}

		fmt.Println(subSectionStyle.Render(string(category)))
		printPackageTable(packages)
	}
}

func printPackageTable(packages []Package) {
	rows := [][]string{}
	for _, pkg := range packages {
		installed, version, extra := CheckPackage(pkg)

		var status string
		if installed {
			status = success.Render("✓")
			if extra != "" {
				if extra == "running" {
					status += " " + success.Render(extra)
				} else {
					status += " " + warning.Render(extra)
				}
			}
		} else {
			status = errorStyle.Render("✗")
		}

		rows = append(rows, []string{pkg.Name, version, pkg.Method.String(), status})
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

	fmt.Println(t.Render())
	fmt.Println()
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
