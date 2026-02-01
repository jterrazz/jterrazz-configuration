package commands

import (
	"os"
	"path/filepath"

	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/tool"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
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

func showStatus() {
	ui.PrintTitle("j status")
	ui.PrintEmpty()

	printSystemInfo()
	printSystemSection()
	printToolsSection()
	printResourcesSection()
}

// =============================================================================
// System Info
// =============================================================================

func printSystemInfo() {
	hostname, _ := os.Hostname()
	osInfo := tool.GetCommandOutput("uname", "-sr")
	arch := tool.GetCommandOutput("uname", "-m")
	user := os.Getenv("USER")
	shell := filepath.Base(os.Getenv("SHELL"))

	ui.PrintSystemHeader(osInfo, arch, hostname, user, shell)
}

// =============================================================================
// System Section
// =============================================================================

func printSystemSection() {
	ui.PrintSection("System")
	ui.PrintEmpty()

	ui.PrintSubSection("Setup")
	printSetupTable()

	ui.PrintSubSection("macOS Security")
	printSecurityTable()

	ui.PrintSubSection("Identity")
	printIdentityTable()
}

func printSetupTable() {
	rows := [][]string{}
	for _, script := range config.Scripts {
		// Skip one-shot scripts (no CheckFn means no checkable state)
		if script.CheckFn == nil {
			continue
		}
		result := config.CheckScript(script)
		rows = append(rows, []string{script.Name, result.Detail, ui.RenderStatusIcon(result.Installed)})
	}
	ui.PrintTable(rows, ui.StatusTableColumns)
	ui.PrintEmpty()
}

func printSecurityTable() {
	rows := [][]string{}
	for _, check := range config.SecurityChecks {
		result := check.CheckFn()
		status := ui.RenderStatusIconWithCondition(result.Installed, check.GoodWhen)
		rows = append(rows, []string{check.Name, check.Description, result.Detail, status})
	}
	ui.PrintTable(rows, ui.CheckTableColumns)
	ui.PrintEmpty()
}

func printIdentityTable() {
	rows := [][]string{}
	for _, check := range config.IdentityChecks {
		result := check.CheckFn()
		status := ui.RenderStatusIconWithCondition(result.Installed, check.GoodWhen)
		rows = append(rows, []string{check.Name, check.Description, result.Detail, status})
	}
	ui.PrintTable(rows, ui.CheckTableColumns)
	ui.PrintEmpty()
}

// =============================================================================
// Tools Section
// =============================================================================

func printToolsSection() {
	ui.PrintSection("Tools")
	ui.PrintEmpty()

	for _, category := range config.ToolCategories {
		tools := config.GetToolsByCategory(category)
		if len(tools) == 0 {
			continue
		}

		ui.PrintSubSection(string(category))
		rows := [][]string{}
		for _, t := range tools {
			result := t.Check()
			status := ui.RenderStatusIcon(result.Installed)
			if result.Installed && result.Status != "" {
				if result.Status == "running" {
					status += " " + ui.RenderSuccess(result.Status)
				} else {
					status += " " + ui.RenderWarning(result.Status)
				}
			}
			rows = append(rows, []string{t.Name, result.Version, t.Method.String(), status})
		}
		ui.PrintTable(rows, ui.ToolTableColumns)
		ui.PrintEmpty()
	}
}

// =============================================================================
// Resources Section
// =============================================================================

func printResourcesSection() {
	ui.PrintSection("Resources")
	ui.PrintEmpty()

	// Network
	ui.PrintSubSection("Network")
	rows := [][]string{}
	for _, check := range config.NetworkChecks {
		result := check.CheckFn()
		if result.Available {
			rows = append(rows, []string{check.Name, styleValue(result.Value, result.Style)})
		}
	}
	if len(rows) > 0 {
		ui.PrintTable(rows, ui.ResourceTableColumns)
	}
	ui.PrintEmpty()

	// Disk Usage
	ui.PrintSubSection("Disk Usage")
	rows = [][]string{}
	for _, check := range config.MainDiskChecks {
		result := check.Check()
		if result.Available {
			rows = append(rows, []string{check.Name, styleValue(result.Value, result.Style)})
		}
	}
	if len(rows) > 0 {
		ui.PrintTable(rows, ui.DiskTableColumns)
		ui.PrintEmpty()
	}

	// Caches & Cleanable
	ui.PrintSubSection("Caches & Cleanable")
	rows = [][]string{}
	for _, check := range config.CacheChecks {
		result := check.Check()
		if result.Available {
			rows = append(rows, []string{check.Name, styleValue(result.Value, result.Style)})
		}
	}
	if len(rows) > 0 {
		ui.PrintTable(rows, ui.CacheTableColumns)
		if config.CommandExists("mo") {
			ui.PrintHint("run ", "j clean", " or ", "mo clean")
		} else {
			ui.PrintHint("run ", "j clean")
		}
	}
	ui.PrintEmpty()
}

func styleValue(value, style string) string {
	switch style {
	case "success":
		return ui.RenderSuccess(value)
	case "warning":
		return ui.RenderWarning(value)
	case "special":
		return ui.RenderSpecial(value)
	default:
		return ui.RenderMuted(value)
	}
}
