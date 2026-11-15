# AGENTS.md

Asuka DS is a source-owned component library for [Asuka](https://github.com/egoist/asuka):
typed theme tokens, 43 components and a CLI that copies plain Go files into
an app. Read [DESIGN.md](DESIGN.md) for the product's boundaries and decisions.

## Hard rules

- **`ds/` is the library.** It is `package ds` at
  `github.com/egoist/asuka-ds/ds` and also the sources the CLI ships. Keep the
  files flat inside `ds/` and keep `package ds`:
  `internal/catalog/components.json` names them by bare filename. The source
  directory is `catalog.SourceDir`; the generator, the drift test and the
  docs site all resolve paths through it.
- **The catalog is the source of truth.** A component exists when
  `internal/catalog/components.json` lists its kit and files. After changing
  a component, run `go generate ./internal/registry` (or
  `go run ./internal/registry/gen`); the drift test compares the embedded
  copies byte-for-byte with the repository's files.
- **Don't hand-edit `internal/registry/files/`.** It is generated.
- **Apps are entries, not packages.** `website/` and `examples/*/` are
  `package main` built with `asuka build <dir>`; they import the `ds` package
  and `internal/catalog`.
- **No generated runtime.** The CLI copies source files; it does not install
  or require a generated dependency.

## Layout

```
ds/                     library (package ds) and its tests
cmd/asuka-ds/           CLI: init / add / list / version
internal/catalog/       component catalog and kit dependency graph
internal/registry/      embedded source snapshot, generator, drift test
website/                documentation site (an Asuka app)
examples/               demo apps (showcase, form)
docs/                   platform matrix and research notes
artwork/screenshots/    screenshots
DESIGN.md               design decisions
rfcs/                   milestone plans
```

## Commands

```sh
go build ./... && go vet ./...
go test -count=1 ./...
gofmt -l .                          # no output
go generate ./internal/registry     # after editing a component
asuka vet ./...
asuka build --platform web website
asuka dev examples/showcase
asuka dev --platform web examples/form
```
