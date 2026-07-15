# Path to your oh-my-zsh configuration.
ZSH=$HOME/.oh-my-zsh

# Skip oh-my-zsh's compaudit insecure-directory scan on startup (~20ms). Safe on
# a single-user machine where we control everything in fpath.
ZSH_DISABLE_COMPFIX="true"

# Path to custom themes and plugins
ZSH_CUSTOM=$HOME/.dotfiles/oh-my-zsh-custom

# Set name of the theme to load.
# Look in ~/.oh-my-zsh/themes/
# Optionally, if you set this to "random", it'll load a random theme each
# time that oh-my-zsh is loaded.
# Conductor's built-in terminal can't render agnoster's powerline/Nerd Font
# glyphs, so use a plain prompt there (set up after oh-my-zsh loads, below).
# Every other terminal (Ghostty, iTerm, ...) keeps agnoster.
if [[ "$__CFBundleIdentifier" == "com.conductor.app" || -n "$CONDUCTOR_INTERNAL_BIN_DIR" ]]; then
    CONDUCTOR_TERMINAL=1
    ZSH_THEME=""
else
    ZSH_THEME="agnoster"
    AGNOSTER_DIR_FG=black
fi

# Hide username in prompt
DEFAULT_USER=`whoami`

# Which plugins would you like to load? (plugins can be found in ~/.oh-my-zsh/plugins/*)
# Custom plugins may be added to ~/.oh-my-zsh/custom/plugins/
# Example format: plugins=(rails git textmate ruby lighthouse)
plugins=(git composer macos)

# Add Homebrew's completions to fpath *before* oh-my-zsh runs compinit, so its
# single completion init picks them up (avoids a second, expensive compinit).
fpath=(/opt/homebrew/share/zsh/site-functions $fpath)

source $ZSH/oh-my-zsh.sh

# Plain prompt for Conductor's terminal: no powerline separators or Nerd Font
# glyphs (those render as tofu boxes there). Uses only glyphs present in the
# default font. The git plugin provides git_prompt_info.
if [[ -n "$CONDUCTOR_TERMINAL" ]]; then
    setopt prompt_subst
    # Mimic agnoster's segmented look without Nerd Font/powerline glyphs
    # (Conductor's terminal renders those as tofu boxes). Solid color blocks
    # stand in for the slanted powerline separators: a blue directory segment
    # and a green git segment, matching the agnoster theme used in Ghostty.
    conductor_git_segment() {
        local branch
        branch=$(command git symbolic-ref --short HEAD 2>/dev/null) || return
        local dirty=""
        [[ -n "$(command git status --porcelain 2>/dev/null)" ]] && dirty=" *"
        print -n "%K{green}%F{black} ${branch}${dirty} %k"
    }
    PROMPT='%K{blue}%F{white} %~ %k$(conductor_git_segment)%f '
fi

# Removed old RVM path
#set numeric keys
# 0 . Enter
bindkey -s "^[Op" "0"
bindkey -s "^[Ol" "."
bindkey -s "^[OM" "^M"
# 1 2 3
bindkey -s "^[Oq" "1"
bindkey -s "^[Or" "2"
bindkey -s "^[Os" "3"
# 4 5 6
bindkey -s "^[Ot" "4"
bindkey -s "^[Ou" "5"
bindkey -s "^[Ov" "6"
# 7 8 9
bindkey -s "^[Ow" "7"
bindkey -s "^[Ox" "8"
bindkey -s "^[Oy" "9"
# + -  * /
bindkey -s "^[Ok" "+"
bindkey -s "^[Om" "-"
bindkey -s "^[Oj" "*"
bindkey -s "^[Oo" "/"

# Load the shell dotfiles, and then some:
# * ~/.dotfiles-custom can be used for other settings you don’t want to commit.
for file in ~/.dotfiles/home/.{exports,aliases,functions}; do
	[ -r "$file" ] && [ -f "$file" ] && source "$file"
done

for file in ~/.dotfiles-custom/shell/.{exports,aliases,functions,zshrc}; do
	[ -r "$file" ] && [ -f "$file" ] && source "$file"
done
unset file

# Directory jumping now handled by zoxide (see modern tools section below)


# Sudoless npm https://github.com/sindresorhus/guides/blob/master/npm-global-without-sudo.md
NPM_PACKAGES="${HOME}/.npm-packages"
export PATH="$PATH:$NPM_PACKAGES/bin"
# Preserve MANPATH if you already defined it somewhere in your config.
# Otherwise, fall back to `manpath` so we can inherit from `/etc/manpath`.
export MANPATH="${MANPATH-$(manpath)}:$NPM_PACKAGES/share/man"

export PATH=$HOME/.dotfiles/bin:$PATH

# Import ssh keys in keychain
ssh-add --apple-use-keychain 2>/dev/null;

# Setup xdebug
export XDEBUG_CONFIG="idekey=VSCODE"

# Enable autosuggestions (installed via brew)
source /opt/homebrew/share/zsh-autosuggestions/zsh-autosuggestions.zsh


# Extra paths
export PATH="$HOME/.composer/vendor/bin:$PATH"
export PATH=/usr/local/bin:$PATH
export PATH="$HOME/.yarn/bin:$PATH"

# Use the Homebrew-path valet so `valet trust` (which whitelists /opt/homebrew/bin/valet)
# lets secure/link/open etc. run without a sudo password prompt.
alias valet="/opt/homebrew/bin/valet"


#export PATH=/Users/Shared/DBngin/postgresql/17.0/bin:$PATH

export PATH=$HOME/bin:~/.config/phpmon/bin:$PATH
export JAVA_HOME="/opt/homebrew/opt/openjdk@17/libexec/openjdk.jdk/Contents/Home"
export ANDROID_HOME="$HOME/Library/Android/sdk"
export PATH="$ANDROID_HOME/platform-tools:$ANDROID_HOME/emulator:$PATH"
export PATH="$HOME/.local/bin:$PATH"

# Initialize modern tools
# zoxide - smarter cd
if command -v zoxide &> /dev/null; then
    eval "$(zoxide init zsh)"
fi

# fnm - Node.js version manager
if command -v fnm &> /dev/null; then
    eval "$(fnm env --use-on-cd)"
fi

# bun completions
[ -s "/Users/freek/.bun/_bun" ] && source "/Users/freek/.bun/_bun"

# bun
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$PATH"

# Raise the open-file limit. macOS defaults NOFILE to "unlimited", which the
# kernel clamps to a tiny legacy value for real open() calls, causing
# "Too many open files" in test runs. A concrete number avoids that clamp.
ulimit -n 65536
