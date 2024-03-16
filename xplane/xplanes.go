package xplane

import (
	"log/slog"
	"net"
)

// Logger is the logger for the xplane package
// Defaults to slog.Default(), but can be overridden by the user
var Logger = slog.Default()

// Xplanes is a handy way to keep track of the XPlaneBeacon instances
type XPlanes map[string]*XPlaneBeacon

// Add adds a new XPlaneBeacon to the XPlanes map
// The beacon IP is the ip of the XPlane with the beacon multicast port (49707), but the port is the port
// that XPlane is listening on which is what we are interested in (usually 49000)
func (xps XPlanes) Add(xpb *XPlaneBeacon) {
	s := xpb.Addr().String()
	if _, ok := xps[s]; ok {
		return
	}
	xps[s] = xpb
}

// List returns a list of the XPlane Human Readable identifiers
func (xps XPlanes) List() []string {
	var list []string
	for _, xpb := range xps {
		s := xpb.String()
		list = append(list, s)
	}
	return list
}

// Find returns the net.UDPAddr for the given XPlane Human Readable identifier
func (xps XPlanes) Find(id string) *net.UDPAddr {
	for _, xpb := range xps {
		s := xpb.String()
		if s == id {
			return xpb.Addr()
		}
	}
	return nil
}
