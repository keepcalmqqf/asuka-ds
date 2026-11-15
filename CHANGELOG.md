# Changelog

All notable changes to Asuka DS are recorded here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project
uses [semantic versioning](https://semver.org/).

## [0.1.0] - 2026-11-15

The first product release: a source-owned component library for Asuka, with
a CLI, an embedded registry and a documentation site.

### Added

- Typed token layer: `Palette`, `Radii`, `TypeScale`, `Motion` and
  `Elevation`, with WCAG 2.1 AA validation through `Theme.Validate`.
- Two built-in themes: `Native` follows platform semantic colors, system
  dark mode and accessibility contrast; `Neutral` ships shadcn's zinc light
  and dark palettes.
- Runtime theming: `SetTheme` and `SetAppearance` switch every mounted
  component through signals, without remounting.
- 43 components across Button, Form, Navigation, Overlay, Data display,
  Feedback and Typography, including drawn and native buttons, native text
  fields in a DS shell, command palette, combobox, data table, date range
  picker, toggle group, dialog, menu, popover, toast and tooltip.
- Product surfaces: `List`/`ListItem` settings rows, `Stat` metrics,
  `ProgressRing`, and a `ColorPicker` over the platform picker.
- `Pressable`, the shared interactive primitive that encodes disabled,
  selected and expanded state for drawn controls.
- `asuka-ds` CLI with `init`, `add`, `list` and `version`. Component sources
  are embedded, dependency-aware, offline and never overwrite edited files
  without `--force`.
- Embedded registry with a drift test that keeps shipped sources identical
  to the repository's components.
- Documentation site built with Asuka itself: guides, token reference, live
  examples and a responsive native menu.
- Demo apps for desktop and web (`demo`, `demo-form`), plus screenshots.
- Release pipeline: GoReleaser builds macOS, Linux and Windows binaries on a
  `v*` tag, with the version stamped through `-X main.Version`.
- Cross-platform CI: Windows, Linux and macOS builds, plus a platform
  support matrix and captured iOS and Android screenshots.
- 76 tests covering components, accessibility, theme tokens, catalog
  integrity and the CLI, with a subtest for every install path.

### Fixed

- Drawn `Button` and `Toggle` no longer tint their own container. On iOS a
  container's `Color` is its tint, so an accent background resolved to the
  label color and the control disappeared; labels and icons keep their own
  color instead.

### Known limitations

- Native controls cannot be deeply restyled; the DS themes their shell.
- Drawn controls have no native disabled or selected trait, so `Pressable`
  expresses those states through hints and event gating.
- Six-platform screenshot snapshots and device audits are still in progress.

[0.1.0]: https://github.com/egoist/asuka-ds/releases/tag/v0.1.0
