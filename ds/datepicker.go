package ds

import (
	"time"

	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// DatePickerProps configures a DatePicker.
type DatePickerProps struct {
	Label       string
	Description string
	Error       string
	// Value is the chosen instant; the platform's picker writes it.
	Value *signal.Signal[time.Time]
	// Mode picks a date, a time or both (the zero value is a date).
	Mode ui.DateMode
	// Min and Max bound the choice (zero: no bound).
	Min, Max time.Time
	// Inline shows the calendar in place instead of behind a compact
	// button, on platforms that support it.
	Inline   bool
	Disabled bool
}

// DatePicker is the platform's own date, time or date-and-time picker in
// the design system's Field: UIKit's compact picker and its calendar,
// Material's pickers, a Mac's clock and calendar. The platform owns the
// editing, validation and accessibility; the Field owns the label, the
// description and the error.
func DatePicker(p DatePickerProps) ui.Node {
	control := func() ui.Node {
		return ui.DatePicker(ui.DatePickerProps{
			Value:    p.Value,
			Mode:     p.Mode,
			Min:      p.Min,
			Max:      p.Max,
			Inline:   p.Inline,
			Disabled: ui.Val(p.Disabled),
			Style: themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Font: t.Type.Body, Color: c.Foreground}
			}),
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
