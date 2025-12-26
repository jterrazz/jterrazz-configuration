package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update system packages (Homebrew + npm global)",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Println(cyan("🔄 Updating system packages..."))

		if commandExists("brew") {
			fmt.Println(cyan("🍺 Updating Homebrew packages..."))
			runCommand("brew", "update")
			runCommand("brew", "upgrade")
		} else {
			printError("Homebrew not found")
			return
		}

		if commandExists("npm") {
			fmt.Println(cyan("📦 Updating npm global packages..."))
			runCommand("npm", "update", "-g")
		} else {
			printWarning("npm not found, skipping global package updates")
		}

		fmt.Println(green("✅ System update completed"))
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
