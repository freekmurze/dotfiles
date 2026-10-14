---
name: macos-auto-update
description: >
  Add Sparkle auto-update support to a native macOS app. Use this skill whenever the user wants to
  add auto-updates, implement Sparkle, set up an appcast feed, add "Check for Updates" to their app,
  or integrate automatic update checking. Also trigger when the user mentions Sparkle, appcast,
  SUFeedURL, EdDSA signing, or update notifications. This covers the full setup: adding the Sparkle
  SPM dependency, creating the UpdaterManager singleton, wiring it into the app delegate and UI,
  configuring Info.plist, and generating the initial appcast.xml.
---

# Sparkle Auto-Update for macOS Apps

This skill adds [Sparkle](https://sparkle-project.org/) auto-update support to a native macOS app. Sparkle is the standard open-source framework for macOS app updates outside the Mac App Store.

## Overview

The implementation has 4 parts:

1. **SPM dependency** -- add the Sparkle package to the Xcode project
2. **UpdaterManager.swift** -- a singleton that wraps `SPUStandardUpdaterController`
3. **Info.plist keys** -- `SUFeedURL` and `SUPublicEDKey`
4. **UI integration** -- "Check for Updates" button in settings/menu bar

## Step 1: Add Sparkle via SPM

In Xcode: File > Add Package Dependencies > enter:

```
https://github.com/sparkle-project/Sparkle
```

Use the "Up to Next Major Version" rule with `2.0.0`. Add the `Sparkle` framework to your app target.

Or add it to `Package.swift` if your project uses one:

```swift
.package(url: "https://github.com/sparkle-project/Sparkle", from: "2.0.0")
```

## Step 2: Create UpdaterManager.swift

Copy `references/UpdaterManager.swift` into your project. This is a singleton that:

- Creates `SPUStandardUpdaterController` early (before `applicationDidFinishLaunching` returns)
- Publishes `canCheckForUpdates` for UI binding
- Exposes `automaticallyChecksForUpdates` toggle
- Skips all update logic in DEBUG builds (so you don't get update prompts during development)
- For menu-bar-only apps: temporarily switches to `.regular` activation policy before showing the update window

The key design decisions in this file:

- **`startingUpdater: false`** in the initializer, then calling `start()` explicitly in `applicationDidFinishLaunching`. This gives you control over timing.
- **DEBUG guards** on `start()` and `checkForUpdates()`. Sparkle should never run in debug builds -- it would try to update your debug app with a release build.
- **`ObservableObject` with `@Published`** (not `@Observable`) because we need the Combine `publisher(for:)` bridge from Sparkle's KVO.

## Step 3: Configure Info.plist

Add these keys to your app's `Info.plist`:

```xml
<key>SUFeedURL</key>
<string>https://raw.githubusercontent.com/OWNER/REPO/main/appcast.xml</string>

<key>SUPublicEDKey</key>
<string>YOUR_PUBLIC_EDDSA_KEY_HERE</string>

<key>SUEnableInstallerLauncherService</key>
<true/>
```

### Generating EdDSA Keys

Sparkle uses EdDSA (Ed25519) signing. Generate a keypair:

```bash
# Find generate_keys in your DerivedData after building the project with Sparkle
find ~/Library/Developer/Xcode/DerivedData -name "generate_keys" -type f 2>/dev/null | head -1
```

Run it:

```bash
/path/to/generate_keys
```

This prints the public key and stores the private key in your Keychain. Put the public key in `SUPublicEDKey` in Info.plist. The private key stays in Keychain and is used by `sign_update` during release.

## Step 4: Wire Into App

### App Delegate

```swift
final class AppDelegate: NSObject, NSApplicationDelegate {
    private let updaterManager = UpdaterManager.shared

    func applicationDidFinishLaunching(_ notification: Notification) {
        updaterManager.start()
    }
}
```

The `UpdaterManager.shared` property must be accessed early so the `SPUStandardUpdaterController` is created before the app finishes launching. Referencing it in the `AppDelegate` property ensures this.

### Settings UI (About Pane)

```swift
struct AboutSettingsPane: View {
    @ObservedObject private var updaterManager = UpdaterManager.shared

    var body: some View {
        Form {
            Section("Updates") {
                Toggle(isOn: Binding(
                    get: { updaterManager.automaticallyChecksForUpdates },
                    set: { updaterManager.automaticallyChecksForUpdates = $0 }
                )) {
                    Text("Automatically check for updates")
                }

                Button("Check for Updates...") {
                    updaterManager.checkForUpdates()
                }
                .disabled(!updaterManager.canCheckForUpdates)
            }
        }
        .formStyle(.grouped)
        .scrollContentBackground(.hidden)
    }
}
```

### Menu Bar (optional)

```swift
Button {
    updaterManager.checkForUpdates()
} label: {
    Label("Check for Updates...", systemImage: "arrow.down.circle")
}
.disabled(!updaterManager.canCheckForUpdates)
```

## Step 5: Create Initial Appcast

Create an `appcast.xml` at the root of your repo. It starts empty and gets populated by your release process:

```xml
<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0" xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle" xmlns:dc="http://purl.org/dc/elements/1.1/">
  <channel>
    <title>YourApp</title>
    <description>Most recent changes for YourApp.</description>
    <language>en</language>
  </channel>
</rss>
```

Host this file on GitHub (raw URL) or any static file host. The URL must match `SUFeedURL` in Info.plist.

## Menu-Bar-Only Apps

If your app runs as `.accessory` (no Dock icon), Sparkle's update window won't appear unless you temporarily switch to `.regular`. The reference `UpdaterManager` handles this in `checkForUpdates()`:

```swift
func checkForUpdates() {
    NSApp.setActivationPolicy(.regular)
    NSApp.activate(ignoringOtherApps: true)
    controller.checkForUpdates(nil)
}
```

The app reverts to `.accessory` when the update window closes (handled by your existing activation policy manager).

## Hardened Runtime Entitlements

If your app uses Hardened Runtime (required for notarization), no special Sparkle entitlements are needed. Sparkle 2.x works with the standard hardened runtime configuration.

## Appcast Item Format

Each release in the appcast looks like this (for reference when building release tooling):

```xml
<item>
  <title>Version 1.2 (Build 5)</title>
  <pubDate>Mon, 26 May 2026 12:00:00 +0000</pubDate>
  <sparkle:version>5</sparkle:version>
  <sparkle:shortVersionString>1.2</sparkle:shortVersionString>
  <sparkle:minimumSystemVersion>14.0</sparkle:minimumSystemVersion>
  <description><![CDATA[<ul><li>New feature</li><li>Bug fix</li></ul>]]></description>
  <enclosure url="https://github.com/OWNER/REPO/releases/download/v1.2/YourApp.dmg"
             type="application/octet-stream"
             sparkle:edSignature="BASE64_EDDSA_SIGNATURE"
             length="FILE_SIZE_BYTES" />
</item>
```

- `sparkle:version` = `CFBundleVersion` (build number)
- `sparkle:shortVersionString` = `CFBundleShortVersionString` (marketing version)
- `sparkle:edSignature` = output of `sign_update YourApp.dmg`
- `length` = file size in bytes
