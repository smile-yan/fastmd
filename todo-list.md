# fast-md todo list

Issues / followups discovered during the `fastmd` CLI install + rename work (commits `82441fe` and `94f9ca4`).

## P0 — data preservation

- [ ] **Migrate `~/Library/Application Support/fast-md/` → `fastmd/` on first launch after upgrade.**
  Renaming the support dir (commit `94f9ca4`) silently invalidates every preview-6 user's `config.json` (language, theme, autosave) and `recent.json` (dock menu). Acceptable for preview → preview, **not** acceptable for a 1.0 release. The path is derived from the literal `appName` constant in `core/app.go:32` and `core/recent_files.go:147`; on startup, detect the old dir, copy contents, then either rename it in place or read-then-rewrite under the new path. A versioned marker file (e.g. `~/Library/Application Support/fastmd/.migrated-from-fast-md`) prevents re-running the migration on every launch.

## P1 — inconsistencies to clean up before 1.0

- [ ] **Test identifier `fast-md-test` in `core/main_test.go` (10 occurrences).** Cosmetic — they're internal Wails `application.Options.Name` strings, no runtime impact — but they're the last in-tree `fast-md` literal. Update to `fastmd-test` for consistency with the rename. (Deliberately skipped in `94f9ca4` because the plan said "internal test identifier, no user impact". Reconsider once we're closer to 1.0.)

- [ ] **Frontend localStorage keys: `fast-md-settings`, `fast-md-theme`, `fast-md-locale`, `fast-md-settings-changed`.** Kept on purpose in `94f9ca4` to preserve users' saved settings across the rename. New users will see `fastmd` everywhere except these internal keys, which is mildly confusing. **Decision needed**: either (a) keep forever, documenting why in a comment, or (b) rename and add a one-shot migration in `useFile.ts` / `useTheme.ts` / `useLocale.ts` to read the old key, write the new, and delete the old. Option (a) is simpler and the keys are invisible.

- [ ] **DMG volume name in CI** (`-volname "fastmd"` in `.github/workflows/release.yml:171`) is now 6 chars and may collide with another app's mounted volume. macOS handles duplicates by appending `-1`, `-2`, etc. Not a real problem, but worth knowing.

## P2 — followups from the planning

- [ ] **Wails `SingleInstance` (P2 from the previous plan).** Confirmed available in v3 alpha.96 (`application.SingleInstanceOptions` with `OnSecondInstanceLaunch` receiving `SecondInstanceData{Args, WorkingDir}`). Not enabled because `open -a` already routes via Apple Events, so the current UX is correct. Enable later if the user wants true single-process behavior (e.g. quit-on-last-window-close semantics become meaningful). Document the trade-off in a comment near `core/run.go:202-214`.

- [ ] **Plan accuracy: tests with hardcoded renamed paths.** The plan for `94f9ca4` marked `core/config_test.go` as "not user-facing, don't modify", but line 99 hardcoded the support-dir path and broke the test. The plan's grep should have caught this. For future renames, run `grep -rn '<old-literal>' **/*_test.go` separately and check whether each hit is a constant compared against runtime output (must change) vs. a comment / test name (can stay).

- [ ] **Plan accuracy: `build/ios/build.sh`.** Not in the initial grep output, only caught in the post-implementation sanity grep. The original grep excluded `.sh` files inside `build/*/` because the find pattern was scoped to the top level. Loosen the grep include to `**/*.sh` (with a few excludes) for the next cross-platform sweep.

- [ ] **`scripts/install-cli.sh` had a command-substitution bug** (backticks inside a `log "..."` argument). Fixed in `82441fe` but worth a `shellcheck` pass on the script directory before the next release. The repo currently has no shell-script linting in CI.

## Already-shipped behaviors worth documenting in CLAUDE.md

- **`open -a fastmd` is the CLI entry point** (not direct binary invocation). macOS LaunchServices routes the kAEOpenDocuments Apple Event to the running instance — no Wails `SingleInstance` needed. Already explained in the `os.Args` fallback comment at `core/run.go:223-234`; consider promoting it to CLAUDE.md so future contributors don't re-enable `SingleInstance` redundantly.

- **Wrapper script lives in `Contents/Resources/fastmd`, not `Contents/MacOS/`.** Deliberately chosen to keep the Go binary's ad-hoc codesign clean (Resources aren't part of `--deep` signing). Documented in `scripts/fastmd-wrapper.sh:2-4`; CLAUDE.md could mention it under the macOS-only features section.
