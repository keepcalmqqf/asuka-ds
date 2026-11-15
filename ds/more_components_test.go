package ds

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/egoist/asuka/host/fake"
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

func TestComboboxFiltersAndSelects(t *testing.T) {
	selected := signal.New("")
	a := uitest.Mount(t, func() ui.Node {
		return Combobox(ComboboxProps{
			Label:       "Framework",
			Placeholder: "Choose…",
			Selected:    selected,
			Items: []ComboboxItem{
				{Value: "asuka", Title: "Asuka", Keywords: "go native"},
				{Value: "swiftui", Title: "SwiftUI"},
				{Value: "compose", Title: "Compose"},
			},
		})
	})
	a.Tap("Choose…")
	if !a.Has("SwiftUI") {
		t.Fatalf("the list did not open:\n%s", a.Snapshot())
	}
	a.Type("Search…", "swift")
	if a.Has("Asuka") || !a.Has("SwiftUI") {
		t.Fatalf("filtering is wrong:\n%s", a.Snapshot())
	}
	a.Tap("SwiftUI")
	if selected.Peek() != "swiftui" {
		t.Errorf("selected = %q, want swiftui", selected.Peek())
	}
	if !a.Has("SwiftUI") {
		t.Errorf("the trigger does not show the choice:\n%s", a.Snapshot())
	}
}

func TestDateRangePickerUsesTwoNativePickers(t *testing.T) {
	start := signal.New(time.Time{})
	end := signal.New(time.Time{})
	a := uitest.Mount(t, func() ui.Node {
		return DateRangePicker(DateRangePickerProps{
			Label: "Billing period", Start: start, End: end,
		})
	})
	if len(a.Host.OfKind(ops.KindDatePicker)) != 2 {
		t.Fatalf("want two date pickers:\n%s", a.Snapshot())
	}
	for _, label := range []string{"Start date", "End date"} {
		if maybeViewWithLabel(a, label) == nil {
			t.Errorf("%s has no accessible name:\n%s", label, a.Snapshot())
		}
	}
}

func TestSegmentedControlUsesNativeSegments(t *testing.T) {
	selected := signal.New(0)
	a := uitest.Mount(t, func() ui.Node {
		return SegmentedControl(SegmentedControlProps{
			Label: "Size", Segments: []string{"Small", "Medium", "Large"},
			Selected: selected,
		})
	})
	if len(a.Host.OfKind(ops.KindSegmented)) != 1 {
		t.Fatalf("want one segmented control:\n%s", a.Snapshot())
	}
	a.Select("Large")
	if selected.Peek() != 2 {
		t.Errorf("selected = %d, want 2", selected.Peek())
	}
}

func TestMenuSupportsSectionsAndSubmenus(t *testing.T) {
	a := uitest.Mount(t, func() ui.Node {
		return Menu(MenuProps{Label: "Actions", Items: []DialogAction{
			{Title: "Edit", Children: []DialogAction{{Title: "Copy"}}},
			{Title: "Danger", Inline: true, Children: []DialogAction{{Title: "Delete"}}},
		}})
	})
	var menu *fake.View
	for _, v := range a.Host.Views {
		if _, ok := v.Props[ops.PropMenu]; ok {
			menu = v
		}
	}
	if menu == nil {
		t.Fatalf("menu trigger has no menu prop:\n%s", a.Snapshot())
	}
	items := ops.DecodeActions(menu.Props[ops.PropMenu].Bytes)
	if len(items) != 2 {
		t.Fatalf("items = %d, want 2", len(items))
	}
	if !items[0].Is(ops.MenuSubmenu) || len(items[0].Items) != 1 ||
		items[0].Items[0].Title != "Copy" {
		t.Errorf("Edit is not a submenu with Copy: %+v", items[0])
	}
	if !items[1].Is(ops.MenuInline) {
		t.Errorf("Danger is not an inline section: %+v", items[1])
	}
}

func TestToggleWritesItsSignal(t *testing.T) {
	pressed := signal.New(false)
	a := uitest.Mount(t, func() ui.Node {
		return Toggle(ToggleProps{Text: "Bold", Pressed: pressed})
	})
	a.Tap("Bold")
	if !pressed.Peek() {
		t.Fatal("toggle did not write its signal")
	}
	if got := viewWithLabel(a, "Bold").String(ops.PropA11yHint); got != "Selected" {
		t.Errorf("hint = %q, want Selected", got)
	}
}

func TestToggleGroupMultiple(t *testing.T) {
	selected := signal.NewFunc([]string{}, slices.Equal)
	a := uitest.Mount(t, func() ui.Node {
		return ToggleGroup(ToggleGroupProps{
			Label: "Formatting",
			Items: []ToggleItem{
				{Value: "bold", Label: "Bold"},
				{Value: "italic", Label: "Italic"},
			},
			Selected: selected,
		})
	})
	a.Tap("Bold")
	a.Tap("Italic")
	if got := selected.Peek(); len(got) != 2 || got[0] != "bold" || got[1] != "italic" {
		t.Errorf("selected = %v, want [bold italic]", got)
	}
	a.Tap("Bold")
	if got := selected.Peek(); len(got) != 1 || got[0] != "italic" {
		t.Errorf("selected after deselect = %v, want [italic]", got)
	}
}

func TestToggleGroupSingleKeepsOne(t *testing.T) {
	selected := signal.NewFunc([]string{}, slices.Equal)
	a := uitest.Mount(t, func() ui.Node {
		return ToggleGroup(ToggleGroupProps{
			Label: "Align", Single: true,
			Items: []ToggleItem{
				{Value: "left", Label: "Left"},
				{Value: "right", Label: "Right"},
			},
			Selected: selected,
		})
	})
	a.Tap("Left")
	a.Tap("Right")
	if got := selected.Peek(); len(got) != 1 || got[0] != "right" {
		t.Errorf("selected = %v, want [right]", got)
	}
	a.Tap("Right")
	if got := selected.Peek(); len(got) != 0 {
		t.Errorf("selected after deselect = %v, want []", got)
	}
}

func TestNumberFieldUsesNativeStepper(t *testing.T) {
	value := signal.New(3.0)
	a := uitest.Mount(t, func() ui.Node {
		return NumberField(NumberFieldProps{Label: "Seats", Value: value, Min: 1, Max: 10})
	})
	if len(a.Host.OfKind(ops.KindStepper)) != 1 {
		t.Fatalf("want one stepper:\n%s", a.Snapshot())
	}
	if !a.Has("3") {
		t.Errorf("value is not shown:\n%s", a.Snapshot())
	}
	value.Set(4)
	a.Settle()
	if !a.Has("4") {
		t.Errorf("value did not follow the signal:\n%s", a.Snapshot())
	}
}

func TestDataTableSortsAndPressesRows(t *testing.T) {
	sortColumn := signal.New(-1)
	descending := signal.New(false)
	pressed := -1
	a := uitest.Mount(t, func() ui.Node {
		return DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Title: "Name", Sortable: true},
				{Title: "Role"},
			},
			Rows: [][]string{
				{"Grace", "Admiral"},
				{"Ada", "Engineer"},
			},
			Sort: sortColumn, Descending: descending,
			OnRowPress: func(index int) { pressed = index },
		})
	})
	a.Tap("Ada, Engineer")
	if pressed != 1 {
		t.Errorf("pressed row = %d, want 1", pressed)
	}
	a.Tap("Sort by Name")
	if sortColumn.Peek() != 0 || descending.Peek() {
		t.Fatalf("sort = %d descending = %v", sortColumn.Peek(), descending.Peek())
	}
	snap := a.Snapshot()
	if strings.Index(snap, "Ada") > strings.Index(snap, "Grace") {
		t.Errorf("rows are not ascending:\n%s", snap)
	}
	a.Tap("Sort by Name, ascending")
	if !descending.Peek() {
		t.Fatal("second press did not reverse the order")
	}
	snap = a.Snapshot()
	if strings.Index(snap, "Grace") > strings.Index(snap, "Ada") {
		t.Errorf("rows are not descending:\n%s", snap)
	}
}
