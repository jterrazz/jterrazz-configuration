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

// SetupItems is the single source of truth for all setup configurations
var SetupItems = []SetupItem{
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
		Name:        "ssh-key",
		Description: "GitHub SSH key",
		CheckFn: func() (bool, string) {
			sshKey := os.Getenv("HOME") + "/.ssh/id_github"
			if _, err := os.Stat(sshKey); err == nil {
				return true, "~/.ssh/id_github"
			}
			return false, ""
		},
		SetupCmd: "git-ssh",
	},
}
