# Top Border Double-Click Toggle Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Double-clicking the window's top resize border toggles the window between its current size and the full visible screen frame (excluding menu bar and Dock).

**Architecture:** Add a `setupTopBorderDoubleClick` Objective-C function in `devtools_darwin.go` that registers an `NSEvent` local monitor. Per-window state (saved frame + expanded flag) is stored via associated objects on the `NSWindow` instance so the pattern scales to multiple windows. Call the setup function once from the existing `WindowShow` hook in `main.go`.

**Tech Stack:** Go, CGO, Objective-C, Cocoa (NSEvent, NSWindow, NSScreen)

---

### Task 1: Add CGO function and Go wrapper in devtools_darwin.go

**Files:**
- Modify: `devtools_darwin.go`

- [ ] **Step 1: Read the current file**

Read `devtools_darwin.go` to confirm the existing CGO preamble ends before `import "C"`.

- [ ] **Step 2: Add Objective-C implementation**

Inside the `/* ... */` CGO comment block (before `import "C"`), append after the existing `positionTrafficLights` function:

```objc
// Keys for per-window associated objects
static const char kSavedFrameKey = 0;
static const char kExpandedKey   = 0;

void setupTopBorderDoubleClick(void *nswin) {
    static BOOL setup = NO;
    if (setup) return;
    setup = YES;

    NSWindow *window = (NSWindow *)nswin;
    [NSEvent addLocalMonitorForEventsMatchingMask:NSEventMaskLeftMouseDown handler:^NSEvent *(NSEvent *event) {
        if (event.window != window || event.clickCount != 2) return event;

        NSPoint loc = [NSEvent mouseLocation];
        NSRect  frame = window.frame;
        CGFloat top   = frame.origin.y + frame.size.height;

        if (loc.x < frame.origin.x || loc.x > frame.origin.x + frame.size.width) return event;
        if (loc.y < top - 5        || loc.y > top + 2)                            return event;

        NSValue *savedVal = objc_getAssociatedObject(window, &kSavedFrameKey);
        BOOL isExpanded   = [objc_getAssociatedObject(window, &kExpandedKey) boolValue];

        if (isExpanded && savedVal) {
            [window setFrame:savedVal.rectValue display:YES animate:YES];
            objc_setAssociatedObject(window, &kExpandedKey, @NO, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
        } else {
            objc_setAssociatedObject(window, &kSavedFrameKey,
                [NSValue valueWithRect:frame], OBJC_ASSOCIATION_RETAIN_NONATOMIC);
            objc_setAssociatedObject(window, &kExpandedKey, @YES, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
            [window setFrame:window.screen.visibleFrame display:YES animate:YES];
        }
        return nil;
    }];
}
```

- [ ] **Step 3: Add the required import for associated objects**

At the top of the CGO preamble (after the existing `#import <Cocoa/Cocoa.h>` line), add:

```objc
#import <objc/runtime.h>
```

- [ ] **Step 4: Add Go wrapper**

After the existing `func positionTrafficLights(...)` Go function, add:

```go
func setupTopBorderDoubleClick(nsWindow unsafe.Pointer) {
	C.setupTopBorderDoubleClick(nsWindow)
}
```

- [ ] **Step 5: Build to verify compilation**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md && go build ./...
```

Expected: no errors.

- [ ] **Step 6: Commit**

```bash
git add devtools_darwin.go
git commit -m "feat: add top-border double-click toggle (CGO)"
```

---

### Task 2: Wire up the setup call in main.go

**Files:**
- Modify: `main.go`

- [ ] **Step 1: Read main.go**

Locate the `WindowShow` hook (around line 68):

```go
window.RegisterHook(events.Mac.WindowShow, func(_ *application.WindowEvent) {
    positionTrafficLights(window.NativeWindow(), 13, 14)
})
```

- [ ] **Step 2: Add the setup call**

Replace that hook body with:

```go
window.RegisterHook(events.Mac.WindowShow, func(_ *application.WindowEvent) {
    positionTrafficLights(window.NativeWindow(), 13, 14)
    setupTopBorderDoubleClick(window.NativeWindow())
})
```

- [ ] **Step 3: Build to verify**

```bash
cd /Users/yanshili/Downloads/md-p-1/fast-md && go build ./...
```

Expected: no errors.

- [ ] **Step 4: Manual smoke test**

Run the app (`task dev` or `./dev.sh`), then:
1. Double-click the very top edge of the window → window expands to fill the screen (minus menu bar and Dock).
2. Double-click the top edge again → window returns to its original size and position.
3. Repeat to confirm the toggle is stable.

- [ ] **Step 5: Commit**

```bash
git add main.go
git commit -m "feat: wire top-border double-click toggle"
```
