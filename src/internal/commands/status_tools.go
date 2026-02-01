package commands

import (
	"fmt"

	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
)

func printToolsSection() {
	fmt.Println(ui.SectionStyle.Render("Tools"))
	fmt.Println()

	categories := []config.ToolCategory{
		config.CategoryPackageManager,
		config.CategoryLanguages,
		config.CategoryInfrastructure,
		config.CategoryAI,
		config.CategoryApps,
		config.CategorySystemTools,
	}

	for _, category := range categories {
		packages := config.GetToolsByCategory(category)
		if len(packages) == 0 {
			continue
		}

		fmt.Println(ui.SubSectionStyle.Render(string(category)))
		printToolTable(packages)
	}
}

func printToolTable(tools []config.Tool) {
	rows := [][]string{}
	for _, tool := range tools {
		result := tool.Check()

		var status string
		if result.Installed {
			status = ui.SuccessStyle.Render("✓")
			if result.Status != "" {
				if result.Status == "running" {
					status += " " + ui.SuccessStyle.Render(result.Status)
				} else {
					status += " " + ui.WarningStyle.Render(result.Status)
				}
			}
		} else {
			status = ui.DangerStyle.Render("✗")
		}

		rows = append(rows, []string{tool.Name, result.Version, tool.Method.String(), status})
	}

	fmt.Println(ui.RenderTable(rows, ui.ToolTableColumns))
	fmt.Println()
}
