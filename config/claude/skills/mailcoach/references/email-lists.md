# Email Lists

## List all email lists

```bash
mailcoach list-email-lists
mailcoach list-email-lists --filter-search="newsletter" --json
```

## Create an email list

```bash
mailcoach create-email-list --input '{
  "name": "Newsletter",
  "default_from_email": "newsletter@example.com",
  "default_from_name": "My Company"
}'
```

## Show, update, delete

```bash
mailcoach show-email-list --email-list=<uuid>
mailcoach update-email-list --email-list=<uuid> --input '{"name": "Updated Name"}'
mailcoach delete-email-list --email-list=<uuid>
```

## Tags

Tags belong to an email list and can be assigned to subscribers.

```bash
# List tags for an email list
mailcoach list-tags --email-list=<uuid>

# Create a tag (returns existing tag if name already exists)
mailcoach create-tag --email-list=<uuid> --input '{"name": "vip"}'

# Show, update, delete
mailcoach show-tag --email-list=<uuid> --tag=<uuid>
mailcoach update-tag --email-list=<uuid> --tag=<uuid> --input '{"name": "premium"}'
mailcoach delete-tag --email-list=<uuid> --tag=<uuid>
```

## Segments

Segments are dynamic groups of subscribers based on tag rules.

```bash
# List segments
mailcoach list-segments --email-list=<uuid>

# Create a segment with tag-based rules
mailcoach create-segment --email-list=<uuid> --input '{
  "name": "Engaged VIPs",
  "all_positive_tags_required": true,
  "positive_tags": ["vip", "engaged"],
  "negative_tags": ["churned"]
}'

# Show, update, delete
mailcoach show-segment --email-list=<uuid> --segment=<uuid>
mailcoach update-segment --email-list=<uuid> --segment=<uuid> --input '{"name": "New Name"}'
mailcoach delete-segment --email-list=<uuid> --segment=<uuid>
```

Segments use positive tags (subscribers must have these) and negative tags (subscribers must not have these). Set `all_positive_tags_required` to `true` to require all positive tags, or `false` to require any.
