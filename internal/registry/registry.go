// Package registry embeds the component sources the CLI installs.
package registry

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/egoist/asuka-ds/internal/catalog"
)

//go:generate go run ./gen

//go:embed files/*.go
var sources embed.FS

// Registry is the catalog and the source files it names.
type Registry struct {
	Catalog catalog.Catalog
	Files   fs.FS
}

// Load returns the embedded registry.
func Load() (*Registry, error) {
	c, err := catalog.Load()
	if err != nil {
		return nil, err
	}
	return &Registry{Catalog: c, Files: sources}, nil
}

// Source returns one component source file.
func (r *Registry) Source(name string) ([]byte, error) {
	data, err := fs.ReadFile(r.Files, "files/"+name)
	if err != nil {
		return nil, fmt.Errorf("registry: %s: %w", name, err)
	}
	return data, nil
}
