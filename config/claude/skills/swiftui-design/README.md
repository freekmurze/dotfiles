<div align="center">

# SwiftUI Design Skill

> *"One prompt. A shippable SwiftUI interface."*

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Agent-Agnostic](https://img.shields.io/badge/Agent-Agnostic-blueviolet)](https://skills.sh)
[![Skills](https://img.shields.io/badge/skills.sh-Compatible-green)](https://skills.sh)

</div>

<p align="center">
  <a href="README.md"><img src="https://img.shields.io/badge/English-blue?style=flat-square" alt="English"></a>
  <a href="README-zh.md"><img src="https://img.shields.io/badge/中文-red?style=flat-square" alt="中文"></a>
</p>

**Type one line in your agent. Get back a SwiftUI interface that doesn't look AI-generated.**

In 3 to 30 minutes, you can design a **signature iOS interface** — not "AI-made-it-looks-okay" quality, but something that looks like it was crafted by a designer with taste.

Anti-AI Slop 6 Iron Rules, Design Direction Advisor, Brand Asset Protocol, 5-Dimension Review — all built in. Give the skill your brand palette and it will read your vibe; give nothing and the 5 built-in design languages will still hold their ground.

Works across agents — Claude Code, Cursor, Codex, OpenCode, Hermes, all supported.

## Install & Go

```bash
npx skills add wholiver/swiftui-design-skill -g -y
```

## What Can It Do

| Capability | Deliverable | Typical Time |
|------------|-------------|--------------|
| SwiftUI Interface Design | Compilable SwiftUI code · Design system tokens · Layout patterns | 10–15 min |
| Design Direction Advisor | 5 schools × N design philosophies · 3 recommended directions · Visual anchor descriptions | 5 min |
| Brand Asset Integration | brand-spec.md · Color palette · Typography · Spacing system | 5–10 min |
| Anti-AI Slop Review | Item-by-item check · Concrete fix code · Alternative approaches | 3–5 min |
| 5-Dimension Review | Radar chart scoring · Keep/Fix/Quick Wins · Actionable fix checklist | 3 min |
| Layout Pattern Library | 9 common layouts · Copy-paste SwiftUI code | Instant |
| Animation Design | Spring/parallax/pull-to-refresh · Compilable code | 5–10 min |

## Core Mechanisms

### Anti-AI Slop 6 Iron Rules

This is the hardest set of rules in the skill. All designs must pass these 6 checks:

| # | Iron Rule | ❌ Forbidden | ✅ Alternative |
|---|-----------|-------------|----------------|
| 1 | Start from context | Inventing from a blank canvas | Ask about design system / UI kit / brand assets first |
| 2 | Junior Designer Mode | Waiting for the perfect solution | Gray placeholder > bad SVG |
| 3 | Give variants, not finals | One "final answer" | 2–3 differentiated directions |
| 4 | Placeholder > bad implementation | AI-generated clip art | Clean gray placeholder + text label |
| 5 | System-first, don't fill | Packing every pixel | Every element must justify its existence |
| 6 | Ban AI slop patterns | Purple-blue gradients, emoji icons, rounded cards + left border | Single warm accent color, SF Symbols, serif headlines |

See `references/anti-ai-slop.md` for details.

### Design Direction Advisor

When requirements are vague, the skill recommends 3 differentiated directions from 5 major design schools:

| School | Characteristics | Signature Style |
|--------|-----------------|-----------------|
| **Informational** | Data-first, chart-dense | Bloomberg Terminal |
| **Editorial** | Magazine layout, serif type, generous whitespace | NYT, Medium |
| **Expressive** | Bold color, asymmetric layout, motion-forward | Stripe, Linear |
| **Functional** | Dense tool feel, monospace accents | Things, OmniFocus |
| **Warm Minimal** | Soft neutrals, rounded corners, subtle texture | Notion, Craft |

### Brand Asset Protocol

A mandatory 5-step hard process when working with a specific brand:

| Step | Action | Purpose |
|------|--------|---------|
| 1 · Ask | Does the user have brand guidelines? | Respect existing assets |
| 2 · Search | Search the brand's official pages | Obtain real materials |
| 3 · Download | Download actual asset files | PNG/SVG logos, fonts |
| 4 · Verify | Verify colors match official sources | Cross-check hex values |
| 5 · Write | Generate brand-spec.md | Record complete design system |

Quality threshold: 5 real brand colors, 10 design tokens, 2 font families, 8pt spacing grid.

### 5-Dimension Review

Every design must pass a 5-dimension review before delivery:

| Dimension | Scoring Criteria | Minimum Score |
|-----------|------------------|---------------|
| 🎯 Philosophy Consistency | Does the design embody the chosen design philosophy? | ≥ 7/10 |
| 📐 Visual Hierarchy | Is the information priority clear? | ≥ 7/10 |
| 🔧 Detail Execution | Are spacing, typography, and colors precise? | ≥ 7/10 |
| ⚡ Functionality | Does the layout serve user goals? | ≥ 7/10 |
| ✨ Originality | Is there at least 1 signature detail? | ≥ 7/10 |

## File Structure

```
swiftui-design-skill/
├── SKILL.md                           ← Core definition (278 lines)
├── references/
│   ├── anti-ai-slop.md                ← Anti-AI Slop detailed rules (268 lines)
│   ├── layout-patterns.md             ← 9 layout patterns + copy-paste code (265 lines)
│   ├── typography-color.md            ← Typography hierarchy + color system (172 lines)
│   ├── design-review.md               ← 5-Dimension Review process (151 lines)
│   └── swift-extensions.md            ← Color/Font/Animation extensions (373 lines, essential)
└── templates/
    └── brand-spec.md                  ← Brand spec template (77 lines)
```

## Difference from swiftui-expert-skill

| | SwiftUI Design Skill | swiftui-expert-skill |
|---|---|---|
| **Focus** | Visual design, aesthetics, brand feel | Code quality, performance, correctness |
| **Question** | "How does it look?" | "Is the code correct?" |
| **Output** | Design direction, color palette, layout, review | Code review, Instruments analysis, API modernization |
| **Complementary** | Use Design first to set direction | Then use Expert to ensure code quality |

The two skills work together: Design decides **what to build**, Expert ensures **how to build it right**.

## License

MIT — use freely, but please keep the original author attribution.
