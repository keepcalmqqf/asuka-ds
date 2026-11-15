package ds

import (
	"fmt"
	"strconv"

	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// PaginationProps configures a Pagination.
type PaginationProps struct {
	// Page is the current page, 1-based; the user's choice writes it.
	Page *signal.Signal[int]
	// Total is the number of pages.
	Total int
	// Siblings is how many pages to show on each side of the current one
	// before collapsing the rest into an ellipsis (0: 1).
	Siblings int
	// Label is the group's accessible name ("Pagination" by default).
	Label string
}

// Pagination is a row of page buttons with previous and next arrows. The
// current page carries the Selected hint, the arrows are unavailable at the
// ends, and the page list rebuilds as Page changes, so the ellipses follow
// the current page.
func Pagination(p PaginationProps) ui.Node {
	total := p.Total
	if total < 1 {
		total = 1
	}
	siblings := p.Siblings
	if siblings <= 0 {
		siblings = 1
	}
	label := p.Label
	if label == "" {
		label = "Pagination"
	}
	current := func() int {
		v := 1
		if p.Page != nil {
			v = p.Page.Get()
		}
		if v < 1 {
			return 1
		}
		if v > total {
			return total
		}
		return v
	}
	goTo := func(n int) {
		if n < 1 {
			n = 1
		}
		if n > total {
			n = total
		}
		if p.Page != nil {
			p.Page.Set(n)
		}
	}
	square := func(hover bool) ui.Styler {
		return themed(func(t Theme, c Palette) ui.Style {
			s := ui.Style{
				Direction: ui.RowDir, Align: ui.AlignCenter, Justify: ui.JustifyCenter,
				Width: ui.Pt(32), Height: ui.Pt(32), CornerRadius: t.Radii.MD,
			}
			if hover {
				s.Hovered = &ui.Style{Background: c.Muted}
			}
			return s
		})
	}
	arrow := func(icon ui.Icon, name string, disabled bool, delta int) ui.Node {
		return Pressable(PressableProps{
			Label:    name,
			Disabled: ui.Val(disabled),
			OnPress:  func() { goTo(current() + delta) },
			Style:    square(true),
		}, ui.Symbol(icon, themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{Font: t.Type.Small, Color: c.Foreground}
		})))
	}
	control := ui.Dynamic(func() ui.Node {
		active := current()
		prev := arrow(ui.IconBack, "Previous page", active <= 1, -1)
		next := arrow(ui.IconForward, "Next page", active >= total, 1)
		windows := pageWindow(active, total, siblings)
		kids := make([]ui.Node, 0, len(windows))
		for _, n := range windows {
			if n < 0 {
				kids = append(kids, ui.View(ui.ViewProps{
					Label: ui.Val("More pages"),
					Style: ui.Style{Width: ui.Pt(24), Height: ui.Pt(32),
						Align: ui.AlignCenter, Justify: ui.JustifyCenter},
				}, ui.StaticText("…", themed(func(t Theme, c Palette) ui.Style {
					return ui.Style{Font: t.Type.Small, Color: c.MutedForeground}
				}))))
				continue
			}
			page, selected := n, n == active
			kids = append(kids, Pressable(PressableProps{
				Label:    fmt.Sprintf("Page %d", page),
				Selected: ui.Val(selected),
				OnPress:  func() { goTo(page) },
				Style: themed(func(t Theme, c Palette) ui.Style {
					s := ui.Style{
						Direction: ui.RowDir, Align: ui.AlignCenter, Justify: ui.JustifyCenter,
						Width: ui.Pt(32), Height: ui.Pt(32), CornerRadius: t.Radii.MD,
					}
					if selected {
						s.Background = c.Primary
					} else {
						s.Hovered = &ui.Style{Background: c.Muted}
					}
					return s
				}),
			}, ui.StaticText(strconv.Itoa(page), themed(func(t Theme, c Palette) ui.Style {
				col := c.Foreground
				if selected {
					col = c.PrimaryForeground
				}
				return ui.Style{Font: t.Type.Small, Color: col}
			}))))
		}
		return ui.Row(ui.Style{Gap: Space2, Align: ui.AlignCenter, Wrap: true},
			prev,
			ui.Row(ui.Style{Gap: Space1, Align: ui.AlignCenter}, kids...),
			next,
		)
	})
	return ui.View(ui.ViewProps{Label: ui.Val(label)}, control)
}

// pageWindow lists the pages to show, with -1 for each gap. Short runs show
// every page; longer ones keep the first, the last and the current page's
// neighbors.
func pageWindow(current, total, siblings int) []int {
	if total <= 5+2*siblings {
		out := make([]int, 0, total)
		for i := 1; i <= total; i++ {
			out = append(out, i)
		}
		return out
	}
	keep := map[int]bool{1: true, total: true}
	for i := current - siblings; i <= current+siblings; i++ {
		if i >= 1 && i <= total {
			keep[i] = true
		}
	}
	var out []int
	prev := 0
	for i := 1; i <= total; i++ {
		if !keep[i] {
			continue
		}
		if prev > 0 && i-prev > 1 {
			out = append(out, -1)
		}
		out = append(out, i)
		prev = i
	}
	return out
}
