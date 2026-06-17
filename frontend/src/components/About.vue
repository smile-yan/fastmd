<template>
  <div class="about-overlay" @click.self="emit('close')">
    <div class="about-panel">
      <div class="about-header">
        <h2>{{ t('menu.about') }}</h2>
        <button class="about-close" @click="emit('close')" :title="t('menu.close')">
          <svg width="16" height="16" viewBox="0 16" fill="none">
            <path d="M4 4l8 8M4 12l8-8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
      </div>
      <div class="about-body">
        <div class="about-mark">
          <div class="about-logo">fastmd</div>
          <div class="about-tagline">{{ t('about.tagline') }}</div>
        </div>
        <dl class="about-meta">
          <div class="about-row">
            <dt>{{ t('about.version') }}</dt>
            <dd>{{ versionLabel }}</dd>
          </div>
          <div v-if="info?.commit" class="about-row">
            <dt>{{ t('about.commit') }}</dt>
            <dd>
              <a class="about-link" :href="commitURL" target="_blank" rel="noreferrer noopener">
                {{ shortCommit }}
              </a>
            </dd>
          </div>
          <div v-if="info?.built" class="about-row">
            <dt>{{ t('about.built') }}</dt>
            <dd>{{ builtLabel }}</dd>
          </div>
          <div class="about-row">
            <dt>{{ t('about.runtime') }}</dt>
            <dd>{{ runtimeLabel }}</dd>
          </div>
        </dl>
        <div class="about-footer">
          <a class="about-link" :href="repoURL" target="_blank" rel="noreferrer noopener">
            {{ t('about.viewOnGitHub') }}
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { GetAppInfo } from '../../bindings/changeme/core/appservice'
import { useLocale } from '../composables/useLocale'

const emit = defineEmits<{
  close: []
}>()

const { t } = useLocale()

const info = ref<{ name: string; version: string; commit: string; built: string; go: string } | null>(null)

onMounted(async () => {
  try {
    info.value = await GetAppInfo()
  } catch {
    // Render with whatever we managed to load; absence of info just hides the rows.
  }
})

const repoURL = 'https://github.com/smile-yan/fast-md'

const versionLabel = computed(() => info.value?.version || 'dev')

const shortCommit = computed(() => {
  const c = info.value?.commit || ''
  return c.length > 7 ? c.slice(0, 7) : c
})

const commitURL = computed(() => {
  const c = info.value?.commit || ''
  if (!c) return repoURL
  return `${repoURL}/commit/${c}`
})

const builtLabel = computed(() => {
  const raw = info.value?.built
  if (!raw) return ''
  // vcs.time is RFC 3339 — show as a locale-friendly short form.
  const d = new Date(raw)
  if (Number.isNaN(d.getTime())) return raw
  return d.toLocaleDateString()
})

const runtimeLabel = computed(() => {
  const go = info.value?.go || ''
  return go ? `Go ${go.replace(/^go/, '')}` : ''
})
</script>

<style scoped>
.about-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(2px);
}

.about-panel {
  width: 380px;
  max-width: 92vw;
  background: var(--bg-primary);
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.35), 0 0 0 1px var(--border-color);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.about-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 16px;
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
  background: var(--bg-secondary);
}

.about-header h2 {
  margin: 0;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
}

.about-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  padding: 0;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  border-radius: 4px;
}

.about-close:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.about-body {
  padding: 20px 24px 18px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.about-mark {
  text-align: center;
}

.about-logo {
  font-size: 22px;
  font-weight: 700;
  letter-spacing: 0.5px;
  color: var(--text-primary);
  background: linear-gradient(120deg, var(--accent-color, #398bff), #b35cff);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.about-tagline {
  margin-top: 4px;
  font-size: 12px;
  color: var(--text-secondary);
}

.about-meta {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin: 0;
  padding: 10px 12px;
  background: var(--bg-secondary);
  border-radius: 8px;
  font-size: 12px;
}

.about-row {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 12px;
}

.about-row dt {
  color: var(--text-secondary);
  flex-shrink: 0;
}

.about-row dd {
  margin: 0;
  color: var(--text-primary);
  text-align: right;
  word-break: break-all;
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 11px;
}

.about-link {
  color: var(--accent-color, #398bff);
  text-decoration: none;
}

.about-link:hover {
  text-decoration: underline;
}

.about-footer {
  display: flex;
  justify-content: center;
}
</style>
