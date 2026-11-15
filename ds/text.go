package ds

import (
	"github.com/egoist/asuka/ui"
)

// Heading is the page's largest text.
func Heading(text string) ui.Node {
	return ui.StaticText(text, themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{Font: t.Type.Display, Color: c.Foreground}
	}))
}

// Title is a section title.
func Title(text string) ui.Node {
	return ui.StaticText(text, themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{Font: t.Type.Title, Color: c.Foreground}
	}))
}

// Text is body text.
func Text(text string) ui.Node {
	return ui.StaticText(text, themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{Font: t.Type.Body, Color: c.Foreground}
	}))
}

// Muted is secondary text: a description, a timestamp, a hint.
func Muted(text string) ui.Node {
	return ui.StaticText(text, themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{Font: t.Type.Small, Color: c.MutedForeground}
	}))
}

// Screen is the page shell: the theme's background, a scroll area, a
// readable maximum width and the page's gutters. Every app writes its own;
// this is the one the components are designed against.
func Screen(children ...ui.Node) ui.Node {
	body := ui.View(ui.ViewProps{Style: ui.Style{
		Width:     ui.Pct(100),
		MaxWidth:  ui.Pt(880),
		AlignSelf: ui.AlignCenter,
		Padding:   ui.All(Space5),
		Gap:       Space6,
		Direction: ui.ColumnDir,
	}}, children...)
	return ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{Grow: 1, Background: c.Background}
	})}, ui.Scroll(ui.ScrollProps{Style: ui.Style{Grow: 1}}, ui.SafeArea(ui.Style{Grow: 1}, body)))
}
