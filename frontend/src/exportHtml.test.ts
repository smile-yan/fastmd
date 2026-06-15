import { describe, expect, it } from 'vitest'
import { buildMarkdownExportHtml, sanitizeBodyHtml } from './exportHtml'

describe('markdown export html', () => {
  it('embeds Typora GitHub default theme styles', () => {
    const html = buildMarkdownExportHtml({
      title: 'note',
      bodyHtml: '<h1>Hello world</h1><p>Body</p>',
    })

    expect(html).toContain('.markdown-body h1')
    expect(html).toContain('@font-face')
    expect(html).toContain('open-sans-v17-latin-ext_latin-regular.woff2')
    expect(html).toContain('border-bottom: 1px solid #eee')
    expect(html).toContain('.markdown-body blockquote')
    expect(html).toContain('border-left: 4px solid #dfe2e5')
    expect(html).toContain('.markdown-body table tr:nth-child(2n), .markdown-body thead')
    expect(html).toContain('.markdown-body code, .markdown-body tt')
    expect(html).toContain('.markdown-body .md-task-list-item > input')
    expect(html).toContain('.markdown-body pre.md-meta-block')
    expect(html).toContain('<body><main class="markdown-body"><h1>Hello world</h1><p>Body</p></main></body>')
  })

  it('emits a Content-Security-Policy meta tag that disables script execution', () => {
    const html = buildMarkdownExportHtml({
      title: 'note',
      bodyHtml: '<p>Body</p>',
    })

    expect(html).toContain('http-equiv="Content-Security-Policy"')
    expect(html).toContain("default-src 'none'")
    expect(html).toContain("script-src 'none'")
    expect(html).toContain("base-uri 'none'")
    expect(html).toContain("form-action 'none'")
  })
})

describe('sanitizeBodyHtml', () => {
  it('passes safe markdown content through unchanged', () => {
    const input = '<h1>Title</h1><p>Body with <strong>bold</strong> and <em>italic</em>.</p><ul><li>one</li><li>two</li></ul>'
    expect(sanitizeBodyHtml(input)).toBe(input)
  })

  it('drops <script> tags entirely', () => {
    const sanitized = sanitizeBodyHtml('<p>safe</p><script>alert(1)</script>')
    expect(sanitized).not.toContain('<script')
    expect(sanitized).not.toContain('alert(1)')
    expect(sanitized).toContain('<p>safe</p>')
  })

  it('strips inline event handler attributes from every tag', () => {
    const sanitized = sanitizeBodyHtml(
      '<p onclick="alert(1)">click me</p><img src="https://example.com/x.png" onerror="alert(2)">'
    )
    expect(sanitized).not.toContain('onclick')
    expect(sanitized).not.toContain('onerror')
    expect(sanitized).toContain('click me')
    // The img keeps its https src — only the onerror is dropped.
    expect(sanitized).toContain('src="https://example.com/x.png"')
  })

  it('strips <style> attributes and inline event handlers from <a>', () => {
    const sanitized = sanitizeBodyHtml(
      '<a href="https://example.com" style="color:red" onclick="evil()">link</a>'
    )
    expect(sanitized).not.toContain('style=')
    expect(sanitized).not.toContain('onclick')
    expect(sanitized).toContain('href="https://example.com"')
  })

  it('rejects javascript: URLs in href', () => {
    const sanitized = sanitizeBodyHtml('<a href="javascript:alert(1)">x</a>')
    expect(sanitized).not.toContain('href=')
    expect(sanitized).not.toContain('javascript:')
  })

  it('rejects javascript: URLs in img src but keeps the tag', () => {
    const sanitized = sanitizeBodyHtml('<img src="javascript:alert(1)" alt="x">')
    expect(sanitized).not.toContain('javascript:')
    expect(sanitized).not.toContain('src=')
    expect(sanitized).toContain('alt="x"')
  })

  it('neutralises whitespace-obfuscated javascript: URLs', () => {
    const inputs = [
      '<a href="java\tscript:alert(1)">x</a>',
      '<a href="java\nscript:alert(1)">x</a>',
      '<a href=" j\na v\na s cript:alert(1)">x</a>',
    ]
    for (const input of inputs) {
      const sanitized = sanitizeBodyHtml(input)
      expect(sanitized).not.toMatch(/href\s*=\s*["'][^"']*javascript/i)
    }
  })

  it('rejects vbscript: URLs', () => {
    const sanitized = sanitizeBodyHtml('<a href="vbscript:msgbox(1)">x</a>')
    expect(sanitized).not.toContain('href=')
  })

  it('rejects data: URLs that are not images', () => {
    const sanitized = sanitizeBodyHtml(
      '<a href="data:text/html,<script>alert(1)</script>">x</a>'
    )
    expect(sanitized).not.toContain('href=')
    expect(sanitized).not.toContain('data:')
  })

  it('allows data:image srcs for inline images', () => {
    const input = '<img src="data:image/png;base64,iVBORw0KGgo=" alt="x">'
    expect(sanitizeBodyHtml(input)).toBe(input)
  })

  it('drops <iframe>, <object>, <embed> with their subtrees', () => {
    const sanitized = sanitizeBodyHtml(
      '<p>before</p><iframe src="https://evil.example.com"><script>alert(1)</script></iframe><object data="x"></object><embed src="y">after'
    )
    expect(sanitized).not.toContain('<iframe')
    expect(sanitized).not.toContain('<object')
    expect(sanitized).not.toContain('<embed')
    expect(sanitized).not.toContain('evil.example.com')
    expect(sanitized).not.toContain('alert(1)')
    expect(sanitized).toContain('before')
    expect(sanitized).toContain('after')
  })

  it('drops <svg> subtrees (which can carry inline script)', () => {
    const sanitized = sanitizeBodyHtml(
      '<svg onload="alert(1)"><script>alert(2)</script></svg><p>kept</p>'
    )
    expect(sanitized).not.toContain('<svg')
    expect(sanitized).not.toContain('onload')
    expect(sanitized).not.toContain('alert')
    expect(sanitized).toContain('kept')
  })

  it('drops <style> tags so attackers cannot inject CSS', () => {
    const sanitized = sanitizeBodyHtml(
      '<style>body{background:url("javascript:alert(1)")}</style><p>kept</p>'
    )
    expect(sanitized).not.toContain('<style')
    expect(sanitized).not.toContain('background')
    expect(sanitized).toContain('kept')
  })

  it('strips the style attribute (the inline-style XSS vector)', () => {
    const sanitized = sanitizeBodyHtml(
      '<p style="background:url(javascript:alert(1))">x</p>'
    )
    expect(sanitized).not.toContain('style=')
    expect(sanitized).not.toContain('javascript:')
  })

  it('keeps a task-list <input type="checkbox"> with its checked/disabled state', () => {
    const input = '<ul><li class="md-task-list-item"><input type="checkbox" checked disabled> item</li></ul>'
    // The DOM round-trip normalises boolean attributes to `checked=""`,
    // `disabled=""` — accept either form.
    const sanitized = sanitizeBodyHtml(input)
    expect(sanitized).toContain('type="checkbox"')
    expect(sanitized).toMatch(/checked(="")?\s/)
    expect(sanitized).toMatch(/disabled(="")?/)
    expect(sanitized).toContain('item')
  })

  it('removes a non-checkbox <input> (form-control smuggling vector)', () => {
    const sanitized = sanitizeBodyHtml(
      '<input type="text" name="x" value="y">after'
    )
    expect(sanitized).not.toContain('<input')
    expect(sanitized).not.toContain('name=')
    expect(sanitized).toContain('after')
  })

  it('keeps a code block with its language class', () => {
    const input = '<pre><code class="language-rust">fn main() {}</code></pre>'
    expect(sanitizeBodyHtml(input)).toBe(input)
  })

  it('keeps mailto: and relative links', () => {
    const input = '<a href="mailto:a@b.c">mail</a><a href="/local/page">local</a><a href="#anchor">frag</a>'
    expect(sanitizeBodyHtml(input)).toBe(input)
  })

  it('drops disallowed attributes like target and srcset', () => {
    const sanitized = sanitizeBodyHtml(
      '<a href="https://example.com" target="_blank" rel="noopener">x</a>'
    )
    expect(sanitized).not.toContain('target=')
    // rel is in the allowlist; keep it.
    expect(sanitized).toContain('rel="noopener"')
  })
})
