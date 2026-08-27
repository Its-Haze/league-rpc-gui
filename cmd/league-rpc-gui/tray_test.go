package main

import (
	"slices"
	"testing"
)

type fakeWindow struct {
	shown   int
	hidden  int
	focused int
}

func (f *fakeWindow) Show()  { f.shown++ }
func (f *fakeWindow) Hide()  { f.hidden++ }
func (f *fakeWindow) Focus() { f.focused++ }

type fakePause struct{ paused bool }

func (f *fakePause) SetPaused(p bool) { f.paused = p }
func (f *fakePause) IsPaused() bool   { return f.paused }

func TestTrayController_CloseHidesAndDoesNotQuit(t *testing.T) {
	win := &fakeWindow{}
	notifs := 0
	c := newTrayController(win, &fakePause{}, func() { notifs++ })

	c.handleClose()
	c.handleClose()
	c.handleClose()

	if win.hidden != 3 {
		t.Fatalf("expected every close to hide the window, got %d hides", win.hidden)
	}
	if win.shown != 0 {
		t.Fatal("closing the window must not show or recreate it")
	}
	if notifs != 1 {
		t.Fatalf("first-hide message must fire exactly once, fired %d times", notifs)
	}
}

func TestTrayController_ShowWindowShowsAndFocuses(t *testing.T) {
	win := &fakeWindow{}
	c := newTrayController(win, &fakePause{}, nil)

	c.showWindow()

	if win.shown != 1 || win.focused != 1 {
		t.Fatalf("showWindow should show and focus once: shown=%d focused=%d", win.shown, win.focused)
	}
}

func TestTrayController_TogglePauseFlipsFlag(t *testing.T) {
	pause := &fakePause{}
	c := newTrayController(&fakeWindow{}, pause, nil)

	if got := c.togglePause(); !got || !pause.paused {
		t.Fatalf("first toggle should pause: returned=%v flag=%v", got, pause.paused)
	}
	if got := c.togglePause(); got || pause.paused {
		t.Fatalf("second toggle should unpause: returned=%v flag=%v", got, pause.paused)
	}
}

func TestTrayController_SetPausedMirrorsToCheckbox(t *testing.T) {
	pause := &fakePause{}
	c := newTrayController(&fakeWindow{}, pause, nil)

	var checked []bool
	c.reflectChecked = func(p bool) { checked = append(checked, p) }

	c.setPaused(true) // frontend-style call, not via the menu
	c.togglePause()   // menu-style call, now unpauses

	if !slices.Equal(checked, []bool{true, false}) {
		t.Fatalf("checkbox did not mirror every pause change: %v", checked)
	}
	if pause.paused {
		t.Fatal("daemon flag should be unpaused after the toggle")
	}
}
