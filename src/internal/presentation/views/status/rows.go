package status

import (
	"fmt"

	"github.com/jterrazz/jterrazz-cli/internal/domain/status"
	"github.com/jterrazz/jterrazz-cli/internal/presentation/theme"
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
		case status.KindNetwork, status.KindDisk, status.KindCache:
			return m.renderResourceRowLoading(item, colWidths)
		default:
			return fmt.Sprintf("  %s  %-*s", m.spinner.View(), colWidths.Name, item.Name)
		}
	}

	switch item.Kind {
	case status.KindSetup:
		return m.renderSetupRow(item, colWidths)
	case status.KindSecurity, status.KindIdentity:
		return m.renderCheckRow(item, colWidths)
	case status.KindTool:
		return m.renderToolRow(item, colWidths)
	case status.KindNetwork, status.KindDisk, status.KindCache:
		return m.renderResourceRow(item, colWidths)
	}

	return ""
}

func (m Model) renderSetupRowLoading(item status.Item, colWidths ColumnWidths) string {
	name := theme.Cell.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))
	desc := theme.Muted.Render(fmt.Sprintf("%-*s", colWidths.Desc, item.Description))
	return fmt.Sprintf("  %s  %s  %s", name, desc, m.spinner.View())
}

func (m Model) renderSetupRow(item status.Item, colWidths ColumnWidths) string {
	name := theme.Cell.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))
	desc := theme.Muted.Render(fmt.Sprintf("%-*s", colWidths.Desc, item.Description))
	statusBadge := badge(item.Installed)
	detail := ""
	if item.Detail != "" {
		detail = theme.Special.Render(item.Detail)
	}
	return fmt.Sprintf("  %s  %s  %s  %s", name, desc, statusBadge, detail)
}

func (m Model) renderCheckRowLoading(item status.Item, colWidths ColumnWidths) string {
	name := theme.Cell.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))
	desc := theme.Muted.Render(fmt.Sprintf("%-*s", colWidths.Desc, item.Description))
	return fmt.Sprintf("  %s  %s  %s", name, desc, m.spinner.View())
}

func (m Model) renderCheckRow(item status.Item, colWidths ColumnWidths) string {
	name := theme.Cell.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))
	desc := theme.Muted.Render(fmt.Sprintf("%-*s", colWidths.Desc, item.Description))
	ok := item.Installed == item.GoodWhen
	statusBadge := badge(ok)
	detail := ""
	if item.Detail != "" {
		detail = theme.Special.Render(item.Detail)
	}
	return fmt.Sprintf("  %s  %s  %s  %s", name, desc, statusBadge, detail)
}

func (m Model) renderToolRowLoading(item status.Item, colWidths ColumnWidths) string {
	name := theme.Cell.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))
	method := theme.Method.Render(fmt.Sprintf("%-*s", colWidths.Method, item.Method))
	return fmt.Sprintf("  %s  %s  %s", name, method, m.spinner.View())
}

func (m Model) renderToolRow(item status.Item, colWidths ColumnWidths) string {
	name := theme.Cell.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))

	// Method very dim (least important)
	method := theme.Method.Render(fmt.Sprintf("%-*s", colWidths.Method, item.Method))

	statusBadge := badge(item.Installed)

	// Version in cyan if present, dimmed dash if not
	versionStr := fmt.Sprintf("%-*s", colWidths.Version, item.Version)
	var version string
	if item.Version != "" {
		version = theme.Special.Render(versionStr)
	} else {
		version = theme.Muted.Render(versionStr)
	}

	// Show service status for docker/ollama
	extra := ""
	if item.Status != "" {
		if item.Status == "running" {
			extra = "  " + theme.ServiceRunning.Render(theme.IconServiceOn) + " " + theme.Success.Render("running")
		} else if item.Status == "stopped" {
			extra = "  " + theme.ServiceStopped.Render(theme.IconServiceOff) + " " + theme.Warning.Render("stopped")
		} else {
			// Other status like "199 formulae, 6 casks" or "2 versions"
			extra = "  " + theme.Muted.Render(item.Status)
		}
	}

	return fmt.Sprintf("  %s  %s  %s  %s%s", name, method, statusBadge, version, extra)
}

func (m Model) renderResourceRowLoading(item status.Item, colWidths ColumnWidths) string {
	name := theme.Muted.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))
	return fmt.Sprintf("  %s  %s", name, m.spinner.View())
}

func (m Model) renderResourceRow(item status.Item, colWidths ColumnWidths) string {
	name := theme.Muted.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))
	value := theme.Render(item.Value, item.Style)
	return fmt.Sprintf("  %s  %s", name, value)
}
