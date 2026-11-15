package ds

import (
	"testing"

	"github.com/egoist/asuka/ops"
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
	"github.com/egoist/asuka/uitest"
)

func TestPopoverOpensFromItsTrigger(t *testing.T) {
	open := signal.New(false)
	a := uitest.Mount(t, func() ui.Node {
		return Popover(PopoverProps{
			Label: "Filters",
			Open:  open,
			Content: func() ui.Node {
				return ui.StaticText("Only this week", ui.Style{})
			},
		})
	})
	if a.Has("Only this week") {
		t.Fatalf("the popover is open before its trigger is pressed:\n%s", a.Snapshot())
	}
	assertRole(t, viewWithLabel(a, "Filters"), ops.RoleButton)
	if got := viewWithLabel(a, "Filters").String(ops.PropA11yHint); got != "Collapsed" {
		t.Errorf("closed hint = %q, want Collapsed", got)
	}
	a.Tap("Filters")
	if !open.Peek() {
		t.Error("the trigger did not open the popover")
	}
	if got := viewWithLabel(a, "Filters").String(ops.PropA11yHint); got != "Expanded" {
		t.Errorf("open hint = %q, want Expanded", got)
	}
	if !a.Has("Only this week") {
		t.Errorf("the content did not open:\n%s", a.Snapshot())
	}
}

func TestPopoverKeepsItsOwnOpenState(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node {
		return Popover(PopoverProps{
			Label: "Details",
			Content: func() ui.Node {
				return ui.StaticText("Deployed 4 minutes ago", ui.Style{})
			},
		})
	})
	a.Tap("Details")
	if !a.Has("Deployed 4 minutes ago") {
		t.Errorf("the content did not open:\n%s", a.Snapshot())
	}
}

func TestListRowsPressAndKeepTheirAction(t *testing.T) {
	pressed := ""
	notify := signal.New(true)
	a := uitest.Mount(t, func() ui.Node {
		return List(ListProps{Divided: true},
			ListItem(ListItemProps{
				Title: "Profile", Description: "Name and photo", Icon: ui.IconPerson,
				Chevron: true, OnPress: func() { pressed = "profile" },
			}),
			ListItem(ListItemProps{
				Title: "Notifications",
				Action: ui.Toggle(ui.ToggleProps{
					Value: notify,
					Style: Themed(func(t Theme, c Palette) ui.Style {
						return ui.Style{Color: c.Primary}
					}),
				}),
			}),
		)
	})
	if !a.Has("Name and photo") {
		t.Errorf("the description is not on screen:\n%s", a.Snapshot())
	}
	a.Tap("Profile")
	if pressed != "profile" {
		t.Errorf("pressed = %q, want profile", pressed)
	}
	a.Toggle("Notifications")
	if notify.Peek() {
		t.Error("the trailing toggle is still on; the row swallowed its press")
	}
}

func TestListItemShowsDetailAndSelection(t *testing.T) {
	selected := signal.New(true)
	a := uitest.Mount(t, func() ui.Node {
		return ListItem(ListItemProps{
			Title: "Plan", Detail: "Pro", Chevron: true,
			Selected: selected, OnPress: func() {},
		})
	})
	if !a.Has("Pro") {
		t.Errorf("the detail is not on screen:\n%s", a.Snapshot())
	}
	if got := viewWithLabel(a, "Plan").String(ops.PropA11yHint); got != "Selected" {
		t.Errorf("hint = %q, want Selected", got)
	}
}

func TestStatShowsMetricAndTrend(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node {
		return Stat(StatProps{
			Label: "Monthly revenue", Value: "$48,210", Delta: "+12.4%",
			Trend: TrendUp, Description: "Compared with last month",
		})
	})
	for _, want := range []string{"Monthly revenue", "$48,210", "+12.4%", "Compared with last month"} {
		if !a.Has(want) {
			t.Errorf("%q is not on screen:\n%s", want, a.Snapshot())
		}
	}
}

func TestProgressRingShowsPercentage(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node {
		return ProgressRing(ProgressRingProps{Value: ui.Val(0.72), Label: "Rollout", ShowValue: true})
	})
	if len(a.Host.OfKind(ops.KindProgressRing)) != 1 {
		t.Fatalf("want one progress ring:\n%s", a.Snapshot())
	}
	if !a.Has("72%") {
		t.Errorf("the percentage is not on screen:\n%s", a.Snapshot())
	}
}

func TestProgressRingWithoutAValueIsIndeterminate(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node {
		return ProgressRing(ProgressRingProps{Size: 24, Label: "Syncing", ShowValue: true})
	})
	rings := a.Host.OfKind(ops.KindProgressRing)
	if len(rings) != 1 {
		t.Fatalf("want one progress ring:\n%s", a.Snapshot())
	}
	if !rings[0].Bool(ops.PropProgressRingIndeterminate) {
		t.Error("a ring with no value is not indeterminate")
	}
}

func TestColorPickerUsesTheNativePicker(t *testing.T) {
	tint := signal.New(ui.Color(0))
	a := uitest.Mount(t, func() ui.Node {
		return ColorPicker(ColorPickerProps{Label: "Accent color", Value: tint})
	})
	pickers := a.Host.OfKind(ops.KindColorPicker)
	if len(pickers) != 1 {
		t.Fatalf("want one color picker:\n%s", a.Snapshot())
	}
	if !a.Has("Accent color") {
		t.Errorf("the label is not on screen:\n%s", a.Snapshot())
	}
	if pickers[0].String(ops.PropColorPickerTitle) != "Accent color" {
		t.Errorf("title = %q, want Accent color", pickers[0].String(ops.PropColorPickerTitle))
	}
}

func TestColorPickerWritesItsValue(t *testing.T) {
	tint := signal.New(ui.Color(0))
	a := uitest.Mount(t, func() ui.Node {
		return ColorPicker(ColorPickerProps{Label: "Accent", Value: tint})
	})
	pickers := a.Host.OfKind(ops.KindColorPicker)
	if len(pickers) != 1 {
		t.Fatalf("want one color picker:\n%s", a.Snapshot())
	}
	a.Host.PickColor(pickers[0], ui.Color(0x3366ccff))
	if tint.Peek() != ui.Color(0x3366ccff) {
		t.Errorf("value = %v, want 0x3366ccff", tint.Peek())
	}
}
