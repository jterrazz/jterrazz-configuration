package commands

import (
	"fmt"

	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run common development commands",
}

func init() {
	// Dynamically build commands from config
	for _, cmdGroup := range config.RunCommands {
		parentCmd := &cobra.Command{
			Use:   cmdGroup.Name,
			Short: cmdGroup.Description,
		}

		for _, sub := range cmdGroup.Subcommands {
			subCmd := createSubcommand(sub)
			parentCmd.AddCommand(subCmd)
		}

		runCmd.AddCommand(parentCmd)
	}

	rootCmd.AddCommand(runCmd)
}

func createSubcommand(sub config.RunSubcommand) *cobra.Command {
	cmd := &cobra.Command{
		Use:   sub.Name,
		Short: sub.Description,
		Run: func(cmd *cobra.Command, args []string) {
			if err := sub.RunFn(args); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		},
	}

	if sub.MinArgs > 0 {
		cmd.Use = fmt.Sprintf("%s [args]", sub.Name)
		cmd.Args = cobra.MinimumNArgs(sub.MinArgs)
	}

	return cmd
}
