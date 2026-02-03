package commands

import (
	"github.com/jterrazz/jterrazz-cli/internal/domain/skill"
	"github.com/jterrazz/jterrazz-cli/internal/presentation/components/tui"
	"github.com/jterrazz/jterrazz-cli/internal/presentation/print"
	setupview "github.com/jterrazz/jterrazz-cli/internal/presentation/views/setup"
)

func runSkillsUI() {
	if !skill.IsInstalled() {
		print.Error("skills CLI not installed. Run: npm install -g skills")
		return
	}

	setupview.InitSkillsState()
	tui.RunOrExit(setupview.SkillsConfig())
}
