package main

import (
	"testing"

	"github.com/egoist/asuka/colors"
	"github.com/egoist/asuka/ui"
)

func TestHighlightPreservesSource(t *testing.T) {
	for _, tc := range []struct {
		language string
		source   string
	}{
		{"Go", "func App() ui.Node {\n\treturn ds.Button(ds.ButtonProps{Text: \"Save\", OnPress: save}) // drawn\n}\n"},
		{"Shell", "asuka-ds add switch --dir internal/ds  # resolve the kit\ngo test ./...\n"},
		{"JSON", "{\n  \"kits\": {\n    \"combobox\": { \"files\": [\"combobox.go\"], \"requires\": [\"overlay\"] }\n  }\n}\n"},
	} {
		if got := highlight(tc.language, tc.source).String(); got != tc.source {
			t.Fatalf("%s highlight changed the snippet:\n%q\nwant:\n%q", tc.language, got, tc.source)
		}
	}
}

func TestHighlightGoKeywordsStringsAndComments(t *testing.T) {
	spans := highlight("Go", `package main // entry
func App() ui.Node { return ds.Text("hi") }
var count int = 3`)
	want := map[string]ui.Color{
		"package":  colors.Purple,
		"func":     colors.Purple,
		"var":      colors.Purple,
		"int":      colors.Blue,
		"App":      colors.Blue,
		"ui":       colors.Blue,
		"Node":     colors.Blue,
		"Text":     colors.Blue,
		`"hi"`:     colors.Green,
		"3":        colors.Orange,
		"// entry": colors.SecondaryLabel,
	}
	for text, color := range want {
		span := findSpan(spans, text)
		if span == nil {
			t.Fatalf("no span for %q in %#v", text, spans)
		}
		if span.Color != color {
			t.Errorf("%q color = %v, want %v", text, span.Color, color)
		}
	}
	if findSpan(spans, "// entry").Italic != true {
		t.Error("comment should be italic")
	}
}

func TestHighlightGoNamesInSnippets(t *testing.T) {
	spans := highlight("Go", `package main
func render() {
	label := ds.ButtonProps{Text: "Save"}
	_ = signal.New("")
}`)
	for text, color := range map[string]ui.Color{
		"render":      colors.Blue,
		"label":       colors.Blue,
		"ds":          colors.Blue,
		"ButtonProps": colors.Blue,
		"Text":        colors.Blue,
		"signal":      colors.Blue,
		"New":         colors.Blue,
	} {
		span := findSpan(spans, text)
		if span == nil {
			t.Fatalf("no span for %q in %#v", text, spans)
		}
		if span.Color != color {
			t.Errorf("%q color = %v, want %v", text, span.Color, color)
		}
	}
	if span := findSpan(spans, "main"); span.Color != 0 {
		t.Errorf("package name main should stay plain, got %v", span.Color)
	}
}

func TestHighlightShellCommandsFlagsAndComments(t *testing.T) {
	spans := highlight("Shell", "asuka-ds add switch --dir internal/ds  # refresh\n")
	for text, color := range map[string]ui.Color{
		"asuka-ds":  colors.Purple,
		"--dir":     colors.Pink,
		"# refresh": colors.SecondaryLabel,
	} {
		span := findSpan(spans, text)
		if span == nil {
			t.Fatalf("no span for %q in %#v", text, spans)
		}
		if span.Color != color {
			t.Errorf("%q color = %v, want %v", text, span.Color, color)
		}
	}
}

func TestHighlightJSONKeysStringsAndNumbers(t *testing.T) {
	spans := highlight("JSON", `{"kits": {"combobox": 2}, "ok": true}`)
	for text, color := range map[string]ui.Color{
		`"kits"`:     colors.Blue,
		`"combobox"`: colors.Blue,
		"2":          colors.Orange,
		"true":       colors.Purple,
	} {
		span := findSpan(spans, text)
		if span == nil {
			t.Fatalf("no span for %q in %#v", text, spans)
		}
		if span.Color != color {
			t.Errorf("%q color = %v, want %v", text, span.Color, color)
		}
	}
}

func findSpan(spans ui.Rich, text string) *ui.Span {
	for i := range spans {
		if spans[i].Text == text {
			return &spans[i]
		}
	}
	return nil
}
