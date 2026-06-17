#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
APP_NAME="${APP_NAME:-fastmd}"
BIN_DIR="${BIN_DIR:-bin}"
APP_BUNDLE="$BIN_DIR/$APP_NAME.app"
DMG_NAME="${DMG_NAME:-$APP_NAME.dmg}"
DMG_TMP="$BIN_DIR/$APP_NAME-tmp.dmg"
DMG_OUT="$BIN_DIR/$DMG_NAME"
VOLUME_NAME="${VOLUME_NAME:-fastmd-build}"
ICON_FILE="${ICON_FILE:-build/darwin/icons.icns}"
UNIVERSAL="${UNIVERSAL:-false}"

cd "$ROOT_DIR"

if [ "$(uname -s)" != "Darwin" ]; then
  echo "Error: macOS DMG packaging must be run on macOS." >&2
  exit 1
fi

for cmd in wails3 hdiutil SetFile sips DeRez Rez; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "Error: $cmd is required but was not found in PATH." >&2
    echo "Install Xcode Command Line Tools if the missing command is a macOS developer tool." >&2
    exit 1
  fi
done

LSREGISTER="/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister"

echo "Packaging macOS app..."
if [ "$UNIVERSAL" = "true" ] || [ "$UNIVERSAL" = "1" ]; then
  wails3 task darwin:package:universal
else
  wails3 task darwin:package
fi

"$LSREGISTER" -u "$APP_BUNDLE" 2>/dev/null || true

echo "Creating DMG..."
rm -f "$DMG_TMP" "$DMG_OUT"

hdiutil create \
  -volname "$VOLUME_NAME" \
  -srcfolder "$APP_BUNDLE" \
  -ov -format UDRW \
  "$DMG_TMP"

MOUNT_DIR="/Volumes/$VOLUME_NAME"
ATTACH_OUT="$(hdiutil attach "$DMG_TMP" -mountpoint "$MOUNT_DIR" -nobrowse)"
DISK_DEV="$(printf '%s\n' "$ATTACH_OUT" | awk '/GUID_partition_scheme/ {print $1; exit}')"

ln -sf /Applications "$MOUNT_DIR/Applications"
cp "$ICON_FILE" "$MOUNT_DIR/.VolumeIcon.icns"
SetFile -a C "$MOUNT_DIR"
"$LSREGISTER" -u "$MOUNT_DIR/$APP_NAME.app" 2>/dev/null || true
hdiutil detach "${DISK_DEV:-$MOUNT_DIR}" -force

hdiutil convert "$DMG_TMP" -format ULFO -o "$DMG_OUT"
rm -f "$DMG_TMP"

ICON_TMP="$BIN_DIR/$APP_NAME-dmg-icon.icns"
ICON_RSRC="$BIN_DIR/$APP_NAME-dmg-icon.rsrc"
cp "$ICON_FILE" "$ICON_TMP"
sips -i "$ICON_TMP" >/dev/null
DeRez -only icns "$ICON_TMP" > "$ICON_RSRC"
Rez -append "$ICON_RSRC" -o "$DMG_OUT"
rm -f "$ICON_TMP" "$ICON_RSRC"
SetFile -a C "$DMG_OUT"

echo "Done: $DMG_OUT"
