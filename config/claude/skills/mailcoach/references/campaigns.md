# Campaigns

## Campaign lifecycle

Campaigns follow a clear workflow: **create draft** -> **set content** -> **send test** -> **send** -> **view statistics**.

## 1. Create a draft campaign

```bash
mailcoach create-campaign --input '{
  "name": "March Newsletter",
  "email_list_uuid": "<email-list-uuid>",
  "template_uuid": "<template-uuid>",
  "segment_uuid": "<segment-uuid>",
  "subject": "Your March Update",
  "from_email": "newsletter@example.com",
  "from_name": "My Company"
}'
```

Only `name` and `email_list_uuid` are required. Use `--json` to capture the campaign UUID from the response.

## 2. Set or update content

```bash
mailcoach update-campaign --campaign=<uuid> --input '{
  "html": "<html><body><h1>Hello {{first_name}}</h1><p>Your monthly update.</p></body></html>",
  "subject": "Updated Subject Line"
}'
```

You can update any campaign field as long as it's still in draft status.

## 3. Send a test email

```bash
mailcoach send-campaign-test --campaign=<uuid> --input '{
  "email": "test@example.com"
}'
```

Sends the campaign to a test recipient. The campaign must be in draft status. You can send to up to 10 test recipients.

## 4. Send the campaign

```bash
mailcoach send-campaign --campaign=<uuid>
```

This triggers sending to all subscribers in the target list/segment. The campaign must be in draft status. This action is irreversible.

## List and filter campaigns

```bash
# List all campaigns (default sorted by most recently sent)
mailcoach list-campaigns

# Filter by status
mailcoach list-campaigns --filter-status=sent

# Filter by email list
mailcoach list-campaigns --filter-email-list=<uuid>
```

Available status filters: `draft`, `sending`, `sent`.

## View campaign statistics

After sending, you can view engagement data:

```bash
# Opens (grouped by subscriber, with open count and timestamps)
mailcoach list-campaign-opens --campaign=<uuid>

# Link clicks (tracked links with unique and total click counts)
mailcoach list-campaign-clicks --campaign=<uuid>

# Unsubscribes
mailcoach list-campaign-unsubscribes --campaign=<uuid>

# Bounces (grouped by subscriber and bounce type)
mailcoach list-campaign-bounces --campaign=<uuid>
```

The `show-campaign` command also includes aggregate statistics:

```bash
mailcoach show-campaign --campaign=<uuid> --json
```

## Delete a campaign

```bash
mailcoach delete-campaign --campaign=<uuid>
```
