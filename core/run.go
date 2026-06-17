package core

import (
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

var (
	Service          *AppService
	allowedClose     = make(map[uint]bool)
	allowedCloseLock sync.Mutex
)

// AllowWindowClose marks the given window ID as authorized to actually close.
// The events.Common.WindowClosing hook checks this map; if the window is not
// allowed, the close is canceled and the frontend is asked what to do.
func AllowWindowClose(windowID uint) {
	allowedCloseLock.Lock()
	allowedClose[windowID] = true
	allowedCloseLock.Unlock()
}

// DeveloperToolsShortcut is the global accelerator for toggling DevTools.
const DeveloperToolsShortcut = "F12"

// NewEditorWindow creates a new editor window without an initial file.
func NewEditorWindow(app *application.App) *application.WebviewWindow {
	return NewEditorWindowWithFile(app, "")
}

// NewEditorWindowWithFile creates a new editor window optionally pre-loaded
// with a file path (passed via the "file" query param so the frontend can
// rehydrate state).
func NewEditorWindowWithFile(app *application.App, filePath string) *application.WebviewWindow {
	windowURL := "/"
	if filePath != "" {
		windowURL = "/?file=" + url.QueryEscape(filePath)
	}
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:          "fast-md",
		Width:          1280,
		Height:         800,
		MinWidth:       600,
		MinHeight:      400,
		URL:            windowURL,
		EnableFileDrop: true,
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropNormal,
			TitleBar:                application.MacTitleBarHiddenInset,
			InvisibleTitleBarHeight: 28,
		},
	})

	window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
		files := event.Context().DroppedFiles()
		for _, f := range files {
			lower := strings.ToLower(f)
			if strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".markdown") {
				// The user dragged a real file from Finder onto the
				// window — the directory it lives in is now trusted.
				if Service != nil {
					if err := Service.trustDir(filepath.Dir(f)); err != nil {
						log.Printf("trustDir(%s) failed: %v", filepath.Dir(f), err)
					}
				}
				window.EmitEvent("file:open", f)
				break
			}
		}
	})

	window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		allowedCloseLock.Lock()
		ok := allowedClose[window.ID()]
		allowedCloseLock.Unlock()
		if !ok {
			event.Cancel()
			window.ExecJS("window.dispatchEvent(new CustomEvent('window:closeRequested'))")
		} else {
			allowedCloseLock.Lock()
			delete(allowedClose, window.ID())
			allowedCloseLock.Unlock()
		}
	})

	window.RegisterHook(events.Mac.WindowShow, func(_ *application.WindowEvent) {
		setupTopBorderDoubleClick(window.NativeWindow())
	})

	return window
}

// EmitToFocused forwards an event to whichever window is currently focused.
func EmitToFocused(app *application.App, name string, data ...any) {
	if w := app.Window.Current(); w != nil {
		w.EmitEvent(name, data...)
	}
}

// FullscreenToggler is satisfied by any window that can enter/exit fullscreen.
type FullscreenToggler interface {
	ToggleFullscreen()
}

// ToggleFocusedFullscreen toggles fullscreen on the given window (nil-safe).
func ToggleFocusedFullscreen(window FullscreenToggler) {
	if window != nil {
		window.ToggleFullscreen()
	}
}

// RequestAppQuit triggers the quit coordinator. Safe to call when no
// service/app has been initialized.
func RequestAppQuit() {
	if Service != nil {
		Service.RequestQuit()
	}
}

// OpenDeveloperTools opens DevTools on the given window (nil-safe).
func OpenDeveloperTools(window application.Window) {
	if window != nil {
		window.OpenDevTools()
	}
}

func openFocusedDeveloperTools(app *application.App) {
	OpenDeveloperTools(app.Window.Current())
}

// firstFileFromArgs scans argv (typically os.Args[1:]) for the first path
// whose extension is .md or .markdown (case-insensitive), expanding a
// leading "~/" or lone "~" via os.UserHomeDir, and returns its absolute
// form. Returns "" when no candidate is found.
//
// The parser is intentionally narrow — only one positional file is
// accepted, mirroring what macOS Launch Services hands the app via
// ApplicationOpenedWithFile. Flags and unrelated args are ignored so a
// stray -psn_0 from Launch Services or a parent process's argv never
// confuses the heuristic.
func firstFileFromArgs(argv []string) string {
	home, _ := os.UserHomeDir()
	for _, arg := range argv {
		if arg == "" {
			continue
		}
		path := arg
		switch {
		case path == "~" && home != "":
			path = home
		case strings.HasPrefix(path, "~/") && home != "":
			path = filepath.Join(home, path[2:])
		}
		lower := strings.ToLower(path)
		if !strings.HasSuffix(lower, ".md") && !strings.HasSuffix(lower, ".markdown") {
			continue
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			continue
		}
		return abs
	}
	return ""
}

// RegisterDeveloperToolsShortcut binds F12 to toggle DevTools on the
// currently-focused window.
func RegisterDeveloperToolsShortcut(app *application.App) {
	app.KeyBinding.Add(DeveloperToolsShortcut, func(window application.Window) {
		OpenDeveloperTools(window)
	})
}

// Run wires up the Wails app and starts its event loop. assets must be the
// //go:embed all:frontend/dist docs/help FS from the root main package; it
// provides the bundled frontend and the help documents.
func Run(assets fs.FS) error {
	SetAssets(assets)

	Service = &AppService{}

	// The recent-files store is created here (not lazily) because the
	// dock menu is built further down and wants the persisted list on
	// first launch. A load failure (corrupt file, permission denied)
	// falls back to an empty list — the user's recents are a nice-to-
	// have, not load-bearing.
	if store, err := newRecentFilesStore(recentFilesPath()); err != nil {
		log.Printf("recent files store init failed: %v", err)
	} else {
		Service.recent = store
	}

	app := application.New(application.Options{
		Name:        "fast-md",
		Description: "A fast Markdown editor",
		Services: []application.Service{
			application.NewService(Service),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})
	RegisterDeveloperToolsShortcut(app)

	Service.app = app
	Service.quit = newQuitCoordinator(AllowWindowClose, app.Quit)

	cfg := LoadConfig()
	SetLocale(cfg.Language)

	// If the user invoked the embedded binary directly (bypassing the
	// Contents/Resources/fastmd wrapper), honour the .md arg from
	// os.Args. The wrapper at Contents/Resources/fastmd is the
	// canonical entry point — it routes through `open -a` so
	// kAEOpenDocuments hits the ApplicationOpenedWithFile handler
	// below and the file opens in a (possibly existing) instance. This
	// direct-binary branch is the fallback for users who invoke the
	// bundle executable manually.
	if initialFile := firstFileFromArgs(os.Args[1:]); initialFile != "" {
		if err := Service.trustDir(filepath.Dir(initialFile)); err != nil {
			log.Printf("trustDir(%s) failed: %v", filepath.Dir(initialFile), err)
		}
		NewEditorWindowWithFile(app, initialFile)
	} else {
		NewEditorWindow(app)
	}

	app.Event.OnApplicationEvent(events.Mac.ApplicationShouldHandleReopen, func(event *application.ApplicationEvent) {
		if !event.Context().HasVisibleWindows() {
			NewEditorWindow(app)
		}
	})

	app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
		path := event.Context().Filename()
		if path == "" {
			return
		}
		// The user double-clicked a file in Finder or used "Open With"
		// from another app — the OS expects a new window per request,
		// never to overwrite whatever the focused window is editing.
		if Service != nil {
			if err := Service.trustDir(filepath.Dir(path)); err != nil {
				log.Printf("trustDir(%s) failed: %v", filepath.Dir(path), err)
			}
		}
		NewEditorWindowWithFile(app, path)
	})

	buildMenuI18n(app)
	setupDockMenu()

	if err := app.Run(); err != nil {
		log.Printf("fast-md: app.Run failed: %v", err)
		return err
	}
	return nil
}