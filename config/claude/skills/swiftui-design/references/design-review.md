# 5-Dimension SwiftUI Design Review

A structured framework for evaluating SwiftUI UI quality. Score each dimension 1-10.

## How to Use

1. Take a screenshot of the UI
2. Evaluate each dimension below
3. Score honestly (most AI-generated UI scores 3-5)
4. Fix issues scoring < 7
5. Re-evaluate after fixes

---

## Dimension 1: Philosophy Alignment (1-10)

Does the design follow a coherent philosophy, or is it a random collection of trends?

| Score | Criteria                                                             |
| ----- | -------------------------------------------------------------------- |
| 1-3   | Random elements, no coherent style, looks like template mashup       |
| 4-5   | Some consistency, but generic (could be any app)                     |
| 6-7   | Clear style direction, recognizable identity                         |
| 8-10  | Strong philosophical foundation, every element supports the vision   |

**Check:**
- [ ] Design follows one of the 5 schools (Informational / Editorial / Expressive / Functional / Warm Minimal)
- [ ] Every element supports the chosen philosophy
- [ ] No contradictory design choices (e.g., playful icons + corporate typography)
- [ ] The design would be recognizable if the brand name were removed

---

## Dimension 2: Visual Hierarchy (1-10)

Can the user immediately understand what's important?

| Score | Criteria                                                             |
| ----- | -------------------------------------------------------------------- |
| 1-3   | Everything looks the same size/weight, no clear reading order        |
| 4-5   | Some hierarchy, but competing elements                               |
| 6-7   | Clear primary/secondary/tertiary levels                              |
| 8-10  | Instant comprehension, natural eye flow, nothing competing           |

**Check:**
- [ ] Screen title is the largest text element
- [ ] Primary action is visually dominant
- [ ] Secondary information is visually subordinate
- [ ] Whitespace creates natural grouping
- [ ] The eye knows where to go first, second, third

---

## Dimension 3: Craft Quality (1-10)

Are the details right? This is where AI-generated UI fails most.

| Score | Criteria                                                             |
| ----- | -------------------------------------------------------------------- |
| 1-3   | Random spacing, inconsistent sizing, alignment issues                |
| 4-5   | Generally okay, but noticeable rough edges                           |
| 6-7   | Consistent spacing, good alignment, clean typography                 |
| 8-10  | Pixel-perfect, every detail intentional, feels handcrafted           |

**Check:**
- [ ] Spacing follows 8pt grid consistently
- [ ] Font sizes are from the defined scale (not random)
- [ ] Colors are from the design token system
- [ ] Corner radii are consistent
- [ ] Icons are the same weight/style throughout
- [ ] No orphaned text (single word on a line)
- [ ] Images have consistent aspect ratios

---

## Dimension 4: Functionality (1-10)

Does the design serve the user's actual needs?

| Score | Criteria                                                             |
| ----- | -------------------------------------------------------------------- |
| 1-3   | Form over function, important actions hidden, confusing flow         |
| 4-5   | Usable but not efficient, some friction points                       |
| 6-7   | Good usability, clear actions, logical flow                          |
| 8-10  | Effortless to use, anticipates needs, zero friction                  |

**Check:**
- [ ] Primary action is immediately visible (not hidden in menus)
- [ ] Touch targets are ≥ 44×44pt
- [ ] Loading states are handled (no blank screens)
- [ ] Error states are handled with helpful messages
- [ ] Empty states guide the user on what to do
- [ ] Navigation is predictable (back goes back)
- [ ] Form inputs have clear labels and validation

---

## Dimension 5: Originality (1-10)

Does this look like a real product, or an AI template?

| Score | Criteria                                                             |
| ----- | -------------------------------------------------------------------- |
| 1-3   | Clearly AI-generated, uses banned slop patterns, template-like      |
| 4-5   | Generic but not offensive, could be any app                          |
| 6-7   | Has personality, some unique elements, feels intentional             |
| 8-10  | Distinctive, memorable, would stand out in App Store                 |

**Check:**
- [ ] No banned AI-slop patterns (see anti-ai-slop.md)
- [ ] At least one signature detail per screen
- [ ] Typography has personality (not just system defaults)
- [ ] Color palette is distinctive (not generic blue/purple)
- [ ] Layout is asymmetric where appropriate
- [ ] The design would not be confused with a template

---

## Review Summary Template

```
## Design Review: [Screen Name]

| Dimension         | Score | Notes                    |
| ----------------- | ----- | ------------------------ |
| Philosophy        | ?/10  | [specific feedback]      |
| Visual Hierarchy  | ?/10  | [specific feedback]      |
| Craft Quality     | ?/10  | [specific feedback]      |
| Functionality     | ?/10  | [specific feedback]      |
| Originality       | ?/10  | [specific feedback]      |
| **Average**       | **?/10** |                        |

### Issues to Fix (score < 7)
1. [Issue]: [Specific fix]
2. [Issue]: [Specific fix]

### What's Working Well
1. [Strength]
2. [Strength]
```

## Common Score Ranges

| Range   | Meaning                                            |
| ------- | -------------------------------------------------- |
| 1-3     | AI slop — needs complete redesign                  |
| 4-5     | Generic — functional but no personality            |
| 6-7     | Good — professional, minor issues                  |
| 8-10    | Excellent — distinctive, crafted, ship-ready       |

**Shipping threshold**: All dimensions ≥ 7, no dimension below 5.
