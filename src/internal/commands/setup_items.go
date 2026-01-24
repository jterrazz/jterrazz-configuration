package commands

import (
	"os"
)

// SetupItem represents a configuration item that can be checked and setup
type SetupItem struct {
	Name        string
	Description string
	CheckFn     func() (installed bool, detail string)
	SetupCmd    string // The j setup subcommand to run
}

// SetupItems is the single source of truth for all setup configurations (alphabetical)
var SetupItems = []SetupItem{
	{
		Name:        "ghostty",
		Description: "Ghostty terminal config",
		CheckFn: func() (bool, string) {
			configPath := os.Getenv("HOME") + "/.config/ghostty/config"
			if _, err := os.Stat(configPath); err == nil {
				return true, "~/.config/ghostty/config"
			}
			return false, ""
		},
		SetupCmd: "ghostty",
	},
	{
		Name:        "hushlogin",
		Description: "Silence terminal login message",
		CheckFn: func() (bool, string) {
			hushPath := os.Getenv("HOME") + "/.hushlogin"
			if _, err := os.Stat(hushPath); err == nil {
				return true, "~/.hushlogin"
			}
			return false, ""
		},
		SetupCmd: "hushlogin",
	},
	{
		Name:        "oh-my-zsh",
		Description: "Oh My Zsh shell configuration",
		CheckFn: func() (bool, string) {
			omzPath := os.Getenv("HOME") + "/.oh-my-zsh"
			if _, err := os.Stat(omzPath); err == nil {
				return true, "~/.oh-my-zsh"
			}
			return false, ""
		},
		SetupCmd: "ohmyzsh",
	},
	{
		Name:        "zed",
		Description: "Zed editor configuration",
		CheckFn: func() (bool, string) {
			configPath := os.Getenv("HOME") + "/.config/zed/settings.json"
			if _, err := os.Stat(configPath); err == nil {
				return true, "~/.config/zed/settings.json"
			}
			return false, ""
		},
		SetupCmd: "zed",
	},
}
