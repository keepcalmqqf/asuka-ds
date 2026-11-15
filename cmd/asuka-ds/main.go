// Command asuka-ds installs Asuka DS components into a Go project.
//
//	asuka-ds init [--dir internal/ds] [--package ds]
//	asuka-ds add button input [--dir internal/ds] [--package ds] [--force]
//	asuka-ds list
//
// The component sources are embedded in the binary, so adding a component
// needs no network access and no checkout of this repository.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/egoist/asuka-ds/internal/registry"
)

// Version is the CLI and library version. Release builds replace it with
// the tag through -ldflags "-X main.Version=...".
var Version = "0.1.0"

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "asuka-ds:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		usage(stdout)
		return nil
	}
	switch args[0] {
	case "init":
		return cmdInit(args[1:], stdout, stderr)
	case "add":
		return cmdAdd(args[1:], stdout, stderr)
	case "list":
		return cmdList(args[1:], stdout, stderr)
	case "version", "-v", "--version":
		fmt.Fprintln(stdout, Version)
		return nil
	case "help", "-h", "--help":
		usage(stdout)
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func usage(w io.Writer) {
	fmt.Fprint(w, `asuka-ds installs components into a Go project.

usage:
  asuka-ds init [--dir internal/ds] [--package ds] [--force]
  asuka-ds add <component>... [--dir internal/ds] [--package ds] [--force]
  asuka-ds list
  asuka-ds version

Components are copied as source you own; re-running add never overwrites
your edits unless --force is given.
`)
}

type installFlags struct {
	dir     string
	pkg     string
	force   bool
	command string
}

func newInstallFlags(command string, stderr io.Writer) (*flag.FlagSet, *installFlags) {
	opts := &installFlags{dir: filepath.Join("internal", "ds"), pkg: "ds", command: command}
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&opts.dir, "dir", opts.dir, "the directory to install into")
	fs.StringVar(&opts.pkg, "package", opts.pkg, "the Go package name")
	fs.BoolVar(&opts.force, "force", false, "overwrite files that already exist")
	return fs, opts
}

func cmdInit(args []string, stdout, stderr io.Writer) error {
	fs, opts := newInstallFlags("init", stderr)
	rest, err := parseInstall(fs, args)
	if err != nil {
		return err
	}
	if len(rest) > 0 {
		return fmt.Errorf("init takes no arguments, got %q", strings.Join(rest, " "))
	}
	r, err := registry.Load()
	if err != nil {
		return err
	}
	files, err := r.Catalog.Files("text")
	if err != nil {
		return err
	}
	fmt.Fprintln(stdout, "Asuka DS")
	if err := install(r, files, *opts, stdout); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "\nNext: asuka-ds add button input --dir %s\n", opts.dir)
	return nil
}

func cmdAdd(args []string, stdout, stderr io.Writer) error {
	fs, opts := newInstallFlags("add", stderr)
	names, err := parseInstall(fs, args)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return fmt.Errorf("add needs at least one component; try asuka-ds list")
	}
	r, err := registry.Load()
	if err != nil {
		return err
	}
	files, err := r.Catalog.Files(names...)
	if err != nil {
		return err
	}
	return install(r, files, *opts, stdout)
}

// parseInstall accepts flags before or after the component names, as
// `asuka-ds add button --dir internal/ds` reads better than the standard
// flag package's flags-first rule.
func parseInstall(fs *flag.FlagSet, args []string) ([]string, error) {
	var flags, rest []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--dir" || arg == "-dir" || arg == "--package" || arg == "-package":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("%s needs a value", arg)
			}
			flags = append(flags, arg, args[i+1])
			i++
		case strings.HasPrefix(arg, "--dir=") || strings.HasPrefix(arg, "-dir=") ||
			strings.HasPrefix(arg, "--package=") || strings.HasPrefix(arg, "-package="):
			flags = append(flags, arg)
		case arg == "--force" || arg == "-force":
			flags = append(flags, arg)
		case strings.HasPrefix(arg, "-"):
			return nil, fmt.Errorf("unknown flag %q", arg)
		default:
			rest = append(rest, arg)
		}
	}
	if err := fs.Parse(flags); err != nil {
		return nil, err
	}
	return rest, nil
}

func cmdList(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		return err
	}
	r, err := registry.Load()
	if err != nil {
		return err
	}
	for _, category := range r.Catalog.Categories() {
		fmt.Fprintln(stdout, category.Name)
		for _, name := range category.Components {
			component := r.Catalog.Components[name]
			fmt.Fprintf(stdout, "  %-24s %s\n", name, component.Description)
		}
	}
	return nil
}

func install(r *registry.Registry, files []string, opts installFlags, stdout io.Writer) error {
	if opts.pkg == "" || !token.IsIdentifier(opts.pkg) {
		return fmt.Errorf("%q is not a valid Go package name", opts.pkg)
	}
	if err := os.MkdirAll(opts.dir, 0o755); err != nil {
		return err
	}
	written, skipped := 0, 0
	for _, name := range files {
		src, err := r.Source(name)
		if err != nil {
			return err
		}
		src = setPackage(src, opts.pkg)
		dst := filepath.Join(opts.dir, name)
		existing, readErr := os.ReadFile(dst)
		if readErr == nil && bytes.Equal(existing, src) {
			fmt.Fprintf(stdout, "  = %s\n", dst)
			skipped++
			continue
		}
		if readErr == nil && !opts.force {
			fmt.Fprintf(stdout, "  ! %s exists; use --force to replace it\n", dst)
			skipped++
			continue
		}
		if err := os.WriteFile(dst, src, 0o644); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "  + %s\n", dst)
		written++
	}
	fmt.Fprintf(stdout, "\n%d written, %d skipped in %s (package %s)\n",
		written, skipped, opts.dir, opts.pkg)
	return nil
}

func setPackage(src []byte, pkg string) []byte {
	if pkg == "ds" {
		return src
	}
	lines := strings.SplitAfter(string(src), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "package ") {
			lines[i] = "package " + pkg + "\n"
			return []byte(strings.Join(lines, ""))
		}
	}
	return src
}
