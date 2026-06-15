package core

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// All tests in this file exercise the trust-root sandbox that gates
// ReadFile / WriteFile / ListDirectory. The contract:
//
//   - A path is allowed only if it is inside a directory the app has
//     been told to trust (via the dialogs or the file-open channels).
//   - "Inside" is checked via filepath.Rel: a path that resolves ".."
//     outside the root is rejected.
//   - os.Root enforces kernel-level containment, including a refusal to
//     follow symlinks that point outside the root.

func newServiceWithTrusted(t *testing.T, dir string) *AppService {
	t.Helper()
	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("Abs: %v", err)
	}
	svc := &AppService{}
	if err := svc.trustDir(abs); err != nil {
		t.Fatalf("trustDir(%s): %v", abs, err)
	}
	return svc
}

func TestReadFileRejectsPathOutsideAnyTrustedRoot(t *testing.T) {
	trusted := t.TempDir()
	outside := t.TempDir()
	outsideFile := filepath.Join(outside, "secret.md")
	if err := os.WriteFile(outsideFile, []byte("nope"), 0644); err != nil {
		t.Fatal(err)
	}

	svc := newServiceWithTrusted(t, trusted)

	_, err := svc.ReadFile(outsideFile)
	if err == nil {
		t.Fatal("expected ReadFile to reject path outside trusted root")
	}
	if !strings.Contains(err.Error(), "not inside a trusted directory") {
		t.Fatalf("expected trust error, got %v", err)
	}
}

func TestWriteFileRejectsPathOutsideAnyTrustedRoot(t *testing.T) {
	trusted := t.TempDir()
	outside := t.TempDir()

	svc := newServiceWithTrusted(t, trusted)

	err := svc.WriteFile(filepath.Join(outside, "evil.md"), "pwned")
	if err == nil {
		t.Fatal("expected WriteFile to reject path outside trusted root")
	}
	if !strings.Contains(err.Error(), "not inside a trusted directory") {
		t.Fatalf("expected trust error, got %v", err)
	}
}

func TestListDirectoryRejectsPathOutsideAnyTrustedRoot(t *testing.T) {
	trusted := t.TempDir()
	outside := t.TempDir()
	// Plant a sensitive-looking file outside the trusted root; if the
	// trust check is broken, the test will see it in the listing.
	if err := os.WriteFile(filepath.Join(outside, "creds.md"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	svc := newServiceWithTrusted(t, trusted)

	_, err := svc.ListDirectory(outside)
	if err == nil {
		t.Fatal("expected ListDirectory to reject path outside trusted root")
	}
	if !strings.Contains(err.Error(), "not inside a trusted directory") {
		t.Fatalf("expected trust error, got %v", err)
	}
}

func TestReadFileAcceptsPathInsideTrustedRoot(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "note.md")
	if err := os.WriteFile(target, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	svc := newServiceWithTrusted(t, dir)

	got, err := svc.ReadFile(target)
	if err != nil {
		t.Fatalf("ReadFile inside trusted root: %v", err)
	}
	if got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestReadFileAcceptsDeeplyNestedPathInsideTrustedRoot(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "a", "b", "c")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(sub, "deep.md")
	if err := os.WriteFile(target, []byte("buried"), 0644); err != nil {
		t.Fatal(err)
	}

	svc := newServiceWithTrusted(t, dir)
	got, err := svc.ReadFile(target)
	if err != nil {
		t.Fatalf("ReadFile nested: %v", err)
	}
	if got != "buried" {
		t.Fatalf("expected %q, got %q", "buried", got)
	}
}

func TestListDirectoryInsideTrustedRootReturnsEntries(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.md"), nil, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ignored.txt"), nil, 0644); err != nil {
		t.Fatal(err)
	}

	svc := newServiceWithTrusted(t, dir)
	entries, err := svc.ListDirectory(dir)
	if err != nil {
		t.Fatalf("ListDirectory: %v", err)
	}

	names := make(map[string]bool)
	for _, e := range entries {
		names[e.Name] = true
	}
	if !names["a.md"] || !names["sub"] {
		t.Fatalf("expected a.md and sub, got %v", names)
	}
	if names["ignored.txt"] {
		t.Fatalf("non-markdown file should be filtered out, got %v", names)
	}
}

func TestTrustDirIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	svc := &AppService{}
	if err := svc.trustDir(dir); err != nil {
		t.Fatal(err)
	}
	if err := svc.trustDir(dir); err != nil {
		t.Fatalf("second trustDir should be a no-op, got %v", err)
	}
	svc.rootsMu.RLock()
	defer svc.rootsMu.RUnlock()
	if got := len(svc.roots); got != 1 {
		t.Fatalf("expected 1 registered root, got %d", got)
	}
}

func TestMultipleTrustedRootsAllWork(t *testing.T) {
	a := t.TempDir()
	b := t.TempDir()
	fileA := filepath.Join(a, "in-a.md")
	fileB := filepath.Join(b, "in-b.md")
	if err := os.WriteFile(fileA, []byte("A"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fileB, []byte("B"), 0644); err != nil {
		t.Fatal(err)
	}

	svc := &AppService{}
	if err := svc.trustDir(a); err != nil {
		t.Fatal(err)
	}
	if err := svc.trustDir(b); err != nil {
		t.Fatal(err)
	}

	if got, err := svc.ReadFile(fileA); err != nil || got != "A" {
		t.Fatalf("ReadFile(a): got=%q err=%v", got, err)
	}
	if got, err := svc.ReadFile(fileB); err != nil || got != "B" {
		t.Fatalf("ReadFile(b): got=%q err=%v", got, err)
	}
}

func TestSymlinkEscapeIsBlockedByKernel(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink semantics differ on Windows; covered by ReadFileAcceptsPathInsideTrustedRoot")
	}
	trusted := t.TempDir()
	outside := t.TempDir()
	outsideFile := filepath.Join(outside, "secret.md")
	if err := os.WriteFile(outsideFile, []byte("classified"), 0644); err != nil {
		t.Fatal(err)
	}

	// Plant a symlink inside the trusted root that points at the
	// outside file. If the trust check is broken or relies only on
	// filepath.Rel, the read would succeed and leak "classified".
	link := filepath.Join(trusted, "leak.md")
	if err := os.Symlink(outsideFile, link); err != nil {
		t.Skipf("Symlink not supported in this environment: %v", err)
	}

	svc := newServiceWithTrusted(t, trusted)
	got, err := svc.ReadFile(link)
	if err == nil {
		t.Fatalf("expected symlink escape to be blocked, got content %q", got)
	}
	if got == "classified" {
		t.Fatalf("symlink escape leaked file content: %q", got)
	}
}

func TestWriteFileInsideTrustedRootCreatesFile(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "new.md")

	svc := newServiceWithTrusted(t, dir)
	if err := svc.WriteFile(target, "written"); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "written" {
		t.Fatalf("expected %q, got %q", "written", got)
	}
}

func TestWriteFileRejectsPathTraversalAttempt(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("path traversal semantics differ on Windows")
	}
	trusted := t.TempDir()
	outside := t.TempDir()
	outsideFile := filepath.Join(outside, "victim.md")

	svc := newServiceWithTrusted(t, trusted)

	// Even if the frontend tries to be clever with "../" segments, the
	// trust check must reject the path before any file is touched.
	escaped := filepath.Join(trusted, "..", "..", filepath.Base(outside), "victim.md")
	if err := svc.WriteFile(escaped, "pwned"); err == nil {
		t.Fatal("expected path-traversal WriteFile to be rejected")
	}
	if _, err := os.Stat(outsideFile); err == nil {
		t.Fatal("path-traversal WriteFile actually created a file outside the root")
	}
}

func TestReadFileRejectsAbsolutePathNotInAnyRoot(t *testing.T) {
	// Regression: the untrusted read should not be able to slip in by
	// using an absolute path with a different prefix than the trusted
	// directory.
	trusted := t.TempDir()
	svc := newServiceWithTrusted(t, trusted)
	if _, err := svc.ReadFile("/etc/passwd"); err == nil {
		t.Fatal("expected /etc/passwd to be rejected when /etc is not trusted")
	}
}

func TestTrustDirRejectsNonExistentPath(t *testing.T) {
	svc := &AppService{}
	if err := svc.trustDir(filepath.Join(t.TempDir(), "does-not-exist")); err == nil {
		t.Fatal("expected trustDir to fail on a non-existent path")
	}
}

func TestTrustDirRejectsFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "a-file.md")
	if err := os.WriteFile(file, nil, 0644); err != nil {
		t.Fatal(err)
	}
	svc := &AppService{}
	if err := svc.trustDir(file); err == nil {
		t.Fatal("expected trustDir to fail on a non-directory")
	}
}
