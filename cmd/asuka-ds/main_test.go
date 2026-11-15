package main

import (
	"bytes"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/egoist/asuka-ds/internal/catalog"
)

func runTest(t *testing.T, args ...string) string {
	t.Helper()
	var out bytes.Buffer
	if err := run(args, &out, io.Discard); err != nil {
		t.Fatalf("run %v: %v", args, err)
	}
	return out.String()
}

func TestInitInstallsTheBaseKit(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "ds")
	out := runTest(t, "init", "--dir", dir)
	if !strings.Contains(out, "Next: asuka-ds add") {
		t.Errorf("init output has no next step:\n%s", out)
	}
	for _, name := range []string{"theme.go", "text.go", "pressable.go", "field.go", "variant.go"} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		assertPackage(t, name, data, "ds")
	}
	if _, err := os.Stat(filepath.Join(dir, "button.go")); !os.IsNotExist(err) {
		t.Error("init should not install Button")
	}
}

func TestAddResolvesDependencies(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "ui")
	runTest(t, "add", "switch", "--dir", dir, "--package", "ui")
	for _, name := range []string{"controls.go", "overlay.go", "field.go", "theme.go"} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		assertPackage(t, name, data, "ui")
	}
}

func assertPackage(t *testing.T, name string, src []byte, want string) {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), name, src, parser.PackageClauseOnly)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	if file.Name.Name != want {
		t.Errorf("%s is package %s, want %s", name, file.Name.Name, want)
	}
}

func TestAddNeverOverwritesWithoutForce(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "ds")
	runTest(t, "add", "button", "--dir", dir)
	path := filepath.Join(dir, "button.go")
	const mine = "package ds\n\n// my button\n"
	if err := os.WriteFile(path, []byte(mine), 0o644); err != nil {
		t.Fatal(err)
	}
	out := runTest(t, "add", "button", "--dir", dir)
	if !strings.Contains(out, "use --force") {
		t.Errorf("add did not report the existing file:\n%s", out)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != mine {
		t.Error("add overwrote a file without --force")
	}
	runTest(t, "add", "button", "--dir", dir, "--force")
	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) == mine {
		t.Error("--force did not replace the file")
	}
}

func TestUnknownComponent(t *testing.T) {
	var out bytes.Buffer
	err := run([]string{"add", "not-a-component", "--dir", t.TempDir()}, &out, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "unknown component") {
		t.Fatalf("err = %v, want unknown component", err)
	}
}

func TestListShowsCategories(t *testing.T) {
	out := runTest(t, "list")
	for _, want := range []string{"Button", "Form", "button", "combobox", "Data display"} {
		if !strings.Contains(out, want) {
			t.Errorf("list is missing %q:\n%s", want, out)
		}
	}
}

func TestVersion(t *testing.T) {
	for _, arg := range []string{"version", "--version"} {
		out := runTest(t, arg)
		if !strings.Contains(out, Version) {
			t.Errorf("%s printed %q, want %s", arg, out, Version)
		}
	}
}

func TestEveryCatalogComponentInstalls(t *testing.T) {
	c, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range c.Names() {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			runTest(t, "add", name, "--dir", dir)
			files, err := c.Files(name)
			if err != nil {
				t.Fatal(err)
			}
			for _, file := range files {
				data, err := os.ReadFile(filepath.Join(dir, file))
				if err != nil {
					t.Fatalf("%s: %v", file, err)
				}
				assertPackage(t, file, data, "ds")
			}
		})
	}
}
