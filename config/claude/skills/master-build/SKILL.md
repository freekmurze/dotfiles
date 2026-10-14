---
name: master-build
description: Keep a stable, installed build of a desktop app that the owner can use and test while agents keep working in the repository. Use when agents building or launching the app are getting in the owner's way, when the working tree is often mid edit and does not compile, when asked for a "master build", or when the owner says they want to use the app themselves while work continues.
---

# Master build

A build the owner uses, decoupled from the repository agents are working in.

## The problem this solves

When several agents work on a desktop app at once, the app the owner is trying to use stops
being usable, for three separate reasons. Fixing one without the others does not help.

1. **Agents launch instances.** Windows appear, steal activation, and sometimes take over the
   owner's own app when the bundle id matches. `set frontmost` on an agent's build resolves to
   whatever is frontmost, not the process it addressed.
2. **The working tree is usually mid edit.** With agents committing continuously, at any given
   moment the tree may not compile. A build taken from the tree is hostage to the second it was
   taken in.
3. **A shared build directory races.** Two agents building into one scratch path produce a binary
   that is a mixture of two states. Symptoms are worse than a failure: it runs, and it is wrong.

## The shape of the fix

Three rules, and an install script.

**Agents never launch the app.** State it as an absolute in every brief, not a preference. They
verify by building, by running the test suite, and by reading. Anything that genuinely needed a
window goes in their report under "could not verify", with a note on what they would have checked.
That is the expected outcome, not a failure. The owner tests when they are ready.

**Agents that must run something re-identify it.** Their own bundle id, their own scratch build
path, their own copy of any database, their own preferences domain, deleted afterwards. Never the
real bundle id, because that is what makes an agent's automation land on the owner's window.

**The owner's build comes from a commit, not from the tree.** This is the load bearing rule.
A detached worktree at a ref, built in its own scratch path, installed to a stable location.
Nothing uncommitted comes along and no concurrent build can swap objects underneath it.

## The script

Live example: `Tools/master.sh` in the Bloom repository, also reachable there as `make master`.
Adapt the build command; the structure is the point.

```sh
REF=${1:-HEAD}
RESOLVED="$(git rev-parse --short "$REF")"

git worktree add --detach /tmp/app-master-src "$RESOLVED"
( cd /tmp/app-master-src && <build release, own scratch path> )

rm -rf "$HOME/Applications/App.app"
cp -R "<built bundle>" "$HOME/Applications/App.app"
git worktree remove --force /tmp/app-master-src
```

Then, only when relaunching, kill **by full path** so a debug build and any agent's build are
left alone:

```sh
pkill -f "$HOME/Applications/App.app/Contents/MacOS/App"
open "$HOME/Applications/App.app"
```

## Details that matter

**Keep the real bundle id.** The point is that the owner's database, preferences and window state
carry over. Agents are the ones who re-identify, not this build.

**Stamp the commit into the bundle**, so a stale install can be identified rather than guessed at:

```sh
defaults write "$DEST/Contents/Info.plist" MasterCommit -string "$RESOLVED"
```

**Build release.** Not for speed (measured on one SwiftUI app: 34.9ms against 35.5ms per frame,
a two percent difference, because the cost was in AppKit and the attribute graph rather than in
the app's own code) but because it is what ships, and because debug numbers reported as if they
were the shipping experience mislead everyone downstream.

**Nothing updates it automatically.** The owner decides when. An app that changes under someone
mid test is the problem this exists to remove, so do not helpfully rebuild it after landing work.
Say what landed and let them run the script.

**Offer a pin.** `./Tools/master.sh <ref>` lets them go back to a known good commit when a later one
turns out worse, without touching the repository.

## Reporting to the owner

When work lands, tell them what changed and that the update is one command away. Do not claim
the app is updated when it is not.

Flag anything in the new build that no agent could see. If agents were barred from launching,
some of what they wrote has never been on a screen, and the owner should know which parts those
are before they go looking. Hover behaviour, animation, positioning near a pane edge and anything
that depends on a pointer are the usual candidates.
