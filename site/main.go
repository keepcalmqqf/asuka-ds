package main

import (
	ds "github.com/egoist/asuka-ds"
	"github.com/egoist/asuka-ds/internal/catalog"
	"github.com/egoist/asuka/app"
	"github.com/egoist/asuka/fonts"
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

func main() { app.Run(App) }

type docs struct {
	page    *signal.Signal[string]
	catalog catalog.Catalog
}

func App() ui.Node {
	c, err := catalog.Load()
	if err != nil {
		return ds.Alert(ds.AlertProps{Title: "Catalog error", Description: err.Error(), Variant: ds.Destructive})
	}
	d := &docs{page: signal.New("intro"), catalog: c}
	themeIndex := signal.New(0)
	appearance := signal.New(0)
	signal.Effect(func() {
		if themeIndex.Get() == 0 {
			ds.SetTheme(ds.Native())
			return
		}
		ds.SetTheme(ds.Neutral())
	})
	signal.Effect(func() { ds.SetAppearance(ui.ColorScheme(appearance.Get())) })

	return ui.View(ui.ViewProps{Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
		return ui.Style{Grow: 1, Background: c.Background}
	})},
		ui.SafeArea(ui.Style{Grow: 1},
			ui.Column(ui.Style{Grow: 1},
				header(d, themeIndex, appearance),
				ui.Row(ui.Style{Grow: 1},
					sidebar(d),
					ui.Scroll(ui.ScrollProps{Style: ui.Style{Grow: 1}},
						ui.View(ui.ViewProps{Style: ui.Style{
							Width: ui.Pct(100), MaxWidth: ui.Pt(900), AlignSelf: ui.AlignCenter,
							Padding: ui.All(ds.Space6), Gap: ds.Space6, Direction: ui.ColumnDir,
						}}, ui.Dynamic(func() ui.Node { return pageFor(d, d.page.Get()) })),
					),
				),
			),
		),
	)
}

func header(d *docs, themeIndex, appearance *signal.Signal[int]) ui.Node {
	return ui.View(ui.ViewProps{Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
		return ui.Style{
			Direction: ui.RowDir, Wrap: true, Align: ui.AlignCenter,
			Justify: ui.JustifySpaceBetween, Gap: ds.Space3,
			Padding:    ui.Symmetric(ds.Space3, ds.Space5),
			Background: c.Card, BorderWidth: 1, BorderColor: c.Border,
		}
	})},
		ui.Row(ui.Style{Gap: ds.Space3, Align: ui.AlignCenter},
			ui.View(ui.ViewProps{Decorative: true, Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
				return ui.Style{
					Width: ui.Pt(30), Height: ui.Pt(30), CornerRadius: 8,
					Align: ui.AlignCenter, Justify: ui.JustifyCenter, Background: c.Primary,
				}
			})}, ui.Symbol("square.grid.2x2.fill", ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
				return ui.Style{Font: t.Type.Caption, Color: c.PrimaryForeground}
			}))),
			ui.Column(ui.Style{Gap: 0},
				ds.Text("Asuka DS"),
				ds.Muted("Components for Asuka"),
			),
		),
		ui.Row(ui.Style{Gap: ds.Space3, Align: ui.AlignCenter, Wrap: true},
			compactBrowse(d),
			ui.View(ui.ViewProps{Style: ui.Style{Width: ui.Pt(200)}},
				ds.Tabs(ds.TabsProps{Items: []string{"Native", "Neutral"}, Selected: themeIndex})),
			ui.View(ui.ViewProps{Style: ui.Style{Width: ui.Pt(250)}},
				ds.Tabs(ds.TabsProps{Items: []string{"System", "Light", "Dark"}, Selected: appearance})),
			link("GitHub", "https://github.com/egoist/asuka-ds"),
		),
	)
}

// compactBrowse is the phone navigation: one native menu with a Guide
// section and one inline section per component category.
func compactBrowse(d *docs) ui.Node {
	env := ui.UseEnvironment()
	return ui.Show(signal.Func(func() bool { return env.Get().WidthClass == ui.SizeClassCompact }), func() ui.Node {
		var items []ds.DialogAction
		guide := []ds.DialogAction{
			{Title: "Introduction", OnPress: func() { d.page.Set("intro") }},
			{Title: "Installation", OnPress: func() { d.page.Set("install") }},
			{Title: "Theming", OnPress: func() { d.page.Set("themes") }},
			{Title: "CLI & registry", OnPress: func() { d.page.Set("cli") }},
			{Title: "Quality & status", OnPress: func() { d.page.Set("quality") }},
			{Title: "Platforms", OnPress: func() { d.page.Set("platforms") }},
		}
		items = append(items, ds.DialogAction{Title: "Guide", Inline: true, Children: guide})
		for _, category := range d.catalog.Categories() {
			var children []ds.DialogAction
			for _, name := range category.Components {
				component := d.catalog.Components[name]
				id := "c:" + name
				children = append(children, ds.DialogAction{
					Title: component.Title,
					OnPress: func() {
						d.page.Set(id)
					},
				})
			}
			items = append(items, ds.DialogAction{Title: category.Name, Inline: true, Children: children})
		}
		return ds.Menu(ds.MenuProps{Label: "Browse", Items: items},
			ds.Button(ds.ButtonProps{Text: "Browse", Variant: ds.Outline, Size: ds.Small}))
	})
}

func sidebar(d *docs) ui.Node {
	env := ui.UseEnvironment()
	return ui.Show(signal.Func(func() bool { return env.Get().WidthClass != ui.SizeClassCompact }), func() ui.Node {
		var kids []ui.Node
		kids = append(kids, navSection("Guide",
			navItem(d, "intro", "Introduction"),
			navItem(d, "install", "Installation"),
			navItem(d, "themes", "Theming"),
			navItem(d, "cli", "CLI & registry"),
			navItem(d, "quality", "Quality & status"),
			navItem(d, "platforms", "Platforms"),
		))
		for _, category := range d.catalog.Categories() {
			items := make([]ui.Node, 0, len(category.Components))
			for _, name := range category.Components {
				items = append(items, navItem(d, "c:"+name, d.catalog.Components[name].Title))
			}
			kids = append(kids, navSection(category.Name, items...))
		}
		return ui.Row(ui.Style{Height: ui.Pct(100)},
			ui.View(ui.ViewProps{Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
				return ui.Style{Width: ui.Pt(250), Shrink: ui.NoShrink, Padding: ui.All(ds.Space4), Background: c.Card}
			})}, ui.Scroll(ui.ScrollProps{Style: ui.Style{Grow: 1}},
				ui.Column(ui.Style{Gap: ds.Space5}, kids...))),
			ds.Separator(true),
		)
	})
}

func navSection(title string, items ...ui.Node) ui.Node {
	return ui.Column(ui.Style{Gap: ds.Space1},
		ui.View(ui.ViewProps{Style: ui.Style{Padding: ui.Symmetric(0, ds.Space2)}},
			ds.Muted(title)),
		ui.Column(ui.Style{Gap: 2}, items...),
	)
}

func navItem(d *docs, id, label string) ui.Node {
	selected := signal.Func(func() bool { return d.page.Get() == id })
	return ds.Pressable(ds.PressableProps{
		Label:    label,
		Selected: selected,
		OnPress:  func() { d.page.Set(id) },
		Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
			s := ui.Style{Padding: ui.Symmetric(ds.Space2, ds.Space2), CornerRadius: t.Radii.SM}
			if selected.Get() {
				s.Background = c.Muted
			} else {
				s.Hovered = &ui.Style{Background: c.Muted}
			}
			return s
		}),
	}, ui.StaticText(label, ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
		col := c.MutedForeground
		if selected.Get() {
			col = c.Foreground
		}
		return ui.Style{Font: t.Type.Small, Color: col, Lines: 1}
	})))
}

func link(label, href string) ui.Node {
	return ui.View(ui.ViewProps{
		Label: ui.Val(label),
		Role:  ui.RoleLink,
		Href:  ui.Href(href),
		Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
			return ui.Style{Padding: ui.Symmetric(ds.Space2, ds.Space3), CornerRadius: t.Radii.SM,
				Hovered: &ui.Style{Background: c.Muted}}
		}),
	}, ds.Text(label))
}

func code(src string) ui.Node {
	return ui.View(ui.ViewProps{Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
		return ui.Style{
			Padding: ui.All(ds.Space4), Background: c.Muted,
			BorderWidth: 1, BorderColor: c.Border, CornerRadius: t.Radii.MD,
			Gap: ds.Space3,
		}
	})},
		ui.StaticText(src, ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
			return ui.Style{Font: fontsMono(t, 12), Color: c.Foreground}
		})),
		ui.Row(ui.Style{Justify: ui.JustifyEnd},
			ds.Button(ds.ButtonProps{Text: "Copy", Variant: ds.Ghost, Size: ds.Small,
				OnPress: func() {
					ui.Copy(src)
					ds.Toast("Copied", nil)
				}}),
		),
	)
}

func fontsMono(t ds.Theme, size float32) (f fonts.Font) {
	return fonts.Custom("monospace", fonts.System(size))
}

func pageFor(d *docs, id string) ui.Node {
	if id == "intro" {
		return introPage(d)
	}
	if id == "install" {
		return installPage()
	}
	if id == "themes" {
		return themesPage()
	}
	if id == "cli" {
		return cliPage()
	}
	if id == "quality" {
		return qualityPage()
	}
	if id == "platforms" {
		return platformsPage()
	}
	if len(id) > 2 && id[:2] == "c:" {
		return componentPage(d, id[2:])
	}
	return ds.Alert(ds.AlertProps{Title: "Page not found", Description: id, Variant: ds.Destructive})
}
