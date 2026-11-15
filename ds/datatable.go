package ds

import (
	"sort"
	"strconv"
	"strings"

	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// DataTableColumn describes a DataTable column.
type DataTableColumn struct {
	Title string
	// Width is a fraction of the table's width; columns share what's left.
	Width    float32
	Sortable bool
}

// DataTableProps configures a DataTable.
type DataTableProps struct {
	Columns []DataTableColumn
	Rows    [][]string
	// Sort is the sorted column's index (-1: none); the header writes it.
	Sort *signal.Signal[int]
	// Descending is the sort's direction; a second press on the sorted
	// column reverses it.
	Descending *signal.Signal[bool]
	Striped    bool
	// OnRowPress makes each row pressable; index is the row's place in
	// Rows, before sorting.
	OnRowPress func(index int)
	// RowLabel is a row's accessible name; the cells joined with ", " by
	// default.
	RowLabel func(cells []string) string
}

// DataTable is a Table with sortable headers and pressable rows. It is
// deliberately not virtualized: for dozens of rows it is fine; for
// thousands use Asuka's ui.List, which virtualizes on every platform.
func DataTable(p DataTableProps) ui.Node {
	sortColumn := p.Sort
	if sortColumn == nil {
		sortColumn = signal.New(-1)
	}
	descending := p.Descending
	if descending == nil {
		descending = signal.New(false)
	}
	rowLabel := p.RowLabel
	if rowLabel == nil {
		rowLabel = func(cells []string) string { return strings.Join(cells, ", ") }
	}
	return ui.Dynamic(func() ui.Node {
		column, down := sortColumn.Get(), descending.Get()
		rows := orderedRows(p.Rows, column, down)

		header := make([]ui.Node, 0, len(p.Columns))
		for index, col := range p.Columns {
			title, i := col.Title, index
			width := col.Width
			if width <= 0 {
				width = 1
			}
			text := func() ui.Node {
				kids := []ui.Node{ui.StaticText(title, themed(func(t Theme, c Palette) ui.Style {
					return ui.Style{Font: t.Type.Caption, Color: c.MutedForeground, Lines: 1}
				}))}
				if col.Sortable && column == i {
					icon := ui.IconUp
					if down {
						icon = ui.IconDown
					}
					kids = append(kids, ui.Symbol(icon, themed(func(t Theme, c Palette) ui.Style {
						return ui.Style{Font: t.Type.Caption, Color: c.MutedForeground}
					})))
				}
				return ui.Row(ui.Style{Gap: Space1, Align: ui.AlignCenter}, kids...)
			}
			if !col.Sortable {
				header = append(header, ui.View(ui.ViewProps{Style: ui.Style{Grow: width}}, text()))
				continue
			}
			label := "Sort by " + title
			if column == i {
				if down {
					label += ", descending"
				} else {
					label += ", ascending"
				}
			}
			header = append(header, Pressable(PressableProps{
				Label:    label,
				Selected: ui.Val(column == i),
				OnPress: func() {
					if column == i {
						descending.Set(!descending.Peek())
						return
					}
					sortColumn.Set(i)
					descending.Set(false)
				},
				Style: themed(func(t Theme, c Palette) ui.Style {
					return ui.Style{Grow: width, Padding: ui.Symmetric(2, 0),
						CornerRadius: 4, Hovered: &ui.Style{Background: c.Border}}
				}),
			}, text()))
		}

		nodes := make([]ui.Node, 0, len(rows)+1)
		nodes = append(nodes, ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Direction: ui.RowDir, Gap: Space3, Padding: ui.Symmetric(Space2, Space3),
				BorderWidth: 1, BorderColor: c.Border, Background: c.Muted,
				CornerRadius: 6, Width: ui.Pct(100)}
		})}, header...))
		for _, row := range rows {
			entry := row
			cells := make([]ui.Node, 0, len(p.Columns))
			for index, col := range p.Columns {
				value := ""
				if index < len(entry.cells) {
					value = entry.cells[index]
				}
				width := col.Width
				if width <= 0 {
					width = 1
				}
				cells = append(cells, ui.View(ui.ViewProps{Style: ui.Style{Grow: width}},
					ui.StaticText(value, themed(func(t Theme, c Palette) ui.Style {
						return ui.Style{Font: t.Type.Small, Color: c.Foreground, Lines: 1}
					}))))
			}
			style := themed(func(t Theme, c Palette) ui.Style {
				s := ui.Style{Direction: ui.RowDir, Gap: Space3, Width: ui.Pct(100),
					Padding: ui.Symmetric(Space2, Space3), CornerRadius: 6}
				if p.Striped && entry.index%2 == 1 {
					s.Background = c.Muted
				}
				if p.OnRowPress != nil {
					s.Hovered = &ui.Style{Background: c.Accent}
				}
				return s
			})
			if p.OnRowPress == nil {
				nodes = append(nodes, ui.View(ui.ViewProps{Style: style}, cells...))
				continue
			}
			index := entry.index
			nodes = append(nodes, Pressable(PressableProps{
				Label:   rowLabel(entry.cells),
				OnPress: func() { p.OnRowPress(index) },
				Style:   style,
			}, cells...))
		}
		return ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Gap: Space1, Padding: ui.All(Space1), Width: ui.Pct(100),
				Background: c.Card, BorderWidth: 1, BorderColor: c.Border,
				CornerRadius: t.Radii.LG}
		})}, nodes...)
	})
}

type tableRow struct {
	cells []string
	index int
}

func orderedRows(rows [][]string, column int, descending bool) []tableRow {
	out := make([]tableRow, 0, len(rows))
	for index, cells := range rows {
		out = append(out, tableRow{cells: cells, index: index})
	}
	if column < 0 {
		return out
	}
	sort.SliceStable(out, func(i, j int) bool {
		cmp := compareCells(out[i].cells, out[j].cells, column)
		if descending {
			return cmp > 0
		}
		return cmp < 0
	})
	return out
}

func compareCells(a, b []string, column int) int {
	left, right := "", ""
	if column < len(a) {
		left = a[column]
	}
	if column < len(b) {
		right = b[column]
	}
	lf, lerr := strconv.ParseFloat(strings.TrimSpace(left), 64)
	rf, rerr := strconv.ParseFloat(strings.TrimSpace(right), 64)
	if lerr == nil && rerr == nil {
		switch {
		case lf < rf:
			return -1
		case lf > rf:
			return 1
		default:
			return 0
		}
	}
	return strings.Compare(strings.ToLower(left), strings.ToLower(right))
}
