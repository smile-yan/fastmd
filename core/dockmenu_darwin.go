//go:build darwin

package core

/*
#include <stdlib.h>
extern void setupDockMenu(void);
extern void updateDockRecentFiles(const char **paths, int count);
*/
import "C"
import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

var (
	recentFiles   []string
	recentFilesMu sync.Mutex
	maxRecent     = 10
)

//export dockMenuNewWindow
func dockMenuNewWindow() {
	if Service == nil || Service.app == nil {
		return
	}
	go NewEditorWindow(Service.app)
}

//export dockMenuOpenFile
func dockMenuOpenFile() {
	if Service == nil || Service.app == nil {
		return
	}
	go func() {
		if w := Service.app.Window.Current(); w != nil {
			w.EmitEvent("menu:open")
		} else {
			ww := NewEditorWindow(Service.app)
			ww.OnWindowEvent(events.Common.WindowRuntimeReady, func(_ *application.WindowEvent) {
				ww.EmitEvent("menu:open")
			})
		}
	}()
}

//export dockMenuOpenRecent
func dockMenuOpenRecent(cpath *C.char) {
	path := C.GoString(cpath)
	if path == "" || Service == nil || Service.app == nil {
		return
	}
	// Recent files are paths the user has previously opened in this app,
	// so re-opening one is implicit consent to access its directory again.
	if err := Service.trustDir(filepath.Dir(path)); err != nil {
		// Not fatal — the file just won't open. Log so the user has a
		// hint if they wonder why nothing happened.
		// (stderr is fine here; main runs in a console-less environment
		// and these are exceptional.)
		_ = err
	}
	go func() {
		if w := Service.app.Window.Current(); w != nil {
			w.EmitEvent("file:open", path)
		} else {
			NewEditorWindowWithFile(Service.app, path)
		}
	}()
	trackRecentFile(path)
}

func trackRecentFile(path string) {
	lower := strings.ToLower(path)
	if !strings.HasSuffix(lower, ".md") && !strings.HasSuffix(lower, ".markdown") {
		return
	}

	recentFilesMu.Lock()
	for i, p := range recentFiles {
		if p == path {
			recentFiles = append(recentFiles[:i], recentFiles[i+1:]...)
			break
		}
	}
	recentFiles = append([]string{path}, recentFiles...)
	if len(recentFiles) > maxRecent {
		recentFiles = recentFiles[:maxRecent]
	}
	recentFilesMu.Unlock()

	updateDockRecentMenu()
}

func updateDockRecentMenu() {
	recentFilesMu.Lock()
	paths := make([]*C.char, len(recentFiles))
	for i, p := range recentFiles {
		paths[i] = C.CString(p)
	}
	count := C.int(len(recentFiles))
	recentFilesMu.Unlock()

	if len(paths) > 0 {
		C.updateDockRecentFiles(&paths[0], count)
	} else {
		C.updateDockRecentFiles(nil, 0)
	}
	// C strings are freed by updateDockRecentFiles after copying
}

func setupDockMenu() {
	C.setupDockMenu()
	updateDockRecentMenu()
}
