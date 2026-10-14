# Swift Extensions for SwiftUI Design

These extensions are required for the code examples in this skill to compile. Copy them into your project.

## Color Extensions

### Color(hex:)

Converts hex string to SwiftUI Color. Supports 3, 4, 6, and 8 character hex strings.

```swift
import SwiftUI

extension Color {
    init(hex: String) {
        let hex = hex.trimmingCharacters(in: CharacterSet.alphanumerics.inverted)
        var int: UInt64 = 0
        Scanner(string: hex).scanHexInt64(&int)
        let a, r, g, b: UInt64
        switch hex.count {
        case 3: // RGB (12-bit)
            (a, r, g, b) = (255, (int >> 8) * 17, (int >> 4 & 0xF) * 17, (int & 0xF) * 17)
        case 6: // RGB (24-bit)
            (a, r, g, b) = (255, int >> 16, int >> 8 & 0xFF, int & 0xFF)
        case 8: // ARGB (32-bit)
            (a, r, g, b) = (int >> 24, int >> 16 & 0xFF, int >> 8 & 0xFF, int & 0xFF)
        default:
            (a, r, g, b) = (255, 0, 0, 0)
        }
        self.init(
            .sRGB,
            red: Double(r) / 255,
            green: Double(g) / 255,
            blue: Double(b) / 255,
            opacity: Double(a) / 255
        )
    }
}
```

### Color(light:dark:) — Adaptive Colors

Creates colors that automatically adapt to light/dark mode.

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

### Predefined Design Colors

```swift
extension Color {
    // Warm accent (coral)
    static let appAccent = Color(hex: "E85D3A")
    static let appAccentLight = Color(hex: "FFF0EB")
    static let appAccentDark = Color(hex: "C44A2E")
    
    // Warm neutrals
    static let appBackground = Color(light: "FAFAF8", dark: "1C1C1E")
    static let appSurface = Color(light: "FFFFFF", dark: "2C2C2E")
    static let appTextPrimary = Color(light: "1A1A1A", dark: "F2F2F7")
    static let appTextSecondary = Color(light: "6B6B6B", dark: "8E8E93")
    static let appDivider = Color(light: "E5E5E5", dark: "38383A")
}
```

---

## Font Extensions

### Correct "New York" Font Usage

SwiftUI has built-in serif support. Use `design: .serif` instead of `Font.custom("New York", ...)`:

```swift
// ✅ CORRECT — iOS 17+
Font.system(size: 28, weight: .bold, design: .serif)

// ✅ CORRECT — Shorthand
Font.system(.title, design: .serif).bold()

// ❌ WRONG — "New York" is not a custom font name
Font.custom("New York", size: 28)
```

### AppFont Enum (Corrected)

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

### Using Custom Fonts (If Needed)

If you must use a custom font (e.g., from a brand), use the PostScript name:

```swift
// Find the PostScript name in Font Book.app → Select font → Info tab
Font.custom("NewYork-Regular", size: 17)
Font.custom("NewYork-Bold", size: 28)
```

---

## Animation Design Patterns

### Spring Animation (Card Expand)

```swift
struct ExpandableCard: View {
    @State private var isExpanded = false
    
    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Title")
                .font(AppFont.heading)
            
            if isExpanded {
                Text("Detailed content that appears with a spring animation.")
                    .font(AppFont.body)
                    .foregroundStyle(.secondary)
                    .transition(.opacity.combined(with: .move(edge: .top)))
            }
        }
        .padding(16)
        .frame(maxWidth: .infinity, alignment: .leading)
        .background(Color.appSurface)
        .clipShape(RoundedRectangle(cornerRadius: 12))
        .onTapGesture {
            withAnimation(.spring(response: 0.4, dampingFraction: 0.8)) {
                isExpanded.toggle()
            }
        }
    }
}
```

### Pull-to-Refresh with Brand Color

```swift
struct BrandRefreshable: View {
    @State private var items: [String] = ["Item 1", "Item 2"]
    
    var body: some View {
        List(items, id: \.self) { item in
            Text(item)
        }
        .refreshable {
            // Simulate refresh
            try? await Task.sleep(for: .seconds(1))
        }
        .tint(Color.appAccent) // Brand-colored refresh indicator
    }
}
```

### Subtle Parallax on Scroll

```swift
struct ParallaxHeader: View {
    var body: some View {
        ScrollView {
            GeometryReader { geo in
                let offset = geo.frame(in: .global).minY
                Image("header-photo")
                    .resizable()
                    .aspectRatio(contentMode: .fill)
                    .frame(height: 300 + max(0, -offset))
                    .offset(y: max(0, -offset) * 0.5) // Parallax factor
                    .clipped()
            }
            .frame(height: 300)
            
            // Content below
            VStack(spacing: 16) {
                Text("Content")
            }
            .padding()
        }
    }
}
```

### Tab Selection Animation

```swift
struct AnimatedTabBar: View {
    @State private var selectedTab = 0
    
    var body: some View {
        TabView(selection: $selectedTab) {
            HomeView()
                .tabItem { Label("Home", systemImage: "house") }
                .tag(0)
            
            SearchView()
                .tabItem { Label("Search", systemImage: "magnifyingglass") }
                .tag(1)
        }
        .tint(Color.appAccent)
    }
}
```

---

## SF Symbols Selection Guide

### Weight and Scale

SF Symbols have 9 weights: ultralight, light, thin, regular, medium, semibold, bold, heavy, black.

```swift
// Use weight to match text hierarchy
Image(systemName: "star.fill")
    .font(.system(size: 20, weight: .semibold)) // Match heading weight

// Use scale for size adjustment
Image(systemName: "star.fill")
    .imageScale(.small)   // 16pt
    .imageScale(.medium)  // 20pt (default)
    .imageScale(.large)   // 28pt
```

### Fill vs Outline

Use `.fill` variants for:
- Active/selected states
- Tab bar active icons
- Toggle states

Use outline (default) for:
- Inactive states
- Navigation bar buttons
- Decorative icons

```swift
// Active state
Image(systemName: selected ? "heart.fill" : "heart")
    .foregroundStyle(selected ? Color.appAccent : Color.appTextSecondary)
```

### Common SF Symbols by Category

| Category   | Symbols                                                                                    |
| ---------- | ------------------------------------------------------------------------------------------ |
| Navigation | `chevron.left`, `chevron.right`, `chevron.down`, `arrow.left`, `arrow.right`                 |
| Actions    | `plus`, `minus`, `pencil`, `trash`, `square.and.arrow.up`, `doc.on.doc`                      |
| Status     | `checkmark.circle.fill`, `xmark.circle.fill`, `exclamationmark.triangle.fill`, `info.circle` |
| Media      | `play.fill`, `pause.fill`, `forward.fill`, `backward.fill`, `speaker.wave.2.fill`             |
| People     | `person`, `person.fill`, `person.2`, `person.circle.fill`, `person.badge.plus`                |
| Objects    | `star`, `star.fill`, `heart`, `heart.fill`, `bell`, `bell.fill`, `gear`, `gearshape`          |

### Rendering Mode

```swift
// Monochrome (default) — uses foregroundStyle
Image(systemName: "star.fill")
    .foregroundStyle(Color.appAccent)

// Hierarchical — adds depth with opacity
Image(systemName: "star.fill")
    .symbolRenderingMode(.hierarchical)
    .foregroundStyle(Color.appAccent)

// Palette — multiple colors
Image(systemName: "star.fill")
    .symbolRenderingMode(.palette)
    .foregroundStyle(Color.appAccent, Color.appAccentLight)

// Multicolor — uses original symbol colors
Image(systemName: "star.fill")
    .symbolRenderingMode(.multicolor)
```

---

## Image & Asset Handling

### AsyncImage (Loading Remote Images)

```swift
struct RemoteImage: View {
    let url: URL?
    
    var body: some View {
        AsyncImage(url: url) { phase in
            switch phase {
            case .success(let image):
                image
                    .resizable()
                    .aspectRatio(contentMode: .fill)
            case .failure:
                // Fallback: clean placeholder
                Rectangle()
                    .fill(Color.appDivider)
                    .overlay {
                        Image(systemName: "photo")
                            .foregroundStyle(Color.appTextSecondary)
                    }
            case .empty:
                // Loading: subtle shimmer or progress
                Rectangle()
                    .fill(Color.appDivider)
                    .overlay {
                        ProgressView()
                    }
            @unknown default:
                EmptyView()
            }
        }
    }
}
```

### Placeholder Best Practices

```swift
// ✅ GOOD: Clean placeholder with icon
Rectangle()
    .fill(Color.appDivider)
    .overlay {
        Image(systemName: "photo")
            .font(.system(size: 32))
            .foregroundStyle(Color.appTextSecondary.opacity(0.5))
    }
    .clipShape(RoundedRectangle(cornerRadius: 12))

// ❌ BAD: Emoji placeholder
Text("🖼️")

// ❌ BAD: AI-generated SVG
// (inline SVG paths)
```

### Image Caching Strategy

For production apps, consider:
- `Kingfisher` (popular, well-maintained) — `KFImage(url)`
- `Nuke` (high performance)
- Or roll your own with `NSCache` + `AsyncImage`

Never load images synchronously on the main thread.
