package ds

import (
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// CardProps configures a Card.
type CardProps struct {
	// Padding overrides the card's own, for a card that holds a list or an
	// image that reaches its edges.
	Padding *ui.Edges
	// Flat drops the border and the shadow: a card on a surface that is
	// already a card.
	Flat bool
}

// Card is a surface: the theme's card color, a 1pt border and a small
// shadow. Compose it with CardHeader, CardContent and CardFooter.
func Card(p CardProps, children ...ui.Node) ui.Node {
	d := Space5
	pad := ui.All(d)
	if p.Padding != nil {
		pad = *p.Padding
	}
	return ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
		s := ui.Style{
			Padding:      pad,
			Gap:          Space4,
			Background:   c.Card,
			CornerRadius: t.Radii.LG,
			Shadow:       c.Elevation2.Shadow,
		}
		if !p.Flat {
			s.BorderWidth = 1
			s.BorderColor = c.Border
		}
		return s
	})}, children...)
}

// CardHeader is a card's title and description, with an optional action on
// the trailing edge.
func CardHeader(title, description string, action ...ui.Node) ui.Node {
	titleStyle := themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{Font: t.Type.Heading, Color: c.CardForeground, Lines: 1}
	})
	var head ui.Node
	if description == "" {
		head = ui.StaticText(title, titleStyle)
	} else {
		head = ui.Column(ui.Style{Gap: Space1},
			ui.StaticText(title, titleStyle),
			ui.StaticText(description, themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Font: t.Type.Small, Color: c.MutedForeground}
			})),
		)
	}
	if len(action) == 0 {
		return head
	}
	return ui.Row(ui.Style{Align: ui.AlignStart, Justify: ui.JustifySpaceBetween, Gap: Space4},
		ui.View(ui.ViewProps{Style: ui.Style{Grow: 1}}, head),
		action[0],
	)
}

// CardTitle is a card's title alone.
func CardTitle(text signal.Accessor[string]) ui.Node {
	return ui.Text(text, themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{Font: t.Type.Heading, Color: c.CardForeground}
	}))
}

// CardDescription is muted text under a title.
func CardDescription(text string) ui.Node {
	return ui.StaticText(text, themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{Font: t.Type.Small, Color: c.MutedForeground}
	}))
}

// CardContent stacks a card's body.
func CardContent(children ...ui.Node) ui.Node {
	return ui.Column(ui.Style{Gap: Space3}, children...)
}

// CardFooter is a card's actions, on its trailing edge.
func CardFooter(children ...ui.Node) ui.Node {
	return ui.Row(ui.Style{Gap: Space2, Align: ui.AlignCenter, Justify: ui.JustifyEnd, Wrap: true}, children...)
}
