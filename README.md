# My dotfiles

![Terminal](images/terminal.png)

Personal dotfiles with modern shell tooling, optimized for Laravel/PHP development. Features fast startup times, smart directory navigation, and modern CLI tools.

## Key Features

- **Custom Agnoster Theme** - Clean powerline prompt with no branch symbols, `•` for changes
- **Version-Controlled Skills & Agents** - All Claude Code skills and agents synced via dotfiles
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
| `~/.claude/CLAUDE.md` | `~/.dotfiles/config/claude/CLAUDE.md` | Claude Code configuration |
| `~/.claude/laravel-php-guidelines.md` | `~/.dotfiles/config/claude/laravel-php-guidelines.md` | Laravel coding standards |
| `~/.claude/settings.json` | `~/.dotfiles/config/claude/settings.json` | Claude Code settings |
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
- **Global npm**: agent-browser
- **Global Composer**: laravel/envoy, spatie/phpunit-watcher, laravel/valet

---

## Claude Code Integration

### Quick Install (Standalone)

Install just Claude Code without the full dotfiles:

```bash
curl -fsSL https://raw.githubusercontent.com/freekmurze/dotfiles/main/bin/install-claude-code | bash
```

### What's Included

- **Claude Code CLI** - Installed via Homebrew
- **Custom configuration** - CLAUDE.md with coding guidelines, laravel-php-guidelines.md
- **Version-controlled skills** - Entire `~/.claude/skills` directory symlinked to dotfiles
- **Version-controlled agents** - Entire `~/.claude/agents` directory symlinked to dotfiles

### Skills (Version Controlled)

All skills are stored in `config/claude/skills/` and version-controlled with your dotfiles. When you run the installer on a new Mac, all skills are immediately available.

**Custom Skills:**
- `ray-skill` - Ray debugging integration
- `fix-github-issue` - GitHub issue automation
- `convert-issue-to-discussion` - GitHub workflow helpers

**Community Skills:**
- `vercel-labs/agent-skills` - Web design guidelines and React best practices
- `anthropics/skills` - Frontend design and skill creation tools
- `vercel-labs/agent-browser` - Browser automation
- `expo/skills` - React Native with Expo
- `callstackincubator/agent-skills` - React Native performance
- `coreyhaines31/marketingskills` - Copywriting and programmatic SEO
- `copy-editing` - Marketing copy editing
- `copywriting` - Marketing copywriting
- `frontend-design` - Frontend design patterns
- `pdf` - PDF manipulation
- `seo-audit` - SEO auditing
- `web-design-guidelines` - Web design best practices

### Adding New Skills

```bash
# Install a new skill (adds directly to your dotfiles)
npx skills add <owner/repo>

# Commit to version control
cd ~/.dotfiles
git add config/claude/skills/
git commit -m "Add new skill"
git push
```

Browse more skills at [skills.sh](https://skills.sh)

### Agents (Version Controlled)

All custom agents are stored in `config/claude/agents/` and version-controlled with your dotfiles. When you run the installer on a new Mac, all agents are immediately available.

**Custom Agents:**
- `laravel-simplifier` - Simplifies and refines PHP/Laravel code for clarity and maintainability
- `laravel-debugger` - Diagnoses and fixes issues in Laravel applications
- `laravel-feature-builder` - Implements new features in Laravel applications
- `task-planner` - Breaks down complex tasks into actionable steps

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
- **install-claude-code** - Standalone Claude Code installer
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
7. **Custom Theme**: Custom agnoster theme stored in `oh-my-zsh-custom/themes/`

---

## Credits

Created by [Freek Van der Herten](https://github.com/freekmurze). Used by many at [Spatie](https://spatie.be).

See `config/Brewfile` for complete package list.
