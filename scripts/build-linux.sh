#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
APP_NAME="${APP_NAME:-fastmd}"
ARCH="${ARCH:-amd64}"
FORMAT="${FORMAT:-all}"

cd "$ROOT_DIR"

if ! command -v wails3 >/dev/null 2>&1; then
  echo "Error: wails3 is required but was not found in PATH." >&2
  exit 1
fi

if [ "$(uname -s)" != "Linux" ]; then
  if ! command -v docker >/dev/null 2>&1; then
    echo "Error: Docker is required to package Linux from a non-Linux host." >&2
    exit 1
  fi
  if ! docker image inspect wails-cross >/dev/null 2>&1; then
    echo "Error: Docker image 'wails-cross' was not found." >&2
    echo "Build it first with: wails3 task setup:docker" >&2
    exit 1
  fi
fi

echo "Packaging Linux app..."
echo "  ARCH=$ARCH"
echo "  FORMAT=$FORMAT"

case "$FORMAT" in
  all)
    GOOS=linux GOARCH="$ARCH" CGO_ENABLED=1 wails3 task linux:package ARCH="$ARCH"
    ;;
  binary)
    GOOS=linux GOARCH="$ARCH" CGO_ENABLED=1 wails3 task linux:build ARCH="$ARCH"
    ;;
  appimage)
    GOOS=linux GOARCH="$ARCH" CGO_ENABLED=1 wails3 task linux:create:appimage ARCH="$ARCH"
    ;;
  deb)
    GOOS=linux GOARCH="$ARCH" CGO_ENABLED=1 wails3 task linux:create:deb ARCH="$ARCH"
    ;;
  rpm)
    GOOS=linux GOARCH="$ARCH" CGO_ENABLED=1 wails3 task linux:create:rpm ARCH="$ARCH"
    ;;
  aur)
    GOOS=linux GOARCH="$ARCH" CGO_ENABLED=1 wails3 task linux:create:aur ARCH="$ARCH"
    ;;
  *)
    echo "Error: unknown FORMAT '$FORMAT'." >&2
    echo "Supported values: all, binary, appimage, deb, rpm, aur" >&2
    exit 1
    ;;
esac

echo "Done: Linux package output is under bin/"
