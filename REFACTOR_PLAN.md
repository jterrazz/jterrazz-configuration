# Presentation Layer Refactoring Plan

## Overview

This plan consolidates duplicated presentation logic into common components to improve maintainability and consistency.

## Priority 1: Badge Consolidation (High Impact, Low Effort)

### Current State
- `components/badge.go` has `BadgeOK()`, `BadgeError()`, `Badge()`, `BadgeLoading()`, `ServiceBadge()`
- `theme/icons.go` has duplicate `StatusIcon()` and `ServiceIcon()` functions
- `views/status/sections.go` has local `badge()` function duplicating `Badge()`

### Migration Steps

1. **Remove duplicate from sections.go**
   - File: `src/internal/presentation/views/status/sections.go`
   - Delete local `badge()` function (lines 178-183)
   - Import and use `components.Badge()` instead

2. **Remove duplicates from icons.go**
   - File: `src/internal/presentation/theme/icons.go`
   - Delete `StatusIcon()` function (lines 65-70)
   - Delete `ServiceIcon()` function (lines 73-78)
   - These are now in `components/badge.go`

---

## Priority 2: Cell Rendering Consistency (High Impact, Medium Effort)

### Current State
- `components/text.go` has `CellNormal()`, `CellMuted()`, `CellMethod()`, `CellSpecial()`
- `views/status/rows.go` uses inline `fmt.Sprintf("%-*s", width, text)` with `theme.X.Render()` instead

### Migration Steps

1. **Update rows.go to use Cell functions**
   - File: `src/internal/presentation/views/status/rows.go`
   - Replace all inline column formatting:
     ```go
     // Before
     name := theme.Cell.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))

     // After
     name := text.CellNormal(item.Name, colWidths.Name)
     ```
   - Apply to all render functions

2. **Add missing Cell function**
   - File: `src/internal/presentation/components/text.go`
   - Add `CellSpinner()` for loading states

---

## Priority 3: Row Prefix Extraction (Medium Impact, Low Effort)

### Current State
- `tui_item.go` has same prefix logic repeated 4 times in renderNavigation, renderAction, renderToggle, renderExpandable

### Migration Steps

1. **Extract prefix builder**
   - File: `src/internal/presentation/components/tui_item.go`
   - Add helper function:
     ```go
     func (i Item) buildPrefix(selected bool) string {
         indent := i.indentPrefix()
         if selected {
             return indent + " " + theme.IconSelected + " "
         }
         return indent + "   "
     }
     ```
   - Replace 4 duplicate blocks with `prefix := i.buildPrefix(selected)`

---

## Priority 4: Spacing Constants (Medium Impact, Low Effort)

### Current State
- Magic strings like `"  "` (column separator) and `" "` (row prefix) scattered across files

### Migration Steps

1. **Define spacing constants**
   - File: `src/internal/presentation/theme/spacing.go` (new file)
   ```go
   package theme

   const (
       ColumnSeparator = "  "  // 2 spaces between columns
       RowPrefix       = " "   // 1 space at row start (inside box)
       ItemPrefix      = "   " // 3 spaces for list items (aligns with section headers)
   )
   ```

2. **Update usages across files**
   - `views/status/rows.go`: Replace `" %s  %s"` with constants
   - `components/tui_item.go`: Replace `"   "` with `theme.ItemPrefix`
   - `components/box.go`: Use constants for padding

---

## Priority 5: Row Builder Pattern (High Impact, High Effort)

### Current State
- `views/status/rows.go` has 8 similar functions with repetitive structure
- Each function: extract fields, pad columns, combine with fmt.Sprintf

### Migration Steps

1. **Create RowBuilder component**
   - File: `src/internal/presentation/components/row_builder.go` (new file)
   ```go
   package components

   type RowBuilder struct {
       columns []string
   }

   func NewRow() *RowBuilder {
       return &RowBuilder{}
   }

   func (r *RowBuilder) AddCell(text string, width int, style Style) *RowBuilder {
       // Format and style the cell
       return r
   }

   func (r *RowBuilder) AddBadge(ok bool) *RowBuilder {
       // Add status badge
       return r
   }

   func (r *RowBuilder) AddSpinner(frame string) *RowBuilder {
       // Add loading spinner
       return r
   }

   func (r *RowBuilder) Build() string {
       return RowPrefix + strings.Join(r.columns, ColumnSeparator)
   }
   ```

2. **Refactor rows.go to use RowBuilder**
   ```go
   // Before
   func (m Model) renderSetupRow(item status.Item, colWidths ColumnWidths) string {
       name := theme.Cell.Render(fmt.Sprintf("%-*s", colWidths.Name, item.Name))
       desc := theme.Muted.Render(fmt.Sprintf("%-*s", colWidths.Desc, item.Description))
       statusBadge := badge(item.Installed)
       detail := ""
       if item.Detail != "" {
           detail = theme.Muted.Render(item.Detail)
       }
       return fmt.Sprintf(" %s  %s  %s  %s", name, desc, statusBadge, detail)
   }

   // After
   func (m Model) renderSetupRow(item status.Item, colWidths ColumnWidths) string {
       return NewRow().
           AddCell(item.Name, colWidths.Name, StyleNormal).
           AddCell(item.Description, colWidths.Desc, StyleMuted).
           AddBadge(item.Installed).
           AddCell(item.Detail, 0, StyleMuted).
           Build()
   }
   ```

---

## Priority 6: Description Rendering Helper (Low Impact, Low Effort)

### Current State
- Pattern `theme.Muted.Render("  "+description)` repeated in tui_item.go

### Migration Steps

1. **Add description helper**
   - File: `src/internal/presentation/components/text.go`
   ```go
   func RenderDescription(text string) string {
       if text == "" {
           return ""
       }
       return theme.Muted.Render(ColumnSeparator + text)
   }
   ```

2. **Update tui_item.go**
   - Replace all `theme.Muted.Render("  "+i.Description)` with `text.RenderDescription(i.Description)`

---

## Priority 7: Box Line Padding (Low Impact, Low Effort)

### Current State
- `components/box.go` has identical padding logic in SubsectionBox and SimpleBox

### Migration Steps

1. **Extract padBoxLine helper**
   - File: `src/internal/presentation/components/box.go`
   ```go
   func padBoxLine(line string, innerWidth int, borderStyle lipgloss.Style) string {
       padding := innerWidth - VisibleLen(line)
       if padding < 0 {
           padding = 0
       }
       return borderStyle.Render(theme.BoxRoundedVertical+" ") +
              line +
              strings.Repeat(" ", padding) +
              borderStyle.Render(" "+theme.BoxRoundedVertical)
   }
   ```

2. **Use in SubsectionBox and SimpleBox**

---

## Implementation Order

1. **Phase 1 - Quick Wins** (can be done immediately)
   - Priority 1: Badge consolidation
   - Priority 3: Row prefix extraction
   - Priority 6: Description helper
   - Priority 7: Box line padding

2. **Phase 2 - Cell Functions** (after Phase 1)
   - Priority 2: Use existing Cell functions in rows.go
   - Priority 4: Spacing constants

3. **Phase 3 - Row Builder** (larger refactor)
   - Priority 5: Implement RowBuilder pattern
   - Refactor all row rendering functions

---

## Files Summary

### Files to Modify
- `src/internal/presentation/views/status/sections.go` - Remove local badge()
- `src/internal/presentation/views/status/rows.go` - Use Cell functions, then RowBuilder
- `src/internal/presentation/theme/icons.go` - Remove StatusIcon, ServiceIcon
- `src/internal/presentation/components/tui_item.go` - Extract prefix builder, use description helper
- `src/internal/presentation/components/box.go` - Extract padBoxLine
- `src/internal/presentation/components/text.go` - Add RenderDescription

### New Files
- `src/internal/presentation/theme/spacing.go` - Spacing constants
- `src/internal/presentation/components/row_builder.go` - Row builder pattern (Phase 3)

---

## Expected Benefits

1. **Reduced code duplication** - ~30% less repeated formatting logic
2. **Consistent styling** - All badges, cells, rows use same components
3. **Easier maintenance** - Change spacing/formatting in one place
4. **Better testability** - Individual components can be unit tested
5. **Clearer intent** - `AddBadge(ok)` vs `theme.BadgeOK.Render(theme.IconCheck)`
