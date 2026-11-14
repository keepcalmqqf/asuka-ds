package ds

import (
	"testing"

	"github.com/egoist/asuka/host"
	"github.com/egoist/asuka/ops"
	"github.com/egoist/asuka/ui"
	"github.com/egoist/asuka/uitest"
)

// background mounts a view styled by Themed and returns its background.
func background(t *testing.T, opts ...uitest.Option) (*uitest.App, func() int64) {
	t.Helper()
	haveEnvDark, ignoringEnv = false, false
	lastEnvDark, systemDark = false, false
	a := uitest.Mount(t, func() ui.Node {
		return ui.View(ui.ViewProps{Style: Themed(func(_ Theme, c Palette) ui.Style {
			return ui.Style{Background: c.Background}
		})})
	}, opts...)
	return a, func() int64 {
		t.Helper()
		for _, view := range a.Host.Views {
			if view.Kind != ops.KindView {
				continue
			}
			if b, ok := view.Props[ops.PropBackground]; ok {
				return b.Int
			}
		}
		t.Fatal("no background prop")
		return 0
	}
}

// Themed follows the host's dark mode as well as SetTheme: both are read
// inside the StyleFunc, so a change to either restyles without a remount.
func TestThemedFollowsDarkEnvironment(t *testing.T) {
	SetTheme(Neutral())
	SetAppearance(ui.ColorSchemeSystem)
	t.Cleanup(func() {
		SetTheme(Native())
		SetAppearance(ui.ColorSchemeSystem)
	})
	a, background := background(t)
	if got, want := background(), int64(Neutral().Light.Background); got != want {
		t.Errorf("light background = %d, want %d", got, want)
	}
	a.Host.SetEnvironment(host.Environment{Dark: true, Scale: 3, FontScale: 1})
	a.Settle()
	if got, want := background(), int64(Neutral().Dark.Background); got != want {
		t.Errorf("dark background = %d, want %d", got, want)
	}
	SetTheme(Native())
	a.Settle()
	if got, want := background(), int64(semanticPalette().Background); got != want {
		t.Errorf("native background = %d, want %d", got, want)
	}
}

// An app-level Light/Dark override restyles an explicit palette even before
// the host echoes the new appearance in its environment (macOS reports the
// app's appearance only after the window's own redraw).
func TestSetAppearanceOverridesEnvironment(t *testing.T) {
	SetTheme(Neutral())
	SetAppearance(ui.ColorSchemeSystem)
	t.Cleanup(func() {
		SetTheme(Native())
		SetAppearance(ui.ColorSchemeSystem)
	})
	a, background := background(t, uitest.Dark())
	if got, want := background(), int64(Neutral().Dark.Background); got != want {
		t.Fatalf("system dark background = %d, want %d", got, want)
	}
	SetAppearance(ui.ColorSchemeLight)
	a.Settle()
	if got, want := background(), int64(Neutral().Light.Background); got != want {
		t.Errorf("forced light background = %d, want %d", got, want)
	}
	SetAppearance(ui.ColorSchemeDark)
	a.Settle()
	if got, want := background(), int64(Neutral().Dark.Background); got != want {
		t.Errorf("forced dark background = %d, want %d", got, want)
	}
	SetAppearance(ui.ColorSchemeSystem)
	a.Settle()
	if got, want := background(), int64(Neutral().Dark.Background); got != want {
		t.Errorf("back to system dark background = %d, want %d", got, want)
	}
}

// Returning to System after a forced Light on a dark system keeps the
// system's appearance: macOS reports the forced one late, and that report
// must not be taken for the system's (the user sees a dark app that never
// goes back).
func TestSetAppearanceSystemIgnoresLateForcedReport(t *testing.T) {
	SetTheme(Neutral())
	SetAppearance(ui.ColorSchemeSystem)
	t.Cleanup(func() {
		SetTheme(Native())
		SetAppearance(ui.ColorSchemeSystem)
	})
	a, background := background(t, uitest.Dark())
	if got, want := background(), int64(Neutral().Dark.Background); got != want {
		t.Fatalf("system dark background = %d, want %d", got, want)
	}
	SetAppearance(ui.ColorSchemeLight)
	a.Settle()
	if got, want := background(), int64(Neutral().Light.Background); got != want {
		t.Fatalf("forced light background = %d, want %d", got, want)
	}
	// No report of the forced appearance has arrived yet.
	SetAppearance(ui.ColorSchemeSystem)
	// A repeated System setting (the selected tab clicked again) must not
	// stop ignoring the report that is still on its way.
	SetAppearance(ui.ColorSchemeSystem)
	a.Host.SetEnvironment(host.Environment{Dark: false, Scale: 3, FontScale: 1})
	if got, want := background(), int64(Neutral().Dark.Background); got != want {
		t.Errorf("system after a late forced report = %d, want %d", got, want)
	}
}

// macOS can report the forced appearance while the app is still forced:
// that report has already arrived when the user goes back to System, and
// System must still use the appearance seen before it.
func TestSetAppearanceSystemAfterLateForcedReportArrived(t *testing.T) {
	SetTheme(Neutral())
	SetAppearance(ui.ColorSchemeSystem)
	t.Cleanup(func() {
		SetTheme(Native())
		SetAppearance(ui.ColorSchemeSystem)
	})
	a, background := background(t, uitest.Dark())
	SetAppearance(ui.ColorSchemeLight)
	a.Settle()
	a.Host.SetEnvironment(host.Environment{Dark: false, Scale: 3, FontScale: 1})
	if got, want := background(), int64(Neutral().Light.Background); got != want {
		t.Fatalf("forced light after its report = %d, want %d", got, want)
	}
	SetAppearance(ui.ColorSchemeSystem)
	a.Settle()
	if got, want := background(), int64(Neutral().Dark.Background); got != want {
		t.Errorf("system dark background = %d, want %d", got, want)
	}
}

// The environment stays on the forced appearance after returning to System,
// so a second Light must not take that stale report for the system's.
func TestSetAppearanceRepeatedForcedCyclesKeepSystemDark(t *testing.T) {
	SetTheme(Neutral())
	SetAppearance(ui.ColorSchemeSystem)
	t.Cleanup(func() {
		SetTheme(Native())
		SetAppearance(ui.ColorSchemeSystem)
	})
	a, background := background(t, uitest.Dark())
	for range 2 {
		SetAppearance(ui.ColorSchemeLight)
		a.Settle()
		a.Host.SetEnvironment(host.Environment{Dark: false, Scale: 3, FontScale: 1})
		if got, want := background(), int64(Neutral().Light.Background); got != want {
			t.Fatalf("forced light background = %d, want %d", got, want)
		}
		SetAppearance(ui.ColorSchemeSystem)
		a.Settle()
		if got, want := background(), int64(Neutral().Dark.Background); got != want {
			t.Errorf("system dark background = %d, want %d", got, want)
		}
	}
}

// A forced appearance on the semantic theme draws the design system's own
// light or dark colors: no host round-trip, so it switches on the click.
func TestNativeForcedAppearanceUsesExplicitPalette(t *testing.T) {
	SetTheme(Native())
	SetAppearance(ui.ColorSchemeSystem)
	t.Cleanup(func() {
		SetTheme(Native())
		SetAppearance(ui.ColorSchemeSystem)
	})
	a, background := background(t, uitest.Dark())
	if got, want := background(), int64(semanticPalette().Background); got != want {
		t.Fatalf("system dark background = %d, want %d", got, want)
	}
	SetAppearance(ui.ColorSchemeLight)
	a.Settle()
	if got, want := background(), int64(nativePalette(false).Background); got != want {
		t.Errorf("forced light background = %d, want %d", got, want)
	}
	SetAppearance(ui.ColorSchemeDark)
	a.Settle()
	if got, want := background(), int64(nativePalette(true).Background); got != want {
		t.Errorf("forced dark background = %d, want %d", got, want)
	}
}
