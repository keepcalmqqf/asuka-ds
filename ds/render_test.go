package ds

import (
	"testing"

	"github.com/egoist/asuka/colors"
	"github.com/egoist/asuka/ops"
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
	"github.com/egoist/asuka/uitest"
)

// A container's Color tints it on iOS, and a semantic background made from
// the tint color resolves against that tint. A drawn button must color its
// label and icon, not its own container, or an accent background turns into
// the label color.
func TestButtonDoesNotTintItsOwnBackground(t *testing.T) {
	SetTheme(Native())
	t.Cleanup(func() { SetTheme(Native()) })
	a := uitest.Mount(t, func() ui.Node {
		return Button(ButtonProps{Text: "Save", Variant: Primary, OnPress: func() {}})
	})
	button := viewWithLabel(a, "Save")
	if got := ops.Color(button.Int(ops.PropBackground)); got != colors.Accent {
		t.Errorf("background = %v, want the semantic accent", got)
	}
	if _, ok := button.Props[ops.PropColor]; ok {
		t.Error("the button tints its container; on iOS that retints its accent background")
	}
	found := false
	for _, text := range a.Host.OfKind(ops.KindText) {
		if text.String(ops.PropTextText) != "Save" {
			continue
		}
		found = true
		if got := ops.Color(text.Int(ops.PropColor)); got != colors.OnAccent {
			t.Errorf("label color = %v, want the on-accent color", got)
		}
	}
	if !found {
		t.Fatal("the label is not a text view")
	}
}

func TestToggleDoesNotTintItsOwnBackground(t *testing.T) {
	SetTheme(Native())
	t.Cleanup(func() { SetTheme(Native()) })
	a := uitest.Mount(t, func() ui.Node {
		return Toggle(ToggleProps{Text: "Bold", Variant: Primary, Pressed: signal.New(true)})
	})
	toggle := viewWithLabel(a, "Bold")
	if got := ops.Color(toggle.Int(ops.PropBackground)); got != colors.Accent.Opacity(0.12) {
		t.Errorf("background = %v, want the selected accent tint", got)
	}
	if _, ok := toggle.Props[ops.PropColor]; ok {
		t.Error("the toggle tints its container; on iOS that retints its accent background")
	}
}
