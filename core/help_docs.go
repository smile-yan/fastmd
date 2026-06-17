package core

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// helpDocumentFS is the FS that backs help document reads. It is initialized
// by SetAssets to fs.Sub(assets, "docs/help"), so all paths inside it are
// bare filenames ("quick-start.md", etc.).
//
// Declared in assets.go.

type helpDocumentEntry struct {
	label    string
	filename string
	icon     menuIconName
}

func helpDocumentEntries(ms menuI18n) []helpDocumentEntry {
	return []helpDocumentEntry{
		{label: ms.helpQuickStart, filename: "quick-start.md", icon: menuIconQuickStart},
		{label: ms.helpKeyboardShortcuts, filename: "keyboard-shortcuts.md", icon: menuIconShortcuts},
		{label: ms.helpMarkdownBasics, filename: "markdown-basics.md", icon: menuIconMarkdown},
		{label: ms.helpMathBasics, filename: "math-formula-basics.md", icon: menuIconMath},
	}
}

func addHelpDocumentMenuItems(app *application.App, helpMenu *application.Menu, entries []helpDocumentEntry) {
	for _, entry := range entries {
		entry := entry
		setMenuIcon(helpMenu.Add(entry.label).OnClick(func(_ *application.Context) {
			if err := openHelpDocument(app, entry.filename); err != nil {
				log.Printf("failed to open help document %q: %v", entry.filename, err)
			}
		}), entry.icon)
	}
}

func openHelpDocument(app *application.App, filename string) error {
	path, err := materializeHelpDocument(filename, helpDocumentCacheDir())
	if err != nil {
		return err
	}
	NewEditorWindowWithFile(app, path)
	return nil
}

func helpDocumentCacheDir() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil || cacheDir == "" {
		cacheDir = os.TempDir()
	}
	return filepath.Join(cacheDir, "fastmd", "help")
}

func materializeHelpDocument(filename, outputDir string) (string, error) {
	if filename == "" || filepath.Base(filename) != filename {
		return "", fmt.Errorf("invalid help document filename %q", filename)
	}
	if helpDocumentFS == nil {
		return "", fmt.Errorf("help document assets not initialized")
	}
	data, err := fs.ReadFile(helpDocumentFS, filename)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}
	outputPath := filepath.Join(outputDir, filename)
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return "", err
	}
	return outputPath, nil
}