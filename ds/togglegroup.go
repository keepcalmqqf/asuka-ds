package ds

import (
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// ToggleItem is one choice of a ToggleGroup.
type ToggleItem struct {
	Value    string
	Label    string
	Icon     ui.Icon
	Disabled bool
}

// ToggleGroupProps configures a ToggleGroup.
type ToggleGroupProps struct {
	Label string
	Items []ToggleItem
	// Selected holds the chosen values. Single keeps at most one.
	// A slice isn't comparable, so build it with
	// signal.NewFunc(..., slices.Equal).
	Selected *signal.Signal[[]string]
	Single   bool
	Disabled bool
}

// ToggleGroup is a row of Toggles sharing one selection. Single behaves
// like a segmented control drawn from Pressables; multiple is a set of
// independent switches (text formatting, filters).
func ToggleGroup(p ToggleGroupProps) ui.Node {
	has := func(value string) bool {
		if p.Selected == nil {
			return false
		}
		for _, v := range p.Selected.Peek() {
			if v == value {
				return true
			}
		}
		return false
	}
	toggle := func(value string) {
		if p.Selected == nil {
			return
		}
		if p.Single {
			if has(value) {
				p.Selected.Set(nil)
			} else {
				p.Selected.Set([]string{value})
			}
			return
		}
		cur := append([]string(nil), p.Selected.Peek()...)
		if has(value) {
			out := cur[:0]
			for _, v := range cur {
				if v != value {
					out = append(out, v)
				}
			}
			p.Selected.Set(out)
			return
		}
		p.Selected.Set(append(cur, value))
	}
	control := func() ui.Node {
		kids := make([]ui.Node, 0, len(p.Items))
		for _, item := range p.Items {
			entry := item
			active := signal.Func(func() bool {
				if p.Selected == nil {
					return false
				}
				for _, v := range p.Selected.Get() {
					if v == entry.Value {
						return true
					}
				}
				return false
			})
			var content []ui.Node
			if entry.Icon != "" {
				content = append(content, ui.Symbol(entry.Icon, themed(func(t Theme, c Palette) ui.Style {
					col := c.Foreground
					if active.Get() {
						col = c.PrimaryForeground
					}
					return ui.Style{Font: t.Type.Small, Color: col}
				})))
			}
			content = append(content, ui.StaticText(entry.Label, themed(func(t Theme, c Palette) ui.Style {
				col := c.Foreground
				if active.Get() {
					col = c.PrimaryForeground
				}
				return ui.Style{Font: t.Type.Small, Color: col, Lines: 1}
			})))
			kids = append(kids, Pressable(PressableProps{
				Label:    entry.Label,
				Selected: active,
				Disabled: ui.Val(p.Disabled || entry.Disabled),
				OnPress:  func() { toggle(entry.Value) },
				Style: themed(func(t Theme, c Palette) ui.Style {
					s := ui.Style{
						Direction: ui.RowDir, Align: ui.AlignCenter, Gap: Space2,
						Padding: ui.Symmetric(Space2, Space3), CornerRadius: t.Radii.MD,
						Background: c.Background, BorderWidth: 1, BorderColor: c.Border,
					}
					if active.Get() {
						s.Background = c.Primary
						s.BorderColor = c.Primary
					} else {
						s.Hovered = &ui.Style{Background: c.Muted}
					}
					return s
				}),
			}, content...))
		}
		return ui.Row(ui.Style{Gap: Space1, Wrap: true, Align: ui.AlignCenter}, kids...)
	}
	return Field(FieldProps{
		Label:    p.Label,
		Disabled: ui.Val(p.Disabled),
		Control:  control,
	})
}
