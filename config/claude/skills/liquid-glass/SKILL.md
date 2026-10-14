---
name: liquid-glass
description: Implement Liquid Glass design using .glassEffect() API for iOS/macOS 26+. Use when creating modern glass-based UI effects.
allowed-tools: [Read, Write, Edit, Glob, Grep, AskUserQuestion]
---

# Liquid Glass Design

Implement Apple's Liquid Glass design language using the modern `.glassEffect()` API.

## When to Use

- User wants glass/blur effects on views
- User asks about Liquid Glass or modern Apple design
- User needs transparent, interactive UI elements
- User wants morphing transitions between views

## Quick Start (SwiftUI)

### Basic Glass Effect

```swift
import SwiftUI

Text("Hello, World!")
    .font(.title)
    .padding()
    .glassEffect()  // Capsule shape by default
```

### Custom Shape

```swift
Text("Hello")
    .padding()
    .glassEffect(in: .rect(cornerRadius: 16))

// Available shapes:
// .capsule (default)
// .rect(cornerRadius: CGFloat)
// .circle
```

### Interactive Glass

```swift
Button("Tap Me") {
    // action
}
.padding()
.glassEffect(.regular.interactive())
```

### Tinted Glass

```swift
Text("Important")
    .padding()
    .glassEffect(.regular.tint(.blue))
```

## Glass Configuration Options

| Option | Description | Example |
|--------|-------------|---------|
| `.regular` | Standard glass effect | `.glassEffect(.regular)` |
| `.tint(Color)` | Add color tint | `.glassEffect(.regular.tint(.orange))` |
| `.interactive()` | React to touch/hover | `.glassEffect(.regular.interactive())` |

## Multiple Glass Effects

### GlassEffectContainer

When using multiple glass elements, wrap them in `GlassEffectContainer` for:
- Better rendering performance
- Proper blending between effects
- Morphing transitions

```swift
GlassEffectContainer(spacing: 40.0) {
    HStack(spacing: 40.0) {
        Image(systemName: "star.fill")
            .frame(width: 80, height: 80)
            .font(.system(size: 36))
            .glassEffect()

        Image(systemName: "heart.fill")
            .frame(width: 80, height: 80)
            .font(.system(size: 36))
            .glassEffect()
    }
}
```

**Spacing Parameter:**
- Controls when effects merge
- Smaller spacing = views must be closer to merge
- Larger spacing = effects merge at greater distances

### Uniting Glass Effects

Combine views into a single glass effect using `glassEffectUnion`:

```swift
@Namespace private var namespace

GlassEffectContainer(spacing: 20.0) {
    HStack(spacing: 20.0) {
        ForEach(items.indices, id: \.self) { index in
            Image(systemName: items[index])
                .frame(width: 60, height: 60)
                .glassEffect()
                .glassEffectUnion(
                    id: index < 2 ? "group1" : "group2",
                    namespace: namespace
                )
        }
    }
}
```

## Morphing Transitions

Create fluid morphing effects when views appear/disappear.

### Setup

1. Create a namespace
2. Assign glass effect IDs
3. Use animations on state changes

```swift
struct MorphingToolbar: View {
    @State private var isExpanded = false
    @Namespace private var namespace

    var body: some View {
        GlassEffectContainer(spacing: 40.0) {
            HStack(spacing: 40.0) {
                // Always visible
                Image(systemName: "pencil")
                    .frame(width: 60, height: 60)
                    .glassEffect()
                    .glassEffectID("pencil", in: namespace)

                // Conditionally visible - will morph in/out
                if isExpanded {
                    Image(systemName: "eraser")
                        .frame(width: 60, height: 60)
                        .glassEffect()
                        .glassEffectID("eraser", in: namespace)

                    Image(systemName: "ruler")
                        .frame(width: 60, height: 60)
                        .glassEffect()
                        .glassEffectID("ruler", in: namespace)
                }
            }
        }

        Button("Toggle") {
            withAnimation(.spring(response: 0.4, dampingFraction: 0.8)) {
                isExpanded.toggle()
            }
        }
        .buttonStyle(.glass)
    }
}
```

## Button Styles

### Glass Button

```swift
Button("Standard") {
    // action
}
.buttonStyle(.glass)
```

### Glass Prominent Button

```swift
Button("Primary Action") {
    // action
}
.buttonStyle(.glassProminent)
```

## Advanced Techniques

### Background Extension

Stretch content under sidebar or inspector:

```swift
NavigationSplitView {
    SidebarView()
} detail: {
    DetailView()
        .background {
            Image("wallpaper")
                .resizable()
                .ignoresSafeArea()
        }
}
```

### Horizontal Scroll Under Sidebar

```swift
ScrollView(.horizontal) {
    HStack {
        ForEach(items) { item in
            ItemView(item: item)
        }
    }
}
.scrollExtensionMode(.underSidebar)
```

## AppKit Implementation

### NSGlassEffectView

```swift
import AppKit

// Create glass effect view
let glassView = NSGlassEffectView(frame: NSRect(x: 20, y: 20, width: 200, height: 100))
glassView.cornerRadius = 16.0
glassView.tintColor = NSColor.systemBlue.withAlphaComponent(0.3)

// Create content
let label = NSTextField(labelWithString: "Glass Content")
label.translatesAutoresizingMaskIntoConstraints = false

// Set content view
glassView.contentView = label

// Add constraints
if let contentView = glassView.contentView {
    NSLayoutConstraint.activate([
        label.centerXAnchor.constraint(equalTo: contentView.centerXAnchor),
        label.centerYAnchor.constraint(equalTo: contentView.centerYAnchor)
    ])
}
```

### NSGlassEffectContainerView

```swift
// Create container
let container = NSGlassEffectContainerView(frame: bounds)
container.spacing = 40.0

// Create content view
let contentView = NSView(frame: container.bounds)
container.contentView = contentView

// Add glass views to content
let glass1 = NSGlassEffectView(frame: NSRect(x: 20, y: 50, width: 150, height: 100))
let glass2 = NSGlassEffectView(frame: NSRect(x: 190, y: 50, width: 150, height: 100))

contentView.addSubview(glass1)
contentView.addSubview(glass2)
```

### Interactive AppKit Glass

```swift
class InteractiveGlassView: NSGlassEffectView {
    override init(frame: NSRect) {
        super.init(frame: frame)
        setupTracking()
    }

    required init?(coder: NSCoder) {
        super.init(coder: coder)
        setupTracking()
    }

    private func setupTracking() {
        let options: NSTrackingArea.Options = [
            .mouseEnteredAndExited,
            .activeInActiveApp
        ]
        let trackingArea = NSTrackingArea(
            rect: bounds,
            options: options,
            owner: self,
            userInfo: nil
        )
        addTrackingArea(trackingArea)
    }

    override func mouseEntered(with event: NSEvent) {
        super.mouseEntered(with: event)
        NSAnimationContext.runAnimationGroup { context in
            context.duration = 0.2
            animator().tintColor = NSColor.systemBlue.withAlphaComponent(0.2)
        }
    }

    override func mouseExited(with event: NSEvent) {
        super.mouseExited(with: event)
        NSAnimationContext.runAnimationGroup { context in
            context.duration = 0.2
            animator().tintColor = nil
        }
    }
}
```

## Common Patterns

### Floating Action Bar

```swift
struct FloatingActionBar: View {
    @Namespace private var namespace

    var body: some View {
        GlassEffectContainer(spacing: 20) {
            HStack(spacing: 16) {
                ForEach(actions) { action in
                    Button {
                        action.perform()
                    } label: {
                        Image(systemName: action.icon)
                            .font(.title2)
                    }
                    .frame(width: 44, height: 44)
                    .glassEffect(.regular.interactive())
                    .glassEffectID(action.id, in: namespace)
                }
            }
            .padding(.horizontal, 8)
            .padding(.vertical, 4)
        }
    }
}
```

### Card with Glass Effect

```swift
struct GlassCard: View {
    let title: String
    let subtitle: String
    let icon: String

    var body: some View {
        HStack(spacing: 16) {
            Image(systemName: icon)
                .font(.title)
                .frame(width: 50, height: 50)
                .glassEffect(.regular.tint(.blue))

            VStack(alignment: .leading) {
                Text(title)
                    .font(.headline)
                Text(subtitle)
                    .font(.subheadline)
                    .foregroundStyle(.secondary)
            }

            Spacer()
        }
        .padding()
        .glassEffect(in: .rect(cornerRadius: 16))
    }
}
```

### Tab Bar with Morphing

```swift
struct GlassTabBar: View {
    @Binding var selection: Int
    @Namespace private var namespace

    let tabs = [
        ("house", "Home"),
        ("magnifyingglass", "Search"),
        ("person", "Profile")
    ]

    var body: some View {
        GlassEffectContainer(spacing: 30) {
            HStack(spacing: 30) {
                ForEach(tabs.indices, id: \.self) { index in
                    Button {
                        withAnimation(.spring(response: 0.35, dampingFraction: 0.8)) {
                            selection = index
                        }
                    } label: {
                        VStack(spacing: 4) {
                            Image(systemName: tabs[index].0)
                                .font(.title2)
                            Text(tabs[index].1)
                                .font(.caption)
                        }
                        .frame(width: 70, height: 60)
                    }
                    .glassEffect(
                        selection == index
                            ? .regular.tint(.blue).interactive()
                            : .regular.interactive()
                    )
                    .glassEffectID("tab\(index)", in: namespace)
                }
            }
        }
    }
}
```

## Migration from Old API

### Before (Old Approach)

```swift
// Old: Using materials directly
VStack {
    Text("Content")
}
.padding()
.background(.ultraThinMaterial)
.cornerRadius(16)
```

### After (New API)

```swift
// New: Using glassEffect modifier
VStack {
    Text("Content")
}
.padding()
.glassEffect(in: .rect(cornerRadius: 16))
```

### Key Differences

| Old Approach | New API |
|--------------|---------|
| `.background(.material)` | `.glassEffect()` |
| Manual corner radius | Shape parameter |
| No interactivity | `.interactive()` modifier |
| Manual tinting | `.tint(Color)` modifier |
| No morphing | `glassEffectID` + `@Namespace` |
| No container grouping | `GlassEffectContainer` |

## Best Practices

1. **Use GlassEffectContainer** for multiple glass views
   - Improves rendering performance
   - Enables morphing transitions

2. **Apply glass effect last** in modifier chain
   - After frame, padding, and content modifiers

3. **Choose appropriate spacing** in containers
   - Controls when effects blend together

4. **Use animations** for state changes
   - Enables smooth morphing transitions

5. **Add interactivity** for touchable elements
   - `.interactive()` for buttons and controls

6. **Tint strategically** to indicate state
   - Selected items, primary actions

7. **Consistent shapes** across your app
   - Establish a shape language (all capsules, or all rounded rects)

## Checklist

- [ ] Use `.glassEffect()` instead of `.background(.material)`
- [ ] Wrap multiple glass views in `GlassEffectContainer`
- [ ] Add `@Namespace` for morphing transitions
- [ ] Use `.glassEffectID()` on views that appear/disappear
- [ ] Add `.interactive()` for touchable elements
- [ ] Use `.buttonStyle(.glass)` for glass buttons
- [ ] Test animations for smooth morphing
- [ ] Consider performance with many glass effects
- [ ] Support both light and dark appearances

## References

- [Applying Liquid Glass to custom views](https://developer.apple.com/documentation/SwiftUI/Applying-Liquid-Glass-to-custom-views)
- [Landmarks: Building an app with Liquid Glass](https://developer.apple.com/documentation/SwiftUI/Landmarks-Building-an-app-with-Liquid-Glass)
- [SwiftUI GlassEffectContainer](https://developer.apple.com/documentation/SwiftUI/GlassEffectContainer)
- [AppKit NSGlassEffectView](https://developer.apple.com/documentation/AppKit/NSGlassEffectView)
