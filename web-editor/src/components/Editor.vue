<template>
  <div ref="editorRef" class="editor-container" :class="{ 'hide-line-numbers': !showLineNumbers }" @click="handleClick" />
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { Crepe, CrepeFeature } from '@milkdown/crepe'
import { refractor } from 'refractor'
import markdown from 'refractor/lang/markdown'
import javascript from 'refractor/lang/javascript'
import typescript from 'refractor/lang/typescript'
import jsx from 'refractor/lang/jsx'
import tsx from 'refractor/lang/tsx'
import css from 'refractor/lang/css'
import python from 'refractor/lang/python'
import bash from 'refractor/lang/bash'
import json from 'refractor/lang/json'
import yaml from 'refractor/lang/yaml'
import markup from 'refractor/lang/markup'
import go from 'refractor/lang/go'
import rust from 'refractor/lang/rust'
import java from 'refractor/lang/java'
import cpp from 'refractor/lang/cpp'
import c from 'refractor/lang/c'
import sql from 'refractor/lang/sql'
import docker from 'refractor/lang/docker'
import nginx from 'refractor/lang/nginx'
import { prism, prismConfig } from '@milkdown/plugin-prism'
import {
  addBlockTypeCommand,
  clearTextInCurrentBlockCommand,
  codeBlockSchema,
  createCodeBlockCommand,
  downgradeHeadingCommand,
  insertImageCommand,
  paragraphSchema,
  selectTextNearPosCommand,
  toggleEmphasisCommand,
  toggleInlineCodeCommand,
  toggleStrongCommand,
  turnIntoTextCommand,
  wrapInBlockquoteCommand,
  wrapInBulletListCommand,
  wrapInHeadingCommand,
  wrapInOrderedListCommand,
  liftListItemCommand,
  sinkListItemCommand,
} from '@milkdown/kit/preset/commonmark'
import {
  insertTableCommand,
  selectColCommand,
  selectRowCommand,
  selectTableCommand,
  toggleStrikethroughCommand,
} from '@milkdown/kit/preset/gfm'
import { imageBlockSchema } from '@milkdown/kit/component/image-block'
import { toggleLinkCommand } from '@milkdown/kit/component/link-tooltip'
import { $shortcut, replaceAll } from '@milkdown/kit/utils'
import { commandsCtx, editorViewCtx, type CmdKey } from '@milkdown/kit/core'
import type { Ctx } from '@milkdown/kit/ctx'
import { AllSelection, TextSelection, type Command } from '@milkdown/kit/prose/state'
import { deleteRow, isInTable, selectedRect } from '@milkdown/kit/prose/tables'
import { applyEditorSettingsVariables } from '../composables/useEditorSettings'
import '@milkdown/crepe/theme/common/style.css'
import '@milkdown/crepe/theme/frame.css'
import 'prismjs/themes/prism-tomorrow.css'

refractor.register(markdown)
refractor.register(javascript)
refractor.register(typescript)
refractor.register(jsx)
refractor.register(tsx)
refractor.register(css)
refractor.register(python)
refractor.register(bash)
refractor.register(json)
refractor.register(yaml)
refractor.register(markup)
refractor.register(go)
refractor.register(rust)
refractor.register(java)
refractor.register(cpp)
refractor.register(c)
refractor.register(sql)
refractor.register(docker)
refractor.register(nginx)

const STORAGE_KEY = 'fast-md-settings'

function getShowLineNumbers(): boolean {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) {
      const settings = JSON.parse(raw)
      return settings.showLineNumbers ?? false
    }
  } catch { /* ignore */ }
  return false
}

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const editorRef = ref<HTMLElement | null>(null)
const showLineNumbers = ref(getShowLineNumbers())
let crepe: Crepe | null = null
let suppressEmitUntil = 0

// Listen for settings change events
function handleSettingsChange() {
  showLineNumbers.value = getShowLineNumbers()
  applyEditorSettingsVariables()
}

function runCommand<T>(slice: CmdKey<T>, payload?: T): (ctx: Ctx) => Command {
  return (ctx) => () => {
    return ctx.get(commandsCtx).call(slice, payload)
  }
}

function handledNoopCommand(): Command {
  return () => true
}

function clearFormatCommand(ctx: Ctx): Command {
  return (state, dispatch) => {
    if (!dispatch) return true

    const paragraph = paragraphSchema.type(ctx)
    const tr = state.tr
    const { from, to, empty } = state.selection
    const isAllSelected = state.selection instanceof AllSelection
    const docFrom = isAllSelected ? 0 : from
    const docTo = isAllSelected ? state.doc.content.size : to

    try {
      tr.setBlockType(docFrom, docTo, paragraph)
    } catch {
      // Some block selections, such as list containers, cannot be flattened this way.
    }

    Object.values(state.schema.marks).forEach((markType) => {
      tr.removeMark(docFrom, docTo, markType)
      if (empty) tr.removeStoredMark(markType)
    })

    if (tr.docChanged || tr.storedMarksSet) dispatch(tr.scrollIntoView())
    return true
  }
}

function selectCurrentBlockCommand(): Command {
  return (state, dispatch) => {
    if (!dispatch) return true

    if (isInTable(state)) {
      const selected = crepe?.editor.action((ctx) => {
        ctx.get(commandsCtx).call(selectRowCommand.key, { index: getCurrentTableRowIndex(state) })
        return true
      }) ?? false
      return selected
    }

    const { $from } = state.selection
    const from = $from.start()
    const to = $from.end()
    dispatch(state.tr.setSelection(TextSelection.create(state.doc, from, to)).scrollIntoView())
    return true
  }
}

function selectStyleScopeCommand(): Command {
  return (state, dispatch) => {
    if (!dispatch) return true

    if (isInTable(state)) {
      const selected = crepe?.editor.action((ctx) => {
        ctx.get(commandsCtx).call(selectColCommand.key, { index: getCurrentTableColIndex(state) })
        return true
      }) ?? false
      return selected
    }

    const { $from } = state.selection
    const mark = $from.marks().at(-1)
    if (!mark) return selectCurrentBlockCommand()(state, dispatch)

    let from = $from.pos
    let to = $from.pos
    const blockStart = $from.start()

    $from.parent.forEach((node, offset) => {
      if (!mark.isInSet(node.marks)) return

      const start = blockStart + offset
      const end = start + node.nodeSize
      if (start <= $from.pos && $from.pos <= end) {
        from = start
        to = end
      }
    })

    dispatch(state.tr.setSelection(TextSelection.create(state.doc, from, to)).scrollIntoView())
    return true
  }
}

function getCurrentTableRowIndex(state: Parameters<Command>[0]): number {
  return selectedRect(state).top
}

function getCurrentTableColIndex(state: Parameters<Command>[0]): number {
  return selectedRect(state).left
}

function selectCurrentTableCommand(ctx: Ctx): Command {
  return (state, dispatch) => {
    if (!isInTable(state)) return false
    if (!dispatch) return true
    ctx.get(commandsCtx).call(selectTableCommand.key)
    return true
  }
}

function deleteCurrentTableRowCommand(): Command {
  return (state, dispatch) => deleteRow(state, dispatch)
}

function increaseHeadingCommand(ctx: Ctx): Command {
  return (state) => {
    const level = state.selection.$from.parent.type.name === 'heading'
      ? Math.min(Number(state.selection.$from.parent.attrs.level) + 1, 6)
      : 1
    return ctx.get(commandsCtx).call(wrapInHeadingCommand.key, level)
  }
}

function insertBlockImageCommand(ctx: Ctx): Command {
  return () => {
    const commands = ctx.get(commandsCtx)
    const imageBlock = imageBlockSchema.type(ctx)
    return commands.call(addBlockTypeCommand.key, { nodeType: imageBlock })
      || commands.call(insertImageCommand.key, { src: '' })
  }
}

function insertDefaultTableCommand(ctx: Ctx): Command {
  return () => ctx.get(commandsCtx).call(insertTableCommand.key, { row: 3, col: 3 })
}

function insertMathBlockCommand(ctx: Ctx): Command {
  return () => {
    const commands = ctx.get(commandsCtx)
    const view = ctx.get(editorViewCtx)
    const { from } = view.state.selection

    commands.call(clearTextInCurrentBlockCommand.key)
    const inserted = commands.call(addBlockTypeCommand.key, {
      nodeType: codeBlockSchema.type(ctx),
      attrs: { language: 'LaTeX' },
    })

    if (inserted) {
      commands.call(selectTextNearPosCommand.key, { pos: from })
      return true
    }

    return commands.call(createCodeBlockCommand.key, 'LaTeX')
  }
}

const typoraMacOSKeymap = $shortcut(() => ({
  'Mod-0': { key: 'Mod-0', onRun: runCommand(turnIntoTextCommand.key), priority: 100 },
  'Mod-1': { key: 'Mod-1', onRun: runCommand(wrapInHeadingCommand.key, 1), priority: 100 },
  'Mod-2': { key: 'Mod-2', onRun: runCommand(wrapInHeadingCommand.key, 2), priority: 100 },
  'Mod-3': { key: 'Mod-3', onRun: runCommand(wrapInHeadingCommand.key, 3), priority: 100 },
  'Mod-4': { key: 'Mod-4', onRun: runCommand(wrapInHeadingCommand.key, 4), priority: 100 },
  'Mod-5': { key: 'Mod-5', onRun: runCommand(wrapInHeadingCommand.key, 5), priority: 100 },
  'Mod-6': { key: 'Mod-6', onRun: runCommand(wrapInHeadingCommand.key, 6), priority: 100 },
  'Mod-Alt-q': { key: 'Mod-Alt-q', onRun: runCommand(wrapInBlockquoteCommand.key), priority: 100 },
  'Mod-Alt-t': { key: 'Mod-Alt-t', onRun: insertDefaultTableCommand, priority: 100 },
  'Ctrl-Mod-i': { key: 'Ctrl-Mod-i', onRun: insertBlockImageCommand, priority: 100 },
  'Mod-Alt-b': { key: 'Mod-Alt-b', onRun: insertMathBlockCommand, priority: 100 },
  'Mod-Alt-o': { key: 'Mod-Alt-o', onRun: runCommand(wrapInOrderedListCommand.key), priority: 100 },
  'Mod-Alt-u': { key: 'Mod-Alt-u', onRun: runCommand(wrapInBulletListCommand.key), priority: 100 },
  'Mod-l': { key: 'Mod-l', onRun: selectCurrentBlockCommand, priority: 100 },
  'Mod-e': { key: 'Mod-e', onRun: selectStyleScopeCommand, priority: 100 },
  'Mod-a': { key: 'Mod-a', onRun: selectCurrentTableCommand, priority: 100 },
  'Mod-Shift-Backspace': { key: 'Mod-Shift-Backspace', onRun: deleteCurrentTableRowCommand, priority: 100 },
  'Mod-b': { key: 'Mod-b', onRun: runCommand(toggleStrongCommand.key), priority: 100 },
  'Mod-i': { key: 'Mod-i', onRun: runCommand(toggleEmphasisCommand.key), priority: 100 },
  'Shift-Mod-`': { key: 'Shift-Mod-`', onRun: runCommand(toggleInlineCodeCommand.key), priority: 100 },
  'Shift-Ctrl-`': { key: 'Shift-Ctrl-`', onRun: runCommand(toggleStrikethroughCommand.key), priority: 100 },
  'Mod-Alt-c': { key: 'Mod-Alt-c', onRun: runCommand(createCodeBlockCommand.key), priority: 100 },
  'Mod-k': { key: 'Mod-k', onRun: runCommand(toggleLinkCommand.key), priority: 100 },
  'Mod-\\': { key: 'Mod-\\', onRun: clearFormatCommand, priority: 100 },
  'Mod-=': { key: 'Mod-=', onRun: increaseHeadingCommand, priority: 100 },
  'Mod--': { key: 'Mod--', onRun: runCommand(downgradeHeadingCommand.key), priority: 100 },
  'Mod-[': { key: 'Mod-[', onRun: runCommand(sinkListItemCommand.key), priority: 110 },
  'Mod-]': { key: 'Mod-]', onRun: runCommand(liftListItemCommand.key), priority: 110 },
  'Mod-Alt-0': { key: 'Mod-Alt-0', onRun: handledNoopCommand, priority: 100 },
  'Mod-Alt-1': { key: 'Mod-Alt-1', onRun: handledNoopCommand, priority: 100 },
  'Mod-Alt-2': { key: 'Mod-Alt-2', onRun: handledNoopCommand, priority: 100 },
  'Mod-Alt-3': { key: 'Mod-Alt-3', onRun: handledNoopCommand, priority: 100 },
  'Mod-Alt-4': { key: 'Mod-Alt-4', onRun: handledNoopCommand, priority: 100 },
  'Mod-Alt-5': { key: 'Mod-Alt-5', onRun: handledNoopCommand, priority: 100 },
  'Mod-Alt-6': { key: 'Mod-Alt-6', onRun: handledNoopCommand, priority: 100 },
  'Mod-Alt-7': { key: 'Mod-Alt-7', onRun: handledNoopCommand, priority: 100 },
  'Mod-Alt-8': { key: 'Mod-Alt-8', onRun: handledNoopCommand, priority: 100 },
  'Mod-Alt-x': { key: 'Mod-Alt-x', onRun: handledNoopCommand, priority: 100 },
  'Mod-Shift-b': { key: 'Mod-Shift-b', onRun: handledNoopCommand, priority: 100 },
}))

onMounted(async () => {
  applyEditorSettingsVariables()
  window.addEventListener('fast-md-settings-changed', handleSettingsChange)

  if (!editorRef.value) return

  crepe = new Crepe({
    root: editorRef.value,
    defaultValue: props.modelValue,
    features: {
      [CrepeFeature.CodeMirror]: true,
      [CrepeFeature.Latex]: false,
      [CrepeFeature.Toolbar]: true,
      [CrepeFeature.BlockEdit]: true,
      [CrepeFeature.Table]: true,
      [CrepeFeature.Cursor]: true,
      [CrepeFeature.ListItem]: true,
      [CrepeFeature.LinkTooltip]: true,
      [CrepeFeature.Placeholder]: true,
    },
  })

  crepe.editor
    .use(typoraMacOSKeymap)
    .use(prism)
    .config((ctx) => {
      ctx.set(prismConfig.key, { configureRefractor: () => refractor })
    })

  crepe.on((api) => {
    api.markdownUpdated((_ctx, markdown) => {
      if (Date.now() >= suppressEmitUntil) {
        emit('update:modelValue', markdown)
      }
    })
  })

  await crepe.create()

  // Auto-focus the editor on mount (Typora-like: ready to type immediately)
  requestAnimationFrame(() => {
    crepe?.editor.action((ctx) => {
      ctx.get(editorViewCtx).focus()
    })
  })
})

function handleClick(e: MouseEvent) {
  if (!crepe) return
  const target = e.target as HTMLElement

  crepe.editor.action((ctx) => {
    const view = ctx.get(editorViewCtx)

    // Click on a content node inside ProseMirror → let ProseMirror handle natively
    if (target !== view.dom && target.closest('.ProseMirror')) return

    // Click landed on empty space (ProseMirror container itself or outside it entirely).
    // Sublime Text behaviour: keep the cursor exactly where it is, just re-focus.
    e.preventDefault()
    view.focus()
  })
}

watch(
  () => props.modelValue,
  (newVal) => {
    if (!crepe) return
    const current = crepe.getMarkdown()
    if (current === newVal) return
    suppressEmitUntil = Date.now() + 300
    crepe.editor.action(replaceAll(newVal))
  },
)

onUnmounted(() => {
  crepe?.destroy()
  crepe = null
  window.removeEventListener('fast-md-settings-changed', handleSettingsChange)
})
</script>

<style scoped>
.editor-container {
  flex: 1;
  min-width: 0;
  overflow-y: auto;
  height: 100%;
  background: var(--bg-primary);
  cursor: text;
}

:deep(.milkdown .ProseMirror) {
  font-family: var(--editor-font-family) !important;
  font-size: var(--editor-font-size, 16px) !important;
}

:deep(.milkdown .ProseMirror p) {
  font-size: var(--editor-font-size, 16px) !important;
}

:deep(.milkdown .ProseMirror pre) {
  background: #2d2d2d;
  border-radius: 6px;
  padding: 16px;
  overflow-x: auto;
}

:deep(.milkdown .ProseMirror pre code) {
  font-family: var(--code-font-family, 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace);
  font-size: 0.9em;
  background: transparent;
  padding: 0;
  color: #ccc;
}

:deep(.milkdown .ProseMirror .code-block) {
  position: relative;
}

.hide-line-numbers :deep(.cm-lineNumbers),
.hide-line-numbers :deep(.cm-gutters) {
  display: none !important;
}
</style>
