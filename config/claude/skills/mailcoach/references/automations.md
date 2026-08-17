# Automations

Automations are configured in the Mailcoach UI (triggers, actions, delays, conditions). The CLI can trigger automations that use a **webhook trigger**.

## Trigger an automation

```bash
mailcoach trigger-automation --automation=<uuid> --input '{
  "subscribers": ["<subscriber-uuid-1>", "<subscriber-uuid-2>"]
}'
```

The automation must have a webhook trigger configured in Mailcoach. The specified subscribers will enter the automation flow.

## Finding the automation UUID

Automation management (create, list, update) is done through the Mailcoach web UI. You can find the automation UUID in the URL when viewing an automation, or check the Mailcoach documentation for your instance.
