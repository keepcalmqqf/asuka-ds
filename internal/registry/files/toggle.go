package ds

import (
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// ToggleProps configures a Toggle.
type ToggleProps struct {
	Text string
	// Label is what a screen reader says; Text by default.
	Label    string
	Icon     ui.Icon
	Pressed  *signal.Signal[bool]
	Variant  Variant
	Size     Size
	Disabled bool
}

// Toggle is a button that stays pressed: a formatting control, a filter, a
// view switch. Its pressed state is announced through the Pressable
// Selected hint, since a View has no selected trait of its own.
func Toggle(p ToggleProps) ui.Node {
	if p.Label == "" {
		p.Label = p.Text
	}
	height, padX := metrics(p.Size)
	selected := func() bool { return p.Pressed != nil && p.Pressed.Get() }
	var selectedAccessor signal.Accessor[bool]
	if p.Pressed != nil {
		selectedAccessor = signal.Func(selected)
	}
	var content []ui.Node
	if p.Icon != "" {
		content = append(content, ui.Symbol(p.Icon, themed(func(t Theme, c Palette) ui.Style {
			_, fg, _ := variantColors(c, p.Variant)
			if selected() {
				fg = c.AccentForeground
			}
			return ui.Style{Font: fontFor(t.Type, p.Size), Color: fg}
		})))
	}
	if p.Text != "" {
		content = append(content, ui.StaticText(p.Text, themed(func(t Theme, c Palette) ui.Style {
			_, fg, _ := variantColors(c, p.Variant)
			if selected() {
				fg = c.AccentForeground
			}
			return ui.Style{Font: fontFor(t.Type, p.Size), Color: fg, Lines: 1}
		})))
	}
	return Pressable(PressableProps{
		Label:    p.Label,
		Disabled: ui.Val(p.Disabled),
		Selected: selectedAccessor,
		OnPress: func() {
			if p.Pressed != nil {
				p.Pressed.Set(!p.Pressed.Peek())
			}
		},
		Style: themed(func(t Theme, c Palette) ui.Style {
			bg, _, border := variantColors(c, p.Variant)
			if selected() {
				bg, _, border = c.Accent, c.AccentForeground, nil
			}
			s := ui.Style{
				Direction: ui.RowDir, Align: ui.AlignCenter,
				Justify: ui.JustifyCenter, Gap: Space2,
				Padding: ui.Symmetric(0, padX), Height: ui.Pt(height),
				Background: bg, CornerRadius: t.Radii.MD,
			}
			if border != nil {
				s.BorderWidth = 1
				s.BorderColor = *border
			}
			if !selected() {
				s.Hovered = &ui.Style{Background: c.Muted}
			}
			return s
		}),
	}, content...)
}
