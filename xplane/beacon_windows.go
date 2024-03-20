//go:build windows
// +build windows

package xplane

import (
	"fmt"
	"net"

	"golang.org/x/net/ipv4"
)

// getConnection returns a UDP connection to the given multicast address.
// This is basically a wrapper around net.ListenMulticastUDP, but it also sets
// the MulticastLoopback option to true which is needed only on Windows.
func getConnection(gaddr *net.UDPAddr) (*net.UDPConn, error) {
	conn, err := net.ListenMulticastUDP("udp", nil, gaddr)
	if err != nil {
		return conn, err
	}

	pc := ipv4.NewPacketConn(conn)

	// added for windows
	if loop, err := pc.MulticastLoopback(); err == nil {
		fmt.Printf("MulticastLoopback status:%v\n", loop)
		if !loop {
			if err := pc.SetMulticastLoopback(true); err != nil {
				fmt.Printf("SetMulticastLoopback error:%v\n", err)
			}
		}
	}

	return conn, nil
}
