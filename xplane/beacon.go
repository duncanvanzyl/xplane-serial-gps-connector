package xplane

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	MCAST_PORT = 49707
	RETRIES    = 10
)

var ErrInvalidBeacon = errors.New("invalid beacon")

// XPlaneBasicBeacon is the basic beacon information
type XPlaneBasicBeacon struct {
	BeaconMajorVersion uint8  // 1 at the time of X-Plane 10.40, 11.55
	BeaconMinorVersion uint8  // 1 at the time of X-Plane 10.40, 2 for 11.55
	ApplicationHostID  int32  // 1 for X-Plane, 2 for PlaneMaker
	VersionNumber      int32  // 104014 is X-Plane 10.40b14, 115501 is 11.55r2
	Role               uint32 // 1 for master, 2 for extern visual, 3 for IOS
	Port               uint16 // port number X-Plane is listening on
}

// XPlaneBeacon is the full beacon information
type XPlaneBeacon struct {
	XPlaneBasicBeacon
	IP           *net.UDPAddr // the IP address of the computer
	ComputerName string       // the hostname of the computer
	RaknetPort   uint16       // port number the X-Plane Raknet clinet is listening on
}

// Addr returns the UDP address of the beacon
// Since the beacon IP contains the computer IP but the beacon port, this is a new UDP address with the
// X-Plane port
func (xpb *XPlaneBeacon) Addr() *net.UDPAddr {
	return &net.UDPAddr{
		IP:   xpb.IP.IP,
		Port: int(xpb.Port),
	}
}

// RaknetAddr returns the UDP address of the Raknet client
// Since the beacon IP contains the computer IP but the beacon port, this is a new UDP address with the
// Raknet port
func (xpb *XPlaneBeacon) RaknetAddr() *net.UDPAddr {
	return &net.UDPAddr{
		IP:   xpb.IP.IP,
		Port: int(xpb.RaknetPort),
	}
}

// BeaconVersion returns the beacon version as a string
func (xpb *XPlaneBeacon) BeaconVersion() string {
	return fmt.Sprintf("%d.%d", xpb.BeaconMajorVersion, xpb.BeaconMinorVersion)
}

// ApplicationType returns the application type as a string
func (xpb *XPlaneBeacon) ApplicationType() string {
	switch xpb.ApplicationHostID {
	case 1:
		return "X-Plane"
	case 2:
		return "PlaneMaker"
	default:
		return "Unknown"
	}
}

// RoleType returns the role type as a string
func (xpb *XPlaneBeacon) RoleType() string {
	switch xpb.Role {
	case 1:
		return "Master"
	case 2:
		return "Extern visual"
	case 3:
		return "IOS"
	default:
		return "Unknown"
	}
}

// String returns a string representation of the beacon
// This is the normal, human readable identifier for the beacon
func (xpb *XPlaneBeacon) String() string {
	return fmt.Sprintf("%s %s on %s", xpb.ApplicationType(), xpb.RoleType(), xpb.ComputerName)
}

// Details returns a detailed string representation of the beacon
func (xpb *XPlaneBeacon) Details() string {
	return fmt.Sprintf("X-Plane Beacon Version: %s %s(%d %s) on %s UDP Port: %s Raknet: %s",
		xpb.BeaconVersion(),
		xpb.ApplicationType(), xpb.VersionNumber, xpb.RoleType(),
		xpb.ComputerName,
		xpb.Addr(),
		xpb.RaknetAddr(),
	)
}

// IsMaster returns true if the beacon is from the master X-Plane instance
func (xpb *XPlaneBeacon) IsMaster() bool {
	return xpb.ApplicationHostID == 1 && xpb.Role == 1
}

// decodeBeacon will decode the beacon from the byte slice
func decodeBeacon(bs []byte, addr *net.UDPAddr) (*XPlaneBeacon, error) {
	sb := XPlaneBasicBeacon{}
	buf := bytes.NewBuffer(bs)
	err := binary.Read(buf, binary.LittleEndian, &sb)
	if err != nil {
		return nil, fmt.Errorf("could not read beacon: %v", err)
	}

	// read the computer name from the buffer until the first null byte
	var computerName = bs[binary.Size(sb) : len(bs)-2]
	// strip everything after the first null byte
	computerName = bytes.Trim(computerName, "\x00")

	// the port is the last 2 bytes
	raknetPort := binary.LittleEndian.Uint16(bs[len(bs)-2:])

	return &XPlaneBeacon{
		XPlaneBasicBeacon: sb,
		IP:                addr,
		ComputerName:      string(computerName),
		RaknetPort:        raknetPort,
	}, nil
}

// listenForBeacon will listen for a beacon on the multicast group and return the beacon or an error if no
// beacon is found
func listenForBeacon(timeout time.Time, gaddr *net.UDPAddr) (*XPlaneBeacon, error) {
	conn, err := getConnection(gaddr)
	if err != nil {
		return nil, fmt.Errorf("could not listen to UDP address: %v", err)
	}
	defer conn.Close()

	// join the multicast group
	if err := conn.SetReadBuffer(1024 * 1024); err != nil {
		return nil, fmt.Errorf("could not set read buffer: %v", err)
	}
	if err := conn.SetReadDeadline(timeout); err != nil {
		return nil, fmt.Errorf("could not set read deadline: %v", err)
	}

	Logger.Debug("Listening for X-Plane beacons...", "local", conn.LocalAddr(), "remote", conn.RemoteAddr())

	// read the beacon
	buf := make([]byte, 512)
	n, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		return nil, fmt.Errorf("could not read from UDP address: %v", err)
	}

	// beacons should start with "BECN"
	if bytes.Equal(buf[:4], []byte("BECN\x00")) {
		Logger.Warn("Unknown Beacon", "addr", addr, "beacon", buf[:n])
		return nil, ErrInvalidBeacon
	}

	// decode the beacon
	return decodeBeacon(buf[5:n], addr)
}

// FindXplane will listen for X-Plane beacons on the multicast group and return the first master X-Plane
// instance found or an error if no master X-Plane instance is found
// wait is the time to wait for a beacon
func FindXplane(wait time.Duration) (*XPlaneBeacon, error) {
	// listen to the multicast group on the MCAST_PORT

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("239.255.1.1:%d", MCAST_PORT))
	if err != nil {
		return nil, fmt.Errorf("could not resolve UDP address: %v", err)
	}

	for i := 0; i < RETRIES; i++ {
		beacon, err := listenForBeacon(time.Now().Add(wait), addr)
		if err == ErrInvalidBeacon {
			Logger.Warn("Invalid Beacon", "err", err)
			continue
		}
		if err != nil {
			return nil, err
		}

		if err == nil && !beacon.IsMaster() {
			Logger.Info("Found non-Master X-Plane", "beacon", beacon)
			continue
		}

		return beacon, nil

	}

	return nil, nil
}
