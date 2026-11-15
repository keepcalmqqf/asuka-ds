package main

import (
	ds "github.com/egoist/asuka-ds/ds"
	"github.com/egoist/asuka-ds/internal/catalog"
	"github.com/egoist/asuka/app"
	"github.com/egoist/asuka/colors"
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

	env := ui.UseEnvironment()
	return ui.View(ui.ViewProps{Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
		return ui.Style{Grow: 1, Background: c.Background}
	})},
		ui.SafeArea(ui.Style{Grow: 1},
			ui.Column(ui.Style{Grow: 1},
				header(d, themeIndex, appearance),
				ui.Row(ui.Style{Grow: 1, Align: ui.AlignStretch},
					sidebar(d),
					ui.Scroll(ui.ScrollProps{Style: ui.Style{Grow: 1}},
						ui.View(ui.ViewProps{Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
							pad := ui.Edges{Top: ds.Space7, Right: ds.Space6, Bottom: ds.Space8, Left: ds.Space6}
							if env.Get().WidthClass == ui.SizeClassCompact {
								pad = ui.Edges{Top: ds.Space5, Right: ds.Space4, Bottom: ds.Space6, Left: ds.Space4}
							}
							return ui.Style{
								Width: ui.Pct(100), MaxWidth: ui.Pt(820), AlignSelf: ui.AlignCenter,
								Padding: pad, Direction: ui.ColumnDir,
							}
						})}, ui.Dynamic(func() ui.Node { return pageFor(d, d.page.Get()) })),
					),
				),
			),
		),
	)
}

func header(d *docs, themeIndex, appearance *signal.Signal[int]) ui.Node {
	env := ui.UseEnvironment()
	return ui.Column(ui.Style{},
		ui.View(ui.ViewProps{Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
			style := ui.Style{
				Direction: ui.RowDir, Wrap: true, Align: ui.AlignCenter,
				Justify: ui.JustifySpaceBetween, Gap: ds.Space4,
				Padding: ui.Symmetric(ds.Space3, ds.Space5),
			}
			if env.Get().WidthClass == ui.SizeClassCompact {
				style.Gap = ds.Space2
				style.Padding = ui.Symmetric(ds.Space2, ds.Space4)
			}
			return style
		})},
			ui.Row(ui.Style{Gap: ds.Space3, Align: ui.AlignCenter},
				mark(),
				ui.StaticText("Asuka DS", ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
					return ui.Style{Font: t.Type.Title, Color: c.Foreground, Lines: 1}
				})),
			),
			ui.Row(ui.Style{Gap: ds.Space2, Align: ui.AlignCenter, Wrap: true},
				compactBrowse(d),
				ui.View(ui.ViewProps{Style: ui.Style{Width: ui.Pt(compactWidth(176, env))}},
					ds.Tabs(ds.TabsProps{Items: []string{"Native", "Neutral"}, Selected: themeIndex})),
				ui.View(ui.ViewProps{Style: ui.Style{Width: ui.Pt(compactWidth(228, env))}},
					ds.Tabs(ds.TabsProps{Items: []string{"System", "Light", "Dark"}, Selected: appearance})),
				link("GitHub", "https://github.com/egoist/asuka-ds"),
			),
		),
		ds.Separator(false),
	)
}

func compactWidth(expanded float32, env signal.Accessor[ui.Environment]) float32 {
	if env.Get().WidthClass == ui.SizeClassCompact {
		return expanded * 0.76
	}
	return expanded
}

// mark is four live theme swatches. It changes with Native/Neutral and
// light/dark, so the wordmark is itself a token specimen.
func mark() ui.Node {
	swatch := func(color func(ds.Palette) ui.Color) ui.Node {
		return ui.View(ui.ViewProps{Decorative: true, Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
			return ui.Style{Width: ui.Pt(8), Height: ui.Pt(8), CornerRadius: 2, Background: color(c)}
		})})
	}
	return ui.View(ui.ViewProps{Decorative: true, Style: ui.Style{Direction: ui.RowDir, Gap: 2}},
		ui.Column(ui.Style{Gap: 2},
			swatch(func(c ds.Palette) ui.Color { return c.Primary }),
			swatch(func(c ds.Palette) ui.Color { return c.Foreground }),
		),
		ui.Column(ui.Style{Gap: 2},
			swatch(func(c ds.Palette) ui.Color { return c.Focus }),
			swatch(func(c ds.Palette) ui.Color { return c.Border }),
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
			ui.View(ui.ViewProps{Style: ui.Style{
				Width: ui.Pt(232), Shrink: ui.NoShrink, Height: ui.Pct(100),
				Padding: ui.Symmetric(ds.Space5, ds.Space3),
			}}, ui.Scroll(ui.ScrollProps{Style: ui.Style{Grow: 1}},
				ui.Column(ui.Style{Gap: ds.Space5}, kids...))),
			ds.Separator(true),
		)
	})
}

func navSection(title string, items ...ui.Node) ui.Node {
	return ui.Column(ui.Style{Gap: ds.Space1},
		ui.View(ui.ViewProps{Style: ui.Style{Padding: ui.Edges{Left: 19, Right: ds.Space2, Bottom: 2}}},
			ui.StaticText(title, ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
				return ui.Style{Font: t.Type.Caption.Weighted(fonts.Medium), Color: c.MutedForeground, Lines: 1}
			}))),
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
			s := ui.Style{
				Direction: ui.RowDir, Align: ui.AlignCenter, Gap: ds.Space2, Width: ui.Pct(100),
				Padding:      ui.Edges{Top: 7, Right: ds.Space2, Bottom: 7, Left: ds.Space2},
				CornerRadius: t.Radii.SM,
			}
			if selected.Get() {
				s.Background = c.Muted
			} else {
				s.Hovered = &ui.Style{Background: c.Muted}
			}
			return s
		}),
	},
		ui.View(ui.ViewProps{Decorative: true, Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
			s := ui.Style{Width: ui.Pt(3), Height: ui.Pt(14), CornerRadius: 2, Background: colors.Clear}
			if selected.Get() {
				s.Background = c.Primary
			}
			return s
		})}),
		ui.StaticText(label, ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
			font := t.Type.Small
			col := c.MutedForeground
			if selected.Get() {
				font = font.Weighted(fonts.Medium)
				col = c.Foreground
			}
			return ui.Style{Font: font, Color: col, Lines: 1}
		})),
	)
}

func link(label, href string) ui.Node {
	return ui.View(ui.ViewProps{
		Label: ui.Val(label),
		Role:  ui.RoleLink,
		Href:  ui.Href(href),
		Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
			return ui.Style{
				Padding: ui.Symmetric(6, ds.Space3), CornerRadius: t.Radii.SM,
				BorderWidth: 1, BorderColor: c.Border,
				Hovered: &ui.Style{Background: c.Muted},
			}
		}),
	}, ui.StaticText(label, ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
		return ui.Style{Font: t.Type.Small, Color: c.Foreground, Lines: 1}
	})))
}

func code(language, src string) ui.Node {
	return ui.View(ui.ViewProps{Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
		return ui.Style{
			Padding: ui.All(ds.Space4), Gap: ds.Space2, Clip: true,
			Background: c.Card, BorderWidth: 1, BorderColor: c.Border, CornerRadius: t.Radii.MD,
		}
	})},
		ui.Row(ui.Style{Justify: ui.JustifySpaceBetween, Align: ui.AlignCenter, Width: ui.Pct(100)},
			ui.StaticText(language, ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
				return ui.Style{Font: fontsMono(t, 12).Weighted(fonts.Medium), Color: c.MutedForeground, Lines: 1}
			})),
			copyControl(src),
		),
		ds.Separator(false),
		ui.RichText(ui.RichTextProps{
			Text: signal.Func(func() ui.Rich { return highlight(language, src) }),
			Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
				return ui.Style{Font: fontsMono(t, 13), Color: c.Foreground}
			}),
		}),
	)
}

func copyControl(src string) ui.Node {
	return ds.Pressable(ds.PressableProps{
		Label: "Copy",
		OnPress: func() {
			ui.Copy(src)
			ds.Toast("Copied", nil)
		},
		Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
			return ui.Style{
				Padding: ui.Symmetric(ds.Space1, ds.Space2), CornerRadius: t.Radii.SM,
				Hovered: &ui.Style{Background: c.Muted},
			}
		}),
	}, ui.StaticText("Copy", ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
		return ui.Style{Font: t.Type.Small.Weighted(fonts.Medium), Color: c.Foreground, Lines: 1}
	})))
}

// stage is the preview ground: a muted field so components, which are
// usually card-colored, read as objects on a surface.
func stage(child ui.Node) ui.Node {
	return ui.View(ui.ViewProps{Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
		return ui.Style{
			Padding: ui.All(ds.Space6), Background: c.Muted,
			BorderWidth: 1, BorderColor: c.Border, CornerRadius: t.Radii.LG,
		}
	})}, child)
}

// lead is a page's supporting sentence, held to a reading measure.
func lead(text string) ui.Node {
	return ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(640)}},
		ui.StaticText(text, ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
			return ui.Style{Font: t.Type.Body, Color: c.MutedForeground}
		})))
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
