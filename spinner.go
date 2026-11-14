package ds

import "github.com/egoist/asuka/ui"

// SpinnerProps configures a Spinner.
type SpinnerProps struct {
	// Label is shown beside the ring and names it; empty means the ring
	// alone (still announced as progress).
	Label string
	// Size is the ring's side, in points (0: 20).
	Size float32
}

// Spinner is the platform's indeterminate circular progress: an
// activity indicator while work whose duration is unknown is under way.
func Spinner(p SpinnerProps) ui.Node {
	size := p.Size
	if size <= 0 {
		size = 20
	}
	ring := ui.ProgressRing(nil, themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{Width: ui.Pt(size), Height: ui.Pt(size), Color: c.Primary}
	}))
	if p.Label == "" {
		return ring
	}
	return ui.View(ui.ViewProps{Label: ui.Val(p.Label)},
		ui.Row(ui.Style{Gap: Space2, Align: ui.AlignCenter},
			ring,
			ui.StaticText(p.Label, themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Font: t.Type.Small, Color: c.MutedForeground}
			})),
		))
}
