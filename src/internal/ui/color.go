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

// =============================================================================
// Basic Print Functions
// =============================================================================

// Print prints a line
func Print(s string) {
	fmt.Println(s)
}

// Printf prints a formatted line
func Printf(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

// PrintEmpty prints an empty line
func PrintEmpty() {
	fmt.Println()
}

// =============================================================================
// Message Print Functions
// =============================================================================

// PrintError prints an error message
func PrintError(msg string) {
	fmt.Printf("%s %s\n", Red("Error:"), msg)
}

// PrintWarning prints a warning message
func PrintWarning(msg string) {
	fmt.Printf("%s %s\n", Yellow("Warning:"), msg)
}

// PrintSuccess prints a success message with checkmark
func PrintSuccess(msg string) {
	fmt.Printf("%s %s\n", Green(IconCheck), msg)
}

// PrintInfo prints an info message in cyan
func PrintInfo(msg string) {
	fmt.Println(Cyan(msg))
}

// PrintDim prints a dimmed/muted message
func PrintDim(msg string) {
	fmt.Println(Dim(msg))
}

// =============================================================================
// Action Print Functions
// =============================================================================

// PrintAction prints an action being performed (e.g., "🔄 Updating...")
func PrintAction(emoji, msg string) {
	fmt.Println(Cyan(emoji + " " + msg))
}

// PrintDone prints a completion message
func PrintDone(msg string) {
	fmt.Println(Green("✅ " + msg))
}

// PrintInstalling prints an installing message
func PrintInstalling(name string) {
	fmt.Printf("  📥 Installing %s...\n", name)
}

// PrintInstallingVia prints an installing message with method
func PrintInstallingVia(name, method string) {
	fmt.Printf("  📥 Installing %s (via %s)...\n", name, method)
}

// =============================================================================
// Section Print Functions
// =============================================================================

// PrintTitle prints a styled title
func PrintTitle(title string) {
	fmt.Println(TitleStyle.MarginTop(1).Render(title))
}

// PrintSection prints a section header
func PrintSection(title string) {
	fmt.Println(SectionStyle.Render(title))
}

// PrintSubSection prints a subsection header
func PrintSubSection(title string) {
	fmt.Println(SubSectionStyle.Render(title))
}

// PrintCategory prints a category header (dimmed)
func PrintCategory(name string) {
	fmt.Println(Dim(name))
}

// =============================================================================
// List/Row Print Functions
// =============================================================================

// PrintRow prints a status row (icon + label + detail)
func PrintRow(ok bool, label, detail string) {
	fmt.Println(RenderRow(ok, label, detail))
}

// PrintTable prints a table
func PrintTable(rows [][]string, columns []ColumnConfig) {
	fmt.Println(RenderTable(rows, columns))
}

// =============================================================================
// Usage Print Functions
// =============================================================================

// PrintUsage prints usage instructions
func PrintUsage(lines ...string) {
	for _, line := range lines {
		fmt.Println(Dim(line))
	}
}

// =============================================================================
// Inline Render Functions (return string, don't print)
// =============================================================================

// RenderSpecial renders text in special/cyan style
func RenderSpecial(s string) string {
	return SpecialStyle.Render(s)
}

// RenderMuted renders text in muted style
func RenderMuted(s string) string {
	return MutedStyle.Render(s)
}

// RenderSuccess renders text in success style
func RenderSuccess(s string) string {
	return SuccessStyle.Render(s)
}

// RenderWarning renders text in warning style
func RenderWarning(s string) string {
	return WarningStyle.Render(s)
}

// RenderDanger renders text in danger style
func RenderDanger(s string) string {
	return DangerStyle.Render(s)
}

// RenderStatusIcon renders a status icon (check, cross, warning)
func RenderStatusIcon(ok bool) string {
	if ok {
		return SuccessStyle.Render(IconCheck)
	}
	return DangerStyle.Render(IconCross)
}

// RenderStatusIconWithCondition renders icon based on condition matching expected
func RenderStatusIconWithCondition(value, expected bool) string {
	if value == expected {
		return SuccessStyle.Render(IconCheck)
	}
	return WarningStyle.Render(IconWarning)
}

// =============================================================================
// System Info Helpers
// =============================================================================

// PrintSystemHeader prints OS info header (e.g., "Darwin 25.2.0 • arm64")
func PrintSystemHeader(osInfo, arch, hostname, user, shell string) {
	fmt.Printf("%s • %s\n", SpecialStyle.Render(osInfo), MutedStyle.Render(arch))
	fmt.Printf("%s • %s • %s\n\n", MutedStyle.Render(hostname), MutedStyle.Render(user), MutedStyle.Render(shell))
}

// PrintHint prints a hint line (e.g., "run j clean or mo clean")
func PrintHint(parts ...string) {
	result := ""
	for i, p := range parts {
		if i%2 == 0 {
			result += MutedStyle.Render(p)
		} else {
			result += SpecialStyle.Render(p)
		}
	}
	fmt.Println(result)
}
