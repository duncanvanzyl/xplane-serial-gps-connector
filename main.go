package main

import (
	"context"
	"log/slog"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/duncanvanzyl/xplane-serial-gps-connector/outputters"
	"github.com/duncanvanzyl/xplane-serial-gps-connector/serial"
	"github.com/duncanvanzyl/xplane-serial-gps-connector/xplane"
)

func main() {
	// Create the logger
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// Create the app
	a := &App{
		Serial: serial.NewSerial([]outputters.Outputter{
			&outputters.GGA{},
			&outputters.VTG{},
		}),
		Logger: logger,
	}

	// Set the serial and xplanes Loggers
	xplane.Logger = logger.With("src", "XPlane")
	serial.Logger = logger.With("src", "Serial")

	// Create the UI
	gui := app.New()
	w := gui.NewWindow("X-Plane GPS Simulator")
	ui := NewAppUI(a, logger.With("src", "AppUI"))

	// watch the app for changes to show in the UI
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go ui.Watch(ctx)

	// Start the UI
	w.SetContent(ui.GetContent())
	w.SetMainMenu(fyne.NewMainMenu(ui.SettingsMenu(w)))
	w.Resize(fyne.NewSize(400, 100))
	w.ShowAndRun()
	if ui.cancelCtx != nil {
		ui.cancelCtx()
	}
}
