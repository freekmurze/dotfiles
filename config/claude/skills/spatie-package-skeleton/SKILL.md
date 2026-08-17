---
name: spatie-package-skeleton
description: Guide for creating PHP and Laravel packages using Spatie's package-skeleton-laravel and package-skeleton-php templates. Use when the user wants to create a new PHP or Laravel package, scaffold a package. Also use when building customizable packages — covers proven patterns for extensibility (events, configurable models/jobs, action classes) instead of config option creep.
metadata:
  author: Spatie
  tags:
    - php
    - laravel
    - package
    - skeleton
    - open-source
---

# Creating a Laravel Package with Spatie's Skeleton

## Prerequisites

- `gh` CLI installed and authenticated
- `php` available in PATH
- `composer` available in PATH

## Workflow

### 1. Gather Package Details

Ask the user for:
- **Vendor name** (e.g. `spatie`) — the GitHub org or username
- **Package name** (e.g. `laravel-cool-feature`) — the repo/package name
- **Package description** — one-liner for composer.json
- **Visibility** — public or private (default: public)

Use defaults where sensible:
- Author name: from `git config user.name`
- Author email: from `git config user.email`
- Author username: from `gh auth status`
- Vendor namespace: PascalCase of vendor name (e.g. `Spatie`)
- Class name: TitleCase of package name without `laravel-` prefix (e.g. `CoolFeature`)

### 2. Create the Repository from Template

```bash
gh repo create <vendor>/<package-name> --template spatie/package-skeleton-laravel --public --clone
cd <package-name>
```

If the user wants a private repo, use `--private` instead of `--public`.

### 3. Configure the Package (Manual Replacement)

**WARNING**: Do NOT pipe stdin to `configure.php`. The script's child processes (`gh auth status`, `git log`, `git config`) consume lines from the piped stdin, causing inputs to shift and produce garbled results. Instead, do the replacements manually:

1. Run `sed` to replace all placeholder strings across the repo:

```bash
find . -type f -not -path './.git/*' -not -path './vendor/*' -not -name 'configure.php' -exec sed -i '' \
  -e 's/:author_name/Author Name/g' \
  -e 's/:author_username/authorusername/g' \
  -e 's/author@domain\.com/author@email.com/g' \
  -e 's/:vendor_name/Vendor Name/g' \
  -e 's/:vendor_slug/vendorslug/g' \
  -e 's/VendorName/VendorNamespace/g' \
  -e 's/:package_slug_without_prefix/package-without-prefix/g' \
  -e 's/:package_slug/package-name/g' \
  -e 's/:package_name/package-name/g' \
  -e 's/:package_description/Package description here/g' \
  -e 's/Skeleton/ClassName/g' \
  -e 's/skeleton/package-name/g' \
  -e 's/migration_table_name/package_without_prefix/g' \
  -e 's/variable/variableName/g' \
  {} +
```

**Important**: The order of `-e` flags matters. Replace `:package_slug_without_prefix` before `:package_slug` to avoid partial matches. Replace `Skeleton` (PascalCase) before `skeleton` (lowercase).

2. Rename the skeleton files:

```bash
mv src/Skeleton.php src/ClassName.php
mv src/SkeletonServiceProvider.php src/ClassNameServiceProvider.php
mv src/Facades/Skeleton.php src/Facades/ClassName.php
mv src/Commands/SkeletonCommand.php src/Commands/ClassNameCommand.php
mv config/skeleton.php config/package-without-prefix.php
mv database/migrations/create_skeleton_table.php.stub database/migrations/create_package_without_prefix_table.php.stub
```

3. Delete `configure.php` and run `composer install`:

```bash
rm configure.php
composer install
```

Use a longer timeout (5 minutes) for `composer install`.

### 4. Verify Setup

After the script completes:

```bash
# Check the directory structure
ls -la src/
# Verify composer.json looks correct
cat composer.json | head -20
# Check tests passed during setup
```

### 5. Initial Commit and Push

The configure script modifies all files but doesn't commit. Create the initial commit:

```bash
git add -A
git commit -m "Configure package skeleton"
git push -u origin main
```

### 6. Report to User

Tell the user:
- The repo URL (e.g. `https://github.com/<vendor>/<package-name>`)
- The namespace (e.g. `VendorNamespace\ClassName`)
- Key files to start editing:
  - `src/<ClassName>.php` — main package class
  - `src/<ClassName>ServiceProvider.php` — service provider
  - `config/<package-slug>.php` — configuration
  - `tests/` — test directory

## Detailed references

Load these as needed:

- **Post-setup reference**: `references/post-setup.md`. What the skeleton generated and what to do with it.
- **Package patterns**: `references/package-patterns.md`. API design principles, and proven extensibility patterns (events, configurable models and jobs, action classes) instead of config option creep.
