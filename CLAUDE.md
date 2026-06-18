# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

fast-md is a Typora-like Markdown editor for macOS built on Wails 3 (Go) + Vue 3 + Milkdown (Crepe). The Go process embeds the Vite-built frontend via `//go:embed all:frontend/dist` and exposes an `AppService` to JS through Wails-generated bindings.

**Module name is `changeme`** — the placeholder from the Wails template. It still appears in import paths and in `frontend/bindings/changeme/...`. Don't "fix" it; the bindings generator and tooling depend on it.

## Rules for Claude

- **DO NOT add `Co-Authored-By: Claude <noreply@anthropic.com>` or any similar attribution to commits.** The user explicitly forbids Claude from appearing as a contributor in this repository's git history. Commit messages should only contain the user's own content.
- **DO NOT sign commits or tags on behalf of Claude.** Use the user's configured git identity only.
- **Use the repository's default git user and email for all commits and tags.** Do not use any other identity.
- **When creating GitHub releases, publish as official releases by default.** Only use `--draft` flag when the user explicitly specifies they want a draft release. All other releases should be set as `latest`.

## Build / Run / Test

The root `Taskfile.yml` includes per-platform subtasks from `build/`. Top-level entry points:

```bash
task dev           # wails3 dev on port 9245 (frontend watch + Go rebuild)
task build         # delegates to {OS}:build
task package       # delegates to {OS}:package
./dev.sh           # kills any stale fast-md/wails3 on :9245, then runs `wails3 dev`
```

Frontend-only (run from `frontend/`):

```bash
npm ci --include=dev
npm run build      # vue-tsc + vite build (production)
npm run build:dev  # vue-tsc + vite build (development, unminified)
npm test           # vitest run (jsdom env, see vite.config.ts)
```

Go (run from repo root; matches CI):

```bash
go test -v -race ./...                                   # what CI runs
GOCACHE=/private/tmp/fast-md-go-cache go test ./...      # if sandboxed
wails3 generate bindings -f '-tags production' -clean=true -ts   # regenerate frontend/bindings/
```

Always run `npm run build` before `go test`/`go build` when `frontend/dist` changes — Go embeds that directory at compile time and tests fail if `frontend/dist/index.html` is missing.

## High-Level Architecture

### Go side (`*.go` at repo root, single `package main`)

- **`main.go`** — Wails app bootstrap, window factory (`newEditorWindow`, `newEditorWindowWithFile`), menu wiring, devtools key binding, file-drop routing. Embeds `frontend/dist` via `//go:embed all:frontend/dist` and reads `frontend/dist/themes/index.json` to populate the View → Theme submenu.
- **`app.go`** — `AppService` (the bound type exposed to JS). Methods: `OpenFileDialog`, `SaveFileDialog`, `OpenFolderDialog`, `ReadFile`, `WriteFile`, `ListDirectory`, `CloseWindow`, `RequestQuit`/`ConfirmQuitWindow`/`CancelQuit`, `SetUILocale`/`GetUILocale`, `GetConfig`/`SaveConfig`, `ShowSaveDialog` (unsaved-changes prompt), `ExportPDF` (uses `github.com/carlos7ags/folio` for HTML→PDF). Also owns `loadConfig`/`saveConfig` writing to `~/Library/Application Support/fast-md/config.json`.
- **`quit.go`** — `quitCoordinator`: gates window-close requests, queues windows, prompts each in turn via the `app:confirmQuitWindow` event, and only calls `app.Quit()` after every window confirms. The `allowedClose` map (in `main.go`) is the single source of truth for "may this window actually close now?".
- **`menu_i18n.go`** — `Locale` (`zh` | `en`), `menuStrings` table, `SetLocale`, `buildMenuI18n`. `zh` is the default; `SetUILocale` from the frontend rebuilds the menu in place.
- **`menu_icons.go`** — SF-symbol-style icons for every menu item via `setMenuIcon`.
- **`dockmenu_darwin.go`** + **`dockmenu_impl_darwin.m`** (CGO) — native dock menu with recent-files list (last 10 `.md`/`.markdown` paths, tracked in `trackRecentFile`). `_darwin.go`/`_other.go` build-tag split keeps non-darwin builds compiling.
- **`savedialog_darwin.go`** + **`savedialog_impl_darwin.m`** — wraps NSSavePanel for PDF export (uses UTF-8 allowed set).
- **`devtools_darwin.go`** / **`devtools_other.go`** — F12 devtools handler.
- **`help_docs.go`** — builds the Help menu from Markdown files in `docs/help/` (4 entries: quick-start, keyboard-shortcuts, markdown-basics, math-formula-basics).

### Frontend (`frontend/src/`)

- **`App.vue`** — top-level layout: `Sidebar` + `Editor` (Milkdown Crepe) + `StatusBar` + `Settings` modal. Owns dirty-state, source-mode toggle, and routes Go events (`menu:newFile`, `menu:open`, `menu:save`, `menu:saveAs`, `menu:exportHTML`, `menu:exportPDF`, `menu:toggleSidebar`, `menu:setTheme`, `menu:settings`, `file:open`, `window:closeRequested`, `app:aboutRequested`).
- **`components/`** — `Editor.vue` (Milkdown Crepe + a long list of `wrapIn*Command`/`insert*Command` imports for keyboard shortcuts), `Sidebar.vue` (folder tree + recent files), `StatusBar.vue`, `Settings.vue` (theme, language, autosave).
- **`composables/`** — `useFile` (file state + autosave loop reading `localStorage['fast-md-settings']`), `useLocale` (zh/en strings + `t()`), `useTheme` (light/dark), `useContentTheme` (loads CSS from `frontend/public/themes/<name>.css`), `useEditorSettings`.
- **`exportHtml.ts`** — standalone HTML export (self-contained document with inlined CSS), tested in `exportHtml.test.ts`.
- **`bindings/changeme/appservice.ts`** — Wails-generated TS bindings; the import path `../../bindings/changeme/appservice` is what the composables/components use. Regenerate with `wails3 generate bindings ...` whenever a Go method on `AppService` changes.
- **`public/themes/`** — content theme CSS files + `index.json` manifest. Adding a theme = add a CSS file + append an entry to `index.json` (Go's `loadThemes` reads this on startup to populate the menu).

### IPC pattern

Go → JS: `app.Event.Emit("app:xxx")` or `window.EmitEvent("menu:xxx", payload)`; JS listens via `Events` from `@wailsio/runtime`. The Go menu and the Vue handlers share the same event-name constants — see `App.vue`'s event listeners and the `emitToFocused(...)` calls in `main.go`.

JS → Go: imports from `frontend/bindings/changeme/appservice` (regenerated by `wails3 generate bindings`).

### Window lifecycle / "close vs quit"

Closing one window does not quit. Closing the last window keeps the app alive (`ApplicationShouldTerminateAfterLastWindowClosed: false`) so the dock menu still works. The real quit is gated by `quitCoordinator`:
1. User triggers `RequestQuit` → `Begin(queue)` with all open windows.
2. For each window: `Focus()` + emit `app:confirmQuitWindow` (frontend shows its unsaved-changes UI).
3. Frontend replies via `ConfirmQuitWindow` (close-and-save) or `CancelQuit` (abort whole quit). `Confirm` adds the window ID to `allowedClose`, then `Close()`; the `events.Common.WindowClosing` hook sees the flag and lets it through.
4. When the queue empties, `app.Quit()` runs.

Single-window close: frontend gets `window:closeRequested`, prompts via `ShowSaveDialog`, then calls `CloseWindow()` which marks `allowedClose` and closes.

## Conventions / Gotchas

- **Module name `changeme`** is intentional. Do not rename without regenerating every binding path.
- **Frontend must be built before `go test`/`go build`** if any frontend file changed — `main.go` embeds `frontend/dist`.
- **macOS-only features** (dock menu, NSSavePanel, devtools handler) live in `*_darwin.go` files paired with `*_other.go` stubs so non-darwin builds still compile. Use the same split for any new platform-specific feature.
- **Wails 3 is alpha (`v3.0.0-alpha.96`)** — APIs and CLI flags may differ from Wails v2; check `wails3 --help` rather than guessing from Wails v2 docs.
- **CI matrix** (`.github/workflows/ci.yml` + `release.yml`): Go 1.25, Node 24, ubuntu-latest for tests. Linux deps: `libgtk-4-dev libwebkitgtk-6.0-dev libsoup-3.0-dev`. Releases are tagged `v*`, cross-built for `linux/amd64`, `windows/amd64`, `darwin/{arm64,amd64,universal}`, packaged as DMG via `scripts/build-macos.sh` + `hdiutil`.
- **Vite port is fixed** (`strictPort: true`, default `9245`) — `dev.sh` kills anything bound there before starting.
- **TypeScript is `^4.9.5`** while `vue-tsc` is `^1.8.27` and `vitest` is `^4.x` — these versions are pinned for a reason; bumping may break the milkdown type-check. CI runs `node node_modules/vue-tsc/bin/vue-tsc.js` directly (not `npx`) to avoid version drift.
- **`.env` does not exist in this repo**; the workspace-level `gails/.env` holds a `GITHUB_TOKEN` for that other project — not related to fast-md, do not read or echo it here.
- **The `gails/` and `gails-web/` directories at the workspace root** are separate projects (a Wails fork and a marketing page). They are not dependencies and not imported by fast-md.