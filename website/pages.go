package main

import (
	ds "github.com/egoist/asuka-ds/ds"
	"github.com/egoist/asuka-ds/internal/catalog"
	"github.com/egoist/asuka/ui"
)

func pageHeader(title, subtitle string) ui.Node {
	return ui.Column(ui.Style{Gap: ds.Space2},
		ds.Heading(title),
		ds.Text(subtitle),
	)
}

func section(title string, children ...ui.Node) ui.Node {
	return ui.Column(ui.Style{Gap: ds.Space3},
		ds.Title(title),
		ui.Column(ui.Style{Gap: ds.Space3}, children...),
	)
}

func introPage(d *docs) ui.Node {
	return ui.Column(ui.Style{Gap: ds.Space6},
		ui.Column(ui.Style{Gap: ds.Space3},
			ui.Row(ui.Style{Gap: ds.Space2, Align: ui.AlignCenter, Wrap: true},
				ds.Badge(ds.BadgeProps{Text: "v0.1.0", Variant: ds.Primary}),
				ds.Badge(ds.BadgeProps{Text: "Go 1.26+", Variant: ds.Outline}),
				ds.Badge(ds.BadgeProps{Text: "Asuka 0.1.1", Variant: ds.Outline}),
			),
			ds.Heading("Components that feel native on every platform"),
			ds.Text("Asuka DS is a source-owned component library for Asuka: shadcn's token structure and composable surfaces, built on platform controls wherever the platform does it better."),
			ui.Row(ui.Style{Gap: ds.Space2, Wrap: true},
				ds.Button(ds.ButtonProps{Text: "Get started", Variant: ds.Primary,
					OnPress: func() { d.page.Set("install") }}),
				ds.Button(ds.ButtonProps{Text: "Browse components", Variant: ds.Outline,
					OnPress: func() { d.page.Set("c:button") }}),
				link("GitHub", "https://github.com/egoist/asuka-ds"),
			),
		),
		ui.Row(ui.Style{Gap: ds.Space3, Wrap: true, Align: ui.AlignStretch},
			feature("Native by default", "Buttons, inputs, pickers and overlays use the platform's own controls, so IME, selection, focus and screen readers behave."),
			feature("Source you own", "The CLI copies Go source into your project. Change the code, keep the API, or vendor the whole design system."),
			feature("Tokens are typed", "Palette, radii, type and motion are structs. Themes switch at runtime through signals, without remounting the page."),
		),
		section("One API, two worlds",
			ds.Text("Use drawn components for brand consistency and native components on the paths where platform behavior matters more."),
			code(`ds.Button(ds.ButtonProps{Text: "Save", OnPress: save})        // drawn
ds.NativeButton(ds.ButtonProps{Text: "Pay", OnPress: pay})  // platform
ds.Input(ds.InputProps{Value: email, Label: "Email"})       // native field, DS shell`),
		),
		section("A quick look",
			ui.Row(ui.Style{Gap: ds.Space3, Wrap: true, Align: ui.AlignStart},
				ds.Card(ds.CardProps{}, ui.Column(ui.Style{Gap: ds.Space3},
					ds.CardHeader("Deploy", "Your app is live in three regions."),
					ui.Row(ui.Style{Gap: ds.Space2},
						ds.Badge(ds.BadgeProps{Text: "Live", Variant: ds.Outline, Dot: true}),
						ds.Badge(ds.BadgeProps{Text: "v1.4", Variant: ds.Secondary}),
					),
					ds.Progress(ds.ProgressProps{Value: ui.Val(0.72), Label: "Rollout", ShowValue: true}),
				)),
				ds.Card(ds.CardProps{}, ui.Column(ui.Style{Gap: ds.Space3},
					ds.CardHeader("Search", "A command palette, not a menu."),
					ds.Command(ds.CommandProps{Items: []ds.CommandItem{
						{Title: "New project", Icon: ui.IconPlus, Keywords: "create"},
						{Title: "Open settings", Icon: ui.IconSettings, Keywords: "preferences"},
					}}),
				)),
			),
		),
		section("Status",
			ds.Table(ds.TableProps{
				Columns: []ds.TableColumn{{Title: "Milestone", Width: 2}, {Title: "State"}, {Title: "What shipped"}},
				Rows: [][]string{
					{"M0–M3", "Done", "Tokens, themes, primitives, 43 components"},
					{"M4", "Done", "Registry, CLI, docs site"},
					{"M5", "In progress", "Release pipeline done; 4 of 6 platforms captured"},
				},
				Striped: true,
			}),
		),
	)
}

func feature(title, body string) ui.Node {
	return ui.View(ui.ViewProps{Style: ui.Style{Grow: 1, Basis: ui.Pt(220)}},
		ds.Card(ds.CardProps{}, ui.Column(ui.Style{Gap: ds.Space2},
			ds.CardTitle(ui.Val(title)),
			ds.Muted(body),
		)))
}

func installPage() ui.Node {
	return ui.Column(ui.Style{Gap: ds.Space6},
		pageHeader("Installation", "Install the CLI, start a design-system package, then add only the components you need."),
		section("1. Install the CLI",
			code(`go install github.com/egoist/asuka-ds/cmd/asuka-ds@latest`),
			ds.Muted("Or use the module directly: go get github.com/egoist/asuka-ds/ds"),
		),
		section("2. Create a package",
			code(`asuka-ds init --dir internal/ds --package ds`),
			ds.Text("init writes the base layer: theme, text, Pressable, Field and the shared variant vocabulary. No components are added yet."),
		),
		section("3. Add components",
			code(`asuka-ds add button input card switch
asuka-ds add date-range-picker --dir internal/ds`),
			ds.Text("Dependencies are resolved for you: add switch also installs overlay.go, and add command installs button.go."),
		),
		section("4. Use them",
			code(`package main

import (
    ds "example.com/app/internal/ds"
    "github.com/egoist/asuka/app"
    "github.com/egoist/asuka/signal"
    "github.com/egoist/asuka/ui"
)

func main() { app.Run(App) }

func App() ui.Node {
    email := signal.New("")
    return ds.Screen(
        ds.Heading("Settings"),
        ds.Input(ds.InputProps{Value: email, Label: "Email"}),
    )
}`),
		),
		section("Rules that keep it yours",
			ds.Alert(ds.AlertProps{Title: "Source-owned, not locked in", Variant: ds.Primary,
				Description: "The CLI never overwrites edited files unless you pass --force. The copied code is ordinary Go with no generated runtime."}),
			ds.Text("You can also depend on the module directly, which is the better choice when you want upgrades without vendoring."),
		),
	)
}

func themesPage() ui.Node {
	return ui.Column(ui.Style{Gap: ds.Space6},
		pageHeader("Theming", "Two built-in themes and one typed token layer. Switch them at runtime with signals."),
		section("Native or Neutral",
			ds.Table(ds.TableProps{
				Columns: []ds.TableColumn{{Title: ""}, {Title: "Native"}, {Title: "Neutral"}},
				Rows: [][]string{
					{"Colors", "Platform semantic colors", "shadcn zinc, light and dark"},
					{"Type", "Platform text styles", "Fixed 14pt body"},
					{"Dark mode", "Automatic", "Explicit palettes"},
					{"Best for", "Native feel", "One brand on six platforms"},
				},
				Striped: true,
			}),
			ds.Text("Use the Native/Neutral and System/Light/Dark tabs in the header; the whole site is the preview."),
		),
		section("Switch at runtime",
			code(`ds.SetTheme(ds.Native())   // or ds.Neutral()
ds.SetAppearance(ui.ColorSchemeDark)`),
			ds.Text("SetAppearance updates the host and the design system's own appearance signal, so drawn components change on the click rather than after the host's next redraw."),
		),
		section("The palette",
			palettePreview("Native", ds.Native().PaletteFor(false)),
			palettePreview("Neutral", ds.Neutral().PaletteFor(false)),
		),
		section("Escape hatch",
			code(`ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
    return ui.Style{
        Background:   c.Card,
        Color:        c.Foreground,
        CornerRadius: t.Radii.LG,
    }
})`),
			ds.Text("Styles written with Themed follow SetTheme and SetAppearance, so app-level chrome participates in live theming."),
		),
	)
}

func palettePreview(title string, p ds.Palette) ui.Node {
	roles := []struct {
		name string
		col  ui.Color
	}{
		{"Background", p.Background},
		{"Foreground", p.Foreground},
		{"Primary", p.Primary},
		{"Secondary", p.Secondary},
		{"Muted", p.Muted},
		{"Accent", p.Accent},
		{"Destructive", p.Destructive},
		{"Border", p.Border},
	}
	swatches := make([]ui.Node, 0, len(roles))
	for _, role := range roles {
		swatches = append(swatches, ui.Column(ui.Style{Gap: ds.Space1},
			ui.View(ui.ViewProps{Style: ui.Style{
				Height: ui.Pt(42), CornerRadius: 6, Background: role.col,
				BorderWidth: 1, BorderColor: p.Border,
			}}),
			ds.Muted(role.name),
		))
	}
	return ui.Column(ui.Style{Gap: ds.Space2},
		ds.Title(title),
		ui.Row(ui.Style{Gap: ds.Space2, Wrap: true}, swatches...),
	)
}

func cliPage() ui.Node {
	return ui.Column(ui.Style{Gap: ds.Space6},
		pageHeader("CLI & registry", "The registry is a typed catalog plus embedded source. The CLI is a thin, offline installer."),
		section("Commands",
			code(`asuka-ds init
asuka-ds add button input card
asuka-ds add switch --dir internal/ds
asuka-ds list
asuka-ds version`),
		),
		section("What add does",
			ui.Column(ui.Style{Gap: ds.Space2},
				ds.Text("1. Resolves the component's kit and every kit it requires."),
				ds.Text("2. Copies unique source files, rewriting the package clause if you asked for another name."),
				ds.Text("3. Leaves existing files alone unless --force is given."),
			),
			ds.Alert(ds.AlertProps{Title: "No network at install time", Variant: ds.Primary,
				Description: "Component sources are embedded in the CLI binary. Installing works offline and does not need a checkout."}),
		),
		section("Catalog",
			ds.Text("The public catalog is internal/catalog/components.json. It names the flat files in ds/ by bare filename, and is the source of truth for the docs, the CLI and dependency resolution."),
			code(`{
  "kits": {
    "combobox": { "files": ["combobox.go"], "requires": ["overlay"] }
  },
  "components": {
    "combobox": { "kit": "combobox", "category": "Form", ... }
  }
}`),
		),
		section("Contributing",
			code(`go generate ./internal/registry   # refresh embedded source
go test ./...`),
			ds.Text("A drift test fails if the embedded sources no longer match the files in the repository."),
		),
	)
}

func qualityPage() ui.Node {
	return ui.Column(ui.Style{Gap: ds.Space6},
		pageHeader("Quality & status", "What is tested today, and which platform gaps are known."),
		section("Checks",
			ds.Table(ds.TableProps{
				Columns: []ds.TableColumn{{Title: "Check", Width: 2}, {Title: "Command"}, {Title: "State"}},
				Rows: [][]string{
					{"Unit & component tests", "go test ./...", "76 tests, 9 packages"},
					{"Static analysis", "go vet ./...", "Pass"},
					{"Asuka checks", "asuka vet ./...", "Pass"},
					{"Docs site", "asuka build --platform web website", "Pass"},
				},
				Striped: true,
			}),
		),
		section("Known boundaries",
			ds.Alert(ds.AlertProps{Title: "Native controls keep their chrome", Variant: ds.Destructive,
				Description: "A native button, switch, slider or picker can be tinted and typed, but not redrawn. Drawn components keep the brand but cannot use Liquid Glass or Material ripple."}),
			ds.Text("Drawn controls have no native Disabled, selected or expanded trait. Pressable expresses them through behavior and Hint; critical paths should use NativeButton or another platform control."),
		),
		section("Roadmap",
			ds.Table(ds.TableProps{
				Columns: []ds.TableColumn{{Title: "Phase", Width: 1}, {Title: "State"}, {Title: "Scope"}},
				Rows: [][]string{
					{"M0–M3", "Done", "Tokens, themes, primitives, components, demos"},
					{"M4", "Done", "Registry, CLI, docs site, module identity"},
					{"M5", "In progress", "Release pipeline done; 4 of 6 platforms captured"},
				},
				Striped: true,
			}),
		),
		section("API stability",
			ds.Text("The token layer is frozen as a typed API. Component props are pre-1.0: additive changes ship in minor releases, breaking changes get a migration note and a minor bump until 1.0."),
		),
	)
}

func platformsPage() ui.Node {
	return ui.Column(ui.Style{Gap: ds.Space6},
		pageHeader("Platforms", "What was built and seen on each host, and what is still a preview."),
		section("Verified on 2026-11-15",
			ds.Table(ds.TableProps{
				Columns: []ds.TableColumn{{Title: "Platform", Width: 2}, {Title: "Toolchain"}, {Title: "Status"}},
				Rows: [][]string{
					{"macOS 27.2.1", "Xcode 27.2", "Built, screenshots"},
					{"iOS 26.5", "iPhone 17 Pro simulator", "Installed, screenshots"},
					{"Android 16", "Pixel 9a emulator", "Installed, screenshots"},
					{"Web", "asuka build --platform web", "Built, screenshots"},
					{"Windows", "GOOS=windows go build", "Cross-build only"},
					{"Linux", "GOOS=linux go build", "Cross-build only"},
				},
				Striped: true,
			}),
			ds.Muted("Windows and Linux compile in CI but have no recorded run yet."),
		),
		section("What the device run caught",
			ds.Alert(ds.AlertProps{Title: "iOS accent backgrounds", Variant: ds.Primary,
				Description: "A container's Color is its tint on iOS, so a semantic accent background resolved to the label color. Drawn Button and Toggle now color their label and icon, not the container; render tests pin the rule."}),
			ds.Text("The same run confirmed semantic colors: iOS blue, Android Material dynamic color, and macOS accent, in both Native and Neutral themes."),
		),
		section("Known gaps",
			ds.Table(ds.TableProps{
				Columns: []ds.TableColumn{{Title: "Platform", Width: 1}, {Title: "Gap"}},
				Rows: [][]string{
					{"iOS", "Platform controls keep their chrome; container Color is a tint"},
					{"Android", "Dynamic color changes the accent; use Neutral for a fixed brand"},
					{"Web", "Native pickers and alerts become browser controls"},
					{"Windows / Linux", "No host run yet: focus, IME and screenshots pending"},
				},
				Striped: true,
			}),
		),
		section("Reproduce",
			code(`asuka build --platform mac demo
asuka build --platform ios --simulator demo
asuka dev --platform android --once -d emulator-5554 demo
asuka build --platform web website`),
			ds.Muted("Screenshots live in artwork/screenshots and are named <platform>-<theme>.png."),
		),
	)
}

func componentPage(d *docs, name string) ui.Node {
	component, ok := d.catalog.Components[name]
	if !ok {
		return ds.Alert(ds.AlertProps{Title: "Unknown component", Description: name, Variant: ds.Destructive})
	}
	example, snippet, hasExample := exampleFor(name)
	kit := d.catalog.Kits[component.Kit]
	var sourceLinks []ui.Node
	for _, file := range kit.Files {
		sourceLinks = append(sourceLinks, link(file, "https://github.com/egoist/asuka-ds/blob/main/"+catalog.SourceDir+"/"+file))
	}
	kids := []ui.Node{
		ui.Column(ui.Style{Gap: ds.Space2},
			ui.Row(ui.Style{Gap: ds.Space2, Align: ui.AlignCenter, Wrap: true},
				ds.Badge(ds.BadgeProps{Text: component.Category, Variant: ds.Secondary}),
				ds.Muted("install: asuka-ds add "+name),
			),
			ds.Heading(component.Title),
			ds.Text(component.Description),
		),
	}
	if hasExample {
		kids = append(kids, section("Preview", ui.View(ui.ViewProps{Style: ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
			return ui.Style{Padding: ui.All(ds.Space5), Background: c.Background,
				BorderWidth: 1, BorderColor: c.Border, CornerRadius: t.Radii.LG}
		})}, example)))
	}
	if snippet != "" {
		kids = append(kids, section("Usage", code(snippet)))
	}
	kids = append(kids, section("Behavior",
		ds.Text(behaviorFor(component)),
		ds.Text("Install with asuka-ds add "+name+" and edit the copied source as your own."),
	))
	if len(sourceLinks) > 0 {
		kids = append(kids, section("Source", ui.Row(ui.Style{Gap: ds.Space2, Wrap: true}, sourceLinks...)))
	}
	return ui.Column(ui.Style{Gap: ds.Space6}, kids...)
}

func behaviorFor(component catalog.Component) string {
	switch component.Kit {
	case "controls", "overlay", "input", "datepicker", "daterange", "segmented", "numberfield",
		"popover", "colorpicker", "progressring":
		return "This component wraps a platform primitive: the platform owns editing, dismissal, keyboard focus and accessibility, while the design system owns the label, chrome and tokens."
	default:
		return "This component is drawn from Views and Pressables: the same look on every platform. Disabled, selected and expanded states are expressed through Pressable behavior and its accessibility Hint."
	}
}
