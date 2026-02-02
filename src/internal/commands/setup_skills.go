package commands

import (
	"github.com/jterrazz/jterrazz-cli/internal/skill"
	"github.com/jterrazz/jterrazz-cli/internal/ui/components/tui"
	"github.com/jterrazz/jterrazz-cli/internal/ui/print"
	setupview "github.com/jterrazz/jterrazz-cli/internal/ui/views/setup"
)

func runSkillsUI() {
	if !skill.IsInstalled() {
		print.Error("skills CLI not installed. Run: npm install -g skills")
		return
	}

	setupview.InitSkillsState()
	tui.RunOrExit(setupview.SkillsConfig())
}
