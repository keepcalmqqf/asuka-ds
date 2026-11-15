package ds

import (
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// ColorPickerProps configures a ColorPicker.
type ColorPickerProps struct {
	// Value is the chosen color; the platform's picker writes it.
	Value *signal.Signal[ui.Color]
	// Label is the visible name, and the picker's accessible title.
	Label string
	// Alpha lets the user choose an opacity.
	Alpha bool
	// Disabled dims the swatch and stops the picker from opening.
	Disabled signal.Accessor[bool]
}

// ColorPicker is a labeled swatch that opens the platform's color picker.
// Pair it with a themed label for a tint setting, an accent or a tag color.
func ColorPicker(p ColorPickerProps) ui.Node {
	value := p.Value
	if value == nil {
		value = signal.New(ui.Color(0))
	}
	swatch := ui.ColorPicker(ui.ColorPickerProps{
		Value:    value,
		Alpha:    p.Alpha,
		Title:    ui.Val(p.Label),
		Disabled: p.Disabled,
		Style:    ui.Style{Width: ui.Pt(28), Height: ui.Pt(28)},
	})
	if p.Label == "" {
		return swatch
	}
	return ui.Row(ui.Style{Align: ui.AlignCenter, Gap: Space3},
		ui.StaticText(p.Label, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Small, Color: c.Foreground, Lines: 1, Grow: 1}
		})),
		swatch,
	)
}
