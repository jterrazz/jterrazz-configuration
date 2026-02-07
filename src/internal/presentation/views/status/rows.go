package status

import (
	"github.com/jterrazz/jterrazz-cli/src/internal/domain/status"
	"github.com/jterrazz/jterrazz-cli/src/internal/presentation/components"
	"github.com/jterrazz/jterrazz-cli/src/internal/presentation/theme"
)

func (m Model) renderTableRow(item status.Item, colWidths ColumnWidths) string {
	if !item.Loaded {
		switch item.Kind {
		case status.KindSetup:
			return m.renderSetupRowLoading(item, colWidths)
		case status.KindSecurity, status.KindIdentity:
			return m.renderCheckRowLoading(item, colWidths)
		case status.KindTool:
			return m.renderToolRowLoading(item, colWidths)
		case status.KindNetwork, status.KindCache:
			return m.renderResourceRowLoading(item, colWidths)
		default:
			return components.RowPrefix + m.spinner.View() + components.ColumnSeparator + components.CellNormal(item.Name, colWidths.Name)
		}
	}

	switch item.Kind {
	case status.KindSetup:
		return m.renderSetupRow(item, colWidths)
	case status.KindSecurity, status.KindIdentity:
		return m.renderCheckRow(item, colWidths)
	case status.KindTool:
		return m.renderToolRow(item, colWidths)
	case status.KindNetwork, status.KindCache:
		return m.renderResourceRow(item, colWidths)
	}

	return ""
}

func (m Model) renderSetupRowLoading(item status.Item, colWidths ColumnWidths) string {
	name := components.CellNormal(item.Name, colWidths.Name)
	desc := components.CellMuted(item.Description, colWidths.Desc)
	return components.RowPrefix + name + components.ColumnSeparator + desc + components.ColumnSeparator + m.spinner.View()
}

func (m Model) renderSetupRow(item status.Item, colWidths ColumnWidths) string {
	name := components.CellNormal(item.Name, colWidths.Name)
	desc := components.CellMuted(item.Description, colWidths.Desc)
	statusBadge := components.Badge(item.Installed)
	detail := ""
	if item.Detail != "" {
		detail = components.Muted(item.Detail)
	}
	return components.RowPrefix + name + components.ColumnSeparator + desc + components.ColumnSeparator + statusBadge + components.ColumnSeparator + detail
}

func (m Model) renderCheckRowLoading(item status.Item, colWidths ColumnWidths) string {
	name := components.CellNormal(item.Name, colWidths.Name)
	desc := components.CellMuted(item.Description, colWidths.Desc)
	return components.RowPrefix + name + components.ColumnSeparator + desc + components.ColumnSeparator + m.spinner.View()
}

func (m Model) renderCheckRow(item status.Item, colWidths ColumnWidths) string {
	name := components.CellNormal(item.Name, colWidths.Name)
	desc := components.CellMuted(item.Description, colWidths.Desc)
	ok := item.Installed == item.GoodWhen
	statusBadge := components.Badge(ok)
	detail := ""
	if item.Detail != "" {
		detail = components.Muted(item.Detail)
	}
	return components.RowPrefix + name + components.ColumnSeparator + desc + components.ColumnSeparator + statusBadge + components.ColumnSeparator + detail
}

func (m Model) renderToolRowLoading(item status.Item, colWidths ColumnWidths) string {
	name := components.CellNormal(item.Name, colWidths.Name)
	method := components.CellMethod(item.Method, colWidths.Method)
	return components.RowPrefix + name + components.ColumnSeparator + method + components.ColumnSeparator + m.spinner.View()
}

func (m Model) renderToolRow(item status.Item, colWidths ColumnWidths) string {
	name := components.CellNormal(item.Name, colWidths.Name)
	method := components.CellMethod(item.Method, colWidths.Method)
	statusBadge := components.Badge(item.Installed)
	version := components.CellSpecial(item.Version, colWidths.Version)

	// Show service status for docker/ollama
	extra := ""
	if item.Status != "" {
		if item.Status == "running" {
			extra = components.ColumnSeparator + theme.ServiceRunning.Render(theme.IconServiceOn) + " " + components.Success("running")
		} else if item.Status == "stopped" {
			extra = components.ColumnSeparator + components.Muted("stopped")
		} else {
			// Other status like "199 formulae, 6 casks" or "2 versions"
			extra = components.ColumnSeparator + components.Muted(item.Status)
		}
	}

	return components.RowPrefix + name + components.ColumnSeparator + method + components.ColumnSeparator + statusBadge + components.ColumnSeparator + version + extra
}

func (m Model) renderResourceRowLoading(item status.Item, colWidths ColumnWidths) string {
	name := components.CellNormal(item.Name, colWidths.Name)
	return components.RowPrefix + name + components.ColumnSeparator + m.spinner.View()
}

func (m Model) renderResourceRow(item status.Item, colWidths ColumnWidths) string {
	name := components.CellNormal(item.Name, colWidths.Name)
	value := components.StyledValue(item.Value, item.Style)
	return components.RowPrefix + name + components.ColumnSeparator + value
}

func (m Model) renderProcessRows(item status.Item) []string {
	// Add category header
	var header string
	if item.Name == "top cpu" {
		header = components.RowPrefix + components.Muted("CPU")
	} else if item.Name == "top memory" {
		header = components.RowPrefix + components.Muted("Memory")
	}

	if len(item.Processes) == 0 {
		if !item.Loaded {
			return []string{header + components.ColumnSeparator + m.spinner.View()}
		}
		return []string{header + components.ColumnSeparator + components.Muted("no data")}
	}

	var rows []string
	rows = append(rows, header)

	const nameWidth = 28

	for i, p := range item.Processes {
		if i >= 5 { // Show top 5
			break
		}
		// Truncate name if too long
		name := p.Name
		if len(name) > nameWidth {
			name = name[:nameWidth-3] + "..."
		}

		row := components.RowPrefix + components.CellNormal(name, nameWidth) + components.ColumnSeparator + components.CellRight(p.Value, 6)
		rows = append(rows, row)
	}
	return rows
}
