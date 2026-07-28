---
name: ableton
description: Control Ableton Live through the AbletonMCP tools. Use when the user wants to build tracks, patches, clips, drums, automation, or generative music in Ableton, or asks to change any device, mixer, or transport setting in their Live set. Triggers on "Ableton", "Live set", "make a patch", "build a track", "add a synth", "write a clip", "sidechain", "automate", or any request to change how something sounds in their DAW.
---

# Controlling Ableton Live

You are driving a real Live set that a musician cares about. Everything you do
lands immediately in their project. Act accordingly.

## The rule that governs everything: you cannot hear

You have no audio input. You can set a filter to 0.37 and confirm the value
took, but you have no idea what came out of the speakers.

So never tell the user something sounds good, warm, punchy, or lush. You don't
know. Describe what you built and what you verified, then let them judge.

Where you're genuinely useful is the tedious, mechanical work: 200 automation
breakpoints, coprime loop lengths, phase offsets across six parameters,
session scaffolding. Where you're useless is taste. Don't pretend otherwise.

## Verify everything. Status codes lie.

A `success` response does not mean the change happened. Always read the value
back.

This is not paranoia. Several commands return success while doing nothing, or
while doing something different from what was asked. Read back after every
change that matters:

```python
# set it
cmd("set_device_parameter", {"track_index": 0, "device_index": 1,
                             "parameter_index": 5, "value": 0.42})
# then PROVE it
ps = cmd("get_device_parameters", {"track_index": 0, "device_index": 1})
```

If a read-back disagrees with what you asked for, say so plainly rather than
reporting success.

## Read the device before you set it

Never guess parameter names, indices, or ranges. Every device differs, and
similar devices differ in surprising ways.

Call `get_device_parameters` first. It returns `name`, `value`, `min`, `max`,
and `value_string`. That last field is the text Live prints on the knob
("Off", "1/16", "440 Hz"), and it's the only way to know what a bare float
like `Transpose Mode = 0.0` actually means.

Ranges are not always 0 to 1. `Gate` on the Arpeggiator runs 1 to 200.
`Transpose` on Drift runs -48 to 48. `Depth` on the M4L LFO runs 0 to 100.
Read first.

## Dropdowns are unreachable. Ask the user to click.

Live doesn't expose routing choices as automatable parameters. If it's a
dropdown or a Map button in the UI, the API cannot set it. Known cases:

The Max for Live LFO's Map button. Every other knob on it (Rate, Depth,
Jitter, Smooth, Phase, Shape) is yours, but the target isn't. Whenever you use a
Max for Live LFO to modulate a parameter, you MUST stop and ask the user to
click Map and pick the target. Do not proceed as if it is done, and do not
pretend a workaround makes it unnecessary. Set every knob on the LFO first, then
say exactly which knob to Map onto (device and parameter name). The click is
once per LFO; after it, you own the rest.

A Compressor's sidechain Audio From source. You can set `S/C On`, threshold,
ratio, attack, release, and the whole sidechain EQ. You cannot pick the source
track.

Drift's modulation matrix source and destination selectors. The amounts
(`Mod Matrix Amt 1`, `LP Mod Amt 1`) are settable; the routing is not.

Grouping tracks (Cmd+G). The Live API cannot group existing tracks: it only
reads group membership, never creates a group. Do not promise grouping. When
the goal is a shared effect across tracks, which is the usual reason to group,
make a return track instead: create_return_track, load the effect on it, then
raise each source track's send. That is fully scriptable and is the right home
for a send effect like a shared echo or reverb. Only a real group needs the user
to select the tracks and press Cmd+G.

Don't fake your way around these. Set everything you can, then tell the user
exactly which control to click, naming the device and the parameter.

## Free-running modulation of an arbitrary parameter needs the Map click

There is no combination of tools that is free-running AND aims at any parameter
AND needs no click. Be honest about which you are giving up:

A native device LFO (Auto Filter's LFO to cutoff, Auto Pan) is free-running and
fully scriptable, but only reaches its own fixed target. Reach for these first
when the target is one they cover, because they need no click.

`add_shaped_automation` draws a waveform onto any parameter, but it is clip
automation: it repeats every loop and is phase-locked to the bar. It is not
free-running. A long or coprime clip length hides the repetition but does not
remove it. Use it when a loop-locked, editable curve is acceptable.

The Max for Live LFO is the only way to get free-running modulation of an
arbitrary parameter, and it requires the one Map click. If the user wants that
combination, ask for the click. Never imply `add_shaped_automation` is
equivalent to a real LFO; it loops, a real LFO does not.

## Parameter names that will bite you

These commands silently defaulted on a wrong argument name in older versions,
and the names are not what you'd guess:

`set_track_volume` wants `volume`, not `value`.
`set_track_pan` wants `pan`, not `value`.
`set_send_level` wants `level`, not `value`.
`load_browser_item` wants `item_uri`, not `uri`.
`load_browser_item_to_return` wants `return_index`, not `return_track_index`.
`set_clip_automation` wants `parameter_name` and `envelope_data`, not indices
and points.

The MCP tool layer uses the correct names. If you're calling the socket
directly, check the handler signature first.

## Note probability: use the right command

`add_notes_to_clip` uses Live's legacy 5-tuple API and cannot carry
probability. It will silently drop it.

For anything generative, use `add_notes_with_probability`. Each note takes
`pitch`, `start_time`, `duration`, `velocity`, `mute`, plus optional
`probability` (0 to 1) and `velocity_deviation`. It returns
`verified_first_note` so you can prove probability actually landed.

Weight probability rather than randomising it flat. Structural notes (bass
roots, the downbeat kick) belong at 0.9 to 1.0 so the music keeps its skeleton.
Decorative notes (hats, melody, arp tones) sit at 0.3 to 0.7 so the texture
shifts. Uniform random probability everywhere reads as broken, not generative.

## Reaching inside racks

`get_device_parameters` on a drum rack or instrument rack returns only the
rack's own macros, which are usually unassigned and useless. The devices
inside are invisible to it.

Use `get_rack_chains` to list the chains, then `get_chain_device_parameters`
and `set_chain_device_parameter` with `chain_index` and `chain_device_index`
to reach the actual device. That's how you set a single drum pad's decay.

Note that MIDI note length does nothing on a drum rack. Drum racks are one-shot
triggers. If a drum hit is too long, the fix is the pad's decay envelope inside
the rack, never the note duration in the clip.

## Clip automation targets by name. Be specific.

`set_clip_automation` matches a parameter name across devices on the track.
Names collide constantly: `Frequency` exists on Auto Filter, Erosion, Reverb,
and Grain Delay simultaneously.

Always pass `device_index` to scope it. Without it, an ambiguous name will
raise (good) or, in older versions, silently automate the wrong device (bad).

Verify with `get_clip_automation`, which samples the envelope and returns the
real curve plus `is_flat`. If `is_flat` is true, your automation didn't land.

## Generative patches: what actually works

Coprime loop lengths are the strongest tool you have. Clips at 3, 5, 7, 11,
13, and 17 bars won't realign for over a million bars. Same clips, endlessly
shifting relationship, no automation needed.

Unsynced LFOs drift against the grid. Set `LFO T Mode` to free (not synced)
and a very slow `LFO Freq`. On the M4L LFO, `Jitter` makes the period itself
wander, so it never repeats even against itself. A synced LFO is just a loop.

Drawn automation is the opposite of generative. It repeats identically every
pass. Reach for LFOs and probability instead when the goal is "never the same
twice."

Live's `Random` MIDI effect has `Chance` at 0 by default, so it does nothing
until you raise it. Pair it with the `Scale` device to keep scattered pitches
in key. Set Scale's `Map 0` through `Map 11` explicitly rather than using
`Use Current Scale`, which depends on the song scale and will surprise you.

## Sound design facts worth knowing

Note length is inaudible when sustain is near zero. Every note decays to
silence immediately, so a 4 beat note and a half beat note sound identical.
If the user asks for varied note lengths, raise sustain first or nothing
changes.

A pure tone plus high resonance plus a long release is a bell, by definition.
If the user says something sounds like a bell and doesn't want that, remove the
pitch (use noise) or the resonance. You can't fix it by tweaking the release.

Big and staccato are opposites. Length and reverb make things big; both destroy
staccato. Don't try to deliver both from the same parameters.

Reverb is usually the reason something isn't staccato, not note length. Check
both the track's own reverb and its return sends.

Timing is what makes something sound digital. Notes on exact grid positions
read as machine-made no matter how much you vary velocity. Offset each note by
a small random amount, biased late (players drag behind the beat far more often
than they rush), and it reads as human.

An arpeggiator discards incoming note timing and emits at its own fixed rate.
Editing clip note lengths on an arp track does nothing. Automate `Gate` if you
want note length to vary.

## Safety

Ask before destructive work. Deleting tracks, clearing a set, or overwriting
clips the user built by hand deserves a confirmation.

Live won't delete the last remaining track. Create the new one first, then
delete the old.

There's no command for a new Live set. That's File, New Live Set, by hand.

Prefer working in Session view for anything that loops or phases. Arrangement
view is a fixed timeline, and launching a scene restarts every clip, which
resets any phasing you set up.
