package status

import (
	"sort"
	"strings"

	"github.com/jterrazz/jterrazz-cli/src/internal/domain/status"
	"github.com/jterrazz/jterrazz-cli/src/internal/presentation/components"
)

// renderContent renders all sections
func (m Model) renderContent() string {
	var b strings.Builder

	sections := m.groupBySection()
	boxWidth := m.width
	if boxWidth < 40 {
		boxWidth = 40
	}

	isFirst := true
	for _, section := range []string{"Setup", "System", "Tools", "Resources"} {
		subsections, ok := sections[section]
		if !ok {
			continue
		}

		// Section header with decorative line (no leading newline for first section)
		if !isFirst {
			b.WriteString("\n")
		}
		isFirst = false
		sectionHeader := components.SectionHeader(section, boxWidth)
		b.WriteString(sectionHeader)
		b.WriteString("\n")

		// Collect all items in this section for column width calculation
		var allSectionItems []status.Item
		for _, subsection := range getSubsectionOrder(section) {
			items, ok := subsections[subsection]
			if !ok {
				continue
			}
			for _, item := range items {
				if item.Kind == status.KindNetwork || item.Kind == status.KindCache {
					if !item.Loaded || item.Available {
						allSectionItems = append(allSectionItems, item)
					}
				} else {
					allSectionItems = append(allSectionItems, item)
				}
			}
		}

		// Calculate column widths for the entire section
		colWidths := calculateColumnWidths(allSectionItems)

		// Render subsections as boxes
		for _, subsection := range getSubsectionOrder(section) {
			items, ok := subsections[subsection]
			if !ok {
				continue
			}

			// Filter out unavailable items
			var visibleItems []status.Item
			for _, item := range items {
				if item.Kind == status.KindNetwork || item.Kind == status.KindCache {
					if !item.Loaded || item.Available {
						visibleItems = append(visibleItems, item)
					}
				} else {
					visibleItems = append(visibleItems, item)
				}
			}

			if len(visibleItems) == 0 {
				continue
			}

			box := m.renderSubsectionBox(subsection, visibleItems, boxWidth, colWidths)
			b.WriteString(box)
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (m Model) groupBySection() map[string]map[string][]status.Item {
	sections := make(map[string]map[string][]status.Item)

	for _, baseItem := range m.itemOrder {
		item := m.items[baseItem.ID]
		if item.Kind == status.KindHeader || item.Kind == status.KindSystemInfo {
			continue
		}

		if sections[item.Section] == nil {
			sections[item.Section] = make(map[string][]status.Item)
		}
		sections[item.Section][item.SubSection] = append(sections[item.Section][item.SubSection], item)
	}

	// Sort items A-Z by name within each subsection
	for _, subsections := range sections {
		for subsection, items := range subsections {
			sort.Slice(items, func(i, j int) bool {
				return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
			})
			subsections[subsection] = items
		}
	}

	return sections
}

func getSubsectionOrder(section string) []string {
	switch section {
	case "Setup":
		return []string{"Setup"}
	case "System":
		return []string{"Security", "Identity"}
	case "Tools":
		return []string{"Package Managers", "Languages", "Infrastructure", "AI", "Apps", "System Tools"}
	case "Resources":
		return []string{"Top Processes", "Network", "Caches & Cleanable"}
	}
	return nil
}

func (m Model) renderSubsectionBox(title string, items []status.Item, width int, colWidths ColumnWidths) string {
	// Render rows
	var rows []string
	for _, item := range items {
		if item.Kind == status.KindProcess {
			// Process items render multiple rows
			processRows := m.renderProcessRows(item)
			rows = append(rows, processRows...)
		} else {
			row := m.renderTableRow(item, colWidths)
			rows = append(rows, row)
		}
	}

	return components.SubsectionBox(title, rows, width)
}

// ColumnWidths holds calculated column widths for alignment
type ColumnWidths struct {
	Name    int
	Desc    int
	Version int
	Method  int
	Detail  int
}

func calculateColumnWidths(items []status.Item) ColumnWidths {
	widths := ColumnWidths{}
	for _, item := range items {
		if len(item.Name) > widths.Name {
			widths.Name = len(item.Name)
		}
		if len(item.Description) > widths.Desc {
			widths.Desc = len(item.Description)
		}
		if len(item.Version) > widths.Version {
			widths.Version = len(item.Version)
		}
		if len(item.Method) > widths.Method {
			widths.Method = len(item.Method)
		}
		if len(item.Detail) > widths.Detail {
			widths.Detail = len(item.Detail)
		}
	}
	return widths
}
