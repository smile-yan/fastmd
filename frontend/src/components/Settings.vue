<template>
  <div class="settings-overlay" :class="{ 'hide-line-numbers': !settings.showLineNumbers }" @click.self="emit('close')">
    <div class="settings-panel">
      <div class="settings-header">
        <h2>{{ t('settings.title') }}</h2>
        <button class="settings-close" @click="emit('close')" :title="t('menu.close')">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M4 4l8 8M4 12l8-8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
      </div>
      <div class="settings-body">
        <nav class="settings-sidebar">
          <button
            v-for="cat in categories"
            :key="cat.id"
            class="settings-category"
            :class="{ active: activeCategory === cat.id }"
            @click="activeCategory = cat.id"
          >
            <span class="settings-category-icon"><i :class="cat.icon"></i></span>
            <span class="settings-category-label">{{ t('settings.' + cat.id) }}</span>
          </button>
        </nav>
        <div class="settings-content">
          <!-- General -->
          <template v-if="activeCategory === 'general'">
            <section class="settings-section">
              <div class="settings-item">
                <span class="settings-item-label">{{ t('settings.language') }}</span>
                <select v-model="settings.language" class="settings-select" @change="handleLocaleChange">
                  <option value="zh">{{ t('settings.spellZh') }}</option>
                  <option value="en">{{ t('settings.spellEn') }}</option>
                </select>
              </div>
              <div class="settings-item">
                <span class="settings-item-label">{{ t('settings.update') }}</span>
                <div class="settings-item-right">
                  <label class="settings-radio">
                    <input v-model="settings.autoUpdate" type="radio" :value="true" />
                    <span>{{ t('settings.autoUpdate') }}</span>
                  </label>
                  <button class="settings-btn">{{ t('settings.checkUpdate') }}</button>
                </div>
              </div>
              <div class="settings-item">
                <span class="settings-item-label">{{ t('settings.save') }}</span>
                <div class="settings-item-right">
                  <select v-model="settings.autoSave" class="settings-select">
                    <option :value="false">{{ t('settings.manualSave') }}</option>
                    <option :value="true">{{ t('settings.autoSave') }}</option>
                  </select>
                  <template v-if="settings.autoSave">
                    <input v-model.number="settings.autoSaveInterval" type="number" min="5" max="300" class="settings-input" />
                    <span class="settings-input-suffix">{{ t('settings.seconds') }}</span>
                  </template>
                </div>
              </div>
            </section>
          </template>

          <!-- Editor -->
          <template v-if="activeCategory === 'editor'">
            <div class="settings-item">
              <span class="settings-item-label">{{ t('settings.font') }}</span>
              <select v-model="settings.fontFamily" class="settings-select">
                <option value="system">{{ t('settings.fontSystem') }}</option>
                <option value="serif">{{ t('settings.fontSerif') }}</option>
                <option value="sans-serif">{{ t('settings.fontSansSerif') }}</option>
                <option value="monospace">{{ t('settings.fontMonospace') }}</option>
              </select>
            </div>
            <div class="settings-item">
              <span class="settings-item-label">{{ t('settings.fontSize') }}</span>
              <div class="settings-item-right">
                <input v-model.number="settings.fontSize" type="number" min="10" max="32" class="settings-input" />
                <span class="settings-input-suffix">px</span>
              </div>
            </div>
            <div class="settings-item">
              <span class="settings-item-label">{{ t('settings.lineHeight') }}</span>
              <select v-model="settings.lineHeight" class="settings-select">
                <option :value="1.4">{{ t('settings.lineHeightCompact') }} (1.4)</option>
                <option :value="1.7">{{ t('settings.lineHeightNormal') }} (1.7)</option>
                <option :value="2.0">{{ t('settings.lineHeightRelaxed') }} (2.0)</option>
              </select>
            </div>
            <div class="settings-item">
              <span class="settings-item-label">{{ t('settings.showLineNumbers') }}</span>
              <select v-model="settings.showLineNumbers" class="settings-select">
                <option :value="true">{{ t('settings.on') }}</option>
                <option :value="false">{{ t('settings.off') }}</option>
              </select>
            </div>
            <div class="settings-item">
              <span class="settings-item-label">{{ t('settings.bracketMatch') }}</span>
              <select v-model="settings.bracketMatch" class="settings-select">
                <option value="true">{{ t('settings.on') }}</option>
                <option value="false">{{ t('settings.off') }}</option>
              </select>
            </div>
            <div class="settings-item">
              <span class="settings-item-label">{{ t('settings.tabSpaces') }}</span>
              <div class="settings-item-right">
                <input v-model.number="settings.tabSpaces" type="number" min="1" max="8" class="settings-input" />
                <span class="settings-input-suffix">{{ t('settings.seconds') }}</span>
              </div>
            </div>
          </template>

          <!-- Image -->
          <template v-if="activeCategory === 'image'">
            <div class="settings-item">
              <span class="settings-item-label">{{ t('settings.imageInsert') }}</span>
              <select v-model="settings.imageInsert" class="settings-select">
                <option value="upload">{{ t('settings.imageUpload') }}</option>
                <option value="local">{{ t('settings.imageLocal') }}</option>
                <option value="copy">{{ t('settings.imageCopy') }}</option>
              </select>
            </div>
            <div class="settings-item">
              <span class="settings-item-label">{{ t('settings.imageWidth') }}</span>
              <div class="settings-item-right">
                <input v-model.number="settings.imageWidth" type="number" min="0" max="100" class="settings-input" />
                <span class="settings-input-suffix">%</span>
              </div>
            </div>
            <div class="settings-item">
              <span class="settings-item-label">{{ t('settings.imageDisplay') }}</span>
              <select v-model="settings.imageDisplay" class="settings-select">
                <option value="auto">{{ t('settings.imageAuto') }}</option>
                <option value="block">{{ t('settings.imageBlock') }}</option>
                <option value="inline">{{ t('settings.imageInline') }}</option>
              </select>
            </div>
          </template>

          <!-- Appearance -->
          <template v-if="activeCategory === 'appearance'">
            <div class="settings-item">
              <span class="settings-item-label">{{ t('settings.appTheme') }}</span>
              <select v-model="settings.theme" class="settings-select">
                <option value="system">{{ t('settings.appThemeSystem') }}</option>
                <option value="light">{{ t('settings.appThemeLight') }}</option>
                <option value="dark">{{ t('settings.appThemeDark') }}</option>
              </select>
            </div>
            <div class="settings-item">
              <span class="settings-item-label">{{ t('settings.contentTheme') }}</span>
              <select v-model="settings.contentTheme" class="settings-select" @change="handleContentThemeChange">
                <option value="github">GitHub</option>
              </select>
            </div>
            <div class="settings-item">
              <span class="settings-item-label">{{ t('settings.sidebar') }}</span>
              <select v-model="settings.showSidebarByDefault" class="settings-select">
                <option :value="true">{{ t('settings.sidebarShow') }}</option>
                <option :value="false">{{ t('settings.sidebarHide') }}</option>
              </select>
            </div>
          </template>

        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { useLocale } from '../composables/useLocale'
import { SaveConfig, RestartApp } from '../../bindings/changeme/core/appservice'

const emit = defineEmits<{
  close: []
  themeChange: [theme: string]
  contentThemeChange: [theme: string]
}>()

const { t, setLocale } = useLocale()

async function handleLocaleChange() {
  const newLocale = settings.language as 'zh' | 'en'
  setLocale(newLocale)

  // Save to config file
  await SaveConfig({ language: newLocale })

  // Ask user to restart
  const confirmed = confirm(t('dialog.restartForLanguageChange'))

  if (confirmed) {
    RestartApp()
  }
}

interface SettingsState {
  autoSave: boolean
  autoSaveInterval: number
  confirmOnClose: boolean
  fontSize: number
  lineHeight: number
  showLineNumbers: boolean
  theme: string
  contentTheme: string
  showSidebarByDefault: boolean
  language: string
  autoUpdate: boolean
  fontFamily: string
  bracketMatch: boolean
  imageInsert: string
  imageWidth: number
  imageDisplay: string
  tabSpaces: number
}

const STORAGE_KEY = 'fast-md-settings'
const CONTENT_THEME_GITHUB = 'github'

function loadSettings(): SettingsState {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) return normalizeSettings({ ...getDefaults(), ...JSON.parse(raw) })
  } catch { /* ignore */ }
  return getDefaults()
}

function normalizeSettings(loaded: SettingsState): SettingsState {
  return { ...loaded, contentTheme: CONTENT_THEME_GITHUB }
}

function getDefaults(): SettingsState {
  return {
    autoSave: true,
    autoSaveInterval: 10,
    confirmOnClose: true,
    fontSize: 16,
    lineHeight: 1.7,
    showLineNumbers: true,
    theme: 'system',
    contentTheme: CONTENT_THEME_GITHUB,
    showSidebarByDefault: false,
    language: 'zh',
    autoUpdate: false,
    fontFamily: 'system',
    bracketMatch: true,
    imageInsert: 'upload',
    imageWidth: 80,
    imageDisplay: 'auto',
    tabSpaces: 4,
  }
}

function saveSettings() {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(settings))
  emit('themeChange', settings.theme)
  window.dispatchEvent(new Event('fast-md-settings-changed'))
}

function handleContentThemeChange() {
  emit('contentThemeChange', settings.contentTheme)
}

const settings = reactive<SettingsState>(loadSettings())

// Auto-save on any change
watch(settings, () => saveSettings(), { deep: true })

const categories = [
  { id: 'general', icon: 'fa-solid fa-gear' },
  { id: 'editor', icon: 'fa-solid fa-pen' },
  { id: 'image', icon: 'fa-solid fa-image' },
  { id: 'appearance', icon: 'fa-solid fa-palette' },
]

const activeCategory = ref('general')
</script>

<style scoped>
.settings-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(2px);
}

.settings-panel {
  width: 640px;
  height: 480px;
  max-width: 92vw;
  max-height: 88vh;
  background: var(--bg-primary);
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.35), 0 0 0 1px var(--border-color);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.settings-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 20px;
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
  background: var(--bg-secondary);
}

.settings-header h2 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.settings-close {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--text-muted);
  padding: 6px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s, color 0.15s;
}
.settings-close:hover {
  background: var(--border-color);
  color: var(--text-primary);
}

.settings-body {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.settings-sidebar {
  width: 28%;
  min-width: 176px;
  flex-shrink: 0;
  border-right: 1px solid var(--border-color);
  padding: 16px 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  background: var(--bg-secondary);
}

.settings-category {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px 18px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-size: 14px;
  text-align: left;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
  width: 100%;
  min-height: 44px;
  box-sizing: border-box;
}
.settings-category:hover {
  background: var(--border-color);
  color: var(--text-primary);
}
.settings-category.active {
  background: var(--accent-color);
  color: #fff;
}
.settings-category-icon {
  width: 20px;
  text-align: center;
  flex-shrink: 0;
}
.settings-category-icon i {
  font-size: 16px;
}
.settings-category-label {
  flex-shrink: 0;
  white-space: nowrap;
}

.settings-content {
  flex: 3;
  padding: 20px 24px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.settings-section {
  flex: 1;
}

.settings-section h3 {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 16px;
}

.settings-row {
  margin-bottom: 16px;
}

.settings-label {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  font-size: 12px;
  color: var(--text-primary);
  cursor: pointer;
}

.settings-label input[type="checkbox"] {
  width: 16px;
  height: 16px;
  accent-color: var(--accent-color);
  cursor: pointer;
  flex-shrink: 0;
}

.settings-desc {
  margin-top: 6px;
  margin-left: 0;
  font-size: 12px;
  color: var(--text-muted);
  line-height: 1.4;
}

.settings-input-group {
  display: flex;
  align-items: center;
  gap: 6px;
}

.settings-input-suffix {
  font-size: 12px;
  color: var(--text-muted);
}

.settings-input {
  width: 72px;
  padding: 6px 10px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 13px;
  outline: none;
  transition: border-color 0.15s;
}
.settings-input:focus {
  border-color: var(--accent-color);
}

.settings-select {
  padding: 6px 10px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 13px;
  outline: none;
  cursor: pointer;
  min-width: 130px;
  transition: border-color 0.15s;
}
.settings-select:focus {
  border-color: var(--accent-color);
}

.settings-placeholder {
  color: var(--text-muted);
  font-size: 12px;
  padding: 20px 0;
}

.settings-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
}
.settings-item:last-child {
  border-bottom: none;
}
.settings-item-label {
  font-size: 13px;
  color: var(--text-primary);
}
.settings-item-right {
  display: flex;
  align-items: center;
  gap: 16px;
}
.settings-radio {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--text-primary);
  cursor: pointer;
}
.settings-radio input {
  accent-color: var(--accent-color);
}
.settings-btn {
  padding: 6px 14px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 12px;
  cursor: pointer;
  transition: border-color 0.15s;
}
.settings-btn:hover {
  border-color: var(--accent-color);
}

/* Hide line numbers when disabled */
:deep(.cm-lineNumbers),
:deep(.cm-gutters) {
  display: none !important;
}

.settings-section .settings-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
}
.settings-section .settings-item:last-child {
  border-bottom: none;
}
.settings-section .settings-item-label {
  font-size: 13px;
  color: var(--text-primary);
}
.settings-section .settings-item-right {
  display: flex;
  align-items: center;
  gap: 16px;
}
.settings-section .settings-radio {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--text-primary);
  cursor: pointer;
}
.settings-section .settings-radio input {
  accent-color: var(--accent-color);
}
.settings-section .settings-btn {
  padding: 6px 14px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 12px;
  cursor: pointer;
  transition: border-color 0.15s;
}
.settings-section .settings-btn:hover {
  border-color: var(--accent-color);
}

</style>
