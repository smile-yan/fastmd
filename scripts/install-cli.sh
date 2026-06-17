#!/bin/sh
# Install the `fastmd` CLI command for the current user.
#
# Creates a symlink in a writable $PATH directory that points at the
# wrapper script bundled inside fastmd.app. The wrapper itself is a
# tiny shim that runs `open -a fastmd "$@"`, which routes through
# macOS LaunchServices so kAEOpenDocuments hits the running instance
# (or cold-launches one). Re-running this script is safe: existing
# symlinks are replaced.
#
# Usage:   bash scripts/install-cli.sh
# Env:     APP_PATH   Override the .app location (default: /Applications/fastmd.app)

set -eu

APP_NAME="fastmd"
WRAPPER_NAME="fastmd"
APP_PATH="${APP_PATH:-/Applications/$APP_NAME.app}"
WRAPPER_REL="Contents/Resources/$WRAPPER_NAME"

log() { printf 'install-cli: %s\n' "$*"; }
err() { printf 'install-cli: error: %s\n' "$*" >&2; }

# 1. Validate the .app and its bundled wrapper.
if [ ! -d "$APP_PATH" ]; then
  err "$APP_PATH not found."
  err "Drag fastmd.app to /Applications first, or set APP_PATH to its location."
  exit 1
fi

WRAPPER="$APP_PATH/$WRAPPER_REL"
if [ ! -f "$WRAPPER" ]; then
  err "wrapper not found at $WRAPPER."
  err "The bundled .app is missing the CLI shim — rebuild with 'task build' (or reinstall fastmd)."
  exit 1
fi

# 2. Pick a writable PATH directory. Priority order:
#    - /opt/homebrew/bin (Apple Silicon Homebrew default)
#    - /usr/local/bin     (Intel Homebrew / custom prefix)
#    - $HOME/.local/bin  (no-sudo fallback)
#    No-sudo is preferred; we only suggest sudo as a manual override.
CANDIDATE_DIRS="/opt/homebrew/bin /usr/local/bin $HOME/.local/bin"

chosen=""
for dir in $CANDIDATE_DIRS; do
  case "$dir" in
    */.local/bin) mkdir -p "$dir" ;;
  esac
  # Check writability without actually creating anything yet: try to
  # create and remove a probe file. This is more accurate than -w on
  # some macOS configurations where the directory has the sticky bit.
  probe="$dir/.install-cli-probe-$$"
  if (umask 077 && : > "$probe") 2>/dev/null; then
    rm -f "$probe"
    chosen="$dir"
    break
  fi
done

if [ -z "$chosen" ]; then
  err "no writable PATH directory found (tried: $CANDIDATE_DIRS)."
  err "Run with sudo to install into /usr/local/bin, or create ~/.local/bin and add it to PATH."
  exit 1
fi

link="$chosen/$WRAPPER_NAME"

# 3. Refuse to clobber an unrelated file/symlink with the same name.
if [ -e "$link" ] && [ ! -L "$link" ]; then
  err "$link already exists and is not a symlink."
  err "Move or remove it, then re-run this script."
  exit 1
fi

# 4. Create the symlink (or refresh it).
ln -sf "$WRAPPER" "$link"
log "linked $link -> $WRAPPER"

# 5. Verify resolution.
if resolved="$(command -v "$WRAPPER_NAME" 2>/dev/null)" && [ "$resolved" = "$link" ]; then
  log "$WRAPPER_NAME is now on PATH at $resolved"
else
  log "$link was created, but $WRAPPER_NAME is not on PATH."
  log "Add this to your shell profile and restart the shell:"
  log "  export PATH=\"$chosen:\$PATH\""
fi

# 6. Final hint if LaunchServices hasn't seen the bundle yet.
if [ ! -d "$HOME/Library/Application Support/$APP_NAME" ]; then
  log "first-time install: launch fastmd once (open the .app, then quit) so"
  log "LaunchServices can route the file to it via Apple Events."
else
  tilde='~'
  log "open a terminal and try:  $WRAPPER_NAME $tilde/notes/hello.md"
fi
