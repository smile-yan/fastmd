package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFirstFileFromArgs(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		// On some CI environments HOME is unset; the parser treats that
		// as "no tilde expansion", so we just substitute an empty home
		// and the tilde cases still verify the fallback behavior.
		home = ""
	}

	mustAbs := func(rel string) string {
		t.Helper()
		abs, err := filepath.Abs(rel)
		if err != nil {
			t.Fatalf("filepath.Abs(%q): %v", rel, err)
		}
		return abs
	}

	cases := []struct {
		name string
		argv []string
		want string // exact expected absolute path, or "" if no file matched
	}{
		{
			name: "no args",
			argv: nil,
			want: "",
		},
		{
			name: "single md",
			argv: []string{"notes.md"},
			want: mustAbs("notes.md"),
		},
		{
			name: "single markdown",
			argv: []string{"essay.markdown"},
			want: mustAbs("essay.markdown"),
		},
		{
			name: "case-insensitive extension",
			argv: []string{"TODO.MD"},
			want: mustAbs("TODO.MD"),
		},
		{
			name: "tilde expansion",
			argv: []string{"~/from-home.md"},
			want: filepath.Join(home, "from-home.md"),
		},
		{
			name: "non-md ignored",
			argv: []string{"--flag", "report.txt", "real.md"},
			want: mustAbs("real.md"),
		},
		{
			name: "multiple md takes the first",
			argv: []string{"first.md", "second.markdown"},
			want: mustAbs("first.md"),
		},
		{
			name: "only non-md args yields nothing",
			argv: []string{"-psn_0", "--debug", "image.png"},
			want: "",
		},
		{
			name: "empty entries skipped",
			argv: []string{"", "  ", "real.md"},
			want: mustAbs("real.md"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := firstFileFromArgs(tc.argv)
			if got != tc.want {
				t.Errorf("firstFileFromArgs(%v) = %q, want %q", tc.argv, got, tc.want)
			}
		})
	}
}

// TestFirstFileFromArgsPreservesAbsolute verifies that an already-absolute
// path is returned untouched. filepath.Abs on an absolute path on Unix is
// effectively a no-op (modulo cleaning), so the result should round-trip.
func TestFirstFileFromArgsPreservesAbsolute(t *testing.T) {
	arg := "/tmp/already-absolute.md"
	got := firstFileFromArgs([]string{arg})
	want, err := filepath.Abs(arg)
	if err != nil {
		t.Fatalf("filepath.Abs(%q): %v", arg, err)
	}
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
