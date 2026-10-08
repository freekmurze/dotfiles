---
name: laravel-cloud-ops
description: Day-to-day Laravel Cloud operations with the `cloud` CLI and the Cloud REST API, plus the conventions for setting up a new Cloud app. Use when asked to create or set up a new Cloud app or environment, look up a Cloud app or environment, read or set env vars, deploy or check a deployment, run tinker or an artisan command on Cloud, fix a stuck or failed domain, purge the edge cache, attach a bucket, check usage, cost or bandwidth, report which apps sleep or stay awake, change hibernation or instance size, configure a database cluster (suspend, backups, public endpoint), or set up deploy notifications.
---

# Laravel Cloud ops

Two tools:

- `cloud` CLI (`cloud list`, `cloud <command> --help`). Most commands take an environment ID or name and support `--json`.
- REST API, base `https://cloud.laravel.com/api`, bearer token. Use it for things the CLI lacks.

```bash
cloudapi() {  # usage: cloudapi METHOD PATH [non-secret JSON]
  if [ -n "${3:-}" ]; then
    printf '%s' "$3" | "$HOME/.dotfiles/bin/spatie-ops" request cloud "$1" "$2" --input -
  else
    "$HOME/.dotfiles/bin/spatie-ops" request cloud "$1" "$2"
  fi
}
```

Read [spatie-ops](../spatie-ops/SKILL.md) first for private authentication and safe API output. The helper uses the Cloud CLI login or a private 1Password reference; its response is wrapped in `status` and `response`. For sensitive payloads, construct JSON in memory and send it on stdin rather than passing secret values as shell arguments. Use `--private-output` or the in-memory client when automation needs unredacted data.

## Never print secrets

- `cloud ... --json` masks env var values. `--show-sensitive` reveals them: only use it piped into a file (chmod 600) or a script, never into the conversation.
- `GET /buckets/{id}/keys` and `/bucket-keys/{id}` return `access_key_secret`. Only print `id`, `name`, `permission`.
- Before deleting env vars, save their values to a chmod 600 backup file without printing them.
- Any API response containing an environment can include its secret values, including a cache-purge response. Use the shared helper's redacted output instead of displaying raw HTTP JSON.

## Setting up a new app (conventions)

The defaults we settled on after moving ~40 apps. Deviate only with a reason. For moving an existing site, follow the move-site-to-laravel-cloud skill; this section is the target setup.

**Naming and source**
- App name = primary domain (`example.com`). Name the GitHub repo after the domain too. If you rename the repo later, point the app at it: `PATCH /applications/<id>` with `{"repository":"org/example.com","source_control_provider_type":"github"}`, then trigger a deploy to prove it works.
- One environment, `production`, deploying `main` with push to deploy. No staging unless asked.
- Region closest to the users (`eu-central-1` for Belgian sites).
- Set the app icon (logo) in the dashboard so it is recognisable on the canvas.

**Runtime**
- Latest Laravel and PHP 8.5. Upgrade old apps first rather than running legacy versions on Cloud.
- Build: `composer install --no-dev --no-interaction --prefer-dist --optimize-autoloader`, then `npm ci && npm run build`, then `php artisan optimize`. Deploy command: `php artisan migrate --force`.
- Cloud deletes `node_modules` after the build. If PHP shells out to a Node tool at runtime (for example Shiki), install it into a separate directory during the build.
- Compute: the smallest Flex size for plain sites. Apps that process images (Glide, media conversions on large originals) or run Statamic need 1 GB; 512 MB runs out of memory. Hibernation on with a 2 minute timeout.
- Keep requests short. We hit a 20 s gateway timeout on cold, expensive pages and a hard limit on long downloads. Anything slow goes to a queued job; cache warm expensive pages before the DNS switch.
- `cloud command:run` runs inside the web container. Heavy one-offs (cache warming, image generation, big imports) compete with web traffic and can OOM the app: dispatch them to the queue, and never run two in parallel.
- Queues: Cloud managed queue on Flex (scales to zero, 90 s job cap). Remove Horizon. Split long jobs into chained or self-dispatching chunks.
- Scheduler: Cloud wakes the app for every cron expression, so a per-minute task means the app never sleeps. Prefer daily or hourly schedules.
- No headless Chromium on Cloud. Use Cloudflare Browser Rendering (REST `/browser-rendering/content`) instead.

**Database**
- Internal and personal apps share one cluster per owner (one for the company's internal apps, one per person's hobby projects). Client sites always get their own cluster.
- MySQL Flex, smallest size, `suspend_seconds: 120`, daily backups kept 7 days, public endpoint off (on only during an import).
- Keep `utf8mb4` and turn MySQL strict mode on. With strict mode off, values longer than a column (for example a 191-char string column) are silently truncated.

**Storage**
- Every bucket is attached to the environment (`filesystem_keys`), so it shows on the canvas and Cloud injects credentials. Restore `throw => true` and any `root`/`visibility`/`url` in a service provider.
- Separate public assets bucket for versioned `public/` assets if stale CSS/JS after deploys matters.
- Never write content to the local disk at runtime; it is wiped on deploy. For Statamic, use the Eloquent driver and turn git automation off.

**Edge cache**
- Guest HTML: `Cache-Control: public, max-age=60, s-maxage=300` for client sites. No `stale-while-revalidate` (ignored).
- For a longer edge TTL, purge on save: a scoped Cloud API token (purge only, one environment) in an env var, a queued job calling `POST /environments/<env>/purge-edge-cache`. Note the token's expiry date somewhere you will see it.
- Never cache the back office, logged-in responses, form posts or live preview. Send `private, no-store` there.

**Domains**
- Add with `wildcard_enabled: false`, `www_redirect: www_to_root` (apex) and `verification_method: real_time`, and add the pre-verification records for both apex and `www`, even if `www` does not resolve today. Without the `www` records the canvas shows "Not connected" although the apex works.
- Lower the TTL to 600 or less well before the switch.

**Observability and mail**
- Flare with performance tracing (sample rate 0.1, or 0.02 for busy sites; `minimal_log_level` error), deploy notifications to Slack (dashboard setting).
- Mail through a Postmark server per app with a verified sender. Store tokens in the password manager, set them with `env:variables`, never print them.

**Oh Dear**

Every request Oh Dear makes can wake a hibernating app, and an edge-cached page can hide an outage. Per app:

- Monitor URL: after the move, make sure it is `https://` and the real domain, so the certificate check sees Cloud's certificate.
- Uptime check: point it at `/up` (Laravel's health route, never edge cached). Checking an edge-cached homepage only notices downtime after the cache expires. For apps meant to sleep, use a long interval (720 or 1440 minutes); busy apps that never sleep can stay at 1 minute. Look for extra monitors on other paths of the same site, they wake it too.
- Other checks (broken links, mixed content, performance, Lighthouse) also hit the app: give them long intervals on sleeping apps.
- Application health (`spatie/laravel-health`): Oh Dear polls the endpoint every few minutes, so the app never sleeps. Use it only on apps that are awake anyway, otherwise turn it off. Drop checks that make no sense on Cloud (Redis, Horizon, used disk space).
- Scheduled task monitoring (`spatie/laravel-schedule-monitor`): the app pings Oh Dear after each task, so it wakes nothing. Set the monitor id (the package reads `monitor_id`), add `php artisan schedule-monitor:sync` to the deploy commands, and raise grace times to about 5 minutes for cold starts.
- Cutover: start a maintenance period on the monitor via the Oh Dear API before freezing the old site and stop it after verification, so the switch doesn't page anyone.
- Check the `OH_DEAR_API_TOKEN` actually works (a 401 means cron sync and maintenance windows silently fail).
- Decommissioning: delete or repoint monitors of removed sites and servers.

**Node apps**
- `X-Forwarded-Proto` is `http`; use `CF-Visitor` for the scheme and `CF-Connecting-IP` for the client IP.
- Cloud's scheduler only runs `php artisan schedule:run`. For Node, expose secret-protected job endpoints and trigger them from a GitHub Actions schedule (or a Cloudflare cron trigger).

## Lookup

```bash
cloud app:list --json
cloud env:list <app> --json
cloud env:get <env> --json            # branch, deploy commands, settings
cloud instance:list <env> --json
cloudapi GET '/applications?per_page=100'
cloudapi GET /applications/<app-id>/environments
cloudapi GET /environments/<env-id>/instances
```

## Env vars

```bash
cloud env:variables <env> --action=set --key=FOO --value=bar --force      # one var
cloudapi POST /environments/<env>/variables '{"method":"set","variables":[{"key":"FOO","value":"bar"}]}'
cloudapi POST /environments/<env>/variables/delete '{"keys":["FOO"]}'
```

- `--action=replace` replaces **all** variables. Avoid unless that is the intent.
- Changes only reach containers on the next deploy.
- Do not use `AWS_*` names for your own settings; they collide with the managed queue's SQS config.

## Deploys

```bash
cloud deploy                                     # from a repo configured with `cloud repo:config`
cloudapi POST /environments/<env>/deployments '{}'
cloud deployment:list <env> --json
cloud deployment:get <deployment-id> --json      # status: pending, ..., deployment.succeeded / failed
cloudapi GET /deployments/<deployment-id>/logs
cloud env:update <env> --branch=main --force
```

Poll `deployment:get` every 10 to 15 s until the status contains `succeeded`, `failed` or `cancel`. Run waits in a background agent.

## Commands and tinker

```bash
cloud command:run <env> --cmd="php artisan migrate:status"
cloud command:run <env> --cmd='php artisan tinker --execute="echo App\Models\User::count();"'
cloud tinker                                     # interactive
cloudapi POST /environments/<env>/commands '{"command":"php artisan about"}'   # then GET /commands/<id>
```

For larger scripts: base64-encode locally and `echo <b64> | base64 -d > /tmp/x.php && php /tmp/x.php` inside `--cmd`.

## Logs

```bash
cloud env:logs <app> <env> --from=2026-01-01T10:00:00Z --to=2026-01-01T11:00:00Z --json   # max ~100 lines per call
```

## Domains

```bash
cloud domain:list <env> --json                   # hostname_status, ssl_status, origin_status, dns_records
cloud domain:get <domain-id> --json
cloud domain:verify <domain-id>
cloud domain:create <env> --name=example.com --wildcard-enabled=false \
  --www-redirect=www_to_root --verification-method=real_time
cloudapi DELETE /domains/<domain-id>
```

- Apex domains: always pass `www_redirect` (`www_to_root`, or the reverse if the site lives on www). Without it the `www` host is dropped.
- Subdomains: `wildcard_enabled: false`, no www.
- Wildcard cannot be changed after creation; update only supports `verification_method`. Delete and re-add instead.
- Pre-verification: add the `_cf-custom-hostname` TXT and `_acme-challenge` CNAME records shown in `dns_records` (also for `www`) so SSL is issued before DNS points at Cloud.
- **Stuck or failed:** domains can sit in `pending` forever or flip to `failed` even with correct DNS. Fix: save the domain JSON, delete, re-add with the same options, re-add the verification records. Audit all domains periodically.

## Edge cache

```bash
cloudapi POST /environments/<env>/purge-edge-cache '{}'
```

- The purge may not clear anything on domains whose `origin_status` is still pending (DNS not yet at Cloud). A redeploy does not purge either.
- The edge honours `s-maxage` but ignores `stale-while-revalidate`.
- Files in `public/` get about a 2-hour edge TTL. Serve versioned assets from a bucket if that is a problem. Bucket CDNs also cache 404s for hours.
- Node apps receive `X-Forwarded-Proto: http`; use `CF-Visitor` to detect https.

## Buckets

```bash
cloud bucket:list --json
cloud bucket-key:list <bucket-id> --json         # contains secrets, filter to id/name/permission
cloudapi PATCH /environments/<env> '{"filesystem_keys":[{"id":"<bucket-key-id>","disk":"media","is_default_disk":false}]}'
cloudapi GET '/environments/<env>?include=buckets'
```

- `filesystem_keys` is the full set: include buckets that are already attached, or they get detached.
- Attached disks are rebuilt by Cloud at runtime: whole disk config replaced, `throw` false, custom `root`/`visibility`/`url`/`options` lost. Restore what you need in a service provider. Needs `league/flysystem-aws-s3-v3`.
- Attach as default only if the app's default disk already is that bucket. Attaching adds the env's domains to the bucket CORS.
- Inspect disks without secrets: `php artisan tinker --execute='foreach (config("filesystems.disks") as $n => $d) { echo $n, " ", $d["driver"] ?? "", " ", $d["bucket"] ?? "", " throw=", var_export($d["throw"] ?? null, true), PHP_EOL; }'`

## Usage and cost

```bash
cloudapi GET /usage                              # data.summary, resources, application_totals, environment_usage
cloudapi GET '/usage?period=1'                   # older billing periods; meta.available_periods lists them
```

No per-app bandwidth breakdown. For bandwidth per path, use the access logs (`data.bytes_sent` in `env:logs`).

## Metrics (sleep/awake reports)

```bash
cloudapi GET '/environments/<env>/metrics?period=7d'          # 6h, 24h, 3d, 7d, 30d
cloudapi GET '/databases/clusters/<cluster-id>/metrics?period=7d'
```

Environment metrics: `cpu_usage`, `memory_usage`, `http_response_count`, `replica_count`, `web_workers_count`. A `replica_count` of 0 means the app is hibernating; the share of non-zero points is its awake time. Cluster metrics include `compute_hours`. Common wake-up causes: uptime monitors, frequent scheduled tasks (Cloud wakes apps for cron expressions), crawlers, pollers.

## Hibernation and instance sizes

```bash
cloud instance:sizes
cloudapi PATCH /instances/<inst-id> '{"hibernation_timeout":2}'
cloudapi PATCH /instances/<inst-id> '{"size":"flex.g-1vcpu-512mb"}'
cloudapi POST /environments/<env>/stop     # or /start
cloudapi POST /instances/<inst-id>/pause   # pause a queue instance, /resume to undo
```

- Instance changes apply after the next deploy.
- The dashboard dropdown may only show legacy sizes; the API accepts the newer `flex.*` sizes.
- Flex queues scale to zero with a 90 s job cap; Pro queues keep one worker running.

## Database clusters

```bash
cloud db-cluster:list --json
cloud db:list <cluster-id> --json
cloudapi PATCH /databases/<cluster-id> '{"config":{"suspend_seconds":120}}'
cloudapi PATCH /databases/<cluster-id> '{"config":{"retention_days":7}}'
cloudapi PATCH /databases/<cluster-id> '{"config":{"is_public":true}}'    # only during an import, then false
cloud db:open <cluster> <database>                                        # local client
```

- Client apps get their own cluster; do not share clusters (they share a default user and settings).
- Deleting: delete the schemas first (`DELETE /databases/clusters/<cluster>/databases/<schema-id>`), then the cluster. Needs the owner's OK.
- Cloud's own snapshots are not off-site backups. Dump separately if you need that.

## Notifications

Slack deploy notifications are configured in the Cloud dashboard (organization or application settings), not via the API. They cover deploy success and failure only: no "deploy started", and there are no generic outgoing webhooks.
