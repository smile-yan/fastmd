# Top Border Double-Click Toggle Design

## Overview

Double-clicking the window's top resize border toggles between the maximum visible screen size and the original window size.

## Behavior

- Double-click top border → expand to `NSScreen.visibleFrame` (excludes menu bar and Dock)
- Double-click again → restore saved original frame
- Repeatable toggle

## Implementation

**Files changed:** `devtools_darwin.go`, `main.go`

### `devtools_darwin.go`

Add Objective-C function `setupTopBorderDoubleClick(void *nswin)`:

- Registers an `NSEvent` local monitor for `NSEventMaskLeftMouseDown`
- Guard with `static BOOL setup` to prevent duplicate registration
- On each event: skip if not this window, not a double-click, or not in top border zone
- Top border zone: mouse Y within `[windowTop - 5, windowTop + 2]` in screen coordinates
- Toggle: save frame → set `visibleFrame` (expand) or restore saved frame (collapse)
- Returns `nil` to consume the event

State: `static NSRect _savedFrame` and `static BOOL _isExpanded`.

Add Go wrapper: `func setupTopBorderDoubleClick(nsWindow unsafe.Pointer)`

### `main.go`

In the existing `WindowShow` hook, add one call to `setupTopBorderDoubleClick(window.NativeWindow())`.
