package ds

import (
	"github.com/egoist/asuka/ui"
)

// BadgeProps configures a Badge.
type BadgeProps struct {
	Text    string
	Variant Variant
	// Dot shows a colored dot before the text, for a status.
	Dot bool
}

// Badge is a small status label: a pill of caption text.
func Badge(p BadgeProps) ui.Node {
	var kids []ui.Node
	if p.Dot {
		kids = append(kids, ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
			_, fg, _ := variantColors(c, p.Variant)
			return ui.Style{Width: ui.Pt(6), Height: ui.Pt(6), CornerRadius: 3, Background: fg}
		})}))
	}
	kids = append(kids, ui.StaticText(p.Text, themed(func(t Theme, c Palette) ui.Style {
		_, fg, _ := variantColors(c, p.Variant)
		return ui.Style{Font: t.Type.Caption, Color: fg, Lines: 1}
	})))
	return ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
		bg, _, border := variantColors(c, p.Variant)
		s := ui.Style{
			Direction:    ui.RowDir,
			Align:        ui.AlignCenter,
			Gap:          Space1,
			Padding:      ui.Symmetric(2, Space2),
			CornerRadius: t.Radii.SM + 2,
			Background:   bg,
		}
		if border != nil {
			s.BorderWidth = 1
			s.BorderColor = *border
		}
		return s
	})}, kids...)
}

// Separator is a hairline rule, horizontal unless Vertical.
func Separator(vertical bool) ui.Node {
	s := ui.Style{Height: ui.Pt(1), Width: ui.Pct(100)}
	if vertical {
		s = ui.Style{Width: ui.Pt(1), Height: ui.Pct(100)}
	}
	return ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
		s.Background = c.Border
		return s
	})})
}

// Skeleton is a placeholder block for content that is loading.
func Skeleton(width ui.Dim, height float32) ui.Node {
	return ui.View(ui.ViewProps{Style: themed(func(t Theme, c Palette) ui.Style {
		return ui.Style{Width: width, Height: ui.Pt(height), CornerRadius: t.Radii.SM, Background: c.Muted}
	})})
}
