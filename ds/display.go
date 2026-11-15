package ds

import (
	"fmt"
	"strings"

	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// AvatarProps configures an Avatar.
type AvatarProps struct {
	// Name is shown as initials when there is no image. It is also the
	// accessible label when the avatar has meaning.
	Name string
	// Source is an image URL or an asset path (see ui.Image).
	Source string
	Size   float32 // 0: 40
}

// Avatar is a round image with initials as its fallback.
func Avatar(p AvatarProps) ui.Node {
	size := p.Size
	if size <= 0 {
		size = 40
	}
	label := p.Name
	if label == "" {
		label = "Avatar"
	}
	var content ui.Node
	if p.Source != "" {
		content = ui.Image(ui.ImageProps{
			Source: ui.Val(p.Source), Label: ui.Val(label), Mode: ui.ImageFill,
			Style: ui.Style{Width: ui.Pt(size), Height: ui.Pt(size)},
		})
	} else {
		content = ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Align: ui.AlignCenter, Justify: ui.JustifyCenter,
				Width: ui.Pt(size), Height: ui.Pt(size), CornerRadius: size / 2, Background: c.Muted}
		})}, ui.StaticText(initials(p.Name), themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Subheading, Color: c.MutedForeground, Lines: 1}
		})))
	}
	return ui.View(ui.ViewProps{
		Label: ui.Val(label), Role: ui.RoleImage,
		Style: ui.Style{Width: ui.Pt(size), Height: ui.Pt(size), CornerRadius: size / 2, Clip: true},
	}, content)
}

func initials(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "?"
	}
	if len(parts) == 1 {
		return strings.ToUpper(parts[0][:1])
	}
	return strings.ToUpper(parts[0][:1] + parts[len(parts)-1][:1])
}

// ProgressProps configures a Progress bar.
type ProgressProps struct {
	// Value is 0..1. A nil or negative value makes an indeterminate bar.
	Value signal.Accessor[float64]
	Label string
	// ShowValue prints the percentage at the trailing edge.
	ShowValue bool
}

// Progress is a themed progress bar. The bar itself is the platform's own
// progress view, so it is announced correctly; the track color comes from
// the theme.
func Progress(p ProgressProps) ui.Node {
	bar := ui.Progress(p.Value, themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{Height: ui.Pt(8), CornerRadius: 4, Color: c.Primary, Grow: 1}
	}))
	if p.Label == "" && !p.ShowValue {
		return bar
	}
	var trail []ui.Node
	if p.ShowValue {
		trail = append(trail, ui.Text(signal.Func(func() string {
			if p.Value == nil {
				return ""
			}
			v := p.Value.Get()
			if v < 0 {
				return "…"
			}
			return fmt.Sprintf("%d%%", int(v*100+0.5))
		}), themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Caption, Color: c.MutedForeground}
		})))
	}
	return ui.Column(ui.Style{Gap: Space2, Width: ui.Pct(100)},
		ui.Row(ui.Style{Justify: ui.JustifySpaceBetween, Align: ui.AlignCenter},
			append([]ui.Node{ui.StaticText(p.Label, themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Font: t.Type.Small, Color: c.Foreground}
			}))}, trail...)...,
		),
		bar,
	)
}

// AlertProps configures an Alert banner (an inline message, not a dialog).
type AlertProps struct {
	Title       string
	Description string
	// Variant picks the tint: Destructive for an error, Primary for info.
	Variant Variant
	Icon    ui.Icon
	Action  func() ui.Node
}

// Alert is an inline, non-modal message. For a modal confirmation use
// Dialog; for a brief message use Toast.
func Alert(p AlertProps) ui.Node {
	icon := p.Icon
	if icon == "" {
		icon = ui.IconInfo
	}
	children := []ui.Node{
		ui.Symbol(icon, themed(func(t Theme, c Palette) ui.Style {
			_, fg, _ := variantColors(c, p.Variant)
			return ui.Style{Font: t.Type.Body, Color: fg}
		})),
		ui.Column(ui.Style{Gap: Space1, Grow: 1},
			ui.StaticText(p.Title, themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Font: t.Type.Subheading, Color: c.Foreground}
			})),
			ui.StaticText(p.Description, themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Font: t.Type.Small, Color: c.MutedForeground}
			})),
		),
	}
	if p.Action != nil {
		children = append(children, p.Action())
	}
	return ui.View(ui.ViewProps{
		Label: ui.Val(p.Title + ". " + p.Description),
		Style: themed(func(t Theme, c Palette) ui.Style {
			s := ui.Style{
				Direction: ui.RowDir, Align: ui.AlignStart, Gap: Space3,
				Padding: ui.All(Space4), CornerRadius: t.Radii.LG,
				Background: c.Card, BorderWidth: 1, BorderColor: c.Border,
			}
			switch p.Variant {
			case Destructive:
				s.Background = c.Destructive.Opacity(0.1)
				s.BorderColor = c.Destructive.Opacity(0.35)
			case Primary:
				s.Background = c.Primary.Opacity(0.08)
				s.BorderColor = c.Primary.Opacity(0.3)
			}
			return s
		}),
	}, children...)
}

// TooltipProps configures a Tooltip.
type TooltipProps struct {
	Text string
}

// Tooltip wraps a child with the platform's native tooltip (title on the
// web, UIToolTipInteraction on iPad, a tooltip on desktop). It must never
// be the only place important information is said: touch-only users never
// see it.
func Tooltip(p TooltipProps, child ui.Node) ui.Node {
	return ui.View(ui.ViewProps{Tooltip: ui.Val(p.Text)}, child)
}
