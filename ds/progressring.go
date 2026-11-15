package ds

import (
	"fmt"

	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// ProgressRingProps configures a ProgressRing.
type ProgressRingProps struct {
	// Value is 0..1. A nil or negative value is indeterminate.
	Value signal.Accessor[float64]
	// Size is the ring's diameter in points (0: 48).
	Size float32
	// Label is the ring's accessible name. It is not drawn.
	Label string
	// ShowValue prints the percentage beside the ring.
	ShowValue bool
}

// ProgressRing is circular progress in the theme's primary color. It is the
// compact counterpart of Progress: use it in a toolbar, a card corner or a
// row, where a bar would take too much room.
func ProgressRing(p ProgressRingProps) ui.Node {
	size := p.Size
	if size <= 0 {
		size = 48
	}
	ring := ui.ProgressRing(p.Value, themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{Width: ui.Pt(size), Height: ui.Pt(size), Color: c.Primary}
	}))
	if p.Label != "" {
		ring = ui.View(ui.ViewProps{Label: ui.Val(p.Label)}, ring)
	}
	if !p.ShowValue || p.Value == nil {
		return ring
	}
	value := ui.Text(signal.Func(func() string {
		v := p.Value.Get()
		if v < 0 {
			return "…"
		}
		return fmt.Sprintf("%d%%", int(v*100+0.5))
	}), themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{Font: t.Type.Small, Color: c.MutedForeground, Lines: 1}
	}))
	return ui.Row(ui.Style{Align: ui.AlignCenter, Gap: Space2}, ring, value)
}
