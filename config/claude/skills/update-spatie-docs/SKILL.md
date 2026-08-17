---
name: update-spatie-docs
description: "Update Spatie package documentation on spatie.be by re-importing docs for a given repo. Use after merging a PR that touches docs/, or when the user says 'update docs for [package]', 'rebuild docs for [repo]', 'import docs for [repo]', or invokes /update-spatie-docs [repo-name]. Takes a repo name argument (e.g., 'backup', 'laravel-medialibrary', 'laravel-pdf')."
---

# Update Spatie Docs

The spatie.be docs site does not rebuild automatically after a PR merges. To publish doc changes, run the `docs:import` artisan command against the repo.

## Command

Run this twice to be sure (the second run catches any version that lagged on the first pass):

```bash
cd ~/dev/code/spatie.be && php artisan remote "docs:import --repo=spatie/{repo_name}"
cd ~/dev/code/spatie.be && php artisan remote "docs:import --repo=spatie/{repo_name}"
```

`php artisan remote` runs the command against the production spatie.be instance. Write it out in full. The `ar` alias is equivalent when typed by hand, but aliases are not defined in non-interactive shells, so a subagent or scheduled run would fail silently.

If no repo name argument is provided, ask the user for one.

## Verification

After running, curl a docs URL to confirm the new content is live. The import sometimes returns success before the published pages reflect it, so check a specific page you changed:

```bash
curl -sL "https://spatie.be/docs/{repo_name}/v2/some-page-you-changed" | grep -oE '<title>[^<]*</title>'
```

If the title matches the fallback (Introduction page) instead of the page you expected, the cache has not caught up yet. Wait a moment and curl again.
