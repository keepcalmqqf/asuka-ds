package ds

import (
	"fmt"

	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// SwitchProps configures a Switch or a Checkbox.
type SwitchProps struct {
	Label       string
	Description string
	Value       *signal.Signal[bool]
	// Checkbox makes it a checkbox: a list item, an agreement, a choice
	// among several that aren't exclusive.
	Checkbox bool
	Disabled bool
}

// Switch is a native switch (or checkbox) with the design system's label,
// description and tint. The control itself is the platform's own, so its
// state is announced correctly by screen readers.
func Switch(p SwitchProps) ui.Node {
	control := func() ui.Node {
		return ui.View(ui.ViewProps{Label: ui.Val(p.Label), Role: ui.RoleButton},
			ui.Row(ui.Style{Gap: Space3, Align: ui.AlignCenter, Width: ui.Pct(100)},
				ui.Toggle(ui.ToggleProps{
					Value:    p.Value,
					Checkbox: p.Checkbox,
					Disabled: ui.Val(p.Disabled),
					Style: themed(func(t Theme, c Palette) ui.Style {
						return ui.Style{Color: c.Primary}
					}),
				}),
				ui.StaticText(p.Label, themed(func(t Theme, c Palette) ui.Style {
					col := c.Foreground
					if p.Disabled {
						col = c.MutedForeground
					}
					return ui.Style{Font: t.Type.Body, Color: col}
				})),
			))
	}
	return Field(FieldProps{
		Label: p.Label, Description: p.Description, HideLabel: true,
		Disabled: ui.Val(p.Disabled), Control: control,
	})
}

// RadioGroupProps configures a RadioGroup.
type RadioGroupProps struct {
	Label       string
	Description string
	Options     []string
	Selected    *signal.Signal[int]
	Disabled    bool
}

// RadioGroup is a set of exclusive choices drawn from Views. Each option is
// a Pressable whose Hint says Selected, because a View has no selected trait.
func RadioGroup(p RadioGroupProps) ui.Node {
	control := func() ui.Node {
		kids := make([]ui.Node, 0, len(p.Options))
		for i, option := range p.Options {
			index, text := i, option
			kids = append(kids, Pressable(PressableProps{
				Label:    text,
				Selected: signal.Func(func() bool { return p.Selected != nil && p.Selected.Get() == index }),
				Disabled: ui.Val(p.Disabled),
				OnPress:  func() { p.Selected.Set(index) },
				Style: themed(func(t Theme, c Palette) ui.Style {
					return ui.Style{Direction: ui.RowDir, Align: ui.AlignCenter, Gap: Space2,
						Padding: ui.Symmetric(Space2, 0), Width: ui.Pct(100)}
				}),
			},
				ui.View(ui.ViewProps{Decorative: true, Style: themed(func(t Theme, c Palette) ui.Style {
					selected := p.Selected != nil && p.Selected.Get() == index
					s := ui.Style{Width: ui.Pt(16), Height: ui.Pt(16), CornerRadius: 8,
						BorderWidth: 1, BorderColor: c.Border}
					if selected {
						s.BorderColor = c.Primary
						s.Background = c.Primary
					}
					if p.Disabled {
						s.Opacity = 0.5
					}
					return s
				})}),
				ui.StaticText(text, themed(func(t Theme, c Palette) ui.Style {
					return ui.Style{Font: t.Type.Body, Color: c.Foreground}
				})),
			))
		}
		return ui.Column(ui.Style{Gap: Space1}, kids...)
	}
	return Field(FieldProps{Label: p.Label, Description: p.Description,
		Disabled: ui.Val(p.Disabled), Control: control})
}

// SelectProps configures a Select.
type SelectProps struct {
	Label       string
	Description string
	Placeholder string
	Options     []string
	Selected    *signal.Signal[int]
	Disabled    bool
}

// Select is a themed button that opens the platform's own menu of options
// (UIMenu, PopupMenu, NSMenu, a popover). The native menu gives keyboard
// navigation and checked items; the trigger stays brand-consistent.
func Select(p SelectProps) ui.Node {
	current := func() string {
		if p.Selected == nil || p.Selected.Get() < 0 || p.Selected.Get() >= len(p.Options) {
			return p.Placeholder
		}
		return p.Options[p.Selected.Get()]
	}
	items := make([]DialogAction, 0, len(p.Options))
	for i, option := range p.Options {
		index := i
		items = append(items, DialogAction{
			Title:   option,
			Checked: signal.Func(func() bool { return p.Selected != nil && p.Selected.Get() == index }),
			OnPress: func() { p.Selected.Set(index) },
		})
	}
	control := func() ui.Node {
		return Menu(MenuProps{Items: items, Label: p.Label, Disabled: p.Disabled},
			ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
				s := ui.Style{
					Direction: ui.RowDir, Align: ui.AlignCenter, Gap: Space2,
					Padding: ui.Symmetric(0, Space3), Height: ui.Pt(36),
					Width:      ui.Pct(100),
					Background: c.Background, BorderWidth: 1, BorderColor: c.Input,
					CornerRadius: t.Radii.MD,
				}
				if p.Disabled {
					s.Opacity = 0.5
				}
				return s
			})},
				ui.Text(signal.Func(current), themed(func(t Theme, c Palette) ui.Style {
					col := c.Foreground
					if p.Disabled {
						col = c.MutedForeground
					}
					return ui.Style{Font: t.Type.Body, Color: col, Lines: 1, Grow: 1}
				})),
				ui.Symbol(ui.IconDown, themed(func(t Theme, c Palette) ui.Style {
					return ui.Style{Font: t.Type.Caption, Color: c.MutedForeground}
				})),
			),
		)
	}
	return Field(FieldProps{Label: p.Label, Description: p.Description,
		Disabled: ui.Val(p.Disabled), Control: control, HideLabel: true})
}

// SliderProps configures a Slider.
type SliderProps struct {
	Label       string
	Description string
	Value       *signal.Signal[float64]
	Min, Max    float64
	Disabled    bool
	// Format turns the value into the text beside the label (default: an
	// integer percent between Min and Max).
	Format func(float64) string
}

// Slider is the platform's own slider with a themed label and live value.
func Slider(p SliderProps) ui.Node {
	format := p.Format
	if format == nil {
		format = func(v float64) string {
			lo, hi := p.Min, p.Max
			if lo == 0 && hi == 0 {
				hi = 1
			}
			pct := 0
			if hi > lo {
				pct = int((v-lo)/(hi-lo)*100 + 0.5)
			}
			return fmt.Sprintf("%d%%", pct)
		}
	}
	control := func() ui.Node {
		return ui.Column(ui.Style{Gap: Space2, Width: ui.Pct(100)},
			ui.Row(ui.Style{Justify: ui.JustifySpaceBetween},
				ui.StaticText(p.Label, themed(func(t Theme, c Palette) ui.Style {
					return ui.Style{Font: t.Type.Subheading, Color: c.Foreground}
				})),
				ui.Text(signal.Func(func() string {
					v := 0.0
					if p.Value != nil {
						v = p.Value.Get()
					}
					return format(v)
				}), themed(func(t Theme, c Palette) ui.Style {
					return ui.Style{Font: t.Type.Small, Color: c.MutedForeground}
				})),
			),
			ui.Slider(ui.SliderProps{
				Value: p.Value, Min: p.Min, Max: p.Max,
				Disabled: ui.Val(p.Disabled),
				Style: themed(func(t Theme, c Palette) ui.Style {
					return ui.Style{Width: ui.Pct(100)}
				}),
			}),
		)
	}
	return Field(FieldProps{Label: p.Label, Description: p.Description,
		Disabled: ui.Val(p.Disabled), Control: control, HideLabel: true})
}
