<template>
  <Transition name="toast">
    <div
      v-if="active"
      :key="active.nonce"
      class="toast"
      role="status"
      aria-live="polite"
    >
      <span class="toast-message" :title="active.message">{{ active.message }}</span>
      <button
        v-if="active.action"
        class="toast-action"
        type="button"
        @click="onAction"
      >
        {{ active.action.label }}
      </button>
      <button
        class="toast-close"
        type="button"
        :aria-label="t('common.close')"
        @click="dismiss"
      >
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
          <path
            d="M3 3l8 8M3 11l8-8"
            stroke="currentColor"
            stroke-width="1.5"
            stroke-linecap="round"
          />
        </svg>
      </button>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { onBeforeUnmount, watch } from 'vue'
import { useToast } from '../composables/useToast'
import { useLocale } from '../composables/useLocale'

const { active, dismiss } = useToast()
const { t } = useLocale()

let timer: number | null = null

const clearTimer = () => {
  if (timer !== null) {
    window.clearTimeout(timer)
    timer = null
  }
}

// Restart the dismiss timer every time a fresh toast is posted (the
// `:key="active.nonce"` binding above forces a re-mount on each new
// toast, but we also clear any stale timer from the previous instance
// for symmetry).
watch(
  () => active.value,
  (current) => {
    clearTimer()
    if (!current) return
    const duration = current.duration ?? 3500
    if (duration > 0) {
      timer = window.setTimeout(dismiss, duration)
    }
  },
  { immediate: true },
)

const onAction = () => {
  const a = active.value?.action
  if (a) a.onClick()
  dismiss()
}

onBeforeUnmount(clearTimer)
</script>

<style scoped>
.toast {
  position: fixed;
  right: 24px;
  bottom: 36px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  background: var(--color-bg-elevated, #2a2a2a);
  color: var(--color-text-primary, #f0f0f0);
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
  font-size: 13px;
  line-height: 1.4;
  z-index: 9999;
  max-width: 480px;
  pointer-events: auto;
}
.toast-message {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1 1 auto;
  min-width: 0;
}
.toast-action,
.toast-close {
  background: transparent;
  color: inherit;
  border: none;
  font: inherit;
  cursor: pointer;
  padding: 0;
  display: inline-flex;
  align-items: center;
}
.toast-action {
  color: var(--color-accent, #4f9eff);
  white-space: nowrap;
}
.toast-action:hover {
  text-decoration: underline;
}
.toast-close {
  opacity: 0.6;
}
.toast-close:hover {
  opacity: 1;
}

.toast-enter-active,
.toast-leave-active {
  transition: transform 0.2s ease, opacity 0.2s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
</style>
