# SwiftUI Layout Patterns

Practical layout patterns for common SwiftUI design scenarios. Each pattern includes code and when to use it.

## Core Layout Principles

1. **Vertical flow first** — SwiftUI naturally flows top-to-bottom
2. **Semantic grouping** — Group related content, not just visual elements
3. **Progressive disclosure** — Show summary first, details on demand
4. **Content-driven sizing** — Let content determine heights, not arbitrary frames

## Pattern 1: Screen Structure

Every screen follows this skeleton:

```swift
struct ScreenView: View {
    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: DesignTokens.spaceL) {
                // 1. Header (small, tasteful)
                ScreenHeader(title: "Screen Title")
                
                // 2. Primary content
                PrimaryContentSection()
                
                // 3. Secondary content
                SecondaryContentSection()
            }
            .padding(DesignTokens.spaceM)
        }
        .background(DesignTokens.background)
    }
}
```

**Key rules:**
- Never use `GeometryReader` for basic layout — use `containerRelativeFrame` or `frame(maxWidth: .infinity)`
- Always wrap scrollable content in `ScrollView`
- Use `.scrollTargetBehavior(.paging)` for paged scrolling
- Use `.scrollClipDisabled()` if content needs to overflow scroll bounds

## Pattern 2: Card Layout

```swift
struct CardView<Content: View>: View {
    let content: Content
    
    init(@ViewBuilder content: () -> Content) {
        self.content = content()
    }
    
    var body: some View {
        content
            .padding(DesignTokens.spaceM)
            .frame(maxWidth: .infinity, alignment: .leading)
            .background(DesignTokens.surface)
            .clipShape(RoundedRectangle(cornerRadius: DesignTokens.radiusM))
            .shadow(color: .black.opacity(0.04), radius: 8, y: 2)
    }
}
```

**Key rules:**
- Cards should have subtle shadow, not heavy borders
- Corner radius should be consistent (12-16pt for cards)
- Always use `clipShape(RoundedRectangle(...))` not `.cornerRadius()`
- Content inside cards should have clear hierarchy

## Pattern 3: List Layout

```swift
struct ItemListView: View {
    let items: [Item]
    
    var body: some View {
        LazyVStack(spacing: 0) {
            ForEach(items) { item in
                ItemRow(item: item)
                
                if item.id != items.last?.id {
                    Divider()
                        .padding(.leading, DesignTokens.spaceM)
                }
            }
        }
        .background(DesignTokens.surface)
        .clipShape(RoundedRectangle(cornerRadius: DesignTokens.radiusM))
    }
}
```

**Key rules:**
- Use `LazyVStack` for moderate lists, `List` for very long lists
- Separators should be inset, not full-width
- Group related items visually (with section headers, not just dividers)

## Pattern 4: Grid Layout

```swift
// 2-column adaptive grid
let columns = [
    GridItem(.adaptive(minimum: 160), spacing: DesignTokens.spaceM)
]

LazyVGrid(columns: columns, spacing: DesignTokens.spaceM) {
    ForEach(items) { item in
        GridItemView(item: item)
    }
}
```

**Key rules:**
- Use `GridItem(.adaptive(minimum:))` for responsive grids
- Avoid fixed column counts — let content determine layout
- Never do symmetric 3-column feature grids (anti-slop rule)
- Grid items should have consistent heights within a row

## Pattern 5: Header + Content

```swift
struct ScreenHeader: View {
    let title: String
    let subtitle: String?
    
    var body: some View {
        VStack(alignment: .leading, spacing: DesignTokens.spaceXS) {
            Text(title)
                .font(.system(size: 28, weight: .bold, design: .serif))
            
            if let subtitle {
                Text(subtitle)
                    .font(.body)
                    .foregroundStyle(.secondary)
            }
        }
        .padding(.bottom, DesignTokens.spaceS)
    }
}
```

**Key rules:**
- Headers should be left-aligned (not centered) for most apps
- Use serif display font for main titles
- Keep headers compact — don't waste vertical space
- Large title navigation: use `.navigationBarTitleDisplayMode(.large)` only for root views

## Pattern 6: Tab Bar Layout

```swift
TabView {
    HomeView()
        .tabItem {
            Label("Home", systemImage: "house")
        }
    
    SearchView()
        .tabItem {
            Label("Search", systemImage: "magnifyingglass")
        }
    
    ProfileView()
        .tabItem {
            Label("Profile", systemImage: "person")
        }
}
.tint(DesignTokens.accent)
```

**Key rules:**
- 3-5 tabs maximum
- Use SF Symbols for tab icons
- Active tab uses the accent color
- Don't put settings in tabs — use a gear icon in the navigation bar

## Pattern 7: Modal / Sheet

```swift
.sheet(isPresented: $showSheet) {
    NavigationStack {
        SheetContent()
            .navigationTitle("Details")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Cancel") { showSheet = false }
                }
            }
    }
    .presentationDetents([.medium, .large])
    .presentationDragIndicator(.visible)
}
```

**Key rules:**
- Always wrap sheets in `NavigationStack`
- Use `.presentationDetents` for half-sheets
- Provide a clear dismiss action
- Use `.presentationDragIndicator(.visible)` for discoverability

## Pattern 8: Empty State

```swift
struct EmptyStateView: View {
    let icon: String
    let title: String
    let message: String
    
    var body: some View {
        ContentUnavailableView {
            Label(title, systemImage: icon)
        } description: {
            Text(message)
        }
    }
}
```

**Key rules:**
- Use `ContentUnavailableView` (iOS 17+) for empty states
- Don't use emoji — use SF Symbols
- Keep messaging concise and helpful
- Include a call-to-action if appropriate

## Pattern 9: Loading State

```swift
struct LoadingView: View {
    var body: some View {
        ProgressView()
            .frame(maxWidth: .infinity, maxHeight: .infinity)
            .background(DesignTokens.background)
    }
}
```

**Key rules:**
- Use `ProgressView()` for simple loading
- For skeleton screens, use `.redacted(reason: .placeholder)`
- Never show a blank screen during loading
- Consider shimmer effects for premium feel

## Pattern 10: Responsive Layout (iPhone + iPad)

```swift
struct ResponsiveView: View {
    @Environment(\.horizontalSizeClass) var sizeClass
    
    var body: some View {
        if sizeClass == .compact {
            // iPhone layout
            VStack { ... }
        } else {
            // iPad layout
            HStack { ... }
        }
    }
}
```

**Key rules:**
- Use `horizontalSizeClass` for layout adaptation
- Don't just scale up — use the extra space meaningfully
- iPad should show more content, not bigger content
- Use `NavigationSplitView` for multi-column on iPad
