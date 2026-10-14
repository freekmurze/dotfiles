---
name: typography
description: UI typography on Apple platforms — text styles and Dynamic Type, the San Francisco family (Pro/Rounded/Compact/Mono/New York + width axis), optical sizes, tracking vs kerning, leading adjustments, and custom-font scaling. Use when choosing fonts, building type hierarchy, fixing truncation/legibility, or making custom fonts respect Dynamic Type.
allowed-tools: [Read, Write, Edit, Glob, Grep]
last_verified: 2026-07-16
review_by: 2027-06-22
---

# Typography

Type is most of your UI. The system does the hard parts — optical sizes, tracking tables,
Dynamic Type — *if* you use its APIs. This skill covers when to trust the system, how to build
hierarchy deliberately, and what custom fonts owe you back.

## When This Skill Activates

- Choosing fonts or building a type hierarchy for a screen or app
- Text truncating, cramping, or breaking under larger Dynamic Type sizes
- Adopting a custom/brand font (and keeping accessibility)
- Display typesetting: hero numbers, stats, editorial headlines

## Rule 1: text styles first

Predefined styles (Large Title → Caption 2) are weight + size + leading combos with Dynamic
Type support built in — and different styles scale *differently* (body grows more than
footnote), which manual font sizes can't replicate.

- Build hierarchy from 2–3 text styles + emphasized variants (bold trait — the actual weight
  varies per style: some go medium→semibold, others bold→heavy) before reaching for new fonts.
- Semantic colors + text styles together give dark mode, contrast, and Dynamic Type for free.
- macOS supports text styles (without Dynamic Type); design Mac type at 100% — no iOS-style
  scaling assumptions.

## Rule 2: let optical sizes and tracking work

- San Francisco blends Text→Display designs **continuously between 17 and 28pt** (below: sturdier
  letterforms, looser tracking for legibility; above: tighter, more refined). System font APIs do
  this automatically — hardcoded single-cut fonts don't.
- **Tracking, not kerning**, for letterspacing — tracking is size-specific and lets the OS
  disable clashing features (ligatures). Override system tracking only in exceptional cases.
- Truncation pressure? Use `allowsTightening` (default-tightening-for-truncation) instead of
  manually squeezing — and prefer wrapping to truncating (line limit 0 on labels that matter).
- Line height moves only via leading traits: tight = −2pt, loose = +2pt (±1pt on watchOS); the
  system adds leading automatically for tall scripts (Arabic, Devanagari).

## The Dynamic Type ladder (WWDC24 10074)

12 sizes: 7 default + 5 accessibility (AX1–AX5). `.body` runs 17pt at default up to 53pt
at AX5 — roughly 3× taller, and layouts must expect that. Escalate in order:

1. **Text styles** — `.font(.title)` / `preferredFont(forTextStyle:)` +
   adjusts-for-content-size, line count 0 so text wraps instead of truncating.
2. **`@ScaledMetric`** for the non-text riding alongside (icon frames, spacing); SF
   Symbols scale via `UIImage.SymbolConfiguration(textStyle:)`. Prioritize scaling
   essential content over decoration.
3. **Switch layout axis at accessibility sizes** — branch on
   `dynamicTypeSize.isAccessibilitySize` with `AnyLayout(HStackLayout())` /
   `AnyLayout(VStackLayout())` (UIKit: flip the stack axis on
   `preferredContentSizeCategory.isAccessibilityCategory`); give text the full line
   width and relax `lineLimit`.
4. **Large Content Viewer** — only for bars that legitimately can't grow: a tab bar takes
   under 10% of screen height, and scaled to accessibility sizes it would eat almost a
   quarter (WWDC24 10074). Scaling is always preferred; the viewer is the fallback, not
   the fix.

Test at all 12 sizes (Xcode Previews → Variants → Dynamic Type Variants), not just the
biggest — mid-range accessibility sizes catch different wrap points.

## The San Francisco family — pick by job

| Face | Job |
|---|---|
| SF Pro | The default; UI text everywhere |
| SF Pro Rounded | Friendlier numerals/labels — widgets, health/fitness data |
| SF Compact | watchOS (space-efficient counterpart) |
| SF Mono | Code, tabular alignment |
| New York | Serif — reading experiences, editorial contrast |
| SF Arabic / SF Arabic Rounded | Arabic script with its own optical sizes |

**The width axis** (Condensed / Regular / Compressed / Expanded):

- Default Regular; every non-Regular choice is a legibility decision — check it.
- Condensed: fit more text comfortably (long headlines wrap one line fewer).
- Compressed: display-only density, not body text.
- Expanded: display typesetting AND small secondary labels (wide + loose tracking).
- Width is a **fourth hierarchy lever** beside weight, size, color; all widths share identical
  vertical proportions, so mixing them never misaligns baselines. 2–3 styles are enough —
  pair one width with contrasting weights, one weight with contrasting widths, or oppose both
  for maximum display contrast.

## Custom fonts: what you owe back

System fonts get Dynamic Type free; a brand font makes *you* the type engine:

- Scale with `UIFontMetrics(forTextStyle:).scaledFont(for:)` + adjusts-for-content-size, or
  SwiftUI `Font.custom(_:size:relativeTo:)`; scale layout constants with `@ScaledMetric`.
- Test at the largest accessibility sizes — wrap to multiple lines rather than truncate.
- Verify script coverage for your locales; have a fallback plan for writing systems the brand
  font doesn't cover (the system stacks: `system-ui`, `ui-rounded`, `ui-serif`, `ui-monospace`
  on the web).
- Budget for the maintenance: tracking tables, optical sizing, and weight ramps are things SF
  does that a licensed font usually doesn't.

❌ Anti-patterns: hardcoded point sizes on user-facing text · kerning APIs for tracking jobs ·
custom font without `relativeTo:` (frozen size) · Compressed body text · truncation where
wrapping was possible.

## Output Format

Type review: `Element | Current (font/size/style) | Issue (hierarchy/scaling/legibility) | Fix`
— check Dynamic Type at XXL + AX sizes before signing off; pairs with the contrast rules in
`generators/accessibility-generator`.

## References

- https://developer.apple.com/videos/play/wwdc2020/10175/ (The details of UI typography)
- https://developer.apple.com/videos/play/wwdc2022/110381/ (The expanded San Francisco family)
- https://developer.apple.com/fonts/ (SF + New York downloads)
- https://developer.apple.com/design/human-interface-guidelines/typography
- Related skills: `design/sf-symbols` (symbols align to text styles), `foundation/attributed-string` (rich text mechanics), `generators/accessibility-generator`
