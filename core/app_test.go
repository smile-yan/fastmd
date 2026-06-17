package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadFile(t *testing.T) {
	tmp, err := os.CreateTemp("", "*.md")
	if err != nil {
		t.Fatal(err)
	}
	tmp.WriteString("# Hello World")
	tmp.Close()
	defer os.Remove(tmp.Name())

	// ReadFile requires the file to be inside a trusted root; the temp
	// file's parent directory is the smallest scope that satisfies the
	// check.
	svc := &AppService{}
	if err := svc.trustDir(filepath.Dir(tmp.Name())); err != nil {
		t.Fatalf("trustDir: %v", err)
	}
	content, err := svc.ReadFile(tmp.Name())
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if content != "# Hello World" {
		t.Errorf("expected '# Hello World', got %q", content)
	}
}

func TestWriteFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")

	svc := &AppService{}
	if err := svc.trustDir(dir); err != nil {
		t.Fatalf("trustDir: %v", err)
	}
	if err := svc.WriteFile(path, "# Test"); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}
	data, _ := os.ReadFile(path)
	if string(data) != "# Test" {
		t.Errorf("expected '# Test', got %q", string(data))
	}
}

func TestListDirectory(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.md"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "c.markdown"), []byte(""), 0644)
	os.Mkdir(filepath.Join(dir, "subdir"), 0755)

	svc := &AppService{}
	if err := svc.trustDir(dir); err != nil {
		t.Fatalf("trustDir: %v", err)
	}
	files, err := svc.ListDirectory(dir)
	if err != nil {
		t.Fatalf("ListDirectory error: %v", err)
	}

	names := make(map[string]bool)
	for _, f := range files {
		names[f.Name] = true
	}
	if !names["a.md"] {
		t.Error("expected a.md")
	}
	if !names["c.markdown"] {
		t.Error("expected c.markdown")
	}
	if !names["subdir"] {
		t.Error("expected subdir")
	}
	if names["b.txt"] {
		t.Error("b.txt should be excluded")
	}
}

func TestGetHomePath(t *testing.T) {
	svc := &AppService{}
	home := svc.GetHomePath()
	if home == "" {
		t.Error("expected non-empty home path")
	}
	if !strings.HasPrefix(home, "/") {
		t.Errorf("expected absolute path, got %q", home)
	}
}

func TestRevealInFinderEmptyPath(t *testing.T) {
	// The app-level guard around revealInFinder: an empty path must
	// produce an error before any platform shell-out is attempted.
	svc := &AppService{}
	if err := svc.RevealInFinder(""); err == nil {
		t.Error("expected error for empty path, got nil")
	}
}
