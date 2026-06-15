# Typora Theme Support Design

## Overview

Two independent sub-systems:
1. **Python converter** (`tools/typora_theme_converter.py`) — offline tool that converts a Typora `.css` file into an app-compatible CSS file
2. **Frontend theme switcher** — runtime theme loading + View menu integration

---

## Sub-system 1: Python Converter

**File:** `tools/typora_theme_converter.py`

**Input:** path to a Typora `.css` file
**Output:** `frontend/public/themes/<name>.css`

**Usage:**
```bash
python tools/typora_theme_converter.py path/to/github.css
# → frontend/public/themes/github.css
# also updates frontend/public/themes/index.json
```

### Conversion rules (applied in order)

1. **Strip irrelevant rules** — drop any rule whose selector contains `.ty-`, `#typora-`, `.md-toc`, `.md-diagram`, `.md-math`, `#write-info-panel`, `@font-face`, `@import`
2. **Selector prefix rewrite:**
   - `#write` → `.milkdown .ProseMirror`
   - `body` (standalone) → `.milkdown .ProseMirror` (only font/color/line-height properties kept)
   - `html` → drop
   - `.md-fences` → `.milkdown .ProseMirror pre`
   - `.md-lang` → `.milkdown .ProseMirror pre .language-label` (best-effort)
   - Any remaining `.md-*` selector → drop
3. **Blockquote structural fix** — Milkdown renders the left bar as `blockquote::before` (absolute-positioned pseudo-element), not `border-left`. After rewriting the selector, if the rule contains `border-left`, extract the color and emit:
   ```css
   .milkdown .ProseMirror blockquote { padding-left: 1em; border-left: none; }
   .milkdown .ProseMirror blockquote::before { background: <extracted-color>; }
   ```
4. **Prepend structural patch block** — inject at top of output:
   ```css
   /* structural patch: elements Milkdown doesn't style by default */
   .milkdown .ProseMirror table { border-collapse: collapse; width: 100%; }
   .milkdown .ProseMirror table th,
   .milkdown .ProseMirror table td { border: 1px solid; padding: 6px 13px; }
   .milkdown .ProseMirror table tr:nth-child(2n) { background-color: rgba(0,0,0,.04); }
   ```
5. **Preserve `:root` variables** — keep as-is (color tokens used by the theme)
6. **Update `index.json`** — read/create `frontend/public/themes/index.json`, add `{"name": "<name>", "label": "<Title Case name>"}` if not already present

---

## Sub-system 2: Frontend Theme Switcher

### Theme storage

`localStorage` key `fast-md-content-theme` (separate from existing `fast-md-theme` which controls light/dark).
Default: `'default'` (no extra stylesheet = current `style.css` styles).

### `useContentTheme.ts` (new composable)

```ts
// frontend/src/composables/useContentTheme.ts
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
    linkEl.href = ''
    linkEl.disabled = true
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

### `App.vue` changes

- Import `useContentTheme`
- Listen for `Events.On('menu:setTheme', (ev) => applyContentTheme(ev.data as string))`

### `main.go` changes

- Remove existing "Toggle Theme" menu item from View menu
- Add "Theme" submenu with:
  - "Default" item → emits `menu:setTheme` with `"default"`
  - Separator
  - One item per theme in `frontend/dist/themes/index.json` (read at startup, fall back to empty list if file missing)
  - Each item emits `menu:setTheme` with the theme name

### `index.json` format

```json
[
  {"name": "github", "label": "GitHub"},
  {"name": "newsprint", "label": "Newsprint"}
]
```

---

## File structure

```
tools/
  typora_theme_converter.py     ← new
frontend/
  public/
    themes/
      index.json                ← new (managed by converter)
      github.css                ← example output
  src/
    composables/
      useContentTheme.ts        ← new
    App.vue                     ← add event listener
main.go                         ← update View menu
```

---

## Out of scope

- Dark-mode variants of Typora themes (Typora ships separate `*-dark.css` files; user runs converter on each)
- Math block styling (Milkdown widget, no CSS selector mapping)
- Task list checkbox styling (Milkdown-specific, not in Typora themes)
