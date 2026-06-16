package core

import (
	"io/fs"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// TestMain wires the bundled frontend/help assets into the package-level
// FS variables so tests that exercise menu building and help materialization
// don't have to set them up individually. The frontend FS keeps the
// "frontend/dist/..." prefix because loadThemes reads from that path; the
// help FS is rooted at "docs/help" because materializeHelpDocument reads
// bare filenames. Missing assets are silently ignored — those tests will
// skip or fail with a clear "assets not initialized" message.
func TestMain(m *testing.M) {
	if dir, err := os.Stat("../frontend/dist"); err == nil && dir.IsDir() {
		assetsFS = os.DirFS("..")
	}
	if dir, err := os.Stat("../docs/help"); err == nil && dir.IsDir() {
		if sub, err := fs.Sub(os.DirFS(".."), "docs/help"); err == nil {
			helpDocumentFS = sub
		}
	}
	os.Exit(m.Run())
}

type fakeFullscreenWindow struct {
	toggleCount int
}

func (w *fakeFullscreenWindow) ToggleFullscreen() {
	w.toggleCount++
}

func TestToggleFocusedFullscreenUsesNativeWindowFullscreen(t *testing.T) {
	window := &fakeFullscreenWindow{}

	ToggleFocusedFullscreen(window)

	if window.toggleCount != 1 {
		t.Fatalf("expected native fullscreen toggle to be called once, got %d", window.toggleCount)
	}
}

func TestToggleFocusedFullscreenIgnoresMissingWindow(t *testing.T) {
	ToggleFocusedFullscreen(nil)
}

func TestShortcutAcceleratorsMatchTyporaMacOS(t *testing.T) {
	app := application.New(application.Options{Name: "fast-md-test"})
	SetLocale("en")
	buildMenuI18n(app)

	menu := app.Menu.GetApplicationMenu()
	assertMenuAccelerator(t, menu, "New File", "Cmd+N")
	assertMenuAccelerator(t, menu, "New Window", "Cmd+Shift+N")
	assertMenuAccelerator(t, menu, "Toggle Sidebar", "Cmd+Shift+L")
	assertMenuAccelerator(t, menu, "Enter Full Screen", "Cmd+Ctrl+F")
	assertMenuAccelerator(t, menu, "Developer Tools", "")
}

func TestDeveloperToolsShortcutRegisteredAsGlobalKeyBinding(t *testing.T) {
	app := application.New(application.Options{Name: "fast-md-test"})
	RegisterDeveloperToolsShortcut(app)

	assertKeyBindingRegistered(t, app, "F12")
}

func TestViewMenuFullscreenItemStaysCustomAction(t *testing.T) {
	app := application.New(application.Options{Name: "fast-md-test"})
	SetLocale("en")
	buildMenuI18n(app)

	menu := app.Menu.GetApplicationMenu()
	if item := menu.FindByRole(application.ToggleFullscreen); item != nil {
		t.Fatalf("expected View fullscreen item to stay custom, got native role item %q", item.Label())
	}
}

func TestSystemMenuCleanerInstalledForEditorAndView(t *testing.T) {
	var got []string
	previous := systemMenuCleanerInstaller
	systemMenuCleanerInstaller = func(menuTitle string) {
		got = append(got, menuTitle)
	}
	defer func() {
		systemMenuCleanerInstaller = previous
	}()

	installConfiguredSystemMenuCleaners("Editor", "View")

	want := []string{"Editor", "View"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected system menu cleaners %#v, got %#v", want, got)
	}
}

func TestBuildMenuI18nInstallsSystemMenuCleanersForCurrentLocale(t *testing.T) {
	var got []string
	previous := systemMenuCleanerInstaller
	systemMenuCleanerInstaller = func(menuTitle string) {
		got = append(got, menuTitle)
	}
	defer func() {
		systemMenuCleanerInstaller = previous
	}()

	app := application.New(application.Options{Name: "fast-md-test"})
	SetLocale("zh")
	buildMenuI18n(app)

	want := []string{"编辑", "视图"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected localized system menu cleaners %#v, got %#v", want, got)
	}
}

func TestViewMenuItemsHaveIcons(t *testing.T) {
	app := application.New(application.Options{Name: "fast-md-test"})
	SetLocale("en")
	buildMenuI18n(app)

	viewItem := app.Menu.GetApplicationMenu().FindByLabel("View")
	if viewItem == nil {
		t.Fatal("expected View menu to exist")
	}
	assertMenuItemsHaveIcons(t, viewItem.GetSubmenu())
}

func TestViewThemeMenuOnlyContainsPackagedThemes(t *testing.T) {
	app := application.New(application.Options{Name: "fast-md-test"})
	SetLocale("en")
	buildMenuI18n(app)

	themeItem := app.Menu.GetApplicationMenu().FindByLabel("Theme")
	if themeItem == nil {
		t.Fatal("expected Theme menu to exist")
	}

	got := menuLabels(themeItem.GetSubmenu())
	want := []string{"Github"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected Theme menu labels %#v, got %#v", want, got)
	}
}

func TestFileMenuItemsHaveIcons(t *testing.T) {
	app := application.New(application.Options{Name: "fast-md-test"})
	SetLocale("en")
	buildMenuI18n(app)

	fileItem := app.Menu.GetApplicationMenu().FindByLabel("File")
	if fileItem == nil {
		t.Fatal("expected File menu to exist")
	}
	assertMenuItemsHaveIcons(t, fileItem.GetSubmenu())
}

func TestHelpMenuItemsHaveIcons(t *testing.T) {
	app := application.New(application.Options{Name: "fast-md-test"})
	SetLocale("en")
	buildMenuI18n(app)

	helpItem := app.Menu.GetApplicationMenu().FindByLabel("Help")
	if helpItem == nil {
		t.Fatal("expected Help menu to exist")
	}
	assertMenuItemsHaveIcons(t, helpItem.GetSubmenu())
}

func TestHelpMenuContainsOnlyDocumentationItems(t *testing.T) {
	app := application.New(application.Options{Name: "fast-md-test"})
	SetLocale("zh")
	buildMenuI18n(app)

	helpItem := app.Menu.GetApplicationMenu().FindByLabel("帮助")
	if helpItem == nil {
		t.Fatal("expected Help menu to exist")
	}

	got := menuLabels(helpItem.GetSubmenu())
	want := []string{"快速开始", "快捷键说明", "markdown入门", "数学公式入门"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected Help menu labels %#v, got %#v", want, got)
	}
}

func TestHelpDocumentsMaterializeMarkdownFiles(t *testing.T) {
	dir := t.TempDir()

	for _, doc := range helpDocumentEntries(getMenuStrings()) {
		path, err := materializeHelpDocument(doc.filename, dir)
		if err != nil {
			t.Fatalf("expected %s to materialize: %v", doc.filename, err)
		}
		if !strings.HasSuffix(path, doc.filename) {
			t.Fatalf("expected materialized path to end with %q, got %q", doc.filename, path)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("expected to read materialized help document %q: %v", path, err)
		}
		if !strings.HasPrefix(string(content), "# ") {
			t.Fatalf("expected %q to be a markdown document with a heading", doc.filename)
		}
	}
}

func assertMenuItemsHaveIcons(t *testing.T, menu *application.Menu) {
	t.Helper()
	if menu == nil {
		t.Fatal("expected submenu to exist")
	}
	for i := 0; ; i++ {
		item := menu.ItemAt(i)
		if item == nil {
			return
		}
		if item.IsSeparator() {
			continue
		}
		if !menuItemHasBitmap(item) {
			t.Fatalf("expected View menu item %q to have an icon", item.Label())
		}
		if submenu := item.GetSubmenu(); submenu != nil {
			assertMenuItemsHaveIcons(t, submenu)
		}
	}
}

func menuLabels(menu *application.Menu) []string {
	var labels []string
	for i := 0; ; i++ {
		item := menu.ItemAt(i)
		if item == nil {
			return labels
		}
		if item.IsSeparator() {
			continue
		}
		labels = append(labels, item.Label())
	}
}

func menuItemHasBitmap(item *application.MenuItem) bool {
	value := reflect.ValueOf(item).Elem().FieldByName("bitmap")
	return value.IsValid() && value.Kind() == reflect.Slice && value.Len() > 0
}

func assertMenuAccelerator(t *testing.T, menu *application.Menu, label string, want string) {
	t.Helper()
	item := menu.FindByLabel(label)
	if item == nil {
		t.Fatalf("expected menu item %q to exist", label)
	}
	if got := item.GetAccelerator(); got != want {
		t.Fatalf("expected %q accelerator %q, got %q", label, want, got)
	}
}

func assertKeyBindingRegistered(t *testing.T, app *application.App, accelerator string) {
	t.Helper()
	for _, binding := range app.KeyBinding.GetAll() {
		if binding.Accelerator == accelerator {
			return
		}
	}
	t.Fatalf("expected key binding %q to be registered", accelerator)
}

// Regression guard for the "second Finder Open With overwrites the first
// window" bug: macOS Finder "Open With" / double-click must always spawn
// a new editor window per request, even when another editor window is
// already focused. The ApplicationOpenedWithFile handler in run.go must
// therefore never be wired through a "route to current window" helper.
// If you re-introduce such a helper, the event handler must not call it.
func TestApplicationOpenedWithFileAlwaysSpawnsNewWindow(t *testing.T) {
	// Read the source as a string and assert the handler does not call
	// RouteOpenedFile / app.Window.Current() before opening a window.
	src, err := os.ReadFile("run.go")
	if err != nil {
		t.Fatalf("read run.go: %v", err)
	}
	text := string(src)

	// Locate the ApplicationOpenedWithFile handler.
	start := strings.Index(text, "OnApplicationEvent(events.Common.ApplicationOpenedWithFile")
	if start < 0 {
		t.Fatal("ApplicationOpenedWithFile handler not found in run.go")
	}
	end := strings.Index(text[start:], "})")
	if end < 0 {
		t.Fatal("could not find end of ApplicationOpenedWithFile handler")
	}
	handler := text[start : start+end+2]

	for _, banned := range []string{"RouteOpenedFile", "Window.Current()"} {
		if strings.Contains(handler, banned) {
			t.Errorf("ApplicationOpenedWithFile handler must not reference %q — every Finder "+
				"file open must spawn a new window, never overwrite the focused one", banned)
		}
	}
	if !strings.Contains(handler, "NewEditorWindowWithFile") {
		t.Error("ApplicationOpenedWithFile handler must call NewEditorWindowWithFile to open a new window")
	}
}
