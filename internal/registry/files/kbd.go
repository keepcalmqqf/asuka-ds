package ds

import (
	"github.com/egoist/asuka/ui"
)

// Kbd renders one or more keys as keycaps: Kbd("⌘", "K") is ⌘ + K. It is
// for documentation and shortcut hints, not for input.
func Kbd(keys ...string) ui.Node {
	kids := make([]ui.Node, 0, len(keys)*2)
	for i, key := range keys {
		if i > 0 {
			kids = append(kids, ui.StaticText("+", themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Font: t.Type.Caption, Color: c.MutedForeground}
			})))
		}
		kids = append(kids, ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{
				Padding:    ui.Symmetric(2, Space1+2),
				Background: c.Muted, BorderWidth: 1, BorderColor: c.Border,
				CornerRadius: t.Radii.SM,
			}
		})}, ui.StaticText(key, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Caption, Color: c.Foreground, Lines: 1}
		}))))
	}
	return ui.Row(ui.Style{Gap: Space1, Align: ui.AlignCenter}, kids...)
}
