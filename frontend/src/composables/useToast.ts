import { ref } from 'vue'

export interface ToastAction {
  label: string
  onClick: () => void
}

export interface ToastSpec {
  message: string
  action?: ToastAction
  /** Auto-dismiss after this many ms. Pass 0 to require manual close. */
  duration?: number
}

// Module-level singleton so any caller (composable, store, component) can
// post a toast without prop-drilling a ref through the component tree.
interface ActiveToast extends ToastSpec {
  /** Bumped on every show() so a fresh toast replaces the visible one and
   *  re-starts its dismiss timer (otherwise the same spec would not
   *  retrigger a watcher). */
  nonce: number
}

const active = ref<ActiveToast | null>(null)
let nextNonce = 0

export function useToast() {
  return {
    active,
    show(spec: ToastSpec) {
      active.value = { ...spec, nonce: ++nextNonce }
    },
    dismiss() {
      active.value = null
    },
  }
}
