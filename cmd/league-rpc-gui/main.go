// Command league-rpc-gui runs the daemon and the desktop GUI in one process.
package main

import (
	"context"
	"log"
	"runtime"

	"github.com/its-haze/league-rpc/frontend"
	"github.com/its-haze/league-rpc/internal/app"
	"github.com/its-haze/league-rpc/internal/config"
	"github.com/its-haze/league-rpc/internal/daemon"
	"github.com/its-haze/league-rpc/internal/logging"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

// singleInstanceID keeps a second launch from starting its own daemon; it
// signals the running instance to surface its window instead.
const singleInstanceID = "com.its-haze.league-rpc"

func main() {
	cfg, err := config.LoadOrCreate()
	if err != nil {
		log.Fatalf("load configuration: %v", err)
	}

	sink, err := logging.New(logging.Options{Debug: cfg.Advanced.DebugMode})
	if err != nil {
		log.Fatalf("init logging: %v", err)
	}
	defer sink.Close()

	store := config.NewStore(cfg)
	d := daemon.Wire(store, sink.Logger)

	ctx, cancel := context.WithCancel(context.Background())
	daemonDone := make(chan struct{})

	svc := newGUIService(app.New(store, d))

	// Assigned once the window exists; the single-instance callback reads it.
	var mainWindow *application.WebviewWindow

	wailsApp := application.New(application.Options{
		Name:        "League RPC",
		Description: "League of Legends Discord Rich Presence",
		Services:    []application.Service{application.NewService(svc)},
		Assets:      application.AssetOptions{Handler: application.AssetFileServerFS(frontend.Assets())},
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: singleInstanceID,
			OnSecondInstanceLaunch: func(application.SecondInstanceData) {
				if mainWindow != nil {
					mainWindow.Show()
					mainWindow.Focus()
				}
			},
		},
		OnShutdown: func() {
			cancel()
			<-daemonDone
		},
	})

	// Start the daemon only after application.New has run the single-instance
	// guard; a second launch exits inside New and never touches Discord.
	go func() {
		defer close(daemonDone)
		d.Run(ctx)
	}()

	// Bridge config.Store changes to a frontend event so screens can react
	// without polling. Runs until the app shuts down.
	go svc.publishConfigChanges(ctx, wailsApp)

	mainWindow = wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "League RPC",
		Width:            960,
		Height:           640,
		BackgroundColour: application.NewRGB(15, 17, 23),
		URL:              "/",
	})

	tray := newTrayController(
		windowAdapter{mainWindow},
		d,
		func() {
			wailsApp.Dialog.Info().
				SetTitle("League RPC is still running").
				SetMessage("Closing the window keeps presence running in the background. Right-click the tray icon and choose Quit to stop it.").
				Show()
		},
	)

	// Close the window -> hide it, keep the daemon running. Only tray Quit exits.
	mainWindow.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		tray.handleClose()
		e.Cancel()
	})

	systemTray := wailsApp.SystemTray.New()
	systemTray.SetIcon(icons.SystrayLight)
	systemTray.SetDarkModeIcon(icons.SystrayDark)
	systemTray.SetTooltip("League RPC")
	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
	}
	systemTray.OnClick(tray.showWindow)
	systemTray.SetMenu(buildTrayMenu(wailsApp, tray))

	// Frontend pause toggles flow through the same path as tray toggles, so
	// the daemon flag and the tray checkbox stay in agreement.
	svc.pauseHook = tray.setPaused

	if err := wailsApp.Run(); err != nil {
		log.Fatal(err)
	}
}

// buildTrayMenu assembles the right-click menu: Open, Pause presence, Quit.
func buildTrayMenu(wailsApp *application.App, tray *trayController) *application.Menu {
	menu := wailsApp.NewMenu()
	menu.Add("Open").OnClick(func(*application.Context) { tray.showWindow() })

	pauseItem := menu.AddCheckbox("Pause presence", tray.pause.IsPaused())
	// A frontend toggle runs off the main thread; marshal the native menu update.
	tray.reflectChecked = func(paused bool) {
		application.InvokeAsync(func() { pauseItem.SetChecked(paused) })
	}
	pauseItem.OnClick(func(*application.Context) { tray.togglePause() })

	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(*application.Context) { wailsApp.Quit() })
	return menu
}

// windowAdapter fits *application.WebviewWindow to the tray's windowController;
// the window's Show/Hide return a value the interface does not.
type windowAdapter struct{ w *application.WebviewWindow }

func (a windowAdapter) Show()  { a.w.Show() }
func (a windowAdapter) Hide()  { a.w.Hide() }
func (a windowAdapter) Focus() { a.w.Focus() }
