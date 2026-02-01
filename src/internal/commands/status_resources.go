package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/tool"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
)

func printResourcesSection() {
	fmt.Println(ui.SectionStyle.Render("Resources"))
	fmt.Println()

	// Network subsection
	fmt.Println(ui.SubSectionStyle.Render("Network"))
	printNetworkTable()

	// Disk Usage subsection
	fmt.Println(ui.SubSectionStyle.Render("Disk Usage"))
	printDiskTable()
}

func printNetworkTable() {
	rows := [][]string{}

	// WiFi network name
	wifiOut, _ := exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport", "-I").Output()
	for _, line := range strings.Split(string(wifiOut), "\n") {
		if strings.Contains(line, " SSID:") {
			ssid := strings.TrimSpace(strings.TrimPrefix(line, " SSID:"))
			rows = append(rows, []string{"wifi", ui.SpecialStyle.Render(ssid)})
			break
		}
	}

	// VPN status
	vpnOut, _ := exec.Command("scutil", "--nc", "list").Output()
	vpnConnected := strings.Contains(string(vpnOut), "(Connected)")
	if vpnConnected {
		rows = append(rows, []string{"vpn", ui.SuccessStyle.Render("connected")})
	} else {
		rows = append(rows, []string{"vpn", ui.MutedStyle.Render("disconnected")})
	}

	// Local IP
	ifconfigOut, _ := exec.Command("ipconfig", "getifaddr", "en0").Output()
	localIP := strings.TrimSpace(string(ifconfigOut))
	if localIP != "" {
		rows = append(rows, []string{"local ip", ui.MutedStyle.Render(localIP)})
	}

	// Public IP (with timeout, prefer IPv4)
	publicIP := ""
	cmd := exec.Command("curl", "-s", "--max-time", "2", "-4", "ifconfig.me")
	if out, err := cmd.Output(); err == nil {
		publicIP = strings.TrimSpace(string(out))
	}
	if publicIP != "" {
		rows = append(rows, []string{"public ip", ui.MutedStyle.Render(publicIP)})
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
		rows = append(rows, []string{"dns", ui.MutedStyle.Render(strings.Join(dnsServers, ", "))})
	}

	// Listening ports count
	lsofOut, _ := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-P", "-n").Output()
	listenLines := strings.Split(strings.TrimSpace(string(lsofOut)), "\n")
	listenCount := 0
	if len(listenLines) > 1 {
		listenCount = len(listenLines) - 1 // subtract header
	}
	if listenCount > 0 {
		rows = append(rows, []string{"listening", ui.WarningStyle.Render(fmt.Sprintf("%d ports", listenCount))})
	} else {
		rows = append(rows, []string{"listening", ui.MutedStyle.Render("0 ports")})
	}

	// Active connections count
	netstatOut, _ := exec.Command("netstat", "-an").Output()
	establishedCount := 0
	for _, line := range strings.Split(string(netstatOut), "\n") {
		if strings.Contains(line, "ESTABLISHED") {
			establishedCount++
		}
	}
	rows = append(rows, []string{"connections", ui.MutedStyle.Render(fmt.Sprintf("%d active", establishedCount))})

	if len(rows) > 0 {
		fmt.Println(ui.RenderTable(rows, ui.ResourceTableColumns))
	}
	fmt.Println()
}

func printDiskTable() {
	home := os.Getenv("HOME")

	// Main directories section
	mainRows := [][]string{}

	// Developer folder
	developerPath := home + "/Developer"
	if size := config.GetDirSize(developerPath); size > 0 {
		mainRows = append(mainRows, []string{"~/Developer", ui.SpecialStyle.Render(tool.FormatBytes(size))})
	}

	// Applications
	appsPath := "/Applications"
	if size := config.GetDirSize(appsPath); size > 0 {
		mainRows = append(mainRows, []string{"/Applications", ui.MutedStyle.Render(tool.FormatBytes(size))})
	}

	// Documents
	documentsPath := home + "/Documents"
	if size := config.GetDirSize(documentsPath); size > 0 {
		mainRows = append(mainRows, []string{"~/Documents", ui.MutedStyle.Render(tool.FormatBytes(size))})
	}

	// Downloads
	downloadsPath := home + "/Downloads"
	if size := config.GetDirSize(downloadsPath); size > 0 {
		mainRows = append(mainRows, []string{"~/Downloads", ui.WarningStyle.Render(tool.FormatBytes(size))})
	}

	if len(mainRows) > 0 {
		fmt.Println(ui.RenderTable(mainRows, ui.DiskTableColumns))
		fmt.Println()
	}

	// Caches & cleanable section
	fmt.Println(ui.SubSectionStyle.Render("Caches & Cleanable"))
	cacheRows := [][]string{}

	// Docker
	if config.CommandExists("docker") {
		out, _ := exec.Command("docker", "system", "df", "--format", "{{.Size}}").Output()
		dockerLines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(dockerLines) > 0 && dockerLines[0] != "" {
			cacheRows = append(cacheRows, []string{"docker", ui.WarningStyle.Render(strings.Join(dockerLines, " + "))})
		}
	}

	// Xcode derived data
	xcodeDerivedData := home + "/Library/Developer/Xcode/DerivedData"
	if size := config.GetDirSize(xcodeDerivedData); size > 0 {
		cacheRows = append(cacheRows, []string{"xcode derived", ui.WarningStyle.Render(tool.FormatBytes(size))})
	}

	// Xcode archives
	xcodeArchives := home + "/Library/Developer/Xcode/Archives"
	if size := config.GetDirSize(xcodeArchives); size > 0 {
		cacheRows = append(cacheRows, []string{"xcode archives", ui.WarningStyle.Render(tool.FormatBytes(size))})
	}

	// iOS Device Support
	iosDeviceSupport := home + "/Library/Developer/Xcode/iOS DeviceSupport"
	if size := config.GetDirSize(iosDeviceSupport); size > 0 {
		cacheRows = append(cacheRows, []string{"ios device support", ui.WarningStyle.Render(tool.FormatBytes(size))})
	}

	// CocoaPods cache
	cocoapodsCache := home + "/Library/Caches/CocoaPods"
	if size := config.GetDirSize(cocoapodsCache); size > 0 {
		cacheRows = append(cacheRows, []string{"cocoapods cache", ui.WarningStyle.Render(tool.FormatBytes(size))})
	}

	// Homebrew cache
	brewCache := home + "/Library/Caches/Homebrew"
	if size := config.GetDirSize(brewCache); size > 0 {
		cacheRows = append(cacheRows, []string{"homebrew cache", ui.WarningStyle.Render(tool.FormatBytes(size))})
	}

	// Multipass
	if config.CommandExists("multipass") {
		multipassData := home + "/Library/Application Support/multipassd"
		if size := config.GetDirSize(multipassData); size > 0 {
			cacheRows = append(cacheRows, []string{"multipass", ui.WarningStyle.Render(tool.FormatBytes(size))})
		}
	}

	// npm cache
	npmCache := home + "/.npm"
	if size := config.GetDirSize(npmCache); size > 0 {
		cacheRows = append(cacheRows, []string{"npm cache", ui.WarningStyle.Render(tool.FormatBytes(size))})
	}

	// pnpm cache
	pnpmCache := home + "/Library/pnpm"
	if size := config.GetDirSize(pnpmCache); size > 0 {
		cacheRows = append(cacheRows, []string{"pnpm cache", ui.WarningStyle.Render(tool.FormatBytes(size))})
	}

	// Yarn cache
	yarnCache := home + "/Library/Caches/Yarn"
	if size := config.GetDirSize(yarnCache); size > 0 {
		cacheRows = append(cacheRows, []string{"yarn cache", ui.WarningStyle.Render(tool.FormatBytes(size))})
	}

	// Go module cache
	goCache := home + "/go/pkg/mod"
	if size := config.GetDirSize(goCache); size > 0 {
		cacheRows = append(cacheRows, []string{"go modules", ui.WarningStyle.Render(tool.FormatBytes(size))})
	}

	// Gradle cache
	gradleCache := home + "/.gradle/caches"
	if size := config.GetDirSize(gradleCache); size > 0 {
		cacheRows = append(cacheRows, []string{"gradle cache", ui.WarningStyle.Render(tool.FormatBytes(size))})
	}

	// System logs
	systemLogs := "/var/log"
	if size := config.GetDirSize(systemLogs); size > 0 {
		cacheRows = append(cacheRows, []string{"system logs", ui.MutedStyle.Render(tool.FormatBytes(size))})
	}

	// User logs
	userLogs := home + "/Library/Logs"
	if size := config.GetDirSize(userLogs); size > 0 {
		cacheRows = append(cacheRows, []string{"user logs", ui.WarningStyle.Render(tool.FormatBytes(size))})
	}

	// Trash
	trashPath := home + "/.Trash"
	if size := config.GetDirSize(trashPath); size > 0 {
		cacheRows = append(cacheRows, []string{"trash", ui.WarningStyle.Render(tool.FormatBytes(size))})
	}

	if len(cacheRows) > 0 {
		fmt.Println(ui.RenderTable(cacheRows, ui.CacheTableColumns))
		fmt.Printf("%s %s", ui.MutedStyle.Render("run"), ui.SpecialStyle.Render("j clean"))
		if config.CommandExists("mo") {
			fmt.Printf(" %s %s", ui.MutedStyle.Render("or"), ui.SpecialStyle.Render("mo clean"))
		}
		fmt.Println()
	}
	fmt.Println()
}
