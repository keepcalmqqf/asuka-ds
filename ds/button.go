package ds

import "github.com/egoist/asuka/ui"

// ButtonProps configures a Button.
type ButtonProps struct {
	Text string
	// Label is what a screen reader says; Text by default.
	Label string
	Icon  ui.Icon
	// IconAfter is a symbol after the text, for a "chevron.down" or "arrow.right".
	IconAfter ui.Icon
	Variant   Variant
	Size      Size
	Disabled  bool
	// FullWidth makes the button as wide as its parent.
	FullWidth bool
	OnPress   func()
	Href      ui.Href
}

// Button is the design system's button: a pressable View, drawn with the
// theme's tokens, so it looks the same on every platform. Use NativeButton
// where the platform's own button (and its behavior) matters more than the
// brand.
func Button(p ButtonProps) ui.Node {
	if p.Label == "" {
		p.Label = p.Text
	}
	height, padX := metrics(p.Size)

	fgStyle := themed(func(t Theme, c Palette) ui.Style {
		_, fg, _ := variantColors(c, p.Variant)
		return ui.Style{Color: fg, Font: fontFor(t.Type, p.Size)}
	})
	labelStyle := themed(func(t Theme, c Palette) ui.Style {
		_, fg, _ := variantColors(c, p.Variant)
		s := ui.Style{Color: fg, Font: fontFor(t.Type, p.Size)}
		if p.FullWidth && p.Icon == "" && p.IconAfter == "" {
			// A full-width button's text is otherwise sized to its own
			// intrinsic width; some hosts then wrap or clip a short word.
			s.Width = ui.Pct(100)
			s.TextAlign = ui.TextCenter
			s.Lines = 1
		}
		return s
	})
	var kids []ui.Node
	if p.Icon != "" {
		kids = append(kids, ui.Symbol(p.Icon, fgStyle))
	}
	if p.Text != "" {
		kids = append(kids, ui.Text(ui.Val(p.Text), labelStyle))
	}
	if p.IconAfter != "" {
		kids = append(kids, ui.Symbol(p.IconAfter, fgStyle))
	}

	props := ui.ViewProps{
		Label: ui.Val(p.Label),
		Style: themed(func(t Theme, c Palette) ui.Style {
			s := ui.Style{
				Direction:    ui.RowDir,
				Align:        ui.AlignCenter,
				Justify:      ui.JustifyCenter,
				Gap:          Space2,
				Padding:      ui.Symmetric(0, padX),
				Height:       ui.Pt(height),
				CornerRadius: t.Radii.MD,
				Font:         fontFor(t.Type, p.Size),
			}
			if p.FullWidth {
				s.Width = ui.Pct(100)
			}
			bg, _, border := variantColors(c, p.Variant)
			s.Background = bg
			if border != nil {
				s.BorderWidth = 1
				s.BorderColor = *border
			}
			if p.OnPress != nil && !p.Disabled {
				s.Hovered = &ui.Style{Background: hoverColors(c, p.Variant, bg)}
				s.Pressed = &ui.Style{Background: hoverColors(c, p.Variant, bg), Transform: ui.Transform{Scale: 0.98}}
			}
			if p.Disabled {
				s.Opacity = 0.5
			}
			return s
		}),
	}
	if p.OnPress != nil && !p.Disabled {
		props.OnPress = p.OnPress
		props.Href = p.Href
	} else {
		// A disabled button keeps its role and name, so it's still
		// announced; without OnPress it isn't actionable.
		props.Role = ui.RoleButton
	}
	if p.Disabled {
		props.Hint = ui.Val("Unavailable")
	}
	return ui.View(props, kids...)
}

// NativeButton is the platform's own button with the design system's type and
// tint: the right choice on iOS and Android when you want the native press
// behavior, the ripple or the keyboard focus ring.
func NativeButton(p ButtonProps) ui.Node {
	variant := ui.ButtonPlain
	switch p.Variant {
	case Primary:
		variant = ui.ButtonFilled
	case Secondary:
		variant = ui.ButtonGray
	case Outline:
		variant = ui.ButtonBordered
	case Destructive:
		variant = ui.ButtonTinted
	}
	props := ui.ButtonProps{
		Title:    ui.Val(p.Text),
		Icon:     p.Icon,
		Variant:  variant,
		Disabled: ui.Val(p.Disabled),
		OnPress:  p.OnPress,
		Href:     p.Href,
		Style: themed(func(t Theme, c Palette) ui.Style {
			s := ui.Style{Font: fontFor(t.Type, p.Size)}
			if p.Variant == Destructive {
				s.Color = c.Destructive
			}
			return s
		}),
	}
	if p.Label != "" {
		props.Label = ui.Val(p.Label)
	}
	return ui.Button(props)
}
