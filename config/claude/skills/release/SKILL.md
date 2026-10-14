---
name: macos-release
description: >
  Release a native macOS app to GitHub with DMG packaging and Sparkle appcast updates. Use this
  skill whenever the user wants to publish a new version, create a release, ship an update, push
  a release to GitHub, or update the appcast. Also trigger when the user mentions DMG creation,
  Sparkle signing, notarization, archiving, or anything related to distributing a new version of
  their macOS app. This covers the full release pipeline: archive, notarize, export, create DMG,
  sign with Sparkle EdDSA, update appcast.xml, git push, and GitHub release creation.
---

# Release macOS App

This skill covers the full release pipeline for distributing a native macOS app outside the Mac App Store via GitHub Releases with Sparkle auto-update support.

## Release Pipeline Overview

```
Bump version → Archive → Notarize → Export → Create DMG → Sign DMG → Update appcast.xml → Git push → GitHub Release
```

## Prerequisites

The user needs these tools installed:

| Tool | Install | Purpose |
|------|---------|---------|
| `create-dmg` | `brew install create-dmg` | Creates the DMG installer |
| `gh` | `brew install gh` | Creates GitHub releases |
| `git` | Built-in | Pushes appcast changes |
| Sparkle `sign_update` | Built automatically when the project is built with Sparkle | EdDSA-signs the DMG |

The `sign_update` binary lives in DerivedData after building the project in Xcode:

```bash
find ~/Library/Developer/Xcode/DerivedData -name "sign_update" -type f 2>/dev/null | head -1
```

## Step-by-Step Release Guide

### 1. Bump Version Numbers

In Xcode, update:
- `MARKETING_VERSION` (e.g., `1.2`) -- the user-facing version
- `CURRENT_PROJECT_VERSION` (e.g., `5`) -- the build number (must be unique per release)

Or via command line:

```bash
# Check current values
grep -E "MARKETING_VERSION|CURRENT_PROJECT_VERSION" YourApp.xcodeproj/project.pbxproj | head -4
```

### 2. Archive in Xcode

Product > Archive. This creates a release build with the proper signing identity.

### 3. Notarize and Export

In the Archives organizer:
1. Select the archive > Distribute App
2. Choose "Direct Distribution" (or "Developer ID" for notarization)
3. Wait for notarization to complete
4. Export to `~/Downloads/YourApp.app`

### 4. Create DMG

```bash
create-dmg \
  --volname "YourApp" \
  --window-pos 200 120 \
  --window-size 660 400 \
  --icon-size 160 \
  --icon "YourApp.app" 180 170 \
  --app-drop-link 480 170 \
  --hide-extension "YourApp.app" \
  ~/Downloads/YourApp.dmg \
  ~/Downloads/YourApp.app
```

### 5. Sign DMG with Sparkle

```bash
/path/to/sign_update ~/Downloads/YourApp.dmg
```

This outputs the EdDSA signature and file length:

```
sparkle:edSignature="BASE64..." length="12345"
```

Save both values for the appcast.

### 6. Update appcast.xml

Add a new `<item>` at the top of the `<channel>` in your `appcast.xml`:

```xml
<item>
  <title>Version 1.2 (Build 5)</title>
  <pubDate>Mon, 26 May 2026 12:00:00 +0000</pubDate>
  <sparkle:version>5</sparkle:version>
  <sparkle:shortVersionString>1.2</sparkle:shortVersionString>
  <sparkle:minimumSystemVersion>14.0</sparkle:minimumSystemVersion>
  <description><![CDATA[<ul><li>Feature one</li><li>Bug fix two</li></ul>]]></description>
  <enclosure url="https://github.com/OWNER/REPO/releases/download/v1.2/YourApp.dmg"
             type="application/octet-stream"
             sparkle:edSignature="THE_SIGNATURE_FROM_STEP_5"
             length="THE_LENGTH_FROM_STEP_5" />
</item>
```

The `pubDate` should be RFC 2822 format. Generate it:

```bash
date -R
```

### 7. Commit and Push Appcast

```bash
git add appcast.xml
git commit -m "Release v1.2 appcast"
git push origin main
```

### 8. Create GitHub Release

```bash
gh release create v1.2 \
  ~/Downloads/YourApp.dmg \
  --title "v1.2" \
  --notes "- Feature one
- Bug fix two"
```

## Automating the Pipeline

For frequent releases, build a CLI tool that automates steps 4-8. See `references/release-pipeline.md` for a template Go CLI that handles DMG creation, signing, appcast updates, and GitHub release creation in one command.

The CLI should:
- Find `sign_update` in DerivedData automatically
- Read version/build from the exported app's Info.plist via `plutil`
- Parse and update appcast.xml (preserving existing entries)
- Interactively collect release notes
- Show a summary and ask for confirmation before proceeding
- Create the GitHub release with the DMG attached

## Release CLI Tool

This repo includes a Go CLI that automates steps 4-8 (DMG creation through GitHub release) in a single command. It reads configuration from a `release.json` file in your project root.

### Setup

1. Create a `release.json` in your project root (see `cli/release.example.json`):

```json
{
  "app_name": "MyApp",
  "github_repo": "owner/myapp"
}
```

Only `app_name` and `github_repo` are required. Everything else has sensible defaults:

| Field | Default | Description |
|-------|---------|-------------|
| `bundle_name` | `{app_name}.app` | The .app bundle filename |
| `dmg_name` | `{app_name}.dmg` | Output DMG filename |
| `git_branch` | `main` | Branch to push appcast to |
| `min_system_version` | `14.0` | Sparkle minimum macOS version |
| `appcast_file` | `appcast.xml` | Appcast filename in repo root |
| `derived_data_prefixes` | `["{app_name}-"]` | DerivedData prefixes to search for `sign_update` |

2. Run the CLI from your project directory:

```bash
go run github.com/fayazara/macos-app-skills/release/cli@latest
```

Or clone this repo and run locally:

```bash
go run ./release/cli
```

The CLI is interactive -- it prompts for release notes and asks for confirmation before proceeding. It must be run in a terminal the user can interact with.

### What the CLI Does

1. Finds `release.json` by walking up from the current directory
2. Checks that `create-dmg`, `gh`, `git`, and Sparkle's `sign_update` are available
3. Validates the exported app in `~/Downloads/` (reads version, build, Sparkle keys from Info.plist)
4. Warns if the build number already exists in the appcast
5. Collects release notes interactively (one bullet per line, empty line to finish)
6. Shows a release summary and asks for confirmation
7. Creates the DMG via `create-dmg`
8. Signs the DMG with Sparkle's `sign_update` (EdDSA)
9. Updates `appcast.xml` with the new release entry
10. Commits and pushes the appcast
11. Creates a GitHub release with the DMG attached

## Common Issues

| Problem | Solution |
|---------|----------|
| `sign_update` not found | Build the project in Xcode first so DerivedData has the Sparkle artifacts |
| `gh` auth failure | Run `gh auth login` |
| Duplicate build number in appcast | Bump `CURRENT_PROJECT_VERSION` before archiving |
| Notarization fails | Check signing identity, entitlements, and hardened runtime settings |
| DMG is too large | Check for debug symbols or unnecessary frameworks in the export |
| App won't update | Verify `SUFeedURL` points to the raw appcast URL, not the GitHub page URL |

## Appcast Hosting

The simplest hosting: commit `appcast.xml` to your GitHub repo and use the raw URL:

```
https://raw.githubusercontent.com/OWNER/REPO/main/appcast.xml
```

This must match `SUFeedURL` in your app's Info.plist.
