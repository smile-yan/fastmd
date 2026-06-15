package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/carlos7ags/folio/document"
	"github.com/carlos7ags/folio/html"
	"github.com/wailsapp/wails/v3/pkg/application"
)

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
	return filepath.Join(home, "Library", "Application Support", "fast-md", "config.json")
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
	return d.PromptForSingleSelection()
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
	return d.PromptForSingleSelection()
}

func (s *AppService) OpenFolderDialog() (string, error) {
	return s.app.Dialog.OpenFile().
		AttachToWindow(s.focusedWindow()).
		SetTitle("Open Folder").
		CanChooseDirectories(true).
		CanChooseFiles(false).
		CanCreateDirectories(true).
		PromptForSingleSelection()
}

func (s *AppService) ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	trackRecentFile(path)
	return string(data), nil
}

func (s *AppService) WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func (s *AppService) ListDirectory(path string) ([]FileInfo, error) {
	entries, err := os.ReadDir(path)
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
	doc.Info.Title = "fast-md export"

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
