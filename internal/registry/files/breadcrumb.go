package ds

import "github.com/egoist/asuka/ui"

// BreadcrumbItem is one step of a Breadcrumb. The last item is the current
// page and is not pressable.
type BreadcrumbItem struct {
	Title string
	Href  ui.Href
	// OnPress handles a crumb in the app instead of navigating to Href.
	OnPress func()
}

// BreadcrumbProps configures a Breadcrumb.
type BreadcrumbProps struct {
	Items []BreadcrumbItem
	// Label is the trail's accessible name ("Breadcrumb" by default).
	Label string
}

// Breadcrumb is the path to the current page: a row of links separated by
// chevrons, ending in plain text. Links are Pressables, so their Hint and
// role work on every host; the current page is a text view and is not
// pressable.
func Breadcrumb(p BreadcrumbProps) ui.Node {
	label := p.Label
	if label == "" {
		label = "Breadcrumb"
	}
	kids := make([]ui.Node, 0, len(p.Items)*2)
	for i, item := range p.Items {
		index, crumb := i, item
		if index > 0 {
			kids = append(kids, ui.Symbol(ui.IconForward, themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Font: t.Type.Caption, Color: c.MutedForeground}
			})))
		}
		if index == len(p.Items)-1 {
			kids = append(kids, ui.View(ui.ViewProps{
				Label: ui.Val(crumb.Title),
				Role:  ui.RoleText,
			}, ui.StaticText(crumb.Title, themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Font: t.Type.Small, Color: c.Foreground}
			}))))
			continue
		}
		role := ui.RoleButton
		if crumb.Href != "" {
			role = ui.RoleLink
		}
		kids = append(kids, Pressable(PressableProps{
			Label:   crumb.Title,
			Role:    role,
			Href:    crumb.Href,
			OnPress: crumb.OnPress,
			Style: themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{
					Padding: ui.Symmetric(2, 0),
					Color:   c.MutedForeground,
					Hovered: &ui.Style{Color: c.Foreground},
				}
			}),
		}, ui.StaticText(crumb.Title, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Small, Color: c.MutedForeground}
		}))))
	}
	return ui.View(ui.ViewProps{Label: ui.Val(label)},
		ui.Row(ui.Style{Gap: Space2, Align: ui.AlignCenter, Wrap: true}, kids...))
}
