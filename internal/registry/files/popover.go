package ds

import (
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
)

// PopoverProps configures a Popover. Its children are the trigger's
// contents: a label, an icon, or a small row. Popover makes the trigger
// pressable, so children must not be another button.
type PopoverProps struct {
	// Label is the trigger's accessible name, and the text it shows when
	// there are no children.
	Label string
	// Open controls the popover. When nil, the trigger keeps its own state;
	// pass a signal to open it from elsewhere.
	Open *signal.Signal[bool]
	// Content is built when the popover opens and disposed when it closes.
	Content func() ui.Node
	// MaxWidth is how wide the bubble may grow (0: the platform's default).
	MaxWidth float32
	// Style overrides the trigger's style.
	Style     ui.Styler
	OnDismiss func()
}

// Popover attaches a bubble of content to a trigger: a few options, a
// detail or a small form, dismissed by tapping outside. The bubble is the
// platform's own (NSPopover, UIPopoverPresentationController, a Material
// popup), so it keeps the platform's placement, focus and keyboard
// behavior. Use Menu for commands and Dialog for confirmations.
func Popover(p PopoverProps, children ...ui.Node) ui.Node {
	open := p.Open
	if open == nil {
		open = signal.New(false)
	}
	toggle := func() { open.Set(!open.Peek()) }

	var body []ui.Node
	if len(children) > 0 {
		body = children
	} else {
		body = []ui.Node{
			ui.StaticText(p.Label, themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Font: t.Type.Subheading, Color: c.Foreground, Lines: 1, Grow: 1}
			})),
			ui.Symbol(ui.IconDown, themed(func(t Theme, c Palette) ui.Style {
				return ui.Style{Font: t.Type.Caption, Color: c.MutedForeground}
			})),
		}
	}

	style := p.Style
	if style == nil {
		style = themed(func(t Theme, c Palette) ui.Style {
			return ui.Style{
				Direction:    ui.RowDir,
				Align:        ui.AlignCenter,
				Gap:          Space2,
				Padding:      ui.Symmetric(0, Space3),
				Height:       ui.Pt(36),
				MinWidth:     ui.Pt(160),
				Background:   c.Background,
				BorderWidth:  1,
				BorderColor:  c.Input,
				CornerRadius: t.Radii.MD,
				Hovered:      &ui.Style{Background: c.Accent},
			}
		})
	}
	trigger := Pressable(PressableProps{
		Label:    p.Label,
		Style:    style,
		Expanded: open,
		OnPress:  toggle,
	}, body...)

	content := p.Content
	if content == nil {
		content = func() ui.Node { return ui.View(ui.ViewProps{}) }
	}
	panel := func() ui.Node {
		return ui.View(ui.ViewProps{Style: ui.Style{
			Padding:  ui.All(Space3),
			Gap:      Space2,
			MinWidth: ui.Pt(200),
		}}, content())
	}
	return ui.Popover(ui.PopoverProps{
		Open:      open,
		Content:   panel,
		MaxWidth:  p.MaxWidth,
		OnDismiss: p.OnDismiss,
	}, trigger)
}
