package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ── Trusted filesystem root registry ───────────────────────────────────────
//
// ReadFile, WriteFile, and ListDirectory only operate on paths that live
// inside a directory the user has explicitly handed us through a native
// dialog (OpenFileDialog, OpenFolderDialog, SaveFileDialog, ShowCloseSheet)
// or through the OS-level file-open / drag-drop channels. Everything else
// is rejected with an error.
//
// Containment is enforced at the kernel layer via os.Root (Go 1.24+):
//   - os.Root methods follow symbolic links but refuse to follow any that
//     point outside the root.
//   - Symlinks must be relative (absolute symlinks are rejected).
//   - The root holds an open file descriptor, so operations stay scoped
//     to the directory even if it's moved or replaced underneath us.
//
// The registry is in-memory only. Trust does not survive across launches,
// which matches user expectation: a session is a working set of folders,
// not a permanent grant.

func (s *AppService) trustDir(dir string) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", abs)
	}

	// Lazy init: the existing zero-value construction in run.go and in
	// tests (quit_test.go uses &AppService{quit: coordinator}) leaves
	// roots nil; a sync.Once is the cheapest way to handle both.
	s.rootsOnce.Do(func() {
		s.roots = make(map[string]*os.Root)
	})

	s.rootsMu.RLock()
	if _, ok := s.roots[abs]; ok {
		s.rootsMu.RUnlock()
		return nil
	}
	s.rootsMu.RUnlock()

	// os.OpenRoot on Unix uses open(O_DIRECTORY|O_RDONLY|O_NOFOLLOW), so
	// passing a symlink as the root itself fails. We let that surface as
	// an error — the caller is expected to pick a real directory.
	root, err := os.OpenRoot(abs)
	if err != nil {
		return err
	}

	s.rootsMu.Lock()
	defer s.rootsMu.Unlock()
	if existing, ok := s.roots[abs]; ok {
		// Lost the race; close the duplicate and return the existing one.
		_ = existing.Close() //nolint:errcheck // best-effort cleanup
		_ = root.Close()     //nolint:errcheck // best-effort cleanup
		return nil
	}
	s.roots[abs] = root
	return nil
}

// resolveTrusted splits absPath into the registered root that contains it
// and the path relative to that root. Returns an error when no registered
// root contains the path. Symlink escapes are caught later, when the
// caller actually invokes the root method (the kernel refuses the open
// with EACCES or similar).
func (s *AppService) resolveTrusted(absPath string) (*os.Root, string, error) {
	abs, err := filepath.Abs(absPath)
	if err != nil {
		return nil, "", err
	}
	s.rootsMu.RLock()
	defer s.rootsMu.RUnlock()
	if len(s.roots) == 0 {
		return nil, "", fmt.Errorf("no trusted directories registered")
	}
	for rootPath, root := range s.roots {
		rel, err := filepath.Rel(rootPath, abs)
		if err != nil {
			continue
		}
		// Reject paths that escape the root via "..". filepath.Rel only
		// produces ".." when the paths are on different roots (e.g. one
		// is on a different drive on Windows), but guarding here is cheap.
		if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			continue
		}
		return root, rel, nil
	}
	return nil, "", fmt.Errorf("path is not inside a trusted directory: %s", abs)
}
