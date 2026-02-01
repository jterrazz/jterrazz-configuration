package config

import (
	"os/exec"
	"strings"
)

// SecurityCheck represents a system security verification
type SecurityCheck struct {
	Name        string
	Description string
	CheckFn     func() CheckResult
	GoodWhen    bool // true = check passes when Installed=true, false = check passes when Installed=false
}

// SecurityChecks is the list of macOS security checks
var SecurityChecks = []SecurityCheck{
	{
		Name:        "filevault",
		Description: "Full disk encryption",
		CheckFn: func() CheckResult {
			out, _ := exec.Command("fdesetup", "status").Output()
			return CheckResult{Installed: strings.Contains(string(out), "FileVault is On")}
		},
		GoodWhen: true,
	},
	{
		Name:        "firewall",
		Description: "Block incoming connections",
		CheckFn: func() CheckResult {
			out, _ := exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--getglobalstate").Output()
			return CheckResult{Installed: strings.Contains(string(out), "enabled")}
		},
		GoodWhen: true,
	},
	{
		Name:        "sip",
		Description: "System Integrity Protection",
		CheckFn: func() CheckResult {
			out, _ := exec.Command("csrutil", "status").Output()
			return CheckResult{Installed: strings.Contains(string(out), "enabled")}
		},
		GoodWhen: true,
	},
	{
		Name:        "gatekeeper",
		Description: "App signature verification",
		CheckFn: func() CheckResult {
			out, _ := exec.Command("spctl", "--status").Output()
			return CheckResult{Installed: strings.Contains(string(out), "enabled")}
		},
		GoodWhen: true,
	},
	{
		Name:        "remote-login",
		Description: "SSH server disabled",
		CheckFn: func() CheckResult {
			out, _ := exec.Command("launchctl", "list").Output()
			sshRunning := strings.Contains(string(out), "com.openssh.sshd")
			return CheckResult{Installed: !sshRunning}
		},
		GoodWhen: true,
	},
}
