package main

import (
	"context"
	"log/slog"
	"net"
	"sync"

	"github.com/duncanvanzyl/xplane-serial-gps-connector/serial"
	"github.com/duncanvanzyl/xplane-serial-gps-connector/xplane"
)

// AppState the state of the app
type AppState uint8

// Possible app states
const (
	// Incomplete is the state when the app is missing a X-Plane or a Serial port
	Incomplete AppState = iota
	// Running is the state when the app is running
	Running
	// Runable is the state when the app is ready to run
	Runable
)

// App is the main application
type App struct {
	mu           sync.RWMutex
	XPlane       *net.UDPAddr
	Serial       serial.Sender
	PositionFreq uint
	Running      bool
	Logger       *slog.Logger
}

// State returns the current state of the app
// The app can not run without a X-Plane and a Serial port
func (a *App) State() AppState {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.XPlane != nil && a.Serial.Configured() {
		if a.Running {
			return Running
		}
		return Runable
	}
	return Incomplete
}

// SetXPlane sets the X-Plane address
func (a *App) SetXPlane(addr *net.UDPAddr) {
	a.Logger.Debug("Set XPlane", "addr", addr)
	a.mu.Lock()
	defer a.mu.Unlock()
	a.XPlane = addr
}

// SetSerialPort sets the serial port
func (a *App) SetSerialPort(port string) {
	a.Logger.Debug("Set SerialPort", "port", port)
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Serial.SetPort(port)
}

// SetBaudRate sets the baud rate
func (a *App) SetBaudRate(baud int) {
	a.Logger.Debug("Set BaudRate", "baud", baud)
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Serial.SetBaud(baud)
}

// SetPositionFreq sets the position frequency
func (a *App) SetPositionFreq(freq uint) {
	a.Logger.Debug("Set PositionFreq", "freq", freq)
	a.mu.Lock()
	defer a.mu.Unlock()
	a.PositionFreq = freq
}

// Run will start the app
// It will request positions from X-Plane and send them to the serial port.
// It will stop when the context is canceled.
func (a *App) Run(ctx context.Context, feedback chan<- string) {
	var wg sync.WaitGroup
	a.Running = true
	defer func() {
		a.Logger.Debug("Stopping")
		a.Running = false
		close(feedback)
	}()
	c := make(chan xplane.Position)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(c)
		xplane.RequestPositions(ctx, a.XPlane, a.PositionFreq, c, feedback)
		a.Logger.Debug("RequestPositions Done")
	}()
	wg.Add(1)
	go func() {
		err := a.Serial.SendPositions(c, feedback)
		if err != nil {
			a.Logger.Info("SendPositions failed", "err", err)
			feedback <- "XXX"
		}
		for range c {
		}
		a.Logger.Debug("SendPositions Done, channel drained")
		wg.Done()
	}()

	wg.Wait()
	a.Logger.Debug("Run Done")
}
