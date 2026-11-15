package catalog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCatalogIsConsistent(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for name, kit := range c.Kits {
		if len(kit.Files) == 0 {
			t.Errorf("kit %q installs no files", name)
		}
		for _, file := range kit.Files {
			if !strings.HasSuffix(file, ".go") {
				t.Errorf("kit %q has a non-Go file %q", name, file)
			}
			if _, err := os.Stat(filepath.Join("..", "..", SourceDir, file)); err != nil {
				t.Errorf("kit %q: %v", name, err)
			}
		}
		for _, require := range kit.Requires {
			if _, ok := c.Kits[require]; !ok {
				t.Errorf("kit %q requires unknown kit %q", name, require)
			}
		}
	}
	for name, component := range c.Components {
		if component.Title == "" || component.Category == "" || component.Description == "" {
			t.Errorf("component %q is missing a title, category or description", name)
		}
		if _, ok := c.Kits[component.Kit]; !ok {
			t.Errorf("component %q uses unknown kit %q", name, component.Kit)
		}
	}
}

func TestCategoriesCoverEveryComponent(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, category := range c.Categories() {
		for _, name := range category.Components {
			if seen[name] {
				t.Errorf("component %q is in two categories", name)
			}
			seen[name] = true
		}
	}
	for _, name := range c.Names() {
		if !seen[name] {
			t.Errorf("component %q is in no category", name)
		}
	}
}

func TestFilesFollowDependencies(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	files, err := c.Files("command")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"command.go", "button.go", "overlay.go", "theme.go"} {
		if !contains(files, want) {
			t.Errorf("command does not install %s: %v", want, files)
		}
	}
	if _, err := c.Files("not-a-component"); err == nil {
		t.Error("an unknown component resolved")
	}
}

func contains(list []string, want string) bool {
	for _, item := range list {
		if item == want {
			return true
		}
	}
	return false
}
