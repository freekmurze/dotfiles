---
name: conductor
description: Working with Conductor (conductor.build), the parallel-agent app. Creating local and cloud workspaces, the conductor CLI and public API, auth, sessions and messages, deep links, settings and setup scripts, environment variables, and writing briefs for handed-off agents. Use when asked to open, create, rename, archive, or inspect a workspace or session, spin up parallel agents, kick off or hand off work in Conductor, query conductor sql, debug a deep link that did not work, fix Conductor auth, or explain the local-versus-cloud distinction.
---

# Conductor

## The one fact that decides everything

**Local workspaces are created with a `conductor://` deep link. Cloud workspaces are created with
the CLI or the API.** The CLI cannot create a local workspace; do not go looking for a subcommand.

| They want | Do this |
|---|---|
| A local workspace (the normal case, Mac) | `open "conductor://prompt=<urlencoded>&path=<urlencoded repo root>"` |
| A cloud workspace, explicitly (runs on Linux) | `conductor workspaces create`, after checking the caveats below |

## Creating a local workspace

Write the brief to a file, then build the URL — do not try to inline a multi-line brief into a
shell string.

```bash
python3 -c "
import urllib.parse
p = open('/tmp/brief.txt').read()
print('conductor://prompt=' + urllib.parse.quote(p, safe='')
      + '&path=' + urllib.parse.quote('/path/to/repo/root', safe=''))
" > /tmp/ws.url

open "$(cat /tmp/ws.url)"
```

`safe=''` matters: the default `quote` leaves `/` unescaped, which is fine for `path` but wrong
inside `prompt`. (`printf %s "$BRIEF" | jq -sRr @uri` produces equivalent output if you prefer jq,
though the Python form above is the one verified end to end.)

Percent-encode both values. An unencoded `&`, `#`, `?` or `+` in the brief silently truncates the
URL — this is the single most common cause of "the deep link doesn't work". Briefs of ~4-6 KB
encoded work fine; that is roughly a full page of markdown.

`path` is the repo root, exactly, no trailing slash. It must be the root of the git repo Conductor
has registered, not a subdirectory and not an existing worktree.

What happens: the app creates a git worktree immediately, derives a branch name **from the prompt
text**, and runs `scripts.setup`. Verified 2026-08-16: 13 workspaces created this way in one
session, all present in `git worktree list` with branch names matching their brief subjects.

**The prompt is only PRE-FILLED in the composer. A human must press enter.** The session sits
`idle` with zero `session_messages` rows until they do (confirmed by DB inspection after setup
completed). There is no way to make a local deep link auto-run. Say this when you hand the
workspace over, or the user will assume the agent is already working.

Other flat forms: `conductor://linear_id=<id>&prompt=...`, and `conductor://async?repo=<name>&plan=<base64 md>`
which stages `[ASYNC]-plan.md` into `.context/attachments/` (also does not auto-run). Flat variants
take `key=value&key=value` straight after `conductor://` with no hostname; only `async` uses a real
hostname. Docs: https://www.conductor.build/docs/reference/deep-links

### `open` exit status proves nothing

`open "conductor://..."` exits 0 whenever macOS finds a handler registered for the scheme. A
malformed URL, an unknown key, or a nonexistent path all still exit 0.

**Never report success from the `open` exit code.** Verify with `git worktree list` in the repo
root, or look at the app. If you tell a user you created something on the strength of a 0 exit,
you may be wrong and have no way to know.

### When the deep link appears not to work

In order: is the app running; is `path` the exact repo root with no trailing slash; is the URL
actually percent-encoded; does a worktree appear in `git worktree list`. If the worktree exists and
the agent simply is not doing anything, that is the pre-fill behaviour above, not a failure.

## Locating the CLI (there is a name collision)

`/usr/local/bin/conductor` is often **not** Conductor.build. It is commonly the Orkes / Netflix
Conductor workflow-orchestration CLI, which has completely unrelated subcommands
(`workflow`, `task`, `worker`, `agent`). Check before use.

The real one:

```
/Applications/Conductor.app/Contents/Resources/bin/conductor
```

That is a `sh` wrapper around `.internal/conductor-runtime cli`. Verify with `conductor --help`:
the Conductor.build CLI lists `auth`, `projects`, `workspaces`, `sessions`, `messages`, `models`, `sql`.

The app binary and the runtime contain no literal `conductor://` strings, so do not bother
reverse-engineering deep links out of them. That was a dead end.

## Know where you are running

```
CONDUCTOR_IS_LOCAL      1 = user's Mac, 0 = cloud Linux sandbox
CONDUCTOR_WORKSPACE_NAME
CONDUCTOR_WORKSPACE_ID  the CLI uses this when a workspace arg is omitted
CONDUCTOR_WORKSPACE_PATH
CONDUCTOR_ROOT_PATH     repo root; equals workspace path in cloud
CONDUCTOR_DEFAULT_BRANCH  local only
CONDUCTOR_PORT          local only, first of ten allocated ports
```

## Auth

Cloud commands need a token in the macOS Keychain.

```
conductor auth status     # does an entry exist (does not print it)
conductor auth whoami     # verify against the API
conductor auth login --token <token>
```

`! No keychain entry for https://api.conductor.build.` means unauthenticated. Do not ask the user
to paste a token to you. Tell them to run `conductor auth login` themselves.

API directly: base `https://api.conductor.build/v0`, `Authorization: Bearer <api key>`,
spec at `https://api.conductor.build/v0/openapi.json`.

## Creating a cloud workspace

```
conductor workspaces create \
  [--project-id <id>] [--repo-url <url>] [--branch <name>] \
  [--name <name>] [--session-name <name>] \
  [--agent claude|codex|cursor|acp] [--model <model>] [--effort <level>] \
  [--env KEY=value]... [--channel prod|beta]
```

Creates the workspace **and its first session**. To add more sessions later:

```
conductor sessions create --workspace <id> --agent claude --message "<prompt>" [--fast-mode]
conductor messages create --session <id> --message "<text>"
```

Check `conductor models` for each agent's valid model ids and effort levels before assuming a model
is unavailable.

### Cloud caveats to check first

1. **Cloud is Linux.** If the repo's `scripts.setup` is Mac-specific it will partially fail. Read
   `<repo>/.conductor/settings.toml` (or a legacy `conductor.json`) before recommending cloud. Real
   example: a setup script that does `cp "$CONDUCTOR_ROOT_PATH/.env" .env` and links a Herd/Valet
   site is local-only, and `CONDUCTOR_ROOT_PATH` is the workspace itself in cloud, so that copy is
   a no-op at best.
2. **Cloud setup reads settings from the branch the workspace is created from**, unlike the Mac app
   which reads shared settings from the default branch on the remote.
3. `CONDUCTOR_PORT` is unset in cloud, so run scripts that need it must be
   `available_in = [ "local" ]`.

## Deep links returned by the API (cloud)

Workspace and session objects carry a **`deepLink`** field (a plain `string` in the OpenAPI spec:
no `pattern`, no `example`, no documented format). It appears on the workspace and session schemas
returned by create/get/list.

**For cloud objects, read `deepLink` from the response rather than constructing one** — it encodes
ids you would otherwise have to assemble, and the app-opening grammar for an existing remote
workspace is not the same as the documented flat local form.

This does **not** mean deep links can never be written by hand. The flat
`conductor://prompt=...&path=...` form for creating a *local* workspace is documented and verified
working — see "Creating a local workspace" above. The two are different jobs: one opens an object
that already exists in the cloud, the other creates a worktree on this Mac.

`--channel <prod|beta>` selects which desktop-app build the returned deep link should open. It is a
query param on every link-returning endpoint:

```
GET  /v0/projects/{projectId}/workspaces      POST /v0/workspaces
GET  /v0/workspaces/{workspaceId}             POST /v0/workspaces/{workspaceId}/rename
GET  /v0/workspaces/{workspaceId}/sessions    POST /v0/sessions
GET  /v0/sessions/{sessionId}                 POST /v0/sessions/{sessionId}/rename
POST /v0/sessions/{sessionId}/messages
```

The CLI's own guidance: *"Use deep links instead of IDs to direct users to work you've started;
deep links are clickable, while IDs are not."* So when you do create cloud work, surface the
`deepLink`, not the id.

## Other CLI commands worth knowing

```
conductor workspaces get|rename|archive|sessions|status [workspaceId]
conductor sessions get|rename|archive|messages|status|cancel <sessionId>
conductor sessions messages <id> --after <messageId>   # incremental polling
conductor projects list|get|workspaces
conductor sql "<read-only SELECT>"                     # over session_transcripts_view
```

- Add `--json` for machine-readable output (auth commands are human-readable only).
- List commands paginate with a small default page. Pass `--limit` (max 100) and step `--offset`
  while the JSON response reports `hasMore`. Do not assume page one is everything.
- To find a specific workspace or session, prefer `conductor sql` over paging list commands.
  `conductor sql --help` documents the view's columns. Workspaces with no transcripts yet do not
  appear in that view, so on zero matches fall back to `workspaces`/`projects` list commands before
  concluding something does not exist.
- Exit codes: `0` success, `1` runtime, `2` usage, `3` auth, `4` server.
- Prefer internal subagents for read-only work. Spend a Conductor workspace or session only when the
  user asks, or when a specific agent or model is only reachable through Conductor.

## Writing the brief

The brief is the deliverable, whether you pass it via `prompt=` in a deep link or paste it for the
user. Since the prompt only pre-fills, a bad brief costs the user a full agent run to discover.

1. **Write one self-contained brief per workspace.** Independent fixes get separate workspaces;
   that is the whole point of Conductor. Do not batch unrelated fixes into one prompt.
2. **Include the trap you already found.** A brief is worth writing mostly because you can encode
   the thing that would break a naive fix. Verify the fix is safe *before* writing the brief, not
   after. Example: guarding a route on `environment('local')` looks obviously right and silently
   breaks every test that hits the route, because `APP_ENV=testing`; the brief must say
   `['local', 'testing']` and name the affected test files.
3. **State the scope guard.** What not to touch, and which tempting adjacent refactor belongs in a
   different PR.
4. **Mind where the brief lives.** `.context/` is per-workspace and gitignored, so a brief written
   there does **not** travel to a new workspace. Either paste the brief as the first message, or
   write it under `$CONDUCTOR_ROOT_PATH` (the repo root, stable and readable from any workspace).
   Conductor's "files to copy" defaults to `.env*` only.
5. Suggest workspace names, and give the base branch.

## Settings and scripts

Precedence: managed > repo local > repo shared > user > defaults. TOML outranks legacy JSON at every
layer.

```
<repo>/.conductor/settings.local.toml    machine-local, one repo
<repo>/.conductor/settings.toml          shared with the team
~/.conductor/settings.toml               this user, all repos
~/.conductor/settings.managed.toml       org-controlled
```

The Mac client ignores a repo-level `conductor.json` once `.conductor/settings.toml` exists, but
cloud setup still falls back to `conductor.json` for a setup script. If a repo has both, that is a
migration bug worth flagging. Migration map: `scripts.run` becomes `scripts.run.<id>.command`,
`runScriptMode` becomes `scripts.run_mode`.

Scripts run in non-interactive shells (`zsh` on Mac, `bash` in cloud), so put toolchain setup in the
script rather than relying on shell startup files.

## Docs

- https://conductor.build/docs/concepts/workspaces-and-branches
- https://conductor.build/docs/concepts/parallel-agents
- https://conductor.build/docs/api
- https://conductor.build/docs/reference/settings
- https://conductor.build/docs/reference/scripts
- https://conductor.build/docs/reference/files-to-copy

There is no `/docs/reference/cli` page (404 as of CLI 0.80.1). The CLI's own
`conductor --help` command reference is more complete than the site.
