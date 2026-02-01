package commands

import (
	"fmt"

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
	fmt.Println(ui.TitleStyle.MarginTop(1).Render("j status"))
	fmt.Println()

	printSystemInfo()
	printSystemSection()
	printToolsSection()
	printResourcesSection()
}
