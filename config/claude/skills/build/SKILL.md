---
name: macos-build
description: >
  Build a native macOS app using xcodebuild from the command line. Use this skill whenever the user
  asks to build, compile, or check if their macOS project compiles successfully. Also use it when
  the user asks to fix build errors, verify changes compile, or run a debug build. Trigger on
  phrases like "build the app", "does it compile", "run xcodebuild", "fix build errors", or even
  just "build".
---

# Build macOS App

Build any native macOS Xcode project from the command line using `xcodebuild`.

## Finding the Project

Before building, locate the Xcode project or workspace:

```bash
find . -maxdepth 2 -name "*.xcodeproj" -o -name "*.xcworkspace" | head -5
```

Then list available schemes:

```bash
xcodebuild -list -project "YourApp.xcodeproj" 2>/dev/null | grep -A 20 "Schemes:"
```

## Build Command

Use this command template, replacing the project path and scheme:

```bash
DEVELOPER_DIR=/Applications/Xcode.app/Contents/Developer xcodebuild build \
  -project "YourApp.xcodeproj" \
  -scheme "YourApp" \
  -configuration Debug \
  -destination "platform=macOS" \
  2>&1 | grep -E "(BUILD SUCCEEDED|BUILD FAILED|error:)" | head -20
```

For workspaces (projects with SPM dependencies or CocoaPods):

```bash
DEVELOPER_DIR=/Applications/Xcode.app/Contents/Developer xcodebuild build \
  -workspace "YourApp.xcworkspace" \
  -scheme "YourApp" \
  -configuration Debug \
  -destination "platform=macOS" \
  2>&1 | grep -E "(BUILD SUCCEEDED|BUILD FAILED|error:)" | head -20
```

## Interpreting Results

- **BUILD SUCCEEDED** -- the build passed, report success to the user.
- **BUILD FAILED** with `error:` lines -- read each error, identify the source file and line, and help the user fix them. After fixing, re-run the build to verify.
- If the output is empty or unclear, re-run without the grep filter to get full output for diagnosis.

## When to Build

- After making code changes, if the user asks to verify they compile
- When the user explicitly says "build", "compile", or "check if it builds"
- After fixing build errors, to confirm the fix worked

## Xcode Beta Toolchains

If the project targets a beta SDK (e.g., macOS 26 Tahoe), you may need to point to the beta Xcode:

```bash
DEVELOPER_DIR=/Applications/Xcode-beta.app/Contents/Developer xcodebuild build ...
```

Check which Xcode is available:

```bash
ls /Applications/ | grep -i xcode
```

## Common Build Failures

| Error | Fix |
|-------|-----|
| `no such module 'Sparkle'` | SPM dependency not resolved. Try `xcodebuild -resolvePackageDependencies` first |
| `no signing identity found` | Set `CODE_SIGN_IDENTITY=""` and `CODE_SIGNING_ALLOWED=NO` for command-line builds |
| `SDK "macosx" cannot be located` | Wrong `DEVELOPER_DIR`. Check Xcode installation path |
| `scheme not found` | Run `xcodebuild -list` to see available schemes |
