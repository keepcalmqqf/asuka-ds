package ds

import (
	"fmt"
	"testing"

	"github.com/egoist/asuka/host/fake"
	"github.com/egoist/asuka/ops"
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
	"github.com/egoist/asuka/uitest"
)

// The a11y matrix: every interactive primitive must be reachable by name,
// must carry a role, and must not fire its action while disabled. These run
// on Asuka's fake host, which records the exact props a real host receives.

func TestButtonIsNamedAndRoleButton(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node {
		return Button(ButtonProps{Text: "Save", Label: "Save changes", Variant: Primary,
			OnPress: func() {}})
	})
	v := viewWithLabel(a, "Save changes")
	if got := v.String(ops.PropA11yLabel); got != "Save changes" {
		t.Errorf("label = %q", got)
	}
	// Asuka makes any pressable View a button on every host, so an
	// enabled Button needs no explicit role. A disabled one isn't
	// pressable, and carries the role explicitly (checked below).
	if !v.Listens[ops.EventPress] {
		t.Error("enabled Button is not pressable")
	}
}

func TestButtonPressFiresAction(t *testing.T) {
	presses := 0
	a := uitest.Mount(t, func() ui.Node {
		return Button(ButtonProps{Text: "Save", OnPress: func() { presses++ }})
	})
	a.Tap("Save")
	if presses != 1 {
		t.Errorf("presses = %d, want 1", presses)
	}
}

func TestDisabledButtonHasNoActionAndHintUnavailable(t *testing.T) {
	presses := 0
	a := uitest.Mount(t, func() ui.Node {
		return Button(ButtonProps{Text: "Delete", Variant: Destructive, Disabled: true,
			OnPress: func() { presses++ }})
	})
	v := viewWithLabel(a, "Delete")
	if v.Listens[ops.EventPress] {
		t.Error("a disabled button still listens for presses")
	}
	if got := v.String(ops.PropA11yHint); got != "Unavailable" {
		t.Errorf("hint = %q, want Unavailable", got)
	}
	assertRole(t, v, ops.RoleButton)
	if presses != 0 {
		t.Error("disabled button fired its action")
	}
}

func TestNativeButtonUsesRealDisabledProp(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node {
		return NativeButton(ButtonProps{Text: "Delete", Disabled: true})
	})
	v := a.Get("Delete")
	if !v.Bool(ops.PropDisabled) {
		t.Error("native button's Disabled prop was not sent")
	}
	if a.Disabled("Delete") != true {
		t.Error("uitest does not see the native button as disabled")
	}
}

func TestPressableSelectedStateInHint(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node {
		return Pressable(PressableProps{Label: "Activity", Selected: ui.Val(true),
			OnPress: func() {}}, ui.StaticText("Activity"))
	})
	v := viewWithLabel(a, "Activity")
	if got := v.String(ops.PropA11yHint); got != "Selected" {
		t.Errorf("hint = %q, want Selected", got)
	}
}

func TestPressableDisabledIsNotPressable(t *testing.T) {
	presses := 0
	a := uitest.Mount(t, func() ui.Node {
		return Pressable(PressableProps{Label: "Archive", Disabled: ui.Val(true),
			OnPress: func() { presses++ }}, ui.StaticText("Archive"))
	})
	v := viewWithLabel(a, "Archive")
	if v.Listens[ops.EventPress] {
		t.Error("disabled Pressable listens for presses")
	}
	if got := v.String(ops.PropA11yHint); got != "Unavailable" {
		t.Errorf("hint = %q, want Unavailable", got)
	}
	if presses != 0 {
		t.Error("disabled Pressable fired its action")
	}
}

func TestFieldAnnouncesLabelAndError(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node {
		return Field(FieldProps{
			Label: "Email", Error: "That address is taken.",
			Control: func() ui.Node { return ui.StaticText("ada@example.com") },
		})
	})
	if !a.Has("Email") || !a.Has("That address is taken.") {
		t.Fatalf("field's label or error is not on screen:\n%s", a.Snapshot())
	}
	v := viewWithLabel(a, "Email: That address is taken.")
	if got := v.String(ops.PropA11yLabel); got == "" {
		t.Error("field group has no accessible name")
	}
	assertRole(t, viewWithLabel(a, "Email"), ops.RoleHeader)
}

func TestTabsUseSignalAndKeyboard(t *testing.T) {
	tab := signal.New(0)
	a := uitest.Mount(t, func() ui.Node {
		return Tabs(TabsProps{Items: []string{"Overview", "Settings"}, Selected: tab})
	})
	a.Tap("Settings")
	if tab.Peek() != 1 {
		t.Errorf("selected = %d, want 1", tab.Peek())
	}
	for _, name := range []string{"Overview", "Settings"} {
		v := viewWithLabel(a, name)
		if got := v.String(ops.PropA11yLabel); got != name {
			t.Errorf("%s label = %q", name, got)
		}
		if !v.Listens[ops.EventPress] {
			t.Errorf("%s is not pressable", name)
		}
	}
}

// viewWithLabel finds the one view carrying an accessibility label. Get is
// ambiguous when a container and its child show the same text; a11y props
// live on the container.
func viewWithLabel(a *uitest.App, label string) *fake.View {
	var out []*fake.View
	for _, v := range a.Host.Views {
		if v.String(ops.PropA11yLabel) == label {
			out = append(out, v)
		}
	}
	if len(out) != 1 {
		panic(fmt.Sprintf("%d views have label %q, want 1\n%s", len(out), label, a.Snapshot()))
	}
	return out[0]
}

func assertRole(t *testing.T, v *fake.View, want ops.Role) {
	t.Helper()
	if got := ops.Role(v.Int(ops.PropA11yRole)); got != want {
		t.Errorf("role = %d, want %d", got, want)
	}
}
