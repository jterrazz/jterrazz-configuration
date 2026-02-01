package ui

import (
	"fmt"

	"github.com/fatih/color"
)

// Color helpers - pre-configured color functions
var (
	Cyan   = color.New(color.FgCyan).SprintFunc()
	Green  = color.New(color.FgGreen).SprintFunc()
	Red    = color.New(color.FgRed).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Dim    = color.New(color.FgHiBlack).SprintFunc()
)

// PrintError prints an error message to stdout
func PrintError(msg string) {
	fmt.Printf("%s %s\n", Red("Error:"), msg)
}

// PrintWarning prints a warning message to stdout
func PrintWarning(msg string) {
	fmt.Printf("%s %s\n", Yellow("Warning:"), msg)
}

// PrintSuccess prints a success message to stdout
func PrintSuccess(msg string) {
	fmt.Printf("%s %s\n", Green("✓"), msg)
}

// PrintInfo prints an info message to stdout
func PrintInfo(msg string) {
	fmt.Println(Cyan(msg))
}
