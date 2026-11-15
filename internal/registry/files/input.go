package ds

import (
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// InputProps configures an Input.
type InputProps struct {
	// Value is the field's text; the widget writes it as the user types.
	Value       *signal.Signal[string]
	Placeholder string
	Label       string
	Description string
	Error       string
	Icon        ui.Icon
	Disabled    bool
	Multiline   bool
	Lines       int
}

// Input is a text field with a label, a description and an error, drawn as a
// bordered box. The field itself is the platform's own (so IME, autofill,
// selection and screen readers behave), with its own chrome turned off
// (TextFieldProps.Plain) and the design system's box around it.
func Input(p InputProps) ui.Node {
	focused := signal.New(false)
	control := func() ui.Node {
		box := themed(func(t Theme, c Palette) ui.Style {
			s := ui.Style{
				Direction:    ui.RowDir,
				Align:        ui.AlignCenter,
				Gap:          Space2,
				Padding:      ui.Symmetric(0, Space3),
				MinHeight:    ui.Pt(36),
				Background:   c.Background,
				BorderWidth:  1,
				BorderColor:  c.Input,
				CornerRadius: t.Radii.MD,
			}
			if p.Multiline {
				s.Padding = ui.Symmetric(Space2, Space3)
			}
			if focused.Get() {
				s.BorderColor = c.Focus
			}
			if p.Error != "" {
				s.BorderColor = c.Destructive
			}
			if p.Disabled {
				s.Background = c.Muted
				s.BorderColor = c.Border
			}
			return s
		})
		fieldStyle := themed(func(t Theme, c Palette) ui.Style {
			col := c.Foreground
			if p.Disabled {
				col = c.MutedForeground
			}
			return ui.Style{Grow: 1, Font: t.Type.Body, Color: col, Lines: p.Lines}
		})
		field := ui.TextField(ui.TextFieldProps{
			Value:       p.Value,
			Placeholder: ui.Val(p.Placeholder),
			Plain:       true,
			Multiline:   p.Multiline,
			Disabled:    ui.Val(p.Disabled),
			OnFocus:     func() { focused.Set(true) },
			OnBlur:      func() { focused.Set(false) },
			Style:       fieldStyle,
		})
		kids := []ui.Node{field}
		if p.Icon != "" {
			kids = append([]ui.Node{ui.Symbol(p.Icon, themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Font: t.Type.Small, Color: c.MutedForeground}
			}))}, kids...)
		}
		return ui.View(ui.ViewProps{Style: box}, kids...)
	}
	return Field(FieldProps{
		Label:                p.Label,
		Description:          p.Description,
		DescriptionWithError: true,
		Error:                p.Error,
		Disabled:             ui.Val(p.Disabled),
		Control:              control,
	})
}
