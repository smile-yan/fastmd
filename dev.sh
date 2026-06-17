#!/usr/bin/env bash
set -e
cd "$(dirname "$0")"

# Ensure the GOPATH bin (where `go install` puts `wails3`) is on PATH.
# Most shells don't include it by default.
if ! command -v wails3 >/dev/null 2>&1; then
    if command -v go >/dev/null 2>&1; then
        gopath_bin="$(go env GOPATH)/bin"
        if [ -x "$gopath_bin/wails3" ]; then
            export PATH="$gopath_bin:$PATH"
        fi
    fi
fi

if ! command -v wails3 >/dev/null 2>&1; then
    echo "Error: wails3 not found. Install it with:" >&2
    echo "  go install github.com/wailsapp/wails/v3/cmd/wails3@latest" >&2
    exit 1
fi

echo "→ 停止已有进程..."
pkill -f "fastmd" 2>/dev/null || true
pkill -f "wails3" 2>/dev/null || true
lsof -ti tcp:9245 | xargs kill -9 2>/dev/null || true
sleep 1

echo "→ 启动 wails3 dev..."
exec wails3 dev
