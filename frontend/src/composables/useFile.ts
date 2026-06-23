import { ref, computed } from 'vue'
import {
  OpenFileDialog,
  SaveFileDialog,
  OpenFolderDialog,
  ReadFile,
  WriteFile,
} from '../../bindings/changeme/core/appservice'
import { useLocale } from './useLocale'

const filePath = ref<string>('')
const content = ref<string>('')
const isDirty = ref<boolean>(false)
const lastSaved = ref<Date | null>(null)
const isSaving = ref<boolean>(false)
// saveError holds the most recent failure from saveFile / saveAs / saveToPath
// / auto-save. The previous design let WriteFile rejections bubble up as
// unhandled promise rejections, so the user saw "Saving..." resolve into
// "Unsaved changes" with no explanation. surfacing the error here lets the
// StatusBar render it next to the save indicator.
const saveError = ref<string | null>(null)
let autoSaveTimer: ReturnType<typeof setTimeout> | null = null

const defaultSettings = {
  autoSave: true,
  autoSaveInterval: 10,
}

function toErrorMessage(err: unknown): string {
  if (err instanceof Error) return err.message
  return String(err)
}

export function setContent(newContent: string) {
  content.value = newContent
  isDirty.value = true
  scheduleAutoSave()
}

function getSettings() {
  try {
    return { ...defaultSettings, ...JSON.parse(localStorage.getItem('fast-md-settings') || '{}') }
  } catch {
    return defaultSettings
  }
}

function scheduleAutoSave() {
  if (autoSaveTimer) clearTimeout(autoSaveTimer)
  if (!filePath.value) return

  const settings = getSettings()
  if (!settings.autoSave) return

  const interval = settings.autoSaveInterval * 1000
  autoSaveTimer = setTimeout(async () => {
    if (isDirty.value && filePath.value) {
      isSaving.value = true
      try {
        await WriteFile(filePath.value, content.value)
        isDirty.value = false
        lastSaved.value = new Date()
        saveError.value = null
      } catch (err) {
        // Disk full, permission denied, trust-check failure, etc. Surface
        // it via the StatusBar; the next successful save clears it.
        saveError.value = toErrorMessage(err)
      } finally {
        isSaving.value = false
      }
    }
  }, interval)
}

export function useFile() {
  const { t } = useLocale()

  function confirmDiscardIfDirty(): boolean {
    if (!isDirty.value) return true
    const name = filePath.value ? filePath.value.split('/').pop() : t('untitled')
    return window.confirm(t('dialog.discardUnsavedFile').replace('{name}', name ?? t('untitled')))
  }

  function resetFile() {
    content.value = ''
    filePath.value = ''
    isDirty.value = false
    lastSaved.value = null
    if (autoSaveTimer) {
      clearTimeout(autoSaveTimer)
      autoSaveTimer = null
    }
  }

  function newFile() {
    if (!confirmDiscardIfDirty()) return
    resetFile()
  }

  async function openFile(path?: string) {
    if (!confirmDiscardIfDirty()) return
    const targetPath = path ?? (await OpenFileDialog())
    if (!targetPath) return
    const fileContent = await ReadFile(targetPath)
    filePath.value = targetPath
    content.value = fileContent
    isDirty.value = false
    lastSaved.value = new Date()
  }

  async function saveFile(): Promise<{ path: string } | null> {
    if (isSaving.value) return null
    if (!filePath.value) {
      return await saveAs()
    }
    isSaving.value = true
    try {
      await WriteFile(filePath.value, content.value)
      isDirty.value = false
      lastSaved.value = new Date()
      saveError.value = null
      return { path: filePath.value }
    } catch (err) {
      saveError.value = toErrorMessage(err)
      return null
    } finally {
      isSaving.value = false
    }
  }

  async function saveAs(): Promise<{ path: string } | null> {
    if (isSaving.value) return null
    const newPath = await SaveFileDialog(filePath.value)
    if (!newPath) return null
    filePath.value = newPath
    isSaving.value = true
    try {
      await WriteFile(newPath, content.value)
      isDirty.value = false
      lastSaved.value = new Date()
      saveError.value = null
      return { path: newPath }
    } catch (err) {
      saveError.value = toErrorMessage(err)
      return null
    } finally {
      isSaving.value = false
    }
  }

  async function saveToPath(path: string): Promise<{ path: string } | null> {
    if (isSaving.value) return null
    isSaving.value = true
    try {
      await WriteFile(path, content.value)
      filePath.value = path
      isDirty.value = false
      lastSaved.value = new Date()
      saveError.value = null
      return { path }
    } catch (err) {
      saveError.value = toErrorMessage(err)
      return null
    } finally {
      isSaving.value = false
    }
  }

  async function openFolder(): Promise<string> {
    return await OpenFolderDialog()
  }

  const fileName = computed(() => {
    if (!filePath.value) return t('untitled')
    return filePath.value.split('/').pop() ?? t('untitled')
  })

  const saveStatus = computed(() => {
    if (isSaving.value) return t('file.saving')
    if (!isDirty.value && lastSaved.value) return t('file.saved')
    if (isDirty.value) return t('file.unsavedChanges')
    return ''
  })

  return {
    filePath,
    content,
    isDirty,
    lastSaved,
    isSaving,
    saveError,
    saveStatus,
    fileName,
    setContent,
    resetFile,
    newFile,
    openFile,
    saveFile,
    saveAs,
    saveToPath,
    openFolder,
  }
}
