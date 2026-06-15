//go:build darwin

package core

/*
#cgo CFLAGS: -mmacosx-version-min=12.0
#cgo LDFLAGS: -framework Cocoa

#include <stdlib.h>

// Only declarations here — no definitions — required when using //export.
extern void showCloseSheet(unsigned int did, void *winPtr,
                            const char *fn, const char *dir);
*/
import "C"

import (
	"path/filepath"
	"sync"
	"unsafe"
)

type _closeSheetResult struct {
	action int
	path   string
}

var (
	_csMu  sync.Mutex
	_csChs = map[uint]chan _closeSheetResult{}
	_csID  uint
)

//export goCloseSheetResult
func goCloseSheetResult(id C.uint, action C.int, path unsafe.Pointer) {
	_csMu.Lock()
	ch := _csChs[uint(id)]
	delete(_csChs, uint(id))
	_csMu.Unlock()
	if ch != nil {
		p := ""
		if path != nil {
			p = C.GoString((*C.char)(path))
		}
		ch <- _closeSheetResult{int(action), p}
	}
}

// ShowCloseSheet returns "cancel", "discard", or "save:<absolute-path>" for
// an unsaved new document.
func (s *AppService) ShowCloseSheet(filename, lastDir string) string {
	_csMu.Lock()
	_csID++
	id := _csID
	ch := make(chan _closeSheetResult, 1)
	_csChs[id] = ch
	_csMu.Unlock()

	fn := C.CString(filename)
	dir := C.CString(lastDir)
	// showCloseSheet copies fn/dir into NSStrings before returning,
	// so it is safe to free them immediately after the call.
	C.showCloseSheet(C.uint(id), unsafe.Pointer(s.focusedWindow().NativeWindow()), fn, dir)
	C.free(unsafe.Pointer(fn))
	C.free(unsafe.Pointer(dir))

	r := <-ch
	switch r.action {
	case 1:
		return "discard"
	case 2:
		// The user picked a save destination through NSSavePanel — that
		// destination is now a trusted directory. (On non-darwin, the
		// ShowCloseSheet stub calls SaveFileDialog, which already
		// registers trust.)
		if r.path != "" {
			if err := s.trustDir(filepath.Dir(r.path)); err != nil {
				_ = err // best-effort; dialog still returns the path
			}
		}
		return "save:" + r.path
	default:
		return "cancel"
	}
}
