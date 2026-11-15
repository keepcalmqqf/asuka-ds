package main

import (
	"testing"

	ds "github.com/egoist/asuka-ds"
	"github.com/egoist/asuka/ui"
	"github.com/egoist/asuka/uitest"
)

func TestSiteLoadsAndNavigates(t *testing.T) {
	ds.SetAppearance(ui.ColorSchemeSystem)
	t.Cleanup(func() { ds.SetAppearance(ui.ColorSchemeSystem) })
	a := uitest.Mount(t, App)
	a.SetEnvironment(func(e *ui.Environment) {
		e.WidthClass = ui.SizeClassExpanded
		e.Width = 1200
		e.Height = 900
	})
	for _, want := range []string{
		"Asuka DS",
		"Components that feel native on every platform",
		"Get started",
	} {
		if !a.Has(want) {
			t.Fatalf("%q is not on screen:\n%s", want, a.Snapshot())
		}
	}
	a.Tap("Installation")
	if !a.Has("1. Install the CLI") {
		t.Fatalf("installation page did not open:\n%s", a.Snapshot())
	}
	a.Tap("Combobox")
	if !a.Has("A searchable select: a native search field and a keyboard-navigable list.") {
		t.Fatalf("component page did not open:\n%s", a.Snapshot())
	}
	a.Tap("List")
	if !a.Has("Settings-style rows with icons, descriptions, details and trailing actions.") {
		t.Fatalf("the List page did not open:\n%s", a.Snapshot())
	}
	a.Tap("Popover")
	if !a.Has("A bubble of content anchored to a trigger, dismissed by tapping outside.") {
		t.Fatalf("the Popover page did not open:\n%s", a.Snapshot())
	}
	a.Tap("Platforms")
	if !a.Has("What was built and seen on each host, and what is still a preview.") {
		t.Fatalf("the Platforms page did not open:\n%s", a.Snapshot())
	}
}
