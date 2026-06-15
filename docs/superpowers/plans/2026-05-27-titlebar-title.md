# Titlebar Document Title Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Show a centered document title in the titlebar (same row as macOS traffic-light buttons), with a gray " — 已编辑" suffix when the document has unsaved changes.

**Architecture:** Single-file change to `App.vue` — add a `computed` property for the title string, render it as an absolutely-positioned `<div>` inside the existing `.editor-titlebar`, and add two CSS rules.

**Tech Stack:** Vue 3 (Composition API), existing `useFile` composable (`filePath`, `isDirty`)

---

### Task 1: Add titlebar title to App.vue

**Files:**
- Modify: `frontend/src/App.vue`

- [ ] **Step 1: Read the current file**

Read `frontend/src/App.vue` to confirm the exact current state before editing.

- [ ] **Step 2: Add the `titlebarTitle` computed property**

In the `<script setup>` block, after the existing `const { filePath, content, isDirty, ... } = useFile()` line, add:

```ts
import { ref, computed, onMounted, onUnmounted } from 'vue'
```

Replace the existing `import { ref, onMounted, onUnmounted } from 'vue'` with the line above, then add after the `useFile()` destructure:

```ts
const titlebarTitle = computed(() =>
  filePath.value ? filePath.value.split('/').pop()! : '未命名'
)
```

- [ ] **Step 3: Add the title element to the template**

Inside `.editor-titlebar`, after the closing `</button>` (or after the `v-if` button block), add:

```html
<div class="titlebar-title">
  {{ titlebarTitle }}<span v-if="isDirty" class="titlebar-edited"> — 已编辑</span>
</div>
```

The full `.editor-titlebar` block should look like:

```html
<div class="editor-titlebar">
  <button v-if="!sidebarOpen" class="sidebar-expand-btn" @click="sidebarOpen = true" title="Show Sidebar">
    <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
      <path d="M10 2l4 6-4 6V2z"/>
    </svg>
  </button>
  <div class="titlebar-title">
    {{ titlebarTitle }}<span v-if="isDirty" class="titlebar-edited"> — 已编辑</span>
  </div>
</div>
```

- [ ] **Step 4: Add CSS rules**

In the `<style>` block of `App.vue`, add `position: relative` to `.editor-titlebar` and add the two new rules:

```css
.editor-titlebar {
  height: 38px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  padding-left: 78px;
  -webkit-app-region: drag;
  background: var(--bg-primary);
  position: relative;
}
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

- [ ] **Step 5: Build to verify**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md/frontend && npm run build 2>&1 | tail -5
```

Expected: build completes with no errors (warnings about unused vars are fine).

- [ ] **Step 6: Commit**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md
git add frontend/src/App.vue
git commit -m "feat: show centered document title in titlebar"
```
