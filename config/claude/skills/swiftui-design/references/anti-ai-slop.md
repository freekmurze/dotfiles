# Anti-AI-Slop Rules for SwiftUI

These are specific patterns that make UI look "AI-generated" or generic. Each rule has a banned pattern, why it's bad, and what to do instead.

## Rule 1: No Purple-Blue Gradient Backgrounds

**Banned:**
```swift
LinearGradient(
    colors: [.purple, .blue],
    startPoint: .topLeading,
    endPoint: .bottomTrailing
)
```

**Why:** Every AI tool defaults to this gradient. It signals "I didn't think about the design."

**Instead:** Use a solid warm neutral background:
```swift
Color(hex: "FAFAF8")  // Warm white
// or
Color(hex: "F5F3EF")  // Light warm gray
```

If you must use a gradient, make it subtle and brand-specific:
```swift
LinearGradient(
    colors: [Color(hex: "FFF5F0"), Color(hex: "FAFAF8")],
    startPoint: .top,
    endPoint: .bottom
)
```

## Rule 2: No Emoji as Icons

**Banned:**
```swift
Text("🏠")  // As a tab icon
Text("❤️")  // As a like button
```

**Why:** Emoji rendering varies by platform, can't be styled, and looks unprofessional.

**Instead:** Use SF Symbols:
```swift
Image(systemName: "house.fill")
Image(systemName: "heart.fill")
```

For custom needs, use a proper icon set (SF Symbols has 5,000+ icons).

## Rule 3: No Rounded Cards with Left Border Accent

**Banned:**
```swift
HStack {
    Rectangle()
        .fill(Color.purple)
        .frame(width: 4)
    VStack(alignment: .leading) {
        Text("Title")
        Text("Description")
    }
}
.padding()
.background(Color(.systemBackground))
.cornerRadius(12)
```

**Why:** This is the most common AI-generated card pattern. It's lazy and generic.

**Instead:** Use proper card design:
```swift
VStack(alignment: .leading, spacing: DesignTokens.spaceS) {
    Text("Title")
        .font(.system(size: 18, weight: .semibold, design: .serif))
    Text("Description")
        .font(.body)
        .foregroundStyle(.secondary)
}
.padding(DesignTokens.spaceM)
.frame(maxWidth: .infinity, alignment: .leading)
.background(DesignTokens.surface)
.clipShape(RoundedRectangle(cornerRadius: DesignTokens.radiusM))
```

## Rule 4: No Generic Hero Sections

**Banned:**
```swift
VStack(spacing: 16) {
    Text("Welcome to App")
        .font(.largeTitle.bold())
    Text("The best app for everything")
        .font(.title3)
        .foregroundColor(.gray)
    Button("Get Started") { }
        .buttonStyle(.borderedProminent)
}
.frame(maxWidth: .infinity)
.padding(.vertical, 80)
.background(
    LinearGradient(colors: [.purple, .blue], ...)
)
```

**Why:** This pattern is everywhere. It wastes the most valuable screen real estate on marketing copy nobody reads.

**Instead:** Show actual content immediately:
```swift
VStack(alignment: .leading, spacing: DesignTokens.spaceL) {
    // Small, tasteful header
    Text("Today")
        .font(.system(size: 28, weight: .bold, design: .serif))
    
    // Real content
    ForEach(items) { item in
        ItemRow(item: item)
    }
}
.padding(DesignTokens.spaceM)
```

## Rule 5: No Inter/Roboto on Apple Platforms

**Banned:**
```swift
Text("Hello")
    .font(.custom("Inter", size: 16))
```

**Why:** Inter and Roboto are Google's fonts for Android/web. Using them on Apple platforms looks like a web app ported to native.

**Instead:** Use Apple's system fonts:
```swift
// Default: SF Pro (sans-serif)
Text("Hello")
    .font(.body)

// Display: New York (serif) for headings
Text("Hello")
    .font(.system(size: 28, weight: .bold, design: .serif))

// If you need a custom font, choose one that fits Apple's ecosystem:
// - Charter, Georgia, Iowan (serif alternatives)
// - Avenir, Futura, Gill Sans (geometric sans)
```

## Rule 6: No Neon-on-Dark Cyber Aesthetic

**Banned:**
```swift
Color(hex: "0D1117")  // GitHub dark
// with neon accent colors like #00FF88, #00D4FF
```

**Why:** This is the "I built a crypto dashboard" default. It's overused and hard to read.

**Instead:** Use Apple's native dark mode colors:
```swift
// Let the system handle dark/light
Color(.systemBackground)
Color(.secondarySystemBackground)

// Or define warm dark colors
Color(hex: "1C1C1E")  // iOS dark background
Color(hex: "2C2C2E")  // iOS dark secondary
```

## Rule 7: No Symmetric 3-Column Feature Grids

**Banned:**
```swift
LazyVGrid(columns: Array(repeating: GridItem(.flexible()), count: 3)) {
    FeatureCard(icon: "star", title: "Fast")
    FeatureCard(icon: "shield", title: "Secure")
    FeatureCard(icon: "bolt", title: "Easy")
}
```

**Why:** This is the SaaS landing page default. It doesn't work for native apps.

**Instead:** Use asymmetric, content-driven layouts:
```swift
VStack(spacing: DesignTokens.spaceL) {
    // Primary feature — large
    FeaturedCard(item: featuredItem)
        .frame(height: 200)
    
    // Secondary features — smaller, horizontal
    HStack(spacing: DesignTokens.spaceM) {
        SecondaryCard(item: item1)
        SecondaryCard(item: item2)
    }
    .frame(height: 120)
}
```

## Rule 8: No AI-Drawn SVG Illustrations

**Banned:**
- Inline SVG paths that draw abstract "people working" illustrations
- CSS-drawn silhouettes
- Any illustration that looks like it came from unDraw or Humaaans

**Why:** These are immediately recognizable as AI-generated stock content.

**Instead:**
- Use real photography (from Unsplash, Pexels, or the brand's own photos)
- Use SF Symbols at large scale as visual elements
- Use clean geometric shapes or patterns from the brand
- Use a simple placeholder with a note: "Replace with brand photo"

## Rule 9: No Default Blue Tint Everywhere

**Banned:**
```swift
// Everything uses the default blue accent
Button("Action") { }  // Blue
Toggle("Setting", isOn: $value)  // Blue
Link("Visit", url: url)  // Blue
```

**Why:** Default blue means "I didn't customize anything."

**Instead:** Define a custom accent color:
```swift
@main
struct MyApp: App {
    var body: some Scene {
        WindowGroup {
            ContentView()
                .tint(DesignTokens.accent)  // Your brand color
        }
    }
}
```

## Rule 10: No Random Spacing

**Banned:**
```swift
.padding(.horizontal, 13)
.padding(.vertical, 7)
.frame(height: 43)
```

**Why:** Odd numbers look accidental. They create visual inconsistency.

**Instead:** Use an 8pt grid:
```swift
.padding(.horizontal, DesignTokens.spaceM)  // 16
.padding(.vertical, DesignTokens.spaceS)    // 8
.frame(height: 44)  // Minimum tap target
```

## Self-Check: Is Your UI AI Slop?

Ask these questions:
1. Would this look at home on a generic SaaS landing page?
2. Are there purple-blue gradients?
3. Are emoji used as functional icons?
4. Is the spacing random or systematic?
5. Could you swap the brand name and it would still look the same?
6. Is there a "Welcome to [App]" hero section?
7. Are all cards the same rounded rectangle with a left accent border?

If you answered "yes" to 3+ questions, the design needs work.
