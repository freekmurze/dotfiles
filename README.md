# My dotfiles

![Terminal](images/terminal.png)

Personal dotfiles with modern shell tooling, optimized for Laravel/PHP development. Features fast startup times, smart directory navigation, and modern CLI tools.

## Contents

**Getting started**

- [Key Features](#key-features)
- [Quick Start](#quick-start) - clone and run `bin/install`
- [What's Included](#whats-included) - [Shell & Prompt](#shell--prompt) · [Modern CLI Tools](#modern-cli-tools) · [Development Tools](#development-tools) · [QuickLook Plugins](#quicklook-plugins)

**Reference**

- [How It Works](#how-it-works) - [Symlinked Files](#symlinked-files) · [Sourced Files](#sourced-files) · [Custom Agnoster Theme](#custom-agnoster-theme)
- [Daily Usage](#daily-usage) - [Smart Navigation](#smart-navigation) · [Laravel/PHP Shortcuts](#laravelphp-shortcuts) · [Data Processing](#data-processing) · [Maintenance Commands](#maintenance-commands)
- [Version Management](#version-management) - [Node.js via fnm](#nodejs-via-fnm) · [PHP & Composer via Homebrew](#php--composer-via-homebrew)
- [Package Management](#package-management) - the Brewfile and global npm and Composer packages

**AI setup**

- [AI Development Setup](#ai-development-setup) - one config for Claude Code and Codex
  - [Quick Install (Standalone)](#quick-install-standalone) - the AI setup without the rest of the dotfiles
  - [Skills](#skills) - all 19, grouped by purpose
  - [Scoped Plugins](#scoped-plugins) - skills that load only in the repos that need them
  - [Code Intelligence](#code-intelligence) - Laravel LSP, Intelephense, TypeScript
  - [Agents](#agents) - custom subagents
  - [The Review Workflow](#the-review-workflow) - the six lanes behind `review-code` and `review-pr`
  - [Hooks](#hooks) - Pint on every edited PHP file
  - [Settings Worth Knowing](#settings-worth-knowing)
  - [Sharing With Codex](#sharing-with-codex) - how one source feeds both harnesses
  - [Adding New Skills](#adding-new-skills) - and the two rules that decide whether a skill gets used

**Everything else**

- [Customization](#customization) - [Personal Aliases & Functions](#personal-aliases--functions) · [Project-Specific Variables](#project-specific-variables)
- [Post-Installation](#post-installation) · [Tool Comparisons](#tool-comparisons) · [Utilities](#utilities) · [Migration Notes](#migration-notes) · [Credits](#credits)

---

## Key Features

- **Custom Agnoster Theme** - Clean powerline prompt with no branch symbols, `•` for changes
- **Version-Controlled AI Setup** - Skills, agents, settings, and instructions for both Claude Code and Codex, from one source
- **Framework-Aware Code Intelligence** - Laravel LSP, Intelephense, and TypeScript language servers wired into the agent
- **Fast Tools** - fnm, zoxide, ripgrep, bat, eza (all Rust-based for speed)
- **One Command Install** - `bin/install` sets up everything including Claude Code

---

## Quick Start

```bash
git clone git@github.com:freekmurze/dotfiles.git ~/.dotfiles
cd ~/.dotfiles
bin/install
```

---

## What's Included

### Shell & Prompt

- **Oh My Zsh** - Framework for managing Zsh configuration (with agnoster theme by default)
- **zoxide** - Smart directory jumping based on frecency
- **fzf** - Fuzzy finder for files and history
- **direnv** - Automatic environment variables per directory

### Modern CLI Tools

- **fnm** - Fast Node.js version manager
- **bat** - Cat with syntax highlighting
- **eza** - Modern ls replacement with icons
- **ripgrep** - Fast grep alternative
- **fd** - Fast find alternative
- **git-delta** - Better git diffs
- **jq** - JSON processor and formatter
- **yq** - YAML processor and formatter
- **bottom** - Modern system monitor

### Development Tools

- **PHP** - Latest version via Homebrew
- **Composer** - Dependency manager via Homebrew
- **Node.js** - LTS version managed via fnm
- **Laravel Valet** - Local development server
- **MySQL** - Database with auto-start

### QuickLook Plugins

Instant file previews in Finder: code files, markdown, JSON, CSV, patches, and archives.

---

## How It Works

### Symlinked Files

The installation creates symlinks from your home directory to the dotfiles repository. This allows you to version control your configuration while keeping files in their expected locations.

| Symlink Location | Points To | Purpose |
|-----------------|-----------|---------|
| `~/.zshrc` | `~/.dotfiles/home/.zshrc` | Main Zsh configuration (Oh My Zsh with custom agnoster theme) |
| `~/.gitconfig` | `~/.dotfiles/home/.gitconfig` | Git configuration with delta diff viewer |
| `~/.global-gitignore` | `~/.dotfiles/home/.global-gitignore` | Global Git ignore patterns |
| `~/.vimrc` | `~/.dotfiles/home/.vimrc` | Vim configuration |
| `~/.vim/` | `~/.dotfiles/home/.vim/` | Vim runtime files |
| `~/.mackup.cfg` | `~/.dotfiles/macos/.mackup.cfg` | Mackup backup configuration |
| `~/.claude/skills` | `~/.dotfiles/config/claude/skills/` | All Claude Code skills (version-controlled) |
| `~/.claude/agents` | `~/.dotfiles/config/claude/agents/` | All Claude Code agents (version-controlled) |
| `~/.claude/CLAUDE.md` | `~/.dotfiles/config/claude/AGENTS.md` | Agent instructions (Claude reads `CLAUDE.md`, not `AGENTS.md`) |
| `~/.claude/settings.json` | `~/.dotfiles/config/claude/settings.json` | Claude Code settings |
| `~/.codex/AGENTS.md` | `~/.dotfiles/config/claude/AGENTS.md` | The same instructions, read natively by Codex |
| `~/.codex/skills/*` | `~/.dotfiles/config/claude/skills/*` | One symlink per shared skill, see `bin/link-agent-skills` |
| `~/.config/zed/settings.json` | `~/.dotfiles/config/zed/settings.json` | Zed editor settings |
| `~/.config/zed/keymap.json` | `~/.dotfiles/config/zed/keymap.json` | Zed custom keybindings |
| `~/.config/ghostty/config` | `~/.dotfiles/config/ghostty/config` | Ghostty terminal settings |

To manually symlink the Zed configuration (if not using `bin/install`):

```bash
mkdir -p ~/.config/zed
ln -sf ~/.dotfiles/config/zed/settings.json ~/.config/zed/settings.json
ln -sf ~/.dotfiles/config/zed/keymap.json ~/.config/zed/keymap.json
```

### Sourced Files

These files are loaded by `.zshrc` but remain in the dotfiles directory:

- `home/.aliases` - Shell command aliases
- `home/.functions` - Custom shell functions
- `home/.exports` - Environment variables

### Custom Agnoster Theme

The default configuration uses a customized agnoster theme stored in `oh-my-zsh-custom/themes/agnoster.zsh-theme`:

**Customizations:**
- No git branch symbol (cleaner look)
- Uses `•` for unstaged changes instead of `±`
- Powerline arrows for segment separators
- Requires a font with powerline glyphs

**Git Status Symbols:**
- `✚` - Staged changes (files added with `git add`)
- `•` - Unstaged changes (modified files not yet staged)
- Yellow background - Uncommitted changes
- Green background - Clean working directory

---

## Daily Usage

### Smart Navigation

```bash
z dotfiles          # Jump to frequently used directories
zi                  # Interactive directory picker
Ctrl+R              # Fuzzy search command history
Ctrl+T              # Fuzzy find files
Alt+C               # Fuzzy change directory
```

### Laravel/PHP Shortcuts

```bash
a                   # php artisan
p                   # Run Pest/PHPUnit tests
c                   # composer
mfs                 # php artisan migrate:fresh --seed
nah                 # git reset --hard; git clean -df
```

### Data Processing

```bash
# JSON processing with jq
curl api.github.com/users/freekmurze | jq
cat composer.json | jq '.require'
php artisan tinker --execute="echo json_encode(User::first());" | jq

# YAML processing with yq
yq '.jobs' .github/workflows/ci.yml
yq -o json docker-compose.yml

# System monitoring
btm                 # Modern system monitor (aliased from top/htop)
```

### Maintenance Commands

```bash
bin/update          # Update all packages and tools
```

---

## Version Management

### Node.js (via fnm)

```bash
fnm install --lts     # Install latest LTS
fnm use lts-latest    # Use latest LTS
fnm install 20        # Install specific version
fnm use 20            # Switch to specific version
fnm list              # Show installed versions
```

### PHP & Composer (via Homebrew)

```bash
brew upgrade php      # Update PHP to latest
brew upgrade composer # Update Composer
```

---

## Package Management

All Homebrew packages are declared in `config/Brewfile`. To add a new tool:

```bash
echo 'brew "neovim"' >> ~/.dotfiles/config/Brewfile
brew bundle --file=~/.dotfiles/config/Brewfile
```

**Complete package list:**

- **Core**: node, php, composer, pkg-config, wget, httpie, ncdu, hub, ack, doctl, 1password-cli, git-secret, imagemagick, mysql, yarn, ghostscript, mackup
- **Modern CLI**: zoxide, bat, eza, ripgrep, fd, git-delta, fnm, fzf, direnv, jq, yq, bottom, zsh-autosuggestions
- **QuickLook**: qlcolorcode, qlstephen, qlmarkdown, quicklook-json, qlprettypatch, quicklook-csv, betterzip, suspicious-package
- **PHP Extensions**: imagick, memcached, xdebug, redis
- **Global npm**: agent-browser, intelephense, typescript-language-server, typescript
- **Global Composer**: laravel/envoy, spatie/phpunit-watcher, laravel/valet, laravel/lsp

---

## AI Development Setup

Everything the agents need lives in `config/claude/`, and both Claude Code and Codex read it from there. Nothing is duplicated per tool.

```
config/claude/
├── AGENTS.md          the instructions, read by both harnesses
├── settings.json      Claude Code settings, permissions, hooks
├── agents/            custom subagents
└── skills/            19 skills, plus 3 scoped plugins
```

### Quick Install (Standalone)

Install just the AI setup without the full dotfiles:

```bash
curl -fsSL https://raw.githubusercontent.com/freekmurze/dotfiles/main/bin/install-claude-code | bash
```

That installs the Claude Code CLI, symlinks the config, and runs `bin/link-agent-skills` to share the harness-neutral skills with Codex.

### Skills

All skills live in `config/claude/skills/` and are version-controlled. On a new Mac they are available immediately after the installer runs.

**Code and conventions** (loaded automatically when relevant)

| Skill | What it does |
|-------|--------------|
| `spatie-guidelines` | The single source for Spatie PHP, Laravel, and JavaScript conventions. Style, docblocks, control flow, naming, validation, Blade, package architecture, Pest, git workflow |
| `react-best-practices` | React hooks, effects, refs, and component design. Scoped to `**/*.{tsx,jsx}` |
| `laravel-inertia-react-structure` | Frontend directory structure for Laravel Inertia React apps. Scoped to React files |
| `livewire-4` | Livewire 4 components, single-file and multi-file |
| `spatie-package-skeleton` | Scaffolding a package from `package-skeleton-laravel`, plus extensibility patterns |
| `speeding-up-laravel-tests` | Making slow Pest suites fast. Factories, fakes, config caching, `LazilyRefreshDatabase` |

**Review and audit**

| Skill | What it does |
|-------|--------------|
| `review-code` | Runs 6 review lanes in parallel over the working tree, then applies the fixes |
| `review-pr` | Reviews a GitHub PR, gates on CI, merges, thanks the author. Releases only when asked |
| `audit-architecture` | Whole-codebase audit of how state and data are modelled. Read-only, ranked P0 to P3 |

**Products and publishing** (CLI wrappers for things Spatie runs)

| Skill | What it does |
|-------|--------------|
| `flare` | Triage errors and performance data on flareapp.io |
| `mailcoach` | Email lists, subscribers, campaigns, automations |
| `there-there` | Helpdesk tickets, contacts, channels |
| `update-spatie-docs` | Re-import package docs to spatie.be after a docs PR merges |
| `write-freek-dev-blogpost` | Draft a post in the freek.dev voice |
| `typefully` | Draft and schedule social posts |
| `code-snippet-images` | Render code screenshots for social media |

**Tooling**

| Skill | What it does |
|-------|--------------|
| `conductor` | Conductor workspaces and sessions, the CLI, the API, deep links, settings, writing briefs |
| `ui` | Explore, build, and refine UI through the ui.sh MCP server |
| `grill-me` | Interrogate a plan until every branch of the decision tree is resolved |

### Scoped Plugins

Any directory under `skills/` with a `.claude-plugin/plugin.json` loads as a plugin named `<name>@skills-dir`. Plugins can be turned off globally and enabled per repository, so situational skills cost nothing in unrelated sessions.

| Plugin | Contents | Enabled in |
|--------|----------|------------|
| `marketing` | 29 marketing, CRO, and SEO skills | freek.dev, spatie.be, flareapp.io, mailcoach, ohdear, there-there.app |
| `music` | Ableton Live control | `~/dev/code/music` |
| `laravel-lsp` | The Laravel language server, see below | everywhere |

To enable one in a repository, add it to that repo's `.claude/settings.json`:

```json
{
  "enabledPlugins": {
    "marketing@skills-dir": true
  }
}
```

Project-scope plugins only load from the directory you launch from, so start the CLI at the repository root.

### Code Intelligence

Three language servers give the agent diagnostics after every edit and real symbol navigation instead of grep.

| Server | Provides | Install |
|--------|----------|---------|
| `laravel-lsp@skills-dir` | Config keys, route names, view paths, translation strings, middleware aliases, and container bindings, in `.blade.php` | `composer global require laravel/lsp` |
| `php-lsp` | PHP types, symbols, references, signatures | `npm i -g intelephense` |
| `typescript-lsp` | TypeScript and TSX intelligence | `npm i -g typescript-language-server typescript` |

Laravel LSP is first-party (announced at Laracon US 2026) and has no official Claude plugin, so `skills/laravel-lsp/` wraps it in a small `.lsp.json`. Its `phpEnvironment` defaults to `auto`, which finds Herd and Valet without configuration.

Claude Code registers **one server per file extension**, so the two PHP servers split the work: Intelephense takes `.php` for types, undefined methods, and references, and Laravel LSP takes `.blade.php` for route names, view paths, and translation strings. To flip that priority, disable `php-lsp` and add `".php"` back to `skills/laravel-lsp/.lsp.json`.

### Agents

Custom subagents live in `config/claude/agents/`.

- `laravel-feature-builder` - Implements new features across models, controllers, migrations, and views

Most delegation now happens through plugins instead, notably the `laravel-simplifier:laravel-simplifier` agent from Taylor's `laravel` marketplace, which is lane 2 of the review workflow.

### The Review Workflow

`review-code` and `review-pr` share one definition of what a review is, in `skills/review-code/references/lanes.md`. Six lanes run in parallel:

1. **Correctness** - bugs, via `/code-review`
2. **PHP simplification** - via the `laravel-simplifier` agent
3. **Spatie conventions** - via `spatie-guidelines`
4. **Laravel practices** - via Laravel Boost's per-repo `laravel-best-practices`
5. **React** - via `react-best-practices`, only when JS or TS changed
6. **Security** - authorization, mass assignment, injection, XSS, secrets, SSRF, PII in logs

Findings are deduplicated, and security outranks correctness, which outranks conventions. Both defect lanes require a concrete failure scenario, so "consider adding a null check" does not count as a finding.

The lanes file has a per-harness table, so the same review runs under Codex using its `review-agent` skill.

### Hooks

A `PostToolUse` hook runs Pint on every PHP file the agent edits. It walks up from the file to the nearest `vendor/bin/pint` and formats just that file, so repositories without Pint are a no-op. Formatting is deterministic rather than something the agent has to remember.

### Settings Worth Knowing

| Setting | Why |
|---------|-----|
| `disableClaudeAiConnectors` | Keeps claude.ai account connectors in the desktop app, out of the terminal |
| `skillOverrides` | Hides 8 bundled skills from the model while keeping them typeable, saving context |
| `cleanupPeriodDays` | Transcript retention, set to a year |
| `defaultMode: auto` | A classifier handles approvals instead of prompting on every step |

### Sharing With Codex

Codex reads `AGENTS.md` natively and Claude Code reads `CLAUDE.md`, so one file is symlinked under both names. Skills are a different story: Codex keeps its own `~/.codex/skills/` alongside its built-in `.system` skills, so a directory symlink would destroy those. `bin/link-agent-skills` links the harness-neutral skills one by one instead.

```bash
bin/link-agent-skills
```

Excluded from sharing: `ui` (needs a Claude MCP tool), `typefully` (Claude-specific `allowed-tools`), and the plugin directories, which have no top-level `SKILL.md`.

### Adding New Skills

```bash
# Install a skill (adds directly to your dotfiles)
npx skills add <owner/repo>

# Share it with Codex too, if it is harness-neutral
# (add it to the SHARED list in bin/link-agent-skills)
bin/link-agent-skills

cd ~/.dotfiles
git add config/claude/skills/
git commit -m "Add new skill"
```

Two rules that decide whether a skill ever gets used:

1. **The name and description determine everything.** A skill named after a person, or after a quarter of what it does, will not be found. Write the description around the phrases you would actually type.
2. **Keep the body short and push detail into `references/`.** A skill body stays in context for the rest of the session once loaded, while reference files load only when needed.

Browse more skills at [skills.sh](https://skills.sh)

---

## Customization

### Personal Aliases & Functions

Create custom configurations that won't be committed:

```bash
mkdir -p ~/.dotfiles-custom/shell
vim ~/.dotfiles-custom/shell/.aliases
```

These files are automatically loaded by `.zshrc` if they exist.

### Project-Specific Variables

Use `direnv` for automatic environment loading:

```bash
cd my-project
echo 'export DEBUG=true' > .envrc
direnv allow
```

Variables load when you enter the directory and unload when you leave.

---

## Post-Installation

1. **Restore settings** (optional): Run `mackup restore` if you have backups

2. **Migrate history** (upgrading only): Run `migration/migrate-z-to-zoxide.sh` if you have `~/.z`

---

## Tool Comparisons

| Old Tool | New Tool | Why Better |
|----------|----------|------------|
| z.sh / autojump | zoxide | Smarter frecency algorithm, Rust speed |
| nvm | fnm | 40x faster, simpler, Rust-based |
| cat | bat | Syntax highlighting, git integration |
| ls | eza | Icons, tree view, git status |
| grep | ripgrep | 5-10x faster, respects .gitignore |
| find | fd | Simpler syntax, 10x faster |
| diff | delta | Side-by-side diffs, syntax highlighting |
| htop | bottom | Better UI, graphs, Rust-based |

---

## Utilities

The `bin/` directory contains helper scripts:

- **install** - Main installation script (idempotent, safe to re-run)
- **install-claude-code** - Standalone installer for the AI setup: the CLI, the symlinks, and the Codex links
- **link-agent-skills** - Symlink the harness-neutral skills and `AGENTS.md` into Codex, leaving Codex's own built-in skills alone
- **exclude-from-spotlight** - Drop a `.metadata_never_index` marker into data heavy directories so Spotlight skips them. Local database directories (DBngin and friends) hold hundreds of thousands of constantly rewritten files, which keeps `mds_stores` busy indefinitely.
- **update** - Update dotfiles, Homebrew, npm, and Composer packages
- **doctor** - Health check and diagnostic tool
- **conductor-merge** - Fast-forward the current Conductor workspace branch into `main` (which lives in another git worktree). Use `--push` to also push `main` to `origin`, which clears Conductor's "Changes" view (it diffs against `origin/main`).

---

## Migration Notes

If upgrading from an older setup:

1. **Directory history**: Run `migration/migrate-z-to-zoxide.sh` to import your `~/.z` data
2. **Prompt**: The default is now Oh My Zsh with custom agnoster theme
3. **Version managers**:
   - fnm replaces nvm for Node.js
   - Homebrew manages PHP/Composer (no more compilation or mise)
4. **Terminal**: Ghostty replaces iTerm2 (config symlinked from dotfiles)
5. **Claude Code Skills**: Now version-controlled in `config/claude/skills/` and symlinked to `~/.claude/skills`
6. **Claude Code Agents**: Now version-controlled in `config/claude/agents/` and symlinked to `~/.claude/agents`
7. **Agent instructions**: `CLAUDE.md` was renamed to `AGENTS.md` and is symlinked under both names, so Codex reads the same file. `laravel-php-guidelines.md` was folded into the `spatie-guidelines` skill and removed
8. **Custom Theme**: Custom agnoster theme stored in `oh-my-zsh-custom/themes/`

---

## Credits

Created by [Freek Van der Herten](https://github.com/freekmurze). Used by many at [Spatie](https://spatie.be).

See `config/Brewfile` for complete package list.
