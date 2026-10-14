---
name: ux-writing
description: Interface copy that works — voice/tone, the PACE framework, alert anatomy, naming features, empty states, error messages, and the small-word edits that measurably improve UX. Distilled from Apple's writing sessions (2017–2026). Use when writing or reviewing any user-facing text — labels, alerts, onboarding, notifications, empty states, paywalls — or naming a feature.
allowed-tools: [Read, Write, Edit, Glob, Grep]
last_verified: 2026-07-16
review_by: 2027-06-22
---

# UX Writing

Interface text is interface design. This skill covers writing it: the frameworks Apple's
writers use, the anatomy of the hardest components (alerts, errors, empty states, names), and
the mechanical edits that improve any screen's copy in minutes.

## When This Skill Activates

- Writing or reviewing user-facing copy: labels, buttons, alerts, onboarding, notifications
- Naming a new feature, tier, or setting
- Error messages that frustrate, empty states that dead-end, permission prompts that get denied
- Localization prep (copy that survives translation)

## The PACE framework (Writing for Interfaces)

Run every screen's copy through four questions:

- **Purpose** — what must this screen communicate? Headers and buttons should carry the screen
  alone; know what to leave out. Hierarchy in words mirrors hierarchy in layout.
- **Anticipation** — what will the user ask or do next? Write for that ("What comes next?").
- **Context** — where is the user right now (quiet home vs. busy airport, mid-task vs. idle)?
  That decides how much to say.
- **Empathy** — write for a person having a problem, not for the system reporting one.

## Voice vs. tone

- **Finding voice**: describe the app as a person, cluster the personality adjectives, keep 2–3
  as your voice attributes — constant across the app, while tone flexes per moment (more
  directness for high-stakes info, sparing warmth in celebrations).
- The test that recurs in every Apple session: **read it aloud**. Does it sound like the best
  version of your app talking?
- ❌ Jargon ("it can make people feel left out"); ❌ personality at the cost of usefulness.

## The four small edits (biggest impact per minute)

1. **Cut filler**: "easily," "quickly," "simply," "just" — ✅ "Enter your license plate to pay
   for parking" ❌ "Simply enter your license plate number to quickly pay for parking."
2. **Kill repetition** — merge redundant sentences: ❌ "We're running late. Your driver won't
   make it on time. They'll be there in 10 minutes." ✅ "Delivery delayed 10 minutes."
3. **Lead with the why** — benefit before action: ✅ "To get reservation updates, enter your
   phone number." Notifications too: ✅ "Keep your streak going by solving today's crossword."
4. **Keep a word list** — three columns: approved term · alternatives to avoid · definition
   (e.g. "alias — not handle/username/title — the player's in-game name, shown to others").
   Reuse identical button labels ("Next") everywhere; consistency builds trust.

## Alert anatomy (Writing Great Alerts)

First: should this be an alert at all? Alerts interrupt by design — reserve them for things
the user must know *now*.

- **Title**: the main point, one sentence or less. **Message**: only if needed.
- **The scannability test**: title + buttons alone must convey the situation.
- ❌ Vague hedging ("You may need to…") — name the file, name the fix.

## Naming features (Craft Clear Names, WWDC26)

Judge every candidate on three criteria: it **belongs** (fits the app and its location), it
**sets expectations** (readers predict what they'll get), it **works everywhere** (languages,
markets, platforms).

- Process: **Think / Feel / Do** — what should the audience think, feel, and do? Group the
  answers into themes; name from the themes; test candidates in natural sentences ("just
  search for ___").
- ✅ Industry-standard words in high-stakes domains: "Balance," not "Spending Power."
- ✅ Compounds are fine when the parts self-explain ("AutoMix"); emotional names win where the
  moment is emotional ("Memories").
- ❌ Implementation jargon — describe the outcome, not the mechanism.

## Empty states, permissions, localization

- **Empty states**: never dead-end — say what will appear and how, ideally with example input.
- **Permission prompts**: explain value first, ask at the contextually right moment (mechanics
  in `generators/permission-priming`).
- **Localization**: plain language over idiom/humor; leave layout room and accessibility text
  for every element.

## Output Format

Copy review table:

```
Location | Current | Problem (PACE / filler / vague button / jargon) | Rewrite
```

ordered by frequency-of-exposure (navigation labels before deep-screen prose).

## References

- https://developer.apple.com/videos/play/wwdc2022/10037/ (Writing for interfaces)
- https://developer.apple.com/videos/play/wwdc2024/10140/ (Add personality through UX writing)
- https://developer.apple.com/videos/play/wwdc2025/404/ (Small writing changes)
- https://developer.apple.com/videos/play/wwdc2026/290/ (Craft clear names)
- https://developer.apple.com/videos/play/wwdc2017/813/ (Writing Great Alerts)
- Related skills: `generators/permission-priming`, `generators/push-notifications` (notification copy), `app-store/app-description-writer` (store copy)
