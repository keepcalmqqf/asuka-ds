package ds

import (
	"strconv"

	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// NumberFieldProps configures a NumberField.
type NumberFieldProps struct {
	Label       string
	Description string
	Value       *signal.Signal[float64]
	Min, Max    float64
	Step        float64
	// Format turns the value into its text ("%g" by default).
	Format   func(float64) string
	Disabled bool
}

// NumberField is a value beside the platform's stepper (UIStepper, two
// Material icon buttons): a quantity, a count, a small bounded number. The
// value is shown because a stepper has no label of its own.
func NumberField(p NumberFieldProps) ui.Node {
	format := p.Format
	if format == nil {
		format = func(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }
	}
	control := func() ui.Node {
		return ui.Row(ui.Style{Justify: ui.JustifySpaceBetween, Align: ui.AlignCenter, Width: ui.Pct(100)},
			ui.Text(signal.Func(func() string {
				v := 0.0
				if p.Value != nil {
					v = p.Value.Get()
				}
				return format(v)
			}), themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Font: t.Type.Subheading, Color: c.Foreground}
			})),
			ui.Stepper(ui.StepperProps{
				Value:    p.Value,
				Min:      p.Min,
				Max:      p.Max,
				Step:     p.Step,
				Disabled: ui.Val(p.Disabled),
			}),
		)
	}
	return Field(FieldProps{
		Label:       p.Label,
		Description: p.Description,
		Disabled:    ui.Val(p.Disabled),
		Control:     control,
	})
}
