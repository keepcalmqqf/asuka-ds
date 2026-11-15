package ds

import (
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// PressableProps configures a Pressable: the design system's interactive
// primitive. Asuka's View has no Disabled and no selected state, so
// Pressable encodes those with the tools it does have: a disabled target is
// not pressable, its Hint says "Unavailable", and its Role stays Button so
// screen readers still announce what it is.
type PressableProps struct {
	Label string
	Hint  string
	Role  ui.Role
	// Disabled stops the action and dims the view.
	Disabled signal.Accessor[bool]
	// Selected is for a tab, a list row or a choice. It is announced in
	// Hint because Asuka has no selected trait for a View.
	Selected signal.Accessor[bool]
	// Expanded is for a disclosure (an accordion trigger).
	Expanded signal.Accessor[bool]
	// Style is the view's style; it can be a StyleFunc that reads the
	// interaction states through Interaction.
	Style ui.Styler
	// Interaction receives the view's hover/pressed states for children.
	Interaction *ui.Interaction
	Keys        ui.KeyMap
	OnPress     func()
	OnLongPress func()
	Href        ui.Href
}

// Pressable builds a pressable View with the system's a11y conventions.
// It is the base of Button, Tabs and the interactive rows of M3.
func Pressable(p PressableProps, children ...ui.Node) ui.Node {
	if p.Role == ui.RoleNone {
		p.Role = ui.RoleButton
	}
	hint := signal.Func(func() string {
		return stateHint(p)
	})
	if p.Disabled == nil && p.Selected == nil && p.Expanded == nil {
		hint = ui.Val(p.Hint)
	}
	press := func() {
		if p.Disabled != nil && untrackBool(p.Disabled) {
			return
		}
		if p.OnPress != nil {
			p.OnPress()
		}
	}
	longPress := func() {
		if p.Disabled != nil && untrackBool(p.Disabled) {
			return
		}
		if p.OnLongPress != nil {
			p.OnLongPress()
		}
	}
	props := ui.ViewProps{
		Label:       ui.Val(p.Label),
		Hint:        hint,
		Role:        p.Role,
		Style:       pressableStyle(p),
		Interaction: p.Interaction,
		Keys:        p.Keys,
		Href:        p.Href,
	}
	if p.OnPress != nil && (p.Disabled == nil || !signal.Untrack(p.Disabled.Get)) {
		props.OnPress = press
	}
	if p.OnLongPress != nil && (p.Disabled == nil || !signal.Untrack(p.Disabled.Get)) {
		props.OnLongPress = longPress
	}
	return ui.View(props, children...)
}

func stateHint(p PressableProps) string {
	if p.Disabled != nil && p.Disabled.Get() {
		return "Unavailable"
	}
	if p.Selected != nil && p.Selected.Get() {
		return "Selected"
	}
	if p.Expanded != nil && p.Disabled == nil {
		if p.Expanded.Get() {
			return "Expanded"
		}
		return "Collapsed"
	}
	return p.Hint
}

func pressableStyle(p PressableProps) ui.Styler {
	if p.Disabled == nil {
		return p.Style
	}
	return ui.StyleFunc(func() ui.Style {
		s := ui.StyleValue(p.Style)
		if untrackBool(p.Disabled) {
			s.Opacity = 0.5
		}
		return s
	})
}

func untrackBool(a signal.Accessor[bool]) bool {
	if v, ok := signal.StaticValue(a); ok {
		return v
	}
	return signal.Untrack(a.Get)
}
