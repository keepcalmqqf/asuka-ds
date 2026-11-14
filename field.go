package ds

import (
	"strings"

	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// FieldProps configures a Field: the label/description/error structure
// around any control. Asuka has no label-for-control association, so the
// label is a header role and the container carries an accessibility label
// that names the control and its current state.
type FieldProps struct {
	Label       string
	Description string
	Error       string
	// Required adds "Required" to the label's text.
	Required bool
	// HideLabel skips the visible label: the control draws its own label
	// beside it (a switch, a slider with a live value). The group still
	// carries the accessible name.
	HideLabel bool
	// Disabled dims the label and description.
	Disabled signal.Accessor[bool]
	// DescriptionWithError keeps the description visible alongside an
	// error. It is off by default: an error replaces the description.
	DescriptionWithError bool
	// Control is the field's control; it is built after the label so a
	// screen reader reaches the label first.
	Control func() ui.Node
}

// Field is the form primitive used by Input, Select, Switch and Checkbox.
// It keeps error/description mutually exclusive and gives the whole group
// one accessible name that includes the error when there is one.
func Field(p FieldProps) ui.Node {
	labelText := p.Label
	if p.Required {
		labelText += " (required)"
	}
	groupLabel := labelText
	var details []string
	if p.Error != "" {
		details = append(details, strings.TrimSuffix(p.Error, "."))
	}
	if p.Description != "" && (p.Error == "" || p.DescriptionWithError) {
		details = append(details, strings.TrimSuffix(p.Description, "."))
	}
	if len(details) > 0 {
		groupLabel += ": " + strings.Join(details, ". ") + "."
	}

	var parts []ui.Node
	if !p.HideLabel {
		parts = append(parts, ui.View(ui.ViewProps{
			Label: ui.Val(labelText),
			Role:  ui.RoleHeader,
			Style: ui.Style{Width: ui.Pct(100)},
		}, ui.StaticText(labelText, themed(func(t Theme, c Palette) ui.Style {
			col := c.Foreground
			if p.Disabled != nil && untrackBool(p.Disabled) {
				col = c.MutedForeground
			}
			return ui.Style{Font: t.Type.Subheading, Color: col, Lines: 1}
		}))))
	}

	if p.Control != nil {
		parts = append(parts, p.Control())
	}
	if p.Error != "" {
		parts = append(parts, ui.View(ui.ViewProps{
			Label: ui.Val("Error: " + p.Error),
		}, ui.StaticText(p.Error, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Caption, Color: c.Destructive}
		}))))
	}
	if p.Description != "" && (p.Error == "" || p.DescriptionWithError) {
		parts = append(parts, ui.View(ui.ViewProps{
			Label: ui.Val(p.Description),
		}, ui.StaticText(p.Description, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Caption, Color: c.MutedForeground}
		}))))
	}
	return ui.View(ui.ViewProps{Label: ui.Val(groupLabel)},
		ui.Column(ui.Style{Gap: Space2, Width: ui.Pct(100)}, parts...))
}
