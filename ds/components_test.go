package ds

import (
	"testing"

	"github.com/egoist/asuka/ops"
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
	"github.com/egoist/asuka/uitest"
)

func TestSwitchUsesNativeToggle(t *testing.T) {
	value := signal.New(false)
	a := uitest.Mount(t, func() ui.Node {
		return Switch(SwitchProps{Label: "Notifications", Value: value})
	})
	if len(a.Host.OfKind(ops.KindToggle)) != 1 {
		t.Fatalf("want one native toggle, got %d", len(a.Host.OfKind(ops.KindToggle)))
	}
	a.Toggle("Notifications")
	if !value.Peek() {
		t.Error("toggling did not update the signal")
	}
}

func TestRadioGroupSelectsOne(t *testing.T) {
	selected := signal.New(0)
	a := uitest.Mount(t, func() ui.Node {
		return RadioGroup(RadioGroupProps{
			Label: "Plan", Options: []string{"Free", "Pro"}, Selected: selected,
		})
	})
	a.Tap("Pro")
	if selected.Peek() != 1 {
		t.Errorf("selected = %d, want 1", selected.Peek())
	}
	if got := viewWithLabel(a, "Pro").String(ops.PropA11yHint); got != "Selected" {
		t.Errorf("Pro hint = %q, want Selected", got)
	}
}

func TestAccordionBuildsContentOnlyWhileOpen(t *testing.T) {
	open := signal.New(-1)
	built := 0
	a := uitest.Mount(t, func() ui.Node {
		return Accordion(AccordionProps{
			Open: open,
			Items: []AccordionItem{{
				Title: "Details", Content: func() ui.Node {
					built++
					return ui.StaticText("The fine print")
				},
			}},
		})
	})
	if built != 0 {
		t.Fatalf("content was built while closed")
	}
	a.Tap("Details")
	if built == 0 || !a.Has("The fine print") {
		t.Fatalf("opening did not build content:\n%s", a.Snapshot())
	}
	if got := viewWithLabel(a, "Details").String(ops.PropA11yHint); got != "Expanded" {
		t.Errorf("hint = %q, want Expanded", got)
	}
	a.Tap("Details")
	if a.Has("The fine print") {
		t.Error("closing did not remove content")
	}
}

func TestSelectOpensPlatformMenu(t *testing.T) {
	selected := signal.New(0)
	a := uitest.Mount(t, func() ui.Node {
		return Select(SelectProps{
			Label: "Framework", Options: []string{"Asuka", "SwiftUI"}, Selected: selected,
		})
	})
	menus := a.Host.OfKind(ops.KindView)
	withMenu := 0
	for _, v := range menus {
		if _, ok := v.Props[ops.PropMenu]; ok {
			withMenu++
		}
	}
	if withMenu == 0 {
		t.Error("Select's trigger has no native menu prop")
	}
}

func TestTableShowsHeadersAndRows(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node {
		return Table(TableProps{
			Columns: []TableColumn{{Title: "Name"}, {Title: "Role"}},
			Rows:    [][]string{{"Ada", "Engineer"}},
		})
	})
	for _, s := range []string{"Name", "Role", "Ada", "Engineer"} {
		if !a.Has(s) {
			t.Errorf("%q is not on screen:\n%s", s, a.Snapshot())
		}
	}
}

func TestAvatarFallsBackToInitials(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node { return Avatar(AvatarProps{Name: "Ada Lovelace"}) })
	if !a.Has("AL") {
		t.Errorf("initials are not on screen:\n%s", a.Snapshot())
	}
	if len(a.Host.OfKind(ops.KindImage)) != 0 {
		t.Error("avatar without a source should not create an image view")
	}
}

func TestProgressIsNative(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node {
		return Progress(ProgressProps{Value: ui.Val(0.42), Label: "Upload", ShowValue: true})
	})
	if len(a.Host.OfKind(ops.KindProgress)) != 1 {
		t.Error("want one native progress view")
	}
	if !a.Has("42%") {
		t.Errorf("value is not shown:\n%s", a.Snapshot())
	}
}

func TestDialogUsesPlatformAlert(t *testing.T) {
	pressed := 0
	a := uitest.Mount(t, func() ui.Node {
		return Button(ButtonProps{Text: "Delete account", Variant: Destructive,
			OnPress: func() {
				Dialog(DialogProps{Title: "Are you sure?", Actions: []DialogAction{
					{Title: "Cancel", Role: ui.ActionCancel},
					{Title: "Delete", Role: ui.ActionDestructive, OnPress: func() { pressed++ }},
				}})
			}})
	})
	a.Tap("Delete account")
	if len(a.Host.OfKind(ops.KindAlert)) != 1 {
		t.Fatalf("dialog did not create an alert view:\n%s", a.Snapshot())
	}
	a.Tap("Delete")
	if pressed != 1 {
		t.Errorf("confirm pressed %d times, want 1", pressed)
	}
}

func TestToastUsesPlatformSnackbar(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node {
		return Button(ButtonProps{Text: "Save", OnPress: func() {
			Toast("Saved", &DialogAction{Title: "Undo"})
		}})
	})
	a.Tap("Save")
	if len(a.Host.OfKind(ops.KindSnackbar)) != 1 {
		t.Fatalf("toast did not create a snackbar view:\n%s", a.Snapshot())
	}
	if !a.Has("Saved") || !a.Has("Undo") {
		t.Errorf("toast text is not on screen:\n%s", a.Snapshot())
	}
}
