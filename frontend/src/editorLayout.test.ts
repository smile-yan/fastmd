import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('editor heading layout', () => {
  it('loads the Typora GitHub Open Sans font assets locally', () => {
    const css = readFileSync(resolve(__dirname, 'style.css'), 'utf-8')

    expect(css).toContain('@font-face')
    expect(css).toContain('font-family: "Open Sans"')
    expect(css).toContain('/themes/github/open-sans-v17-latin-ext_latin-regular.woff2')
    expect(css).toContain('/themes/github/open-sans-v17-latin-ext_latin-700italic.woff2')
  })

  it('uses configurable editor font variables for the document body', () => {
    const css = readFileSync(resolve(__dirname, 'style.css'), 'utf-8')

    expect(css).toMatch(/\.milkdown\s+\.ProseMirror\s*{[^}]*font-family:\s*var\(--editor-font-family\)[^;]*!important;[^}]*font-size:\s*var\(--editor-font-size,\s*16px\)\s*!important;/s)
  })

  it('uses dark theme colors for the editable document surface', () => {
    const css = readFileSync(resolve(__dirname, 'style.css'), 'utf-8')

    expect(css).toMatch(/html\.dark\s+\.milkdown\s+\.ProseMirror\s*{[^}]*background:\s*var\(--bg-primary\)\s*!important;[^}]*color:\s*var\(--text-primary\)\s*!important;/s)
  })

  it('overrides the later Crepe theme font rules inside the editor component', () => {
    const source = readFileSync(resolve(__dirname, 'components/Editor.vue'), 'utf-8')

    expect(source).toMatch(/:deep\(\.milkdown\s+\.ProseMirror\)\s*{[^}]*font-family:\s*var\(--editor-font-family\)\s*!important;[^}]*font-size:\s*var\(--editor-font-size,\s*16px\)\s*!important;/s)
    expect(source).toMatch(/:deep\(\.milkdown\s+\.ProseMirror\s+p\)\s*{[^}]*font-size:\s*var\(--editor-font-size,\s*16px\)\s*!important;/s)
  })

  it('keeps the first heading flush when the virtual cursor is rendered before it', () => {
    const css = readFileSync(resolve(__dirname, 'style.css'), 'utf-8')

    expect(css).toMatch(
      /\.milkdown\s+\.ProseMirror\s*>\s*\.prosemirror-virtual-cursor:first-child\s*\+\s*h1\s*{\s*margin-top:\s*0\s*!important;\s*}/,
    )
  })

  it('uses Typora GitHub heading rules for h1 and h2', () => {
    const css = readFileSync(resolve(__dirname, 'style.css'), 'utf-8')

    expect(css).toMatch(/\.milkdown\s+\.ProseMirror\s+h1,[\s\S]*?\.milkdown\s+\.ProseMirror\s+h6\s*{[^}]*font-family:\s*inherit\s*!important;[^}]*font-weight:\s*bold\s*!important;/s)
    expect(css).toMatch(/\.milkdown\s+\.ProseMirror\s+h1\s*{[^}]*font-size:\s*2\.25em\s*!important;[^}]*line-height:\s*1\.2\s*!important;[^}]*border-bottom:\s*1px\s+solid\s+#eee\s*!important;/s)
    expect(css).toMatch(/\.milkdown\s+\.ProseMirror\s+h2\s*{[^}]*font-size:\s*1\.75em\s*!important;[^}]*line-height:\s*1\.225\s*!important;[^}]*border-bottom:\s*1px\s+solid\s+#eee\s*!important;/s)
  })
})

describe('editor blockquote layout', () => {
  it('matches the Typora GitHub blockquote spacing and divider', () => {
    const css = readFileSync(resolve(__dirname, 'style.css'), 'utf-8')

    expect(css).toMatch(/\.milkdown\s+\.ProseMirror\s+blockquote\s*{[^}]*margin:\s*0\.8em\s+0\s*!important;[^}]*padding:\s*0\s+15px\s*!important;[^}]*color:\s*#777777\s*!important;[^}]*border-left:\s*4px\s+solid\s+#dfe2e5\s*!important;/s)
    expect(css).toMatch(/\.milkdown\s+\.ProseMirror\s+blockquote::before\s*{[^}]*display:\s*none\s*!important;/s)
  })
})

describe('editor Typora GitHub body layout', () => {
  it('matches Typora GitHub table, code block and task-list rules', () => {
    const css = readFileSync(resolve(__dirname, 'style.css'), 'utf-8')

    expect(css).toMatch(/\.milkdown\s+\.ProseMirror\s+table\s*{[^}]*border-collapse:\s*collapse\s*!important;[^}]*word-break:\s*initial\s*!important;/s)
    expect(css).toMatch(/\.milkdown\s+\.ProseMirror\s+table\s+th\s*{[^}]*border:\s*1px\s+solid\s+#dfe2e5\s*!important;[^}]*padding:\s*6px\s+13px\s*!important;/s)
    expect(css).toMatch(/\.milkdown\s+\.ProseMirror\s+pre\s*{[^}]*margin-top:\s*15px\s*!important;[^}]*background-color:\s*#f8f8f8\s*!important;/s)
    expect(css).toMatch(/\.editor-container\s+\.milkdown\s+\.ProseMirror\s+\.milkdown-code-block\s*{[^}]*border:\s*1px\s+solid\s+#e7eaed\s*!important;[^}]*background-color:\s*#f8f8f8\s*!important;/s)
    expect(css).toMatch(/\.milkdown\s+\.ProseMirror\s+\.md-task-list-item\s*>\s*input\s*{[^}]*margin-left:\s*-1\.3em\s*!important;/s)
  })

  it('uses dark theme colors for code blocks and CodeMirror internals', () => {
    const css = readFileSync(resolve(__dirname, 'style.css'), 'utf-8')

    expect(css).toMatch(/html\.dark\s+\.milkdown\s+\.ProseMirror\s+pre,[\s\S]*?html\.dark\s+\.milkdown\s+\.ProseMirror\s+\.milkdown-code-block\s*{[^}]*background-color:\s*var\(--bg-secondary\)\s*!important;[^}]*border-color:\s*var\(--border-color\)\s*!important;/s)
    expect(css).toMatch(/html\.dark\s+\.milkdown\s+\.ProseMirror\s+code,[\s\S]*?html\.dark\s+\.milkdown\s+\.ProseMirror\s+tt\s*{[^}]*background-color:\s*var\(--bg-secondary\)\s*!important;[^}]*border-color:\s*var\(--border-color\)\s*!important;/s)
    expect(css).toMatch(/html\.dark\s+\.milkdown\s+\.ProseMirror\s+\.milkdown-code-block\s+\.cm-editor,[\s\S]*?html\.dark\s+\.milkdown\s+\.ProseMirror\s+\.milkdown-code-block\s+\.cm-activeLineGutter\s*{[^}]*background-color:\s*var\(--bg-secondary\)\s*!important;/s)
  })

  it('keeps Milkdown list item widgets from inheriting normal Typora list spacing', () => {
    const css = readFileSync(resolve(__dirname, 'style.css'), 'utf-8')

    expect(css).toMatch(/\.milkdown\s+\.ProseMirror\s+\.milkdown-list-item-block\s*{[^}]*margin:\s*0\s*!important;[^}]*padding:\s*0\s*!important;/s)
    expect(css).toMatch(/\.milkdown\s+\.ProseMirror\s+\.milkdown-list-item-block\s+ul,[\s\S]*?\.milkdown\s+\.ProseMirror\s+\.milkdown-list-item-block\s+ol\s*{[^}]*margin:\s*0\s*!important;[^}]*padding-left:\s*0\s*!important;/s)
    expect(css).toMatch(/\.milkdown\s+\.ProseMirror\s+\.milkdown-list-item-block\s+p\s*{[^}]*margin:\s*0\s*!important;/s)
    expect(css).toMatch(/\.milkdown\s+\.ProseMirror\s+\.milkdown-list-item-block\s*>\s*\.list-item\s*{[^}]*display:\s*flex\s*!important;[^}]*align-items:\s*flex-start\s*!important;[^}]*gap:\s*10px\s*!important;/s)
    expect(css).toMatch(/\.milkdown\s+\.ProseMirror\s+\.milkdown-list-item-block\s*>\s*\.list-item\s*>\s*\.children\s*{[^}]*min-width:\s*0\s*!important;[^}]*flex:\s*1\s*!important;/s)
    expect(css).toMatch(/\.milkdown\s+\.ProseMirror\s+\.milkdown-list-item-block\s+\.label-wrapper\s*{[^}]*width:\s*24px\s*!important;[^}]*height:\s*32px\s*!important;/s)
  })
})
