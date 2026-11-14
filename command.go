package ds

import (
	"strings"

	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// CommandItem is an entry of a Command palette.
type CommandItem struct {
	Title    string
	Subtitle string
	Icon     ui.Icon
	Keywords string
	OnSelect func()
}

// CommandProps configures a Command palette.
type CommandProps struct {
	Items       []CommandItem
	Placeholder string
	// Empty is shown when nothing matches; a sensible default is used.
	Empty string
}

// Command is a searchable list with keyboard navigation: Up/Down move the
// selection, Enter chooses it, Escape closes the overlay. It is built from a
// native TextField, Pressable rows and Asuka's KeyMap, so it works with a
// hardware keyboard on desktop, iPad and the web.
func Command(p CommandProps) ui.Node {
	query := signal.New("")
	selected := signal.New(0)
	open := signal.New(false)
	empty := p.Empty
	if empty == "" {
		empty = "No results"
	}
	placeholder := p.Placeholder
	if placeholder == "" {
		placeholder = "Search…"
	}

	matches := signal.Func(func() []CommandItem {
		q := strings.ToLower(strings.TrimSpace(query.Get()))
		var out []CommandItem
		for _, item := range p.Items {
			if q == "" || strings.Contains(strings.ToLower(item.Title+" "+item.Keywords), q) {
				out = append(out, item)
			}
		}
		return out
	})
	signal.Effect(func() {
		items := matches.Get()
		if selected.Peek() >= len(items) {
			selected.Set(0)
		}
	})

	rows := func() ui.Node {
		return ui.Column(ui.Style{Gap: Space1, Width: ui.Pct(100)},
			ui.Dynamic(func() ui.Node {
				items := matches.Get()
				if len(items) == 0 {
					return ui.View(ui.ViewProps{Style: ui.Style{Padding: ui.All(Space4)}},
						ui.StaticText(empty, themed(func(t Theme, c Palette) ui.Style {
							return ui.Style{Font: t.Type.Small, Color: c.MutedForeground}
						})))
				}
				nodes := make([]ui.Node, 0, len(items))
				for i, item := range items {
					index, entry := i, item
					active := signal.Func(func() bool { return selected.Get() == index })
					nodes = append(nodes, Pressable(PressableProps{
						Label:    entry.Title,
						Selected: active,
						OnPress: func() {
							open.Set(false)
							if entry.OnSelect != nil {
								entry.OnSelect()
							}
						},
						Style: themed(func(t Theme, c Palette) ui.Style {
							s := ui.Style{
								Direction: ui.RowDir, Align: ui.AlignCenter, Gap: Space3,
								Padding: ui.Symmetric(Space2, Space3), Width: ui.Pct(100),
								CornerRadius: 6,
							}
							if active.Get() {
								s.Background = c.Accent
							} else {
								s.Hovered = &ui.Style{Background: c.Muted}
							}
							return s
						}),
					},
						ui.Symbol(entry.Icon, themed(func(t Theme, c Palette) ui.Style {
							return ui.Style{Font: t.Type.Body, Color: c.MutedForeground}
						})),
						ui.Column(ui.Style{Gap: 0, Grow: 1},
							ui.StaticText(entry.Title, themed(func(t Theme, c Palette) ui.Style {
								return ui.Style{Font: t.Type.Body, Color: c.Foreground, Lines: 1}
							})),
							ui.StaticText(entry.Subtitle, themed(func(t Theme, c Palette) ui.Style {
								return ui.Style{Font: t.Type.Caption, Color: c.MutedForeground, Lines: 1}
							})),
						),
					))
				}
				return ui.Column(ui.Style{Gap: Space1}, nodes...)
			}),
		)
	}

	choose := func() {
		items := signal.Untrack(matches.Get)
		i := selected.Peek()
		if i >= 0 && i < len(items) {
			open.Set(false)
			if items[i].OnSelect != nil {
				items[i].OnSelect()
			}
		}
	}

	field := ui.TextField(ui.TextFieldProps{
		Value:       query,
		Placeholder: ui.Val(placeholder),
		Plain:       true,
		Keys: ui.KeyMap{
			{Key: "down"}: func() {
				if selected.Peek() < len(signal.Untrack(matches.Get))-1 {
					selected.Set(selected.Peek() + 1)
				}
			},
			{Key: "up"}: func() {
				if selected.Peek() > 0 {
					selected.Set(selected.Peek() - 1)
				}
			},
			{Key: "enter"}:  choose,
			{Key: "escape"}: func() { open.Set(false) },
		},
		Style: themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Body, Color: c.Foreground, Width: ui.Pct(100)}
		}),
	})

	trigger := Button(ButtonProps{
		Text: placeholder, Icon: ui.IconSearch, Variant: Outline,
		OnPress: func() { open.Set(true) },
	})

	content := func() ui.Node {
		return ui.Column(ui.Style{Gap: Space3, Width: ui.Pct(100), MinWidth: ui.Pt(320)},
			ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Padding: ui.Symmetric(Space3, Space4), BorderWidth: 1,
					BorderColor: c.Border, CornerRadius: t.Radii.MD, Width: ui.Pct(100)}
			})}, field),
			rows(),
		)
	}
	return ui.View(ui.ViewProps{}, trigger,
		Overlay(OverlayProps{Open: open, Content: content, Grabber: true}))
}
