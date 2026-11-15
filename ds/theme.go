// Package ds is a shadcn-style design system for Asuka: a token layer, and
// components built from Asuka's own primitives so apps can copy, own and
// change them.
package ds

import (
	"fmt"
	"math"

	"github.com/egoist/asuka/colors"
	"github.com/egoist/asuka/fonts"
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// Spacing scale, in points. 1 is a hairline gap, 6 is a page gutter.
const (
	Space1 float32 = 4
	Space2 float32 = 8
	Space3 float32 = 12
	Space4 float32 = 16
	Space5 float32 = 24
	Space6 float32 = 32
	Space7 float32 = 40
	Space8 float32 = 48
)

// Elevation is a surface's depth. Asuka's semantic backgrounds collapse into
// the window background on desktop hosts, so a design system needs its own
// elevation tokens: an overlay of the foreground color, at a known alpha,
// over a known base. Shadow is the same idea cast onto what's behind.
type Elevation struct {
	// Overlay is the foreground color drawn on top of Base.
	Overlay ui.Color
	// Base is the surface the overlay is composited onto.
	Base ui.Color
	// Alpha is the overlay's opacity, 0 to 1.
	Alpha float32
	// Shadow is the drop shadow under the surface.
	Shadow ui.Shadow
}

// Shadow makes the elevation's shadow.
func (e Elevation) ShadowStyle() ui.Shadow { return e.Shadow }

// Motion is the system's animation vocabulary. Durations are seconds; a
// component reads them from the theme rather than hard-coding them, so an app
// can slow everything down or turn motion off in one place.
type Motion struct {
	Instant float32
	Fast    float32
	Normal  float32
	Slow    float32
	// Curve is the default easing.
	Curve ui.Curve
	// Disabled is set by apps that want no motion at all, on top of the
	// platform's ReduceMotion environment value.
	Disabled bool
}

// MotionFor returns the motion tokens for the given environment. When the
// user asked for reduced motion, every duration is zero and the caller can
// skip building animated views.
func (t Theme) MotionFor(reduce bool) Motion {
	m := t.Motion
	if reduce || m.Disabled {
		m.Instant, m.Fast, m.Normal, m.Slow = 0, 0, 0, 0
	}
	return m
}

// Palette is the semantic color set of the design system: the same roles
// shadcn uses as CSS variables, resolved to Asuka colors. Elevation1..3 are
// the layers above Background; they are explicit tokens, not aliases of a
// semantic color, so a card still has a face on macOS and Windows.
type Palette struct {
	Background          ui.Color
	Foreground          ui.Color
	Card                ui.Color
	CardForeground      ui.Color
	Popover             ui.Color
	PopoverForeground   ui.Color
	Primary             ui.Color
	PrimaryForeground   ui.Color
	Secondary           ui.Color
	SecondaryForeground ui.Color
	Muted               ui.Color
	MutedForeground     ui.Color
	Accent              ui.Color
	AccentForeground    ui.Color
	Destructive         ui.Color
	OnDestructive       ui.Color
	Border              ui.Color
	Input               ui.Color
	Ring                ui.Color
	// Focus is the keyboard focus ring. WCAG requires 3:1 against the
	// adjacent background; Ring may be a softer tint used for decoration.
	Focus ui.Color
	// Elevation1 is a subtle raised surface (a track, a hovered row), 2 is
	// a card or a tab's selected item, 3 is a popover, a menu or a dialog.
	Elevation1 Elevation
	Elevation2 Elevation
	Elevation3 Elevation
}

// Radii are corner radii, in points.
type Radii struct {
	SM, MD, LG, XL float32
}

// TypeScale is the theme's type: Asuka's semantic text styles by default, so
// text follows the platform and the user's text size; a fixed scale when the
// look matters more than the platform.
type TypeScale struct {
	Display, Title, Heading, Subheading, Body, Small, Caption fonts.Font
}

// Theme is a set of tokens. Semantic themes follow the platform's colors and
// dark mode on their own; explicit themes carry one palette per appearance.
type Theme struct {
	Name  string
	Light Palette
	Dark  Palette
	// Semantic uses the platform's own colors instead of Light and Dark:
	// the app looks native and follows the system's appearance,
	// accessibility contrast and (on Android) dynamic color.
	Semantic bool
	Radii    Radii
	Type     TypeScale
	Motion   Motion
}

// PaletteFor picks the palette for an appearance.
func (t Theme) PaletteFor(dark bool) Palette {
	if t.Semantic {
		return semanticPalette()
	}
	if dark {
		return t.Dark
	}
	return t.Light
}

// paletteFor picks the palette for an app-level appearance. A semantic
// theme follows the platform while it follows the system; a forced Light
// or Dark uses explicit colors, because a host isn't guaranteed to echo
// SetColorScheme back in its environment (macOS reports the app's new
// appearance only after the window's own redraw, and layer colors are
// fixed CGColors by then).
func (t Theme) paletteFor(scheme ui.ColorScheme, dark bool) Palette {
	if t.Semantic {
		if scheme == ui.ColorSchemeSystem {
			return semanticPalette()
		}
		return nativePalette(dark)
	}
	if dark {
		return t.Dark
	}
	return t.Light
}

// Validate checks a theme's explicit palettes for missing roles and for
// WCAG 2.1 AA contrast (4.5:1 for normal text, 3:1 for large text). Semantic
// palettes resolve on the host, so they are checked for completeness only.
// It returns one message per problem, in role order.
func (t Theme) Validate() []string {
	var out []string
	if t.Name == "" {
		out = append(out, "Theme.Name is empty")
	}
	if t.Radii.MD <= 0 || t.Type.Body == (fonts.Font{}) {
		out = append(out, "Radii.MD and Type.Body must be set")
	}
	if t.Semantic {
		return out
	}
	out = append(out, validatePalette("Light", t.Light)...)
	out = append(out, validatePalette("Dark", t.Dark)...)
	return out
}

var textRoles = []struct {
	name  string
	fg    func(Palette) ui.Color
	bg    func(Palette) ui.Color
	large bool
}{
	{"Foreground/Background", func(p Palette) ui.Color { return p.Foreground }, func(p Palette) ui.Color { return p.Background }, false},
	{"CardForeground/Card", func(p Palette) ui.Color { return p.CardForeground }, func(p Palette) ui.Color { return p.Card }, false},
	{"PopoverForeground/Popover", func(p Palette) ui.Color { return p.PopoverForeground }, func(p Palette) ui.Color { return p.Popover }, false},
	{"PrimaryForeground/Primary", func(p Palette) ui.Color { return p.PrimaryForeground }, func(p Palette) ui.Color { return p.Primary }, false},
	{"SecondaryForeground/Secondary", func(p Palette) ui.Color { return p.SecondaryForeground }, func(p Palette) ui.Color { return p.Secondary }, false},
	{"MutedForeground/Muted", func(p Palette) ui.Color { return p.MutedForeground }, func(p Palette) ui.Color { return p.Muted }, false},
	{"AccentForeground/Accent", func(p Palette) ui.Color { return p.AccentForeground }, func(p Palette) ui.Color { return p.Accent }, false},
	{"OnDestructive/Destructive", func(p Palette) ui.Color { return p.OnDestructive }, func(p Palette) ui.Color { return p.Destructive }, false},
	{"Focus/Background", func(p Palette) ui.Color { return p.Focus }, func(p Palette) ui.Color { return p.Background }, true},
}

func validatePalette(name string, p Palette) []string {
	var out []string
	required := []struct {
		key string
		v   ui.Color
	}{
		{"Background", p.Background}, {"Foreground", p.Foreground},
		{"Card", p.Card}, {"CardForeground", p.CardForeground},
		{"Popover", p.Popover}, {"PopoverForeground", p.PopoverForeground},
		{"Primary", p.Primary}, {"PrimaryForeground", p.PrimaryForeground},
		{"Secondary", p.Secondary}, {"SecondaryForeground", p.SecondaryForeground},
		{"Muted", p.Muted}, {"MutedForeground", p.MutedForeground},
		{"Accent", p.Accent}, {"AccentForeground", p.AccentForeground},
		{"Destructive", p.Destructive}, {"OnDestructive", p.OnDestructive},
		{"Border", p.Border}, {"Input", p.Input}, {"Ring", p.Ring}, {"Focus", p.Focus},
	}
	for _, r := range required {
		if r.v == 0 {
			out = append(out, name+"."+r.key+" is unset (the zero color means the platform default)")
		}
	}
	for _, r := range textRoles {
		fg, bg := r.fg(p), r.bg(p)
		if _, _, _, _, ok := fg.RGBA(); !ok {
			out = append(out, name+"."+r.name+": foreground is a semantic color, cannot measure contrast")
			continue
		}
		if _, _, _, _, ok := bg.RGBA(); !ok {
			out = append(out, name+"."+r.name+": background is a semantic color, cannot measure contrast")
			continue
		}
		ratio := ContrastRatio(fg, bg)
		min := float32(4.5)
		if r.large {
			min = 3.0
		}
		if ratio < min {
			out = append(out, name+"."+r.name+": contrast "+formatRatio(ratio)+":1 is below "+formatRatio(min)+":1")
		}
	}
	return out
}

// ContrastRatio is the WCAG contrast ratio between two opaque Asuka colors.
// Both must be RGBA colors; semantic colors resolve on the host and return 0.
func ContrastRatio(a, b ui.Color) float32 {
	la, ok := relativeLuminance(a)
	if !ok {
		return 0
	}
	lb, ok := relativeLuminance(b)
	if !ok {
		return 0
	}
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

func relativeLuminance(c ui.Color) (float32, bool) {
	r, g, b, a, ok := c.RGBA()
	if !ok {
		return 0, false
	}
	if a == 0 {
		return 0, false
	}
	// Composite over white: a translucent color's drawn value depends on
	// its backdrop, and white is the common case for a token check.
	over := func(v, alpha uint8) float32 {
		return float32(v)*float32(alpha)/255 + 255*(1-float32(alpha)/255)
	}
	return luminance(over(r, a), over(g, a), over(b, a)), true
}

func luminance(r, g, b float32) float32 {
	lin := func(v float32) float32 {
		v /= 255
		if v <= 0.04045 {
			return v / 12.92
		}
		return float32(math.Pow(float64((v+0.055)/1.055), 2.4))
	}
	return 0.2126*lin(r) + 0.7152*lin(g) + 0.0722*lin(b)
}

func formatRatio(v float32) string {
	s := []byte(fmt.Sprintf("%.2f", v))
	for len(s) > 0 && s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	return string(s)
}

// current is the app's theme. It's a signal so a component's StyleFunc
// restyles when the theme changes: live theme switching, no remount.
var current = signal.New(Native())

var (
	// appearance is the app's chosen light/dark scheme. It mirrors
	// ui.SetColorScheme: hosts report dark mode through the environment,
	// but some do so only after the window's own redraw (macOS), and an
	// explicit theme must switch its palette on the click, not one run
	// loop later.
	appearance = signal.New(ui.ColorSchemeSystem)
	// systemDark is the last system appearance seen while following the
	// system. A host reports its effective appearance, which a forced
	// Light or Dark changes too, so it is only recorded in System mode.
	systemDark bool
	// lastEnvDark is the last appearance the host reported, and ignoringEnv
	// is set when the app returns to the system's: the host's next report
	// can still describe the forced appearance it is leaving (macOS does),
	// and must not overwrite systemDark.
	lastEnvDark bool
	haveEnvDark bool
	ignoringEnv bool
)

// SetTheme changes the theme of everything already on screen.
func SetTheme(t Theme) { current.Set(t) }

// SetAppearance sets the app's own light/dark setting: the host's color
// scheme, and the appearance the design system's explicit palettes use.
// System follows the environment's Dark again.
func SetAppearance(s ui.ColorScheme) {
	prev := appearance.Peek()
	if s != prev {
		ignoringEnv = prev != ui.ColorSchemeSystem && s == ui.ColorSchemeSystem
	}
	ui.SetColorScheme(s)
	appearance.Set(s)
}

// ThemeNow returns the current theme, and tracks it when read in a
// StyleFunc.
func ThemeNow() Theme { return current.Get() }

// PaletteNow returns the current palette, and tracks the theme and the
// appearance when read in a StyleFunc.
func PaletteNow() Palette {
	t := current.Get()
	return t.paletteFor(appearance.Get(), darkNow())
}

// Native is the default: shadcn's token structure with the platform's own
// colors, type and radii. It looks like the platform, not like a brand, and
// gets dark mode, contrast and dynamic color for free.
func Native() Theme {
	return Theme{
		Name:     "native",
		Semantic: true,
		Radii:    Radii{SM: 6, MD: 8, LG: 12, XL: 16},
		Motion:   defaultMotion(),
		Type: TypeScale{
			Display:    fonts.LargeTitle,
			Title:      fonts.Title2,
			Heading:    fonts.Headline,
			Subheading: fonts.Subheadline,
			Body:       fonts.Body,
			Small:      fonts.Subheadline,
			Caption:    fonts.Caption,
		},
	}
}

func defaultMotion() Motion {
	return Motion{
		Instant: 0.12,
		Fast:    0.18,
		Normal:  0.25,
		Slow:    0.4,
		Curve:   ui.EaseInOut,
	}
}

// semanticPalette maps shadcn's roles onto Asuka's semantic colors.
func semanticPalette() Palette {
	p := Palette{
		Background:          colors.Background,
		Foreground:          colors.Label,
		Card:                colors.SecondaryGroupedBackground,
		CardForeground:      colors.Label,
		Popover:             colors.SecondaryBackground,
		PopoverForeground:   colors.Label,
		Primary:             colors.Accent,
		PrimaryForeground:   colors.OnAccent,
		Secondary:           colors.SecondaryBackground,
		SecondaryForeground: colors.Label,
		Muted:               colors.SecondaryBackground,
		MutedForeground:     colors.SecondaryLabel,
		Accent:              colors.Accent.Opacity(0.12),
		AccentForeground:    colors.Label,
		Destructive:         colors.Red,
		OnDestructive:       colors.RGB(255, 255, 255),
		Border:              colors.Separator,
		Input:               colors.Separator,
		Ring:                colors.Accent,
		Focus:               colors.Accent,
	}
	p.Elevation1 = Elevation{Overlay: colors.Label, Base: colors.Background, Alpha: 0.05,
		Shadow: ui.Shadow{Color: colors.Separator, Radius: 2, Y: 1}}
	p.Elevation2 = Elevation{Overlay: colors.Label, Base: colors.Background, Alpha: 0.1,
		Shadow: ui.Shadow{Color: colors.Separator, Radius: 8, Y: 2}}
	p.Elevation3 = Elevation{Overlay: colors.Label, Base: colors.Background, Alpha: 0.16,
		Shadow: ui.Shadow{Color: colors.Separator, Radius: 20, Y: 8}}
	return p
}

// nativePalette is the explicit light or dark palette a forced appearance
// uses: macOS's system colors, close enough to the platform's, and drawn
// as real colors rather than resolved by the host so a switch needs no
// environment round-trip.
func nativePalette(dark bool) Palette {
	if dark {
		p := Palette{
			Background:          colors.Hex(0x1e1e1e),
			Foreground:          colors.Hex(0xffffff),
			Card:                colors.Hex(0x2a2a2c),
			CardForeground:      colors.Hex(0xffffff),
			Popover:             colors.Hex(0x2a2a2c),
			PopoverForeground:   colors.Hex(0xffffff),
			Primary:             colors.Hex(0x0a66e0),
			PrimaryForeground:   colors.Hex(0xffffff),
			Secondary:           colors.Hex(0x2c2c2e),
			SecondaryForeground: colors.Hex(0xffffff),
			Muted:               colors.Hex(0x2c2c2e),
			MutedForeground:     colors.Hex(0x98989d),
			Accent:              colors.Hex(0x1f3a5f),
			AccentForeground:    colors.Hex(0xffffff),
			Destructive:         colors.Hex(0xd70015),
			OnDestructive:       colors.Hex(0xffffff),
			Border:              colors.RGBA(255, 255, 255, 41),
			Input:               colors.RGBA(255, 255, 255, 61),
			Ring:                colors.Hex(0x0a84ff),
			Focus:               colors.Hex(0x0a84ff),
		}
		neutralize(&p, false)
		return p
	}
	p := Palette{
		Background:          colors.Hex(0xffffff),
		Foreground:          colors.Hex(0x1d1d1f),
		Card:                colors.Hex(0xffffff),
		CardForeground:      colors.Hex(0x1d1d1f),
		Popover:             colors.Hex(0xffffff),
		PopoverForeground:   colors.Hex(0x1d1d1f),
		Primary:             colors.Hex(0x0a66e0),
		PrimaryForeground:   colors.Hex(0xffffff),
		Secondary:           colors.Hex(0xf2f2f7),
		SecondaryForeground: colors.Hex(0x1d1d1f),
		Muted:               colors.Hex(0xf2f2f7),
		MutedForeground:     colors.Hex(0x6b6b70),
		Accent:              colors.RGBA(10, 102, 224, 31),
		AccentForeground:    colors.Hex(0x1d1d1f),
		Destructive:         colors.Hex(0xd70015),
		OnDestructive:       colors.Hex(0xffffff),
		Border:              colors.RGBA(0, 0, 0, 31),
		Input:               colors.RGBA(0, 0, 0, 51),
		Ring:                colors.Hex(0x0a66e0),
		Focus:               colors.Hex(0x0a66e0),
	}
	neutralize(&p, true)
	return p
}

// Neutral is shadcn's default zinc palette with a 14pt body: the same look on
// every platform, at the cost of looking foreign on each of them.
func Neutral() Theme {
	t := Theme{
		Name: "neutral",
		Light: Palette{
			Background:          colors.Hex(0xffffff),
			Foreground:          colors.Hex(0x0a0a0a),
			Card:                colors.Hex(0xffffff),
			CardForeground:      colors.Hex(0x0a0a0a),
			Popover:             colors.Hex(0xffffff),
			PopoverForeground:   colors.Hex(0x0a0a0a),
			Primary:             colors.Hex(0x171717),
			PrimaryForeground:   colors.Hex(0xfafafa),
			Secondary:           colors.Hex(0xf5f5f5),
			SecondaryForeground: colors.Hex(0x171717),
			Muted:               colors.Hex(0xf5f5f5),
			MutedForeground:     colors.Hex(0x6b6b6b),
			Accent:              colors.Hex(0xf5f5f5),
			AccentForeground:    colors.Hex(0x171717),
			Destructive:         colors.Hex(0xe7000b),
			OnDestructive:       colors.Hex(0xffffff),
			Border:              colors.Hex(0xe5e5e5),
			Input:               colors.Hex(0xe5e5e5),
			Ring:                colors.Hex(0xa1a1a1),
			Focus:               colors.Hex(0x3b82f6),
		},
		Dark: Palette{
			Background:          colors.Hex(0x0a0a0a),
			Foreground:          colors.Hex(0xfafafa),
			Card:                colors.Hex(0x171717),
			CardForeground:      colors.Hex(0xfafafa),
			Popover:             colors.Hex(0x171717),
			PopoverForeground:   colors.Hex(0xfafafa),
			Primary:             colors.Hex(0xfafafa),
			PrimaryForeground:   colors.Hex(0x171717),
			Secondary:           colors.Hex(0x262626),
			SecondaryForeground: colors.Hex(0xfafafa),
			Muted:               colors.Hex(0x262626),
			MutedForeground:     colors.Hex(0xa1a1a1),
			Accent:              colors.Hex(0x262626),
			AccentForeground:    colors.Hex(0xfafafa),
			Destructive:         colors.Hex(0xff6467),
			OnDestructive:       colors.Hex(0x1a0505),
			Border:              colors.RGBA(255, 255, 255, 26),
			Input:               colors.RGBA(255, 255, 255, 38),
			Ring:                colors.Hex(0x737373),
			Focus:               colors.Hex(0x60a5fa),
		},
		Radii:  Radii{SM: 6, MD: 8, LG: 10, XL: 14},
		Motion: defaultMotion(),
		Type: TypeScale{
			Display:    fonts.System(30).Weighted(fonts.Bold),
			Title:      fonts.System(20).Weighted(fonts.Semibold),
			Heading:    fonts.System(16).Weighted(fonts.Semibold),
			Subheading: fonts.System(14).Weighted(fonts.Medium),
			Body:       fonts.System(14),
			Small:      fonts.System(13),
			Caption:    fonts.System(12),
		},
	}
	neutralize(&t.Light, true)
	neutralize(&t.Dark, false)
	return t
}

// neutralize derives elevation tokens from a palette's own surfaces, so a
// component can use c.Elevation2 without knowing which theme is active.
func neutralize(p *Palette, light bool) {
	overlay := colors.RGB(255, 255, 255)
	if !light {
		overlay = colors.RGB(0, 0, 0)
	}
	p.Elevation1 = Elevation{Overlay: overlay, Base: p.Secondary, Alpha: 0.04,
		Shadow: ui.Shadow{Color: p.Border, Radius: 2, Y: 1}}
	p.Elevation2 = Elevation{Overlay: overlay, Base: p.Card, Alpha: 0.06,
		Shadow: ui.Shadow{Color: p.Border, Radius: 8, Y: 2}}
	p.Elevation3 = Elevation{Overlay: overlay, Base: p.Popover, Alpha: 0.1,
		Shadow: ui.Shadow{Color: p.Border, Radius: 20, Y: 8}}
}

// themed makes a style that follows the theme and the appearance. Call it
// while building a component: it captures the screen's environment once.
func themed(f func(Theme, Palette) ui.Style) ui.StyleFunc {
	env := ui.UseEnvironment()
	return func() ui.Style {
		theme := current.Get()
		scheme := appearance.Get()
		return f(theme, theme.paletteFor(scheme, appearanceDark(scheme, env.Get().Dark)))
	}
}

func darkNow() bool {
	scheme := appearance.Get()
	return appearanceDark(scheme, ui.UseEnvironment().Get().Dark)
}

// appearanceDark resolves the app's scheme against the environment: Light
// and Dark override it, System follows it.
func appearanceDark(scheme ui.ColorScheme, envDark bool) bool {
	if !haveEnvDark || envDark != lastEnvDark {
		lastEnvDark, haveEnvDark = envDark, true
		if scheme == ui.ColorSchemeSystem {
			if ignoringEnv {
				ignoringEnv = false
			} else {
				systemDark = envDark
			}
		}
	}
	switch scheme {
	case ui.ColorSchemeLight:
		return false
	case ui.ColorSchemeDark:
		return true
	default:
		return systemDark
	}
}

// Themed builds a style from the theme, for the parts of an app the library
// doesn't cover: the same tokens, the same live switching.
//
//	ds.Themed(func(t ds.Theme, c ds.Palette) ui.Style {
//		return ui.Style{Background: c.Muted, CornerRadius: t.Radii.SM}
//	})
func Themed(f func(Theme, Palette) ui.Style) ui.Styler { return themed(f) }
