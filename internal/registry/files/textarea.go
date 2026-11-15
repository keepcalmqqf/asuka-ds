package ds

import (
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// TextareaProps configures a Textarea.
type TextareaProps struct {
	Value       *signal.Signal[string]
	Placeholder string
	Label       string
	Description string
	Error       string
	// Lines is the field's visible height, in lines (0: 4).
	Lines    int
	Disabled bool
}

// Textarea is a multi-line Input: a comment, a bio, a message. It is the
// platform's own text field, so IME, selection, autofill and screen readers
// behave, wrapped in the design system's border and Field semantics.
func Textarea(p TextareaProps) ui.Node {
	lines := p.Lines
	if lines <= 0 {
		lines = 4
	}
	return Input(InputProps{
		Value:       p.Value,
		Placeholder: p.Placeholder,
		Label:       p.Label,
		Description: p.Description,
		Error:       p.Error,
		Disabled:    p.Disabled,
		Multiline:   true,
		Lines:       lines,
	})
}
