---
name: macos-patterns
description: >
  Essential native macOS development patterns that web developers don't know about. Use this skill
  whenever the user is building a macOS app and needs guidance on native patterns, or when they ask
  about menu bar apps, floating panels, window levels, keyboard shortcuts, file pickers, clipboard,
  drag and drop, screen capture, activation policy, Quick Look, launch at login, or any macOS-specific
  API. Also trigger when the user seems to be applying web development patterns to macOS (e.g., using
  z-index thinking for windows, expecting simple clipboard APIs, or not understanding focus/activation).
  This is the "how things actually work on macOS" reference that prevents the AI from generating
  confident but wrong code. Use this skill proactively whenever building any native macOS app.
---

# Native macOS Patterns for Web Developers

This is a reference guide for the macOS-specific patterns that have no web equivalent. When building a native macOS app, these are the things that trip up everyone coming from web development. The AI should consult this before generating macOS code to avoid confidently producing patterns that don't work.

## Menu Bar Apps

There are two approaches. Use `MenuBarExtra` for simple menus, `NSStatusItem` for full control.

### SwiftUI MenuBarExtra (simple)

```swift
@main
struct MyApp: App {
    var body: some Scene {
        MenuBarExtra("MyApp", image: "MenuBarIcon") {
            MenuBarView()
        }
    }
}
```

This creates a menu-bar-only app. The menu content is a standard SwiftUI view. For a popover-style menu bar app (richer UI than a plain menu), use the `.window` style:

```swift
MenuBarExtra("MyApp", image: "MenuBarIcon") {
    PopoverContentView()
}
.menuBarExtraStyle(.window)
```

### AppKit NSStatusItem (full control)

For custom menus, dynamic icons, or complex interactions:

```swift
let statusItem = NSStatusBar.system.statusItem(withLength: NSStatusItem.squareLength)
if let button = statusItem.button {
    let icon = NSImage(named: "MenuBarIcon")!
    icon.size = NSSize(width: 18, height: 18)
    icon.isTemplate = true  // CRITICAL: adapts to dark/light mode automatically
    button.image = icon
}

let menu = NSMenu()
menu.addItem(NSMenuItem(title: "Settings...", action: #selector(openSettings), keyEquivalent: ","))
statusItem.menu = menu
```

`isTemplate = true` is essential. Without it, your icon will be invisible in light mode or look wrong in dark mode.

### NSPopover for Rich Menu Bar Content

For a popover attached to the menu bar icon (like Bartender, iStatMenus):

```swift
let popover = NSPopover()
popover.contentSize = NSSize(width: 300, height: 400)
popover.behavior = .transient  // auto-closes when clicking outside
popover.contentViewController = NSHostingController(rootView: MyPopoverView())

// Show from the status item button:
if let button = statusItem.button {
    popover.show(relativeTo: button.bounds, of: button, preferredEdge: .minY)
}
```

## Activation Policy -- Dock Icon Toggling

A macOS app can dynamically show/hide its Dock icon and Cmd-Tab presence at runtime. Menu-bar-only apps start as `.accessory` (invisible in Dock) and temporarily become `.regular` when they open windows like Settings.

```swift
// At launch -- hide from Dock:
NSApp.setActivationPolicy(.accessory)

// When opening a window -- show in Dock:
NSApp.setActivationPolicy(.regular)
NSApp.activate(ignoringOtherApps: true)

// When closing the last window -- hide again:
NSApp.setActivationPolicy(.accessory)
```

Use reference counting if multiple windows can be open simultaneously:

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
        Task { @MainActor in NSApp.setActivationPolicy(.accessory) }
    }
}
```

## NSPanel vs NSWindow

`NSPanel` is a subclass of `NSWindow` for auxiliary content that should not steal focus. Use it for:

- Floating overlays (preview cards, recording indicators)
- Palettes and tool windows
- Inspectors
- Anything that should stay visible while the user works in another app

Key configuration:

```swift
let panel = NSPanel(
    contentRect: .zero,
    styleMask: [.borderless, .nonactivatingPanel],  // Does NOT steal focus
    backing: .buffered,
    defer: false
)
panel.isOpaque = false
panel.backgroundColor = .clear
panel.hasShadow = false
panel.hidesOnDeactivate = false     // Stay visible when app loses focus
panel.isFloatingPanel = true        // Float above normal windows
panel.ignoresMouseEvents = true     // Clicks pass through to windows below
panel.collectionBehavior = [
    .canJoinAllSpaces,              // Visible on all virtual desktops
    .fullScreenAuxiliary,           // Visible over fullscreen apps
]
```

To host SwiftUI content in the panel:

```swift
let hostingView = NSHostingView(rootView: MySwiftUIView())
panel.contentView = hostingView
```

Override focus behavior when needed:

```swift
class MyPanel: NSPanel {
    override var canBecomeKey: Bool { true }     // Can receive keyboard input
    override var canBecomeMain: Bool { false }   // Never becomes the "main" window
}
```

### Screen Capture Exclusion

If your app shows floating UI that shouldn't appear in screenshots/recordings:

```swift
panel.sharingType = .none  // Invisible to screen capture APIs
```

## Window Levels

macOS has a multi-tier window level system that controls where windows appear relative to the entire OS (not just your app):

```
normal (0)           -- standard app windows
floating (3)         -- tool palettes, floating panels
modalPanel (8)       -- modal dialogs
mainMenu (24)        -- menu bar
screenSaver (1000)   -- above everything except...
CGShieldingWindowLevel() -- above absolutely everything (kiosk mode)
```

Set via:

```swift
window.level = .floating
// or for extreme cases:
window.level = NSWindow.Level(rawValue: Int(CGShieldingWindowLevel()))
```

### Collection Behaviors

Control how windows interact with Spaces, fullscreen, and Cmd-Tab:

```swift
window.collectionBehavior = [
    .canJoinAllSpaces,       // Visible on all virtual desktops
    .fullScreenAuxiliary,    // Visible over fullscreen apps
    .stationary,             // Doesn't move with Space transitions
    .ignoresCycle,           // Hidden from Cmd-` window cycling
    .transient,              // Removed when app is hidden
]
```

## Screen Geometry

macOS uses **bottom-left origin** coordinates. The Y axis is flipped compared to the web.

```swift
let screen = NSScreen.main!

// Full screen rectangle (includes menu bar and Dock area):
screen.frame  // e.g., (0, 0, 1728, 1117)

// Usable area (excludes menu bar and Dock):
screen.visibleFrame  // e.g., (0, 0, 1728, 1055)

// Notch detection (MacBook with notch has top safe area):
screen.safeAreaInsets.top > 0  // true on notch Macs
```

Use `frame` when positioning over the menu bar (notch overlays). Use `visibleFrame` for normal window placement.

### Quartz vs AppKit Y-Axis

Core Graphics / Quartz uses **top-left origin**. AppKit uses **bottom-left origin**. When converting between the two:

```swift
let desktopFrame = NSScreen.screens.reduce(CGRect.null) { $0.union($1.frame) }
let appKitY = desktopFrame.maxY - quartzY - height  // Quartz → AppKit
let quartzY = desktopFrame.maxY - appKitY - height  // AppKit → Quartz
```

### Multi-Monitor

Never assume a single screen. Always handle the case where `NSScreen.main` is not the only display:

```swift
// Find the screen containing the mouse pointer:
let mouseLocation = NSEvent.mouseLocation
let screen = NSScreen.screens.first { $0.frame.contains(mouseLocation) }

// Find the screen containing a specific window:
let screen = window.screen
```

## Keyboard Shortcuts

There are 3 tiers, each for different use cases.

### Tier 1: SwiftUI keyboard shortcuts (in-app, when focused)

```swift
Button("Settings") { openSettings() }
    .keyboardShortcut(",", modifiers: [.command])  // Cmd+,
```

### Tier 2: NSEvent monitors (app-wide or global)

```swift
// Local: fires only when YOUR app is active
let monitor = NSEvent.addLocalMonitorForEvents(matching: .keyDown) { event in
    if event.keyCode == 53 { /* Escape */ return nil /* consume */ }
    return event  // pass through
}

// Global: fires even when ANOTHER app is active
let monitor = NSEvent.addGlobalMonitorForEvents(matching: .keyDown) { event in
    // Cannot consume the event -- read-only
}

// CRITICAL: Always remove monitors when done
NSEvent.removeMonitor(monitor)
```

Returning `nil` from a local monitor **consumes** the event (stops propagation). Global monitors **cannot** consume events.

### Tier 3: Carbon hotkeys (system-wide, works when app is not focused)

The only way to register true global keyboard shortcuts. Uses the Carbon API (1990s-era, still not deprecated):

```swift
import Carbon.HIToolbox

var hotKeyRef: EventHotKeyRef?
let hotKeyID = EventHotKeyID(signature: OSType(0x4D594150), id: 1) // "MYAP"

RegisterEventHotKey(
    UInt32(kVK_ANSI_1),      // Key code
    UInt32(optionKey),         // Modifiers
    hotKeyID,
    GetApplicationEventTarget(),
    0,
    &hotKeyRef
)
```

The callback runs on an arbitrary thread -- always dispatch to main:

```swift
DispatchQueue.main.async { self.handleHotKey() }
```

## File Pickers

### Open (select files/directories)

```swift
let panel = NSOpenPanel()
panel.allowedContentTypes = [.image, .movie]  // UTType filter
panel.allowsMultipleSelection = true
panel.canChooseDirectories = false
panel.canChooseFiles = true
panel.canCreateDirectories = true
panel.title = "Choose Images"

// Modal (blocks the thread):
if panel.runModal() == .OK {
    let urls = panel.urls
}

// Or sheet-modal (attached to a window, async):
panel.beginSheetModal(for: window) { response in
    guard response == .OK else { return }
    let urls = panel.urls
}
```

### Save

```swift
let panel = NSSavePanel()
panel.allowedContentTypes = [.png]
panel.nameFieldStringValue = "screenshot.png"
panel.canCreateDirectories = true

panel.begin { response in
    guard response == .OK, let url = panel.url else { return }
    // write file to url
}
```

### Directory picker

```swift
let panel = NSOpenPanel()
panel.canChooseFiles = false
panel.canChooseDirectories = true
panel.directoryURL = URL(fileURLWithPath: NSHomeDirectory())
```

## Clipboard / Pasteboard

macOS's pasteboard is fundamentally different from the web's `navigator.clipboard`. It is a multi-item, multi-type container.

```swift
let pasteboard = NSPasteboard.general

// CRITICAL: You MUST clear before writing. Forgetting this is a common bug.
pasteboard.clearContents()

// Write text:
pasteboard.setString("hello", forType: .string)

// Write image data:
let pngData = try Data(contentsOf: imageURL)
pasteboard.setData(pngData, forType: .png)

// Write a file URL (for Finder paste):
pasteboard.writeObjects([url as NSURL])

// Read (check what types are available):
if let string = pasteboard.string(forType: .string) { ... }
if let data = pasteboard.data(forType: .png) { ... }
```

A single pasteboard item can advertise multiple types. When reading, check types in priority order.

## Drag and Drop

### SwiftUI (simple)

```swift
// Make something draggable:
Image(nsImage: image)
    .draggable(fileURL) {
        // Custom drag preview
        Image(nsImage: image).frame(width: 100, height: 75)
    }

// Accept drops:
.dropDestination(for: URL.self) { urls, location in
    handleDroppedFiles(urls)
    return true
}
```

### AppKit (full control)

Register as a drop target:

```swift
class MyView: NSView {
    override init(frame: NSRect) {
        super.init(frame: frame)
        registerForDraggedTypes([.fileURL, .png, .tiff])
    }

    override func draggingEntered(_ sender: NSDraggingInfo) -> NSDragOperation {
        .copy  // Show the green + cursor
    }

    override func performDragOperation(_ sender: NSDraggingInfo) -> Bool {
        guard let urls = sender.draggingPasteboard.readObjects(forClasses: [NSURL.self]) as? [URL] else {
            return false
        }
        handleFiles(urls)
        return true
    }
}
```

## NavigationSplitView + Inspector Layout

The native macOS pattern for a sidebar + detail + inspector layout:

```swift
struct ContentView: View {
    @State private var selection: Item?
    @State private var showInspector = true

    var body: some View {
        NavigationSplitView {
            SidebarView(selection: $selection)
                .navigationSplitViewColumnWidth(min: 200, ideal: 220, max: 300)
        } detail: {
            DetailView(item: selection)
        }
        .inspector(isPresented: $showInspector) {
            InspectorView(item: selection)
                .inspectorColumnWidth(min: 250, ideal: 280, max: 400)
        }
        .toolbar {
            ToolbarItem {
                Button { showInspector.toggle() } label: {
                    Image(systemName: "sidebar.right")
                }
            }
        }
    }
}
```

The `.inspector()` modifier creates a native right-side panel that slides in/out, automatically manages layout, and integrates with the window's toolbar.

## Launch at Login

Use `SMAppService` (macOS 13+):

```swift
import ServiceManagement

// Check status:
switch SMAppService.mainApp.status {
case .enabled: /* running at login */
case .requiresApproval: /* registered but user must approve in System Settings */
case .notRegistered, .notFound: /* not registered */
}

// Register:
try SMAppService.mainApp.register()

// Unregister:
try SMAppService.mainApp.unregister()
```

The `requiresApproval` state is unique to macOS -- the app has asked to launch at login, but the user must manually approve it in System Settings > General > Login Items.

Always disable in debug builds to avoid polluting the login item list during development.

## Quick Look Preview

Show a system Quick Look panel for any file (images, PDFs, videos, documents):

```swift
import QuickLookUI

class PreviewPresenter: NSObject, QLPreviewPanelDataSource {
    var url: URL?

    func show(url: URL) {
        self.url = url
        guard let panel = QLPreviewPanel.shared() else { return }
        NSApp.activate()  // MUST activate first or panel opens behind other windows
        panel.dataSource = self
        panel.reloadData()
        panel.makeKeyAndOrderFront(nil)
    }

    func numberOfPreviewItems(in panel: QLPreviewPanel!) -> Int { url != nil ? 1 : 0 }

    func previewPanel(_ panel: QLPreviewPanel!, previewItemAt index: Int) -> (any QLPreviewItem)! {
        url as? NSURL
    }
}
```

## NSWorkspace -- OS Integration

```swift
// Open URL in default browser:
NSWorkspace.shared.open(URL(string: "https://example.com")!)

// Reveal file in Finder (selects it):
NSWorkspace.shared.activateFileViewerSelecting([fileURL])

// Get the frontmost application:
let app = NSWorkspace.shared.frontmostApplication
let name = app?.localizedName  // e.g., "Safari"

// Check accessibility preferences:
NSWorkspace.shared.accessibilityDisplayShouldReduceMotion  // Respect "Reduce Motion"
NSWorkspace.shared.accessibilityDisplayShouldReduceTransparency
```

## UserDefaults + @AppStorage

### Programmatic access

```swift
// Write:
UserDefaults.standard.set(true, forKey: "autoSave")
UserDefaults.standard.set(0.8, forKey: "quality")

// Read (returns false/0/nil if key doesn't exist):
let autoSave = UserDefaults.standard.bool(forKey: "autoSave")
let quality = UserDefaults.standard.double(forKey: "quality")
```

### Reactive SwiftUI binding

```swift
struct SettingsView: View {
    @AppStorage("autoSave") private var autoSave = false
    @AppStorage("quality") private var quality = 0.8

    var body: some View {
        Toggle("Auto Save", isOn: $autoSave)      // Auto-persisted
        Slider(value: $quality, in: 0...1)         // Auto-persisted
    }
}
```

`@AppStorage` automatically reads from and writes to `UserDefaults`, and triggers SwiftUI view updates when the value changes.

## ScreenCaptureKit

Capture screen content at native Retina resolution:

```swift
import ScreenCaptureKit

let content = try await SCShareableContent.excludingDesktopWindows(false, onScreenWindowsOnly: false)
guard let display = content.displays.first else { return }

let filter = SCContentFilter(display: display, excludingWindows: [])
let config = SCStreamConfiguration()
let scale = CGFloat(filter.pointPixelScale)
config.width = Int(CGFloat(display.width) * scale)   // Points → pixels
config.height = Int(CGFloat(display.height) * scale)
config.showsCursor = false

let image = try await SCScreenshotManager.captureImage(contentFilter: filter, configuration: config)
```

The points-to-pixels conversion is critical. ScreenCaptureKit works in points, but output dimensions must be in pixels for Retina resolution.

## Common Mistakes Web Devs Make

| What they try | Why it fails | What to do instead |
|---|---|---|
| SwiftUI `Window` scene for floating UI | Steals focus, shows in Dock, no transparency | Use `NSPanel` with `.nonactivatingPanel` |
| `z-index` thinking for window ordering | macOS uses discrete window levels, not a flat stack | Set `window.level` to `.floating`, `.screenSaver`, etc. |
| `window.innerHeight` for positioning | Doesn't account for menu bar, Dock, or notch | Use `NSScreen.visibleFrame` or `.frame` depending on context |
| `navigator.clipboard.writeText()` | macOS requires `clearContents()` first | Always call `pasteboard.clearContents()` before writing |
| `addEventListener('keydown')` for global shortcuts | Only works when app is focused | Use Carbon `RegisterEventHotKey` for system-wide |
| `<input type="file">` mental model | macOS has modal, sheet-modal, and async file pickers | Use `NSOpenPanel` with the right presentation mode |
| Single-monitor assumptions | macOS users commonly have 2-3 displays | Always use `NSScreen.screens` and find the right one |
| CSS animation for everything | macOS has spring physics, reduced motion, per-window animation | Use SwiftUI `.animation(.spring(...))` and check `accessibilityDisplayShouldReduceMotion` |
