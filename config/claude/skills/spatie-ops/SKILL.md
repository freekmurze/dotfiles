---
name: spatie-ops
description: Operate Spatie and Freek's hosting with private credentials. Use for requests to move a site to Laravel Cloud, inspect or change DigitalOcean (DO) droplets, manage Forge servers, deploy or troubleshoot Cloud apps, inspect hosting costs or performance, or change Openprovider DNS. Provides authentication and routes to the relevant provider or migration skill.
---

# Spatie ops

Start here for hosting, servers, DNS and migrations. A request such as "move this site to Cloud", "delete this DO droplet", "what runs on this Forge server?" or "check Cloud performance" should reach this skill without the user naming it.

## Choose the workflow

| Request | Read next |
| --- | --- |
| Move a site from Forge/DO to Cloud, switch DNS, retire its old server | [move-site-to-laravel-cloud](../move-site-to-laravel-cloud/SKILL.md) |
| Cloud app, deployment, environment, compute, bucket, database, cost or metrics | [laravel-cloud-ops](../laravel-cloud-ops/SKILL.md) |
| Openprovider DNS changes or propagation | [openprovider-dns](../openprovider-dns/SKILL.md) |
| DO or Forge inventory, server identity or retirement | [DO and Forge](references/servers.md) |
| Application exceptions or request traces | [flare](../flare/SKILL.md) |

Use `spatie-guidelines` when changing PHP/Laravel code and `gh` for GitHub operations.

## Authenticate and inspect

The public helper is `~/.dotfiles/bin/spatie-ops`. The private mapping is `~/.dotfiles-custom/spatie-ops/credentials.json`; it contains credential references, not token values. Read [authentication](references/authentication.md) if setup is missing, access fails, or credentials are needed inside another program.

```sh
OPS="$HOME/.dotfiles/bin/spatie-ops"
"$OPS" status
"$OPS" request do GET '/droplets?per_page=100'
"$OPS" request forge GET /servers
"$OPS" request cloud GET '/applications?per_page=100'
"$OPS" request openprovider GET '/dns/zones/example.com?with_records=true'
```

Paths are relative to each provider's API prefix. The private Forge configuration supplies the organization. Responses have `status` and `response` fields; secrets are redacted before display. The helper performs authenticated HTTPS requests itself, so bearer tokens are not placed in shell arguments.

For a JSON mutation, supply a file or stdin with `--input`; use values from 1Password in memory when constructing sensitive payloads. For example, once the user has authorized the instance change:

```sh
printf '%s' '{"hibernation_timeout":1}' | "$OPS" request cloud PATCH /instances/INSTANCE_ID --input -
```

Use the `cloud` CLI for deployments, command execution and other purpose-built Cloud commands. Use this helper for REST operations and metrics. Do not depend on migration-workspace helpers such as `secret.sh`, `cloudapi.sh` or `openprovider-token.sh`.

## Scope and verification

- Credential access establishes access, not permission to change infrastructure. Follow the user's current and standing authorization; resolve routine reads independently.
- Before a destructive operation, establish the exact named resource, ownership and dependencies. Verify its absence at the provider after deletion.
- Read back changes and verify the resulting behavior. An accepted API request alone is insufficient.
- Keep actual item IDs, vault names, credentials and account-specific mappings outside public dotfiles. API responses can contain newly issued secrets, including responses from apparently unrelated operations such as a cache purge.
- For secret-bearing data needed by automation, use `--private-output` outside Git and inspect it programmatically. See [authentication](references/authentication.md). Do not print the raw file into a chat.
