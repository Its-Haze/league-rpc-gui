package main

import "sync"

// windowController is the subset of the Wails window the tray drives. An
// adapter over *application.WebviewWindow satisfies it; tests use a fake.
type windowController interface {
	Show()
	Hide()
	Focus()
}

// pauseControl is the runtime pause surface the tray toggles.
type pauseControl interface {
	SetPaused(bool)
	IsPaused() bool
}

// trayController holds the window and pause wiring behind the tray menu,
type trayController struct {
	win        windowController
	pause      pauseControl
	firstHide  sync.Once
	notifyHide func()
	// reflectChecked mirrors the pause flag onto the tray menu checkbox.
	// Set by the Wails wiring once the menu item exists.
	reflectChecked func(bool)
}

func newTrayController(win windowController, pause pauseControl, notifyHide func()) *trayController {
	return &trayController{win: win, pause: pause, notifyHide: notifyHide}
}

// showWindow brings the window back and focuses it. Used by the tray icon
// click, the "Open" menu item, and a second launch of the app.
func (c *trayController) showWindow() {
	c.win.Show()
	c.win.Focus()
}

// handleClose redirects a window close to a hide so the app keeps running.
// The first time it fires it also shows the "still running" message once.
func (c *trayController) handleClose() {
	c.win.Hide()
	c.firstHide.Do(func() {
		if c.notifyHide != nil {
			c.notifyHide()
		}
	})
}

// setPaused is the one place pause changes flow through, so the daemon flag
// and the tray checkbox never disagree, whichever surface triggered it.
func (c *trayController) setPaused(paused bool) {
	c.pause.SetPaused(paused)
	if c.reflectChecked != nil {
		c.reflectChecked(paused)
	}
}

// togglePause flips the pause flag and returns the new value so a tray-menu
// click can update the checkbox it was clicked on.
func (c *trayController) togglePause() bool {
	next := !c.pause.IsPaused()
	c.setPaused(next)
	return next
}
