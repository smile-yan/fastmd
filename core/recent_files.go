package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// maxRecentFiles caps the persisted recent-files list. The macOS dock menu
// shows at most a handful before users stop reading, and a small JSON file
// is cheap to load on launch.
const maxRecentFiles = 10

// recentFilesStore tracks the user's recently-opened markdown files and
// persists them across launches. The store is safe for concurrent use;
// the dock-menu update path takes a snapshot under the lock and works
// from that copy, so the C interop never observes a partially-mutated
// slice.
type recentFilesStore struct {
	mu    sync.Mutex
	paths []string
	// file is the absolute path of the JSON persistence file. The store
	// is the only writer; concurrent writers to the same file would
	// race the JSON marshal, so the mutex also covers save().
	file string
}

func newRecentFilesStore(persistPath string) (*recentFilesStore, error) {
	s := &recentFilesStore{file: persistPath}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

// load reads the persisted list from disk. A missing file yields an
// empty list (the user simply hasn't opened anything yet). A corrupt
// file is treated as empty too — the next save overwrites it, so a
// bad write doesn't poison subsequent launches.
func (s *recentFilesStore) load() error {
	data, err := os.ReadFile(s.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var paths []string
	if err := json.Unmarshal(data, &paths); err != nil {
		// Corrupt or wrong schema — start fresh.
		return nil
	}
	s.mu.Lock()
	s.paths = normalizeAndDedupe(paths)
	s.mu.Unlock()
	return nil
}

// save writes the current list to disk. mkdir-p the parent directory
// first; on a fresh machine the Application Support folder may not
// exist yet.
func (s *recentFilesStore) save() error {
	snapshot := s.Snapshot()
	if err := os.MkdirAll(filepath.Dir(s.file), 0755); err != nil {
		return err
	}
	data, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	return os.WriteFile(s.file, data, 0644)
}

// Add prepends path to the list, removes any prior occurrence, caps the
// list at maxRecentFiles, and persists. Returns the updated snapshot
// so the caller (e.g. the dock menu) can refresh itself in one round
// trip. Non-markdown paths are ignored.
func (s *recentFilesStore) Add(path string) ([]string, error) {
	if !isMarkdownPath(path) {
		return s.Snapshot(), nil
	}

	s.mu.Lock()
	filtered := s.paths[:0:0]
	for _, p := range s.paths {
		if p != path {
			filtered = append(filtered, p)
		}
	}
	filtered = append([]string{path}, filtered...)
	if len(filtered) > maxRecentFiles {
		filtered = filtered[:maxRecentFiles]
	}
	s.paths = filtered
	snapshot := append([]string(nil), s.paths...)
	s.mu.Unlock()

	if err := s.save(); err != nil {
		return snapshot, err
	}
	return snapshot, nil
}

// Snapshot returns an independent copy of the current list. Callers can
// iterate it without holding the lock or worrying about the underlying
// slice being mutated underneath them.
func (s *recentFilesStore) Snapshot() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.paths...)
}

func isMarkdownPath(path string) bool {
	lower := strings.ToLower(path)
	return strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".markdown")
}

// normalizeAndDedupe trims empties, removes duplicates (preserving the
// first occurrence so the most-recent position is kept), and caps the
// list. Used by load() to clean up a list that may have been written
// by an older or buggy version.
func normalizeAndDedupe(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		if p == "" {
			continue
		}
		if _, dup := seen[p]; dup {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	if len(out) > maxRecentFiles {
		out = out[:maxRecentFiles]
	}
	return out
}

// recentFilesPath is the canonical location for the persisted recent
// list, alongside the config file.
func recentFilesPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Application Support", "fastmd", "recent.json")
}
