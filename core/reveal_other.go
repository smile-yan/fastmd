//go:build !darwin && !windows

package core

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// revealInFinder opens the parent directory in the user's default file
// manager. There's no portable "select a file" semantic on Linux, so the
// best we can do is open the containing folder via xdg-open.
func revealInFinder(path string) error {
	dir := filepath.Dir(path)
	if err := exec.Command("xdg-open", dir).Run(); err != nil {
		return fmt.Errorf("xdg-open %q failed: %w", dir, err)
	}
	return nil
}
