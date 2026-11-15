package registry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/egoist/asuka-ds/internal/catalog"
)

// The embedded copies must match the sources a reader sees in the repo.
// Run `go generate ./internal/registry` after changing a component.
func TestEmbeddedSourcesMatchRepository(t *testing.T) {
	r, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	names := r.Catalog.Names()
	files, err := r.Catalog.Files(names...)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range files {
		want, err := os.ReadFile(filepath.Join("..", "..", catalog.SourceDir, name))
		if err != nil {
			t.Fatal(err)
		}
		got, err := r.Source(name)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != string(want) {
			t.Errorf("%s is stale: run go generate ./internal/registry", name)
		}
	}
}
