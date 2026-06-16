# fast-md

fast-md is a lightweight Markdown editor built with Wails 3, Go, Vue, TypeScript, and Milkdown.

## Features

- Typora-like Markdown editing experience on macOS.
- Multiple independent editor windows.
- Local file open, save, save as, and folder browsing.
- HTML and PDF export.
- GitHub-style Markdown content theme.
- Localized menus and settings.

## Command Line

Open a file from the terminal by pointing the binary at a `.md`/`.markdown`
path:

```bash
fast-md ~/notes/hello.md
```

On first install, macOS does not put `.app` bundles on `$PATH`. To get
the short `fastmd hello.md` workflow, symlink the bundle's binary into a
directory on your `$PATH`:

```bash
ln -sf /Applications/fast-md.app/Contents/MacOS/fast-md /usr/local/bin/fastmd
# Apple Silicon: /opt/homebrew/bin/fastmd
```

After that, any `.md` or `.markdown` argument is opened in a fresh window.
A leading `~/` is expanded via `$HOME`; relative paths are resolved
against the current working directory. Flags and unrelated arguments are
ignored, so a stray `-psn_0` from Launch Services never confuses the
parser.

If you don't want a symlink, the built-in macOS opener works too:

```bash
open -a fast-md ~/notes/hello.md
```

## Project Structure

- `app.go`, `main.go`, `quit.go`: Go application services, window lifecycle, and quit coordination.
- `menu_*.go`, `dockmenu_*.go`, `devtools_*.go`: macOS menu, dock menu, and developer tooling integration.
- `frontend/src`: Vue application, editor components, composables, styles, and tests.
- `frontend/public/themes`: Markdown content themes and theme assets.
- `frontend/bindings`: Wails-generated frontend bindings.
- `docs/help`: Markdown help documents opened from the Help menu.
- `build`: Wails platform build configuration.
- `tools`: maintenance utilities.

## Development

```bash
task dev
```

The Vite dev server uses port `9245` by default.

## Verification

Run frontend checks:

```bash
cd frontend
npm test
npm run build
```

Run Go checks from the repository root:

```bash
GOCACHE=/private/tmp/fast-md-go-cache go test ./...
GOCACHE=/private/tmp/fast-md-go-cache go build -o /private/tmp/fast-md-window-flow-test .
```

Run `npm run build` before Go tests or builds when `frontend/dist` has changed, because Go embeds that directory.
