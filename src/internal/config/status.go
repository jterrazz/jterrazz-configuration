package config

// StatusSection represents a section in the status display
type StatusSection struct {
	Title    string
	SubTitle string // Optional subsection title
	RenderFn func() // Function to render this section
}

// StatusSections defines all status sections in display order
// This is the single source of truth for the status command layout
var StatusSections = []StatusSection{
	// Setup section (standalone)
	{Title: "Setup", SubTitle: "Setup", RenderFn: nil}, // Uses Scripts

	// System section with subsections
	{Title: "System", SubTitle: "Security", RenderFn: nil}, // Uses SecurityChecks
	{Title: "System", SubTitle: "Identity", RenderFn: nil}, // Uses IdentityChecks

	// Tools section - one subsection per category
	{Title: "Tools", SubTitle: "Package Managers", RenderFn: nil},
	{Title: "Tools", SubTitle: "Runtimes", RenderFn: nil},
	{Title: "Tools", SubTitle: "DevOps", RenderFn: nil},
	{Title: "Tools", SubTitle: "AI", RenderFn: nil},
	{Title: "Tools", SubTitle: "Terminal & Git", RenderFn: nil},
	{Title: "Tools", SubTitle: "GUI Apps", RenderFn: nil},
	{Title: "Tools", SubTitle: "Mac App Store", RenderFn: nil},

	// Resources section with subsections
	{Title: "Resources", SubTitle: "Top Processes", RenderFn: nil},      // Uses ProcessChecks
	{Title: "Resources", SubTitle: "Network", RenderFn: nil},            // Uses NetworkChecks
	{Title: "Resources", SubTitle: "Caches & Cleanable", RenderFn: nil}, // Uses CacheChecks
}

// ToolCategories defines the order of tool categories in status display
var ToolCategories = []ToolCategory{
	CategoryPackageManager,
	CategoryRuntimes,
	CategoryDevOps,
	CategoryAI,
	CategoryTerminalGit,
	CategoryGUIApps,
	CategoryMacAppStore,
}
