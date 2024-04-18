package main

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/duncanvanzyl/xplane-serial-gps-connector/serial"
	"github.com/duncanvanzyl/xplane-serial-gps-connector/xplane"
)

// AppUI is the UI for the App
type AppUI struct {
	app               *App
	XPlanes           xplane.XPlanes
	xplaneSelect      *widget.Select
	xplaneRefresh     *widget.Button
	serialPortsSelect *widget.Select
	serialPortRefresh *widget.Button
	baudRate          *widget.Select
	refreshFreq       *widget.Select
	runButton         *widget.Button
	stopButton        *widget.Button
	status            *widget.Label
	cancelCtx         context.CancelFunc
	Logger            *slog.Logger
}

var (
	// PossibleBaudeRates is the list of possible baud rates
	PossibleBaudeRates = [...]string{"9600", "14400", "19200", "38400", "57600", "115200"}
	// PossiblePosFreqs is the list of possible position frequencies
	// This will determine the rate that the X-Plane position is read
	PossiblePosFreqs = [...]string{"1Hz", "2Hz", "5Hz", "10Hz", "20Hz"}
)

// NewAppUI returns a new AppUI
func NewAppUI(xApp *App, logger *slog.Logger) *AppUI {
	ui := &AppUI{
		app:     xApp,
		XPlanes: make(xplane.XPlanes),
		status:  widget.NewLabel(""),
		Logger:  logger,
	}

	ui.xplaneRefresh = widget.NewButton("Refresh X-Plane List", func() {
		go ui.findXplanes(5 * time.Second)
	})
	ui.xplaneSelect = widget.NewSelect([]string{}, ui.setXPlane(xApp))

	ui.serialPortRefresh = widget.NewButton("Refresh Serial Port List", ui.getSerial)
	ui.serialPortsSelect = widget.NewSelect([]string{}, ui.setSerialPort(xApp))
	ui.baudRate = widget.NewSelect(PossibleBaudeRates[:], func(value string) {
		ui.Logger.Debug("Set BaudRate", "baud", value)

		v, err := strconv.Atoi(value)
		if err != nil {
			ui.Logger.Error("Failed to convert BaudRate", "err", err)
			return
		}
		xApp.SetBaudRate(v)
	})

	ui.refreshFreq = widget.NewSelect(PossiblePosFreqs[:], func(value string) {
		ui.Logger.Debug("Set PositionFreq", "freq", value)
		v, err := strconv.Atoi(value[:len(value)-2])
		if err != nil {
			ui.Logger.Error("Failed to convert PositionFreq", "err", err)
			return
		}
		xApp.SetPositionFreq(uint(v))
	})

	ui.runButton = widget.NewButton("Run", ui.run)
	ui.runButton.Disable()
	ui.stopButton = widget.NewButton("Stop", ui.stop)

	// Populate Serial Port List, XPlane List and set default selects
	go ui.getSerial()
	go ui.findXplanes(5 * time.Second)
	ui.baudRate.SetSelected("38400")
	ui.refreshFreq.SetSelected("10Hz")
	ui.stopButton.Disable()

	return ui
}

// watchGUIState will set the GUI state based on the app state
func (ui *AppUI) watchGUIState(state, last_state AppState) {
	if state == last_state {
		return
	}
	switch state {
	case Running:
		ui.xplaneRefresh.Disable()
		ui.xplaneSelect.Disable()
		ui.serialPortRefresh.Disable()
		ui.serialPortsSelect.Disable()
		ui.baudRate.Disable()
		ui.refreshFreq.Disable()
		ui.runButton.Disable()
		ui.stopButton.Enable()
	case Runable:
		ui.xplaneRefresh.Enable()
		ui.xplaneSelect.Enable()
		ui.serialPortRefresh.Enable()
		ui.serialPortsSelect.Enable()
		ui.baudRate.Enable()
		ui.refreshFreq.Enable()
		ui.runButton.Enable()
		ui.stopButton.Disable()
	case Incomplete:
		ui.xplaneRefresh.Enable()
		ui.xplaneSelect.Enable()
		ui.serialPortRefresh.Enable()
		ui.serialPortsSelect.Enable()
		ui.baudRate.Enable()
		ui.refreshFreq.Enable()
		ui.runButton.Disable()
		ui.stopButton.Disable()
	}
}

// Watch will watch the app state and update the GUI accordingly
func (ui *AppUI) Watch(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	last_state := ui.app.State()
	last_xplane_count := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ui.watchGUIState(ui.app.State(), last_state)
			last_state = ui.app.State()

			if len(ui.XPlanes) != last_xplane_count {
				last_xplane_count = len(ui.XPlanes)
				ui.xplaneSelect.SetOptions(ui.XPlanes.List())
			}
		}
	}
}

// xplaneLayout returns the X-Plane section layout
func (ui *AppUI) xplaneLayout() fyne.CanvasObject {
	title := widget.NewLabel("X-Plane")
	title.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewVBox(
		title,
		container.New(
			layout.NewFormLayout(),
			widget.NewLabel(""), ui.xplaneRefresh,
			widget.NewLabel("X-Plane Instance"), ui.xplaneSelect,
			widget.NewLabel("Position Interval"), ui.refreshFreq,
		),
	)
}

// serialLayout returns the Serial Port section layout
func (ui *AppUI) serialLayout() fyne.CanvasObject {
	title := widget.NewLabel("Serial Port")
	title.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewVBox(
		title,
		container.New(layout.NewFormLayout(),
			widget.NewLabel(""), ui.serialPortRefresh,
			widget.NewLabel("Port"), ui.serialPortsSelect,
			widget.NewLabel("Baud Rate"), ui.baudRate,
		),
	)
}

// GetContent returns the content of the AppUI
func (ui *AppUI) GetContent() fyne.CanvasObject {
	return container.NewVBox(
		ui.xplaneLayout(),
		widget.NewSeparator(),
		ui.serialLayout(),
		widget.NewSeparator(),
		container.NewHBox(ui.runButton, ui.stopButton, ui.status),
	)
}

// getSerial will populate the serialPortsSelect with the serial ports available on the system
func (ui *AppUI) getSerial() {
	ui.Logger.Debug("GetSerial")
	ports := serial.FindPorts()

	ui.Logger.Debug("Serial Ports", "count", len(ports), "ports", ports)
	if len(ports) == 0 {
		return
	}
	ui.serialPortsSelect.SetOptions(ports)
}

// FindXplanes will search for X-Plane beacons and add them to the XPlanes map
// It will stop after the timeout
func (ui *AppUI) findXplanes(timeout time.Duration) {
	ui.xplaneRefresh.Disable()
	defer ui.xplaneRefresh.Enable()
	// spend 5 seconds searching for X-Plane beacons
	after := time.After(timeout)
	ui.XPlanes = make(xplane.XPlanes)
	for {
		select {
		case <-after:
			ui.Logger.Debug("FindXplanes Timeout", "timeout", timeout, "count", len(ui.XPlanes), "xplanes", ui.XPlanes.List())
			return
		default:
			beacon, err := xplane.FindXplane(1 * time.Second)
			if err != nil {
				ui.Logger.Warn("FindXplanes failed", "err", err)
				continue
			}
			ui.XPlanes.Add(beacon)
		}
	}
}

// setSerialPort returns a function that will set the serial port on the app
func (ui *AppUI) setSerialPort(xApp *App) func(string) {
	return xApp.SetSerialPort
}

// setXPlane returns a function that will set the X-Plane on the app
func (ui *AppUI) setXPlane(xApp *App) func(string) {
	return func(xp string) {
		addr := ui.XPlanes.Find(xp)
		xApp.SetXPlane(addr)
	}
}

// run will start the app
func (ui *AppUI) run() {
	ui.Logger.Debug("Run")
	ctx, cancel := context.WithCancel(context.Background())
	ui.cancelCtx = cancel
	feedback := make(chan string, 3)
	feedback <- "Starting"

	// start the app
	go ui.app.Run(ctx, feedback)

	// watch the feedback channel and update the status label
	go func() {
		t := time.Now()
		for msg := range feedback {
			if msg == "" && time.Since(t) < 20*time.Second {
				continue
			}
			if msg == "XXX" {
				ui.Logger.Warn("Quit on Feedback", "msg", msg)
				// stop the app and drain the feedback channel
				ui.stop()
				for range feedback {
				}
				return
			}
			ui.status.SetText(msg)
			ui.Logger.Info("Feedback", "msg", msg)
			t = time.Now()
		}
	}()
}

// stop will stop the app
func (ui *AppUI) stop() {
	if ui.cancelCtx == nil {
		serial.Logger.Info("Can't Stop, no context")
		return
	}
	ui.cancelCtx()
	ui.cancelCtx = nil
	ui.Logger.Debug("Stop")
}
