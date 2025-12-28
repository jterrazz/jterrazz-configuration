package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Docker container and image management",
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run common development commands",
}

var dockerRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove all containers",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Println(cyan("🧹 Removing all Docker containers..."))

		out, err := exec.Command("docker", "ps", "-aq").Output()
		if err != nil || strings.TrimSpace(string(out)) == "" {
			fmt.Println("No containers to remove")
			return
		}

		containers := strings.Fields(string(out))
		runCommand("docker", append([]string{"rm", "-vf"}, containers...)...)
	},
}

var dockerRmiCmd = &cobra.Command{
	Use:   "rmi",
	Short: "Remove all images",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Println(cyan("🗑️  Removing all Docker images..."))

		out, err := exec.Command("docker", "images", "-aq").Output()
		if err != nil || strings.TrimSpace(string(out)) == "" {
			fmt.Println("No images to remove")
			return
		}

		images := strings.Fields(string(out))
		runCommand("docker", append([]string{"rmi", "-f"}, images...)...)
	},
}

var dockerCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up Docker system (prune)",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Println(cyan("🧹 Cleaning up Docker system..."))
		runCommand("docker", "system", "prune", "-af")
	},
}

var dockerResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Remove all containers and images",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Println(cyan("🔄 Resetting Docker (removing containers and images)..."))
		dockerRmCmd.Run(cmd, args)
		dockerRmiCmd.Run(cmd, args)
	},
}

func init() {
	dockerCmd.AddCommand(dockerRmCmd)
	dockerCmd.AddCommand(dockerRmiCmd)
	dockerCmd.AddCommand(dockerCleanCmd)
	dockerCmd.AddCommand(dockerResetCmd)
	runCmd.AddCommand(dockerCmd)
	rootCmd.AddCommand(runCmd)
}
