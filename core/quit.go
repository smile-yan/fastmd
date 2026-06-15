package core

import "sync"

const quitConfirmWindowEvent = "app:confirmQuitWindow"

type quitWindow interface {
	ID() uint
	EmitEvent(name string, data ...any) bool
	Focus()
	Close()
}

type quitCoordinator struct {
	mu         sync.Mutex
	queue      []quitWindow
	current    quitWindow
	active     bool
	allowClose func(uint)
	quit       func()
}

func newQuitCoordinator(allowClose func(uint), quit func()) *quitCoordinator {
	return &quitCoordinator{
		allowClose: allowClose,
		quit:       quit,
	}
}

func (q *quitCoordinator) Begin(windows []quitWindow) {
	q.mu.Lock()
	if q.active {
		q.mu.Unlock()
		return
	}
	q.active = true
	q.queue = append([]quitWindow(nil), windows...)
	q.current = nil
	q.requestNextLocked()
	q.mu.Unlock()
}

func (q *quitCoordinator) Confirm(windowID uint) {
	q.mu.Lock()
	if q.current == nil || q.current.ID() != windowID {
		q.mu.Unlock()
		return
	}

	window := q.current
	q.current = nil
	if q.allowClose != nil {
		q.allowClose(window.ID())
	}
	q.mu.Unlock()

	window.Close()

	q.mu.Lock()
	q.requestNextLocked()
	q.mu.Unlock()
}

func (q *quitCoordinator) Cancel() {
	q.mu.Lock()
	q.queue = nil
	q.current = nil
	q.active = false
	q.mu.Unlock()
}

func (q *quitCoordinator) requestNextLocked() {
	if q.current != nil {
		return
	}

	for len(q.queue) > 0 {
		next := q.queue[0]
		q.queue = q.queue[1:]
		if next == nil {
			continue
		}
		q.current = next
		next.Focus()
		// Pass the window ID in the event payload so the frontend can echo it
		// back to ConfirmQuitWindow. Previously the JS side called the binding
		// with no arguments and relied on Go pulling the window out of
		// ctx.Value(application.WindowKey) — that key is not populated for
		// generic RPC calls, so the confirmation was silently dropped and the
		// quit never advanced.
		next.EmitEvent(quitConfirmWindowEvent, next.ID())
		return
	}

	if q.quit != nil {
		q.active = false
		q.quit()
	}
}
