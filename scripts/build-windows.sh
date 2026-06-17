#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
APP_NAME="${APP_NAME:-fastmd}"
ARCH="${ARCH:-amd64}"
FORMAT="${FORMAT:-nsis}"
INSTALL_SCOPE="${INSTALL_SCOPE:-user}"
CGO_ENABLED="${CGO_ENABLED:-0}"

cd "$ROOT_DIR"

if ! command -v wails3 >/dev/null 2>&1; then
  echo "Error: wails3 is required but was not found in PATH." >&2
  exit 1
fi

if [ "$FORMAT" = "nsis" ] && ! command -v makensis >/dev/null 2>&1; then
  echo "Error: NSIS is required for Windows installer packaging." >&2
  echo "Install it on macOS with: brew install nsis" >&2
  exit 1
fi

echo "Packaging Windows app..."
echo "  ARCH=$ARCH"
echo "  FORMAT=$FORMAT"
echo "  INSTALL_SCOPE=$INSTALL_SCOPE"
echo "  CGO_ENABLED=$CGO_ENABLED"

case "$FORMAT" in
  binary)
    GOOS=windows GOARCH="$ARCH" CGO_ENABLED="$CGO_ENABLED" wails3 task windows:build ARCH="$ARCH"
    echo "Done: bin/$APP_NAME.exe"
    ;;
  nsis)
    GOOS=windows GOARCH="$ARCH" CGO_ENABLED="$CGO_ENABLED" wails3 task windows:package \
      ARCH="$ARCH" \
      FORMAT="$FORMAT" \
      INSTALL_SCOPE="$INSTALL_SCOPE"
    echo "Done: bin/$APP_NAME-$ARCH-installer.exe"
    ;;
  msix)
    GOOS=windows GOARCH="$ARCH" CGO_ENABLED="$CGO_ENABLED" wails3 task windows:package \
      ARCH="$ARCH" \
      FORMAT="$FORMAT" \
      INSTALL_SCOPE="$INSTALL_SCOPE"
    echo "Done: bin/$APP_NAME-$ARCH.msix"
    ;;
  *)
    echo "Error: unknown FORMAT '$FORMAT'." >&2
    echo "Supported values: binary, nsis, msix" >&2
    exit 1
    ;;
esac
