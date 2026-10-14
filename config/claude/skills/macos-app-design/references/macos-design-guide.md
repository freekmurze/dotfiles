# World-Class macOS App Design & Build Guide

Complete reference for designing and implementing native macOS apps.

## 1) Quality Bar: What "World-Class" Means

### "Good Mac Citizen" Test

Mac excellence is less about toolkit and more about embracing Mac conventions:
- Menus and keyboard shortcuts
- Window management and services
- Text behaviors and preferences

### Performance is Non-Negotiable

**Practical sanity test:** Open a large document, scroll hard, resize—must stay fast and responsive ("Moby Dick" test).

### Apple Design Awards Rubric

Use these categories for evaluation:
- **Delight and Fun**: Satisfying micro-interactions
- **Inclusivity**: VoiceOver, keyboard-only, reduced motion
- **Innovation**: Platform tech for meaningful improvement
- **Interaction**: Effortless, tailored to macOS workflows
- **Visuals and Graphics**: Cohesive iconography/typography
- **Social Impact**: Accessibility, broad user benefit

---

## 2) Liquid Glass Design System (macOS Tahoe 26+)

Apple's visual refresh built on **Liquid Glass**—translucent, dynamic material that reflects/refracts surroundings.

### What Liquid Glass Is

A "digital meta-material" that dynamically bends light and behaves fluidly. Reserved primarily for the **navigation/controls layer** floating above content.

### Where to Use Liquid Glass

**Do:**
- Navigation/controls (toolbars, sidebars, bars, key UI chrome)
- Let system components provide built-in behaviors (focus changes, shadowing, interaction glow)

**Don't:**
- Content layer (tables, lists, documents)—muddies hierarchy
- "Glass on glass" stacking—gets cluttered

### Accessibility Adaptations

Liquid Glass automatically adapts to:
- **Reduced Transparency**
- **Increased Contrast**
- **Reduced Motion**

Your app must remain usable with these settings enabled.

### New Design System Rules

- **Remove decorative backgrounds/borders** you added for emphasis
- Express hierarchy via **layout and grouping**, not decoration
- **Crowded toolbar** = signal to remove or demote actions
- Group bar items **by function and frequency**
- Avoid grouping text button + icon button (reads as one control)

---

## 3) macOS-Native IA & UI Patterns

### 3.1 Menu Bar: Your Command Center

**Required baseline:**
- Standard menu layout: **App / File / Edit / View / Window / Help** (+ Format if relevant)
- **Settings…** in App menu with **⌘,** behavior
- Menu items enable/disable based on focus/selection

**SwiftUI implementation:**
- Use `commands` in App, add via `CommandMenu` / `CommandGroup`
- Keyboard shortcuts and discoverability come "for free"

**World-class touches:**
- Help menu with real help content
- Menu organization reflecting user workflows
- Primary tasks near top, destructive tasks grouped and labeled

### 3.2 Keyboard Shortcuts and Full Keyboard Workflows

**Non-negotiable:**
- Every primary command reachable via keyboard
- Standard shortcuts behave normally (⌘C, ⌘V, ⌘Z, ⌘F, ⌘N, ⌘O)

**World-class:**
- Proper tabbing and focus rings on custom controls
- Command palette *in addition* to menus (not instead)

### 3.3 Windows: Multi-Window Platform

**Choose window model for archetype:**
- **Document-based**: Files as primary units, open/save/duplicate/revert
- **Library + detail**: Sidebar lists items, detail/editor in main area
- **Utility**: Lightweight single window, optional menu bar

**Required:**
- Window resizing as first-class interaction—layouts fluidly adapt
- Responsive under resize

**World-class:**
- Multiple windows showing different content (or views of same)
- Respect close/minimize/zoom/toolbar/titlebar, tabbing, fullscreen

### 3.4 Sidebars: Structure, Navigation, Context

macOS sidebars refract/reflect content behind them while maintaining context.

**Do:**
- Top-level navigation with clear hierarchy
- Scannable, stable items (users build muscle memory)
- Let content extend beneath/behind sidebars

**Don't:**
- Overload with dense controls—use inspectors/toolbars for actions

### 3.5 Toolbars and Bar Content

**New system encourages:**
- Remove custom backgrounds/borders
- Group by function/frequency
- Demote secondary actions to "more" menu
- Primary action visually distinct (often tinted)

### 3.6 Text System Behavior

**Users expect:**
- Standard text shortcuts and behaviors
- Correct selection rules
- Undo/redo
- Services menu integration
- Context menu items

**Practical rule:** Use system text components to inherit Mac editing ecosystem.

### 3.7 Services, Share, System Integrations

A "Mac citizen" participates in:
- **Services menu** integrations
- Standard share flows
- Drag & drop, copy/paste
- Quick Look previews (document/media apps)

---

## 4) Visual Language

### 4.1 Typography and Layout

- Prefer system text styles and default spacing
- Adjust only for hierarchy and readability
- Design for large displays, resizable windows, dense information, long sessions

### 4.2 Color and Materials

Use system colors/materials for automatic adaptation to:
- Light/dark environments
- User customization (accent colors, increased contrast)

Brand through content, iconography, subtle accents—not by replacing system structure.

### 4.3 App Icons: Modern, Layered, Multi-Appearance

New look includes multi-layer Liquid Glass-crafted icons with personalization (light/dark, tints, "clear" look).

#### Icon Composer Workflow

**Tool requirements:** macOS Sequoia 15.3+

**Canvas sizes:**
- macOS/iPhone/iPad: 1024px
- Watch: 1088px (optically larger)

**Workflow:**
1. Design as layers—keep source art flat and controllable
2. Add dynamic glass properties in Icon Composer (not baked in)
3. Export vectors as **SVG**, complex assets as **PNG**
4. Convert text to outlines if needed
5. Don't export enclosure mask (system applies cropping)

**Layer groups:** Up to **4 groups/layers** of glass complexity (intentional limit)

**Appearances supported:**
- Default, Dark
- Clear (light/dark)
- Tinted (light/dark)

**Mono legibility:** Set at least one key element to white, tune conversion for contrast.

**Pitfalls:**
- Highly translucent art becomes unclear in monochrome
- Watch blur/gradients and detail at small sizes

### 4.4 SF Symbols

**SF Symbols 7:**
- 6,900+ symbols aligned with San Francisco typeface
- Multiple weights/scales
- Draw animations and annotation tools
- Variable rendering, gradients, Magic Replace
- Requires macOS Sonoma+ to run app

**Guideline:** Use SF Symbols for system concepts. Only design custom symbols when domain requires—keep them SF-like and accessible.

### 4.5 Motion and Micro-Interactions

Motion communicates hierarchy and responsiveness:
- Clarify state change, focus, hierarchy, transitions
- Respect Reduced Motion (system adapts Liquid Glass; your custom animations must too)
- Subtle, high-quality over flashy

---

## 5) Accessibility and Inclusive Design

**Design in, not added later.**

### Baseline Requirements

| Area | Requirement |
|------|-------------|
| **VoiceOver** | Meaningful labels/values/hints for every interactive element |
| **Keyboard** | Everything operable without mouse/trackpad |
| **Contrast** | Readable in light/dark and increased contrast modes |
| **Reduce Motion** | UI usable and understandable with setting enabled |
| **Reduce Transparency** | UI remains functional |
| **Localization** | Layouts handle longer strings, different writing systems |

### Inclusive Craft Moves

- Don't encode meaning by color alone
- Redundancy: icon + label, shape + color, text + sound
- Reduce cognitive load: clear IA, consistent terminology, predictable controls

---

## 6) Engineering for Excellence

### 6.1 Architecture

- Clear separation between **domain model** and **UI state**
- Concurrency keeping UI responsive (avoid blocking main thread)
- Deliberate **undo/redo** strategy, document/versioning, autosave

**Pattern:** Define internal "app intents" (command model) so UI, menus, shortcuts, toolbar all invoke same code paths.

### 6.2 SwiftUI vs AppKit

**Prefer SwiftUI** for new apps—faster iteration, modern patterns.

**Drop to AppKit when:**
- Advanced text editing matching macOS standards
- Deep window customization unavailable in SwiftUI
- Certain pro-workflow integrations

**Rule:** If using AppKit, keep bridging layer thin and well-tested.

### 6.3 Performance Discipline

Repeatable checklist:
- Large-file open, scroll, resize test
- Memory sanity: no unbounded caches, retain cycles; windows close cleanly
- Smooth interaction: scrolling/resizing responsive

### 6.4 Fit and Finish Features

World-class apps include:
- Fully functional menu structure + shortcuts
- Native-feeling Settings (⌘,)
- Window tabbing support and multi-window
- Help content
- Scriptability/automation hooks
- Respect for user preferences (appearance, accent colors)

---

## 7) ADA-Inspired Evaluation Rubric

Use as release gate:

### Delight and Fun
- [ ] Satisfying micro-interactions enhancing understanding?
- [ ] Feels "alive" in restrained, high-quality way?

### Inclusivity
- [ ] VoiceOver user can complete core workflow?
- [ ] Keyboard-only user can do everything?
- [ ] Reduced motion/transparency preserves usability?

### Interaction
- [ ] Control effortless and tailored to macOS workflows?

### Visuals and Graphics
- [ ] UI cohesive with consistent iconography/typography?
- [ ] SF Symbols used appropriately, custom symbols match system?

### Innovation (Differentiating)
- [ ] Leveraging platform tech for meaningfully better experience?

---

## 8) Definition of Done Checklist

### Mac Citizen Essentials

- [ ] Standard menu bar structure; key commands discoverable
- [ ] Settings opens via ⌘, with native Mac feel
- [ ] Undo/redo works reliably
- [ ] Copy/paste/selection behave normally in text surfaces
- [ ] Multi-window behavior sensible; resizing robust and fast

### Liquid Glass Readiness

- [ ] Custom bars/backgrounds removed unless strongly justified
- [ ] Liquid Glass for functional/navigation layers; content clear
- [ ] "Glass on glass" stacking avoided

### Icon + Branding

- [ ] App icon modern and layered (Icon Composer), tested across appearances
- [ ] SF Symbols for standard concepts; custom symbols SF-consistent and accessible

### Accessibility

- [ ] VoiceOver labels/hints for all interactive elements
- [ ] Full keyboard navigation supported
- [ ] Reduced motion/transparency/increased contrast tested

### Performance & Stability

- [ ] Large-document / heavy-data scenarios responsive
- [ ] No major UI hitches while scrolling/resizing

---

## 9) Source References

**Apple Official:**
- [Apple Newsroom - Design Update](https://www.apple.com/newsroom/2025/06/apple-introduces-a-delightful-and-elegant-new-software-design/)
- [WWDC25: Meet Liquid Glass](https://developer.apple.com/videos/play/wwdc2025/219/)
- [WWDC25: Get to know the new design system](https://developer.apple.com/videos/play/wwdc2025/356/)
- [WWDC25: New look of app icons](https://developer.apple.com/videos/play/wwdc2025/220/)
- [WWDC25: Create icons with Icon Composer](https://developer.apple.com/videos/play/wwdc2025/361/)
- [SF Symbols](https://developer.apple.com/sf-symbols/)
- [Icon Composer](https://developer.apple.com/icon-composer/)
- [Apple Design Resources](https://developer.apple.com/design/resources/)
- [Apple Design Awards](https://developer.apple.com/design/awards/)

**Community:**
- [App Feel on Mac](https://coyotetracks.org/blog/app-feel-on-mac/)
- [Best in Class macOS App](https://swiftjectivec.com/What-does-a-best-in-class-macOS-app-look-like/)
- [The macOS App Icon Book](https://flarup.shop/products/the-macos-app-icon-book)
