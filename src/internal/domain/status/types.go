package status

// ItemKind represents the type of status item
type ItemKind int

const (
	KindHeader ItemKind = iota
	KindSetup
	KindSecurity
	KindIdentity
	KindTool
	KindNetwork
	KindDisk
	KindCache
	KindSystemInfo
)

// Item represents a single item in the status display
type Item struct {
	ID          string
	Kind        ItemKind
	Section     string
	SubSection  string
	Name        string
	Description string
	Loaded      bool

	// Result data (populated after loading)
	Installed bool
	Version   string
	Status    string
	Detail    string
	Value     string
	Style     string // Semantic style: "success", "warning", "muted", etc.
	GoodWhen  bool   // For checks: true means Installed=true is good
	Method    string // Install method for tools
	Available bool   // For resources: whether the resource exists
}

// UpdateMsg is sent when a status item finishes loading
type UpdateMsg struct {
	ID   string
	Item Item
}

// AllLoadedMsg is sent when all items have finished loading
type AllLoadedMsg struct{}
