package main

import (
	ds "github.com/egoist/asuka-ds"

	"github.com/egoist/asuka/app"
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

func main() { app.Run(App) }

func App() ui.Node {
	name := signal.New("Ada Lovelace")
	email := signal.New("")
	bio := signal.New("")
	workspace := signal.New("")
	notify := signal.New(true)
	tab := signal.New(1)

	return ds.Screen(
		ds.Heading("Form controls"),
		ds.Muted("Native text fields and switches inside the design system's chrome."),
		ds.Card(ds.CardProps{},
			ui.Column(ui.Style{Gap: ds.Space4},
				ds.CardHeader("Profile", "How the team sees you."),
				ds.Input(ds.InputProps{Value: name, Label: "Name", Placeholder: "Ada Lovelace",
					Description: "Shown on your profile."}),
				ds.Input(ds.InputProps{Value: email, Label: "Email", Placeholder: "ada@example.com",
					Icon: ui.IconMail, Error: "That address is already taken."}),
				ds.Input(ds.InputProps{Value: bio, Label: "Bio", Multiline: true, Lines: 3,
					Placeholder: "A few words about you"}),
				ds.Input(ds.InputProps{Value: workspace, Label: "Workspace", Placeholder: "acme", Disabled: true,
					Description: "Ask an admin to change it."}),
				ui.View(ui.ViewProps{Label: ui.Val("Email notifications")},
					ui.Row(ui.Style{Gap: ds.Space3, Align: ui.AlignCenter},
						ui.Toggle(ui.ToggleProps{Value: notify}),
						ds.Text("Email notifications"),
					),
				),
				ds.Tabs(ds.TabsProps{Items: []string{"Overview", "Activity", "Settings"}, Selected: tab}),
				ds.CardFooter(
					ds.Button(ds.ButtonProps{Text: "Cancel", Variant: ds.Ghost, OnPress: func() {}}),
					ds.Button(ds.ButtonProps{Text: "Save changes", Variant: ds.Primary, OnPress: func() {}}),
				),
			),
		),
	)
}
