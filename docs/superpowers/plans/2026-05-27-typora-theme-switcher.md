# Typora Theme Switcher (Frontend + Menu) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a View → Theme submenu that lets users switch between converted Typora themes at runtime, persisted in localStorage.

**Architecture:** New `useContentTheme.ts` composable manages a `<link>` tag in `<head>` pointing to `/themes/<name>.css`. `App.vue` listens for `menu:setTheme` events. `main.go` reads `frontend/dist/themes/index.json` at startup to build the Theme submenu dynamically; falls back to empty list if file is missing.

**Tech Stack:** Vue 3 (Composition API), Go (Wails v3 events), existing `@wailsio/runtime` Events

**Prerequisite:** Plan `2026-05-27-typora-theme-converter.md` must be complete (needs `frontend/public/themes/index.json` and at least one theme CSS).

---

### Task 1: Create `useContentTheme.ts` composable

**Files:**
- Create: `frontend/src/composables/useContentTheme.ts`

- [ ] **Step 1: Create the composable**

Create `/Users/yanshili/Downloads/md-p-1/fast-md/frontend/src/composables/useContentTheme.ts`:

```ts
import { ref } from 'vue'

const STORAGE_KEY = 'fast-md-content-theme'
const contentTheme = ref<string>(localStorage.getItem(STORAGE_KEY) ?? 'default')
let linkEl: HTMLLinkElement | null = null

function applyContentTheme(name: string) {
  if (!linkEl) {
    linkEl = document.createElement('link')
    linkEl.rel = 'stylesheet'
    linkEl.id = 'content-theme-stylesheet'
    document.head.appendChild(linkEl)
  }
  if (name === 'default') {
    linkEl.disabled = true
    linkEl.href = ''
  } else {
    linkEl.disabled = false
    linkEl.href = `/themes/${name}.css`
  }
  localStorage.setItem(STORAGE_KEY, name)
  contentTheme.value = name
}

applyContentTheme(contentTheme.value)

export function useContentTheme() {
  return { contentTheme, applyContentTheme }
}
```

- [ ] **Step 2: Build to verify no TypeScript errors**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md/frontend && npm run build 2>&1 | grep -E "^.*error" | head -10
```

Expected: no output (no errors).

- [ ] **Step 3: Commit**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md
git add frontend/src/composables/useContentTheme.ts
git commit -m "feat: add useContentTheme composable for runtime theme switching"
```

---

### Task 2: Wire `menu:setTheme` event in `App.vue`

**Files:**
- Modify: `frontend/src/App.vue`

- [ ] **Step 1: Read App.vue**

Read `/Users/yanshili/Downloads/md-p-1/fast-md/frontend/src/App.vue` to confirm the current imports and `onMounted` block.

- [ ] **Step 2: Import `useContentTheme`**

In the `<script setup>` block, add after the existing imports:

```ts
import { useContentTheme } from './composables/useContentTheme'
```

And after the existing `const { toggleTheme } = useTheme()` line, add:

```ts
const { applyContentTheme } = useContentTheme()
```

- [ ] **Step 3: Add event listener in `onMounted`**

Inside the `onMounted` block, after the last `cleanups.push(...)` line, add:

```ts
cleanups.push(Events.On('menu:setTheme', (ev) => applyContentTheme(ev.data as string)))
```

- [ ] **Step 4: Build to verify**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md/frontend && npm run build 2>&1 | grep -E "^.*error" | head -10
```

Expected: no errors.

- [ ] **Step 5: Commit**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md
git add frontend/src/App.vue
git commit -m "feat: listen for menu:setTheme event in App.vue"
```

---

### Task 3: Add Theme submenu in `main.go`

**Files:**
- Modify: `main.go`

The Go side needs to:
1. Read `frontend/dist/themes/index.json` (the built output) at startup
2. Replace "Toggle Theme" with a "Theme" submenu containing "Default" + one item per theme

- [ ] **Step 1: Read `main.go`**

Read `/Users/yanshili/Downloads/md-p-1/fast-md/main.go` to confirm the current View menu block (around lines 121–136).

- [ ] **Step 2: Add imports**

The file already imports `"embed"`, `"log"`, `"runtime"`, `"strings"`. Add `"encoding/json"` and `"io/fs"` to the import block:

```go
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
```

- [ ] **Step 3: Add theme index type and loader**

After the `var assets embed.FS` line, add:

```go
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
```

- [ ] **Step 4: Replace "Toggle Theme" with Theme submenu**

In `buildMenu`, find and replace the existing "Toggle Theme" block:

```go
	viewMenu.AddSeparator()
	viewMenu.Add("Toggle Theme").OnClick(func(_ *application.Context) {
		app.Event.Emit("menu:toggleTheme")
	})
```

Replace with:

```go
	viewMenu.AddSeparator()
	themeMenu := viewMenu.AddSubmenu("Theme")
	themeMenu.Add("Default").OnClick(func(_ *application.Context) {
		app.Event.Emit("menu:setTheme", "default")
	})
	if themes := loadThemes(); len(themes) > 0 {
		themeMenu.AddSeparator()
		for _, t := range themes {
			name := t.Name
			label := t.Label
			themeMenu.Add(label).OnClick(func(_ *application.Context) {
				app.Event.Emit("menu:setTheme", name)
			})
		}
	}
```

- [ ] **Step 5: Build to verify**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md && go build . 2>&1 | grep -v "^ld: warning"
```

Expected: only `# changeme` line, no errors.

- [ ] **Step 6: Commit**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md
git add main.go
git commit -m "feat: add Theme submenu to View menu with dynamic theme list"
```

---

### Task 4: Build frontend and verify end-to-end

- [ ] **Step 1: Build frontend (copies themes to dist)**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md/frontend && npm run build 2>&1 | tail -5
```

Expected: build succeeds, `frontend/dist/themes/` directory exists with `index.json` and `github.css`.

- [ ] **Step 2: Verify dist contains themes**

```bash
ls /Users/yanshili/Downloads/md-p-1/fast-md/frontend/dist/themes/
```

Expected: `github.css  index.json`

- [ ] **Step 3: Build Go binary**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md && go build . 2>&1 | grep -v "^ld: warning"
```

Expected: no errors.

- [ ] **Step 4: Manual smoke test**

Run the app. Open View menu → Theme. Verify:
- "Default" item is present
- "GitHub" item is present (separator between them)
- Clicking "GitHub" changes the editor typography/colors to the GitHub theme
- Clicking "Default" restores the original appearance
- Restarting the app restores the last selected theme (localStorage persistence)

- [ ] **Step 5: Commit**

No code changes in this task — it's a verification step only.
