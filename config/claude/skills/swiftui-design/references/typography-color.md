> **Note**: Requires `swift-extensions.md` for `Color(hex:)` and `Color(light:dark:)` implementations.

# Typography & Color Systems for SwiftUI

## Typography System

### Font Hierarchy

| Level    | Size    | Weight     | Font       | Use For                    |
| -------- | ------- | ---------- | ---------- | -------------------------- |
| Display  | 28-34pt | Bold       | New York   | Screen titles              |
| Heading  | 20-24pt | Semibold   | SF Pro     | Section titles             |
| Subhead  | 17pt    | Medium     | SF Pro     | Card titles, list items    |
| Body     | 15-17pt | Regular    | SF Pro     | Paragraph text             |
| Caption  | 12-13pt | Regular    | SF Pro     | Labels, metadata, timestamps |
| Footnote | 11pt    | Regular    | SF Pro     | Legal text, disclaimers    |

### Font Implementation

```swift
enum AppFont {
    // Display: Serif for personality (use .serif design)
    static let display = Font.system(size: 28, weight: .bold, design: .serif)
    static let displaySmall = Font.system(size: 22, weight: .bold, design: .serif)
    
    // Heading: Sans-serif for clarity
    static let heading = Font.system(size: 20, weight: .semibold, design: .default)
    static let subheading = Font.system(size: 17, weight: .medium, design: .default)
    
    // Body: System default
    static let body = Font.system(size: 16, weight: .regular, design: .default)
    static let bodySmall = Font.system(size: 14, weight: .regular, design: .default)
    
    // Caption: Smaller, lighter
    static let caption = Font.system(size: 12, weight: .regular, design: .default)
    static let captionMedium = Font.system(size: 12, weight: .medium, design: .default)
}
```

### Typography Rules

1. **Maximum 2 font families** per screen (e.g., New York + SF Pro)
2. **Maximum 4 font sizes** per screen (e.g., display + heading + body + caption)
3. **Line height**: 1.4× for body text, 1.2× for headings
4. **Line length**: 45-75 characters per line for body text
5. **Never** use Inter, Roboto, or Poppins on Apple platforms
6. **Always** use `.font()` modifier, not raw `UIFont`/`NSFont`

### Dynamic Type Support

```swift
// Use @ScaledMetric for size-dependent values
@ScaledMetric(relativeTo: .body) var iconSize: CGFloat = 20

Image(systemName: "star")
    .font(.system(size: iconSize))
```

### When to Use Serif (New York)

- App title / branding
- Article/blog post headings
- Editorial content
- Landing page hero text
- When you want to feel "premium" or "literary"

**Don't use serif for:**
- Body text (harder to read at small sizes)
- Data-heavy interfaces
- Technical/admin UIs

---

## Color System

### Color Palette Structure

```swift
enum AppColors {
    // Accent — 1 warm color that defines the brand
    static let accent = Color(hex: "E85D3A")         // Warm coral
    static let accentLight = Color(hex: "FFF0EB")    // Very light accent
    static let accentDark = Color(hex: "C44A2E")     // Dark accent
    
    // Neutrals — warm undertone, not pure gray
    static let background = Color(light: "FAFAF8", dark: "1C1C1E")
    static let surface = Color(light: "FFFFFF", dark: "2C2C2E")
    static let textPrimary = Color(light: "1A1A1A", dark: "F2F2F7")
    static let textSecondary = Color(light: "6B6B6B", dark: "8E8E93")
    static let textTertiary = Color(light: "999999", dark: "636366")
    static let divider = Color(light: "E5E5E5", dark: "38383A")
    
    // Semantic — contextual colors
    static let success = Color(hex: "34C759")
    static let warning = Color(hex: "FF9500")
    static let error = Color(hex: "FF3B30")
    static let info = Color(hex: "007AFF")
}
```

### Color Rules

1. **One accent color** — not a rainbow. One warm color defines the brand.
2. **Warm neutrals** — grays should have a warm undertone (slightly yellow/red), not blue.
3. **Semantic colors** — use context-based names (`textPrimary`) not appearance-based (`darkGray`).
4. **Dark mode** — every color must have a light AND dark variant.
5. **Contrast** — minimum 4.5:1 for body text, 3:1 for large text (WCAG AA).
6. **No pure black** — use `#1A1A1A` or `#1C1C1E`, not `#000000`.
7. **No pure white** — use `#FAFAF8` or `#FFFFFF` with slight warmth.

### Dark Mode Implementation

```swift
extension Color {
    init(light lightHex: String, dark darkHex: String) {
        self.init(uiColor: UIColor(
            light: UIColor(Color(hex: lightHex)),
            dark: UIColor(Color(hex: darkHex))
        ))
    }
}

extension UIColor {
    convenience init(light: UIColor, dark: UIColor) {
        self.init { traitCollection in
            traitCollection.userInterfaceStyle == .dark ? dark : light
        }
    }
}
```

### Accent Color Usage

```
Primary actions:    accent (buttons, links, active states)
Hover/pressed:      accentDark (interactive feedback)
Selection:          accentLight (background highlight)
Decorative:         accent at 10-20% opacity (subtle accents)
```

### Color Anti-Patterns

| Banned                        | Why                        | Instead                    |
| ----------------------------- | -------------------------- | -------------------------- |
| Purple-blue gradient          | AI slop default            | Solid warm accent          |
| Neon green (#00FF88)          | Crypto/dashboard cliché    | Muted success green        |
| Pure black background (#000)  | Too harsh, no depth        | Dark gray (#1C1C1E)        |
| Rainbow accent colors         | No visual identity         | One consistent accent      |
| Gray text on gray background  | Low contrast               | Ensure 4.5:1 ratio         |
| Random hex values             | Inconsistent               | Use design token system    |

### Spacing System (8pt Grid)

```swift
enum Spacing {
    static let xxs: CGFloat = 2    // Hairline gaps
    static let xs: CGFloat = 4     // Tight spacing
    static let s: CGFloat = 8      // Default small
    static let m: CGFloat = 16     // Default medium
    static let l: CGFloat = 24     // Section spacing
    static let xl: CGFloat = 32    // Large section spacing
    static let xxl: CGFloat = 48   // Screen-level spacing
    static let xxxl: CGFloat = 64  // Hero spacing
}
```

**Rules:**
- All spacing values must be multiples of 4 (preferably 8)
- Horizontal padding: 16pt (phone), 20-24pt (tablet)
- Vertical spacing between sections: 24-32pt
- Vertical spacing within sections: 8-16pt
- Minimum touch target: 44×44pt
