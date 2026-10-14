---
name: ui-prototyping
description: Explore multiple *divergent* UI directions for a screen as named Swift #Previews, remix the strongest elements into hybrids, fill them with lived-in sample content and edge-case states, then tune signature animations with a generated tuning panel. Produces real native SwiftUI you carry forward, not throwaway mockups. Use early — after new-app or at the start of any UI-heavy phase, before /apple:plan commits to a single layout. Based on Apple's WWDC method for prototyping with coding agents in Xcode: go wide → remix → make lived-in → tune.
allowed-tools: [Read, Write, Edit, Bash, Glob, Grep, AskUserQuestion]
last_verified: 2026-07-16
review_by: 2027-06-22
---

# UI Prototyping Skill

Task lists and specs commit you to **one** layout before you've seen the alternatives. The cost is quiet but large: you anchor on the first arbitrary structure an agent guessed, then spend the rest of the phase fighting feature creep around it. This skill front-loads *divergent exploration* — many real, named SwiftUI variations you flip between in Xcode's canvas — so the layout you carry into `/apple:plan` is one you **chose**, not one you defaulted into.

> **Agents are collaborators, not designers. You always have final say.** Go wide, remix, repeat.

## Why this exists (the trap it removes)

A vague prompt — *"make a UI for a book club"* — produces one arbitrary layout, silently guesses at features you never asked for (polling? a photo gallery?), and anchors you on a flawed start. By the time you've bent it to the features you *do* want, it looks clunky and inelegant. Three disciplines remove the trap, and this skill enforces all three:

| Discipline | What it means | What it prevents |
|---|---|---|
| **Specificity** | Bake the exact features in; nothing extra | Feature creep, arbitrary navigation elements |
| **Stylistic intent** | Name the mood — warm coffee-shop palette? paper & typography? clean editorial? | Point-of-view-less, generic layouts |
| **Multiplicity** | Ask for *many divergent* options at once | Anchoring on the first guess |

Early is your **only** cheap chance to explore wide. Spend it.

## Input

- `.planning/APP.md` — features, audience, platform. The feature list is exactly what you bake into every variation.
- Optional target screen (arg). If none, prototype the primary screen implied by `<mvp-features>`.
- **If `APP.md` is missing, that's fine** — elicit the essentials first via `AskUserQuestion` (this is the novice front door): the one screen, its 3–5 must-have features, the mood, and 0–2 reference apps. You do not need to know how to write the "good prompt" — the skill assembles it from your answers.

## The method — go wide → remix → make lived-in → tune

### Stage 1 — Go wide (divergent variations)

1. **Assemble the brief** (via `AskUserQuestion` if it isn't already pinned down in `APP.md`): the screen, the 3–5 must-have features (baked in — *this* is what stops feature creep), the mood cue, reference apps.
2. **Generate 6–10 genuinely divergent variations** of that ONE screen. Divergence is measured by *organizing metaphor*, not paint: vary tab vs. grid→detail vs. single-scroll vs. dashboard/standings, the navigation shape, typography (New York / serif vs. SF), density, and color/mood. Two variations that differ only in tint = one wasted slot.
3. **Each variation gets its own `#Preview` with a descriptive, memorable name** — `"Cozy"`, `"Editorial"`, `"Club Hub"`, `"Blueprint Atelier"`. The name is the handle you'll use when you remix. Put each in its own file (or a clearly separated struct) so the named previews render side by side in Xcode's canvas.
   ```swift
   #Preview("Cozy") { CozyHomeView(club: .sample) }
   #Preview("Club Hub") { ClubHubHomeView(club: .sample) }
   ```
4. **Only the requested features appear.** If the app needs no polling or gallery, no variation invents one.
5. **It compiles.** `xcodebuild build` (or the canvas) — you're carrying this forward, not screenshotting a dead mock.

Present the named set; the user flips between them and reacts. **Expect duds** ("…well, it was worth a shot") — range is the goal, not a batting average.

### Stage 2 — Remix + make it lived-in

1. **Remix.** The user names what they liked *by preview name and element* — "the standings board from *Club Hub* + the current-book cover from *Cozy*." Generate new **hybrid** variations from only those elements, each its own new named preview. Go wide → remix → repeat until one direction feels right.
2. **Lived-in content.** Empty scaffolding lies about how a screen feels. Generate **reusable sample models in their own file** (delegate to `generators/preview-data-generator`) so you can edit them, and make the content *plausible for the audience* — a book club's discussions are about books, not lorem ipsum. Plausible also means **messy, not idealized**: prototype with real data ranges, because the layout that flatters a hilly San Francisco route must also survive a pancake-flat New York one (WWDC23). Ideal placeholders hide exactly the cases that break a design.
3. **Edge-case previews.** Think the states through yourself; don't let the agent silently pick the happy path. Cover at minimum:
   - **Empty / blank-slate** — no meeting scheduled yet, zero items. Is there a call-to-action and an account/management entry, or does it just look broken?
   - **Unbounded growth** — many members, long message threads, dozens of past items, long titles. Truncate, or an expand control? (The leaderboard that pushes the discussion off-screen is the classic bug.)
   - **Long input** — does text truncate with `.help()`/detail access, or wrap to multiple lines?

   Delegate the state-matrix generation to the `swiftui-builder` agent (it emits a `#Preview` per state); **you** decide which states matter.

### Stage 3 — Tune key moments (a tuning panel)

Signature interactions — a cover-to-detail transition, a staggered list entrance — are where an app reads as considered or cheap. Don't hand-edit constants scattered across files; have the agent build a **tuning panel**.

1. **Be specific about what you're tuning.** Ease (duration) vs. spring (stiffness / damping / mass)? A single element, or a transition where views enter and leave the hierarchy?
2. **Break the animation into named phases** — "Phase 1: cover transitions to detail. Phase 2: rows stagger in." Phases give you and the agent a shared vocabulary; ask for controls per phase so you can inspect and retune one in isolation.
3. **Lay the panel out side-by-side on a wide window**, not a modal that obstructs the content — toggle a parameter and see the effect without context-switching. Ask for a resize control that moves the panel beside the UI on a larger canvas.
4. **Tuning panels generalize** beyond animation: swap app states, colors, fonts, visual offsets. Use `design/animation-patterns` for the curve knowledge; the panel is its runnable counterpart.
5. **Prototype beyond static layout.** Friction, inertia, device-motion effects (gyroscope parallax), and haptics are all promptable and panel-tunable the same way — and they're *feel*, only judgeable running on a device, not in a screenshot (WWDC26).
6. **Guard it behind DEBUG / a launch arg** so it never ships — same discipline as walkthrough's `-uiTestSeed`.
7. **On device early; review as a demo app.** Don't wait for "done" to leave the canvas — feel parameters (Digital Crown sensitivity, spring values) only tune while running live, so the earlier a prototype is on hardware the cheaper every iteration. And when it's time for a design review, bundle the competing prototypes into one on-device demo app rather than exporting stills: "no slide deck in the world can compare with one spectacular demo" (WWDC23).

## Output: `.planning/PROTOTYPE.md` + real Swift

Record the decisions so the next session — and `/apple:plan` — inherits them:

```markdown
# UI Prototype — [Screen], Phase [N]

## Variations explored (Stage 1)
| Preview name | Organizing idea | Verdict |
|--------------|-----------------|---------|
| Cozy | single-scroll, warm serif, current-book hero | ✅ base direction |
| Club Hub | tab bar + standings board | remix: take the standings board |
| Blueprint Atelier | grid → detail | ✗ too cold |

## Chosen direction & remix (Stage 2)
- Base: *Cozy*; grafted the standings board from *Club Hub*.
- Sample data: `PreviewData/BookClubSamples.swift` (editable, reused across previews).
- Edge cases covered: empty (added CTA + account entry), long member list (rank pinned + expand), long titles (truncate).

## Tuned moments (Stage 3)
- Cover→detail: spring(stiffness: …, damping: …); rows stagger 0.04s. Panel: `Debug/TransitionTuner.swift` (DEBUG only).

## Carry-forward
- Winning layout → `/apple:plan` `<flows>` + `<apple-patterns>`; keep the Swift files.
```

Plus the actual files in the project: the variation views, the sample-data file, and any tuning panel (DEBUG-guarded).

## Honest limits

- **Agents propose; you decide.** Don't delegate "which is best" — that judgment *is* the design work. A pile of variations with no chosen direction is not a result.
- **Divergent ≠ good.** Some variations will flop; that's the method working, not failing.
- **Nothing beats real users.** Lived-in sample content is a head start on the feedback stage, not a substitute for people using your app with their own content.
- **The tuning panel is scaffolding** — DEBUG-only, guarded, and removed or gated before release.

## Cadence & integration

- **Run early, once per signature screen** — after `/apple:new-app`, or at the start of a UI-heavy phase, *before* `/apple:plan` locks a layout. Exploration is cheapest before the "real" screen exists.
- **Feeds forward:** the winning variation is real SwiftUI — hand its structure to `/apple:plan` (`<flows>` / `<apple-patterns>`) and `/apple:build`; don't rebuild it.
- **Bookends with `flow-walkthrough` / `/apple:walkthrough`:** prototype the *screen* (diverge, choose, tune), then walkthrough the *flow* (verify transitions and dead-ends) once it's wired up. Explore → build → verify.
- **Complements `product/ux-spec`:** ux-spec formalizes *one* spec; this skill is how you find the one worth formalizing. Feed the winner in, or skip ux-spec for small apps.
- **Reuses:** `generators/preview-data-generator` (sample content), the `swiftui-builder` agent (edge-case `#Preview` matrix), `design/animation-patterns` (curve reference).
