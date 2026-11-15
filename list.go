package ds

import (
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// ListProps configures a List.
type ListProps struct {
	// Divided draws a hairline between rows.
	Divided bool
	// Flat drops the surface's border and shadow.
	Flat bool
}

// List stacks ListItems on a card-like surface. Rows own their padding and
// their press behavior, so a list can mix settings rows, switches and plain
// information.
func List(p ListProps, children ...ui.Node) ui.Node {
	rows := make([]ui.Node, 0, len(children)*2)
	for i, child := range children {
		if i > 0 && p.Divided {
			rows = append(rows, listDivider())
		}
		rows = append(rows, child)
	}
	return ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
		s := ui.Style{
			Background:   c.Card,
			CornerRadius: t.Radii.LG,
			Shadow:       c.Elevation2.Shadow,
			Clip:         true,
		}
		if !p.Flat {
			s.BorderWidth = 1
			s.BorderColor = c.Border
		}
		return s
	})}, ui.Column(ui.Style{}, rows...))
}

// ListItemProps configures a ListItem.
type ListItemProps struct {
	Title       string
	Description string
	// Icon is shown in a small tile before the title.
	Icon ui.Icon
	// Detail is muted text on the trailing edge, before the chevron.
	Detail string
	// Chevron shows that the row opens something.
	Chevron bool
	// Selected tints the row and is announced.
	Selected signal.Accessor[bool]
	Disabled signal.Accessor[bool]
	// OnPress makes the row a button.
	OnPress func()
	// Action is a trailing control (a Switch, a Badge, a value). It is
	// outside the row's press target, so it keeps its own behavior.
	Action ui.Node
}

// ListItem is a settings-style row: an optional icon, a title and
// description, trailing detail and a chevron.
func ListItem(p ListItemProps) ui.Node {
	leading := make([]ui.Node, 0, 3)
	if p.Icon != "" {
		leading = append(leading, ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{
				Direction:    ui.RowDir,
				Align:        ui.AlignCenter,
				Justify:      ui.JustifyCenter,
				Width:        ui.Pt(32),
				Height:       ui.Pt(32),
				CornerRadius: t.Radii.MD,
				Background:   c.Muted,
			}
		})}, ui.Symbol(p.Icon, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Small, Color: c.Foreground}
		}))))
	}
	text := []ui.Node{ui.StaticText(p.Title, themed(func(t Theme, c Palette) ui.Style {
		s := ui.Style{Font: t.Type.Subheading, Color: c.Foreground, Lines: 1}
		if p.Disabled != nil && p.Disabled.Get() {
			s.Color = c.MutedForeground
		}
		return s
	}))}
	if p.Description != "" {
		text = append(text, ui.StaticText(p.Description, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Small, Color: c.MutedForeground, Lines: 2}
		})))
	}
	leading = append(leading, ui.View(ui.ViewProps{Style: ui.Style{Grow: 1}},
		ui.Column(ui.Style{Gap: 2}, text...)))
	if p.Detail != "" {
		leading = append(leading, ui.StaticText(p.Detail, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Small, Color: c.MutedForeground, Lines: 1}
		})))
	}
	if p.Chevron && p.OnPress != nil {
		leading = append(leading, ui.Symbol(ui.IconForward, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Caption, Color: c.MutedForeground}
		})))
	}

	rowStyle := themed(func(t Theme, c Palette) ui.Style {
		s := ui.Style{
			Direction: ui.RowDir,
			Align:     ui.AlignCenter,
			Gap:       Space3,
			Padding:   ui.Symmetric(Space3, Space4),
			Grow:      1,
		}
		if p.Selected != nil && p.Selected.Get() {
			s.Background = c.Accent
		} else {
			s.Hovered = &ui.Style{Background: c.Muted}
		}
		return s
	})

	var body ui.Node
	if p.OnPress != nil {
		body = Pressable(PressableProps{
			Label:    p.Title,
			Hint:     p.Description,
			Disabled: p.Disabled,
			Selected: p.Selected,
			Style:    rowStyle,
			OnPress:  p.OnPress,
		}, leading...)
	} else {
		body = ui.View(ui.ViewProps{
			Label: ui.Val(p.Title),
			Style: rowStyle,
		}, leading...)
	}

	if p.Action == (ui.Node{}) {
		return body
	}
	return ui.Row(ui.Style{Align: ui.AlignCenter}, body,
		ui.View(ui.ViewProps{Style: ui.Style{Padding: ui.Edges{Right: Space4}}}, p.Action))
}

func listDivider() ui.Node {
	return ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{Height: ui.Pt(1), Width: ui.Pct(100), Background: c.Border}
	})})
}
