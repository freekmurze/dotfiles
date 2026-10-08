# DigitalOcean and Forge

Use `~/.dotfiles/bin/spatie-ops request PROVIDER METHOD PATH`. Authentication and the Forge organization come from the private mapping.

## Identify the server

```sh
OPS="$HOME/.dotfiles/bin/spatie-ops"
"$OPS" request do GET '/droplets?per_page=100'
"$OPS" request do GET /droplets/DROPLET_ID
"$OPS" request forge GET /servers
"$OPS" request forge GET /servers/SERVER_ID
"$OPS" request forge GET /servers/SERVER_ID/sites
```

Match Forge's provider and provider identifier against the DO droplet ID, plus its name and IP addresses. Forge server IDs are a different namespace. A Forge server supplied by Laravel or another provider need not exist in our DO account. Follow pagination to establish absence; an item missing from the first page proves nothing.

For site inventory or migration, read [move-site-to-laravel-cloud](../../move-site-to-laravel-cloud/SKILL.md). SSH uses the owner's existing SSH configuration/agent. Verify site files, databases, cron, daemons and other services before treating a droplet as dedicated to one website.

## DNS and retirement

DO DNS uses `/domains` and `/domains/DOMAIN/records`. Identify authoritative nameservers before changing a zone: a domain may have an inactive shadow zone at another provider. For Openprovider, read [openprovider-dns](../../openprovider-dns/SKILL.md).

Prefer deleting a managed server through Forge once that exact deletion is authorized. On the organization API, send `DELETE /servers/SERVER_ID` with `{"preserve_at_provider":false}` using `--input`. Then verify Forge no longer lists it and `GET /droplets/DROPLET_ID` on DO returns 404. An asynchronous 202 from Forge is only acceptance, not proof of removal.

Retain verified backups and clear DNS dependencies before deletion. Follow the migration skill's decommissioning checks. Powering off a DO droplet still incurs compute charges until it is deleted.
