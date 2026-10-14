---
name: macos-settings-ui
description: >
  Build a proper macOS settings/preferences window with liquid glass support for macOS 26 (Tahoe).
  Use this skill whenever the user asks to create a settings window, preferences UI, settings view,
  or preferences pane for a macOS app. Also trigger when the user mentions "liquid glass settings",
  "NavigationSplitView settings", "grouped Form settings", "macOS settings layout", or wants to
  add/fix/rebuild a settings screen in any native macOS SwiftUI app. This includes phrases like
  "add settings to my app", "create a preferences window", "my settings look wrong", "fix my
  settings UI", or even just "I need a settings view". If the user is building a macOS app and
  mentions settings or preferences in any context, use this skill.
---

# macOS Settings UI with Liquid Glass

This skill produces a native macOS settings window that follows Apple's macOS 26 design guidelines. The result is a sidebar + detail NavigationSplitView with liquid glass window chrome, back/forward toolbar navigation, grouped Form sections with transparent backgrounds, and proper scroll edge effects.

## Architecture Overview

The settings UI is composed of 3 layers:

1. **Window Controller** (`SettingsWindowController.swift`) — An `NSWindowController` that creates the `NSWindow` programmatically with `.fullSizeContentView`. This is what gives the window rounded liquid glass corners. You cannot get this effect from a SwiftUI `Window` scene.

2. **Root View** (`SettingsView.swift`) — A `NavigationSplitView` with a sidebar list and detail pane. Includes back/forward navigation history in the toolbar, which also forces the creation of an `NSToolbar` (required for the liquid glass title bar treatment).

3. **Detail Panes** (one file per tab) — Each pane uses `Form { Section(...) { ... } }.formStyle(.grouped).scrollContentBackground(.hidden)`.

## Why NSWindowController Instead of SwiftUI Window Scene

SwiftUI's declarative `Window` scene does not expose the `NSWindow` style mask. The `.fullSizeContentView` flag must be set at window creation time for macOS 26 to render the liquid glass chrome (rounded corners, translucent sidebar, blurred title bar). Trying to inject it later via `NSViewRepresentable` is unreliable because SwiftUI resets the window's configuration.

The `NSWindowController` approach also lets you control the toolbar style, frame autosave, minimum size, and delegate lifecycle directly.

## Critical Modifiers

Every detail pane MUST have these three modifiers on its `Form`:

```swift
.formStyle(.grouped)
.scrollContentBackground(.hidden)
.contentMargins(.top, 8, for: .scrollContent)
```

- `.formStyle(.grouped)` — gives the native inset rounded-rect section appearance
- `.scrollContentBackground(.hidden)` — makes the form background transparent so the liquid glass window chrome shows through. Without this, you get an opaque white/dark background that breaks the glass effect.
- `.contentMargins(.top, 8)` — adds breathing room between the toolbar and the first section

The sidebar List MUST have:

```swift
.listStyle(.sidebar)
.scrollEdgeEffectStyleSoftIfAvailable()  // macOS 26 progressive blur at scroll edges
.navigationTitle("Settings")
```

The NavigationSplitView MUST have:

```swift
NavigationSplitView(columnVisibility: .constant(.all))  // sidebar always visible
// ...
.navigationTitle("Settings")
.navigationSplitViewStyle(.balanced)
```

## Reference Implementation

Read the reference files for complete, working code:

- `references/SettingsWindowController.swift` — The window controller (copy as-is, adapt the activation policy calls to your app)
- `references/SettingsView.swift` — The root view with sidebar, detail routing, and navigation history
- `references/ExampleDetailPane.swift` — A template detail pane showing common control patterns (Toggle, Picker, Slider, LabeledContent)

## Step-by-Step: Adding Settings to a New App

### 1. Create the Tab Enum

Define your settings categories. Each case needs a title and SF Symbol icon:

```swift
enum SettingsTab: String, CaseIterable, Identifiable {
    case general
    case appearance
    case about

    var id: Self { self }

    var title: String {
        switch self {
        case .general: "General"
        case .appearance: "Appearance"
        case .about: "About"
        }
    }

    var systemImage: String {
        switch self {
        case .general: "gearshape"
        case .appearance: "paintbrush"
        case .about: "info.circle"
        }
    }
}
```

### 2. Create the Window Controller

Copy `references/SettingsWindowController.swift` into your project. Adapt:
- The initial `contentRect` size (default `700x540` is good for most apps)
- The `minSize` (default `620x460`)
- The activation policy calls (`AppActivationPolicy.enter()/leave()`) — if your app isn't a menu-bar-only app, remove these

### 3. Create the Root Settings View

Copy `references/SettingsView.swift`. Adapt:
- The `SettingsTab` enum cases to match your categories
- The `SettingsDetailView` switch to return your panes

### 4. Create Detail Panes

For each tab, create a pane file. Use this template:

```swift
import SwiftUI

struct GeneralSettingsPane: View {
    @AppStorage("someKey") private var someValue = false

    var body: some View {
        Form {
            Section("Section Name") {
                Toggle("Toggle label", isOn: $someValue)
            }
        }
        .formStyle(.grouped)
        .scrollContentBackground(.hidden)
        .contentMargins(.top, 8, for: .scrollContent)
    }
}
```

### 5. Open Settings from Your App

From a menu bar, button, or anywhere:

```swift
SettingsWindowController.show(tab: .general)
```

Do NOT use a SwiftUI `Window` scene. Do NOT use `openWindow(id:)`. The window controller handles everything.

### 6. Remove Any SwiftUI Window Scene

If you previously had a `Window("Settings", id: "SETTINGS")` scene in your `App` struct, remove it entirely. The `SettingsWindowController` replaces it.

## Common Control Patterns Inside Form Sections

**Toggle with description:**
```swift
Toggle(isOn: $value) {
    VStack(alignment: .leading, spacing: 2) {
        Text("Primary label")
        Text("Description text explaining what this does.")
            .font(.caption)
            .foregroundStyle(.secondary)
    }
}
.toggleStyle(.switch)
```

**Picker (dropdown):**
```swift
Picker("Label", selection: $value) {
    Text("Option A").tag(OptionEnum.a)
    Text("Option B").tag(OptionEnum.b)
}
.pickerStyle(.menu)
```

**Picker (segmented):**
```swift
Picker("Label", selection: $value) {
    ForEach(SomeEnum.allCases) { item in
        Text(item.title).tag(item)
    }
}
.pickerStyle(.segmented)
```

**Slider with value label:**
```swift
LabeledContent("Size") {
    HStack(spacing: 12) {
        Slider(value: $size, in: 24...96, step: 2)
            .frame(width: 180)
        Text("\(Int(size)) pt")
            .monospacedDigit()
            .foregroundStyle(.secondary)
            .frame(width: 46, alignment: .trailing)
    }
}
```

**Button row:**
```swift
HStack(spacing: 8) {
    Button("Action") { doSomething() }
        .controlSize(.small)
    Button("Reset") { reset() }
        .controlSize(.small)
        .disabled(isDefault)
}
```

## Menu-Bar-Only Apps

If your app uses `NSApp.setActivationPolicy(.accessory)` (no Dock icon), you need activation policy management so the settings window brings the app to the foreground. The reference `SettingsWindowController` calls `AppActivationPolicy.enter()` on show and `.leave()` on close. Implement this as a simple reference-counting wrapper:

```swift
@MainActor
enum AppActivationPolicy {
    private static var count = 0

    static func enter() {
        count += 1
        NSApp.setActivationPolicy(.regular)
        NSApp.activate(ignoringOtherApps: true)
    }

    static func leave() {
        count = max(0, count - 1)
        guard count == 0 else { return }
        Task { @MainActor in
            NSApp.setActivationPolicy(.accessory)
        }
    }
}
```

If your app always shows in the Dock, remove the `AppActivationPolicy` calls from the window controller and just use `NSApp.activate(ignoringOtherApps: true)` in `showWindow`.

## macOS Version Compatibility

The `scrollEdgeEffectStyle(.soft)` API is macOS 26+ only. Always wrap it in an availability check:

```swift
private extension View {
    @ViewBuilder
    func scrollEdgeEffectStyleSoftIfAvailable() -> some View {
        if #available(macOS 26.0, *) {
            scrollEdgeEffectStyle(.soft, for: .all)
        } else {
            self
        }
    }
}
```

The rest of the pattern (NSWindowController with `.fullSizeContentView`, grouped Form, NavigationSplitView) works on macOS 14+.
