# Haptic Design

Designing a haptic *vocabulary* — semantic, distinct, and centrally owned — rather than sprinkling generator calls through views.

## API Decision Table

| Need | API | Notes |
|------|-----|-------|
| Simple state-driven feedback in SwiftUI | `.sensoryFeedback(_:trigger:)` (iOS 17+) | Declarative; fine for `.success`, `.selection`, `.impact` on value changes |
| One-off imperative feedback, UIKit-style | `UIImpactFeedbackGenerator` etc. | Call `prepare()` first; limited to system presets |
| A custom vocabulary of distinct patterns | **Core Haptics** (`CHHapticEngine`) | The right choice once you have ≥4 semantic events; full intensity/sharpness control |
| Designer-authored waveforms | AHAP files + `CHHapticPatternPlayer` | Only worth it for signature audio-haptic sync moments |

Rule of thumb: the moment you catch yourself choosing between `.success` and `.impact(weight:)` for *game* events, you've outgrown presets — build a Core Haptics vocabulary.

## The Semantic-Vocabulary Architecture

One central service owns the engine; game code speaks in *meanings*, never in patterns:

```swift
enum HapticEvent {
    case start          // round/session begins
    case correct        // success tap
    case pass           // skip / failure tap
    case reveal         // content shown (card, letter, dare)
    case timerPulse     // gentle urgency (e.g. 10s left)
    case timerUrgent    // hard urgency (e.g. 5s left)
    case stop           // the signature interrupt (music stop, freeze, catch)
    case roundEnd       // round complete, score committed
    case winner         // the rarest, biggest pattern
}

@MainActor @Observable
final class HapticsService {
    private var engine: CHHapticEngine?
    private let supportsHaptics =
        CHHapticEngine.capabilitiesForHardware().supportsHaptics

    func fire(_ event: HapticEvent) {
        guard supportsHaptics else { return }          // silent no-op on Simulator/iPad
        startEngineIfNeeded()                          // lazy start on first use
        // build CHHapticPattern from events(for: event), play it
    }
}
```

Architecture rules (each prevents a real failure):

- **Lazy engine start** on first `fire`, not in `init` — engine start is not free and may race app launch.
- **Capability-gate once** (`capabilitiesForHardware().supportsHaptics`) so every call site can fire blindly; no availability checks leak into views.
- **Handle `stoppedHandler` / `resetHandler`** by nulling the engine so the *next* fire restarts it — the system stops engines freely (audio interruptions, thermal).
- **Stop the engine on background** (`scenePhase` change) and let lazy-start recover on foreground — a running engine in background wastes power and can be killed out from under you.
- **Inject as an optional seam** (`haptics: HapticsService? = nil` in stores/engines) — unit tests pass `nil`; previews pass a throwaway.
- Views get it via `@Environment`; only *stores/engines* fire game events (see single-owner rule in SKILL.md).

## Vocabulary Design Rules

1. **Name events by meaning, not by pattern** — `.correct`, not `.doubleTapLight`. The pattern can be retuned later without touching call sites.
2. **Signature moments get signature patterns.** The freeze/stop moment of a music game, the spin of a wheel, the winner — each needs its *own* event. Reusing `.roundEnd` for the winner erases the climax (rarity gradient, SKILL.md).
3. **Escalation tiers are pattern *pairs*** — `.timerPulse` and `.timerUrgent` should be the same shape at different intensity/sharpness so they read as "the same thing, more urgent", and their thresholds must match the visual escalation exactly.
4. **Continuous events for continuous moments.** A spin deserves a `.hapticContinuous` with an intensity ramp (parameter curves), not repeated transients. Transients are for instants.
5. **Cover every event in the domain, then stop.** Draw up the event×channel matrix (feedback-audit.md) first; if an event isn't in the matrix, it doesn't get a haptic.

## Pattern Sketches (starting points, tune on device)

| Event | Shape |
|-------|-------|
| `.start` | One medium transient (intensity 0.6, sharpness 0.4) — "we've begun" |
| `.correct` | One crisp transient (0.8, 0.7) |
| `.pass` | One soft dull transient (0.5, 0.25) |
| `.timerPulse` | Light transient (0.4, 0.3) each second |
| `.timerUrgent` | Same cadence, harder (0.8, 0.6) |
| `.stop` | Heavy transient + 0.3s continuous rumble tail — unmistakable |
| `.roundEnd` | Two transients 150ms apart |
| `.winner` | Three ascending transients + continuous swell — the biggest thing you play |

## Pitfalls

#### ❌ Double-fire
Engine fires on "music stopped" *and* the store fires on the same state change → user feels a stutter-buzz. One owner per event.

#### ❌ Preset semantics for game events
`UINotificationFeedbackGenerator().notificationOccurred(.success)` for a winner — it's the same buzz every form in iOS uses. Signature moments deserve custom patterns.

#### ❌ No user toggle
Core Haptics ignores any system setting — there is no global "haptics off" that covers you. Ship a settings toggle; store it in the service, not in every call site.

#### ❌ Testing feel in the Simulator
The Simulator plays **no haptics at all**. Pattern *code* can be unit-tested (event → pattern mapping), but intensity/sharpness tuning is device-only work — schedule it.

#### ✅ Do
Centralize; speak semantics; gate once; tune the rarest events to be the biggest; give users an off switch.

## References

- [HIG — Playing haptics](https://developer.apple.com/design/human-interface-guidelines/playing-haptics)
- [Core Haptics documentation](https://developer.apple.com/documentation/corehaptics)
- [CHHapticEngine](https://developer.apple.com/documentation/corehaptics/chhapticengine)
