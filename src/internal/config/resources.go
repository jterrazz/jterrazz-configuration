package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jterrazz/jterrazz-cli/internal/tool"
)

// ResourceCheck represents a system resource check (network, disk, cache)
type ResourceCheck struct {
	Name    string
	CheckFn func() ResourceResult
}

// ResourceResult holds the result of a resource check
type ResourceResult struct {
	Value     string // The value to display (e.g., IP address, size)
	Style     string // "success", "warning", "muted", "special"
	Available bool   // Whether this resource is available/relevant
}

// NetworkChecks is the list of network resource checks
var NetworkChecks = []ResourceCheck{
	{
		Name: "wifi",
		CheckFn: func() ResourceResult {
			out, _ := exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport", "-I").Output()
			for _, line := range strings.Split(string(out), "\n") {
				if strings.Contains(line, " SSID:") {
					ssid := strings.TrimSpace(strings.TrimPrefix(line, " SSID:"))
					return ResourceResult{Value: ssid, Style: "muted", Available: true}
				}
			}
			return ResourceResult{Available: false}
		},
	},
	{
		Name: "vpn",
		CheckFn: func() ResourceResult {
			out, _ := exec.Command("scutil", "--nc", "list").Output()
			if strings.Contains(string(out), "(Connected)") {
				return ResourceResult{Value: "connected", Style: "success", Available: true}
			}
			return ResourceResult{Value: "disconnected", Style: "muted", Available: true}
		},
	},
	{
		Name: "local ip",
		CheckFn: func() ResourceResult {
			out, _ := exec.Command("ipconfig", "getifaddr", "en0").Output()
			ip := strings.TrimSpace(string(out))
			if ip != "" {
				return ResourceResult{Value: ip, Style: "muted", Available: true}
			}
			return ResourceResult{Available: false}
		},
	},
	{
		Name: "public ip",
		CheckFn: func() ResourceResult {
			cmd := exec.Command("curl", "-s", "--max-time", "2", "-4", "ifconfig.me")
			out, err := cmd.Output()
			if err == nil {
				ip := strings.TrimSpace(string(out))
				if ip != "" {
					return ResourceResult{Value: ip, Style: "muted", Available: true}
				}
			}
			return ResourceResult{Available: false}
		},
	},
	{
		Name: "dns",
		CheckFn: func() ResourceResult {
			out, _ := exec.Command("scutil", "--dns").Output()
			var servers []string
			for _, line := range strings.Split(string(out), "\n") {
				if strings.Contains(line, "nameserver[") {
					idx := strings.Index(line, "] : ")
					if idx == -1 {
						continue
					}
					server := strings.TrimSpace(line[idx+4:])
					if server == "" || server == "127.0.0.1" || server == "::1" {
						continue
					}
					found := false
					for _, s := range servers {
						if s == server {
							found = true
							break
						}
					}
					if !found && len(servers) < 2 {
						servers = append(servers, server)
					}
				}
			}
			if len(servers) > 0 {
				return ResourceResult{Value: strings.Join(servers, ", "), Style: "muted", Available: true}
			}
			return ResourceResult{Available: false}
		},
	},
	{
		Name: "listening",
		CheckFn: func() ResourceResult {
			out, _ := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-P", "-n").Output()
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			count := 0
			if len(lines) > 1 {
				count = len(lines) - 1 // subtract header
			}
			if count > 0 {
				return ResourceResult{Value: fmt.Sprintf("%d ports", count), Style: "muted", Available: true}
			}
			return ResourceResult{Value: "0 ports", Style: "muted", Available: true}
		},
	},
	{
		Name: "connections",
		CheckFn: func() ResourceResult {
			out, _ := exec.Command("netstat", "-an").Output()
			count := 0
			for _, line := range strings.Split(string(out), "\n") {
				if strings.Contains(line, "ESTABLISHED") {
					count++
				}
			}
			return ResourceResult{Value: fmt.Sprintf("%d active", count), Style: "muted", Available: true}
		},
	},
}

// DiskCheck represents a disk usage check
type DiskCheck struct {
	Name    string
	Path    string                // Path to check (supports ~ expansion)
	Style   string                // Default style for this check
	CheckFn func() ResourceResult // Custom check (overrides Path)
}

// MainDiskChecks shows primary directories
var MainDiskChecks = []DiskCheck{
	{Name: "~/Developer", Path: "~/Developer", Style: "muted"},
	{Name: "/Applications", Path: "/Applications", Style: "muted"},
	{Name: "~/Documents", Path: "~/Documents", Style: "muted"},
	{Name: "~/Downloads", Path: "~/Downloads", Style: "muted"},
}

// CacheChecks shows cleanable caches
var CacheChecks = []DiskCheck{
	{
		Name: "docker",
		CheckFn: func() ResourceResult {
			if !CommandExists("docker") {
				return ResourceResult{Available: false}
			}
			out, _ := exec.Command("docker", "system", "df", "--format", "{{.Size}}").Output()
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			if len(lines) > 0 && lines[0] != "" {
				return ResourceResult{Value: strings.Join(lines, " + "), Style: "muted", Available: true}
			}
			return ResourceResult{Available: false}
		},
	},
	{Name: "xcode derived", Path: "~/Library/Developer/Xcode/DerivedData", Style: "muted"},
	{Name: "xcode archives", Path: "~/Library/Developer/Xcode/Archives", Style: "muted"},
	{Name: "ios device support", Path: "~/Library/Developer/Xcode/iOS DeviceSupport", Style: "muted"},
	{Name: "cocoapods cache", Path: "~/Library/Caches/CocoaPods", Style: "muted"},
	{Name: "homebrew cache", Path: "~/Library/Caches/Homebrew", Style: "muted"},
	{
		Name: "multipass",
		CheckFn: func() ResourceResult {
			if !CommandExists("multipass") {
				return ResourceResult{Available: false}
			}
			path := expandHome("~/Library/Application Support/multipassd")
			if size := GetDirSize(path); size > 0 {
				return ResourceResult{Value: tool.FormatBytes(size), Style: "muted", Available: true}
			}
			return ResourceResult{Available: false}
		},
	},
	{Name: "npm cache", Path: "~/.npm", Style: "muted"},
	{Name: "pnpm cache", Path: "~/Library/pnpm", Style: "muted"},
	{Name: "yarn cache", Path: "~/Library/Caches/Yarn", Style: "muted"},
	{Name: "go modules", Path: "~/go/pkg/mod", Style: "muted"},
	{Name: "gradle cache", Path: "~/.gradle/caches", Style: "muted"},
	{Name: "system logs", Path: "/var/log", Style: "muted"},
	{Name: "user logs", Path: "~/Library/Logs", Style: "muted"},
	{Name: "trash", Path: "~/.Trash", Style: "muted"},
}

// CheckDisk checks a disk path and returns the result
func (d DiskCheck) Check() ResourceResult {
	if d.CheckFn != nil {
		return d.CheckFn()
	}

	path := expandHome(d.Path)
	if size := GetDirSize(path); size > 0 {
		return ResourceResult{Value: tool.FormatBytes(size), Style: d.Style, Available: true}
	}
	return ResourceResult{Available: false}
}

// expandHome expands ~ to the user's home directory
func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(os.Getenv("HOME"), path[2:])
	}
	return path
}
