package core

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/carlos7ags/folio/document"
	"github.com/carlos7ags/folio/html"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// AppInfo carries the metadata shown in the About dialog. Version is sourced
// from runtime/debug.ReadBuildInfo — when the binary is built with
// `-ldflags '-X main.version=...'`, that value surfaces here; otherwise Go
// fills it with vcs.revision and vcs.time from the build settings. The Go
// runtime version is included so users can file useful bug reports.
type AppInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Built   string `json:"built"`
	Go      string `json:"go"`
}

const appName = "fastmd"

func (s *AppService) GetAppInfo() AppInfo {
	info := AppInfo{
		Name:    appName,
		Version: "dev",
		Go:      runtime.Version(),
	}
	if bi, ok := debug.ReadBuildInfo(); ok {
		if v := bi.Main.Version; v != "" && v != "(devel)" {
			info.Version = v
		}
		for _, set := range bi.Settings {
			switch set.Key {
			case "vcs.revision":
				info.Commit = set.Value
			case "vcs.time":
				info.Built = set.Value
			}
		}
	}
	return info
}

type FileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"isDir"`
}

type AppConfig struct {
	Language string `json:"language"`
}

func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Application Support", "fastmd", "config.json")
}

// LoadConfig reads the user's persisted config, falling back to defaults
// when the file is missing or malformed.
func LoadConfig() AppConfig {
	path := getConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return AppConfig{Language: "zh"}
	}
	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return AppConfig{Language: "zh"}
	}
	if cfg.Language != "zh" && cfg.Language != "en" {
		cfg.Language = "zh"
	}
	return cfg
}

// SaveConfig writes the config to disk, creating the directory if needed.
func SaveConfig(cfg AppConfig) error {
	path := getConfigPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

type AppService struct {
	app  *application.App
	quit *quitCoordinator

	// roots is the set of directories the user has explicitly granted the
	// app access to. Lookups and inserts are guarded by rootsMu; rootsOnce
	// lazily allocates the map so callers that never touch the filesystem
	// (e.g. the quit coordinator tests) can keep using zero-value
	// construction.
	rootsOnce sync.Once
	rootsMu   sync.RWMutex
	roots     map[string]*os.Root

	// recent is the persisted recently-opened-files store. It is created
	// in Run() before the dock menu is built so the dock menu reflects
	// the user's history from the previous session on first launch.
	recent *recentFilesStore
}

func (s *AppService) focusedWindow() application.Window {
	return s.app.Window.Current()
}

func (s *AppService) OpenFileDialog() (string, error) {
	d := s.app.Dialog.OpenFile().
		SetTitle("Open Markdown File").
		AddFilter("Markdown Files", "*.md;*.markdown").
		CanChooseFiles(true)
	if w := s.focusedWindow(); w != nil {
		d = d.AttachToWindow(w)
	}
	path, err := d.PromptForSingleSelection()
	// The user explicitly chose this file, so its directory is now
	// trusted for the rest of the session. trustDir is best-effort —
	// a failure here is logged but doesn't block the dialog result.
	if err == nil && path != "" {
		if trustErr := s.trustDir(filepath.Dir(path)); trustErr != nil {
			fmt.Fprintf(os.Stderr, "trustDir(%s) failed: %v\n", filepath.Dir(path), trustErr)
		}
	}
	return path, err
}

func (s *AppService) SaveFileDialog(defaultPath string) (string, error) {
	d := s.app.Dialog.SaveFile().
		AttachToWindow(s.focusedWindow()).
		AddFilter("Markdown Files", "*.md")
	if defaultPath != "" {
		d.SetDirectory(filepath.Dir(defaultPath))
		d.SetFilename(filepath.Base(defaultPath))
	} else {
		d.SetFilename("untitled.md")
	}
	path, err := d.PromptForSingleSelection()
	// Trust the directory of any path the user chose to save into. A save
	// dialog is the user's explicit opt-in to write into a location, so
	// that's the moment the directory becomes part of the trust set.
	if err == nil && path != "" {
		if trustErr := s.trustDir(filepath.Dir(path)); trustErr != nil {
			fmt.Fprintf(os.Stderr, "trustDir(%s) failed: %v\n", filepath.Dir(path), trustErr)
		}
	}
	return path, err
}

func (s *AppService) OpenFolderDialog() (string, error) {
	path, err := s.app.Dialog.OpenFile().
		AttachToWindow(s.focusedWindow()).
		SetTitle("Open Folder").
		CanChooseDirectories(true).
		CanChooseFiles(false).
		CanCreateDirectories(true).
		PromptForSingleSelection()
	// The user picked a folder to browse — the folder itself and every
	// file inside it (and its subdirectories) is now trusted.
	if err == nil && path != "" {
		if trustErr := s.trustDir(path); trustErr != nil {
			fmt.Fprintf(os.Stderr, "trustDir(%s) failed: %v\n", path, trustErr)
		}
	}
	return path, err
}

func (s *AppService) ReadFile(path string) (string, error) {
	root, rel, err := s.resolveTrusted(path)
	if err != nil {
		return "", err
	}
	data, err := root.ReadFile(rel)
	if err != nil {
		return "", err
	}
	trackRecentFile(path)
	return string(data), nil
}

func (s *AppService) WriteFile(path string, content string) error {
	root, rel, err := s.resolveTrusted(path)
	if err != nil {
		return err
	}
	return root.WriteFile(rel, []byte(content), 0644)
}

func (s *AppService) ListDirectory(path string) ([]FileInfo, error) {
	root, rel, err := s.resolveTrusted(path)
	if err != nil {
		return nil, err
	}
	// os.Root has no ReadDir method, but Root.FS() returns an fs.FS that
	// is itself scoped to the root. fs.ReadDir on it goes through the
	// same openat path that ReadFile uses, so symlink escapes are still
	// caught at the kernel layer.
	entries, err := fs.ReadDir(root.FS(), rel)
	if err != nil {
		return nil, err
	}
	files := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		lower := strings.ToLower(name)
		if entry.IsDir() || strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".markdown") {
			files = append(files, FileInfo{
				Name:  name,
				Path:  filepath.Join(path, name),
				IsDir: entry.IsDir(),
			})
		}
	}
	return files, nil
}

func (s *AppService) GetHomePath() string {
	home, _ := os.UserHomeDir()
	return home
}

func (s *AppService) CloseWindow() {
	if w := s.focusedWindow(); w != nil {
		AllowWindowClose(w.ID())
		w.Close()
	}
}

func (s *AppService) QuitApp() {
	if s.app != nil {
		s.app.Quit()
	}
}

func (s *AppService) RequestQuit() {
	if s.app == nil || s.quit == nil {
		return
	}
	windows := s.app.Window.GetAll()
	if len(windows) == 0 {
		s.app.Quit()
		return
	}
	queue := make([]quitWindow, 0, len(windows))
	for _, w := range windows {
		queue = append(queue, w)
	}
	s.quit.Begin(queue)
}

// ConfirmQuitWindow advances the quit coordinator for the given window. The
// window ID is supplied by the frontend, which receives it in the
// app:confirmQuitWindow event payload (see quitCoordinator.requestNextLocked).
// Earlier versions read the window from ctx.Value(application.WindowKey), but
// Wails 3 generic RPC calls do not populate that key, so the call silently
// no-op'd and the app never quit.
func (s *AppService) ConfirmQuitWindow(windowID uint) {
	if s.quit == nil {
		return
	}
	s.quit.Confirm(windowID)
}

func (s *AppService) CancelQuit() {
	if s.quit != nil {
		s.quit.Cancel()
	}
}

func (s *AppService) SetUILocale(locale string) {
	SetLocale(locale)
	if s.app != nil {
		buildMenuI18n(s.app)
	}
}

func (s *AppService) GetUILocale() string {
	return GetLocale()
}

func (s *AppService) RestartApp() {
	if s.app != nil {
		s.app.Quit()
	}
}

func (s *AppService) GetConfig() AppConfig {
	return LoadConfig()
}

func (s *AppService) SaveConfig(cfg AppConfig) error {
	return SaveConfig(cfg)
}

func (s *AppService) ShowSaveDialog(filename string) string {
	ms := getMenuStrings()
	done := make(chan string, 1)

	title := ms.unsavedTitle
	if filename != "" {
		title = fmt.Sprintf(`"%s" %s`, filename, ms.unsavedTitle)
	}

	dlg := s.app.Dialog.Question().
		SetTitle(title).
		SetMessage(ms.closeWithoutSaving).
		AttachToWindow(s.focusedWindow())

	dlg.AddButton("取消").OnClick(func() { done <- "cancel" }).SetAsCancel()
	dlg.AddButton("不保存").OnClick(func() { done <- "discard" })
	dlg.AddButton("保存").OnClick(func() { done <- "save" }).SetAsDefault()
	dlg.Show()

	return <-done
}

func (s *AppService) ExportPDF(htmlContent string, title string) (string, error) {
	path, err := s.app.Dialog.SaveFile().
		AttachToWindow(s.focusedWindow()).
		AddFilter("PDF Files", "*.pdf").
		SetFilename(title + ".pdf").
		PromptForSingleSelection()
	if err != nil || path == "" {
		return "", fmt.Errorf("cancelled")
	}

	if err := exportHTMLToPDF(htmlContent, path); err != nil {
		return "", err
	}
	return path, nil
}

func exportHTMLToPDF(htmlContent string, outputPath string) error {
	doc := document.NewDocument(document.PageSizeA4)
	doc.Info.Title = "fastmd export"

	elems, err := html.Convert(htmlContent, nil)
	if err != nil {
		return fmt.Errorf("failed to convert HTML: %w", err)
	}
	for _, e := range elems {
		doc.Add(e)
	}

	if err := doc.Save(outputPath); err != nil {
		return fmt.Errorf("failed to save PDF: %w", err)
	}
	return nil
}
