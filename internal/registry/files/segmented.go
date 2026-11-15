package ds

import (
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// SegmentedControlProps configures a SegmentedControl.
type SegmentedControlProps struct {
	Label       string
	Description string
	Segments    []string
	// Selected is the chosen segment's index; the user's choice writes it.
	Selected *signal.Signal[int]
	Disabled bool
}

// SegmentedControl is the platform's own segmented control in the design
// system's Field: a small set of exclusive choices side by side (UISegmented
// Control, Material's segmented buttons). Its selection is native, so
// screen readers announce it without the Pressable Selected workaround.
func SegmentedControl(p SegmentedControlProps) ui.Node {
	control := func() ui.Node {
		return ui.Segmented(ui.SegmentedProps{
			Segments: ui.Val(p.Segments),
			Selected: p.Selected,
			Disabled: ui.Val(p.Disabled),
			Style: themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Width: ui.Pct(100)}
			}),
		})
	}
	return Field(FieldProps{
		Label:       p.Label,
		Description: p.Description,
		Disabled:    ui.Val(p.Disabled),
		Control:     control,
	})
}
