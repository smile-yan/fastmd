import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { nextTick } from 'vue'

// useToast stores state at module scope. Each test gets a fresh module
// via vi.resetModules + dynamic import so singletons don't leak between
// cases (see useTheme.test.ts for the same pattern).
describe('useToast', () => {
  beforeEach(() => {
    vi.resetModules()
    vi.useFakeTimers()
  })
  afterEach(() => {
    vi.useRealTimers()
  })

  it('starts with no active toast', async () => {
    const { useToast } = await import('./useToast')
    const { active } = useToast()
    expect(active.value).toBeNull()
  })

  it('show() sets the active toast with a fresh nonce', async () => {
    const { useToast } = await import('./useToast')
    const { active, show } = useToast()
    show({ message: 'hello' })
    expect(active.value?.message).toBe('hello')
    expect(active.value?.nonce).toBe(1)
  })

  it('each show() bumps the nonce', async () => {
    const { useToast } = await import('./useToast')
    const { active, show } = useToast()
    show({ message: 'a' })
    const first = active.value?.nonce
    show({ message: 'b' })
    const second = active.value?.nonce
    expect(second).toBeGreaterThan(first ?? 0)
  })

  it('dismiss() clears the active toast', async () => {
    const { useToast } = await import('./useToast')
    const { active, show, dismiss } = useToast()
    show({ message: 'hello' })
    dismiss()
    expect(active.value).toBeNull()
  })

  it('passes action through unchanged', async () => {
    const { useToast } = await import('./useToast')
    const { active, show } = useToast()
    const onClick = vi.fn()
    show({ message: 'exported', action: { label: 'Reveal', onClick } })
    expect(active.value?.action?.label).toBe('Reveal')
    expect(active.value?.action?.onClick).toBe(onClick)
  })

  it('respects duration: 0 (no auto-dismiss)', async () => {
    const { useToast } = await import('./useToast')
    const { active, show, dismiss } = useToast()
    show({ message: 'persistent', duration: 0 })
    await nextTick()
    vi.advanceTimersByTime(60_000)
    expect(active.value?.message).toBe('persistent')
    dismiss()
  })
})
