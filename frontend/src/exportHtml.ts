export const typoraGithubExportCss = `
@font-face {
  font-family: "Open Sans";
  font-style: normal;
  font-weight: normal;
  src: local("Open Sans Regular"), local("OpenSans-Regular"), url("./github/open-sans-v17-latin-ext_latin-regular.woff2") format("woff2");
  unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+2000-206F, U+2074, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD, U+0100-024F, U+0259, U+1E00-1EFF, U+2020, U+20A0-20AB, U+20AD-20CF, U+2113, U+2C60-2C7F, U+A720-A7FF;
}

@font-face {
  font-family: "Open Sans";
  font-style: italic;
  font-weight: normal;
  src: local("Open Sans Italic"), local("OpenSans-Italic"), url("./github/open-sans-v17-latin-ext_latin-italic.woff2") format("woff2");
  unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+2000-206F, U+2074, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD, U+0100-024F, U+0259, U+1E00-1EFF, U+2020, U+20A0-20AB, U+20AD-20CF, U+2113, U+2C60-2C7F, U+A720-A7FF;
}

@font-face {
  font-family: "Open Sans";
  font-style: normal;
  font-weight: bold;
  src: local("Open Sans Bold"), local("OpenSans-Bold"), url("./github/open-sans-v17-latin-ext_latin-700.woff2") format("woff2");
  unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+2000-206F, U+2074, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD, U+0100-024F, U+0259, U+1E00-1EFF, U+2020, U+20A0-20AB, U+20AD-20CF, U+2113, U+2C60-2C7F, U+A720-A7FF;
}

@font-face {
  font-family: "Open Sans";
  font-style: italic;
  font-weight: bold;
  src: local("Open Sans Bold Italic"), local("OpenSans-BoldItalic"), url("./github/open-sans-v17-latin-ext_latin-700italic.woff2") format("woff2");
  unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+2000-206F, U+2074, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD, U+0100-024F, U+0259, U+1E00-1EFF, U+2020, U+20A0-20AB, U+20AD-20CF, U+2113, U+2C60-2C7F, U+A720-A7FF;
}

html {
  font-size: 16px;
  -webkit-font-smoothing: antialiased;
}

body {
  margin: 0;
  background: #ffffff;
}

.markdown-body {
  max-width: 860px;
  margin: 0 auto;
  padding: 30px;
  padding-bottom: 100px;
  font-family: "Open Sans", "Clear Sans", "Helvetica Neue", Helvetica, Arial, "Segoe UI Emoji", sans-serif;
  color: rgb(51, 51, 51);
  line-height: 1.6;
}

.markdown-body > ul:first-child,
.markdown-body > ol:first-child {
  margin-top: 30px;
}

.markdown-body a {
  color: #4183c4;
}

.markdown-body h1,
.markdown-body h2,
.markdown-body h3,
.markdown-body h4,
.markdown-body h5,
.markdown-body h6 {
  position: relative;
  margin-top: 1rem;
  margin-bottom: 1rem;
  padding: 0;
  font-family: inherit;
  font-weight: bold;
  line-height: 1.4;
}

.markdown-body h1 {
  font-size: 2.25em;
  line-height: 1.2;
  border-bottom: 1px solid #eee;
}

.markdown-body h2 {
  font-size: 1.75em;
  line-height: 1.225;
  border-bottom: 1px solid #eee;
}

.markdown-body h3 {
  font-size: 1.5em;
  line-height: 1.43;
}

.markdown-body h4 {
  font-size: 1.25em;
}

.markdown-body h5 {
  font-size: 1em;
}

.markdown-body h6 {
  font-size: 1em;
  color: #777;
}

.markdown-body h1 tt,
.markdown-body h1 code,
.markdown-body h2 tt,
.markdown-body h2 code,
.markdown-body h3 tt,
.markdown-body h3 code,
.markdown-body h4 tt,
.markdown-body h4 code,
.markdown-body h5 tt,
.markdown-body h5 code,
.markdown-body h6 tt,
.markdown-body h6 code {
  font-size: inherit;
}

.markdown-body p,
.markdown-body blockquote,
.markdown-body ul,
.markdown-body ol,
.markdown-body dl,
.markdown-body table {
  margin: 0.8em 0;
}

.markdown-body li > ol,
.markdown-body li > ul {
  margin: 0;
}

.markdown-body ul,
.markdown-body ol {
  padding-left: 30px;
}

.markdown-body ul:first-child,
.markdown-body ol:first-child {
  margin-top: 0;
}

.markdown-body ul:last-child,
.markdown-body ol:last-child {
  margin-bottom: 0;
}

.markdown-body li p.first {
  display: inline-block;
}

.markdown-body hr {
  height: 2px;
  padding: 0;
  margin: 16px 0;
  background-color: #e7e7e7;
  border: 0 none;
  overflow: hidden;
  box-sizing: content-box;
}

.markdown-body blockquote {
  border-left: 4px solid #dfe2e5;
  padding: 0 15px;
  color: #777777;
}

.markdown-body blockquote blockquote {
  padding-right: 0;
}

.markdown-body table {
  width: 100%;
  padding: 0;
  border-collapse: collapse;
  word-break: initial;
}

.markdown-body table tr {
  border: 1px solid #dfe2e5;
  margin: 0;
  padding: 0;
}

.markdown-body table tr:nth-child(2n), .markdown-body thead {
  background-color: #f8f8f8;
}

.markdown-body table th {
  font-weight: bold;
  border: 1px solid #dfe2e5;
  border-bottom: 0;
  margin: 0;
  padding: 6px 13px;
}

.markdown-body table td {
  border: 1px solid #dfe2e5;
  margin: 0;
  padding: 6px 13px;
}

.markdown-body table th:first-child,
.markdown-body table td:first-child {
  margin-top: 0;
}

.markdown-body table th:last-child,
.markdown-body table td:last-child {
  margin-bottom: 0;
}

.markdown-body code, .markdown-body tt, .markdown-body pre {
  border: 1px solid #e7eaed;
  background-color: #f8f8f8;
  border-radius: 3px;
  padding: 2px 4px 0 4px;
  color: inherit;
  font-size: 0.9em;
}

.markdown-body code {
  background-color: #f3f4f4;
  padding: 0 2px;
}

.markdown-body pre {
  margin-top: 15px;
  margin-bottom: 15px;
  padding-top: 8px;
  padding-bottom: 6px;
  background-color: #f8f8f8;
  overflow-x: auto;
}

.markdown-body pre code {
  border: 0;
  background: transparent;
  padding: 0;
}

.markdown-body img {
  max-width: 100%;
  height: auto;
}

.markdown-body .md-task-list-item > input {
  margin-left: -1.3em;
}

.markdown-body pre.md-meta-block {
  padding: 1rem;
  font-size: 85%;
  line-height: 1.45;
  background-color: #f7f7f7;
  border: 0;
  border-radius: 3px;
  color: #777777;
  margin-top: 0;
}

.markdown-body .md-mathjax-midline {
  background: #fafafa;
}

.markdown-body .md-image > .md-meta {
  border-radius: 3px;
  padding: 2px 0 0 4px;
  font-size: 0.9em;
  color: inherit;
}

.markdown-body .md-tag {
  color: #a7a7a7;
  opacity: 1;
}

.markdown-body .md-lang {
  color: #b4654d;
}

@media print {
  html {
    font-size: 13px;
  }

  pre {
    page-break-inside: avoid;
    word-wrap: break-word;
  }
}
`.trim()

export function buildMarkdownExportHtml({ title, bodyHtml }: { title: string, bodyHtml: string }) {
  // Sanitize the editor's HTML before embedding it in the export. Milkdown
  // lets raw HTML and image/link URLs pass through unchanged, so an attacker
  // who can plant a `<script>`, an `onerror=` handler, or a `javascript:`
  // URL in the source markdown turns the exported file into an XSS payload
  // the moment the victim opens or shares it. The sanitizer runs on every
  // call so the contract is "give me raw body HTML, get back a safe export".
  const safeBody = sanitizeBodyHtml(bodyHtml)
  return `<!DOCTYPE html>
<html><head>
<meta charset="utf-8">
<meta http-equiv="Content-Security-Policy" content="${EXPORT_CSP}">
<title>${escapeHtml(title)}</title>
<style>${typoraGithubExportCss}</style>
</head><body><main class="markdown-body">${safeBody}</main></body></html>`
}

// `script-src 'none'` is the load-bearing directive: even if a `<script>`
// tag slipped past the sanitizer, the browser would refuse to execute it.
// The export inlines its own <style>, so style-src needs 'unsafe-inline'.
// `default-src 'none'` denies everything else (connect, frame, media, …).
const EXPORT_CSP =
  "default-src 'none'; script-src 'none'; style-src 'unsafe-inline'; " +
  "img-src http: https: data:; font-src http: https: data:; " +
  "base-uri 'none'; form-action 'none'"

function escapeHtml(value: string) {
  return value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

// ── Body HTML sanitizer ────────────────────────────────────────────────────
//
// Milkdown produces a ProseMirror DOM subtree that may contain raw HTML the
// author typed inline (CommonMark allows it), plus image srcs, link hrefs,
// and code-block class names. None of it can be trusted: an attacker who
// can plant a `<script>` or a `javascript:` URL in the source markdown wins
// RCE on every machine that opens the export. The sanitizer walks the DOM
// with an allowlist and strips anything not on it. Keep the allowlist tight.

const ALLOWED_TAGS = new Set([
  'a', 'abbr', 'aside', 'b', 'blockquote', 'br', 'caption', 'cite', 'code',
  'dd', 'del', 'details', 'dfn', 'div', 'dl', 'dt', 'em', 'figcaption',
  'figure', 'footer', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'header', 'hr',
  'i', 'img', 'ins', 'kbd', 'li', 'main', 'mark', 'nav', 'ol', 'p', 'pre',
  'q', 's', 'samp', 'section', 'small', 'span', 'strong', 'sub', 'summary',
  'sup', 'table', 'tbody', 'td', 'tfoot', 'th', 'thead', 'time', 'tr',
  'u', 'ul', 'var',
  // Task-list checkbox — only allowed in this constrained form. The HTML
  // spec would normally require an enclosing <form> for an input to render;
  // browsers accept it standalone and Typora's CSS specifically targets it.
  'input',
])

// Tags removed together with their entire subtree. Anything in here can
// either execute script (script, iframe, object) or smuggle it (svg, math).
const DROP_WITH_SUBTREE = new Set([
  'script', 'style', 'noscript', 'iframe', 'object', 'embed', 'frame',
  'frameset', 'noframes', 'applet', 'base', 'form', 'button', 'select',
  'textarea', 'link', 'meta', 'svg', 'math', 'video', 'audio', 'source',
  'track', 'picture', 'template',
])

// Per-tag attribute allowlist. Anything not listed is dropped. This stops
// `on*` event handlers (key for XSS) plus the various HTML/CSS-injection
// vectors (style, srcdoc, …) by omission.
const ALLOWED_ATTRS: Record<string, Set<string>> = {
  a: new Set(['href', 'title', 'name', 'id', 'rel']),
  img: new Set(['src', 'alt', 'title', 'width', 'height']),
  th: new Set(['colspan', 'rowspan', 'scope', 'id', 'class']),
  td: new Set(['colspan', 'rowspan', 'id', 'class']),
  ol: new Set(['start', 'type', 'id', 'class']),
  li: new Set(['value', 'id', 'class']),
  code: new Set(['id', 'class']),
  pre: new Set(['id', 'class']),
  input: new Set(['type', 'checked', 'disabled']),
  '*': new Set(['id', 'class']),
}

// URL schemes we allow in href/src. Anything else is stripped to a
// harmless text node. `data:` is restricted to images below.
const SAFE_URL_SCHEMES = /^(https?:|mailto:|tel:|\/|\.\/|\.\.|[A-Za-z0-9_\-./?#=&%@:+,;~!'])/

export function sanitizeBodyHtml(rawHtml: string): string {
  // The sanitizer only runs in environments that expose DOMParser. When
  // it's not available (older webviews, certain SSR contexts) we fall back
  // to the original html-escape function, which neutralises every tag the
  // browser would otherwise interpret.
  if (typeof DOMParser === 'undefined') {
    return escapeHtml(rawHtml)
  }

  const doc = new DOMParser().parseFromString(`<body>${rawHtml}</body>`, 'text/html')
  walkAndSanitize(doc.body)
  return doc.body.innerHTML
}

function walkAndSanitize(root: Element) {
  // Collect children first — we'll mutate the tree in place.
  const children = Array.from(root.children)
  for (const child of children) {
    const tag = child.tagName.toLowerCase()

    if (DROP_WITH_SUBTREE.has(tag)) {
      child.remove()
      continue
    }

    if (!ALLOWED_TAGS.has(tag)) {
      // Unknown tag — replace with its text content (unwraps the element)
      // so a hostile `<unknown onerror=…>` doesn't sneak through. Children
      // are kept; they get re-evaluated on the next pass via the parent.
      unwrap(child)
      continue
    }

    sanitizeAttrs(child)
    walkAndSanitize(child)
  }
}

function unwrap(el: Element) {
  const parent = el.parentNode
  if (!parent) return
  while (el.firstChild) {
    parent.insertBefore(el.firstChild, el)
  }
  parent.removeChild(el)
}

function sanitizeAttrs(el: Element) {
  const tag = el.tagName.toLowerCase()
  const allowed = new Set([
    ...(ALLOWED_ATTRS[tag] ?? []),
    ...ALLOWED_ATTRS['*'],
  ])

  // Iterate over a snapshot — removeAttribute mutates `attributes`.
  const attrs = Array.from(el.attributes)
  for (const attr of attrs) {
    const name = attr.name.toLowerCase()
    if (!allowed.has(name)) {
      el.removeAttribute(attr.name)
      continue
    }

    if (name === 'href' || name === 'src') {
      // Reject dangerous URL schemes. Stripping the attribute (rather than
      // the whole element) keeps the surrounding anchor/image readable.
      const value = attr.value.trim()
      if (name === 'src' && /^data:image\//i.test(value)) {
        // Inline data: images are common in markdown — allow.
        continue
      }
      if (!isSafeUrl(value)) {
        el.removeAttribute(attr.name)
      }
    }
  }

  // input is a special case: keep it ONLY if it is a checkbox (task list).
  // Any other type is a form-control smuggling vector.
  if (tag === 'input') {
    const type = (el.getAttribute('type') ?? '').toLowerCase()
    if (type !== 'checkbox') {
      unwrap(el)
    }
  }
}

function isSafeUrl(value: string): boolean {
  // Strip control chars and whitespace tricks (`java\tscript:`) before
  // inspecting the scheme.
  const cleaned = value.replace(/[ -\s]/g, '').toLowerCase()
  if (cleaned.startsWith('javascript:') || cleaned.startsWith('vbscript:')) {
    return false
  }
  if (cleaned.startsWith('data:')) {
    // Only data:image is permitted (handled by the caller for src). Any
    // other data: (e.g. data:text/html) is a script-delivery primitive.
    return false
  }
  return SAFE_URL_SCHEMES.test(cleaned)
}
