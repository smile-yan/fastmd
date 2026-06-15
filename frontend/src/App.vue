<template>
  <div class="layout">
    <div class="main-area">
      <Sidebar
        ref="sidebarRef"
        :is-open="sidebarOpen"
        :current-file-path="filePath"
        @open-file="handleOpenFile"
        @folder-opened="() => {}"
      />
      <div class="editor-area" data-file-drop-target>
        <div class="editor-titlebar">
          <div class="titlebar-title">
            {{ titlebarTitle }}<span v-if="isDirty" class="titlebar-edited"> — {{ t('edited') }}</span>
          </div>
        </div>
        <div class="editor-scroll">
          <Editor
            :key="editorKey"
            :model-value="content"
            @update:model-value="handleEditorChange"
          />
          <textarea
            v-show="sourceMode"
            ref="sourceTextarea"
            :value="sourceContent"
            @input="handleSourceInput"
            class="source-textarea"
          />
        </div>
      </div>
    </div>
    <StatusBar
      :file-path="filePath"
      :content="content"
      :is-dirty="isDirty"
      :save-status="saveStatus"
    />
    <Settings
      v-if="showSettings"
      @close="showSettings = false"
      @theme-change="handleThemeChange"
      @content-theme-change="applyContentTheme"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { Events } from '@wailsio/runtime'
import Editor from './components/Editor.vue'
import Sidebar from './components/Sidebar.vue'
import StatusBar from './components/StatusBar.vue'
import Settings from './components/Settings.vue'
import { useFile, setContent } from './composables/useFile'
import { useTheme } from './composables/useTheme'
import { useContentTheme } from './composables/useContentTheme'
import { useLocale } from './composables/useLocale'
import { buildMarkdownExportHtml } from './exportHtml'
import { CancelQuit, CloseWindow, ConfirmQuitWindow, ShowSaveDialog, ShowCloseSheet, ExportPDF } from '../bindings/changeme/core/appservice'

const { filePath, content, isDirty, saveStatus, newFile, openFile, saveFile, saveAs, saveToPath, resetFile } = useFile()
const { t } = useLocale()
const titlebarTitle = computed(() =>
  filePath.value ? filePath.value.split('/').pop()! : t('menu.new')
)
const { toggleTheme, setTheme } = useTheme()
const { applyContentTheme } = useContentTheme()

const sidebarOpen = ref(false)
const sidebarRef = ref<InstanceType<typeof Sidebar> | null>(null)
const showSettings = ref(false)
const sourceMode = ref(false)
const sourceContent = ref('')
const sourceTextarea = ref<HTMLTextAreaElement | null>(null)
const editorKey = ref(0)
let isToggling = false

function toggleSourceMode() {
  if (isToggling) return
  isToggling = true
  setTimeout(() => { isToggling = false }, 200)

  if (sourceMode.value) {
    commitSourceModeEdits()
    sourceMode.value = false
    nextTick(() => {
      editorKey.value++
    })
  } else {
    sourceContent.value = content.value
    sourceMode.value = true
    nextTick(() => {
      sourceTextarea.value?.focus()
      sourceTextarea.value?.select()
    })
  }
}

function handleSourceInput(event: Event) {
  const nextContent = (event.target as HTMLTextAreaElement).value
  sourceContent.value = nextContent
  setContent(nextContent)
}

function commitSourceModeEdits() {
  if (sourceMode.value && sourceContent.value !== content.value) {
    setContent(sourceContent.value)
  }
}

// ── Close / quit dialog ────────────────────────────────────────────────────
// Window ID captured from the app:confirmQuitWindow event payload. Go puts
// the originating window's ID into the event so the frontend can echo it
// back to ConfirmQuitWindow — the binding no longer relies on pulling the
// window out of a request context (which Wails 3 generic RPC does not
// populate), so the quit coordinator is guaranteed to advance.
let pendingConfirmWindowID: number | null = null

function executeClose(action: 'close' | 'quit') {
  if (action === 'quit') {
    if (pendingConfirmWindowID == null) {
      // Defensive: if the event payload ever goes missing, fall back to
      // closing this window directly and cancel the queued quit so the
      // other windows can still be asked.
      console.warn('app:confirmQuitWindow fired without a window ID payload')
      CancelQuit()
      CloseWindow()
      return
    }
    ConfirmQuitWindow(pendingConfirmWindowID)
  } else {
    CloseWindow()
  }
}

function cancelClose(action: 'close' | 'quit') {
  if (action === 'quit') CancelQuit()
}

async function requestClose(action: 'close' | 'quit') {
  commitSourceModeEdits()
  const hasUnsaved = isDirty.value || (content.value.trim() !== '' && !filePath.value)
  if (!hasUnsaved) {
    executeClose(action)
    return
  }

  if (filePath.value) {
    // Existing file: simple one-click alert (save writes directly, no second dialog)
    const filename = filePath.value.split('/').pop() ?? ''
    const result = await ShowSaveDialog(filename)
    if (result === 'save') {
      try {
        await saveFile()
      } catch (err) {
        // Save failed — surface the error and abort the close so the user
        // can retry or pick a different location. Previously the rejection
        // bubbled out of the Events.On handler as an unhandled promise
        // rejection and the user saw nothing.
        console.error('Save failed:', err)
        alert(t('dialog.saveFailed'))
        cancelClose(action)
        return
      }
      if (!isDirty.value) executeClose(action)
      else cancelClose(action)
    } else if (result === 'discard') {
      resetFile()
      executeClose(action)
    } else {
      cancelClose(action)
    }
  } else {
    // New unsaved file: single native sheet with filename + location picker
    const result = await ShowCloseSheet('', '')
    if (result.startsWith('save:')) {
      try {
        await saveToPath(result.slice(5))
      } catch (err) {
        console.error('Save failed:', err)
        alert(t('dialog.saveFailed'))
        cancelClose(action)
        return
      }
      executeClose(action)
    } else if (result === 'discard') {
      resetFile()
      executeClose(action)
    } else {
      cancelClose(action)
    }
  }
}

// ── Rest of app logic ──────────────────────────────────────────────────────
function handleEditorChange(val: string) {
  setContent(val)
}

async function handleOpenFile(path: string) {
  await openFile(path)
}

function handleThemeChange(theme: string) {
  if (theme === 'light' || theme === 'dark') {
    setTheme(theme as 'light' | 'dark')
  } else {
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    setTheme(prefersDark ? 'dark' : 'light')
  }
}

async function handleExportHTML() {
  commitSourceModeEdits()
  await nextTick()
  const proseMirror = document.querySelector('.milkdown .ProseMirror')
  const html = buildMarkdownExportHtml({
    title: filePath.value || t('untitled'),
    bodyHtml: proseMirror?.innerHTML ?? '',
  })
  const blob = new Blob([html], { type: 'text/html' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = (filePath.value?.split('/').pop()?.replace(/\.md$/, '') ?? 'export') + '.html'
  a.click()
  URL.revokeObjectURL(url)
}

async function handleExportPDF() {
  commitSourceModeEdits()
  await nextTick()
  const proseMirror = document.querySelector('.milkdown .ProseMirror')
  const title = filePath.value?.split('/').pop()?.replace(/\.md$/, '') ?? 'export'
  const html = buildMarkdownExportHtml({ title, bodyHtml: proseMirror?.innerHTML ?? '' })
  try {
    await ExportPDF(html, title)
  } catch (err) {
    console.error('Export PDF failed:', err)
  }
}

const cleanups: Array<() => void> = []

onMounted(async () => {
  // Open file passed via URL query param (e.g. from Finder "Open With")
  const initialFile = new URLSearchParams(window.location.search).get('file')
  if (initialFile) await openFile(initialFile)

  // Typora-compatible macOS shortcuts handled outside the native menu.
  window.addEventListener('keydown', (e) => {
    const isCommandOnly = e.metaKey && !e.ctrlKey && !e.altKey && !e.shiftKey
    if (isCommandOnly && e.key.toLowerCase() === 'w') {
      e.preventDefault()
      requestClose('close')
    }
    if (isCommandOnly && e.key === '/') {
      e.preventDefault()
      toggleSourceMode()
    }
  })
  cleanups.push(Events.On('menu:newFile', () => {
    commitSourceModeEdits()
    newFile()
    sourceMode.value = false
    editorKey.value++
  }))
  cleanups.push(Events.On('menu:open', () => {
    openFile()
  }))
  cleanups.push(Events.On('file:open', (ev) => openFile(ev.data as string)))
  cleanups.push(Events.On('menu:save', () => {
    commitSourceModeEdits()
    saveFile()
  }))
  cleanups.push(Events.On('menu:saveAs', () => {
    commitSourceModeEdits()
    saveAs()
  }))
  cleanups.push(Events.On('menu:toggleSidebar', () => { sidebarOpen.value = !sidebarOpen.value }))
  cleanups.push(Events.On('menu:settings', () => { showSettings.value = true }))
  cleanups.push(Events.On('menu:toggleTheme', () => toggleTheme()))
  cleanups.push(Events.On('menu:exportHTML', () => handleExportHTML()))
  cleanups.push(Events.On('menu:exportPDF', () => handleExportPDF()))
  cleanups.push(Events.On('menu:setTheme', (ev) => applyContentTheme(ev.data as string)))

  cleanups.push(Events.On('common:WindowFilesDropped', (ev) => {
    const data = ev.data as unknown as { filenames: string[] }
    const md = data?.filenames?.find((f: string) => f.endsWith('.md'))
    if (md) openFile(md)
  }))

  // Go side emits these instead of directly hiding / quitting
  const onCloseRequested = () => requestClose('close')
  window.addEventListener('window:closeRequested', onCloseRequested)
  cleanups.push(() => window.removeEventListener('window:closeRequested', onCloseRequested))

  cleanups.push(Events.On('app:confirmQuitWindow', (ev) => {
    // Go puts the originating window's ID in the payload. Capture it so
    // executeClose('quit') can echo it back to ConfirmQuitWindow.
    const raw = (ev as { data?: unknown })?.data
    let id: number | null = null
    if (typeof raw === 'number' && Number.isFinite(raw)) {
      id = raw
    } else if (typeof raw === 'string') {
      const parsed = Number(raw)
      if (Number.isFinite(parsed)) id = parsed
    }
    pendingConfirmWindowID = id
    requestClose('quit')
  }))

  // Notify Go that this window's frontend is ready to receive events
})

onUnmounted(() => cleanups.forEach(fn => fn()))
</script>

<style>
.layout {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100vh;
  overflow: hidden;
}
.main-area {
  display: flex;
  flex: 1;
  width: 100%;
  min-width: 0;
  overflow: hidden;
}
.editor-area {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow: hidden;
}
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
  max-width: 60%;
  overflow: hidden;
  text-overflow: ellipsis;
}
.titlebar-edited { color: var(--text-muted); }
.editor-scroll {
  flex: 1;
  overflow-y: auto;
  position: relative;
}
.source-textarea {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  border: none;
  resize: none;
  outline: none;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-family: var(--editor-font-family);
  font-size: var(--editor-font-size, 16px);
  line-height: 1.7;
  z-index: 10;
  max-width: 1060px;
  margin: 0 auto;
  padding: 24px 40px;
  box-sizing: border-box;
}

@media (max-width: 1160px) {
  .source-textarea {
    max-width: none;
    padding: 24px 32px;
  }
}

@media (max-width: 680px) {
  .source-textarea {
    padding: 16px 20px;
  }
}
</style>
