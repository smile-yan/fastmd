import { mount } from '@vue/test-utils'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { defineComponent } from 'vue'

const eventHandlers = new Map<string, (event?: unknown) => unknown>()
const registeredEventNames: string[] = []

vi.mock('@wailsio/runtime', () => ({
  Events: {
    On: vi.fn((name: string, handler: (event?: unknown) => unknown) => {
      registeredEventNames.push(name)
      eventHandlers.set(name, handler)
      return () => eventHandlers.delete(name)
    }),
  },
}))

vi.mock('../bindings/changeme/core/appservice', () => ({
  CancelQuit: vi.fn(),
  CloseWindow: vi.fn(),
  ConfirmQuitWindow: vi.fn(),
  ExportPDF: vi.fn(),
  QuitApp: vi.fn(),
  ShowCloseSheet: vi.fn(),
  ShowSaveDialog: vi.fn(),
  WriteFile: vi.fn(),
  OpenFileDialog: vi.fn(),
  SaveFileDialog: vi.fn(),
  ReadFile: vi.fn(),
  OpenFolderDialog: vi.fn(),
  ListDirectory: vi.fn(),
}))

vi.mock('./components/Editor.vue', () => ({
  default: defineComponent({
    name: 'Editor',
    props: {
      modelValue: {
        type: String,
        default: '',
      },
    },
    emits: ['update:modelValue'],
    methods: {
      handleInput(event: Event) {
        this.$emit('update:modelValue', (event.target as HTMLTextAreaElement).value)
      },
    },
    template: '<textarea class="mock-editor" :value="modelValue" @input="handleInput" />',
  }),
}))

describe('App source mode', () => {
  beforeEach(() => {
    eventHandlers.clear()
    registeredEventNames.length = 0
    localStorage.clear()
    vi.clearAllMocks()
    vi.resetModules()
  })

  it('saves edits made in source mode without requiring a mode toggle first', async () => {
    const bindings = await import('../bindings/changeme/core/appservice')
    vi.mocked(bindings.ReadFile).mockResolvedValue('# Original')
    vi.mocked(bindings.WriteFile).mockResolvedValue(undefined)

    const { default: App } = await import('./App.vue')
    const wrapper = mount(App)

    await eventHandlers.get('file:open')?.({ data: '/tmp/note.md' })
    await wrapper.find('.mock-editor').setValue('# Original')

    window.dispatchEvent(new KeyboardEvent('keydown', { key: '/', metaKey: true }))
    await wrapper.vm.$nextTick()

    await wrapper.find('.source-textarea').setValue('# Source edit')
    await eventHandlers.get('menu:save')?.()

    expect(bindings.WriteFile).toHaveBeenCalledWith('/tmp/note.md', '# Source edit')
  })

  it('uses the Typora macOS source mode shortcut', async () => {
    const { default: App } = await import('./App.vue')
    const wrapper = mount(App)

    window.dispatchEvent(new KeyboardEvent('keydown', { key: '/', metaKey: true }))
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.source-textarea').isVisible()).toBe(true)
  })

  it('does not use the non-Typora source mode shortcut', async () => {
    const { default: App } = await import('./App.vue')
    const wrapper = mount(App)

    window.dispatchEvent(new KeyboardEvent('keydown', { key: '/', metaKey: true, shiftKey: true }))
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.source-textarea').isVisible()).toBe(false)
  })

  it('does not open an external file over unsaved changes when the user cancels', async () => {
    const bindings = await import('../bindings/changeme/core/appservice')
    vi.mocked(bindings.ReadFile)
      .mockResolvedValueOnce('# Original')
      .mockResolvedValueOnce('# External')
    vi.spyOn(window, 'confirm').mockReturnValue(false)

    const { default: App } = await import('./App.vue')
    const wrapper = mount(App)

    await eventHandlers.get('file:open')?.({ data: '/tmp/current.md' })
    await wrapper.find('.mock-editor').setValue('# Unsaved')
    await eventHandlers.get('file:open')?.({ data: '/tmp/external.md' })

    expect(bindings.ReadFile).toHaveBeenCalledTimes(1)
    expect(wrapper.find('.mock-editor').element.value).toBe('# Unsaved')
  })

  it('echoes the window ID from the confirm-quit event back to ConfirmQuitWindow', async () => {
    const bindings = await import('../bindings/changeme/core/appservice')

    const { default: App } = await import('./App.vue')
    mount(App)

    // Go puts the originating window ID in the event payload. The frontend
    // must capture it and pass it back so the quit coordinator advances.
    await eventHandlers.get('app:confirmQuitWindow')?.({ data: 42 })

    expect(bindings.ConfirmQuitWindow).toHaveBeenCalledTimes(1)
    expect(bindings.ConfirmQuitWindow).toHaveBeenCalledWith(42)
    expect(bindings.QuitApp).not.toHaveBeenCalled()
  })

  it('falls back to CancelQuit + CloseWindow if the confirm event has no window ID', async () => {
    const bindings = await import('../bindings/changeme/core/appservice')

    const { default: App } = await import('./App.vue')
    mount(App)

    await eventHandlers.get('app:confirmQuitWindow')?.({ data: null })

    expect(bindings.ConfirmQuitWindow).not.toHaveBeenCalled()
    expect(bindings.CancelQuit).toHaveBeenCalledTimes(1)
    expect(bindings.CloseWindow).toHaveBeenCalledTimes(1)
  })

  it('cancels coordinated quit when the user cancels an unsaved close prompt', async () => {
    const bindings = await import('../bindings/changeme/core/appservice')
    vi.mocked(bindings.ReadFile).mockResolvedValue('# Original')
    vi.mocked(bindings.ShowSaveDialog).mockResolvedValue('cancel')

    const { default: App } = await import('./App.vue')
    const wrapper = mount(App)

    await eventHandlers.get('file:open')?.({ data: '/tmp/note.md' })
    await wrapper.find('.mock-editor').setValue('# Unsaved')
    await eventHandlers.get('app:confirmQuitWindow')?.({ data: 1 })

    expect(bindings.CancelQuit).toHaveBeenCalledTimes(1)
    expect(bindings.ConfirmQuitWindow).not.toHaveBeenCalled()
    expect(bindings.QuitApp).not.toHaveBeenCalled()
  })

  it('does not render the top-left sidebar toggle button', async () => {
    const { default: App } = await import('./App.vue')
    const wrapper = mount(App)

    expect(wrapper.find('.sidebar-expand-btn').exists()).toBe(false)
  })

  it('localizes the edited titlebar marker', async () => {
    localStorage.setItem('fast-md-locale', 'en')

    const { default: App } = await import('./App.vue')
    const wrapper = mount(App)

    await wrapper.find('.mock-editor').setValue('# Unsaved')

    expect(wrapper.find('.titlebar-edited').text()).toBe('— edited')
    expect(wrapper.find('.titlebar-edited').text()).not.toContain('已编辑')
  })

  it('registers external file open events once per app instance', async () => {
    const { default: App } = await import('./App.vue')
    mount(App)

    expect(registeredEventNames.filter((name) => name === 'file:open')).toHaveLength(1)
  })

  it('uses configurable editor font variables in source mode', () => {
    const source = readFileSync(resolve(__dirname, 'App.vue'), 'utf-8')

    expect(source).toMatch(/\.source-textarea\s*{[^}]*font-family:\s*var\(--editor-font-family\)[^;]*;[^}]*font-size:\s*var\(--editor-font-size,\s*16px\);/s)
  })
})
