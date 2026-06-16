package core

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type menuIconName string

const (
	menuIconSidebar     menuIconName = "sidebar"
	menuIconTheme       menuIconName = "theme"
	menuIconFullscreen  menuIconName = "fullscreen"
	menuIconDevTools    menuIconName = "devtools"
	menuIconNewFile     menuIconName = "new-file"
	menuIconNewWindow   menuIconName = "new-window"
	menuIconOpen        menuIconName = "open"
	menuIconSave        menuIconName = "save"
	menuIconSaveAs      menuIconName = "save-as"
	menuIconExportHTML  menuIconName = "export-html"
	menuIconExportPDF   menuIconName = "export-pdf"
	menuIconQuit        menuIconName = "quit"
	menuIconAbout       menuIconName = "about"
	menuIconPreferences menuIconName = "preferences"
	menuIconQuickStart  menuIconName = "quick-start"
	menuIconShortcuts   menuIconName = "shortcuts"
	menuIconMarkdown    menuIconName = "markdown"
	menuIconMath        menuIconName = "math"
)

var menuIconCache = map[menuIconName][]byte{}

func setMenuIcon(item *application.MenuItem, name menuIconName) *application.MenuItem {
	if item != nil {
		item.SetBitmap(menuIcon(name))
	}
	return item
}

func menuIcon(name menuIconName) []byte {
	if bitmap, ok := menuIconCache[name]; ok {
		return bitmap
	}

	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	ink := color.RGBA{R: 62, G: 68, B: 76, A: 255}
	accent := color.RGBA{R: 57, G: 139, B: 255, A: 255}
	muted := color.RGBA{R: 142, G: 150, B: 160, A: 255}

	switch name {
	case menuIconSidebar:
		fillRect(img, 2, 3, 14, 13, ink)
		fillRect(img, 3, 4, 6, 12, muted)
		fillRect(img, 7, 4, 13, 12, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	case menuIconTheme:
		fillRect(img, 3, 3, 8, 8, accent)
		fillRect(img, 8, 3, 13, 8, color.RGBA{R: 255, G: 185, B: 60, A: 255})
		fillRect(img, 3, 8, 8, 13, color.RGBA{R: 68, G: 195, B: 120, A: 255})
		fillRect(img, 8, 8, 13, 13, color.RGBA{R: 211, G: 92, B: 255, A: 255})
	case menuIconFullscreen:
		fillRect(img, 3, 3, 7, 5, ink)
		fillRect(img, 3, 3, 5, 7, ink)
		fillRect(img, 9, 3, 13, 5, ink)
		fillRect(img, 11, 3, 13, 7, ink)
		fillRect(img, 3, 11, 7, 13, ink)
		fillRect(img, 3, 9, 5, 13, ink)
		fillRect(img, 9, 11, 13, 13, ink)
		fillRect(img, 11, 9, 13, 13, ink)
	case menuIconDevTools:
		fillRect(img, 2, 4, 14, 12, ink)
		fillRect(img, 3, 5, 13, 11, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		fillRect(img, 5, 7, 7, 9, accent)
		fillRect(img, 8, 7, 11, 9, muted)
	case menuIconNewFile:
		drawDocument(img, ink, accent)
		fillRect(img, 10, 3, 12, 9, accent)
		fillRect(img, 8, 5, 14, 7, accent)
	case menuIconNewWindow:
		fillRect(img, 2, 4, 11, 12, muted)
		fillRect(img, 3, 5, 10, 11, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		fillRect(img, 5, 2, 14, 10, ink)
		fillRect(img, 6, 3, 13, 9, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	case menuIconOpen:
		fillRect(img, 2, 5, 7, 7, ink)
		fillRect(img, 2, 7, 14, 12, color.RGBA{R: 255, G: 185, B: 60, A: 255})
		fillRect(img, 5, 4, 14, 7, ink)
		fillRect(img, 3, 8, 13, 11, color.RGBA{R: 255, G: 209, B: 92, A: 255})
	case menuIconSave:
		drawFloppy(img, ink, accent)
	case menuIconSaveAs:
		drawFloppy(img, ink, muted)
		fillRect(img, 11, 10, 13, 14, accent)
		fillRect(img, 10, 11, 14, 13, accent)
	case menuIconExportHTML:
		drawDocument(img, ink, accent)
		fillRect(img, 5, 7, 7, 9, accent)
		fillRect(img, 9, 7, 11, 9, accent)
		fillRect(img, 7, 6, 9, 10, muted)
	case menuIconExportPDF:
		drawDocument(img, ink, color.RGBA{R: 230, G: 75, B: 75, A: 255})
		fillRect(img, 5, 7, 11, 9, color.RGBA{R: 230, G: 75, B: 75, A: 255})
		fillRect(img, 5, 10, 9, 12, color.RGBA{R: 230, G: 75, B: 75, A: 255})
	case menuIconQuit:
		fillRect(img, 7, 2, 9, 8, ink)
		fillRect(img, 4, 5, 6, 11, color.RGBA{R: 230, G: 75, B: 75, A: 255})
		fillRect(img, 10, 5, 12, 11, color.RGBA{R: 230, G: 75, B: 75, A: 255})
		fillRect(img, 5, 11, 11, 13, color.RGBA{R: 230, G: 75, B: 75, A: 255})
	case menuIconAbout:
		// Filled circle with a knocked-out "i" in the middle.
		fillRect(img, 6, 2, 10, 3, ink)
		fillRect(img, 4, 3, 12, 4, ink)
		fillRect(img, 3, 4, 13, 5, ink)
		fillRect(img, 2, 5, 14, 11, ink)
		fillRect(img, 3, 11, 13, 12, ink)
		fillRect(img, 4, 12, 12, 13, ink)
		fillRect(img, 6, 13, 10, 14, ink)
		fillRect(img, 7, 4, 9, 5, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		fillRect(img, 7, 6, 9, 7, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		fillRect(img, 7, 8, 9, 12, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	case menuIconPreferences:
		// Three horizontal sliders, each with a knob in a different
		// position. Recognizable at 16x16 and reads as "settings".
		fillRect(img, 2, 4, 14, 5, muted)
		fillRect(img, 9, 3, 12, 7, accent)
		fillRect(img, 2, 8, 14, 9, muted)
		fillRect(img, 4, 7, 7, 11, accent)
		fillRect(img, 2, 12, 14, 13, muted)
		fillRect(img, 10, 11, 13, 15, accent)
	case menuIconQuickStart:
		fillRect(img, 3, 3, 8, 13, ink)
		fillRect(img, 8, 3, 13, 13, muted)
		fillRect(img, 4, 5, 7, 6, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		fillRect(img, 9, 5, 12, 6, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		fillRect(img, 4, 8, 7, 9, accent)
	case menuIconShortcuts:
		fillRect(img, 2, 4, 14, 12, ink)
		fillRect(img, 3, 5, 13, 11, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		fillRect(img, 4, 6, 6, 8, muted)
		fillRect(img, 7, 6, 9, 8, muted)
		fillRect(img, 10, 6, 12, 8, muted)
		fillRect(img, 5, 9, 11, 10, accent)
	case menuIconMarkdown:
		drawDocument(img, ink, accent)
		fillRect(img, 4, 7, 6, 11, accent)
		fillRect(img, 10, 7, 12, 11, accent)
		fillRect(img, 6, 9, 10, 11, accent)
	case menuIconMath:
		drawDocument(img, ink, color.RGBA{R: 211, G: 92, B: 255, A: 255})
		fillRect(img, 5, 6, 11, 8, color.RGBA{R: 211, G: 92, B: 255, A: 255})
		fillRect(img, 7, 4, 9, 10, color.RGBA{R: 211, G: 92, B: 255, A: 255})
		fillRect(img, 5, 12, 11, 13, muted)
	}

	var out bytes.Buffer
	if err := png.Encode(&out, img); err != nil {
		return nil
	}
	bitmap := out.Bytes()
	menuIconCache[name] = bitmap
	return bitmap
}

func fillRect(img draw.Image, x0, y0, x1, y1 int, c color.Color) {
	draw.Draw(img, image.Rect(x0, y0, x1, y1), image.NewUniform(c), image.Point{}, draw.Src)
}

func drawDocument(img draw.Image, ink, accent color.Color) {
	fillRect(img, 4, 2, 11, 14, ink)
	fillRect(img, 5, 3, 10, 13, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	fillRect(img, 9, 2, 12, 5, mutedWhite())
	fillRect(img, 5, 6, 10, 7, accent)
	fillRect(img, 5, 9, 9, 10, accent)
}

func drawFloppy(img draw.Image, ink, accent color.Color) {
	fillRect(img, 3, 3, 13, 13, ink)
	fillRect(img, 5, 4, 11, 7, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	fillRect(img, 6, 10, 11, 13, accent)
	fillRect(img, 7, 11, 10, 13, color.RGBA{R: 255, G: 255, B: 255, A: 255})
}

func mutedWhite() color.Color {
	return color.RGBA{R: 236, G: 239, B: 243, A: 255}
}
