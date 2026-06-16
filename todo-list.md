# fast-md 项目问题清单 (Top 5)

> 基于源码逐行验证,按严重性从高到低排序。
> 生成时间:2026-06-16

---

## 1. [Critical] Quit 流程实际上不会触发 `app.Quit()`

**位置:**
- `core/app.go:173-180`
- `core/quit.go:43-62`
- `frontend/src/App.vue:113-156`

**问题:** Go 端从 `ctx.Value(application.WindowKey)` 取窗口,但前端通过 `ByID` 绑定调用,`ctx` 是否携带 `WindowKey` 取决于 Wails 运行时行为;一旦取不到,`Confirm` 不会被调用 → 队列不前进 → `app.Quit()` 从不执行。Save 路径还有第二处隐患:`App.vue:134-137` 的 `await saveFile()` 未捕获错误,若 `WriteFile` 失败,`isDirty` 保持 true,`executeClose` 被跳过,错误被静默吞掉。

**影响:** Cmd-Q 应用不退出;或看似保存成功实则没保存就关窗。

- [ ] 验证 Wails 3 alpha `ByID` 绑定是否注入 `application.WindowKey` 到 ctx
- [ ] 若不注入,改 `ConfirmQuitWindow()` 签名为接收 `windowID uint`
- [ ] `App.vue:134-137` 的 `await saveFile()` 改为 try/catch 并提示用户
- [ ] 补 quit 流程的 end-to-end 测试(目前完全无覆盖)

---

## 2. [High] HTML/PDF 导出未净化,任意 XSS 落入导出文件

**位置:**
- `frontend/src/App.vue:177-205`
- `frontend/src/exportHtml.ts:298-304`
- `core/app.go:235-267`

**问题:** `escapeHtml` 只对 `title` 用了,`bodyHtml` 是 `.ProseMirror.innerHTML` 原样拼接。Milkdown 允许 raw HTML、图片 `src`、链接 `href` 透传 — 包含 `<script>` / `<img src="javascript:…">` / `<svg onload=…>` 的 markdown 导出后变 RCE payload。导出文件无 CSP。`ExportPDF` 在 `app.go:255` 把同一段 HTML 喂给 folio webview,同样的 sink。

**影响:** 一次粘贴即可获得受害者机器 RCE;导出文件会外发,放大攻击面。

- [ ] 在 `buildMarkdownExportHtml` 中对 `bodyHtml` 做白名单标签/属性过滤
- [ ] 给导出 HTML 加 `<meta http-equiv="Content-Security-Policy" content="default-src 'none'; img-src https: data:; style-src 'unsafe-inline'">`
- [ ] `insertImageCommand` 等 Milkdown 命令处过滤 `javascript:` / `data:` URL
- [ ] 给 `exportHtml.ts` 加测试:含 `<script>` / `onerror=` / `javascript:` 的输入必须被中和

---

## 3. [High] `ListDirectory` 接受任意路径,无沙箱

**位置:** `core/app.go:119-137`

**问题:** `os.ReadDir(path)` 路径完全无校验,后端只对**返回条目**按后缀过滤,目录本身不限制。`WriteFile` 同样裸跑(`app.go:115-117`)。结合 #2 的 XSS,一次成功渲染即可枚举整盘 / 改写可写文件。

**影响:** 当前威胁面小,但**安全姿态为 0**;任何新功能(预览、插件、URL 协议处理器)立即成为完全 FS read/write primitive。

- [ ] `ListDirectory` 增加允许根目录白名单(如用户 home、用户显式打开过的目录)
- [ ] `WriteFile` 同样做范围检查
- [ ] 至少记录"用户打开了非常规目录"的事件日志,作为安全审计基线
- [ ] 评估是否在 Go 端引入最小 chroot 视图

---

## 4. [Medium] `WriteFile` 错误被静默吞掉,Save 状态撒谎

**位置:** `frontend/src/composables/useFile.ts:95-137`

**问题:** `saveFile()` / `saveAs()` / `saveToPath()` 三个方法共用 `try/finally` 但不 `catch`、不 `rethrow`、不通知 UI。`WriteFile` 失败后 `isSaving` 复位,`isDirty` 保持 true,`StatusBar` 显示"未保存" — 用户以为刚才在保存,实际没存。配合 #1,Cmd-Q 选"保存"后被静默"取消退出",无任何反馈。

**影响:** 静默数据丢失,触发条件常见(磁盘满、file lock、iCloud 同步冲突),且 UI 撒谎。

- [ ] `saveFile` / `saveAs` / `saveToPath` 改为 `try/catch` 并通过 composable 返回 error
- [ ] `StatusBar.vue` / `App.vue` 在 save 失败时显示明确错误
- [ ] Cmd-Q 流程区分"保存成功"、"保存失败(请重试或另存为)"、"用户取消"
- [ ] `useFile.test.ts` 加测试:mock `WriteFile` reject,断言错误冒泡到 UI

---

## 5. [Medium] Recent files 不持久化 + 写入时机错误 + 跨 goroutine 状态

**位置:**
- `core/dockmenu_darwin.go:50-103`
- `core/app.go:111`

**问题(三层叠加):**
1. **不持久化** — `recentFiles` 是纯内存 slice,从 `loadConfig` / `saveConfig` 不写盘,重启清零。
2. **错误时机** — `dockMenuOpenRecent` 在 `go func()` 启动后**立刻**调用 `trackRecentFile`,没等 `NewEditorWindowWithFile` 验证文件可读。
3. **跨 goroutine 状态** — `trackRecentFile` 与 `updateDockRecentMenu` 共享 `recentFiles`,`&paths[0]` Go 指针跨越 cgo 边界。

**影响:** macOS dock "最近文件" 集成失去价值(recent 总是空);recent 菜单列出打不开的"幽灵条目";Wails/Go 升级时隐藏 footgun。

- [ ] 把 `recentFiles` 落盘到 `~/Library/Application Support/fast-md/recent.json`,启动时加载
- [ ] `dockMenuOpenRecent` 改为:在 `NewEditorWindowWithFile` 成功(窗口 ready)后再 `trackRecentFile`
- [ ] 抽 `RecentFilesStore` 类型封装 mutex + 持久化,让 `trackRecentFile` 和 `updateDockRecentMenu` 只通过它交互
- [ ] 评估 `&paths[0]` 跨 cgo 边界的 Go 指针安全性,必要时改成显式拷贝

---

## 顺带值得关注(不进 Top 5)

- [ ] `useFile.ts:25` — `setContent` 无条件 `isDirty = true`,导致 `App.vue:124` 的 `content.value.trim() !== '' && !filePath.value` 变成死代码,新建空文档意外触发 unsaved 提示
- [ ] `App.vue:251-255` — 监听 `common:WindowFilesDropped` 的 `data.filenames` 形态可能与 Wails 事件 payload 不一致(Go 端 `run.go:60-69` 自己 emit `file:open`),存在双通道冲突
- [ ] `App.vue:202-204` — `ExportPDF` 返回的 `path` 完全被忽略,`console.error` 也不展示给用户;取消保存时 `fmt.Errorf("cancelled")` 弹控制台
- [ ] `core/app.go:213-233` — `ShowSaveDialog` 三个按钮文案("取消/不保存/保存")**写死中文**(`menu_i18n.go` 未走),切换 `SetUILocale('en')` 后该对话框仍是中文
- [ ] 测试覆盖缺口:`core/dockmenu_darwin.go`、`core/menu_i18n.go`、`core/savedialog_darwin.go` 无 Go 测试;`Editor.vue` / `Sidebar.vue` / `StatusBar.vue` / `Settings.vue` 无组件测试;**quit 流程 end-to-end 无测试**
