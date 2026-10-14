# Sound Design (SFX Layer)

A sound-effect layer distinct from music playback: what to ship, how to play it with near-zero latency, and how to coexist with music — the host's or your own.

## Why sound is not optional for room-facing apps

Haptics reach only the holder; screen motion reaches only whoever is looking. When the phone is passed around, lying flat, or facing away (party games, timers, group activities), **audio is the only channel the room perceives**. A time-up moment that is haptic-only is *silent* to everyone but one person. If the app has any moment the room must notice, that moment needs a sound (or torch).

## The minimal kit

Resist a sound per event. Four to six mastered effects cover a whole app:

| Sound | Fires on | Character |
|-------|----------|-----------|
| Tick | Final-seconds countdown | Short, quiet, pitch-stable |
| Buzzer | Time up / failure | Unmistakable, ~0.5s |
| Ding | Success / correct | Bright, ~0.3s |
| Reveal | Content shown (card, wheel lands) | Soft whoosh/pop |
| Fanfare | Winner — the rarest event | The only "big" sound, 1–2s |

Master all effects to consistent perceived loudness; the fanfare may peak louder. Same rarity gradient as haptics: routine sounds small, rare sounds big.

## Asset format and latency

| Player | Latency | Use for | Caveats |
|--------|---------|---------|---------|
| `AVAudioPlayer` (preloaded + `prepareToPlay()`) | Low | The default for SFX | Allocate **once per sound at startup**, not per play |
| `AVAudioEngine` + `AVAudioPlayerNode` | Lowest | Overlapping/rapid-fire SFX, pitch variation | More setup; buffers stay scheduled |
| `AudioServicesPlaySystemSound` | Low | Fire-and-forget UI blips | No volume control, no session control — avoid for games |
| `AVPlayer` | High (streaming-class) | **Never** for SFX | Built for media, spins up per item |

- **Format: `.caf` containing PCM (or IMA4)** — decodes with effectively zero spin-up. Compressed formats (MP3/AAC) add decoder latency to the first play.
- Convert anything with the built-in tool: `afconvert input.wav output.caf -d ima4 -f caff`
- Keep effects short (≤2s) and small; they live in the app bundle, not remote config.

## Coexisting with music — the AVAudioSession decision

The session category is a *policy decision*; make it deliberately and write it down.

| Situation | Category / options | Result |
|-----------|-------------------|--------|
| SFX over **your own** in-app music | One `.playback` session; mix SFX at engine level | Music and SFX share the session you already own |
| SFX while **the user's** music plays (Apple Music, Spotify) | `.ambient` + `.mixWithOthers` | Effects layer over their music without stopping it |
| SFX must be *heard over* their music | add `.duckOthers` | Their music dips during your effect — use sparingly |
| Quiet utility app | `.ambient` (default respects silent switch) | Effects mute with the ringer switch |

#### ❌ The classic mistake
Activating a fresh `.playback` session just to play a ding — it **stops the user's background music**. If their music should keep playing, you want `.ambient` + `.mixWithOthers`, activated once, not per sound.

#### ⚠️ "Set the category once" fails when another component also sets it
If any other part of the app sets an *exclusive* category later (a music
engine claiming the room for its own playback is the common case), your
once-set mixing option is gone — and the next SFX play activates the stale
exclusive session, pausing the user's music. Field-proven fix: **re-assert
your category+options immediately before every play** (cheap), and never
call `setActive(false)` from the SFX layer — session activation belongs to
whichever component owns the room's audio.

#### Silent-switch policy
`.playback` ignores the ringer switch; `.ambient` respects it. For a game whose *product* is the room hearing the buzzer, playing through the switch during an active session is the defensible choice (players opted into a loud activity) — but pair it with an in-app mute. For everything else, respect the switch. Either way: **document the choice** where reviewers will look, because it *will* be questioned.

## User control

- **In-app mute toggle for SFX**, independent of the haptics toggle and of music volume. Store it in the sound service; check it at fire time.
- If music and SFX both exist, two volumes (or music/SFX toggles) beat one master switch.
- Never gate SFX on Reduce Motion — unrelated setting. There is no system "reduce sounds" to honor; your toggle *is* the accessibility story.

## VoiceOver interplay

Game-critical announcements (whose turn, what happened) should be **VoiceOver announcements**, not just SFX — a buzzer conveys *that* something happened, never *what*. Post `AccessibilityNotification.Announcement` alongside the sound. If effects talk over speech, add `.interruptSpokenAudioAndMixWithOthers` to the options and test with VoiceOver live.

## Service shape

Mirror the haptics architecture (haptic-design.md): one `@Observable` service, semantic events (often the *same* enum the haptics service uses — one `fire(event:)` fanning out to both), preloaded players keyed by event, capability/no-op seams for tests, mute state persisted. Game code fires meaning; the service decides channels.

## Pitfalls

#### ❌ Allocating `AVAudioPlayer` per play — first-play latency and dropped sounds under rapid fire
#### ❌ MP3 assets for latency-critical effects — decoder spin-up on first play
#### ❌ Killing the user's music with a needless `.playback` activation
#### ❌ No mute toggle — the only recourse is deleting the app
#### ❌ SFX as the sole conveyance of game state — pair with VoiceOver announcements
#### ✅ Preload `.caf` PCM at startup, one session policy chosen deliberately, mute toggle, rarity-scaled loudness

## References

- [HIG — Playing audio](https://developer.apple.com/design/human-interface-guidelines/playing-audio)
- [AVAudioSession category/options](https://developer.apple.com/documentation/avfaudio/avaudiosession)
- [AVAudioEngine](https://developer.apple.com/documentation/avfaudio/avaudioengine)
