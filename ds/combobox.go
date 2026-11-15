package ds

import (
	"strings"

	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// ComboboxItem is one choice of a Combobox.
type ComboboxItem struct {
	// Value is what Selected holds; it is never shown.
	Value    string
	Title    string
	Subtitle string
	// Keywords are extra words the search matches.
	Keywords string
	Icon     ui.Icon
}

// ComboboxProps configures a Combobox.
type ComboboxProps struct {
	Label       string
	Description string
	Error       string
	Placeholder string
	Items       []ComboboxItem
	// Selected is the chosen item's Value; empty means none.
	Selected *signal.Signal[string]
	// Empty is shown when nothing matches ("No results" by default).
	Empty    string
	Disabled bool
}

// Combobox is a Select with a search field: a themed trigger opens a
// popover with a native TextField and a filtered, keyboard-navigable list.
// It is built from Pressables and Asuka's TextField, so it works with a
// hardware keyboard on desktop, iPad and the web.
func Combobox(p ComboboxProps) ui.Node {
	query := signal.New("")
	selected := signal.New(0)
	open := signal.New(false)
	placeholder := p.Placeholder
	if placeholder == "" {
		placeholder = "Select…"
	}
	empty := p.Empty
	if empty == "" {
		empty = "No results"
	}

	matches := signal.Func(func() []ComboboxItem {
		q := strings.ToLower(strings.TrimSpace(query.Get()))
		var out []ComboboxItem
		for _, item := range p.Items {
			text := item.Title + " " + item.Subtitle + " " + item.Keywords
			if q == "" || strings.Contains(strings.ToLower(text), q) {
				out = append(out, item)
			}
		}
		return out
	})
	signal.Effect(func() {
		if selected.Peek() >= len(matches.Get()) {
			selected.Set(0)
		}
	})

	currentValue := func() string {
		if p.Selected == nil {
			return ""
		}
		return p.Selected.Get()
	}
	currentTitle := func() string {
		value := currentValue()
		if value == "" {
			return ""
		}
		for _, item := range p.Items {
			if item.Value == value {
				return item.Title
			}
		}
		return value
	}
	choose := func(value string) {
		if p.Selected != nil {
			p.Selected.Set(value)
		}
		open.Set(false)
	}
	chooseIndex := func() {
		items := signal.Untrack(matches.Get)
		if i := selected.Peek(); i >= 0 && i < len(items) {
			choose(items[i].Value)
		}
	}
	show := func() {
		if p.Disabled {
			return
		}
		query.Set("")
		selected.Set(0)
		open.Set(true)
	}

	trigger := Pressable(PressableProps{
		Label:    p.Label,
		Disabled: ui.Val(p.Disabled),
		OnPress:  show,
		Style: themed(func(t Theme, c Palette) ui.Style {
			s := ui.Style{
				Direction: ui.RowDir, Align: ui.AlignCenter, Gap: Space2,
				Padding: ui.Symmetric(0, Space3), Height: ui.Pt(36), Width: ui.Pct(100),
				Background: c.Background, BorderWidth: 1, BorderColor: c.Input,
				CornerRadius: t.Radii.MD,
			}
			if p.Error != "" {
				s.BorderColor = c.Destructive
			}
			return s
		}),
	},
		ui.Text(signal.Func(func() string {
			if title := currentTitle(); title != "" {
				return title
			}
			return placeholder
		}), ui.StyleFunc(func() ui.Style {
			theme, palette := ThemeNow(), PaletteNow()
			col := palette.Foreground
			if currentTitle() == "" {
				col = palette.MutedForeground
			}
			return ui.Style{Font: theme.Type.Body, Color: col, Lines: 1, Grow: 1}
		})),
		ui.Symbol(ui.IconDown, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Caption, Color: c.MutedForeground}
		})),
	)

	rows := func() ui.Node {
		return ui.Dynamic(func() ui.Node {
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
					OnPress:  func() { choose(entry.Value) },
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
		})
	}

	content := func() ui.Node {
		field := ui.TextField(ui.TextFieldProps{
			Value:       query,
			Placeholder: ui.Val("Search…"),
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
				{Key: "enter"}:  chooseIndex,
				{Key: "escape"}: func() { open.Set(false) },
			},
			Style: themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Font: t.Type.Body, Color: c.Foreground, Width: ui.Pct(100)}
			}),
		})
		return ui.Column(ui.Style{Gap: Space3, Width: ui.Pct(100), MinWidth: ui.Pt(260)},
			ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Padding: ui.Symmetric(Space3, Space4), BorderWidth: 1,
					BorderColor: c.Border, CornerRadius: t.Radii.MD, Width: ui.Pct(100)}
			})}, field),
			rows(),
		)
	}

	control := func() ui.Node {
		return Overlay(OverlayProps{
			Open: open, Popover: true, MaxWidth: 360,
			Content:   content,
			OnDismiss: func() { open.Set(false) },
		}, trigger)
	}
	return Field(FieldProps{
		Label:                p.Label,
		Description:          p.Description,
		DescriptionWithError: true,
		Error:                p.Error,
		Disabled:             ui.Val(p.Disabled),
		Control:              control,
	})
}
