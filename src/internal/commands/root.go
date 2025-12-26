package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "j",
	Short: "jterrazz unified command system",
	Long:  "A unified CLI tool for development workflow automation.",
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

func Execute() error {
	return rootCmd.Execute()
}
