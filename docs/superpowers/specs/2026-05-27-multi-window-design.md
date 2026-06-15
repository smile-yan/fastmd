# Multi-Window Support Design

## Goal

Support multiple independent editor windows. Each window has its own file state, sidebar state, and lifecycle. Cmd+N opens a new blank window. Closing the last window keeps the app running in the background.

## Architecture

### Go side (main.go, app.go)

**`newEditorWindow(app *application.App)`** — extracts the current window creation block from `main()` into a standalone function. Sets up all hooks (WindowClosing, WindowShow, WindowFilesDropped) on the new window. Returns `*application.WebviewWindow`.

**Menu event routing** — all `app.Event.Emit(...)` calls in `buildMenu` that target a single window (save, open, new file, sidebar toggle, theme, export, fullscreen, devtools) are replaced with:
```go
w := app.Window.GetFocused()
if w != nil { w.EmitEvent(...) }
```
`app:quitRequested` remains a broadcast since it targets the whole app.

**Cmd+N** — the "New" menu item calls `newEditorWindow(app)` directly in Go instead of emitting `menu:new` to the frontend.

**`AppService`** — remove the `window *application.WebviewWindow` field. All methods that call `AttachToWindow(s.window)` are changed to `AttachToWindow(s.app.Window.GetFocused())`. `HideWindow` is replaced by `CloseWindow()` which calls `s.app.Window.GetFocused().Close()` — in multi-window mode, closing a window destroys it rather than hiding it. `ShowWindow` is removed.

**Dock reopen** — `ApplicationShouldHandleReopen` handler calls `newEditorWindow(app)` instead of `window.Show() / window.Focus()`.

**Window close hook** — each window's `WindowClosing` hook emits `app:closeRequested` via `window.EmitEvent(...)` (already per-window; no change needed in logic, just ensure the hook is set up inside `newEditorWindow`).

**Frontend close flow** — `App.vue` currently calls `HideWindow()` after confirming close. This is replaced by `CloseWindow()` (new binding). `executeClose('hide')` → `CloseWindow()`. `executeClose('quit')` → `QuitApp()` (unchanged).

### Frontend (no changes except close binding)

Each `WebviewWindow` runs an independent WebView process with its own JS heap. Module-level refs in `useFile.ts` (`filePath`, `content`, `isDirty`) are naturally isolated per window. `useContentTheme.ts` reads/writes `localStorage` — shared across windows, which is the desired behavior (theme setting is global).

`App.vue` changes: replace `HideWindow` import with `CloseWindow`, and update `executeClose('hide')` to call `CloseWindow()` instead.

## Behavior

| Action | Result |
|--------|--------|
| Cmd+N | New blank editor window opens |
| Cmd+S in window A | Saves window A's file only |
| Close window (red button) | That window's unsaved-changes dialog runs; other windows unaffected |
| Close last window | App stays in background (Dock icon remains) |
| Click Dock icon when no windows open | New blank window opens |
| Drag .md file onto any window | That window opens the file |

## Out of scope

- Tab support (planned as a separate phase)
- Window state restoration on app restart
- Inter-window communication (e.g. reload if file changed by another window)
