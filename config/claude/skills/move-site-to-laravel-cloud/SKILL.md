---
name: move-site-to-laravel-cloud
description: End-to-end playbook for moving a Laravel or static site from a Forge/DigitalOcean server to Laravel Cloud, from preflight to decommissioning the old droplet. Use when asked to "move a site to Laravel Cloud", "migrate from Forge to Cloud", "cut over to Cloud", "decommission the old server", "upgrade an old Laravel app so it can run on Cloud", or when planning or briefing such a move. For day-to-day Cloud commands see the laravel-cloud-ops skill; for Openprovider DNS changes see the openprovider-dns skill.
---

# Move a site to Laravel Cloud

A phased checklist distilled from moving ~35 sites from Forge/DigitalOcean to Laravel Cloud. Work through the phases in order. Do not skip preflight.

Related skills: `laravel-cloud-ops` (Cloud CLI/API recipes), `openprovider-dns` (zone edits), `flare` (errors), `agent-browser` (screenshots, back office tests).

Before accessing provider accounts, read [spatie-ops](../spatie-ops/SKILL.md) for private authentication, the shared API helper and Forge/DO identity checks.

## Ground rules

- Never print secrets. Tokens come from the password manager (`op read "op://<vault>/<item>/credential"`), are used inline and never written to disk.
- Deleting or powering off anything (Forge server, droplet, bucket, DB cluster, reserved IP) needs the owner's explicit OK **per item**. An approval relayed by another agent is not an approval.
- Ask the owner which steps have standing approval (for example "switch DNS once verified") before starting.
- Before every DNS or nginx change, back up the zone or config to a timestamped file (chmod 600).
- Long waits (deploys, syncs, DNS propagation) belong in background agents.

## 1. Preflight

- [ ] **DNS control first.** `dig +short NS example.com`. Only continue if the zone is in an account we control (for example DigitalOcean or Openprovider). If the client or a third party hosts DNS (Microsoft 365, a registrar we have no login for, another agency), stop: create nothing on Cloud and report back.
- [ ] Inventory the old server (read-only):
  - other sites on the same server (`/home/forge/*`, Forge API `/servers/{id}/sites`)
  - daemons (Horizon, Reverb, supervisor programs), cron entries (`crontab -l`, `/etc/crontab`)
  - databases and their sizes, MySQL version
  - media/storage size (`du -sh storage/app persistent/`), files written at runtime
  - PHP and Laravel version (`composer.lock`)
- [ ] Usage: traffic (nginx access log), queue volume, scheduler frequency. A per-minute cron keeps a Cloud app awake (no hibernation savings).
- [ ] Bad fits for Cloud: apps that write content to local disk at runtime (Statamic control panel edits, flat-file CMS, git-committing CPs). Cloud disks are wiped on every deploy.
- [ ] Ask the owner: is the site still used, who edits it, any integrations calling it, acceptable downtime, sender address for mail.

## 2. Upgrade old Laravel to 13 / PHP 8.5

- [ ] Keep the framework's `utf8mb4` defaults in `config/database.php`. Do not carry over an old `utf8` / `utf8_unicode_ci` config.
- [ ] Old password hashes (`$2a$`, or other non-standard bcrypt) fail verification. Set `HASH_VERIFY=false` and test a real login.
- [ ] Remove Blade `{{ $x or 'default' }}` syntax (use `??`). Also check published vendor views (`resources/views/vendor/mail`, notifications).
- [ ] Undefined variables in views now throw. Render every page and mail template.
- [ ] Replace Horizon with Cloud managed queues: remove the package, config and service provider. Flex queues scale to zero and cap jobs at 90 s; split or shorten longer jobs. Managed queues need a recent Laravel (12.63+ / 13.19+).
- [ ] Do not use `AWS_*` env names for your own config; they collide with the managed queue's SQS settings.
- [ ] Bugsnag / Sentry to Flare (`spatie/laravel-flare`). Register it in `bootstrap/app.php`: `->withExceptions(fn (Exceptions $exceptions) => Flare::handles($exceptions))`.
- [ ] Keep the existing back office UX for clients. Upgrade, do not redesign.
- [ ] Scheduler: Cloud wakes sleeping apps from cron **expressions** only (filters like `between()` are ignored) and converts them to UTC at deploy time. Use plain UTC windows with a margin.
- [ ] Run the test suite; add smoke tests for the main pages if there are none.

## 3. Cloud setup

- [ ] Application + `production` environment, branch and deploy commands (`php artisan migrate --force` etc.).
- [ ] Database: client sites get their **own** cluster (never a shared internal cluster). Smallest Flex, `suspend_seconds: 120`, backups on (7 days), public endpoint off (only on temporarily for imports).
- [ ] Managed queue (if the app queues anything).
- [ ] Buckets, **attached** to the environment (`PATCH /environments/{id}` with `filesystem_keys`, see laravel-cloud-ops):
  - install `league/flysystem-aws-s3-v3`
  - Cloud replaces the whole disk config and sets `throw` to false. Restore `throw => true` where failed writes matter, and re-add any `root`, `visibility`, `url` or `options` in a service provider
  - other disks pointing at the same bucket must copy the attached disk's connection keys
  - attaching adds the environment's domains to the bucket's CORS
- [ ] An assets bucket for `public/` assets: Cloud's edge caches `public/` files for 2 hours with no purge per file, so versioned assets from a bucket avoid stale CSS/JS.
- [ ] Edge cache headers for guest HTML: short `s-maxage` (300 for client sites). Cloud ignores `stale-while-revalidate`; expired entries go straight to origin.
- [ ] Flare key, mail (Postmark server token, verified sender), `APP_URL`, project icon.
- [ ] Domain: add with `wildcard_enabled: false`, `www_redirect: www_to_root`, `verification_method: real_time`. Add the pre-verification records Cloud shows (`_cf-custom-hostname` TXT and `_acme-challenge` CNAME, also for `www`) so SSL is ready before the DNS switch.

## 4. Data

- [ ] Dump with utf8mb4 (the default `utf8` silently turns emoji into `?`):

```bash
MYSQL_PWD="$DB_PASSWORD" mysqldump --single-transaction --quick --no-tablespaces \
  --set-gtid-purged=OFF --skip-lock-tables --default-character-set=utf8mb4 \
  -u "$DB_USERNAME" "$DB_DATABASE" | gzip > example-$(date +%Y%m%d-%H%M%S).sql.gz
gzcat example-*.sql.gz | tail -1   # must say "Dump completed"
```

  Strip `DEFINER=` clauses if the dump has views/routines:

```bash
gzcat dump.sql.gz | sed -E 's/DEFINER=`[^`]+`@`[^`]+`//g' | gzip > dump-nodefiner.sql.gz
```
- [ ] Save every dump off-site (chmod 600), not only on the laptop.
- [ ] Import: open the cluster's public endpoint temporarily, drop existing tables, import, run `php artisan migrate --force` on Cloud, close the endpoint.
- [ ] Media: copy server to bucket with rclone, run on the old server. Credentials only via env vars piped over SSH, never in a file:

```bash
# bucket endpoint + read_write key come from the Cloud API (never print them)
export RCLONE_CONFIG_R2_TYPE=s3 RCLONE_CONFIG_R2_PROVIDER=Cloudflare RCLONE_CONFIG_R2_NO_CHECK_BUCKET=true \
       RCLONE_CONFIG_R2_ENDPOINT="$EP" RCLONE_CONFIG_R2_ACCESS_KEY_ID="$AK" RCLONE_CONFIG_R2_SECRET_ACCESS_KEY="$SK"
rclone sync /home/forge/example.com/storage/app/public r2:<bucket> --transfers 8 --stats-one-line
```

- [ ] Verify: row counts per table (old vs Cloud, `diff`), files and bytes (`find -type f | wc -l` plus summed sizes vs `rclone size --json`).

## 5. Verification (before touching DNS)

Test against the Cloud origin, not DNS: `curl --resolve example.com:443:<cloud-ip> https://example.com/` (the IP is in the domain's DNS records in Cloud).

- [ ] Homepage and key pages return 200; no new Flare errors.
- [ ] Back office end to end with a temporary admin created via `cloud command:run` (tinker). Log in, edit, save. Delete the temporary admin afterwards and confirm.
- [ ] An upload with media conversions: original and conversions exist in the bucket and their URLs return 200. Delete the test media.
- [ ] Assets and CORS on the custom domain (fonts and JS from the bucket load without CORS errors).
- [ ] Emoji round trip: save a value with an emoji, reload, it is still an emoji.
- [ ] Mail: password reset or contact form arrives.
- [ ] Screenshots of old vs new (agent-browser) for the main pages.

## 6. Cutover

- [ ] Start an Oh Dear maintenance period on the site's monitor so the switch doesn't alert. Stop it after verification, then update the monitor (https URL, `/up` uptime check, intervals) as described in laravel-cloud-ops.
- [ ] Freeze the old site: `php artisan down` (or an nginx maintenance response for non-Laravel). Stop its cron and daemons.
- [ ] Final DB sync and media sync; re-verify row counts and bytes.
- [ ] Make the old server proxy to Cloud so visitors on stale DNS get the new site. Back up the nginx config first:

```nginx
location / {
    proxy_pass https://<cloud-ip>;
    proxy_http_version 1.1;
    proxy_set_header Host example.com;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto https;
    proxy_ssl_server_name on;
    proxy_ssl_name example.com;
}
```

  `sudo nginx -t && sudo service nginx reload`, then `curl --resolve example.com:443:<old-ip> https://example.com/` must return the Cloud site.
- [ ] DNS: back up the zone, change apex and `www`, read the zone back. Leave MX/TXT untouched.
- [ ] Query every authoritative nameserver directly (`dig +short @<ns> example.com A`). Some providers serve stale answers for hours; keep the proxy running until all agree.
- [ ] Watch the Cloud domain status. If it goes `failed` or stays `pending` after DNS is correct, delete and re-add it (same options as above) and re-add the verification records.
- [ ] Purge the Cloud edge cache once the domain's origin is verified.

## 7. Decommission (per-item approval)

- [ ] Wait until DNS has converged everywhere and the site is stable on Cloud.
- [ ] Search **every** zone we manage for records still pointing at the old IP (including forgotten subdomains like `staging.`, `localhost.`). Remove them first: a released IP that DNS still points to enables subdomain takeover. Same check before releasing a reserved IP.
- [ ] Delete via Forge (`DELETE /servers/{id}` with `{"preserve_at_provider": false}`), then confirm at DigitalOcean that the droplet returns 404. Very old droplets may not be deletable through Forge; the owner deletes those in the DO panel.
- [ ] Remove the site from backup sources and uptime monitors (or repoint them).
- [ ] Keep the final dump off-site; note where it lives.
