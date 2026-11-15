package ds

import (
	"github.com/egoist/asuka/fonts"
	"github.com/egoist/asuka/ui"
)

// Variant is a component's emphasis.
type Variant uint8

const (
	Primary Variant = iota
	Secondary
	Outline
	Ghost
	Destructive
	Link
)

// Size is a control's size.
type Size uint8

const (
	Small Size = iota
	Medium
	Large
	IconOnly
)

func metrics(s Size) (height, padX float32) {
	switch s {
	case Small:
		return 32, Space3
	case Large:
		return 40, Space5
	case IconOnly:
		return 36, 0
	default:
		return 36, Space4
	}
}

func fontFor(t TypeScale, s Size) fonts.Font {
	switch s {
	case Small:
		return t.Small
	case Large, Medium:
		return t.Subheading
	default:
		return t.Body
	}
}

func variantColors(c Palette, v Variant) (bg, fg ui.Color, border *ui.Color) {
	switch v {
	case Primary:
		return c.Primary, c.PrimaryForeground, nil
	case Secondary:
		return c.Secondary, c.SecondaryForeground, nil
	case Outline:
		b := c.Border
		return 0, c.Foreground, &b
	case Ghost:
		return 0, c.Foreground, nil
	case Destructive:
		return c.Destructive, c.OnDestructive, nil
	case Link:
		return 0, c.Primary, nil
	}
	return c.Primary, c.PrimaryForeground, nil
}

func hoverColors(c Palette, v Variant, bg ui.Color) ui.Color {
	switch v {
	case Primary, Destructive:
		return bg.Opacity(0.88)
	case Outline, Ghost, Link:
		return c.Accent
	}
	return c.Accent
}
