package ds

import (
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// AccordionItem is one section of an Accordion.
type AccordionItem struct {
	Title   string
	Content func() ui.Node
}

// AccordionProps configures an Accordion.
type AccordionProps struct {
	Items []AccordionItem
	// Open is the index of the open section; -1 closes all.
	Open *signal.Signal[int]
}

// Accordion is a stack of disclosure rows. Each trigger is a Pressable whose
// Hint says Expanded/Collapsed; content is built only while open (ui.Show),
// so its signals and effects are disposed when it closes.
func Accordion(p AccordionProps) ui.Node {
	kids := make([]ui.Node, 0, len(p.Items))
	for i, item := range p.Items {
		index, section := i, item
		expanded := signal.Func(func() bool { return p.Open != nil && p.Open.Get() == index })
		trigger := Pressable(PressableProps{
			Label:    section.Title,
			Expanded: expanded,
			OnPress: func() {
				if p.Open.Peek() == index {
					p.Open.Set(-1)
				} else {
					p.Open.Set(index)
				}
			},
			Style: themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{
					Direction: ui.RowDir, Align: ui.AlignCenter,
					Justify: ui.JustifySpaceBetween, Width: ui.Pct(100),
					Padding:    ui.Symmetric(Space3, Space4),
					Background: c.Card, CornerRadius: t.Radii.MD,
				}
			}),
		},
			ui.StaticText(section.Title, themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Font: t.Type.Subheading, Color: c.Foreground}
			})),
			ui.Symbol(ui.IconDown, themed(func(t Theme, c Palette) ui.Style {
				s := ui.Style{Font: t.Type.Caption, Color: c.MutedForeground}
				if expanded.Get() {
					s.Transform = ui.Transform{Rotate: 180}
				}
				return s
			})),
		)
		kids = append(kids, ui.Column(ui.Style{Gap: Space1},
			trigger,
			ui.Show(expanded, func() ui.Node {
				return ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
					return ui.Style{Padding: ui.All(Space4), Background: c.Card,
						CornerRadius: t.Radii.MD}
				})}, section.Content())
			}),
		))
	}
	return ui.Column(ui.Style{Gap: Space2}, kids...)
}

// TableColumn describes a Table column.
type TableColumn struct {
	Title string
	// Width is a fraction of the table's width; columns share what's left.
	Width float32
}

// TableProps configures a Table.
type TableProps struct {
	Columns []TableColumn
	Rows    [][]string
	// Striped alternates row backgrounds.
	Striped bool
}

// Table is a simple, non-virtualized data table. For thousands of rows use
// Asuka's ui.List, which virtualizes on every platform; Table is for dozens.
func Table(p TableProps) ui.Node {
	header := make([]ui.Node, 0, len(p.Columns))
	for _, col := range p.Columns {
		title := col.Title
		width := col.Width
		if width <= 0 {
			width = 1
		}
		header = append(header, ui.View(ui.ViewProps{Style: ui.Style{Grow: width}},
			ui.StaticText(title, themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Font: t.Type.Caption, Color: c.MutedForeground, Lines: 1}
			}))))
	}
	rows := make([]ui.Node, 0, len(p.Rows)+1)
	rows = append(rows, ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{Direction: ui.RowDir, Gap: Space3, Padding: ui.Symmetric(Space2, Space3),
			BorderWidth: 1, BorderColor: c.Border, Background: c.Muted,
			CornerRadius: 6, Width: ui.Pct(100)}
	})}, header...))

	for rowIndex, cells := range p.Rows {
		stripe := p.Striped && rowIndex%2 == 1
		cellNodes := make([]ui.Node, 0, len(p.Columns))
		for columnIndex, col := range p.Columns {
			text := ""
			if columnIndex < len(cells) {
				text = cells[columnIndex]
			}
			width := col.Width
			if width <= 0 {
				width = 1
			}
			cellNodes = append(cellNodes, ui.View(ui.ViewProps{Style: ui.Style{Grow: width}},
				ui.StaticText(text, themed(func(t Theme, c Palette) ui.Style {
					return ui.Style{Font: t.Type.Small, Color: c.Foreground, Lines: 1}
				}))))
		}
		rows = append(rows, ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
			s := ui.Style{Direction: ui.RowDir, Gap: Space3, Width: ui.Pct(100),
				Padding: ui.Symmetric(Space2, Space3), CornerRadius: 6}
			if stripe {
				s.Background = c.Muted
			}
			return s
		})}, cellNodes...))
	}
	return ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{Gap: Space1, Padding: ui.All(Space1), Width: ui.Pct(100),
			Background: c.Card, BorderWidth: 1, BorderColor: c.Border,
			CornerRadius: t.Radii.LG}
	})}, rows...)
}
