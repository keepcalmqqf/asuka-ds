// Package catalog is the product's source of truth: what components exist,
// what files they install, and what they depend on.
package catalog

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
)

//go:embed components.json
var raw []byte

// Kit is a set of files installed together.
type Kit struct {
	Files    []string `json:"files"`
	Requires []string `json:"requires"`
}

// Component is a public, installable component.
type Component struct {
	Title       string `json:"title"`
	Category    string `json:"category"`
	Kit         string `json:"kit"`
	Description string `json:"description"`
}

// Catalog is the decoded components.json.
type Catalog struct {
	Kits       map[string]Kit       `json:"kits"`
	Components map[string]Component `json:"components"`
}

// Load returns the embedded catalog.
func Load() (Catalog, error) {
	var c Catalog
	if err := json.Unmarshal(raw, &c); err != nil {
		return Catalog{}, fmt.Errorf("catalog: %w", err)
	}
	if len(c.Kits) == 0 || len(c.Components) == 0 {
		return Catalog{}, fmt.Errorf("catalog: empty")
	}
	return c, nil
}

// Names returns the component names in a stable order.
func (c Catalog) Names() []string {
	names := make([]string, 0, len(c.Components))
	for name := range c.Components {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Files resolves names to unique source files, following kit requirements.
// A missing component or kit is an error.
func (c Catalog) Files(names ...string) ([]string, error) {
	seenKit := map[string]bool{}
	seenFile := map[string]bool{}
	var files []string
	var visit func(kit string) error
	visit = func(kit string) error {
		if seenKit[kit] {
			return nil
		}
		k, ok := c.Kits[kit]
		if !ok {
			return fmt.Errorf("unknown kit %q", kit)
		}
		seenKit[kit] = true
		for _, require := range k.Requires {
			if err := visit(require); err != nil {
				return err
			}
		}
		for _, file := range k.Files {
			if seenFile[file] {
				continue
			}
			seenFile[file] = true
			files = append(files, file)
		}
		return nil
	}
	for _, name := range names {
		component, ok := c.Components[name]
		if !ok {
			return nil, fmt.Errorf("unknown component %q", name)
		}
		if err := visit(component.Kit); err != nil {
			return nil, err
		}
	}
	sort.Strings(files)
	return files, nil
}

// Categories returns component names grouped by category, in the order
// categories first appear in a fixed product order.
func (c Catalog) Categories() []Category {
	order := []string{
		"Button", "Form", "Navigation", "Overlay", "Data display",
		"Feedback", "Typography",
	}
	groups := map[string][]string{}
	for name, component := range c.Components {
		groups[component.Category] = append(groups[component.Category], name)
	}
	out := make([]Category, 0, len(groups))
	for _, name := range order {
		if len(groups[name]) == 0 {
			continue
		}
		sort.Strings(groups[name])
		out = append(out, Category{Name: name, Components: groups[name]})
	}
	return out
}

// Category is a named group of components.
type Category struct {
	Name       string
	Components []string
}
