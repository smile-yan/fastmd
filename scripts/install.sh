#!/bin/bash
# Script to build and install fastmd on macOS

set -e

cd "$(dirname "$0")/.."

echo "Building fastmd..."
task package

echo "Removing old installation..."
rm -rf /Applications/fastmd.app

echo "Installing to /Applications..."
cp -R bin/fastmd.app /Applications/

# Rebuild Launch Services to fix duplicate UTI entries in "Open With" menu
echo "Rebuilding Launch Services database..."
LSREGISTER="/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister"
"$LSREGISTER" -kill -r -domain local -domain system -domain user 2>/dev/null || true

echo "Done! fastmd has been installed to /Applications/fastmd.app"