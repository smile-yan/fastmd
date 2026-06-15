import { ref, computed } from 'vue'

export type Locale = 'zh' | 'en'

const LOCALE_KEY = 'fast-md-locale'

const translations = {
  zh: {
    // App
    untitled: '未命名',
    edited: '已编辑',

    // Menu
    menu: {
      file: '文件',
      new: '新建',
      open: '打开...',
      save: '保存',
      saveAs: '另存为...',
      exportHtml: '导出 HTML',
      exportPdf: '导出 PDF',
      edit: '编辑',
      undo: '撤销',
      redo: '重做',
      cut: '剪切',
      copy: '复制',
      paste: '粘贴',
      selectAll: '全选',
      view: '视图',
      toggleSidebar: '切换侧边栏',
      toggleSourceMode: '切换源码模式',
      toggleFullscreen: '切换全屏',
      theme: '主题',
      toggleTheme: '切换主题',
      settings: '设置',
      close: '关闭',
      help: '帮助',
      about: '关于',
    },

    // File
    file: {
      saved: '已保存',
      saving: '保存中...',
      unsavedChanges: '未保存的更改',
    },

    // Settings
    settings: {
      title: '设置',
      general: '通用',
      editor: '编辑器',
      image: '图像',
      appearance: '外观',

      // General
      language: '语言',
      update: '更新',
      autoUpdate: '开启自动更新',
      checkUpdate: '检查更新',
      save: '保存',
      autoSave: '自动保存',
      manualSave: '手动保存',
      autoSaveInterval: '自动保存间隔',
      seconds: '秒',
      confirmOnClose: '关闭前确认保存',
      confirmOnCloseDesc: '关闭窗口或退出时提示保存未保存的更改',

      spellZh: '中文',
      spellEn: 'English',
      on: '开启',
      off: '关闭',

      // Editor
      font: '字体',
      fontSystem: '系统字体',
      fontSerif: 'Serif',
      fontSansSerif: 'Sans Serif',
      fontMonospace: '等宽字体',
      fontSize: '字体大小',
      lineHeight: '行高',
      lineHeightCompact: '紧凑',
      lineHeightNormal: '标准',
      lineHeightRelaxed: '宽松',
      showLineNumbers: '显示行号',
      bracketMatch: '括号匹配',
      tabSpaces: 'Tab 空格数',

      // Image
      imageInsert: '插入图片',
      imageUpload: '上传到服务器',
      imageLocal: '使用本地路径',
      imageCopy: '复制到本地',
      imageWidth: '默认宽度',
      imageDisplay: '图片显示',
      imageAuto: '自动',
      imageBlock: '块级',
      imageInline: '行内',

      exportHtml: 'HTML',
      exportPdf: 'PDF',
      exportDocx: 'Word',
      exportMarkdown: 'Markdown',

      // Appearance
      appTheme: '主题',
      appThemeSystem: '跟随系统',
      appThemeLight: '浅色',
      appThemeDark: '深色',
      contentTheme: '内容主题',
      sidebar: '侧边栏',
      sidebarShow: '启动时显示',
      sidebarHide: '启动时隐藏',
    },

    // Status bar
    statusBar: {
      words: '字',
      characters: '字符',
      lines: '行',
    },

    // Dialogs
    dialog: {
      unsavedTitle: '未保存的更改',
      unsavedMessage: '有未保存的更改，确定要放弃吗？',
      discardUnsavedFile: '"{name}" 有未保存的更改。是否放弃并继续？',
      restartForLanguageChange: '语言已切换为中文，需要重启应用才能生效。是否立即重启？',
      saveFailed: '保存失败，请检查文件权限或磁盘空间后重试。',
      exportFailed: '导出 PDF 失败，请重试或检查目标位置是否可写。',
      save: '保存',
      discard: '放弃',
      cancel: '取消',
    },

    // Sidebar
    sidebar: {
      noFolder: '无文件夹',
      openFolder: '打开文件夹以浏览文件',
      noFiles: '未找到 .md 文件',
    },
  },
  en: {
    // App
    untitled: 'Untitled',
    edited: 'edited',

    // Menu
    menu: {
      file: 'File',
      new: 'New',
      open: 'Open...',
      save: 'Save',
      saveAs: 'Save As...',
      exportHtml: 'Export HTML',
      exportPdf: 'Export PDF',
      edit: 'Edit',
      undo: 'Undo',
      redo: 'Redo',
      cut: 'Cut',
      copy: 'Copy',
      paste: 'Paste',
      selectAll: 'Select All',
      view: 'View',
      toggleSidebar: 'Toggle Sidebar',
      toggleSourceMode: 'Toggle Source Mode',
      toggleFullscreen: 'Toggle Fullscreen',
      theme: 'Theme',
      toggleTheme: 'Toggle Theme',
      settings: 'Settings',
      close: 'Close',
      help: 'Help',
      about: 'About',
    },

    // File
    file: {
      saved: 'Saved',
      saving: 'Saving...',
      unsavedChanges: 'Unsaved changes',
    },

    // Settings
    settings: {
      title: 'Settings',
      general: 'General',
      editor: 'Editor',
      image: 'Image',
      appearance: 'Appearance',

      // General
      language: 'Language',
      update: 'Update',
      autoUpdate: 'Enable Auto Update',
      checkUpdate: 'Check for Updates',
      save: 'Save',
      autoSave: 'Auto Save',
      manualSave: 'Manual Save',
      autoSaveInterval: 'Auto Save Interval',
      seconds: 'seconds',
      confirmOnClose: 'Confirm Before Close',
      confirmOnCloseDesc: 'Prompt to save unsaved changes when closing or quitting',

      spellZh: 'Chinese',
      spellEn: 'English',
      on: 'On',
      off: 'Off',

      // Editor
      font: 'Font',
      fontSystem: 'System Font',
      fontSerif: 'Serif',
      fontSansSerif: 'Sans Serif',
      fontMonospace: 'Monospace',
      fontSize: 'Font Size',
      lineHeight: 'Line Height',
      lineHeightCompact: 'Compact',
      lineHeightNormal: 'Normal',
      lineHeightRelaxed: 'Relaxed',
      showLineNumbers: 'Show Line Numbers',
      bracketMatch: 'Bracket Matching',
      tabSpaces: 'Tab Spaces',

      // Image
      imageInsert: 'Insert Image',
      imageUpload: 'Upload to Server',
      imageLocal: 'Use Local Path',
      imageCopy: 'Copy to Local',
      imageWidth: 'Default Width',
      imageDisplay: 'Image Display',
      imageAuto: 'Auto',
      imageBlock: 'Block',
      imageInline: 'Inline',

      exportHtml: 'HTML',
      exportPdf: 'PDF',
      exportDocx: 'Word',
      exportMarkdown: 'Markdown',

      // Appearance
      appTheme: 'Theme',
      appThemeSystem: 'System',
      appThemeLight: 'Light',
      appThemeDark: 'Dark',
      contentTheme: 'Content Theme',
      sidebar: 'Sidebar',
      sidebarShow: 'Show on Startup',
      sidebarHide: 'Hide on Startup',
    },

    // Status bar
    statusBar: {
      words: 'words',
      characters: 'chars',
      lines: 'lines',
    },

    // Dialogs
    dialog: {
      unsavedTitle: 'Unsaved Changes',
      unsavedMessage: 'You have unsaved changes. Discard them?',
      discardUnsavedFile: '"{name}" has unsaved changes. Discard and continue?',
      restartForLanguageChange: 'Language changed to English. Restart to apply?',
      saveFailed: 'Save failed. Please check the file permissions or disk space and try again.',
      exportFailed: 'PDF export failed. Please try again or check whether the destination is writable.',
      save: 'Save',
      discard: 'Discard',
      cancel: 'Cancel',
    },

    // Sidebar
    sidebar: {
      noFolder: 'No folder',
      openFolder: 'Open a folder to browse files',
      noFiles: 'No .md files found',
    },
  },
}

const currentLocale = ref<Locale>((localStorage.getItem(LOCALE_KEY) as Locale) || 'zh')

export function setLocale(locale: Locale) {
  currentLocale.value = locale
  localStorage.setItem(LOCALE_KEY, locale)
  window.dispatchEvent(new Event('locale-changed'))
}

export function useLocale() {
  const t = (key: string): string => {
    const keys = key.split('.')
    let value: any = translations[currentLocale.value]
    for (const k of keys) {
      value = value?.[k]
    }
    return value ?? key
  }

  const locale = computed(() => currentLocale.value)
  const isZh = computed(() => currentLocale.value === 'zh')
  const isEn = computed(() => currentLocale.value === 'en')

  return {
    t,
    locale,
    isZh,
    isEn,
    setLocale,
  }
}
