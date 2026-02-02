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
	// System section with subsections
	{Title: "System", SubTitle: "Setup", RenderFn: nil},          // Uses Scripts
	{Title: "System", SubTitle: "MacOS Security", RenderFn: nil}, // Uses SecurityChecks
	{Title: "System", SubTitle: "Identity", RenderFn: nil},       // Uses IdentityChecks

	// Tools section - one subsection per category
	{Title: "Tools", SubTitle: "Package Managers", RenderFn: nil},
	{Title: "Tools", SubTitle: "Languages", RenderFn: nil},
	{Title: "Tools", SubTitle: "Infrastructure", RenderFn: nil},
	{Title: "Tools", SubTitle: "AI", RenderFn: nil},
	{Title: "Tools", SubTitle: "Apps", RenderFn: nil},
	{Title: "Tools", SubTitle: "System Tools", RenderFn: nil},

	// Resources section with subsections
	{Title: "Resources", SubTitle: "Network", RenderFn: nil},            // Uses NetworkChecks
	{Title: "Resources", SubTitle: "Disk Usage", RenderFn: nil},         // Uses MainDiskChecks
	{Title: "Resources", SubTitle: "Caches & Cleanable", RenderFn: nil}, // Uses CacheChecks
}

// ToolCategories defines the order of tool categories in status display
var ToolCategories = []ToolCategory{
	CategoryPackageManager,
	CategoryLanguages,
	CategoryInfrastructure,
	CategoryAI,
	CategoryApps,
	CategorySystemTools,
}
