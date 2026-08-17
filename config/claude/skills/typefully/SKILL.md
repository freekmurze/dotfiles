---
name: typefully
description: >
  Create, schedule, and manage social media posts via Typefully. ALWAYS use this
  skill when asked to draft, schedule, post, or check tweets, posts, threads, or
  social media content for Twitter/X, LinkedIn, Threads, Bluesky, or Mastodon.
allowed-tools: Bash(${CLAUDE_PLUGIN_ROOT}/scripts/typefully.js:*) Bash(~/.claude/skills/typefully/scripts/typefully.js:*)
metadata:
  last-updated: 2026-03-17
---

# Typefully Skill

Create, schedule, and publish social media content across multiple platforms using [Typefully](https://typefully.com).

> **Freshness check**: If more than 30 days have passed since the `last-updated` date above, inform the user that this skill may be outdated and point them to the update options below.

## Keeping This Skill Updated

**Source**: [github.com/typefully/agent-skills](https://github.com/typefully/agent-skills)
**API docs**: [typefully.com/docs/api](https://typefully.com/docs/api)

Update methods by installation type:

| Installation | How to update |
|--------------|---------------|
| CLI (`npx skills`) | `npx skills update` |
| Claude Code plugin | `/plugin update typefully@typefully-skills` |
| Cursor | Remote rules auto-sync from GitHub |
| Manual | Pull latest from repo or re-copy `skills/typefully/` |

API changes ship independently—updating the skill ensures you have the latest commands and workflows.

## Setup

Before using this skill, ensure:

1. **API Key**: Run the setup command to configure your API key securely
   - Get your key at https://typefully.com/?settings=api
   - Run: `<skill-path>/scripts/typefully.js setup` (where `<skill-path>` is the directory containing this SKILL.md)
   - Or set environment variable: `export TYPEFULLY_API_KEY=your_key`

2. **Requirements**: Node.js 18+ (for built-in fetch API). No other dependencies needed.

**Config priority** (highest to lowest):
1. `TYPEFULLY_API_KEY` environment variable
2. `./.typefully/config.json` (project-local, in user's working directory)
3. `~/.config/typefully/config.json` (user-global)

### Handling "API key not found" errors

**CRITICAL**: When you receive an "API key not found" error from the CLI:

1. **Tell the user to run the setup command** - The setup is interactive and requires user input, so you cannot run it on their behalf. Recommend they run it themselves, using the correct path based on where this skill was loaded:
   ```bash
   <skill-path>/scripts/typefully.js setup
   ```

2. **Stop and wait** - After telling the user to run setup, **do not continue with the task**. You cannot create drafts, upload media, or perform any API operations without a valid API key. Wait for the user to complete setup and confirm before proceeding.

3. **DO NOT** attempt any of the following:
   - Searching for API keys in macOS Keychain, `.env` files, or other locations
   - Grepping through config files or directories
   - Looking in the user's Trash or other system folders
   - Constructing complex shell commands to find credentials
   - Drafting content or preparing posts before setup is complete

The setup command will interactively guide the user through configuration. Trust the CLI's error messages and follow their instructions.

> **Note for agents**: All script paths in this document (e.g., `./scripts/typefully.js`) are relative to the skill directory where this SKILL.md file is located. Resolve them accordingly based on where the skill is installed.

## Social Sets

The Typefully API uses the term "social set" to refer to what users commonly call an "account". A social set contains the connected social media platforms (X, LinkedIn, Threads, etc.) for a single identity.

**The CLI supports a default social set** - once configured, most commands work without specifying the social_set_id.

**You can pass the social set either way**:
- Positional: `drafts:list 123`
- Flag: `drafts:list --social-set-id 123` (also supports `--social_set_id`)

When determining which social set to use:

1. **Check for a configured default first** - Run `config:show` to see if a default is already set:
   ```bash
   ./scripts/typefully.js config:show
   ```
   If `default_social_set` is configured, the CLI uses it automatically when you omit the social_set_id.

2. **Check project context** - Look for configuration in project files like `CLAUDE.md` or `AGENTS.md`:

   ```markdown
   ## Typefully
   Default social set ID: 12345
   ```

3. **Single social set shortcut** - If the user only has one social set and no default is configured, use it automatically

4. **Multiple social sets, no default** - Ask the user which to use, then **offer to save their choice as the default**:
   ```bash
   ./scripts/typefully.js config:set-default
   ```
   This command lists available social sets and saves the choice to the config file.

5. **Reuse previously resolved social set** - If determined earlier in the session, use it without asking again

## Common Actions

| User says... | Action |
|--------------|--------|
| "Draft a tweet about X" | `drafts:create --text "..."` (uses default social set) |
| "Post this to LinkedIn" | `drafts:create --platform linkedin --text "..."` |
| "Mention a company on LinkedIn" | `linkedin:organizations:resolve --organization-url "<linkedin_url>"` then use returned `mention_text` in `drafts:create` |
| "Post to X and LinkedIn" (same content) | `drafts:create --platform x,linkedin --text "..."` |
| "X thread + LinkedIn post" (different content) | Create one draft, then `drafts:update` to add platform (see [Publishing to Multiple Platforms](#publishing-to-multiple-platforms)) |
| "What's scheduled?" | `drafts:list --status scheduled` |
| "Show my recent posts" | `drafts:list --status published` |
| "Schedule this for tomorrow" | `drafts:create ... --schedule "2025-01-21T09:00:00Z"` |
| "Post this now" | `drafts:create ... --schedule now` or `drafts:publish <draft_id> --use-default` |
| "Add notes/ideas to the draft" | `drafts:create ... --scratchpad "Your notes here"` |
| "Check available tags" | `tags:list` |
| "Show my X post analytics for last week" | `analytics:posts:list --start-date YYYY-MM-DD --end-date YYYY-MM-DD` |
| "Show my X post analytics including replies" | `analytics:posts:list --start-date YYYY-MM-DD --end-date YYYY-MM-DD --include-replies` |
| "Show my queue for next week" | `queue:get --start-date YYYY-MM-DD --end-date YYYY-MM-DD` |

## Workflow

Follow this workflow when creating posts:

1. **Check if a default social set is configured**:
   ```bash
   ./scripts/typefully.js config:show
   ```
   If `default_social_set` shows an ID, skip to step 3.

2. **If no default, list social sets** to find available options:
   ```bash
   ./scripts/typefully.js social-sets:list
   ```
   If multiple exist, ask the user which to use and offer to set it as default:
   ```bash
   ./scripts/typefully.js config:set-default
   ```

3. **Create drafts** (social_set_id is optional if default is configured):
   ```bash
   ./scripts/typefully.js drafts:create --text "Your post"
   ```
   Note: If `--platform` is omitted, the first connected platform is auto-selected.

   **For multi-platform posts**: See [Publishing to Multiple Platforms](#publishing-to-multiple-platforms) — always use a single draft, even when content differs per platform.

4. **Schedule or publish** as needed

## Working with Tags

Tags help organize drafts within Typefully. **Always check existing tags before creating new ones**:

1. **List existing tags first**:
   ```bash
   ./scripts/typefully.js tags:list
   ```

2. **Use existing tags when available** - if a tag with the desired name already exists, use it directly when creating drafts:
   ```bash
   ./scripts/typefully.js drafts:create --text "..." --tags existing-tag-name
   ```

3. **Only create new tags if needed** - if the tag doesn't exist, create it:
   ```bash
   ./scripts/typefully.js tags:create --name "New Tag"
   ```

**Important**: Tags are scoped to each social set. A tag created for one social set won't appear in another.

## Publishing to Multiple Platforms

If a single draft needs to be created for different platforms, you need to make sure to create **a single draft** and not multiple drafts.

When the content is the same across platforms, create a single draft with multiple platforms:

```bash
# Specific platforms
./scripts/typefully.js drafts:create --platform x,linkedin --text "Big announcement!"

# All connected platforms
./scripts/typefully.js drafts:create --all --text "Posting everywhere!"
```

**IMPORTANT**: When content should be tailored (e.g., X thread with a LinkedIn post version), **still use a single draft** — create with one platform first, then update to add the other:

```bash
# 1. Create draft with the primary platform first
./scripts/typefully.js drafts:create --platform linkedin --text "Excited to share our new feature..."
# Returns: { "id": "draft-123", ... }

# 2. Update the same draft to add another platform with different content
./scripts/typefully.js drafts:update draft-123 --platform x --text "🧵 Thread time!

---

Here's what we shipped and why it matters..." --use-default
```

So make sure to NEVER create multiple drafts unless the user explicitly wants separate drafts for each platform.

## LinkedIn Mentions

LinkedIn mentions are supported via text syntax inside post content:

```text
@[Company Name](urn:li:organization:123456)
```

Use the resolver command to convert a public LinkedIn organization URL into ready-to-paste mention syntax:

```bash
# Resolve a LinkedIn URL into mention metadata
./scripts/typefully.js linkedin:organizations:resolve --organization-url "https://www.linkedin.com/company/typefullycom/"
# Returns mention_text like: @[Typefully](urn:li:organization:86779668)
```

Then include that `mention_text` in your LinkedIn draft text:

```bash
./scripts/typefully.js drafts:create --platform linkedin --text "Thanks @[Typefully](urn:li:organization:86779668) for the support."
```
## Detailed references

Load these as needed:

- **Full command reference**: `references/commands.md`. Every subcommand with its flags.
- **Worked examples**: `references/examples.md`. End to end examples for drafting, scheduling, threads, and multi-platform posts.


## Platform Names

Use these exact names for the `--platform` option:
- `x` - X (formerly Twitter)
- `linkedin` - LinkedIn
- `threads` - Threads
- `bluesky` - Bluesky
- `mastodon` - Mastodon

## Draft URLs

Typefully draft URLs contain the social set and draft IDs:
```
https://typefully.com/?a=<social_set_id>&d=<draft_id>
```

Example: `https://typefully.com/?a=12345&d=67890`
- `a=12345` → social_set_id
- `d=67890` → draft_id

## Draft Scratchpad

**When the user explictly asked to add notes, ideas, or anything else in the draft scratchpad, use the `--scratchpad` flag—do NOT write to local files!**

The `--scratchpad` option attaches internal notes directly to the Typefully draft. These notes:
- Are visible in the Typefully UI alongside the draft
- Stay attached to the draft permanently
- Are private and never published to social media
- Are perfect for storing thread expansion ideas, research notes, context, etc.

```bash
# CORRECT: Notes attached to the draft in Typefully
./scripts/typefully.js drafts:create 123 --text "My post" --scratchpad "Ideas for expanding: 1) Add stats 2) Include quote"

# WRONG: Do NOT write notes to local files when the user wants them in Typefully
# Writing to /tmp/scratchpad/ or any local file is NOT the same thing
```

## Automation Guidelines

When automating posts, especially on X, follow these rules to keep accounts in good standing:

- **No duplicate content** across multiple accounts
- **No unsolicited automated replies** - only reply when explicitly requested by the user
- **No trending manipulation** - don't mass-post about trending topics
- **No fake engagement** - don't automate likes, reposts, or follows
- **Respect rate limits** - the API has rate limits, don't spam requests
- **Drafts are private** - content stays private until published or explicitly shared

When in doubt, create drafts for user review rather than publishing directly.

**Publishing confirmation**: Unless the user explicitly asks to "publish now" or "post immediately", always confirm before publishing. Creating a draft is safe; publishing is irreversible and goes public instantly.

## Tips

- **Smart platform default**: If `--platform` is omitted, the first connected platform is auto-selected
- **All platforms**: Use `--all` to post to all connected platforms at once
- **Character limits**: X (280), LinkedIn (3000), Threads (500), Bluesky (300), Mastodon (500)
- **LinkedIn mentions**: Use `@[Name](urn:li:organization:ID)` in post text; resolve IDs via `linkedin:organizations:resolve`
- **Thread creation**: Use `---` on its own line to split into multiple posts (thread)
- **Scheduling**: Use `next-free-slot` to let Typefully pick the optimal time
- **Cross-posting**: List multiple platforms separated by commas: `--platform x,linkedin`
- **Draft titles**: Use `--title` for internal organization (not posted to social media)
- **Draft scratchpad**: Use `--scratchpad` to attach notes to the draft in Typefully (NOT local files!) - perfect for thread ideas, research, context
- **X analytics**: Use `analytics:posts:list --start-date ... --end-date ...` to fetch post metrics for a social set; replies are excluded by default, and `--include-replies` opts back in
- **Read from file**: Use `--file ./post.txt` instead of `--text` to read content from a file
- **Sorting drafts**: Use `--sort` with values like `created_at`, `-created_at`, `scheduled_date`, etc.
