//go:build windows

package core

import "os/exec"

// revealInFinder selects the file in Explorer. `/select,` (no space) is
// the documented Explorer.exe switch for "open the parent folder and
// highlight this file".
func revealInFinder(path string) error {
	return exec.Command("explorer", "/select,"+path).Run()
}
