---
name: mailcoach
description: Manage email marketing with the Mailcoach CLI — email lists, subscribers, campaigns, transactional emails, templates, automations, tags, and segments.
---

# Mailcoach CLI

Use this skill when the user wants to manage email marketing through their Mailcoach instance. This includes managing email lists, subscribers, campaigns, transactional emails, templates, automations, tags, segments, and suppressions.

## Prerequisites

The `mailcoach` CLI must be installed and authenticated. If a command returns a 401 error, the user needs to authenticate first.

```bash
# Authenticate with a Mailcoach instance
mailcoach login

# Verify authentication works
mailcoach list-email-lists
```

See [references/authentication.md](references/authentication.md) for the full login flow.

## How commands work

Commands are **auto-generated from the Mailcoach OpenAPI spec**. Always discover commands dynamically rather than guessing names.

```bash
# List all available commands
mailcoach list

# Get help for a specific command (shows all options)
mailcoach <command-name> --help
```

### Command naming

Commands use kebab-case derived from API operation IDs:

- `list-email-lists`, `create-email-list`, `show-email-list`, `update-email-list`, `delete-email-list`
- `list-subscribers`, `create-subscriber`, `show-subscriber`, `update-subscriber`
- `list-campaigns`, `create-campaign`, `send-campaign`, `send-campaign-test`
- `list-transactional-mails`, `send-transactional-mail`
- `list-tags`, `create-tag`, `list-segments`, `create-segment`
- `trigger-automation`

### Parameter naming

- Path parameters like `{emailList}` become **required** options: `--email-list=<uuid>`
- Query parameters like `filter[search]` become **optional** options: `--filter-search=<value>`
- Both `snake_case` and `camelCase` parameter names are converted to `--kebab-case` (e.g., `book_id` and `bookId` both become `--book-id`)

### Sending data

```bash
# JSON body (preferred for creates/updates)
mailcoach create-email-list --input '{"name": "Newsletter", "default_from_email": "hi@example.com"}'

# Form fields (repeatable, simpler for flat data)
mailcoach create-email-list --field name="Newsletter" --field default_from_email="hi@example.com"

# File uploads (prefix path with @)
mailcoach create-subscriber-import --email-list=<uuid> --field csv=@/path/to/subscribers.csv
```

- `--field` values are sent as JSON by default. When any field contains a file (`@` prefix), the entire request switches to `multipart/form-data`.
- You cannot combine `--field` and `--input` in the same command.

### Output formats

```bash
mailcoach list-email-lists              # Human-readable table (default)
mailcoach list-email-lists --json       # Raw JSON (useful for parsing)
mailcoach list-email-lists --yaml       # YAML output
mailcoach list-email-lists --minify     # Minified single-line JSON (implies --json)
mailcoach list-email-lists --H          # Include response headers
mailcoach list-email-lists --output-html # Show HTML response bodies (hidden by default)
```

Always use `--json` when you need to extract UUIDs or data from responses for use in subsequent commands.

### Filtering, sorting, and pagination

```bash
mailcoach list-subscribers --email-list=<uuid> --filter-search="john" --sort=email --filter-per-page=50 --filter-page=2
```

- `--filter-search`: fuzzy text search
- `--sort`: field name, prefix with `-` for descending (e.g., `--sort=-created_at`)
- `--filter-per-page`: results per page (max 100, default 15)
- `--filter-page`: page number

### Error handling

- **401**: Authentication failed. Run `mailcoach login` to re-authenticate.
- **422**: Validation error. The response body contains details about which fields are invalid.
- **404**: Resource not found. Verify the UUID is correct.
- **Missing path parameter**: The CLI tells you which `--option` is required.
- All errors exit with a non-zero code, so you can chain commands with `&&`.
- HTML error responses hide the body by default — add `--output-html` to see them.

## Common workflows

### Command patterns
- [references/command-patterns.md](references/command-patterns.md) — Discovering commands and the shared flag conventions

### Email lists and subscribers
- [references/email-lists.md](references/email-lists.md) — Create and manage email lists, tags, and segments
- [references/subscribers.md](references/subscribers.md) — Add, import, tag, and manage subscribers

### Campaigns
- [references/campaigns.md](references/campaigns.md) — Full campaign lifecycle: create, configure, test, send, and view statistics

### Transactional emails
- [references/transactional.md](references/transactional.md) — Send transactional emails via templates or inline content

### Templates and automations
- [references/templates.md](references/templates.md) — Manage reusable email templates
- [references/automations.md](references/automations.md) — Trigger automations via the CLI

## Important notes

- Mailcoach is self-hosted: each user has a unique instance URL configured during `mailcoach login`
- The OpenAPI spec is cached for 24 hours. Run `mailcoach clear-cache` to force a refresh if commands seem outdated.
- UUIDs are used as resource identifiers throughout the API.
- When chaining commands (e.g., create a list then add subscribers), use `--json` to extract UUIDs from responses.
