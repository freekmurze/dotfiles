---
name: design
description: Design skills for Apple platform UI — Liquid Glass, animations, game feel (haptics, sound, celebrations), UI prototyping, UX writing, SF Symbols, and typography. Use when implementing design language features, adding juice/feedback, writing interface copy, or choosing type and iconography.
allowed-tools: [Read, Write, Edit, Glob, Grep, AskUserQuestion]
last_verified: 2026-07-18
review_by: 2027-06-22
os_version: iOS 27 / macOS 27
---

# Design Skills

Skills for implementing Apple's modern design systems across platforms.

## When This Skill Activates

Use this skill when the user:
- Asks about Liquid Glass design
- Wants to implement modern Apple UI effects
- Needs guidance on visual design patterns
- Asks about materials, transparency, or blur effects
- Wants to create fluid animations
- Asks about **spring**, **bounce**, or **snappy** animations
- Wants **PhaseAnimator** or **KeyframeAnimator** help
- Needs **view transitions**, **matched geometry**, or **hero transitions**
- Wants **SF Symbol effects** (bounce, pulse, wiggle, breathe)
- Asks about **animation completions** or **withAnimation**
- Says the app feels **flat** or wants **juice/game feel** (celebrations, haptics, sound effects)
- Wants a **feedback audit** — do events reach the user on the right channels?
- Needs UX structure help: navigation clarity, discoverability, information architecture
- Is writing or reviewing interface copy, alerts, or feature names
- Asks about SF Symbols usage/authoring, fonts, Dynamic Type, or type hierarchy

## Available Skills

### liquid-glass/
Comprehensive Liquid Glass implementation for iOS 26+, macOS 26+.
- SwiftUI `.glassEffect()` API
- AppKit `NSGlassEffectView`
- GlassEffectContainer patterns
- Morphing transitions
- Interactive effects
- Button styles

### animation-patterns/
SwiftUI animation patterns for iOS 13–18+.
- Spring configurations (3 API generations)
- PhaseAnimator and KeyframeAnimator (iOS 17+)
- View transitions, matched geometry, navigation transitions
- SF Symbol effects
- Animation completions, transactions, timing curves

### game-feel/
Game feel ("juice") as a discipline — the event×channel feedback audit, haptic vocabulary design (Core Haptics service architecture), and sound-effect layers that coexist with music. Routes to animation-patterns for motion technique and generators/milestone-celebration for generated code.

### ui-prototyping/
Explore divergent UI directions for a screen as named, runnable Swift `#Preview` variants — compare, remix, and tune before committing to one design.

### ux-writing/
Interface copy that works — voice/tone, the PACE framework, alert anatomy, feature naming, empty states, and the small-word edits that measurably improve UX.

### sf-symbols/
SF Symbols end-to-end — choosing/configuring system symbols (rendering modes, variable color/draw), authoring custom symbols that interpolate across weights, and the animation preset vocabulary.

### typography/
UI typography — text styles and Dynamic Type, the San Francisco family with its width axis, optical sizes, tracking vs kerning, and custom-font scaling obligations.

## Key Principles

### 1. Platform Consistency
- Follow Apple Human Interface Guidelines
- Use system-provided APIs
- Respect user appearance preferences

### 2. Performance
- Use GlassEffectContainer for multiple effects
- Limit number of glass effects per view
- Consider GPU resources

### 3. Visual Hierarchy
- Glass effects create depth and layering
- Use tints to indicate prominence
- Combine with appropriate shadows

## Reference Documentation

- [Applying Liquid Glass to custom views (SwiftUI)](https://developer.apple.com/documentation/swiftui/applying-liquid-glass-to-custom-views)
- [NSGlassEffectView (AppKit)](https://developer.apple.com/documentation/appkit/nsglasseffectview)
- Local captured docs (optional): if `~/Downloads/docs/SwiftUI-Implementing-Liquid-Glass-Design.md` or `~/Downloads/docs/AppKit-Implementing-Liquid-Glass-Design.md` exists, read it for extra detail; skip silently if absent.
