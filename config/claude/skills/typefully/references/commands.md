
## Commands Reference

### User & Social Sets

| Command | Description |
|---------|-------------|
| `me:get` | Get authenticated user info |
| `social-sets:list` | List all social sets you can access |
| `social-sets:get <id>` | Get social set details including connected platforms |
| `linkedin:organizations:resolve [social_set_id] --organization-url <url>` | Resolve LinkedIn company/school URL into mention metadata (`mention_text`, `urn`) |

### Analytics

All analytics commands support an optional `[social_set_id]` - if omitted, the configured default is used.

The public API currently supports **X post analytics only** on this endpoint. The CLI defaults `--platform` to `x`, so you can usually omit it.

Replies are now **excluded by default** so the result set matches the main published-post view more closely. Add `--include-replies` when you explicitly want reply posts included.

Analytics responses return post-level metrics for the requested inclusive date range, including:
- `impressions`
- engagement totals and breakdowns like `likes`, `comments`, `shares`, `quotes`, `saves`, `profile_clicks`, and `link_clicks`

| Command | Description |
|---------|-------------|
| `analytics:posts:list [social_set_id] --start-date <YYYY-MM-DD> --end-date <YYYY-MM-DD>` | List X posts with normalized analytics metrics for an inclusive date range |
| `analytics:posts:list ... --start_date <YYYY-MM-DD> --end_date <YYYY-MM-DD>` | Snake case aliases for date flags (copied from API docs) |
| `analytics:posts:list ... --include-replies` | Include X replies in the results (excluded by default) |
| `analytics:posts:list ... --include_replies` | Snake case alias for the include-replies flag |
| `analytics:posts:list ... --limit 100 --offset 25` | Paginate through results |
| `analytics:posts:list ... --platform x` | Explicitly request X analytics (currently the only supported platform) |

### Drafts

All drafts commands support an optional `[social_set_id]` - if omitted, the configured default is used.
**Safety note**: For commands that take `[social_set_id] <draft_id>`, if you pass only a single argument (the draft_id) while a default social set is configured, you must add `--use-default` to confirm intent.

| Command | Description |
|---------|-------------|
| `drafts:list [social_set_id]` | List drafts (add `--status scheduled` to filter, `--sort` to order) |
| `drafts:get [social_set_id] <draft_id>` | Get a specific draft with full content (single-arg requires `--use-default` if a default is configured) |
| `drafts:create [social_set_id] --text "..."` | Create a new draft (auto-selects platform) |
| `drafts:create [social_set_id] --platform x --text "..."` | Create a draft for specific platform(s) |
| `drafts:create [social_set_id] --all --text "..."` | Create a draft for all connected platforms |
| `drafts:create [social_set_id] --file <path>` | Create draft from file content |
| `drafts:create ... --media <media_ids>` | Attach media. `,` separates ids on the same post; `\|` delimits per-post groups in a thread (mirrors how `---` delimits text). |
| `drafts:create ... --reply-to <url>` | Reply to an existing X post |
| `drafts:create ... --community <id>` | Post to an X community |
| `drafts:create ... --quote-post-url <url>` | Quote an existing X post URL |
| `drafts:create ... --share` | Generate a public share URL for the draft |
| `drafts:create ... --scratchpad "..."` | Add internal notes/scratchpad to the draft |
| `drafts:update [social_set_id] <draft_id> --text "..."` | Update an existing draft (single-arg requires `--use-default` if a default is configured) |
| `drafts:update ... --quote-post-url <url>` | Update X post(s) in a draft to quote an existing post URL |
| `drafts:update [social_set_id] <draft_id> --tags "tag1,tag2"` | Update tags on an existing draft (content unchanged) |
| `drafts:update ... --share` | Generate a public share URL for the draft |
| `drafts:update ... --scratchpad "..."` | Update internal notes/scratchpad |
| `drafts:update [social_set_id] <draft_id> --append --text "..."` | Append to existing thread |

### Scheduling & Publishing

**Safety note**: These commands require `--use-default` when using the default social set with a single argument (to prevent accidental operations from ambiguous syntax).

| Command | Description |
|---------|-------------|
| `drafts:delete <social_set_id> <draft_id>` | Delete a draft (explicit IDs) |
| `drafts:delete <draft_id> --use-default` | Delete using default social set |
| `drafts:schedule <social_set_id> <draft_id> --time next-free-slot` | Schedule to next available slot |
| `drafts:schedule <draft_id> --time next-free-slot --use-default` | Schedule using default social set |
| `drafts:publish <social_set_id> <draft_id>` | Publish immediately |
| `drafts:publish <draft_id> --use-default` | Publish using default social set |

### Queue

All queue commands support an optional `[social_set_id]` - if omitted, the configured default is used.

The queue is a **social-set-specific timeline** made of:
- Queue slots generated from that social set's queue schedule
- Scheduled drafts/posts that belong to that same social set

Use `queue:get` when the user asks what is already scheduled (or free) for a given account in a date range.

| Command | Description |
|---------|-------------|
| `queue:get [social_set_id] --start-date <YYYY-MM-DD> --end-date <YYYY-MM-DD>` | Get the queue timeline for one social set: free queue slots plus scheduled drafts/posts in a date range |
| `queue:get ... --start_date <YYYY-MM-DD> --end_date <YYYY-MM-DD>` | Snake case aliases for date flags (copied from API docs) |
| `queue:schedule:get [social_set_id]` | Get queue schedule rules |
| `queue:schedule:put [social_set_id] --rules '[{"h":9,"m":30,"days":["mon","wed","fri"]}]'` | Replace queue schedule rules (full replacement) |

### Tags

| Command | Description |
|---------|-------------|
| `tags:list [social_set_id]` | List all tags |
| `tags:create [social_set_id] --name "Tag Name"` | Create a new tag |

### Media

| Command | Description |
|---------|-------------|
| `media:upload [social_set_id] <file_path>` | Upload media, wait for processing, return ready media_id |
| `media:upload ... --no-wait` | Upload and return immediately (use media:status to poll) |
| `media:upload ... --timeout <seconds>` | Set custom timeout (default: 60) |
| `media:status [social_set_id] <media_id>` | Check media upload status |

### Setup & Configuration

| Command | Description |
|---------|-------------|
| `setup` | Interactive setup - prompts for API key, storage location, and default social set |
| `setup --key <key> --location <global\|local>` | Non-interactive setup for scripts/CI (auto-selects default if only one social set) |
| `setup --key <key> --default-social-set <id>` | Non-interactive setup with explicit default social set |
| `setup --key <key> --no-default` | Non-interactive setup, skip default social set selection |
| `config:show` | Show current config, API key source, and default social set |
| `config:set-default [social_set_id]` | Set default social set (interactive if ID omitted) |
