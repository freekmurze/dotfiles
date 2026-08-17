# Subscribers

## List subscribers

```bash
# All subscribers on a list
mailcoach list-subscribers --email-list=<uuid>

# Search and filter
mailcoach list-subscribers --email-list=<uuid> \
  --filter-search="john" \
  --filter-status=subscribed \
  --sort=-created_at
```

Available status filters: `subscribed`, `unconfirmed`, `unsubscribed`.

## Create a subscriber

```bash
mailcoach create-subscriber --email-list=<uuid> --input '{
  "email": "john@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "tags": ["newsletter", "vip"],
  "extra_attributes": {
    "company": "Acme Inc"
  },
  "skip_confirmation": true
}'
```

Set `skip_confirmation` to `true` to subscribe immediately without sending a confirmation email.

## Show, update, delete

```bash
mailcoach show-subscriber --subscriber=<uuid>

mailcoach update-subscriber --subscriber=<uuid> --input '{
  "first_name": "Jane",
  "tags": ["premium"],
  "append_tags": true
}'

mailcoach delete-subscriber --subscriber=<uuid>
```

Set `append_tags` to `true` to add tags without removing existing ones. When `false` (default), the provided tags replace all existing tags.

## Tag management

```bash
# Add tags to a subscriber
mailcoach add-subscriber-tags --subscriber=<uuid> --input '{
  "tags": ["vip", "early-adopter"]
}'

# Remove specific tags
mailcoach remove-subscriber-tags --subscriber=<uuid> --input '{
  "tags": ["churned"]
}'
```

## Subscriber actions

```bash
# Confirm an unconfirmed subscriber
mailcoach confirm-subscriber --subscriber=<uuid>

# Unsubscribe
mailcoach unsubscribe-subscriber --subscriber=<uuid>

# Resubscribe a previously unsubscribed subscriber
mailcoach resubscribe-subscriber --subscriber=<uuid>

# Resend the confirmation email
mailcoach resend-confirmation --subscriber=<uuid>
```

These actions are also available via the email list:

```bash
mailcoach confirm-subscriber-by-list --email-list=<uuid> --input '{"email": "john@example.com"}'
mailcoach unsubscribe-by-list --email-list=<uuid> --input '{"email": "john@example.com"}'
mailcoach resubscribe-by-list --email-list=<uuid> --input '{"email": "john@example.com"}'
```

## Bulk import

For importing many subscribers at once:

```bash
# 1. Create an import job
mailcoach create-subscriber-import --email-list=<uuid> --input '{
  "subscribers_csv": "email,first_name,tags\njohn@example.com,John,\"vip,newsletter\"",
  "subscribe_unsubscribed": false,
  "unsubscribe_others": false,
  "replace_tags": false
}'

# 2. Optionally append more rows
mailcoach append-subscriber-import --subscriber-import=<uuid> --input '{
  "subscribers_csv": "jane@example.com,Jane,premium"
}'

# 3. Start the import
mailcoach start-subscriber-import --subscriber-import=<uuid>
```

Import flags:
- `subscribe_unsubscribed`: re-subscribe previously unsubscribed emails
- `unsubscribe_others`: unsubscribe emails not in the CSV
- `replace_tags`: replace existing tags instead of appending

## Suppressions

Globally suppressed email addresses are blocked from all lists.

```bash
mailcoach list-suppressions
mailcoach create-suppression --input '{"email": "spam@example.com", "reason": "complained"}'
mailcoach delete-suppression --suppression=<uuid>
```
