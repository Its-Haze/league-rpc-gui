// Command league-rpc-gui runs the daemon and the desktop GUI in one process.
package main

import (
	"context"
	"log"

	"github.com/its-haze/league-rpc/frontend"
	"github.com/its-haze/league-rpc/internal/app"
	"github.com/its-haze/league-rpc/internal/config"
	"github.com/its-haze/league-rpc/internal/daemon"
	"github.com/its-haze/league-rpc/internal/logging"
	"github.com/wailsapp/wails/v3/pkg/application"
)

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
	go func() {
		defer close(daemonDone)
		d.Run(ctx)
	}()

	svc := newGUIService(app.New(store))

	wailsApp := application.New(application.Options{
		Name:        "League RPC",
		Description: "League of Legends Discord Rich Presence",
		Services:    []application.Service{application.NewService(svc)},
		Assets:      application.AssetOptions{Handler: application.AssetFileServerFS(frontend.Assets())},
		OnShutdown: func() {
			cancel()
			<-daemonDone
		},
	})

	// Bridge config.Store changes to a frontend event so screens can react
	// without polling. Runs until the app shuts down.
	go svc.publishConfigChanges(ctx, wailsApp)

	wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "League RPC",
		Width:            960,
		Height:           640,
		BackgroundColour: application.NewRGB(15, 17, 23),
		URL:              "/",
	})

	if err := wailsApp.Run(); err != nil {
		log.Fatal(err)
	}
}
