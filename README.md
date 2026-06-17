# fastmd

fastmd is a lightweight Markdown editor built with Wails 3, Go, Vue, TypeScript, and Milkdown.

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
fastmd ~/notes/hello.md
```

For the short `fastmd <file>.md` workflow, run the install script once
after dragging `fastmd.app` to `/Applications`:

```bash
bash scripts/install-cli.sh
```

It symlinks a tiny wrapper at `fastmd.app/Contents/Resources/fastmd`
into the first writable `$PATH` directory it finds
(`/opt/homebrew/bin` → `/usr/local/bin` → `~/.local/bin`) as `fastmd`.
The wrapper runs `open -a fastmd "$@"`, which routes the file through
macOS LaunchServices — if fastmd is already running, the file opens in
a new window of the running instance; otherwise fastmd cold-launches.
The script is idempotent and never needs `sudo` (it falls back to
`~/.local/bin` if neither Homebrew prefix is writable).

After install, any `.md` or `.markdown` argument is opened in a fresh
window. A leading `~/` is expanded via `$HOME`; relative paths are
resolved against the current working directory. Flags and unrelated
arguments are ignored, so a stray `-psn_0` from Launch Services never
confuses the parser.

If you don't want to install the wrapper, the built-in macOS opener
works too:

```bash
open -a fastmd ~/notes/hello.md
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
GOCACHE=/private/tmp/fastmd-go-cache go test ./...
GOCACHE=/private/tmp/fastmd-go-cache go build -o /private/tmp/fastmd-window-flow-test .
```

Run `npm run build` before Go tests or builds when `frontend/dist` has changed, because Go embeds that directory.
