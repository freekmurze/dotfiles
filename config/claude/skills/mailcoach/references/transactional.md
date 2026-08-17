# Transactional Emails

Transactional emails are one-off messages triggered by events (welcome emails, password resets, order confirmations, etc.).

## Send via template

Templates are created in the Mailcoach UI. Reference them by name and pass dynamic replacements:

```bash
mailcoach send-transactional-mail --input '{
  "mail_name": "welcome",
  "to": "john@example.com",
  "replacements": {
    "name": "John",
    "activation_url": "https://example.com/activate/abc123"
  }
}'
```

## Send with inline content

```bash
mailcoach send-transactional-mail --input '{
  "to": "john@example.com",
  "subject": "Your order has shipped",
  "html": "<h1>Order Shipped</h1><p>Your order #1234 is on its way.</p>",
  "from": "orders@example.com",
  "store": true
}'
```

Set `store` to `true` to log the email in Mailcoach for later review.

## Additional send options

The `send-transactional-mail` command also supports:
- `cc`, `bcc`: additional recipients
- `reply_to`: reply-to address
- `attachments`: array of file attachments
- `fake`: set to `true` to test without actually sending

## List sent transactional mails

```bash
mailcoach list-transactional-mails
mailcoach list-transactional-mails --filter-search="order shipped"
```

## View a sent mail

```bash
mailcoach show-transactional-mail-log-item --transactional-mail=<uuid>
```

## Resend a transactional mail

```bash
mailcoach resend-transactional-mail --transactional-mail=<uuid>
```

## Manage templates

```bash
# List all transactional mail templates
mailcoach list-transactional-mail-templates

# View a specific template
mailcoach show-transactional-mail-template --transactional-mail-template=<uuid>
```
