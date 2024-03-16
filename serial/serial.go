package serial

import (
	"log/slog"

	"go.bug.st/serial"

	"github.com/duncanvanzyl/xplane-serial-gps-connector/outputters"
	"github.com/duncanvanzyl/xplane-serial-gps-connector/xplane"
)

var _ Sender = &Serial{}

// Logger is the default logger for the serial package
var Logger = slog.Default()

// Sender is the interface for sending positions to a serial port
type Sender interface {
	// SendPositions will send the positions from the channel wherever the Sender sends them, and send feedback to the feedback channel
	SendPositions(c <-chan xplane.Position, feedback chan<- string) error
	// Configured will return true if the serial port is configured
	Configured() bool
	// SetPort will set the serial port
	SetPort(string)
	// SetBaud will set the baud rate
	SetBaud(int)
}

// Serial is an object that will send positions to a serial port
type Serial struct {
	port       string
	mode       *serial.Mode
	Outputters []outputters.Outputter
}

// NewSerial returns a new Serial
func NewSerial(outputters []outputters.Outputter) *Serial {
	return &Serial{
		mode: &serial.Mode{
			BaudRate: 9600,
			Parity:   serial.NoParity,
			DataBits: 8,
		},
		Outputters: outputters,
	}
}

// SendPositions will send the positions from the channel to the serial port
func (s *Serial) SendPositions(c <-chan xplane.Position, feedback chan<- string) error {
	Logger.Debug("SendPositions Started")

	ser, err := serial.Open(
		s.port,
		s.mode,
	)
	if err != nil {
		Logger.Error("Failed to open serial port", "err", err)
		feedback <- "Failed to open serial port"
		return err
	}
	defer func() {
		ser.Close()
		Logger.Debug("Serial port closed")
	}()
	Logger.Debug("Serial port opened", "port", s.port, "mode", s.mode)

	for pos := range c {
		for _, o := range s.Outputters {
			msg, err := o.Output(pos)
			if err != nil {
				Logger.Warn("Output failed", "err", err)
				feedback <- "Output failed"
				continue
			}
			ser.Write([]byte(msg))
			Logger.Debug("Sent", "msg", msg)
		}
	}

	return nil
}

// Configured will return true if the serial port is configured
func (s *Serial) Configured() bool {
	return s.port != "" && s.mode.BaudRate != 0
}

// SetPort will set the serial port
func (s *Serial) SetPort(port string) {
	Logger.Debug("SetPort", "port", port)
	s.port = port
}

// SetBaud will set the baud rate
func (s *Serial) SetBaud(baud int) {
	Logger.Debug("SetBaud", "baud", baud)
	s.mode.BaudRate = baud
}

// SetDataBits will set the data bits
func (s *Serial) SetDataBits(dataBits int) {
	Logger.Debug("Set DataBits", "databits", dataBits)
	s.mode.DataBits = dataBits
}

// SetParity will set the parity
func (s *Serial) SetParity(parity serial.Parity) {
	Logger.Debug("Set Parity", "parity", parity)
	s.mode.Parity = parity
}

// SetStopBits will set the stop bits
func (s *Serial) SetStopBits(stopBits serial.StopBits) {
	Logger.Debug("Set StopBits", "stopbits", stopBits)
	s.mode.StopBits = stopBits
}

// Mode will return the current mode
func (s *Serial) Mode() serial.Mode { return *s.mode }

// FindPorts will find the available serial ports on the system
func FindPorts() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		panic(err)
	}

	Logger.Debug("Found ports", "ports", ports)
	return ports
}
