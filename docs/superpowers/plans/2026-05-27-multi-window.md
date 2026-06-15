# Multi-Window Support Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Support multiple independent editor windows — Cmd+N opens a new blank window, each window has isolated file state, closing the last window keeps the app running in the background.

**Architecture:** Each Wails `WebviewWindow` runs an independent WebView process, so JS module-level state in `useFile.ts` is naturally isolated per window. The main work is in Go: extract a `newEditorWindow` factory, route menu events to the focused window instead of broadcasting, and replace `HideWindow` with `CloseWindow` in the service layer.

**Tech Stack:** Go (Wails v3 alpha.96), Vue 3 + TypeScript, CGO/Objective-C (macOS)

---

## File Map

| File | Change |
|------|--------|
| `main.go` | Extract `newEditorWindow`, route menu events to focused window, update Dock reopen handler |
| `app.go` | Remove `window` field, add `CloseWindow()`, remove `ShowWindow()`, fix `ShowCloseSheet` to use focused window |
| `savedialog_darwin.go` | Fix `ShowCloseSheet` to get native window from focused window |
| `frontend/bindings/changeme/appservice.ts` | Regenerate (adds `CloseWindow`, removes `HideWindow`/`ShowWindow`) |
| `frontend/src/App.vue` | Replace `HideWindow` import with `CloseWindow`, update `executeClose` |

---

### Task 1: Add `CloseWindow`, remove `HideWindow`/`ShowWindow` from `app.go`

**Files:**
- Modify: `app.go`

**Context:** `AppService` currently holds a `window *application.WebviewWindow` field used by dialogs and `HideWindow`/`ShowWindow`. In multi-window mode, dialogs should attach to whichever window is currently focused. `HideWindow` is replaced by `CloseWindow` which closes (destroys) the focused window.

- [ ] **Step 1: Read `app.go`**

Read `/Users/yanshili/Downloads/md-p-1/fast-md/app.go` to confirm current state.

- [ ] **Step 2: Rewrite `AppService` struct and affected methods**

Replace the entire `app.go` content with:

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type FileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"isDir"`
}

type AppService struct {
	app *application.App
}

func (s *AppService) focusedWindow() *application.WebviewWindow {
	return s.app.Window.GetFocused()
}

func (s *AppService) OpenFileDialog() (string, error) {
	return s.app.Dialog.OpenFile().
		AttachToWindow(s.focusedWindow()).
		SetTitle("Open Markdown File").
		AddFilter("Markdown Files", "*.md;*.markdown").
		CanChooseFiles(true).
		PromptForSingleSelection()
}

func (s *AppService) SaveFileDialog(defaultPath string) (string, error) {
	d := s.app.Dialog.SaveFile().
		AttachToWindow(s.focusedWindow()).
		AddFilter("Markdown Files", "*.md")
	if defaultPath != "" {
		d.SetDirectory(filepath.Dir(defaultPath))
		d.SetFilename(filepath.Base(defaultPath))
	} else {
		d.SetFilename("untitled.md")
	}
	return d.PromptForSingleSelection()
}

func (s *AppService) OpenFolderDialog() (string, error) {
	return s.app.Dialog.OpenFile().
		AttachToWindow(s.focusedWindow()).
		SetTitle("Open Folder").
		CanChooseDirectories(true).
		CanChooseFiles(false).
		CanCreateDirectories(true).
		PromptForSingleSelection()
}

func (s *AppService) ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *AppService) WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func (s *AppService) ListDirectory(path string) ([]FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	files := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		lower := strings.ToLower(name)
		if entry.IsDir() || strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".markdown") {
			files = append(files, FileInfo{
				Name:  name,
				Path:  filepath.Join(path, name),
				IsDir: entry.IsDir(),
			})
		}
	}
	return files, nil
}

func (s *AppService) GetHomePath() string {
	home, _ := os.UserHomeDir()
	return home
}

func (s *AppService) CloseWindow() {
	if w := s.focusedWindow(); w != nil {
		w.Close()
	}
}

func (s *AppService) QuitApp() {
	if s.app != nil {
		s.app.Quit()
	}
}

func (s *AppService) ShowSaveDialog(filename string) string {
	done := make(chan string, 1)

	title := "当前文档有未保存的更改"
	if filename != "" {
		title = fmt.Sprintf(`"%s" 有未保存的更改`, filename)
	}

	dlg := s.app.Dialog.Question().
		SetTitle(title).
		SetMessage("关闭前是否保存更改？").
		AttachToWindow(s.focusedWindow())

	dlg.AddButton("取消").OnClick(func() { done <- "cancel" }).SetAsCancel()
	dlg.AddButton("不保存").OnClick(func() { done <- "discard" })
	dlg.AddButton("保存").OnClick(func() { done <- "save" }).SetAsDefault()
	dlg.Show()

	return <-done
}
```

- [ ] **Step 3: Build to verify no Go errors**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md && go build . 2>&1 | grep -v "^ld: warning"
```

Expected: errors about `s.window` in `savedialog_darwin.go` (we fix that next), no other errors.

- [ ] **Step 4: Commit**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md
git add app.go
git commit -m "refactor: remove window field from AppService, add CloseWindow"
```

---

### Task 2: Fix `savedialog_darwin.go` to use focused window

**Files:**
- Modify: `savedialog_darwin.go`

**Context:** `ShowCloseSheet` currently calls `s.window.NativeWindow()`. With `window` field removed, it must call `s.focusedWindow().NativeWindow()` instead.

- [ ] **Step 1: Read `savedialog_darwin.go`**

Read `/Users/yanshili/Downloads/md-p-1/fast-md/savedialog_darwin.go` to confirm current state.

- [ ] **Step 2: Replace `s.window.NativeWindow()` with `s.focusedWindow().NativeWindow()`**

Find this line in `ShowCloseSheet`:
```go
C.showCloseSheet(C.uint(id), unsafe.Pointer(s.window.NativeWindow()), fn, dir)
```

Replace with:
```go
C.showCloseSheet(C.uint(id), unsafe.Pointer(s.focusedWindow().NativeWindow()), fn, dir)
```

- [ ] **Step 3: Build to verify no errors**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md && go build . 2>&1 | grep -v "^ld: warning"
```

Expected: only `# changeme` line (package name), no errors.

- [ ] **Step 4: Commit**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md
git add savedialog_darwin.go
git commit -m "fix: use focused window in ShowCloseSheet"
```

---

### Task 3: Extract `newEditorWindow` and update `main.go`

**Files:**
- Modify: `main.go`

**Context:** Currently `main()` creates one window inline and stores it in `svc.window`. We need to:
1. Remove `svc.window = window` (field no longer exists)
2. Extract window creation into `newEditorWindow(app, svc)`
3. Change Cmd+N menu item to call `newEditorWindow` directly (not emit `menu:new`)
4. Change all `app.Event.Emit(...)` for single-window events to `app.Window.GetFocused().EmitEvent(...)`
5. Update Dock reopen handler to call `newEditorWindow`

- [ ] **Step 1: Read `main.go`**

Read `/Users/yanshili/Downloads/md-p-1/fast-md/main.go` to confirm current state.

- [ ] **Step 2: Rewrite `main.go`**

Replace the entire `main.go` content with:

```go
package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

type themeEntry struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

func loadThemes() []themeEntry {
	data, err := fs.ReadFile(assets, "frontend/dist/themes/index.json")
	if err != nil {
		return nil
	}
	var entries []themeEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil
	}
	return entries
}

func newEditorWindow(app *application.App) *application.WebviewWindow {
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:          "fast-md",
		Width:          1280,
		Height:         800,
		MinWidth:       600,
		MinHeight:      400,
		URL:            "/",
		EnableFileDrop: true,
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropNormal,
			TitleBar:                application.MacTitleBarHiddenInset,
			InvisibleTitleBarHeight: 28,
		},
	})

	window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
		files := event.Context().DroppedFiles()
		for _, f := range files {
			lower := strings.ToLower(f)
			if strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".markdown") {
				window.EmitEvent("file:open", f)
				break
			}
		}
	})

	window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		event.Cancel()
		window.EmitEvent("app:closeRequested")
	})

	window.RegisterHook(events.Mac.WindowShow, func(_ *application.WindowEvent) {
		positionTrafficLights(window.NativeWindow(), 13, 14)
		setupTopBorderDoubleClick(window.NativeWindow())
	})

	return window
}

func main() {
	svc := &AppService{}

	app := application.New(application.Options{
		Name:        "fast-md",
		Description: "A fast Markdown editor",
		Services: []application.Service{
			application.NewService(svc),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	svc.app = app

	newEditorWindow(app)

	// Dock icon click: open a new window
	app.Event.OnApplicationEvent(events.Mac.ApplicationShouldHandleReopen, func(_ *application.ApplicationEvent) {
		newEditorWindow(app)
	})

	buildMenu(app)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func emitToFocused(app *application.App, name string, data ...any) {
	if w := app.Window.GetFocused(); w != nil {
		w.EmitEvent(name, data...)
	}
}

func buildMenu(app *application.App) {
	menu := app.NewMenu()

	if runtime.GOOS == "darwin" {
		menu.AddRole(application.AppMenu)
	}

	fileMenu := menu.AddSubmenu("File")
	fileMenu.Add("New Window").SetAccelerator("CmdOrCtrl+N").OnClick(func(_ *application.Context) {
		newEditorWindow(app)
	})
	fileMenu.Add("Open...").SetAccelerator("CmdOrCtrl+O").OnClick(func(_ *application.Context) {
		emitToFocused(app, "menu:open")
	})
	fileMenu.AddSeparator()
	fileMenu.Add("Save").SetAccelerator("CmdOrCtrl+S").OnClick(func(_ *application.Context) {
		emitToFocused(app, "menu:save")
	})
	fileMenu.Add("Save As...").SetAccelerator("CmdOrCtrl+Shift+S").OnClick(func(_ *application.Context) {
		emitToFocused(app, "menu:saveAs")
	})
	fileMenu.AddSeparator()
	fileMenu.Add("Export as HTML...").OnClick(func(_ *application.Context) {
		emitToFocused(app, "menu:exportHTML")
	})
	fileMenu.Add("Export as PDF (Print)...").OnClick(func(_ *application.Context) {
		emitToFocused(app, "menu:exportPDF")
	})
	fileMenu.AddSeparator()
	fileMenu.Add("Quit fast-md").SetAccelerator("CmdOrCtrl+Q").OnClick(func(_ *application.Context) {
		app.Event.Emit("app:quitRequested")
	})

	menu.AddRole(application.EditMenu)

	viewMenu := menu.AddSubmenu("View")
	viewMenu.Add("Toggle Sidebar").SetAccelerator("CmdOrCtrl+\\").OnClick(func(_ *application.Context) {
		emitToFocused(app, "menu:toggleSidebar")
	})
	viewMenu.AddSeparator()
	themeMenu := viewMenu.AddSubmenu("Theme")
	themeMenu.Add("Default").OnClick(func(_ *application.Context) {
		emitToFocused(app, "menu:setTheme", "default")
	})
	if themes := loadThemes(); len(themes) > 0 {
		themeMenu.AddSeparator()
		for _, t := range themes {
			name := t.Name
			label := t.Label
			themeMenu.Add(label).OnClick(func(_ *application.Context) {
				emitToFocused(app, "menu:setTheme", name)
			})
		}
	}
	viewMenu.AddSeparator()
	viewMenu.Add("Enter Full Screen").SetAccelerator("Ctrl+CmdOrCtrl+F").OnClick(func(_ *application.Context) {
		emitToFocused(app, "menu:fullscreen")
	})
	viewMenu.AddSeparator()
	viewMenu.Add("Developer Tools").SetAccelerator("F12").OnClick(func(_ *application.Context) {
		if w := app.Window.GetFocused(); w != nil {
			toggleDevTools(w.NativeWindow())
		}
	})

	app.Menu.Set(menu)
}
```

- [ ] **Step 3: Build to verify**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md && go build . 2>&1 | grep -v "^ld: warning"
```

Expected: only `# changeme` line, no errors.

- [ ] **Step 4: Commit**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md
git add main.go
git commit -m "feat: extract newEditorWindow, route menu events to focused window"
```

---

### Task 4: Regenerate bindings and update `App.vue`

**Files:**
- Modify: `frontend/bindings/changeme/appservice.ts` (regenerated)
- Modify: `frontend/src/App.vue`

**Context:** `AppService` now has `CloseWindow` instead of `HideWindow`/`ShowWindow`. Bindings must be regenerated, then `App.vue` updated to call `CloseWindow()` instead of `HideWindow()`.

- [ ] **Step 1: Regenerate bindings**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md && wails3 generate bindings -f '' -clean=true -ts
```

Expected: `frontend/bindings/changeme/appservice.ts` updated — `CloseWindow` present, `HideWindow` and `ShowWindow` absent.

Verify:
```bash
grep -E "CloseWindow|HideWindow|ShowWindow" /Users/yanshili/Downloads/md-p-1/fast-md/frontend/bindings/changeme/appservice.ts
```

Expected output:
```
export function CloseWindow(): $CancellablePromise<void> {
```

- [ ] **Step 2: Read `App.vue`**

Read `/Users/yanshili/Downloads/md-p-1/fast-md/frontend/src/App.vue` to confirm current import and `executeClose` lines.

- [ ] **Step 3: Update import in `App.vue`**

Find:
```ts
import { HideWindow, QuitApp, ShowSaveDialog, ShowCloseSheet } from '../bindings/changeme/appservice'
```

Replace with:
```ts
import { CloseWindow, QuitApp, ShowSaveDialog, ShowCloseSheet } from '../bindings/changeme/appservice'
```

- [ ] **Step 4: Update `executeClose` in `App.vue`**

Find:
```ts
function executeClose(action: 'hide' | 'quit') {
  if (action === 'quit') QuitApp()
  else HideWindow()
}
```

Replace with:
```ts
function executeClose(action: 'hide' | 'quit') {
  if (action === 'quit') QuitApp()
  else CloseWindow()
}
```

- [ ] **Step 5: Build frontend to verify**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md/frontend && npm run build 2>&1 | grep -E "error|✓ built"
```

Expected: `✓ built in ...`

- [ ] **Step 6: Build Go binary to verify full build**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md && go build . 2>&1 | grep -v "^ld: warning"
```

Expected: no errors.

- [ ] **Step 7: Commit**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md
git add frontend/bindings/changeme/appservice.ts frontend/src/App.vue
git commit -m "feat: use CloseWindow binding in App.vue for multi-window close"
```

---

### Task 5: Manual smoke test

No code changes — verification only.

- [ ] **Step 1: Run the app**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md && go run . &
```

- [ ] **Step 2: Verify Cmd+N opens a new independent window**

Press Cmd+N. A second editor window should appear. Both windows should be blank and independent.

- [ ] **Step 3: Verify file isolation**

In window 1, type some text. In window 2, verify it is still blank. Save window 1 (Cmd+S) — window 2 should not be affected.

- [ ] **Step 4: Verify Cmd+S targets focused window**

Open a file in window 1. Edit it. Click window 2 to focus it. Press Cmd+S — only window 2's save dialog should appear (or save silently if window 2 has a file path).

- [ ] **Step 5: Verify close behavior**

Close window 1 (red button). If it has unsaved changes, the save dialog should appear. Window 2 should remain open and unaffected.

- [ ] **Step 6: Verify last-window-close keeps app running**

Close all windows. The Dock icon should remain. Click the Dock icon — a new blank window should open.

- [ ] **Step 7: Verify drag-and-drop targets correct window**

Drag a `.md` file onto window 2. Window 2 should open the file; window 1 should be unaffected.
