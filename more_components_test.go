package ds

import (
	"testing"
	"time"

	"github.com/egoist/asuka/ops"
	"github.com/egoist/asuka/signal"
	"github.com/egoist/asuka/ui"
	"github.com/egoist/asuka/uitest"
)

func TestTextareaIsAMultilineInput(t *testing.T) {
	value := signal.New("")
	a := uitest.Mount(t, func() ui.Node {
		return Textarea(TextareaProps{Label: "Bio", Value: value})
	})
	fields := a.Host.OfKind(ops.KindTextField)
	if len(fields) != 1 {
		t.Fatalf("want one text field, got %d", len(fields))
	}
	if !fields[0].Bool(ops.PropTextFieldMultiline) {
		t.Error("textarea is not multiline")
	}
	if !a.Has("Bio") {
		t.Errorf("label is not on screen:\n%s", a.Snapshot())
	}
}

func TestDatePickerUsesNativeControl(t *testing.T) {
	when := signal.New(time.Time{})
	a := uitest.Mount(t, func() ui.Node {
		return DatePicker(DatePickerProps{Label: "Start", Value: when})
	})
	if len(a.Host.OfKind(ops.KindDatePicker)) != 1 {
		t.Fatalf("want one native date picker:\n%s", a.Snapshot())
	}
	if !a.Has("Start") {
		t.Errorf("label is not on screen:\n%s", a.Snapshot())
	}
}

func TestBreadcrumbLinksEveryCrumbButTheCurrentOne(t *testing.T) {
	pressed := ""
	a := uitest.Mount(t, func() ui.Node {
		return Breadcrumb(BreadcrumbProps{Items: []BreadcrumbItem{
			{Title: "Home", OnPress: func() { pressed = "home" }},
			{Title: "Docs", OnPress: func() { pressed = "docs" }},
			{Title: "Button"},
		}})
	})
	a.Tap("Docs")
	if pressed != "docs" {
		t.Errorf("pressed = %q, want docs", pressed)
	}
	current := viewWithLabel(a, "Button")
	if current.Listens[ops.EventPress] {
		t.Error("the current crumb is pressable")
	}
	assertRole(t, current, ops.RoleText)
}

func TestPaginationMovesBetweenPages(t *testing.T) {
	page := signal.New(1)
	a := uitest.Mount(t, func() ui.Node {
		return Pagination(PaginationProps{Page: page, Total: 10})
	})
	a.Tap("Next page")
	if page.Peek() != 2 {
		t.Fatalf("page = %d, want 2", page.Peek())
	}
	a.Tap("Page 10")
	if page.Peek() != 10 {
		t.Fatalf("page = %d, want 10", page.Peek())
	}
	a.Tap("Previous page")
	if page.Peek() != 9 {
		t.Fatalf("page = %d, want 9", page.Peek())
	}
	a.Tap("Page 10")
	if got := viewWithLabel(a, "Page 10").String(ops.PropA11yHint); got != "Selected" {
		t.Errorf("current page hint = %q, want Selected", got)
	}
	next := viewWithLabel(a, "Next page")
	if next.Listens[ops.EventPress] {
		t.Error("next page is pressable at the end")
	}
	if got := next.String(ops.PropA11yHint); got != "Unavailable" {
		t.Errorf("next at the end hint = %q, want Unavailable", got)
	}
}

func TestSpinnerIsIndeterminateProgress(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node {
		return Spinner(SpinnerProps{Label: "Loading"})
	})
	rings := a.Host.OfKind(ops.KindProgressRing)
	if len(rings) != 1 {
		t.Fatalf("want one progress ring, got %d", len(rings))
	}
	if !rings[0].Bool(ops.PropProgressRingIndeterminate) {
		t.Error("spinner is not indeterminate")
	}
	if !a.Has("Loading") {
		t.Errorf("label is not on screen:\n%s", a.Snapshot())
	}
}

func TestEmptyStateShowsAction(t *testing.T) {
	pressed := false
	a := uitest.Mount(t, func() ui.Node {
		return EmptyState(EmptyStateProps{
			Title:       "No projects",
			Description: "Create one to get started.",
			Action: func() ui.Node {
				return Button(ButtonProps{Text: "New project", OnPress: func() { pressed = true }})
			},
		})
	})
	if !a.Has("No projects") || !a.Has("Create one to get started.") {
		t.Fatalf("empty state copy is missing:\n%s", a.Snapshot())
	}
	a.Tap("New project")
	if !pressed {
		t.Error("empty state action did not fire")
	}
}

func TestKbdShowsEveryKey(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node { return Kbd("⌘", "K") })
	if !a.Has("⌘") || !a.Has("K") {
		t.Errorf("keycaps are missing:\n%s", a.Snapshot())
	}
}
