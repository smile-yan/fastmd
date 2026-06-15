# Titlebar Document Title Design

## Overview

Display a centered document title in the titlebar row (same row as macOS traffic-light buttons), matching Typora's behavior.

## Title States

| Condition | Display |
|---|---|
| No file, not edited | `未命名` |
| No file, edited | `未命名` + gray ` — 已编辑` |
| File open, not dirty | `README.md` |
| File open, dirty | `README.md` + gray ` — 已编辑` |

"Edited" means `isDirty === true`. Deleting all content back to empty resets `isDirty` to false via existing `useFile` logic.

## Implementation

**File changed:** `App.vue` only.

**Computed property** `titlebarTitle` returns the base filename string (`未命名` or `filePath` basename).

**Template:** Inside `.editor-titlebar`, add:
```html
<div class="titlebar-title">
  {{ titlebarTitle }}<span v-if="isDirty" class="titlebar-edited"> — 已编辑</span>
</div>
```

**CSS:**
```css
.editor-titlebar { position: relative; }
.titlebar-title {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  font-size: 13px;
  color: var(--text-primary);
  pointer-events: none;
  white-space: nowrap;
}
.titlebar-edited { color: var(--text-muted); }
```
