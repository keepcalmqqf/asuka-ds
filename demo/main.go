package main

import (
	"fmt"
	"time"

	ds "asuka-ds"

	"github.com/egoist/asuka/app"
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

func main() { app.Run(App) }

func App() ui.Node {
	themeIndex := signal.New(0)
	appearance := signal.New(0)
	tab := signal.New(0)
	email := signal.New("")
	name := signal.New("")
	bio := signal.New("")
	workspace := signal.New("")
	notify := signal.New(true)
	checkbox := signal.New(false)
	plan := signal.New(1)
	framework := signal.New(0)
	volume := signal.New(0.62)
	page := signal.New(3)
	start := signal.New(time.Now())
	accordion := signal.New(0)
	sheetOpen := signal.New(false)

	signal.Effect(func() {
		if themeIndex.Get() == 0 {
			ds.SetTheme(ds.Native())
		} else {
			ds.SetTheme(ds.Neutral())
		}
	})
	signal.Effect(func() {
		ds.SetAppearance(ui.ColorScheme(appearance.Get()))
	})
	return ds.Screen(
		ui.Column(ui.Style{Gap: ds.Space2},
			ds.Heading("Asuka DS"),
			ds.Muted("A shadcn-style component library on Asuka's native primitives. Tokens first, components you own."),
		),

		ui.Row(ui.Style{Gap: ds.Space3, Wrap: true, Align: ui.AlignCenter},
			ui.View(ui.ViewProps{Style: ui.Style{Width: ui.Pt(260)}},
				ds.Tabs(ds.TabsProps{Items: []string{"Native", "Neutral"}, Selected: themeIndex})),
			ui.View(ui.ViewProps{Style: ui.Style{Width: ui.Pt(320)}},
				ds.Tabs(ds.TabsProps{Items: []string{"System", "Light", "Dark"}, Selected: appearance})),
		),

		Section("Buttons",
			ui.Row(ui.Style{Gap: ds.Space2, Wrap: true},
				ds.Button(ds.ButtonProps{Text: "Primary", Variant: ds.Primary, OnPress: func() {}}),
				ds.Button(ds.ButtonProps{Text: "Secondary", Variant: ds.Secondary, OnPress: func() {}}),
				ds.Button(ds.ButtonProps{Text: "Outline", Variant: ds.Outline, OnPress: func() {}}),
				ds.Button(ds.ButtonProps{Text: "Ghost", Variant: ds.Ghost, OnPress: func() {}}),
				ds.Button(ds.ButtonProps{Text: "Destructive", Variant: ds.Destructive, OnPress: func() {}}),
				ds.Button(ds.ButtonProps{Text: "Link", Variant: ds.Link, OnPress: func() {}}),
			),
			ui.Row(ui.Style{Gap: ds.Space2, Wrap: true, Align: ui.AlignCenter},
				ds.Button(ds.ButtonProps{Text: "Small", Size: ds.Small, Variant: ds.Primary, OnPress: func() {}}),
				ds.Button(ds.ButtonProps{Text: "Medium", Variant: ds.Primary, OnPress: func() {}}),
				ds.Button(ds.ButtonProps{Text: "Large", Size: ds.Large, Variant: ds.Primary, OnPress: func() {}}),
				ds.Button(ds.ButtonProps{Text: "With icon", Icon: ui.IconPlus, Variant: ds.Outline, OnPress: func() {}}),
				ds.Button(ds.ButtonProps{Icon: ui.IconMore, Size: ds.IconOnly, Variant: ds.Ghost, Label: "More options", OnPress: func() {}}),
				ds.Button(ds.ButtonProps{Text: "Disabled", Variant: ds.Primary, Disabled: true}),
			),
			ds.Card(ds.CardProps{}, ui.Column(ui.Style{Gap: ds.Space3},
				ds.CardHeader("Native or drawn?", "The same props, two implementations."),
				ui.Row(ui.Style{Gap: ds.Space3, Wrap: true, Align: ui.AlignCenter},
					ds.NativeButton(ds.ButtonProps{Text: "Native", Variant: ds.Primary, OnPress: func() {}}),
					ds.NativeButton(ds.ButtonProps{Text: "Native outline", Variant: ds.Outline, OnPress: func() {}}),
					ds.Button(ds.ButtonProps{Text: "Drawn", Variant: ds.Primary, OnPress: func() {}}),
				),
			)),
		),

		Section("Badges",
			ui.Row(ui.Style{Gap: ds.Space2, Wrap: true, Align: ui.AlignCenter},
				ds.Badge(ds.BadgeProps{Text: "Primary", Variant: ds.Primary}),
				ds.Badge(ds.BadgeProps{Text: "Secondary", Variant: ds.Secondary}),
				ds.Badge(ds.BadgeProps{Text: "Outline", Variant: ds.Outline}),
				ds.Badge(ds.BadgeProps{Text: "Destructive", Variant: ds.Destructive}),
				ds.Badge(ds.BadgeProps{Text: "Stable", Variant: ds.Outline, Dot: true}),
				ds.Badge(ds.BadgeProps{Text: "Beta", Variant: ds.Primary, Dot: true}),
			),
		),

		Section("Form",
			ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(420)}},
				ui.Column(ui.Style{Gap: ds.Space4},
					ds.Input(ds.InputProps{Value: name, Label: "Name", Placeholder: "Ada Lovelace",
						Description: "Shown on your profile."}),
					ds.Input(ds.InputProps{Value: email, Label: "Email", Placeholder: "ada@example.com",
						Icon: ui.IconMail, Error: "That address is already taken."}),
					ds.Textarea(ds.TextareaProps{Value: bio, Label: "Bio",
						Placeholder: "A few words about you"}),
					ds.Input(ds.InputProps{Value: workspace, Label: "Workspace", Placeholder: "acme", Disabled: true,
						Description: "Ask an admin to change it."}),
					ds.DatePicker(ds.DatePickerProps{Label: "Start date", Value: start,
						Description: "The platform's own picker."}),
					ui.View(ui.ViewProps{Label: ui.Val("Email notifications")},
						ui.Row(ui.Style{Gap: ds.Space3, Align: ui.AlignCenter},
							ui.Toggle(ui.ToggleProps{Value: notify}),
							ds.Text("Email notifications"),
						),
					),
					ui.Row(ui.Style{Gap: ds.Space2, Justify: ui.JustifyEnd},
						ds.Button(ds.ButtonProps{Text: "Cancel", Variant: ds.Ghost, OnPress: func() {}}),
						ds.Button(ds.ButtonProps{Text: "Save changes", Variant: ds.Primary, OnPress: func() {
							fmt.Println("saved", name.Peek(), email.Peek())
						}}),
					),
				),
			),
		),

		Section("Tabs",
			ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(420)}},
				ui.Column(ui.Style{Gap: ds.Space4},
					ds.Tabs(ds.TabsProps{Items: []string{"Overview", "Activity", "Settings"}, Selected: tab}),
					ds.Card(ds.CardProps{}, ui.Column(ui.Style{Gap: ds.Space2},
						ds.CardTitle(signal.Func(func() string {
							return []string{"Overview", "Activity", "Settings"}[tab.Get()]
						})),
						ds.Muted("The panel follows the tab's signal: only the text is re-sent."),
					)),
				),
			),
		),

		Section("Breadcrumb & pagination",
			ui.Column(ui.Style{Gap: ds.Space4},
				ds.Breadcrumb(ds.BreadcrumbProps{Items: []ds.BreadcrumbItem{
					{Title: "Projects", OnPress: func() {}},
					{Title: "asuka-ds", OnPress: func() {}},
					{Title: "Settings"},
				}}),
				ds.Pagination(ds.PaginationProps{Page: page, Total: 12}),
			),
		),

		Section("Surfaces",
			ui.Row(ui.Style{Gap: ds.Space4, Wrap: true, Align: ui.AlignStretch},
				ui.View(ui.ViewProps{Style: ui.Style{Grow: 1, Basis: ui.Pt(260)}},
					ds.Card(ds.CardProps{},
						ui.Column(ui.Style{Gap: ds.Space3},
							ui.Row(ui.Style{Justify: ui.JustifySpaceBetween, Align: ui.AlignCenter},
								ds.CardTitle(ui.Val("Deploy")),
								ds.Badge(ds.BadgeProps{Text: "Live", Variant: ds.Outline, Dot: true}),
							),
							ds.Muted("Your app is serving traffic from 3 regions."),
							ds.Skeleton(ui.Pct(100), 8),
							ds.Skeleton(ui.Pct(64), 8),
							ds.CardFooter(
								ds.Button(ds.ButtonProps{Text: "Logs", Variant: ds.Outline, Size: ds.Small, OnPress: func() {}}),
								ds.Button(ds.ButtonProps{Text: "Redeploy", Variant: ds.Primary, Size: ds.Small, OnPress: func() {}}),
							),
						),
					),
				),
				ui.View(ui.ViewProps{Style: ui.Style{Grow: 1, Basis: ui.Pt(260)}},
					ds.Card(ds.CardProps{},
						ui.Column(ui.Style{Gap: ds.Space3},
							ds.CardHeader("Usage", "This billing period"),
							ui.Row(ui.Style{Justify: ui.JustifySpaceBetween, Align: ui.AlignEnd},
								ds.Title("68%"),
								ds.Muted("of 100 GB"),
							),
							ui.View(ui.ViewProps{Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
								return ui.Style{Height: ui.Pt(8), CornerRadius: 4, Background: c.Muted}
							})},
								ui.View(ui.ViewProps{Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
									return ui.Style{Width: ui.Pct(68), Height: ui.Pt(8), CornerRadius: 4, Background: c.Primary}
								})}),
							),
						),
					),
				),
			),
		),

		Section("Selection controls",
			ui.Row(ui.Style{Gap: ds.Space5, Wrap: true, Align: ui.AlignStart},
				ui.View(ui.ViewProps{Style: ui.Style{Grow: 1, Basis: ui.Pt(300)}},
					ui.Column(ui.Style{Gap: ds.Space4},
						ds.Switch(ds.SwitchProps{Label: "Email notifications", Value: notify,
							Description: "Weekly summaries and mentions."}),
						ds.Switch(ds.SwitchProps{Label: "Agree to terms", Value: checkbox, Checkbox: true}),
						ds.RadioGroup(ds.RadioGroupProps{Label: "Plan", Options: []string{"Free", "Pro", "Team"},
							Selected: plan, Description: "Change at any time."}),
						ds.Select(ds.SelectProps{Label: "Framework", Options: []string{"Asuka", "SwiftUI", "Compose"},
							Selected: framework, Placeholder: "Choose…"}),
						ds.Slider(ds.SliderProps{Label: "Volume", Value: volume, Min: 0, Max: 1}),
					),
				),
				ui.View(ui.ViewProps{Style: ui.Style{Grow: 1, Basis: ui.Pt(300)}},
					ui.Column(ui.Style{Gap: ds.Space4},
						ds.Accordion(ds.AccordionProps{Open: accordion,
							Items: []ds.AccordionItem{
								{Title: "What is Asuka?", Content: func() ui.Node {
									return ds.Text("A Go UI framework that renders each platform's own widgets.")
								}},
								{Title: "Why two button backends?", Content: func() ui.Node {
									return ds.Text("Drawn buttons keep a brand consistent; native buttons keep platform feel.")
								}},
							},
						}),
						ds.Alert(ds.AlertProps{Title: "Heads up", Variant: ds.Primary,
							Description: "M1 through M3 are implemented in this library."}),
					),
				),
			),
		),

		Section("Data & feedback",
			ui.Row(ui.Style{Gap: ds.Space5, Wrap: true, Align: ui.AlignStart},
				ui.View(ui.ViewProps{Style: ui.Style{Grow: 1, Basis: ui.Pt(300)}},
					ui.Column(ui.Style{Gap: ds.Space4},
						ui.Row(ui.Style{Gap: ds.Space3, Align: ui.AlignCenter},
							ds.Avatar(ds.AvatarProps{Name: "Ada Lovelace"}),
							ds.Avatar(ds.AvatarProps{Name: "Grace Hopper", Size: 32}),
							ds.Badge(ds.BadgeProps{Text: "v1.0", Variant: ds.Secondary}),
						),
						ds.Progress(ds.ProgressProps{Value: ui.Val(0.68), Label: "Storage", ShowValue: true}),
						ds.Progress(ds.ProgressProps{Value: ui.Val(float64(-1)), Label: "Indexing"}),
						ui.Row(ui.Style{Gap: ds.Space3, Align: ui.AlignCenter},
							ds.Spinner(ds.SpinnerProps{Label: "Syncing"}),
							ds.Spinner(ds.SpinnerProps{Size: 16}),
						),
						ui.Row(ui.Style{Gap: ds.Space2, Wrap: true},
							ds.Button(ds.ButtonProps{Text: "Toast", Variant: ds.Outline, Size: ds.Small,
								OnPress: func() { ds.Toast("Saved", &ds.DialogAction{Title: "Undo"}) }}),
							ds.Button(ds.ButtonProps{Text: "Dialog", Variant: ds.Outline, Size: ds.Small,
								OnPress: func() {
									ds.ConfirmDialog("Delete project?", "This cannot be undone.", "Delete", func() {})
								}}),
							ds.Button(ds.ButtonProps{Text: "Sheet", Variant: ds.Outline, Size: ds.Small,
								OnPress: func() { sheetOpen.Set(true) }}),
						),
						ds.Command(ds.CommandProps{Items: []ds.CommandItem{
							{Title: "New project", Icon: ui.IconPlus, Keywords: "create add"},
							{Title: "Open settings", Icon: ui.IconSettings, Keywords: "preferences"},
							{Title: "Sign out", Icon: ui.IconSignOut, Keywords: "logout"},
						}}),
					),
				),
				ui.View(ui.ViewProps{Style: ui.Style{Grow: 1, Basis: ui.Pt(360)}},
					ds.Table(ds.TableProps{
						Columns: []ds.TableColumn{{Title: "Name", Width: 2}, {Title: "Role"}, {Title: "Status"}},
						Rows: [][]string{
							{"Ada Lovelace", "Engineer", "Active"},
							{"Grace Hopper", "Admiral", "Active"},
							{"Alan Turing", "Researcher", "Away"},
						},
						Striped: true,
					}),
				),
			),
		),

		Section("Empty states & shortcuts",
			ui.Column(ui.Style{Gap: ds.Space4},
				ui.Row(ui.Style{Gap: ds.Space2, Align: ui.AlignCenter},
					ds.Kbd("⌘", "K"),
					ds.Muted("opens the command palette"),
					ds.Kbd("⌘", "⇧", "P"),
				),
				ds.EmptyState(ds.EmptyStateProps{
					Icon:        ui.IconFolder,
					Title:       "No projects yet",
					Description: "Create a project to deploy your first app.",
					Action: func() ui.Node {
						return ds.Button(ds.ButtonProps{Text: "New project", Icon: ui.IconPlus,
							Variant: ds.Primary, OnPress: func() {}})
					},
				}),
			),
		),

		ds.Overlay(ds.OverlayProps{Open: sheetOpen, Medium: true, Grabber: true,
			Content: func() ui.Node {
				return ui.Column(ui.Style{Gap: ds.Space4, Width: ui.Pct(100)},
					ds.Title("Edit profile"),
					ds.Muted("This is a native sheet with design-system content."),
					ds.Button(ds.ButtonProps{Text: "Close", Variant: ds.Primary, FullWidth: true,
						OnPress: func() { sheetOpen.Set(false) }}),
				)
			},
		}),
	)
}

func Section(title string, children ...ui.Node) ui.Node {
	return ui.Column(ui.Style{Gap: ds.Space4},
		ui.Column(ui.Style{Gap: ds.Space1},
			ds.Title(title),
			ds.Separator(false),
		),
		ui.Column(ui.Style{Gap: ds.Space4}, children...),
	)
}
