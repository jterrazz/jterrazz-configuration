package commands

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jterrazz/jterrazz-cli/internal/config"
	"github.com/jterrazz/jterrazz-cli/internal/tool"
)

// =============================================================================
// Status Item Types
// =============================================================================

// StatusItemKind represents the type of status item
type StatusItemKind int

const (
	StatusItemHeader StatusItemKind = iota
	StatusItemSetup
	StatusItemSecurity
	StatusItemIdentity
	StatusItemTool
	StatusItemNetwork
	StatusItemDisk
	StatusItemCache
	StatusItemSystemInfo
)

// StatusItem represents a single item in the status display
type StatusItem struct {
	ID          string
	Kind        StatusItemKind
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
	Style     string
	GoodWhen  bool
	Method    string
	Available bool
}

// StatusUpdateMsg is sent when a status item finishes loading
type StatusUpdateMsg struct {
	ID   string
	Item StatusItem
}

// AllLoadedMsg is sent when all items have finished loading
type AllLoadedMsg struct{}

// =============================================================================
// Status Loader
// =============================================================================

// StatusLoader manages parallel loading of status items
type StatusLoader struct {
	items   []StatusItem
	updates chan StatusUpdateMsg
	started bool
	mu      sync.Mutex
}

// NewStatusLoader creates a new loader with all items in pending state
func NewStatusLoader() *StatusLoader {
	loader := &StatusLoader{
		updates: make(chan StatusUpdateMsg, 100),
	}
	loader.buildItems()
	return loader
}

// GetItems returns a copy of all items
func (l *StatusLoader) GetItems() []StatusItem {
	l.mu.Lock()
	defer l.mu.Unlock()
	items := make([]StatusItem, len(l.items))
	copy(items, l.items)
	return items
}

// GetPendingCount returns the number of items that need loading
func (l *StatusLoader) GetPendingCount() int {
	count := 0
	for _, item := range l.items {
		if !item.Loaded && item.Kind != StatusItemHeader {
			count++
		}
	}
	return count
}

// buildItems creates all status items in display order
func (l *StatusLoader) buildItems() {
	// System Info header
	l.addItem(StatusItem{
		ID:      "sysinfo",
		Kind:    StatusItemSystemInfo,
		Section: "System",
		Name:    "System Info",
	})

	// Setup section
	l.addItem(StatusItem{ID: "header-setup", Kind: StatusItemHeader, Section: "System", SubSection: "Setup", Loaded: true})
	for _, script := range config.Scripts {
		if script.CheckFn == nil {
			continue
		}
		l.addItem(StatusItem{
			ID:          "setup-" + script.Name,
			Kind:        StatusItemSetup,
			Section:     "System",
			SubSection:  "Setup",
			Name:        script.Name,
			Description: script.Description,
		})
	}

	// Security section
	l.addItem(StatusItem{ID: "header-security", Kind: StatusItemHeader, Section: "System", SubSection: "MacOS Security", Loaded: true})
	for _, check := range config.SecurityChecks {
		l.addItem(StatusItem{
			ID:          "security-" + check.Name,
			Kind:        StatusItemSecurity,
			Section:     "System",
			SubSection:  "MacOS Security",
			Name:        check.Name,
			Description: check.Description,
			GoodWhen:    check.GoodWhen,
		})
	}

	// Identity section
	l.addItem(StatusItem{ID: "header-identity", Kind: StatusItemHeader, Section: "System", SubSection: "Identity", Loaded: true})
	for _, check := range config.IdentityChecks {
		l.addItem(StatusItem{
			ID:          "identity-" + check.Name,
			Kind:        StatusItemIdentity,
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
		l.addItem(StatusItem{ID: "header-tools-" + string(category), Kind: StatusItemHeader, Section: "Tools", SubSection: string(category), Loaded: true})
		for _, t := range tools {
			l.addItem(StatusItem{
				ID:         "tool-" + t.Name,
				Kind:       StatusItemTool,
				Section:    "Tools",
				SubSection: string(category),
				Name:       t.Name,
				Method:     t.Method.String(),
			})
		}
	}

	// Network section
	l.addItem(StatusItem{ID: "header-network", Kind: StatusItemHeader, Section: "Resources", SubSection: "Network", Loaded: true})
	for _, check := range config.NetworkChecks {
		l.addItem(StatusItem{
			ID:         "network-" + check.Name,
			Kind:       StatusItemNetwork,
			Section:    "Resources",
			SubSection: "Network",
			Name:       check.Name,
		})
	}

	// Disk section
	l.addItem(StatusItem{ID: "header-disk", Kind: StatusItemHeader, Section: "Resources", SubSection: "Disk Usage", Loaded: true})
	for _, check := range config.MainDiskChecks {
		l.addItem(StatusItem{
			ID:         "disk-" + check.Name,
			Kind:       StatusItemDisk,
			Section:    "Resources",
			SubSection: "Disk Usage",
			Name:       check.Name,
		})
	}

	// Cache section
	l.addItem(StatusItem{ID: "header-cache", Kind: StatusItemHeader, Section: "Resources", SubSection: "Caches & Cleanable", Loaded: true})
	for _, check := range config.CacheChecks {
		l.addItem(StatusItem{
			ID:         "cache-" + check.Name,
			Kind:       StatusItemCache,
			Section:    "Resources",
			SubSection: "Caches & Cleanable",
			Name:       check.Name,
		})
	}
}

func (l *StatusLoader) addItem(item StatusItem) {
	l.items = append(l.items, item)
}

// Start launches all checks in parallel (call only once)
func (l *StatusLoader) Start() {
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
		l.updates <- StatusUpdateMsg{ID: item.ID, Item: item}
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
			item := StatusItem{
				ID:        "setup-" + s.Name,
				Kind:      StatusItemSetup,
				Name:      s.Name,
				Loaded:    true,
				Installed: result.Installed,
				Detail:    result.Detail,
			}
			l.updates <- StatusUpdateMsg{ID: item.ID, Item: item}
		}(script)
	}

	// Security checks
	for _, check := range config.SecurityChecks {
		wg.Add(1)
		go func(c config.SecurityCheck) {
			defer wg.Done()
			result := c.CheckFn()
			item := StatusItem{
				ID:          "security-" + c.Name,
				Kind:        StatusItemSecurity,
				Name:        c.Name,
				Description: c.Description,
				Loaded:      true,
				Installed:   result.Installed,
				Detail:      result.Detail,
				GoodWhen:    c.GoodWhen,
			}
			l.updates <- StatusUpdateMsg{ID: item.ID, Item: item}
		}(check)
	}

	// Identity checks
	for _, check := range config.IdentityChecks {
		wg.Add(1)
		go func(c config.IdentityCheck) {
			defer wg.Done()
			result := c.CheckFn()
			item := StatusItem{
				ID:          "identity-" + c.Name,
				Kind:        StatusItemIdentity,
				Name:        c.Name,
				Description: c.Description,
				Loaded:      true,
				Installed:   result.Installed,
				Detail:      result.Detail,
				GoodWhen:    c.GoodWhen,
			}
			l.updates <- StatusUpdateMsg{ID: item.ID, Item: item}
		}(check)
	}

	// Tool checks
	for _, t := range config.Tools {
		wg.Add(1)
		go func(t config.Tool) {
			defer wg.Done()
			result := t.Check()
			item := StatusItem{
				ID:        "tool-" + t.Name,
				Kind:      StatusItemTool,
				Name:      t.Name,
				Loaded:    true,
				Installed: result.Installed,
				Version:   result.Version,
				Status:    result.Status,
				Method:    t.Method.String(),
			}
			l.updates <- StatusUpdateMsg{ID: item.ID, Item: item}
		}(t)
	}

	// Network checks
	for _, check := range config.NetworkChecks {
		wg.Add(1)
		go func(c config.ResourceCheck) {
			defer wg.Done()
			result := c.CheckFn()
			item := StatusItem{
				ID:        "network-" + c.Name,
				Kind:      StatusItemNetwork,
				Name:      c.Name,
				Loaded:    true,
				Available: result.Available,
				Value:     result.Value,
				Style:     result.Style,
			}
			l.updates <- StatusUpdateMsg{ID: item.ID, Item: item}
		}(check)
	}

	// Disk checks
	for _, check := range config.MainDiskChecks {
		wg.Add(1)
		go func(c config.DiskCheck) {
			defer wg.Done()
			result := c.Check()
			item := StatusItem{
				ID:        "disk-" + c.Name,
				Kind:      StatusItemDisk,
				Name:      c.Name,
				Loaded:    true,
				Available: result.Available,
				Value:     result.Value,
				Style:     result.Style,
			}
			l.updates <- StatusUpdateMsg{ID: item.ID, Item: item}
		}(check)
	}

	// Cache checks
	for _, check := range config.CacheChecks {
		wg.Add(1)
		go func(c config.DiskCheck) {
			defer wg.Done()
			result := c.Check()
			item := StatusItem{
				ID:        "cache-" + c.Name,
				Kind:      StatusItemCache,
				Name:      c.Name,
				Loaded:    true,
				Available: result.Available,
				Value:     result.Value,
				Style:     result.Style,
			}
			l.updates <- StatusUpdateMsg{ID: item.ID, Item: item}
		}(check)
	}

	// Close channel when all done
	go func() {
		wg.Wait()
		close(l.updates)
	}()
}

// WaitForUpdate returns a command that waits for the next update
func (l *StatusLoader) WaitForUpdate() tea.Cmd {
	return func() tea.Msg {
		update, ok := <-l.updates
		if !ok {
			return AllLoadedMsg{}
		}
		return update
	}
}

// loadSystemInfo loads system information
func (l *StatusLoader) loadSystemInfo() StatusItem {
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

	return StatusItem{
		ID:     "sysinfo",
		Kind:   StatusItemSystemInfo,
		Loaded: true,
		Detail: osInfo + " " + arch + " • " + hostname + " • " + user + " • " + shell,
	}
}
