package serial

import (
	"github.com/duncanvanzyl/xplane-serial-gps-connector/outputters"
	"github.com/duncanvanzyl/xplane-serial-gps-connector/xplane"
)

// Dummy is a dummy serial port
// This just logs the output to the logger
type Dummy struct {
	Outputters []outputters.Outputter
}

// SendPositions will send the positions from the channel to the serial port
func (s *Dummy) SendPositions(c <-chan xplane.Position) error {
	for pos := range c {
		Logger.Info("Position", "pos", pos)
		for _, o := range s.Outputters {
			msg, err := o.Output(pos)
			if err != nil {
				Logger.Warn("Output failed", "err", err)
				continue
			}
			Logger.Info("Output", "msg", msg)
		}
	}
	return nil
}

// Configured will return true if the serial port is configured
func (s *Dummy) Configured() bool { return true }

// SetPort will set the serial port
func (s *Dummy) SetPort(port string) {}

// SetBaud will set the baud rate
func (s *Dummy) SetBaud(baud int) {}
