package ds

import (
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// DialogAction is a button of a Dialog or the item of a Menu.
type DialogAction struct {
	Title    string
	Icon     ui.Icon
	Role     ui.ActionRole
	Disabled signal.Accessor[bool]
	Checked  signal.Accessor[bool]
	OnPress  func()
	// Children makes the item a submenu, opened when it is chosen.
	Children []DialogAction
	// Inline shows Children in place as a section with dividers, under
	// the title, instead of in a submenu. Dialog ignores both fields.
	Inline bool
}

// DialogProps configures a Dialog.
type DialogProps struct {
	Title   string
	Message string
	Actions []DialogAction
	// ActionSheet shows the actions as a list instead of a compact alert.
	ActionSheet bool
}

// Dialog shows the platform's own alert (UIAlertController, Material
// dialog, NSAlert, ContentDialog, AdwAlertDialog, <dialog>) with the design
// system's action ordering. Use it for confirmation; use Overlay for custom
// content.
func Dialog(p DialogProps) {
	actions := make([]ui.Action, 0, len(p.Actions))
	for _, a := range p.Actions {
		actions = append(actions, ui.Action{
			Title:    ui.Val(a.Title),
			Icon:     a.Icon,
			Role:     a.Role,
			Disabled: a.Disabled,
			Checked:  a.Checked,
			OnPress:  a.OnPress,
		})
	}
	ui.Alert(ui.AlertProps{
		Title:       p.Title,
		Message:     p.Message,
		Actions:     actions,
		ActionSheet: p.ActionSheet,
	})
}

// ConfirmDialog is a yes/no dialog whose destructive action is last.
func ConfirmDialog(title, message, confirm string, onConfirm func()) {
	Dialog(DialogProps{
		Title:   title,
		Message: message,
		Actions: []DialogAction{
			{Title: "Cancel", Role: ui.ActionCancel},
			{Title: confirm, Role: ui.ActionDestructive, OnPress: onConfirm},
		},
	})
}

// OverlayProps configures an Overlay.
type OverlayProps struct {
	Open *signal.Signal[bool]
	// Content is built when the overlay opens and disposed when it closes.
	Content func() ui.Node
	// Medium starts a sheet at half height (a bottom sheet on a phone).
	Medium bool
	// FullScreen covers the whole screen; the user can't dismiss it by
	// tapping outside.
	FullScreen bool
	// Grabber shows a drag handle.
	Grabber bool
	// Popover points at the trigger instead of covering the screen.
	Popover bool
	// MaxWidth is a popover's maximum width, in points (0: 320).
	MaxWidth  float32
	OnDismiss func()
}

// Overlay is the design system's modal primitive: it wraps Asuka's Sheet
// and Popover, so the platform owns dismissal, Escape, back gestures and
// accessibility focus, while the content is styled by the theme.
func Overlay(p OverlayProps, children ...ui.Node) ui.Node {
	content := p.Content
	if content == nil && len(children) > 0 {
		content = func() ui.Node {
			return ui.Column(ui.Style{Gap: Space4}, children...)
		}
	}
	render := func() ui.Node {
		if content == nil {
			return ui.View(ui.ViewProps{})
		}
		return content()
	}
	wrapped := func() ui.Node {
		return ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{
				Gap:          Space4,
				Padding:      ui.All(Space5),
				Background:   c.Popover,
				CornerRadius: t.Radii.XL,
				Shadow:       c.Elevation3.Shadow,
			}
		})}, render())
	}
	if p.Popover {
		return ui.Popover(ui.PopoverProps{
			Open:      p.Open,
			Content:   wrapped,
			MaxWidth:  p.MaxWidth,
			OnDismiss: p.OnDismiss,
		}, children...)
	}
	detent := ui.DetentLarge
	if p.Medium {
		detent = ui.DetentMedium
	}
	return ui.Sheet(ui.SheetProps{
		Open:       p.Open,
		Content:    wrapped,
		Detents:    detent,
		Grabber:    p.Grabber,
		FullScreen: p.FullScreen,
		OnDismiss:  p.OnDismiss,
	})
}

// MenuProps configures a Menu.
type MenuProps struct {
	Items []DialogAction
	// LongPress makes it a context menu rather than a tap menu.
	LongPress bool
	// Disabled keeps the trigger visible but removes the menu.
	Disabled bool
	// Label is the trigger's accessible name; the first item's title is
	// not a good name for a menu trigger.
	Label string
}

// Menu is the platform's own menu (UIMenu, PopupMenu, NSMenu, popover),
// which gives keyboard navigation, checked items and disabled items for
// free. It is the design system's answer to shadcn's DropdownMenu.
func Menu(p MenuProps, children ...ui.Node) ui.Node {
	actions := make([]ui.Action, 0, len(p.Items))
	for _, a := range p.Items {
		actions = append(actions, menuAction(a))
	}
	trigger := ui.View(ui.ViewProps{
		Label: ui.Val(p.Label),
		Style: themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{
				Direction:    ui.RowDir,
				Align:        ui.AlignCenter,
				Gap:          Space2,
				Padding:      ui.Symmetric(Space2, Space3),
				CornerRadius: t.Radii.MD,
				Color:        c.Foreground,
			}
		}),
	}, children...)
	if p.LongPress {
		return ui.Menu(ui.MenuProps{Actions: actions, LongPress: true,
			Style: ui.Style{}}, trigger)
	}
	if p.Disabled {
		return trigger
	}
	return ui.Menu(ui.MenuProps{Actions: actions, Style: ui.Style{}}, trigger)
}

// menuAction converts a DialogAction, and its children, to the native
// menu's Action. Submenus and sections are the platform's own, so nested
// menus keep keyboard navigation and disabled/checked items for free.
func menuAction(a DialogAction) ui.Action {
	action := ui.Action{
		Title:    ui.Val(a.Title),
		Icon:     a.Icon,
		Role:     a.Role,
		Disabled: a.Disabled,
		Checked:  a.Checked,
		OnPress:  a.OnPress,
		Inline:   a.Inline,
	}
	for _, child := range a.Children {
		action.Children = append(action.Children, menuAction(child))
	}
	return action
}

// Toast shows the platform's snackbar/toast. It is a function, not a
// component: a toast is shown once, not mounted in the tree.
func Toast(message string, action *DialogAction) func() {
	var a *ui.Action
	if action != nil {
		a = &ui.Action{
			Title:    ui.Val(action.Title),
			Disabled: action.Disabled,
			OnPress:  action.OnPress,
		}
	}
	return ui.Snackbar(ui.SnackbarProps{Message: message, Action: a})
}
