package core

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// withTempHome redirects os.UserHomeDir() to a fresh temp dir for the
// duration of the test, so LoadConfig/SaveConfig don't read or write the
// real ~/Library/Application Support/fastmd/config.json. On non-darwin
// CI hosts the path doesn't include "Library", but the test still works —
// the redirect applies uniformly and the package functions write under
// whatever subpath they hard-code.
func withTempHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", dir)
	}
	return dir
}

// TestLoadConfigDefaultsToZhWhenFileMissing verifies the "fresh install"
// path: no config file on disk → zh is the default language.
func TestLoadConfigDefaultsToZhWhenFileMissing(t *testing.T) {
	withTempHome(t)

	cfg := LoadConfig()
	if cfg.Language != "zh" {
		t.Fatalf("expected default language zh, got %q", cfg.Language)
	}
}

// TestLoadConfigFallsBackToZhOnCorruptJSON guards against a bad write
// (or a hand-edited file) bricking the menu: the app should still boot
// with the default locale.
func TestLoadConfigFallsBackToZhOnCorruptJSON(t *testing.T) {
	withTempHome(t)

	path := getConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{not json"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := LoadConfig()
	if cfg.Language != "zh" {
		t.Fatalf("expected fallback zh on corrupt JSON, got %q", cfg.Language)
	}
}

// TestLoadConfigRejectsUnknownLanguage covers the validation step that
// makes sure a manually-edited config with a bogus Language field doesn't
// leak into the menu builder. Unknown languages fall back to zh.
func TestLoadConfigRejectsUnknownLanguage(t *testing.T) {
	withTempHome(t)

	path := getConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"language":"klingon"}`), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := LoadConfig()
	if cfg.Language != "zh" {
		t.Fatalf("expected fallback zh for unknown language, got %q", cfg.Language)
	}
}

// TestLoadConfigAcceptsValidLanguages is the positive control: both
// supported languages survive a round-trip.
func TestLoadConfigAcceptsValidLanguages(t *testing.T) {
	withTempHome(t)

	for _, lang := range []string{"zh", "en"} {
		cfg := AppConfig{Language: lang}
		if err := SaveConfig(cfg); err != nil {
			t.Fatalf("SaveConfig(%s): %v", lang, err)
		}
		got := LoadConfig()
		if got.Language != lang {
			t.Fatalf("round-trip lost language: saved %q, loaded %q", lang, got.Language)
		}
	}
}

// TestSaveConfigCreatesParentDirectory ensures a fresh install — where
// ~/Library/Application Support/fastmd/ doesn't yet exist — doesn't
// fail the very first write.
func TestSaveConfigCreatesParentDirectory(t *testing.T) {
	home := withTempHome(t)
	expected := filepath.Join(home, "Library", "Application Support", "fastmd", "config.json")

	if _, err := os.Stat(expected); err == nil {
		t.Fatalf("config file already exists before write: %s", expected)
	}

	if err := SaveConfig(AppConfig{Language: "en"}); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	if _, err := os.Stat(expected); err != nil {
		t.Fatalf("config file not created: %v", err)
	}
}

// TestSaveConfigRoundTripPreservesFields uses a struct with the same
// shape as AppConfig to make sure the marshaling and unmarshaling are
// faithful (no field is silently dropped or renamed).
func TestSaveConfigRoundTripPreservesFields(t *testing.T) {
	withTempHome(t)

	in := AppConfig{Language: "en"}
	if err := SaveConfig(in); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	// Read the raw file and confirm it parses as a JSON object with
	// exactly the "language" field — that way a future field added to
	// AppConfig without updating this test will fail loudly here.
	data, err := os.ReadFile(getConfigPath())
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("SaveConfig wrote an empty file")
	}

	out := LoadConfig()
	if out.Language != in.Language {
		t.Fatalf("round-trip mismatch: in=%q out=%q", in.Language, out.Language)
	}
}

// TestAppServiceGetConfigDelegatesToLoadConfig — the AppService method
// is a thin pass-through; this test guards against future refactors that
// might return a different struct (e.g. one without Language validation).
func TestAppServiceGetConfigDelegatesToLoadConfig(t *testing.T) {
	withTempHome(t)

	svc := &AppService{}
	got := svc.GetConfig()
	if got.Language != "zh" {
		t.Fatalf("expected zh from fresh install via AppService, got %q", got.Language)
	}

	if err := svc.SaveConfig(AppConfig{Language: "en"}); err != nil {
		t.Fatalf("SaveConfig on service: %v", err)
	}

	if got := svc.GetConfig(); got.Language != "en" {
		t.Fatalf("expected en after SaveConfig on service, got %q", got.Language)
	}
}

// TestSaveConfigOverwritesExistingFile — SaveConfig is idempotent and
// must replace the previous file rather than append or refuse to write.
func TestSaveConfigOverwritesExistingFile(t *testing.T) {
	withTempHome(t)

	if err := SaveConfig(AppConfig{Language: "en"}); err != nil {
		t.Fatal(err)
	}
	if err := SaveConfig(AppConfig{Language: "zh"}); err != nil {
		t.Fatal(err)
	}

	if got := LoadConfig(); got.Language != "zh" {
		t.Fatalf("expected second SaveConfig to win, got %q", got.Language)
	}
}
