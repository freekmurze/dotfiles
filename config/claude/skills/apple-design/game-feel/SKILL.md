---
name: game-feel
description: Game feel ("juice") for Apple apps — celebration choreography, haptic vocabulary design, sound-effect layers, and the event×channel feedback audit. Use when big moments fall flat, when designing haptics or SFX, or when auditing whether every meaningful event actually reaches the user.
allowed-tools: [Read, Glob, Grep]
last_verified: 2026-07-18
review_by: 2027-06-22
os_version: iOS 26 / macOS 26
---

# Game Feel (Juice)

Makes an app's meaningful moments *land* — through coordinated motion, haptics, and sound. Animation technique lives in `design/animation-patterns`; confetti/badge generation lives in `generators/milestone-celebration`. This skill is the **discipline layer**: which events deserve feedback, on which channels, and how to verify nothing important is silent.

## When This Skill Activates

Use this skill when the user:
- Says the app feels **flat**, **lifeless**, or "lacks juice/polish/delight"
- Wants **celebrations** designed (winners, streaks, reveals, score changes)
- Asks to design or review **haptics** (vocabulary, Core Haptics, `.sensoryFeedback`)
- Asks about **sound effects**, buzzers, countdown ticks, or fanfares
- Wants a **feedback audit** — do the right events fire the right channels?
- Is building a **party/group game** where the phone is shared, passed, or watched at a distance

## Core Principles

### 1. Channels have different reach — match the channel to the audience

| Channel | Who perceives it |
|---------|------------------|
| Haptics | Only the person **holding** the device |
| Screen motion | Whoever is **looking** at the screen |
| Audio | The whole **room** |
| Torch / flash | The whole room, even in noise |

A single-user utility can lean on haptics. A shared-phone or room-facing app (party games, kitchen timers, presentations) **must** put its signature moments on a room-scale channel. The most common juice failure: all feedback is haptic, and the phone is lying on a table.

### 2. Feedback scales with rarity

Tick < action confirmed < success < round complete < **winner**. If the rarest event (winning) fires the same pattern as a routine one (round end), the app has no climax. Reserve the biggest pattern, the confetti burst, and the fanfare for the rarest event.

### 3. Escalation must agree across channels

If the timer's ring turns amber at 10s and red at 5s, the haptic tier and any audio tick must escalate at the *same* thresholds. Cross-channel disagreement reads as broken, not layered.

### 4. Single-owner rule

Exactly **one** code layer fires feedback for a given event. When an engine *and* a view-store both react to "music stopped", the user gets a double-buzz. Decide the owner (usually the layer closest to the state machine) and make the other layer silent.

### 5. Juice amplifies meaning — it never decorates

Every effect must be attached to an event the user cares about. Ambient/always-on effects (background particles, looping glows) are *texture*, not feedback — keep them static or slow, and never let them compete with event feedback.

### 6. Every channel needs an accessibility story

- Motion → gate large-scale motion on Reduce Motion, with a static equivalent
- Haptics → user-facing toggle (the system provides no global haptics setting for Core Haptics)
- Sound → mute toggle independent of system volume; a deliberate, documented silent-switch policy

## Decision Tree

```
What do you need?
│
├─ Which events should have feedback, and on which channels?
│  └─ → feedback-audit.md (event×channel matrix + three-lens audit)
│
├─ Design or review haptics (vocabulary, service architecture, APIs)
│  └─ → haptic-design.md
│
├─ Add a sound-effect layer (assets, latency, AVAudioSession coexistence)
│  └─ → sound-design.md
│
├─ HOW to animate a celebration/transition you've already decided on
│  └─ → design/animation-patterns (springs, transitions, symbol effects)
│
├─ Generate confetti / badge / milestone UI code
│  └─ → generators/milestone-celebration
│
└─ Animated SF Symbols for feedback moments
   └─ → design/sf-symbols (preset vocabulary)
```

## Review Checklist

When reviewing an app's game feel, verify:

- [ ] **Signature moments hit ≥2 channels** (e.g. winner = motion + haptic + sound)
- [ ] **Room-facing events use a room-scale channel** (audio, display-size motion, torch)
- [ ] **Rarity gradient exists** — the rarest event has visibly/audibly the biggest feedback
- [ ] **No channel disagreement** — escalation thresholds match across visual/haptic/audio
- [ ] **Single owner per event** — no double-fire from two layers
- [ ] **Score/number changes animate** (`.contentTransition(.numericText())`) — values never teleport
- [ ] **State changes transition** — no hard cuts between phases of a flow
- [ ] **Reduce Motion parity** — every gated effect has a static equivalent conveying the same information
- [ ] **User controls exist** — haptics toggle; mute toggle if there's an SFX layer
- [ ] **Loading/waiting states have anticipation** — determinate progress or a live effect, not a bare spinner at a climactic moment

## Reference Files

| File | Content |
|------|---------|
| [feedback-audit.md](feedback-audit.md) | Event×channel matrix method, three-lens audit (runnable subagent prompts), scoring rubric, output format |
| [haptic-design.md](haptic-design.md) | API decision table, semantic-vocabulary service architecture, pattern design rules, lifecycle, pitfalls |
| [sound-design.md](sound-design.md) | Minimal SFX kit, asset formats and latency, AVAudioSession coexistence with music, silent-switch policy, pitfalls |

## Related Skills

- `design/animation-patterns` — the *how* of every motion technique
- `design/sf-symbols` — symbol effect vocabulary for feedback glyphs
- `generators/milestone-celebration` — production confetti/badge/celebration code
- `ios/accessibility-audit` — Reduce Motion and assistive-tech verification

## References

- [Human Interface Guidelines — Playing haptics](https://developer.apple.com/design/human-interface-guidelines/playing-haptics)
- [Human Interface Guidelines — Playing audio](https://developer.apple.com/design/human-interface-guidelines/playing-audio)
- [Human Interface Guidelines — Motion](https://developer.apple.com/design/human-interface-guidelines/motion)
- [Core Haptics](https://developer.apple.com/documentation/corehaptics)
