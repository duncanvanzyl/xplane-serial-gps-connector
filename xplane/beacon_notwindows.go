//go:build !windows
// +build !windows

package xplane

import (
	"net"
)

// getConnection returns a UDP connection to the given multicast address.
// This is basically a wrapper around net.ListenMulticastUDP, but it's windows version also
// sets the MulticastLoopback option to true which is needed only on Windows.
func getConnection(gaddr *net.UDPAddr) (*net.UDPConn, error) {
	return net.ListenMulticastUDP("udp", nil, gaddr)
}
