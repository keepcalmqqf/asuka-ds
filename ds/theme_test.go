package ds

import (
	"strings"
	"testing"

	"github.com/egoist/asuka/colors"
	"github.com/egoist/asuka/ui"
)

func TestBuiltinThemesValidate(t *testing.T) {
	for _, theme := range []Theme{Native(), Neutral()} {
		if problems := theme.Validate(); len(problems) != 0 {
			t.Errorf("%s: %s", theme.Name, strings.Join(problems, "; "))
		}
	}
}

func TestValidateCatchesMissingRoles(t *testing.T) {
	theme := Neutral()
	theme.Light.MutedForeground = 0
	problems := theme.Validate()
	if !contains(problems, "Light.MutedForeground is unset") {
		t.Fatalf("missing role not reported: %v", problems)
	}
}

func TestValidateCatchesLowContrast(t *testing.T) {
	theme := Neutral()
	theme.Light.Foreground = colors.Hex(0xfafafa)
	problems := theme.Validate()
	if !contains(problems, "Light.Foreground/Background: contrast") {
		t.Fatalf("low contrast not reported: %v", problems)
	}
}

func TestContrastRatio(t *testing.T) {
	cases := []struct {
		a, b ui.Color
		want float32
	}{
		{colors.Hex(0x000000), colors.Hex(0xffffff), 21},
		{colors.Hex(0xffffff), colors.Hex(0xffffff), 1},
		{colors.Hex(0x0a0a0a), colors.Hex(0xffffff), 19.6},
	}
	for _, c := range cases {
		got := ContrastRatio(c.a, c.b)
		if diff := got - c.want; diff > 0.5 || diff < -0.5 {
			t.Errorf("ContrastRatio(%v, %v) = %.2f, want %.2f", c.a, c.b, got, c.want)
		}
	}
}

func TestMotionForReduceMotion(t *testing.T) {
	theme := Neutral()
	m := theme.MotionFor(true)
	if m.Fast != 0 || m.Normal != 0 || m.Slow != 0 {
		t.Errorf("reduced motion should zero durations: %+v", m)
	}
	if theme.MotionFor(false).Normal == 0 {
		t.Error("normal motion lost its duration")
	}
	theme.Motion.Disabled = true
	if theme.MotionFor(false).Fast != 0 {
		t.Error("Disabled=true should zero durations")
	}
}

func TestElevationTokensArePresent(t *testing.T) {
	for _, theme := range []Theme{Native(), Neutral()} {
		if theme.Semantic {
			if theme.PaletteFor(false).Elevation2.Alpha == 0 {
				t.Error("semantic palette lost its elevation tokens")
			}
			continue
		}
		for _, name := range []string{"Light", "Dark"} {
			palette := theme.Light
			if name == "Dark" {
				palette = theme.Dark
			}
			if palette.Elevation2.Alpha == 0 {
				t.Errorf("%s.%s: elevation2 is unset", theme.Name, name)
			}
		}
	}
}

// The explicit palettes a forced appearance uses meet the same contrast
// bar as the named themes.
func TestNativeForcedPalettesValidate(t *testing.T) {
	for name, dark := range map[string]bool{"Light": false, "Dark": true} {
		if problems := validatePalette(name, nativePalette(dark)); len(problems) != 0 {
			t.Errorf("%s: %s", name, strings.Join(problems, "; "))
		}
	}
}

func contains(ss []string, sub string) bool {
	for _, s := range ss {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
