# Feedback Audit

A systematic method for answering: *does every meaningful event actually reach its audience?* Produces an event×channel matrix, a scored report, and a tiered fix list.

## Step 1 — Build the event×channel matrix

List every meaningful event in the app (per mode/flow if the app has modes). For each, record what fires on each channel — **from the code, not from memory** — and who must perceive the event.

| Event | Audience | Visual motion | Haptic | Audio | Verdict |
|-------|----------|--------------|--------|-------|---------|
| Round starts | holder | — | — | — | 🔴 silent everywhere |
| Correct answer | room | — | `.correct` | — | 🟠 room can't perceive it |
| Timer final 5s | room | ring → red | `.timerUrgent` | — | 🟡 visual+haptic, no audio |
| Music stops (freeze) | room | ember flash + torch | `.roundEnd` | (music cut = the cue) | ✅ multi-channel |
| Winner declared | room | — | `.roundEnd` | — | 🔴 climax = routine pattern |

Verdict rules:
- **Signature events need ≥2 channels**; room-audience events need ≥1 room-scale channel (audio, display-size motion, torch)
- 🔴 = signature/room event with zero or wrong-reach feedback, or climax reusing a routine pattern
- 🟠 = meaningful event on one channel only, or channel escalations disagreeing
- 🟡 = adequate but flat (static where motion is expected; generic where signature is expected)
- ✅ = channels match audience and rarity

## Step 2 — The three-lens code sweep

Run three *parallel, read-only* audits (subagents if available; sequential passes otherwise). Facts with `file:line` evidence only — no quality judgments until synthesis. Explicit "NOT FOUND" beats silence.

**Lens 1 — Haptics + sound inventory.** Find the haptics mechanism (Core Haptics service? `.sensoryFeedback`? scattered generators?) and describe its architecture: central or ad hoc, lifecycle handling, injection seams, capability gating. Map every fire-site to its event, per mode. Then find every sound-effect mechanism and bundled audio asset (search for player APIs *and* `.caf/.wav/.m4a` files); report the AVAudioSession category/options and interruption handling; report user-facing toggles. List events that fire nothing.

**Lens 2 — Motion inventory.** Find every animation API use (`withAnimation`, `.animation`, `.transition`, `.contentTransition`, `phaseAnimator`, `keyframeAnimator`, `matchedGeometryEffect`, `.symbolEffect`, `TimelineView`, custom `Animatable`). Summarize the vocabulary: which curves/durations, shared tokens or scattered literals. Inspect each celebration moment (winner, score change, reveal, elimination): what actually moves? Check every timer/progress visualization for stepping vs. interpolation. Verify Reduce Motion gates on all large-scale motion and list ungated cases. Check loading states at climactic moments for bare spinners.

**Lens 3 — Cohesion + reach.** Which screens carry the app's design system vs. stock chrome (stock empty-state views, default button styles, system-typography islands)? For each *room-facing* screen: are the primary readouts at display scale? Cross-check channel escalation thresholds (visual color changes vs. haptic tiers vs. audio ticks) for agreement. Flag double-fire owners (two layers reacting to one event).

## Step 3 — Score and report

Rate each dimension; anchor every score to matrix rows and lens evidence:

| Dimension | Poor (1–3) | Adequate (5–6) | Excellent (8–10) |
|-----------|-----------|----------------|------------------|
| Feedback coverage | Signature events silent | Most events on one channel | Matrix full; channels match audience |
| Rarity gradient | Winner = routine pattern | Some differentiation | Climaxes unmistakably biggest |
| Motion craft | Hard cuts, jumping numbers | Standard transitions | Choreographed, tokenized, interruptible |
| Haptic design | Scattered presets | Central service, generic patterns | Semantic vocabulary, tuned tiers |
| Sound design | None where room needs it | Minimal kit, session-correct | Rarity-scaled kit, deliberate policies |
| Accessibility parity | Ungated motion, no toggles | Reduce Motion gated | Static equivalents + haptic/mute toggles |

## Output Format

```
## Game-Feel Audit — [App]

### Scorecard
| Dimension | Score | Evidence |

### Event×Channel Matrix
[the matrix, worst rows first]

### Findings
🔴 Critical: [silent signature moments, wrong-reach feedback]
🟠 High: [single-channel signature events, double-fires, escalation disagreement]
🟡 Medium: [flat-but-present feedback, stock islands]
🟢 Low: [tokenization, minor polish]

### Fix tiers
Tier 1 — wiring (hours): missing haptic events, numericText, symbol effects
Tier 2 — celebration moments (day): confetti burst, reveals, transitions
Tier 3 — sound layer (day + assets + policy decisions): see sound-design.md
```

Tier 3 always requires **user decisions** (silent-switch policy, asset sourcing) — surface them, never assume.

## Honest limits

A code audit verifies *existence and wiring*, not *feel*. Intensity curves, animation timing, loudness balance, and torch visibility are device-only judgments — end every audit report by scheduling a device pass. The Simulator plays no haptics and its audio timing lies.
