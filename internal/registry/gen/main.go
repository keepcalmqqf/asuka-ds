// Command gen copies the repository's component sources into the registry
// package, where they are embedded in the CLI.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/egoist/asuka-ds/internal/catalog"
)

type manifest struct {
	Kits map[string]struct {
		Files []string `json:"files"`
	} `json:"kits"`
}

func main() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		fail("cannot locate generator")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
	manifestPath := filepath.Join(root, "internal", "catalog", "components.json")
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		fail(err.Error())
	}
	var m manifest
	if err := json.Unmarshal(raw, &m); err != nil {
		fail(err.Error())
	}
	unique := map[string]bool{}
	for _, kit := range m.Kits {
		for _, name := range kit.Files {
			unique[name] = true
		}
	}
	names := make([]string, 0, len(unique))
	for name := range unique {
		names = append(names, name)
	}
	sort.Strings(names)
	out := filepath.Join(root, "internal", "registry", "files")
	if err := os.RemoveAll(out); err != nil {
		fail(err.Error())
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		fail(err.Error())
	}
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(root, catalog.SourceDir, name))
		if err != nil {
			fail(err.Error())
		}
		if err := os.WriteFile(filepath.Join(out, name), data, 0o644); err != nil {
			fail(err.Error())
		}
	}
	fmt.Printf("registry: %d files\n", len(names))
}

func fail(message string) {
	fmt.Fprintln(os.Stderr, "registry:", message)
	os.Exit(1)
}
