package main

import (
	"slices"
	"time"

	ds "github.com/egoist/asuka-ds/ds"
	"github.com/egoist/asuka/colors"
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

func row(children ...ui.Node) ui.Node {
	return ui.Row(ui.Style{Gap: ds.Space2, Wrap: true, Align: ui.AlignCenter}, children...)
}

func col(children ...ui.Node) ui.Node {
	return ui.Column(ui.Style{Gap: ds.Space3}, children...)
}

func exampleFor(name string) (ui.Node, string, bool) {
	switch name {
	case "button":
		return row(
				ds.Button(ds.ButtonProps{Text: "Primary", Variant: ds.Primary, OnPress: func() {}}),
				ds.Button(ds.ButtonProps{Text: "Secondary", Variant: ds.Secondary, OnPress: func() {}}),
				ds.Button(ds.ButtonProps{Text: "Outline", Variant: ds.Outline, OnPress: func() {}}),
				ds.Button(ds.ButtonProps{Text: "Ghost", Variant: ds.Ghost, OnPress: func() {}}),
				ds.Button(ds.ButtonProps{Text: "Delete", Variant: ds.Destructive, OnPress: func() {}}),
				ds.NativeButton(ds.ButtonProps{Text: "Native", Variant: ds.Primary, OnPress: func() {}}),
			), `ds.Button(ds.ButtonProps{Text: "Save", Variant: ds.Primary, OnPress: save})
ds.NativeButton(ds.ButtonProps{Text: "Pay", OnPress: pay})`, true
	case "input":
		email := signal.New("")
		return ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(420)}},
				ds.Input(ds.InputProps{Value: email, Label: "Email", Placeholder: "ada@example.com",
					Description: "We never share it."}),
			), `email := signal.New("")
ds.Input(ds.InputProps{Value: email, Label: "Email",
    Placeholder: "ada@example.com"})`, true
	case "textarea":
		bio := signal.New("")
		return ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(420)}},
				ds.Textarea(ds.TextareaProps{Value: bio, Label: "Bio", Lines: 4}),
			), `bio := signal.New("")
ds.Textarea(ds.TextareaProps{Value: bio, Label: "Bio", Lines: 4})`, true
	case "number-field":
		seats := signal.New(3.0)
		return ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(420)}},
				ds.NumberField(ds.NumberFieldProps{Label: "Seats", Value: seats, Min: 1, Max: 20}),
			), `seats := signal.New(3.0)
ds.NumberField(ds.NumberFieldProps{
    Label: "Seats", Value: seats, Min: 1, Max: 20,
})`, true
	case "date-picker":
		when := signal.New(time.Now())
		return ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(420)}},
				ds.DatePicker(ds.DatePickerProps{Label: "Start date", Value: when}),
			), `when := signal.New(time.Now())
ds.DatePicker(ds.DatePickerProps{Label: "Start date", Value: when})`, true
	case "date-range-picker":
		start, end := signal.New(time.Now()), signal.New(time.Now().AddDate(0, 0, 30))
		return ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(520)}},
				ds.DateRangePicker(ds.DateRangePickerProps{Label: "Billing period", Start: start, End: end}),
			), `start := signal.New(time.Now())
end := signal.New(time.Now().AddDate(0, 0, 30))
ds.DateRangePicker(ds.DateRangePickerProps{
    Label: "Billing period", Start: start, End: end,
})`, true
	case "switch":
		notify := signal.New(true)
		return ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(420)}},
				ds.Switch(ds.SwitchProps{Label: "Notifications", Value: notify,
					Description: "Weekly summaries and mentions."}),
			), `notify := signal.New(true)
ds.Switch(ds.SwitchProps{Label: "Notifications", Value: notify})`, true
	case "checkbox":
		ok := signal.New(true)
		return ds.Switch(ds.SwitchProps{Label: "Agree to terms", Value: ok, Checkbox: true}),
			`ok := signal.New(true)
ds.Switch(ds.SwitchProps{Label: "Agree to terms", Value: ok, Checkbox: true})`, true
	case "radio-group":
		plan := signal.New(1)
		return ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(420)}},
				ds.RadioGroup(ds.RadioGroupProps{Label: "Plan", Options: []string{"Free", "Pro", "Team"}, Selected: plan}),
			), `plan := signal.New(1)
ds.RadioGroup(ds.RadioGroupProps{
    Label: "Plan", Options: []string{"Free", "Pro", "Team"}, Selected: plan,
})`, true
	case "select":
		framework := signal.New(0)
		return ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(420)}},
				ds.Select(ds.SelectProps{Label: "Framework", Options: []string{"Asuka", "SwiftUI", "Compose"},
					Selected: framework, Placeholder: "Choose…"}),
			), `framework := signal.New(0)
ds.Select(ds.SelectProps{
    Label: "Framework", Options: []string{"Asuka", "SwiftUI"},
    Selected: framework, Placeholder: "Choose…",
})`, true
	case "combobox":
		selected := signal.New("")
		return ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(420)}},
				ds.Combobox(ds.ComboboxProps{Label: "Preferred framework", Placeholder: "Search frameworks…",
					Selected: selected,
					Items: []ds.ComboboxItem{
						{Value: "asuka", Title: "Asuka", Subtitle: "Go, native widgets"},
						{Value: "swiftui", Title: "SwiftUI", Subtitle: "Apple platforms"},
						{Value: "compose", Title: "Compose", Subtitle: "Android"},
					}}),
			), `selected := signal.New("")
ds.Combobox(ds.ComboboxProps{
    Label: "Framework", Selected: selected,
    Items: []ds.ComboboxItem{{Value: "asuka", Title: "Asuka"}},
})`, true
	case "color-picker":
		tint := signal.New(colors.Hex(0x3366cc))
		return ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(320)}},
				ds.ColorPicker(ds.ColorPickerProps{Label: "Accent color", Value: tint}),
			), `tint := signal.New(colors.Hex(0x3366cc))
ds.ColorPicker(ds.ColorPickerProps{Label: "Accent color", Value: tint})`, true
	case "segmented-control":
		density := signal.New(1)
		return ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(420)}},
				ds.SegmentedControl(ds.SegmentedControlProps{Label: "Density",
					Segments: []string{"Compact", "Regular", "Comfortable"}, Selected: density}),
			), `density := signal.New(1)
ds.SegmentedControl(ds.SegmentedControlProps{
    Label: "Density", Segments: []string{"Compact", "Regular"}, Selected: density,
})`, true
	case "slider":
		volume := signal.New(0.62)
		return ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(420)}},
				ds.Slider(ds.SliderProps{Label: "Volume", Value: volume, Min: 0, Max: 1}),
			), `volume := signal.New(0.62)
ds.Slider(ds.SliderProps{Label: "Volume", Value: volume, Min: 0, Max: 1})`, true
	case "toggle":
		bold := signal.New(true)
		return row(
				ds.Toggle(ds.ToggleProps{Text: "Bold", Pressed: bold, Variant: ds.Outline}),
				ds.Toggle(ds.ToggleProps{Text: "Italic", Pressed: signal.New(false), Variant: ds.Outline}),
				ds.Toggle(ds.ToggleProps{Icon: ui.IconStar, Label: "Star", Pressed: signal.New(true), Variant: ds.Ghost}),
			), `pressed := signal.New(true)
ds.Toggle(ds.ToggleProps{Text: "Bold", Pressed: pressed, Variant: ds.Outline})`, true
	case "toggle-group":
		formatting := signal.NewFunc([]string{"bold"}, slices.Equal)
		return ds.ToggleGroup(ds.ToggleGroupProps{Label: "Formatting", Selected: formatting,
				Items: []ds.ToggleItem{
					{Value: "bold", Label: "Bold"},
					{Value: "italic", Label: "Italic"},
					{Value: "code", Label: "Code", Icon: ui.IconDocument},
				}}),
			`formatting := signal.NewFunc([]string{"bold"}, slices.Equal)
ds.ToggleGroup(ds.ToggleGroupProps{
    Label: "Formatting", Selected: formatting,
    Items: []ds.ToggleItem{{Value: "bold", Label: "Bold"}},
})`, true
	case "command":
		return ds.Command(ds.CommandProps{Items: []ds.CommandItem{
				{Title: "New project", Icon: ui.IconPlus, Keywords: "create add"},
				{Title: "Open settings", Icon: ui.IconSettings, Keywords: "preferences"},
				{Title: "Sign out", Icon: ui.IconSignOut, Keywords: "logout"},
			}}),
			`ds.Command(ds.CommandProps{Items: []ds.CommandItem{
    {Title: "New project", Icon: ui.IconPlus, Keywords: "create"},
    {Title: "Open settings", Icon: ui.IconSettings},
}})`, true
	case "tabs":
		tab := signal.New(0)
		return ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(460)}},
				ds.Tabs(ds.TabsProps{Items: []string{"Overview", "Activity", "Settings"}, Selected: tab}),
			), `tab := signal.New(0)
ds.Tabs(ds.TabsProps{Items: []string{"Overview", "Activity"}, Selected: tab})`, true
	case "accordion":
		open := signal.New(0)
		return ds.Accordion(ds.AccordionProps{Open: open, Items: []ds.AccordionItem{
				{Title: "What is Asuka?", Content: func() ui.Node {
					return ds.Text("A Go UI framework that renders each platform's own widgets.")
				}},
				{Title: "Why native controls?", Content: func() ui.Node {
					return ds.Text("IME, focus, selection and screen readers come for free.")
				}},
			}}),
			`open := signal.New(0)
ds.Accordion(ds.AccordionProps{
    Open: open,
    Items: []ds.AccordionItem{{Title: "Details", Content: func() ui.Node {
        return ds.Text("The fine print")
    }}},
})`, true
	case "breadcrumb":
		return ds.Breadcrumb(ds.BreadcrumbProps{Items: []ds.BreadcrumbItem{
				{Title: "Projects", OnPress: func() {}},
				{Title: "asuka-ds", OnPress: func() {}},
				{Title: "Button"},
			}}),
			`ds.Breadcrumb(ds.BreadcrumbProps{Items: []ds.BreadcrumbItem{
    {Title: "Projects", OnPress: goProjects},
    {Title: "Button"},
}})`, true
	case "pagination":
		page := signal.New(3)
		return ds.Pagination(ds.PaginationProps{Page: page, Total: 12}),
			`page := signal.New(3)
ds.Pagination(ds.PaginationProps{Page: page, Total: 12})`, true
	case "dialog":
		return ds.Button(ds.ButtonProps{Text: "Delete project", Variant: ds.Destructive, OnPress: func() {
				ds.ConfirmDialog("Delete project?", "This cannot be undone.", "Delete", func() {})
			}}),
			`ds.ConfirmDialog("Delete project?", "This cannot be undone.",
    "Delete", deleteProject)`, true
	case "confirm-dialog":
		return ds.Button(ds.ButtonProps{Text: "Show confirmation", Variant: ds.Outline, OnPress: func() {
				ds.ConfirmDialog("Publish changes?", "They go live immediately.", "Publish", func() {})
			}}),
			`ds.ConfirmDialog("Publish changes?", "They go live immediately.",
    "Publish", publish)`, true
	case "menu":
		return ds.Menu(ds.MenuProps{Label: "Project actions", Items: []ds.DialogAction{
				{Title: "Rename", Icon: ui.IconEdit},
				{Title: "Move to", Children: []ds.DialogAction{
					{Title: "Archive"}, {Title: "Trash", Role: ui.ActionDestructive},
				}},
			}}, ds.Button(ds.ButtonProps{Text: "Project actions", Variant: ds.Outline})),
			`ds.Menu(ds.MenuProps{Label: "Actions", Items: []ds.DialogAction{
    {Title: "Rename", Icon: ui.IconEdit},
    {Title: "Move to", Children: []ds.DialogAction{{Title: "Archive"}}},
}}, ds.Button(ds.ButtonProps{Text: "Actions", Variant: ds.Outline}))`, true
	case "popover":
		return ds.Popover(ds.PopoverProps{
				Label: "Filters",
				Content: func() ui.Node {
					return ui.Column(ui.Style{Gap: ds.Space2},
						ds.Text("Only show items updated in the last week."),
						ds.Button(ds.ButtonProps{Text: "Apply", Variant: ds.Primary, Size: ds.Small}),
					)
				},
			}), `ds.Popover(ds.PopoverProps{
    Label: "Filters",
    Content: func() ui.Node {
        return ds.Text("Only show recent items")
    },
})`, true
	case "toast":
		return ds.Button(ds.ButtonProps{Text: "Show toast", Variant: ds.Outline, OnPress: func() {
				ds.Toast("Saved", &ds.DialogAction{Title: "Undo"})
			}}),
			`ds.Toast("Saved", &ds.DialogAction{Title: "Undo"})`, true
	case "card":
		return ds.Card(ds.CardProps{}, ui.Column(ui.Style{Gap: ds.Space3},
				ds.CardHeader("Deploy", "Your app is live in three regions."),
				ds.Muted("The last deploy was 4 minutes ago."),
				ds.CardFooter(
					ds.Button(ds.ButtonProps{Text: "Logs", Variant: ds.Outline, Size: ds.Small}),
					ds.Button(ds.ButtonProps{Text: "Redeploy", Variant: ds.Primary, Size: ds.Small}),
				),
			)),
			`ds.Card(ds.CardProps{}, ui.Column(ui.Style{Gap: ds.Space3},
    ds.CardHeader("Deploy", "Your app is live."),
    ds.CardFooter(ds.Button(ds.ButtonProps{Text: "Redeploy"})),
))`, true
	case "badge":
		return row(
			ds.Badge(ds.BadgeProps{Text: "Primary", Variant: ds.Primary}),
			ds.Badge(ds.BadgeProps{Text: "Secondary", Variant: ds.Secondary}),
			ds.Badge(ds.BadgeProps{Text: "Outline", Variant: ds.Outline}),
			ds.Badge(ds.BadgeProps{Text: "Live", Variant: ds.Outline, Dot: true}),
		), `ds.Badge(ds.BadgeProps{Text: "Live", Variant: ds.Outline, Dot: true})`, true
	case "avatar":
		return row(
			ds.Avatar(ds.AvatarProps{Name: "Ada Lovelace"}),
			ds.Avatar(ds.AvatarProps{Name: "Grace Hopper", Size: 32}),
		), `ds.Avatar(ds.AvatarProps{Name: "Ada Lovelace"})`, true
	case "table":
		return ds.Table(ds.TableProps{
				Columns: []ds.TableColumn{{Title: "Name", Width: 2}, {Title: "Role"}},
				Rows:    [][]string{{"Ada Lovelace", "Engineer"}, {"Grace Hopper", "Admiral"}},
				Striped: true,
			}), `ds.Table(ds.TableProps{
    Columns: []ds.TableColumn{{Title: "Name", Width: 2}, {Title: "Role"}},
    Rows: [][]string{{"Ada Lovelace", "Engineer"}},
})`, true
	case "data-table":
		sortColumn, descending := signal.New(-1), signal.New(false)
		return ds.DataTable(ds.DataTableProps{
				Columns: []ds.DataTableColumn{
					{Title: "Name", Width: 2, Sortable: true},
					{Title: "Role", Sortable: true},
				},
				Rows: [][]string{{"Ada Lovelace", "Engineer"}, {"Grace Hopper", "Admiral"}},
				Sort: sortColumn, Descending: descending, Striped: true,
			}), `sortColumn, descending := signal.New(-1), signal.New(false)
ds.DataTable(ds.DataTableProps{
    Columns: []ds.DataTableColumn{{Title: "Name", Sortable: true}},
    Rows: [][]string{{"Ada Lovelace", "Engineer"}},
    Sort: sortColumn, Descending: descending,
})`, true
	case "list":
		notify := signal.New(true)
		return ui.View(ui.ViewProps{Style: ui.Style{MaxWidth: ui.Pt(520)}},
				ds.List(ds.ListProps{Divided: true},
					ds.ListItem(ds.ListItemProps{Title: "Profile", Description: "Name, photo and email",
						Icon: ui.IconPerson, Chevron: true, OnPress: func() {}}),
					ds.ListItem(ds.ListItemProps{Title: "Notifications", Description: "Mentions and weekly summaries",
						Icon: ui.IconBell, Action: ui.Toggle(ui.ToggleProps{
							Value: notify,
							Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
								return ui.Style{Color: c.Primary}
							}),
						})}),
					ds.ListItem(ds.ListItemProps{Title: "Plan", Detail: "Pro", Icon: ui.IconStar,
						Chevron: true, OnPress: func() {}}),
				),
			), `ds.List(ds.ListProps{Divided: true},
    ds.ListItem(ds.ListItemProps{Title: "Profile", Icon: ui.IconPerson,
        Chevron: true, OnPress: openProfile}),
    ds.ListItem(ds.ListItemProps{Title: "Plan", Detail: "Pro",
        Chevron: true, OnPress: openPlan}),
)`, true
	case "stat":
		return ui.Row(ui.Style{Gap: ds.Space3, Wrap: true, Align: ui.AlignStretch},
				ui.View(ui.ViewProps{Style: ui.Style{Basis: ui.Pt(190), Grow: 1}},
					ds.Stat(ds.StatProps{Label: "Monthly revenue", Value: "$48,210", Delta: "+12.4%",
						Trend: ds.TrendUp, Description: "Compared with last month"})),
				ui.View(ui.ViewProps{Style: ui.Style{Basis: ui.Pt(190), Grow: 1}},
					ds.Stat(ds.StatProps{Label: "Churn", Value: "1.8%", Delta: "-0.4%",
						Trend: ds.TrendDown, Description: "Fewer cancellations"})),
			), `ds.Stat(ds.StatProps{
    Label: "Monthly revenue", Value: "$48,210",
    Delta: "+12.4%", Trend: ds.TrendUp,
    Description: "Compared with last month",
})`, true
	case "empty-state":
		return ds.EmptyState(ds.EmptyStateProps{
				Icon: ui.IconFolder, Title: "No projects yet",
				Description: "Create a project to deploy your first app.",
				Action: func() ui.Node {
					return ds.Button(ds.ButtonProps{Text: "New project", Icon: ui.IconPlus, Variant: ds.Primary})
				},
			}), `ds.EmptyState(ds.EmptyStateProps{
    Icon: ui.IconFolder, Title: "No projects yet",
    Action: func() ui.Node { return ds.Button(ds.ButtonProps{Text: "New project"}) },
})`, true
	case "alert":
		return ds.Alert(ds.AlertProps{Title: "Heads up", Variant: ds.Primary,
				Description: "M1 through M3 are implemented in this library."}),
			`ds.Alert(ds.AlertProps{
    Title: "Heads up", Variant: ds.Primary,
    Description: "This is an inline message.",
})`, true
	case "progress":
		return col(
			ds.Progress(ds.ProgressProps{Value: ui.Val(0.68), Label: "Storage", ShowValue: true}),
			ds.Progress(ds.ProgressProps{Value: ui.Val(float64(-1)), Label: "Indexing"}),
		), `ds.Progress(ds.ProgressProps{Value: ui.Val(0.68), Label: "Storage", ShowValue: true})`, true
	case "progress-ring":
		return row(
				ds.ProgressRing(ds.ProgressRingProps{Value: ui.Val(0.72), Label: "Rollout", ShowValue: true}),
				ds.ProgressRing(ds.ProgressRingProps{Size: 32, Label: "Syncing"}),
			), `ds.ProgressRing(ds.ProgressRingProps{
    Value: ui.Val(0.72), Label: "Rollout", ShowValue: true,
})`, true
	case "spinner":
		return row(
			ds.Spinner(ds.SpinnerProps{Label: "Syncing"}),
			ds.Spinner(ds.SpinnerProps{Size: 16}),
		), `ds.Spinner(ds.SpinnerProps{Label: "Syncing"})`, true
	case "skeleton":
		return col(
			ds.Skeleton(ui.Pct(100), 8),
			ds.Skeleton(ui.Pct(64), 8),
		), `ds.Skeleton(ui.Pct(100), 8)`, true
	case "tooltip":
		return ds.Tooltip(ds.TooltipProps{Text: "Copy to clipboard"},
				ds.Button(ds.ButtonProps{Text: "Hover me", Variant: ds.Outline})),
			`ds.Tooltip(ds.TooltipProps{Text: "Copy to clipboard"},
    ds.Button(ds.ButtonProps{Text: "Copy"}))`, true
	case "heading":
		return col(
				ds.Heading("Page heading"),
				ds.Title("Section title"),
			), `ds.Heading("Page heading")
ds.Title("Section title")`, true
	case "text":
		return col(
				ds.Text("Body text: the default reading size."),
				ds.Muted("Muted text: descriptions, timestamps and hints."),
			), `ds.Text("Body text")
ds.Muted("Secondary text")`, true
	case "kbd":
		return row(
			ds.Kbd("⌘", "K"),
			ds.Muted("opens the command palette"),
			ds.Kbd("⌘", "⇧", "P"),
		), `ds.Kbd("⌘", "K")`, true
	}
	return ds.Alert(ds.AlertProps{Title: "Preview coming soon", Variant: ds.Outline,
		Description: "This component is installable now; its live preview is on the way."}), "", false
}
