package ds

import (
	"testing"

	"github.com/egoist/asuka/host/fake"
	"github.com/egoist/asuka/ops"
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
	"github.com/egoist/asuka/uitest"
)

func TestPressableTracksDynamicDisabledState(t *testing.T) {
	disabled := signal.New(false)
	presses := 0
	a := uitest.Mount(t, func() ui.Node {
		return Pressable(PressableProps{
			Label:    "Archive",
			Disabled: disabled,
			OnPress:  func() { presses++ },
		}, ui.StaticText("Archive"))
	})
	v := viewWithLabel(a, "Archive")
	disabled.Set(true)
	a.Settle()
	if got := v.String(ops.PropA11yHint); got != "Unavailable" {
		t.Errorf("hint after disabling = %q, want Unavailable", got)
	}
	a.Tap("Archive")
	if presses != 0 {
		t.Errorf("disabled Pressable fired %d times", presses)
	}
}

func TestTabsAnnounceSelectedState(t *testing.T) {
	selected := signal.New(0)
	a := uitest.Mount(t, func() ui.Node {
		return Tabs(TabsProps{Items: []string{"Overview", "Settings"}, Selected: selected})
	})
	a.Tap("Settings")
	if got := viewWithLabel(a, "Overview").String(ops.PropA11yHint); got != "" {
		t.Errorf("unselected tab hint = %q, want empty", got)
	}
	if got := viewWithLabel(a, "Settings").String(ops.PropA11yHint); got != "Selected" {
		t.Errorf("selected tab hint = %q, want Selected", got)
	}
}

func TestOverlayShowsChildrenOnASheet(t *testing.T) {
	open := signal.New(false)
	a := uitest.Mount(t, func() ui.Node {
		return ui.Column(ui.Style{},
			ui.StaticText("Open sheet"),
			Overlay(OverlayProps{Open: open,
				Content: func() ui.Node { return ui.StaticText("Sheet body") },
			}),
		)
	})
	open.Set(true)
	a.Settle()
	if !a.Has("Open sheet") {
		t.Fatalf("overlay children are missing:\n%s", a.Snapshot())
	}
	if !a.Has("Sheet body") {
		t.Fatalf("overlay content is missing:\n%s", a.Snapshot())
	}
}

func TestDisabledSelectDoesNotOpenAMenu(t *testing.T) {
	selected := signal.New(0)
	a := uitest.Mount(t, func() ui.Node {
		return Select(SelectProps{
			Label: "Framework", Options: []string{"Asuka", "SwiftUI"},
			Selected: selected, Disabled: true,
		})
	})
	for _, v := range a.Host.Views {
		if _, ok := v.Props[ops.PropMenu]; ok {
			t.Fatalf("disabled Select still has a menu:\n%s", a.Snapshot())
		}
	}
}

func TestInputShowsDescriptionAndErrorWithFieldSemantics(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node {
		return Input(InputProps{
			Label: "Email", Description: "We never share it.",
			Error: "That address is taken.", Value: signal.New(""),
		})
	})
	if !a.Has("We never share it.") {
		t.Errorf("description is suppressed by error:\n%s", a.Snapshot())
	}
	if maybeViewWithLabel(a, "Email: That address is taken. We never share it.") == nil {
		t.Errorf("field group has no accessible name:\n%s", a.Snapshot())
	}
	label := maybeViewWithLabel(a, "Email")
	if label == nil {
		t.Fatalf("field label has no accessible name:\n%s", a.Snapshot())
	}
	assertRole(t, label, ops.RoleHeader)
}

func TestCommandResetsSelectionWhenQueryNarrows(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node {
		return Command(CommandProps{Items: []CommandItem{
			{Title: "New project", OnSelect: func() {}}, {Title: "Open settings", OnSelect: func() {}}, {Title: "Sign out", OnSelect: func() {}},
		}})
	})
	a.Tap("Search…")
	a.Type("Search…", "")
	for i := 0; i < 2; i++ {
		a.Key("down")
	}
	a.Type("Search…", "sign")
	if !a.Has("Sign out") {
		t.Fatalf("filtered result is missing:\n%s", a.Snapshot())
	}
	selected := 0
	for _, v := range a.Host.Views {
		if got := v.String(ops.PropA11yHint); got == "Selected" {
			selected++
		}
	}
	if selected != 1 {
		t.Errorf("selected hints = %d, want 1:\n%s", selected, a.Snapshot())
	}
}

func maybeViewWithLabel(a *uitest.App, label string) *fake.View {
	var out []*fake.View
	for _, v := range a.Host.Views {
		if v.String(ops.PropA11yLabel) == label {
			out = append(out, v)
		}
	}
	if len(out) != 1 {
		return nil
	}
	return out[0]
}
