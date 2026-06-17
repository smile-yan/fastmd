#!/bin/sh
# Wrapper for the `fastmd` CLI entry point. This file is bundled into
# fastmd.app/Contents/Resources/fastmd by build/darwin/Taskfile.yml and
# exposed on $PATH by scripts/install-cli.sh.
#
# Why a wrapper rather than symlinking the Go binary directly: `open -a`
# routes the file through LaunchServices, which fires the kAEOpenDocuments
# Apple Event on the running instance (or cold-launches one). The Go side
# consumes that via events.Common.ApplicationOpenedWithFile in core/run.go
# and opens the file in a new window. Calling the Go binary directly would
# launch a fresh process every time and bypass that routing.
exec open -a fastmd "$@"
