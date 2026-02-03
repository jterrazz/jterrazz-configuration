package status

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/domain/tool"
)

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

// Loader manages parallel loading of status items
type Loader struct {
	items   []Item
	updates chan UpdateMsg
	started bool
	mu      sync.Mutex
}

// NewLoader creates a new loader with all items in pending state
func NewLoader() *Loader {
	loader := &Loader{
		updates: make(chan UpdateMsg, 100),
	}
	loader.buildItems()
	return loader
}

// GetItems returns a copy of all items
func (l *Loader) GetItems() []Item {
	l.mu.Lock()
	defer l.mu.Unlock()
	items := make([]Item, len(l.items))
	copy(items, l.items)
	return items
}

// GetPendingCount returns the number of items that need loading
func (l *Loader) GetPendingCount() int {
	count := 0
	for _, item := range l.items {
		if !item.Loaded && item.Kind != KindHeader {
			count++
		}
	}
	return count
}

// buildItems creates all status items in display order
func (l *Loader) buildItems() {
	// System Info header
	l.addItem(Item{
		ID:      "sysinfo",
		Kind:    KindSystemInfo,
		Section: "System",
		Name:    "System Info",
	})

	// Setup section
	l.addItem(Item{ID: "header-setup", Kind: KindHeader, Section: "System", SubSection: "Setup", Loaded: true})
	for _, script := range config.Scripts {
		if script.CheckFn == nil {
			continue
		}
		l.addItem(Item{
			ID:          "setup-" + script.Name,
			Kind:        KindSetup,
			Section:     "System",
			SubSection:  "Setup",
			Name:        script.Name,
			Description: script.Description,
		})
	}

	// Security section
	l.addItem(Item{ID: "header-security", Kind: KindHeader, Section: "System", SubSection: "MacOS Security", Loaded: true})
	for _, check := range config.SecurityChecks {
		l.addItem(Item{
			ID:          "security-" + check.Name,
			Kind:        KindSecurity,
			Section:     "System",
			SubSection:  "MacOS Security",
			Name:        check.Name,
			Description: check.Description,
			GoodWhen:    check.GoodWhen,
		})
	}

	// Identity section
	l.addItem(Item{ID: "header-identity", Kind: KindHeader, Section: "System", SubSection: "Identity", Loaded: true})
	for _, check := range config.IdentityChecks {
		l.addItem(Item{
			ID:          "identity-" + check.Name,
			Kind:        KindIdentity,
			Section:     "System",
			SubSection:  "Identity",
			Name:        check.Name,
			Description: check.Description,
			GoodWhen:    check.GoodWhen,
		})
	}

	// Tools sections
	for _, category := range config.ToolCategories {
		tools := config.GetToolsByCategory(category)
		if len(tools) == 0 {
			continue
		}
		l.addItem(Item{ID: "header-tools-" + string(category), Kind: KindHeader, Section: "Tools", SubSection: string(category), Loaded: true})
		for _, t := range tools {
			l.addItem(Item{
				ID:         "tool-" + t.Name,
				Kind:       KindTool,
				Section:    "Tools",
				SubSection: string(category),
				Name:       t.Name,
				Method:     t.Method.String(),
			})
		}
	}

	// Network section
	l.addItem(Item{ID: "header-network", Kind: KindHeader, Section: "Resources", SubSection: "Network", Loaded: true})
	for _, check := range config.NetworkChecks {
		l.addItem(Item{
			ID:         "network-" + check.Name,
			Kind:       KindNetwork,
			Section:    "Resources",
			SubSection: "Network",
			Name:       check.Name,
		})
	}

	// Disk section
	l.addItem(Item{ID: "header-disk", Kind: KindHeader, Section: "Resources", SubSection: "Disk Usage", Loaded: true})
	for _, check := range config.MainDiskChecks {
		l.addItem(Item{
			ID:         "disk-" + check.Name,
			Kind:       KindDisk,
			Section:    "Resources",
			SubSection: "Disk Usage",
			Name:       check.Name,
		})
	}

	// Cache section
	l.addItem(Item{ID: "header-cache", Kind: KindHeader, Section: "Resources", SubSection: "Caches & Cleanable", Loaded: true})
	for _, check := range config.CacheChecks {
		l.addItem(Item{
			ID:         "cache-" + check.Name,
			Kind:       KindCache,
			Section:    "Resources",
			SubSection: "Caches & Cleanable",
			Name:       check.Name,
		})
	}
}

func (l *Loader) addItem(item Item) {
	l.items = append(l.items, item)
}

// Start launches all checks in parallel (call only once)
func (l *Loader) Start() {
	l.mu.Lock()
	if l.started {
		l.mu.Unlock()
		return
	}
	l.started = true
	l.mu.Unlock()

	var wg sync.WaitGroup

	// System info
	wg.Add(1)
	go func() {
		defer wg.Done()
		item := l.loadSystemInfo()
		l.updates <- UpdateMsg{ID: item.ID, Item: item}
	}()

	// Setup checks
	for _, script := range config.Scripts {
		if script.CheckFn == nil {
			continue
		}
		wg.Add(1)
		go func(s config.Script) {
			defer wg.Done()
			result := config.CheckScript(s)
			item := Item{
				ID:        "setup-" + s.Name,
				Kind:      KindSetup,
				Name:      s.Name,
				Loaded:    true,
				Installed: result.Installed,
				Detail:    result.Detail,
			}
			l.updates <- UpdateMsg{ID: item.ID, Item: item}
		}(script)
	}

	// Security checks
	for _, check := range config.SecurityChecks {
		wg.Add(1)
		go func(c config.SecurityCheck) {
			defer wg.Done()
			result := c.CheckFn()
			item := Item{
				ID:          "security-" + c.Name,
				Kind:        KindSecurity,
				Name:        c.Name,
				Description: c.Description,
				Loaded:      true,
				Installed:   result.Installed,
				Detail:      result.Detail,
				GoodWhen:    c.GoodWhen,
			}
			l.updates <- UpdateMsg{ID: item.ID, Item: item}
		}(check)
	}

	// Identity checks
	for _, check := range config.IdentityChecks {
		wg.Add(1)
		go func(c config.IdentityCheck) {
			defer wg.Done()
			result := c.CheckFn()
			item := Item{
				ID:          "identity-" + c.Name,
				Kind:        KindIdentity,
				Name:        c.Name,
				Description: c.Description,
				Loaded:      true,
				Installed:   result.Installed,
				Detail:      result.Detail,
				GoodWhen:    c.GoodWhen,
			}
			l.updates <- UpdateMsg{ID: item.ID, Item: item}
		}(check)
	}

	// Tool checks
	for _, t := range config.Tools {
		wg.Add(1)
		go func(t config.Tool) {
			defer wg.Done()
			result := t.Check()
			item := Item{
				ID:        "tool-" + t.Name,
				Kind:      KindTool,
				Name:      t.Name,
				Loaded:    true,
				Installed: result.Installed,
				Version:   result.Version,
				Status:    result.Status,
				Method:    t.Method.String(),
			}
			l.updates <- UpdateMsg{ID: item.ID, Item: item}
		}(t)
	}

	// Network checks
	for _, check := range config.NetworkChecks {
		wg.Add(1)
		go func(c config.ResourceCheck) {
			defer wg.Done()
			result := c.CheckFn()
			item := Item{
				ID:        "network-" + c.Name,
				Kind:      KindNetwork,
				Name:      c.Name,
				Loaded:    true,
				Available: result.Available,
				Value:     result.Value,
				Style:     result.Style,
			}
			l.updates <- UpdateMsg{ID: item.ID, Item: item}
		}(check)
	}

	// Disk checks
	for _, check := range config.MainDiskChecks {
		wg.Add(1)
		go func(c config.DiskCheck) {
			defer wg.Done()
			result := c.Check()
			item := Item{
				ID:        "disk-" + c.Name,
				Kind:      KindDisk,
				Name:      c.Name,
				Loaded:    true,
				Available: result.Available,
				Value:     result.Value,
				Style:     result.Style,
			}
			l.updates <- UpdateMsg{ID: item.ID, Item: item}
		}(check)
	}

	// Cache checks
	for _, check := range config.CacheChecks {
		wg.Add(1)
		go func(c config.DiskCheck) {
			defer wg.Done()
			result := c.Check()
			item := Item{
				ID:        "cache-" + c.Name,
				Kind:      KindCache,
				Name:      c.Name,
				Loaded:    true,
				Available: result.Available,
				Value:     result.Value,
				Style:     result.Style,
			}
			l.updates <- UpdateMsg{ID: item.ID, Item: item}
		}(check)
	}

	// Close channel when all done
	go func() {
		wg.Wait()
		close(l.updates)
	}()
}

// WaitForUpdate returns a command that waits for the next update
func (l *Loader) WaitForUpdate() tea.Cmd {
	return func() tea.Msg {
		update, ok := <-l.updates
		if !ok {
			return AllLoadedMsg{}
		}
		return update
	}
}

// loadSystemInfo loads system information
func (l *Loader) loadSystemInfo() Item {
	hostname, _ := os.Hostname()
	// Shorten hostname (remove .local suffix and truncate if too long)
	if idx := strings.Index(hostname, "."); idx > 0 {
		hostname = hostname[:idx]
	}
	if len(hostname) > 20 {
		hostname = hostname[:20]
	}

	osInfo := tool.GetCommandOutput("uname", "-sr")
	arch := tool.GetCommandOutput("uname", "-m")
	user := os.Getenv("USER")
	shell := filepath.Base(os.Getenv("SHELL"))

	return Item{
		ID:     "sysinfo",
		Kind:   KindSystemInfo,
		Loaded: true,
		Detail: osInfo + " " + arch + " • " + hostname + " • " + user + " • " + shell,
	}
}
