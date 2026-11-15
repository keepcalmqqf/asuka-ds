package ds

import (
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// TabsProps configures Tabs.
type TabsProps struct {
	Items    []string
	Selected *signal.Signal[int]
}

// Tabs is a segmented switch: a muted track with the selected item on a card.
// It's drawn from Views, so it looks the same everywhere; Asuka's Segmented
// is the native one (and is announced as selected by screen readers).
func Tabs(p TabsProps) ui.Node {
	track := themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{
			Direction:    ui.RowDir,
			Gap:          Space1,
			Padding:      ui.All(Space1),
			Width:        ui.Pct(100),
			Background:   c.Muted,
			CornerRadius: t.Radii.MD,
		}
	})
	kids := make([]ui.Node, 0, len(p.Items))
	for i, item := range p.Items {
		index := i
		label := item
		itemStyle := themed(func(t Theme, c Palette) ui.Style {
			selected := p.Selected != nil && p.Selected.Get() == index
			s := ui.Style{
				Direction:    ui.RowDir,
				Align:        ui.AlignCenter,
				Justify:      ui.JustifyCenter,
				Grow:         1,
				Height:       ui.Pt(32),
				CornerRadius: t.Radii.SM,
			}
			if selected {
				s.Background = c.Card
				s.Shadow = c.Elevation1.Shadow
			} else {
				s.Hovered = &ui.Style{Background: c.Background}
			}
			return s
		})
		labelStyle := themed(func(t Theme, c Palette) ui.Style {
			selected := p.Selected != nil && p.Selected.Get() == index
			col := c.MutedForeground
			if selected {
				col = c.Foreground
			}
			return ui.Style{Font: t.Type.Subheading, Color: col, Lines: 1}
		})
		selected := signal.Func(func() bool {
			return p.Selected != nil && p.Selected.Get() == index
		})
		kids = append(kids, Pressable(PressableProps{
			Label:    label,
			Selected: selected,
			OnPress:  func() { p.Selected.Set(index) },
			Style:    itemStyle,
		}, ui.StaticText(label, labelStyle)))
	}
	return ui.View(ui.ViewProps{Style: track}, kids...)
}
