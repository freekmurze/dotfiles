# Command Patterns

## Discovering commands

```bash
# List all available commands
mailcoach list

# Get detailed help for a command
mailcoach create-campaign --help
```

Always check `--help` before running a command to see required and optional parameters.

## Sending data

### JSON input (preferred for creates and updates)

```bash
mailcoach create-campaign --input '{"name": "March Newsletter", "email_list_uuid": "abc-123"}'
```

Use `--input` when the request body has nested data, arrays, or when you want to match the API schema exactly.

### Form fields (for simple key-value data)

```bash
mailcoach create-email-list --field name="My List" --field default_from_email="hi@example.com"
```

Fields are sent as JSON by default. If the API spec declares `application/x-www-form-urlencoded` as the content type, fields are sent as form data instead. You can repeat `--field` multiple times.

### File uploads

```bash
mailcoach create-subscriber-import --email-list=<uuid> --field csv=@/path/to/file.csv
```

Prefix the file path with `@` to upload it. When any field contains a file, the entire request is sent as `multipart/form-data`.

**You cannot combine `--field` and `--input` in the same command.**

## Output formats

| Flag | Output |
|------|--------|
| (none) | Human-readable formatted output |
| `--json` | Raw JSON |
| `--yaml` | YAML |
| `--minify` | Minified single-line JSON (implies `--json`) |
| `--H` | Include HTTP response headers |
| `--output-html` | Show HTML response bodies (hidden by default) |

Use `--json` when you need to parse the response or extract UUIDs for subsequent commands.

## Filtering and pagination

Most list commands support these query parameters as options:

```bash
mailcoach list-subscribers --email-list=<uuid> \
  --filter-search="john@example.com" \
  --filter-status=subscribed \
  --sort=-created_at \
  --filter-per-page=100 \
  --filter-page=1
```

- `--filter-search`: fuzzy text search across relevant fields
- `--sort`: sort field, prefix `-` for descending
- `--filter-per-page`: items per page (max 100, default 15)
- `--filter-page`: page number (default 1)
- Additional filters vary by endpoint (use `--help` to see them)

## Debugging

Use `-vvv` (very verbose) to see the full HTTP request before it's sent:

```bash
mailcoach create-campaign --input '{"name": "Test"}' -vvv
```

This shows: HTTP method, resolved URL, request headers (Accept, Content-Type, Authorization), and request body.

## Gotchas

- **Path parameters are required**, query parameters are optional. Missing a required path param gives a clear error telling you which `--option` is needed.
- **All errors exit with non-zero codes**, so you can safely chain commands with `&&`.
- **Both `snake_case` and `camelCase`** parameter names become `--kebab-case` options (e.g., `book_id` and `bookId` both become `--book-id`).
- **Bracket notation** in query params is converted to kebab-case: `filter[id]` becomes `--filter-id`.
- **HTML error responses** hide the body by default. Use `--output-html` to see them.
