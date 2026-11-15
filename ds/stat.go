package ds

import "github.com/egoist/asuka/ui"

// Trend is the direction of a Stat's change.
type Trend uint8

const (
	// TrendFlat has no direction and stays muted.
	TrendFlat Trend = iota
	// TrendUp rises: a dark arrow in the foreground color.
	TrendUp
	// TrendDown falls: a red arrow in the destructive color.
	TrendDown
)

// StatProps configures a Stat.
type StatProps struct {
	Label string
	Value string
	// Delta is the change, such as "+12.4%"; Trend chooses its arrow and
	// color.
	Delta string
	Trend Trend
	// Description is a muted line under the value.
	Description string
	Icon        ui.Icon
}

// Stat shows one metric: a label, a large value, an optional trend and a
// description, on a card-like surface. It is a display component, so the
// value is a string the caller formats.
func Stat(p StatProps) ui.Node {
	head := []ui.Node{
		ui.StaticText(p.Label, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Small, Color: c.MutedForeground, Lines: 1, Grow: 1}
		})),
	}
	if p.Icon != "" {
		head = append(head, ui.Symbol(p.Icon, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Small, Color: c.MutedForeground}
		})))
	}

	body := []ui.Node{
		ui.Row(ui.Style{Align: ui.AlignCenter, Gap: Space2}, head...),
		ui.StaticText(p.Value, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Title, Color: c.Foreground, Lines: 1}
		})),
	}
	if p.Delta != "" {
		delta := []ui.Node{ui.StaticText(p.Delta, themed(func(t Theme, c Palette) ui.Style {
			color := c.MutedForeground
			if p.Trend == TrendUp {
				color = c.Foreground
			} else if p.Trend == TrendDown {
				color = c.Destructive
			}
			return ui.Style{Font: t.Type.Small, Color: color, Lines: 1}
		}))}
		switch p.Trend {
		case TrendUp:
			delta = append([]ui.Node{trendIcon(ui.IconUp, p.Trend)}, delta...)
		case TrendDown:
			delta = append([]ui.Node{trendIcon(ui.IconDown, p.Trend)}, delta...)
		}
		body = append(body, ui.Row(ui.Style{Align: ui.AlignCenter, Gap: Space1}, delta...))
	}
	if p.Description != "" {
		body = append(body, ui.StaticText(p.Description, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Small, Color: c.MutedForeground, Lines: 2}
		})))
	}

	return ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{
			Padding:      ui.All(Space4),
			Gap:          Space2,
			Background:   c.Card,
			CornerRadius: t.Radii.LG,
			Shadow:       c.Elevation2.Shadow,
			BorderWidth:  1,
			BorderColor:  c.Border,
		}
	})}, ui.Column(ui.Style{Gap: Space1}, body...))
}

func trendIcon(icon ui.Icon, trend Trend) ui.Node {
	return ui.Symbol(icon, themed(func(t Theme, c Palette) ui.Style {
		color := c.MutedForeground
		if trend == TrendUp {
			color = c.Foreground
		} else if trend == TrendDown {
			color = c.Destructive
		}
		return ui.Style{Font: t.Type.Caption, Color: color}
	}))
}
