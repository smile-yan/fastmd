package core

import (
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Locale string

const (
	LocaleZh Locale = "zh"
	LocaleEn Locale = "en"
)

var currentLocale = LocaleZh

var menuStrings = map[Locale]menuI18n{
	LocaleZh: {
		appName:               "fast-md",
		about:                 "关于 fast-md",
		preferences:           "偏好设置...",
		quit:                  "退出 fast-md",
		help:                  "帮助",
		file:                  "文件",
		newFile:               "新建文件",
		newWindow:             "新建窗口",
		open:                  "打开...",
		save:                  "保存",
		saveAs:                "另存为...",
		exportHtml:            "导出 HTML...",
		exportPdf:             "导出 PDF...",
		editor:                "编辑",
		view:                  "视图",
		toggleSidebar:         "切换侧边栏",
		theme:                 "主题",
		enterFullscreen:       "进入全屏",
		exitFullscreen:        "退出全屏",
		developerTools:        "开发者工具",
		helpQuickStart:        "快速开始",
		helpKeyboardShortcuts: "快捷键说明",
		helpMarkdownBasics:    "markdown入门",
		helpMathBasics:        "数学公式入门",
		dockNewWindow:         "新窗口",
		dockOpenFile:          "打开文件...",
		dockRecentFiles:       "最近文件",
		unsavedTitle:          "未保存的更改",
		unsavedMessage:        "有未保存的更改，确定要放弃吗？",
		closeWithoutSaving:    "是否在不保存更改的情况下关闭？",
		discard:               "不保存",
		cancel:                "取消",
	},
	LocaleEn: {
		appName:               "fast-md",
		about:                 "About fast-md",
		preferences:           "Preferences...",
		quit:                  "Quit fast-md",
		help:                  "Help",
		file:                  "File",
		newFile:               "New File",
		newWindow:             "New Window",
		open:                  "Open...",
		save:                  "Save",
		saveAs:                "Save As...",
		exportHtml:            "Export as HTML...",
		exportPdf:             "Export as PDF...",
		editor:                "Editor",
		view:                  "View",
		toggleSidebar:         "Toggle Sidebar",
		theme:                 "Theme",
		enterFullscreen:       "Enter Full Screen",
		exitFullscreen:        "Exit Full Screen",
		developerTools:        "Developer Tools",
		helpQuickStart:        "Quick Start",
		helpKeyboardShortcuts: "Keyboard Shortcuts",
		helpMarkdownBasics:    "Markdown Basics",
		helpMathBasics:        "Math Formula Basics",
		dockNewWindow:         "New Window",
		dockOpenFile:          "Open File...",
		dockRecentFiles:       "Recent Files",
		unsavedTitle:          "Unsaved Changes",
		unsavedMessage:        "You have unsaved changes. Discard them?",
		closeWithoutSaving:    "Close without saving changes?",
		discard:               "Discard",
		cancel:                "Cancel",
	},
}

type menuI18n struct {
	appName               string
	about                 string
	preferences           string
	quit                  string
	help                  string
	file                  string
	newFile               string
	newWindow             string
	open                  string
	save                  string
	saveAs                string
	exportHtml            string
	exportPdf             string
	editor                string
	view                  string
	toggleSidebar         string
	theme                 string
	enterFullscreen       string
	exitFullscreen        string
	developerTools        string
	helpQuickStart        string
	helpKeyboardShortcuts string
	helpMarkdownBasics    string
	helpMathBasics        string
	dockNewWindow         string
	dockOpenFile          string
	dockRecentFiles       string
	unsavedTitle          string
	unsavedMessage        string
	closeWithoutSaving    string
	discard               string
	cancel                string
}

func getMenuStrings() menuI18n {
	return menuStrings[currentLocale]
}

func SetLocale(locale string) {
	switch Locale(locale) {
	case LocaleEn:
		currentLocale = LocaleEn
	default:
		currentLocale = LocaleZh
	}
}

func GetLocale() string {
	return string(currentLocale)
}

func buildMenuI18n(app *application.App) {
	ms := getMenuStrings()

	menu := app.NewMenu()

	if runtime.GOOS == "darwin" {
		appMenu := menu.AddSubmenu(ms.appName)
		setMenuIcon(appMenu.Add(ms.about).OnClick(func(_ *application.Context) {
			app.Event.Emit("app:aboutRequested")
		}), menuIconAbout)
		appMenu.AddSeparator()
		setMenuIcon(appMenu.Add(ms.preferences).SetAccelerator("Cmd+,").OnClick(func(_ *application.Context) {
			EmitToFocused(app, "menu:settings")
		}), menuIconPreferences)
		appMenu.AddSeparator()
		setMenuIcon(appMenu.Add(ms.quit).SetAccelerator("Cmd+Q").OnClick(func(_ *application.Context) {
			RequestAppQuit()
		}), menuIconQuit)
	} else {
		helpMenu := menu.AddSubmenu(ms.help)
		addHelpDocumentMenuItems(app, helpMenu, helpDocumentEntries(ms))
	}

	fileMenu := menu.AddSubmenu(ms.file)
	setMenuIcon(fileMenu.Add(ms.newFile).SetAccelerator("CmdOrCtrl+N").OnClick(func(_ *application.Context) {
		EmitToFocused(app, "menu:newFile")
	}), menuIconNewFile)
	setMenuIcon(fileMenu.Add(ms.newWindow).SetAccelerator("CmdOrCtrl+Shift+N").OnClick(func(_ *application.Context) {
		NewEditorWindow(app)
	}), menuIconNewWindow)
	setMenuIcon(fileMenu.Add(ms.open).SetAccelerator("CmdOrCtrl+O").OnClick(func(_ *application.Context) {
		EmitToFocused(app, "menu:open")
	}), menuIconOpen)
	fileMenu.AddSeparator()
	setMenuIcon(fileMenu.Add(ms.save).SetAccelerator("CmdOrCtrl+S").OnClick(func(_ *application.Context) {
		EmitToFocused(app, "menu:save")
	}), menuIconSave)
	setMenuIcon(fileMenu.Add(ms.saveAs).SetAccelerator("CmdOrCtrl+Shift+S").OnClick(func(_ *application.Context) {
		EmitToFocused(app, "menu:saveAs")
	}), menuIconSaveAs)
	fileMenu.AddSeparator()
	setMenuIcon(fileMenu.Add(ms.exportHtml).OnClick(func(_ *application.Context) {
		EmitToFocused(app, "menu:exportHTML")
	}), menuIconExportHTML)
	setMenuIcon(fileMenu.Add(ms.exportPdf).OnClick(func(_ *application.Context) {
		EmitToFocused(app, "menu:exportPDF")
	}), menuIconExportPDF)
	fileMenu.AddSeparator()
	setMenuIcon(fileMenu.Add(ms.quit).SetAccelerator("CmdOrCtrl+Q").OnClick(func(_ *application.Context) {
		RequestAppQuit()
	}), menuIconQuit)

	editorMenu := menu.AddSubmenu(ms.editor)
	editorMenu.AddRole(application.Undo)
	editorMenu.AddRole(application.Redo)
	editorMenu.AddSeparator()
	editorMenu.AddRole(application.Cut)
	editorMenu.AddRole(application.Copy)
	editorMenu.AddRole(application.Paste)
	editorMenu.AddSeparator()
	editorMenu.AddRole(application.SelectAll)

	viewMenu := menu.AddSubmenu(ms.view)
	setMenuIcon(viewMenu.Add(ms.toggleSidebar).SetAccelerator("CmdOrCtrl+Shift+L").OnClick(func(_ *application.Context) {
		EmitToFocused(app, "menu:toggleSidebar")
	}), menuIconSidebar)
	viewMenu.AddSeparator()
	themeMenu := viewMenu.AddSubmenu(ms.theme)
	setMenuIcon(viewMenu.FindByLabel(ms.theme), menuIconTheme)
	if themes := loadThemes(); len(themes) > 0 {
		for _, t := range themes {
			name := t.Name
			label := t.Label
			setMenuIcon(themeMenu.Add(label).OnClick(func(_ *application.Context) {
				EmitToFocused(app, "menu:setTheme", name)
			}), menuIconTheme)
		}
	}
	viewMenu.AddSeparator()
	setMenuIcon(viewMenu.Add(ms.enterFullscreen).SetAccelerator("Ctrl+CmdOrCtrl+F").OnClick(func(_ *application.Context) {
		ToggleFocusedFullscreen(app.Window.Current())
	}), menuIconFullscreen)
	viewMenu.AddSeparator()
	setMenuIcon(viewMenu.Add(ms.developerTools).OnClick(func(_ *application.Context) {
		openFocusedDeveloperTools(app)
	}), menuIconDevTools)

	helpMenu := menu.AddSubmenu(ms.help)
	addHelpDocumentMenuItems(app, helpMenu, helpDocumentEntries(ms))

	app.Menu.Set(menu)
	installConfiguredSystemMenuCleaners(ms.editor, ms.view)
	installDeveloperToolsShortcutDisplay(ms.view, ms.developerTools, DeveloperToolsShortcut)
	installCustomHelpMenu(ms.help)
}
