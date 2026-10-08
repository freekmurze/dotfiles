---
name: openprovider-dns
description: Read and change DNS zones hosted at Openprovider through its REST API, including its quirks (prio on original_record, one record change per call, stale anycast nameservers, TTL clamping). Use when asked to "change DNS at Openprovider", "add/remove/update a DNS record", "point a domain to Laravel Cloud", "back up a zone", "check whether the nameservers serve the new record", or when Openprovider nameservers return stale or inconsistent answers.
---

# Openprovider DNS

API base: `https://api.openprovider.eu/v1beta`. Nameservers: `ns1.openprovider.nl`, `ns2.openprovider.be`, `ns3.openprovider.eu`.

## Rules

- Never print the login or the bearer token. Fetch the token inline, never write it to a file.
- Back up the zone before every change and read it back after every change.
- Make **one** record change per API call.

## Auth

Read [spatie-ops](../spatie-ops/SKILL.md) first. Its helper reads the private 1Password mapping, logs in with the API user and retains the resulting token in memory. The token is valid for about 48 hours, but the helper does not persist it between commands.

```bash
OPS="$HOME/.dotfiles/bin/spatie-ops"
```

## Read a zone (and back it up)

```bash
ZONE=example.com
BACKUP="$HOME/.dotfiles-custom/spatie-ops/backups/$ZONE-zone-backup-$(date +%Y%m%d-%H%M%S).json"
"$OPS" request openprovider GET "/dns/zones/$ZONE?with_records=true" \
  --private-output "$BACKUP"

jq -r '.data.records[] | select(.type != "SOA" and .type != "NS") | [.type, .name, .value, .prio, .ttl] | @tsv' "$BACKUP"
```

GET returns full names (`www.example.com`). PUT uses names relative to the zone (`""` for the apex, `"www"`).

## Change records

All changes are a `PUT /dns/zones/{zone}` with `add`, `remove` or `update` under `records`.

```bash
op_put() { printf '%s' "$1" | "$OPS" request openprovider PUT "/dns/zones/$ZONE" --input -; }

# add
op_put '{"name":"example.com","records":{"add":[{"type":"A","name":"","value":"203.0.113.10","ttl":900}]}}'

# remove: must match the existing record exactly, including ttl and "prio":0
op_put '{"name":"example.com","records":{"remove":[{"type":"A","name":"","value":"198.51.100.20","ttl":86400,"prio":0}]}}'

# update: original_record must include "prio":0 (also for A/CNAME/TXT), or it does not match
op_put '{"name":"example.com","records":{"update":[{
  "original_record":{"type":"A","name":"www","value":"198.51.100.20","ttl":86400,"prio":0},
  "record":{"type":"A","name":"www","value":"203.0.113.10","ttl":900}}]}}'
```

Quirks:

- `original_record` and `remove` entries need `"prio":0` for non-MX records.
- Multi-record `remove` calls can return success while removing nothing. Adding an A record can leave the old A in place, so a host ends up with two A records. Do one change per call and read the zone back after each.
- To move a host to a new IP: remove the old record (one call), check it is gone, add the new one (one call), read back.
- TTLs can be clamped: values below 600 are raised to 600, and an existing 86400 sometimes stays 86400 after an update. Lower the TTL a day ahead if you need a fast switch, and verify the TTL in the read-back.

## Verify what the world sees

```bash
for ns in ns1.openprovider.nl ns2.openprovider.be ns3.openprovider.eu; do
  echo "$ns: $(dig +short @$ns example.com A | tr '\n' ' ') www: $(dig +short @$ns www.example.com A | tr '\n' ' ')"
done
dig +short @ns1.openprovider.nl example.com SOA   # compare serials
```

## Stale or inconsistent nameservers

The anycast nameservers can serve stale or mixed answers for hours after a change, even when the API read-back is correct and the SOA serial matches (one nameserver was seen serving a two-year-old serial).

1. Keep the old server proxying to the new origin meanwhile, so stale answers still reach the right site.
2. Recreate the affected records: remove (one call), confirm gone, add again (one call). This sometimes forces a republish.
3. If nameservers still disagree after a few hours, open an Openprovider support ticket listing the zones, the nameservers, and what each returns.
4. Do not decommission the old server until every nameserver returns the new records.
