package main

import (
	"testing"

	ds "asuka-ds"

	"github.com/egoist/asuka/ops"
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
	"github.com/egoist/asuka/uitest"
)

func TestAppearanceTabsSwitchColorScheme(t *testing.T) {
	ds.SetAppearance(ui.ColorSchemeSystem)
	t.Cleanup(func() { ds.SetAppearance(ui.ColorSchemeSystem) })
	a := uitest.Mount(t, func() ui.Node {
		appearance := signal.New(0)
		signal.Effect(func() { ds.SetAppearance(ui.ColorScheme(appearance.Get())) })
		return ds.Tabs(ds.TabsProps{Items: []string{"System", "Light", "Dark"}, Selected: appearance})
	})
	last := func() byte {
		cmds := a.Host.Commands(ops.CmdColorScheme)
		if len(cmds) == 0 {
			return 255
		}
		b := cmds[len(cmds)-1].Value.Bytes
		if len(b) == 0 {
			return 255
		}
		return b[0]
	}
	a.Tap("Light")
	if got := last(); got != byte(ops.ColorSchemeLight) {
		t.Fatalf("after Light: CmdColorScheme = %d", got)
	}
	a.Tap("Dark")
	if got := last(); got != byte(ops.ColorSchemeDark) {
		t.Fatalf("after Dark: CmdColorScheme = %d", got)
	}
	a.Tap("System")
	if got := last(); got != byte(ops.ColorSchemeSystem) {
		t.Fatalf("after System: CmdColorScheme = %d", got)
	}
}
