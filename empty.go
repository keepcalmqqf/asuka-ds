package ds

import (
	"strings"

	"github.com/egoist/asuka/ui"
)

// EmptyStateProps configures an EmptyState.
type EmptyStateProps struct {
	Icon        ui.Icon
	Title       string
	Description string
	// Action is the way out: usually one Button.
	Action func() ui.Node
}

// EmptyState is a placeholder for a view with nothing in it yet: an icon,
// a title, a sentence on what to do, and an action. It is composition, not
// a new primitive; every part uses the theme.
func EmptyState(p EmptyStateProps) ui.Node {
	label := p.Title
	if p.Description != "" {
		label += ": " + strings.TrimSuffix(p.Description, ".")
	}
	var kids []ui.Node
	if p.Icon != "" {
		kids = append(kids, ui.View(ui.ViewProps{Decorative: true, Style: themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{
				Align: ui.AlignCenter, Justify: ui.JustifyCenter,
				Width: ui.Pt(48), Height: ui.Pt(48), CornerRadius: 24,
				Background: c.Muted,
			}
		})}, ui.Symbol(p.Icon, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Title, Color: c.MutedForeground}
		}))))
	}
	if p.Title != "" {
		kids = append(kids, ui.StaticText(p.Title, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Subheading, Color: c.Foreground,
				TextAlign: ui.TextCenter, Lines: 2}
		})))
	}
	if p.Description != "" {
		kids = append(kids, ui.StaticText(p.Description, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Small, Color: c.MutedForeground,
				TextAlign: ui.TextCenter, MaxWidth: ui.Pt(360)}
		})))
	}
	if p.Action != nil {
		kids = append(kids, p.Action())
	}
	return ui.View(ui.ViewProps{Label: ui.Val(label), Style: themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{
			Align: ui.AlignCenter, Gap: Space3,
			Padding:    ui.All(Space6),
			Width:      ui.Pct(100),
			Background: c.Card, BorderWidth: 1, BorderColor: c.Border,
			CornerRadius: t.Radii.LG,
		}
	})}, kids...)
}
