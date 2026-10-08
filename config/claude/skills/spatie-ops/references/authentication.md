# Authentication and private configuration

The helper needs Python 3 and the 1Password `op` CLI. It reads `~/.dotfiles-custom/spatie-ops/credentials.json`, or `SPATIE_OPS_CONFIG` if set. The file must be owned by the current user and have mode `600`; its directory should have mode `700`.

The checked-in [example](credentials.example.json) contains placeholders only. Keep the real mapping outside the dotfiles repository, including vault names, item identifiers and the Forge organization. Actual credentials remain in 1Password or the Cloud CLI's existing private configuration. Copy the private mapping separately when setting up a new Mac.

Run `~/.dotfiles/bin/spatie-ops status` to check configuration without reading credentials. A request may prompt the owner to authorize 1Password access. If authorization times out, report the pending unlock; do not print credentials, search unrelated vaults or silently use a different account.

## Sources

- `digitalocean` and `forge`: `token_ref`, resolved with `op read` and captured in memory.
- `cloud`: first use the configured `cli_config` if it has an authenticated token, otherwise resolve `token_ref`. Use `cloud auth` to repair an invalid CLI login. The helper does not renew or rotate credentials.
- `openprovider`: resolve `username_ref` and `password_ref`, or read only the username/password fields of the private `login_item`. POST them to `/v1beta/auth/login` and retain the resulting bearer token in memory.

Aliases `do`, `digital-ocean` and `laravel-cloud` are accepted. Official API hosts are fixed in the helper, and authenticated redirects are not followed. Inputs cannot redirect credentials to another host.

## Output

Normal output is JSON with `status` and `response`. Sensitive keys, environment-variable values and known credential values are redacted recursively. Use `--fields data.id,data.attributes.name` to narrow the already-redacted response. Error output omits subprocess output and raw HTTP bodies; a provider 401 or 403 is reported as an unsuccessful request.

If a task needs original environment values or newly created bucket keys, capture them privately:

```sh
~/.dotfiles/bin/spatie-ops request cloud GET /environments/ENVIRONMENT_ID \
    --private-output "$HOME/.dotfiles-custom/spatie-ops/environment.json"
```

The unredacted file is written with mode `600`, outside Git, without following file symlinks. Use a program to consume just the needed values, keep them out of command arguments and logs, and delete the file when finished. For programs that can consume data entirely in memory, import `Client` from `scripts/ops.py` and use `Client.from_file().request(provider, method, path, payload)`; sanitize any displayed results with `client.redact(...)`.

Arbitrary log messages can contain sensitive application data. Use narrow fields and avoid dumping whole logs or raw configuration into a conversation, even when known credential keys are redacted.
