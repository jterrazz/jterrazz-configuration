package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/system"
	"github.com/jterrazz/jterrazz-cli/internal/ui"
)

func printSystemInfo() {
	hostname, _ := os.Hostname()
	osInfo := system.GetCommandOutput("uname", "-sr")
	arch := system.GetCommandOutput("uname", "-m")
	user := os.Getenv("USER")
	shell := filepath.Base(os.Getenv("SHELL"))

	fmt.Printf("%s • %s\n", ui.SpecialStyle.Render(osInfo), ui.MutedStyle.Render(arch))
	fmt.Printf("%s • %s • %s\n\n", ui.MutedStyle.Render(hostname), ui.MutedStyle.Render(user), ui.MutedStyle.Render(shell))
}

func printSystemSection() {
	fmt.Println(ui.SectionStyle.Render("System"))
	fmt.Println()

	// Setup subsection
	fmt.Println(ui.SubSectionStyle.Render("Setup"))
	printSetupTable()

	// macOS Security subsection
	fmt.Println(ui.SubSectionStyle.Render("macOS Security"))
	printMacOSSecurityTable()

	// Identity subsection
	fmt.Println(ui.SubSectionStyle.Render("Identity"))
	printIdentityTable()
}

func printSetupTable() {
	rows := [][]string{}

	for _, script := range config.Scripts {
		result := config.CheckScript(script)
		if result.Installed {
			rows = append(rows, []string{script.Name, result.Detail, ui.SuccessStyle.Render("✓")})
		} else {
			rows = append(rows, []string{script.Name, "", ui.DangerStyle.Render("✗")})
		}
	}

	fmt.Println(ui.RenderTable(rows, ui.StatusTableColumns))
	fmt.Println()
}

func printMacOSSecurityTable() {
	rows := [][]string{}
	for _, check := range config.SecurityChecks {
		result := check.CheckFn()
		var status string
		if result.Installed == check.GoodWhen {
			status = ui.SuccessStyle.Render("✓")
		} else {
			status = ui.WarningStyle.Render("!")
		}
		rows = append(rows, []string{check.Name, check.Description, result.Detail, status})
	}

	fmt.Println(ui.RenderTable(rows, ui.CheckTableColumns))
	fmt.Println()
}

func printIdentityTable() {
	rows := [][]string{}
	for _, check := range config.IdentityChecks {
		result := check.CheckFn()
		var status string
		if result.Installed == check.GoodWhen {
			status = ui.SuccessStyle.Render("✓")
		} else {
			status = ui.WarningStyle.Render("!")
		}
		rows = append(rows, []string{check.Name, check.Description, result.Detail, status})
	}

	fmt.Println(ui.RenderTable(rows, ui.CheckTableColumns))
	fmt.Println()
}
