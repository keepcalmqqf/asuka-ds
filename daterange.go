package ds

import (
	"time"

	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// DateRangePickerProps configures a DateRangePicker.
type DateRangePickerProps struct {
	Label       string
	Description string
	Error       string
	// Start and End are the chosen bounds; each picker writes its own.
	Start, End *signal.Signal[time.Time]
	// Min and Max bound the whole range (zero: no bound).
	Min, Max time.Time
	Disabled bool
}

// DateRangePicker is two platform date pickers with the range enforced
// between them: End can't precede Start and Start can't follow End. The
// bounds rebuild as either end changes, so the native controls carry the
// constraint themselves.
func DateRangePicker(p DateRangePickerProps) ui.Node {
	control := func() ui.Node {
		return ui.Dynamic(func() ui.Node {
			var start, end time.Time
			if p.Start != nil {
				start = p.Start.Get()
			}
			if p.End != nil {
				end = p.End.Get()
			}
			startMax, endMin := p.Max, p.Min
			if !end.IsZero() {
				startMax = end
			}
			if !start.IsZero() {
				endMin = start
			}
			picker := func(label string, value *signal.Signal[time.Time], min, max time.Time) ui.Node {
				return ui.View(ui.ViewProps{Label: ui.Val(label)},
					ui.DatePicker(ui.DatePickerProps{
						Value:    value,
						Min:      min,
						Max:      max,
						Disabled: ui.Val(p.Disabled),
						Style: themed(func(t Theme, c Palette) ui.Style {
							return ui.Style{Font: t.Type.Body, Color: c.Foreground}
						}),
					}))
			}
			return ui.Row(ui.Style{Gap: Space3, Align: ui.AlignCenter, Wrap: true},
				picker("Start date", p.Start, p.Min, startMax),
				ui.StaticText("to", themed(func(t Theme, c Palette) ui.Style {
					return ui.Style{Font: t.Type.Small, Color: c.MutedForeground}
				})),
				picker("End date", p.End, endMin, p.Max),
			)
		})
	}
	return Field(FieldProps{
		Label:                p.Label,
		Description:          p.Description,
		DescriptionWithError: true,
		Error:                p.Error,
		Disabled:             ui.Val(p.Disabled),
		Control:              control,
	})
}
