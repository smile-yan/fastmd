import { describe, it, expect, vi, beforeEach } from 'vitest'

vi.mock('../../bindings/changeme/core/appservice', () => ({
  OpenFileDialog: vi.fn(),
  SaveFileDialog: vi.fn(),
  ReadFile: vi.fn(),
  WriteFile: vi.fn(),
  OpenFolderDialog: vi.fn(),
}))

describe('useFile', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.resetModules()
  })

  it('initializes with empty state', async () => {
    const { useFile } = await import('./useFile')
    const { filePath, content, isDirty } = useFile()
    expect(filePath.value).toBe('')
    expect(content.value).toBe('')
    expect(isDirty.value).toBe(false)
  })

  it('setContent marks dirty', async () => {
    const { useFile, setContent } = await import('./useFile')
    const { isDirty } = useFile()
    setContent('# Hello')
    expect(isDirty.value).toBe(true)
  })

  it('openFile reads file and clears dirty', async () => {
    const bindings = await import('../../bindings/changeme/core/appservice')
    vi.mocked(bindings.OpenFileDialog).mockResolvedValue('/tmp/test.md')
    vi.mocked(bindings.ReadFile).mockResolvedValue('# Content')
    const { useFile } = await import('./useFile')
    const { filePath, content, isDirty, openFile } = useFile()
    await openFile()
    expect(filePath.value).toBe('/tmp/test.md')
    expect(content.value).toBe('# Content')
    expect(isDirty.value).toBe(false)
  })

  it('saveFile writes to current path and returns {path}', async () => {
    const bindings = await import('../../bindings/changeme/core/appservice')
    vi.mocked(bindings.WriteFile).mockResolvedValue(undefined)
    const { useFile, setContent } = await import('./useFile')
    const { filePath, saveFile } = useFile()
    filePath.value = '/tmp/test.md'
    setContent('# Updated')
    const result = await saveFile()
    expect(bindings.WriteFile).toHaveBeenCalledWith('/tmp/test.md', '# Updated')
    expect(result).toEqual({ path: '/tmp/test.md' })
  })

  it('saveFile calls SaveFileDialog when no path and returns {path}', async () => {
    const bindings = await import('../../bindings/changeme/core/appservice')
    vi.mocked(bindings.SaveFileDialog).mockResolvedValue('/tmp/new.md')
    vi.mocked(bindings.WriteFile).mockResolvedValue(undefined)
    const { useFile, setContent } = await import('./useFile')
    const { filePath, saveFile } = useFile()
    setContent('# New')
    const result = await saveFile()
    expect(bindings.SaveFileDialog).toHaveBeenCalled()
    expect(filePath.value).toBe('/tmp/new.md')
    expect(result).toEqual({ path: '/tmp/new.md' })
  })

  it('saveAs returns null when the user cancels the dialog', async () => {
    const bindings = await import('../../bindings/changeme/core/appservice')
    vi.mocked(bindings.SaveFileDialog).mockResolvedValue('')
    const { useFile, setContent } = await import('./useFile')
    const { filePath, saveAs } = useFile()
    setContent('# Draft')
    const result = await saveAs()
    expect(result).toBeNull()
    // filePath must NOT change on cancel
    expect(filePath.value).toBe('')
  })

  it('auto saves existing files every 10 seconds by default', async () => {
    vi.useFakeTimers()
    const bindings = await import('../../bindings/changeme/core/appservice')
    vi.mocked(bindings.WriteFile).mockResolvedValue(undefined)

    const { useFile, setContent } = await import('./useFile')
    const { filePath } = useFile()
    filePath.value = '/tmp/auto.md'

    setContent('# Auto')
    await vi.advanceTimersByTimeAsync(9999)
    expect(bindings.WriteFile).not.toHaveBeenCalled()

    await vi.advanceTimersByTimeAsync(1)
    expect(bindings.WriteFile).toHaveBeenCalledWith('/tmp/auto.md', '# Auto')

    vi.useRealTimers()
  })

  it('newFile resets all state', async () => {
    vi.spyOn(window, 'confirm').mockReturnValue(true)
    const { useFile, setContent } = await import('./useFile')
    const { filePath, content, isDirty, newFile } = useFile()
    filePath.value = '/tmp/old.md'
    setContent('# Old')
    newFile()
    expect(filePath.value).toBe('')
    expect(content.value).toBe('')
    expect(isDirty.value).toBe(false)
  })

  it('localizes file name and save status labels', async () => {
    localStorage.setItem('fast-md-locale', 'zh')

    const { useFile } = await import('./useFile')
    const { fileName, saveStatus } = useFile()

    expect(fileName.value).toBe('未命名')
    expect(saveStatus.value).toBe('')
  })

  it('captures the write error in saveError when WriteFile rejects', async () => {
    const bindings = await import('../../bindings/changeme/core/appservice')
    vi.mocked(bindings.WriteFile).mockRejectedValue(new Error('disk full'))

    const { useFile, setContent } = await import('./useFile')
    const { filePath, saveError, saveFile, isDirty } = useFile()
    filePath.value = '/tmp/test.md'
    setContent('# data')
    isDirty.value = true

    await saveFile()

    expect(saveError.value).toBe('disk full')
    // isDirty must stay true so the user knows the change is still in
    // memory and the dot stays visible.
    expect(isDirty.value).toBe(true)
  })

  it('clears saveError on the next successful save', async () => {
    const bindings = await import('../../bindings/changeme/core/appservice')
    vi.mocked(bindings.WriteFile)
      .mockRejectedValueOnce(new Error('first try fails'))
      .mockResolvedValueOnce(undefined)

    const { useFile, setContent } = await import('./useFile')
    const { filePath, saveError, saveFile, isDirty } = useFile()
    filePath.value = '/tmp/test.md'
    setContent('# data')
    isDirty.value = true

    await saveFile()
    expect(saveError.value).toBe('first try fails')

    await saveFile()
    expect(saveError.value).toBeNull()
    expect(isDirty.value).toBe(false)
  })

  it('captures the error from saveAs when WriteFile rejects', async () => {
    const bindings = await import('../../bindings/changeme/core/appservice')
    vi.mocked(bindings.SaveFileDialog).mockResolvedValue('/tmp/new.md')
    vi.mocked(bindings.WriteFile).mockRejectedValue(new Error('permission denied'))

    const { useFile, setContent } = await import('./useFile')
    const { filePath, saveError, saveAs } = useFile()
    setContent('# data')

    await saveAs()

    expect(saveError.value).toBe('permission denied')
    // filePath is set BEFORE the write attempt, so the user can see
    // where the failed save targeted.
    expect(filePath.value).toBe('/tmp/new.md')
  })

  it('captures the error from auto-save when WriteFile rejects', async () => {
    vi.useFakeTimers()
    const bindings = await import('../../bindings/changeme/core/appservice')
    vi.mocked(bindings.WriteFile).mockRejectedValue(new Error('disk full'))

    const { useFile, setContent } = await import('./useFile')
    const { filePath, saveError } = useFile()
    filePath.value = '/tmp/auto.md'
    setContent('# auto')

    await vi.advanceTimersByTimeAsync(10_000)

    expect(saveError.value).toBe('disk full')

    vi.useRealTimers()
  })

  it('saveToPath writes content to the supplied path and updates state', async () => {
    const bindings = await import('../../bindings/changeme/core/appservice')
    vi.mocked(bindings.WriteFile).mockResolvedValue(undefined)

    const { useFile, setContent } = await import('./useFile')
    const { filePath, content, isDirty, saveError, saveToPath } = useFile()
    setContent('# SaveToPath body')

    await saveToPath('/tmp/chosen.md')

    expect(bindings.WriteFile).toHaveBeenCalledWith('/tmp/chosen.md', '# SaveToPath body')
    expect(filePath.value).toBe('/tmp/chosen.md')
    expect(isDirty.value).toBe(false)
    expect(saveError.value).toBeNull()
    // Content stays populated so the editor continues to show what was saved.
    expect(content.value).toBe('# SaveToPath body')
  })

  it('saveToPath captures WriteFile errors into saveError and keeps isDirty', async () => {
    const bindings = await import('../../bindings/changeme/core/appservice')
    vi.mocked(bindings.WriteFile).mockRejectedValue(new Error('permission denied'))

    const { useFile, setContent } = await import('./useFile')
    const { filePath, isDirty, saveError, saveToPath } = useFile()
    setContent('# Body')
    isDirty.value = true

    await saveToPath('/tmp/locked.md')

    expect(saveError.value).toBe('permission denied')
    expect(isDirty.value).toBe(true)
    // filePath is updated AFTER the write, so a failed save leaves the
    // editor pointed at the previous (still-trusted) path.
    expect(filePath.value).toBe('')
  })

  it('saveToPath is a no-op while a save is already in flight', async () => {
    const bindings = await import('../../bindings/changeme/core/appservice')
    let resolveWrite!: () => void
    vi.mocked(bindings.WriteFile).mockImplementation(
      () => new Promise<void>((res) => { resolveWrite = res }),
    )

    const { useFile, setContent } = await import('./useFile')
    const { saveToPath } = useFile()
    setContent('# Body')

    const first = saveToPath('/tmp/first.md')
    // While the first write is still pending, a second call should be ignored.
    await saveToPath('/tmp/second.md')

    expect(bindings.WriteFile).toHaveBeenCalledTimes(1)
    expect(bindings.WriteFile).toHaveBeenCalledWith('/tmp/first.md', '# Body')

    resolveWrite()
    await first
  })

  it('saveToPath clears a prior saveError on the next successful write', async () => {
    const bindings = await import('../../bindings/changeme/core/appservice')
    vi.mocked(bindings.WriteFile)
      .mockRejectedValueOnce(new Error('first fails'))
      .mockResolvedValueOnce(undefined)

    const { useFile, setContent } = await import('./useFile')
    const { saveError, saveToPath } = useFile()
    setContent('# Body')

    await saveToPath('/tmp/a.md')
    expect(saveError.value).toBe('first fails')

    await saveToPath('/tmp/b.md')
    expect(saveError.value).toBeNull()
  })
})
