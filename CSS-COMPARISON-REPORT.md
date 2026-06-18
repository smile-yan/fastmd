# CSS Comparison Report: Typora Export vs fast-md Export

> 除了代码块部分，目标是一致性

## 测试文件

`test-export-comparison.md` — 涵盖 headings、paragraph、list、blockquote、table、task list、image、link、hr 等元素。

**Typora 导出**: `test-export-comparison.html` (25KB, 含5个 style block)
**fast-md 导出**: `test-export-comparison-fastmd.html` (9KB, 含1个 style block)

---

## Typora 导出的 CSS 架构

Typora 导出包含 **5 层 CSS**（按级联顺序）：

| Block | ID | 内容 | 大小 |
|-------|-----|------|------|
| 0 | (inline) | `html {overflow-x: initial !important}` | 38B |
| 1 | `style-base` | 编辑器的 base 样式 | 16KB |
| 2 | `style-theme_css` | GitHub 主题 CSS | 6.2KB |
| 3 | `style-lp` | `ol, ul {padding-left: 40px}` | 27B |
| 4 | `style-mac-print` | 打印样式 | 117B |

fast-md 只有 **1 层 CSS**：`typoraGithubExportCss`（对应 Typora 的 block 2 主题 CSS）。

**关键发现**：Typora 的 base 样式 + lp 样式会覆盖/补充主题样式，这些规则在 fast-md 中缺失。

---

## 逐元素差异分析

### 1. 页面级别 (html/body)

| 属性 | Typora 最终值 | fast-md 值 | 状态 |
|------|-------------|-----------|------|
| `html font-size` | `16px` (theme 覆盖 base 的 `14px`) | `16px` | ✅ 一致 |
| `html -webkit-font-smoothing` | `antialiased` (base) | `antialiased` | ✅ 一致 |
| `body background` | `var(--bg-color)` = `#ffffff` (base) | `#ffffff` | ✅ 一致 |
| `body font-family` | `"Open Sans","Clear Sans", ...` (theme) | 在 `.markdown-body` 上设置 | ⚠️ 位置不同，效果等价 |
| `body color` | `#333333` (base + theme) | 在 `.markdown-body` 上设置 | ⚠️ 位置不同 |
| `body line-height` | `1.6` (theme 覆盖 base 的 `1.428571`) | `1.6` | ✅ 一致 |
| `body padding` | `30px` (`.typora-export`) | `30px` (`.markdown-body`) | ✅ 一致 |
| **`* { box-sizing: border-box }`** | ✅ (base) | ❌ **缺失** | 🔴 差异 |
| `body.typora-export` padding | `30px` (base) | N/A | ⚠️ fast-md 用 `.markdown-body` 容器 |

### 2. 容器

| 属性 | Typora | fast-md |
|------|--------|---------|
| 容器元素 | `<div id='write'>` | `<main class='markdown-body'>` |
| max-width | `860px` (theme `#write`) | `860px` (`.markdown-body`) |
| margin | `0 auto` (theme) | `0 auto` |
| padding | `30px` + `padding-bottom: 100px` | `30px` + `padding-bottom: 100px` |

容器级别的样式基本一致。✅

### 3. 标题 (h1-h6)

| 属性 | Typora | fast-md | 状态 |
|------|--------|---------|------|
| `h1 font-size` | `2.25em` (theme 覆盖 base `2rem`) | `2.25em` | ✅ |
| `h2 font-size` | `1.75em` (theme 覆盖 base `1.8rem`) | `1.75em` | ✅ |
| `h3 font-size` | `1.5em` (theme 覆盖 base `1.6rem`) | `1.5em` | ✅ |
| `h4 font-size` | `1.25em` (theme 覆盖 base `1.4rem`) | `1.25em` | ✅ |
| `h5 font-size` | `1em` (theme 覆盖 base `1.2rem`) | `1em` | ✅ |
| `h6 font-size` | `1em` | `1em` | ✅ |
| `h1, h2 border-bottom` | `1px solid #eee` | `1px solid #eee` | ✅ |
| `margin-top/bottom` | `1rem` (base) → theme 也设 `1rem` | `1rem` | ✅ |
| `h1-h6 break-inside` | `avoid` (base) | ❌ **缺失** | 🟡 打印相关 |
| `code in headings` | `font-size: inherit` | `font-size: inherit` | ✅ |
| `h1-h6 font-weight` | `bold` | `bold` | ✅ |

### 4. 段落

| 属性 | Typora | fast-md | 状态 |
|------|--------|---------|------|
| `p margin` | `0.8em 0` (theme 覆盖 base `1rem`) | `0.8em 0` | ✅ |
| `p orphans` | `4` (base) | ❌ **缺失** | 🟡 打印相关 |

### 5. 列表 (ul/ol/li)

| 属性 | Typora | fast-md | 状态 |
|------|--------|---------|------|
| **`ul, ol padding-left`** | **`40px`** (style-lp 覆盖 theme `30px`) | `30px` | 🔴 **差异** |
| `li > ol, li > ul margin` | `0` | `0` | ✅ |
| `li p margin` | `0.5rem 0px` (base) | ❌ **缺失** | 🔴 列表项内段落间距不同 |
| `ul:first-child margin-top` | `0` | `0` | ✅ |

### 6. 引用块 (blockquote)

| 属性 | Typora | fast-md | 状态 |
|------|--------|---------|------|
| `border-left` | `4px solid #dfe2e5` | `4px solid #dfe2e5` | ✅ |
| `padding` | `0 15px` | `0 15px` | ✅ |
| `color` | `#777777` | `#777777` | ✅ |
| `margin` | `0.8em 0` (theme 覆盖 base `1rem 0px`) | `0.8em 0` | ✅ |
| **`blockquote > :last-child`** | `margin-bottom: 0px` (base) | ❌ **缺失** | 🟡 |
| **`blockquote > :first-child`** | `margin-top: 0px` (base) | ❌ **缺失** | 🟡 |

### 7. 表格

| 属性 | Typora | fast-md | 状态 |
|------|--------|---------|------|
| `border-collapse` | `collapse` (base) | `collapse` | ✅ |
| `width` | `100%` (base) | `100%` | ✅ |
| `tr border` | `1px solid #dfe2e5` | `1px solid #dfe2e5` | ✅ |
| `th, td padding` | `6px 13px` | `6px 13px` | ✅ |
| `th font-weight` | `bold` | `bold` | ✅ |
| `tr:nth-child(2n) bg` | `#f8f8f8` | `#f8f8f8` | ✅ |
| `thead bg` | `#f8f8f8` (base + theme) | `#f8f8f8` | ✅ |
| `table text-align` | `left` (base) | ❌ **缺失** | 🟡 |
| `tr break-inside` | `avoid` (base) | ❌ **缺失** | 🟡 打印相关 |

### 8. 图片

| 属性 | Typora | fast-md | 状态 |
|------|--------|---------|------|
| `max-width` | `100%` | `100%` | ✅ |
| `height` | `auto` (on `img:not([height])`) | `auto` (always) | ⚠️ 逻辑不同 |
| **`vertical-align`** | **`middle`** (base) | ❌ **缺失** | 🔴 差异 |

### 9. 任务列表 (task list)

这是一个 **显著差异区域**：

| 属性 | Typora (base) | fast-md | 状态 |
|------|--------------|---------|------|
| `list-style-type` | `none` (`.md-task-list-item`) | `none` (通过 `.md-task-list-item > input`) | ⚠️ |
| `position` | `relative` (`.md-task-list-item`) | ❌ 缺失 | 🟡 |
| `padding-left` | `0px` (`.task-list-item.md-task-list-item`) | ❌ 缺失 | 🔴 |
| `input position` | `absolute; top: 0; left: 0` | ❌ 缺失 | 🔴 |
| `input margin-left` | `-1.2em` (base) | `-1.3em` | ⚠️ 微小差异 |
| `input margin-top` | `calc(1em - 10px)` | ❌ 缺失 | 🔴 |
| `input type="checkbox" padding` | `0px` (base) | ❌ 缺失 | 🟡 |
| `input[type="checkbox"] line-height` | `normal` (base) | ❌ 缺失 | 🟡 |

### 10. 水平线 (hr)

| 属性 | Typora | fast-md | 状态 |
|------|--------|---------|------|
| `height` | `2px` | `2px` | ✅ |
| `background-color` | `#e7e7e7` | `#e7e7e7` | ✅ |
| `border` | `0 none` | `0 none` | ✅ |
| `margin` | `16px 0` | `16px 0` | ✅ |
| `box-sizing` | `content-box` | `content-box` | ✅ |

### 11. 行内代码 (inline code)

| 属性 | Typora | fast-md | 状态 |
|------|--------|---------|------|
| `border` | `1px solid #e7eaed` | `1px solid #e7eaed` | ✅ |
| `background-color` | `#f3f4f4` (theme 覆盖 base `#f8f8f8`) | `#f3f4f4` | ✅ |
| `border-radius` | `3px` | `3px` | ✅ |
| `padding` | `0 2px` | `0 2px` | ✅ |
| `font-size` | `0.9em` | `0.9em` | ✅ |
| **`font-family`** | **`var(--monospace)`** (base) | ❌ **缺失** | 🔴 代码字体不同 |

### 12. `<kbd>` 标签

Typora 有完整的 `<kbd>` 样式（border, border-radius, box-shadow, background 等），fast-md **完全缺失**。

---

## 缺失的全局样式

以下样式在 fast-md 导出中完全缺失：

```css
/* 盒模型 — 显著影响布局 */
*, ::after, ::before { box-sizing: border-box; }

/* 列表项段落间距 */
li p { margin: 0.5rem 0; }

/* 图片垂直对齐 */
img { vertical-align: middle; }

/* 表格文本对齐 */
table { text-align: left; }

/* 代码等宽字体 */
code, pre, samp, tt { font-family: var(--monospace); }

/* 任务列表完整样式 */
.md-task-list-item { position: relative; list-style-type: none; }
.task-list-item.md-task-list-item { padding-left: 0; }
.md-task-list-item > input { position: absolute; top: 0; left: 0; margin-top: calc(1em - 10px); }
input[type="checkbox"] { line-height: normal; padding: 0; }

/* kbd 键盘按键样式 */
kbd { ... }

/* mark 高亮样式 */
mark { background: #ff0; color: #000; }

/* 块引用子元素边距 */
blockquote > :first-child { margin-top: 0; }
blockquote > :last-child { margin-bottom: 0; }

/* 列表容器边距 */
li > :first-child { margin-top: 0; }

/* 打印相关 */
h1, h2, h3, h4, h5, h6 { break-inside: avoid; }
thead, tr { break-inside: avoid; }
```

---

## 差异影响评分

| 差异项 | 影响程度 | 说明 |
|--------|---------|------|
| `*, ::after, ::before { box-sizing: border-box }` | 🔴 高 | 影响所有元素的盒模型计算 |
| `ol, ul { padding-left: 40px }` vs `30px` | 🔴 高 | 列表缩进明显不同 |
| `li p { margin: 0.5rem 0 }` | 🔴 高 | 列表项内段落间距不同 |
| `img { vertical-align: middle }` | 🟡 中 | 图片与文字对齐 |
| 代码字体族 `var(--monospace)` | 🟡 中 | 代码显示字体 |
| `table { text-align: left }` | 🟡 中 | 表格内容对齐 |
| 任务列表完整样式 | 🔴 高 | 复选框位置偏移 |
| `blockquote > :first/last-child` | 🟡 中 | 嵌套引用块边距 |
| `li > :first-child { margin-top: 0 }` | 🟡 中 | 列表首元素间距 |
| `h1-h6 { break-inside: avoid }` | 🟢 低 | 仅影响打印 |
| `<kbd>` 样式 | 🟢 低 | 使用频率低 |
| `<mark>` 样式 | 🟡 中 | 高亮文本 |
| `input[type="checkbox"]` 样式 | 🔴 高 | 任务复选框外观 |

---

## 总结

fast-md 的 `typoraGithubExportCss` 是基于 Typora 的 **GitHub 主题 CSS** 创建的，但缺少了 Typora 的 **base 样式**和 **style-lp 样式**。这导致大约 **15 个 CSS 规则差异**，其中约 **8 个是高影响差异**。

建议将缺失的 base 样式整合进 fast-md 的导出 CSS 中（代码块相关样式除外），以达到与 Typora 一致的渲染效果。
