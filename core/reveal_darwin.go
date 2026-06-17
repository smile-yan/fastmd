//go:build darwin

package core

import (
	"fmt"
	"os/exec"
)

// revealInFinder selects the file in Finder. `open -R <path>` is the
// macOS-standard "reveal in Finder" command: it tells Finder to open
// its parent directory and highlight the given file.
func revealInFinder(path string) error {
	cmd := exec.Command("open", "-R", path)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("open -R %q failed: %w (%s)", path, err, out)
	}
	return nil
}
