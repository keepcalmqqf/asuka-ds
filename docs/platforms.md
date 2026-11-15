# Platform support

Asuka DS is one Go codebase over Asuka's platform hosts. This page records
what was actually built and seen, on which OS and device, and what is still
unverified.

Verified on 2026-11-15:

| Platform | Toolchain | Verify | Status |
|---|---|---|---|
| macOS 27.2.1 | Xcode 27.2 | `asuka build --platform mac examples/showcase` | ✅ built, screenshots |
| iOS 26.5 | iPhone 17 Pro simulator | `asuka build --platform ios --simulator examples/showcase` | ✅ installed, screenshots |
| Android 16 | Pixel 9a emulator | `asuka dev --platform android --once examples/showcase` | ✅ installed, screenshots |
| Web | `asuka build --platform web` | static site + demo | ✅ built, screenshots |
| Windows | `GOOS=windows GOARCH=amd64 go build ./...` | CI cross-build | ⏳ build-only; no Windows host here |
| Linux | `GOOS=linux GOARCH=arm64 go build ./...` | CI cross-build | ⏳ build-only; no Linux host here |

## What was checked

- Native and Neutral themes, in light and dark where the host supports it.
- Drawn and platform buttons, text fields, switches, sliders, pickers,
  segmented controls, menus, sheets and snackbars.
- Live theme switching on iOS and Android without remounting.
- Semantic colors follow the platform: iOS blue, Android Material dynamic
  color (the Pixel 9a's purple), and macOS accent.
- `Theme.Validate` passes for both explicit palettes before a build.

## Screenshots

| Platform | Light | Dark |
|---|---|---|
| macOS | [native light](../artwork/screenshots/macos-native-light.png), [form](../artwork/screenshots/macos-form.png) | [native dark](../artwork/screenshots/macos-native-dark.png) |
| iOS | [native](../artwork/screenshots/ios-native-light.png) | [native](../artwork/screenshots/ios-native-dark.png), [neutral](../artwork/screenshots/ios-neutral-dark.png) |
| Android | [native](../artwork/screenshots/android-native-light.png) | [native](../artwork/screenshots/android-native-dark.png), [neutral](../artwork/screenshots/android-neutral-dark.png) |
| Web | [neutral](../artwork/screenshots/web-neutral-light.png), [focus ring](../artwork/screenshots/web-focus-ring.png) | [native](../artwork/screenshots/web-neutral-dark.png) |

Captures are device-resolution PNGs: iOS at 3x (1206×2622), Android at
1080×2424, macOS at 2x (1600×1264), and web at 2x. macOS images are taken by
launching the demo under each system appearance; an occluded window can return
its previous backing store to `screencapture`, so bring it to the front first.

## Known gaps

- **iOS**: a container's `Color` is its tint, and a semantic accent
  background resolves against that tint. Drawn components must color their
  label and icon, not the container; `Button` and `Toggle` follow that rule
  and have regression tests. Platform controls keep their own chrome.
- **Android**: Material dynamic color changes the accent by design. Use
  `Neutral()` when the brand color must be identical everywhere.
- **Web**: native pickers, menus and alerts fall back to the browser's own
  controls, so their placement and animation differ from a desktop app.
- **Windows and Linux**: the library cross-compiles in CI, but no host was
  available to run the demo, take screenshots or audit focus and IME
  behavior. Treat those platforms as preview until a real run is recorded.
- **Dynamic type**: text follows the platform's type scale, but a full audit
  at the largest accessibility sizes is still pending.

## Reproduce

```sh
# macOS
asuka build --platform mac examples/showcase
open "examples/showcase/dist/Asuka DS.app"

# iOS simulator
xcrun simctl boot "iPhone 17 Pro"
asuka build --platform ios --simulator examples/showcase
xcrun simctl install booted "examples/showcase/dist/Asuka DS.app"

# Android emulator
emulator -avd Pixel_9a &
asuka dev --platform android --once -d emulator-5554 examples/showcase

# Web
asuka build --platform web website
```

When the demo's layout changes, replace the screenshots at the same names so
the README, this page and the release notes keep pointing at truthful
images.
