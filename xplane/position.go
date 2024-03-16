package xplane

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"time"
)

type Position struct {
	Dat_lon     float64 // float longitude of the aircraft in X-Plane of course, in degrees
	Dat_lat     float64 // float latitude
	Dat_ele     float64 // float elevation above sea level in meters
	Y_agl_mtr   float32 // float elevation above the terrain in meters
	Veh_the_loc float32 // float pitch, degrees
	Veh_psi_loc float32 // float true heading, in degrees
	Veh_phi_loc float32 // float roll, in degrees
	Vx_wrl      float32 // float speed in the x, EAST, direction, in meters per second
	Vy_wrl      float32 // float speed in the y, UP, direction, in meters per second
	Vz_wrl      float32 // float speed in the z, SOUTH, direction, in meters per second
	Prad        float32 // float roll rate in radians per second
	Qrad        float32 // float pitch rate in radians per second
	Rrad        float32 // float yaw rate in radians per second
}

// SOG returns the speed over ground in m/s
func (p *Position) SOG() float64 {
	// speed over ground in m/s
	return math.Sqrt(float64(p.Vx_wrl*p.Vx_wrl) + float64(p.Vz_wrl*p.Vz_wrl))
}

// ReadPosition reads a Position from an io.Reader
func ReadPosition(r io.Reader) (*Position, error) {
	pos := &Position{}
	err := binary.Read(r, binary.LittleEndian, pos)
	if err != nil {
		Logger.Warn("binary.Read failed", "err", err, "size", binary.Size(pos))
		return nil, err
	}
	return pos, nil
}

// getRequest will return a byte slice with the request for positions
// freq is the frequency in Hz. Valid values are numbers up to 60
func getRequest(freq uint) []byte {
	return []byte(fmt.Sprintf("RPOS\x00%d\x00", freq))
}

// RequestPositions will request positions from X-Plane and send them to the channel
// ctx is the context to stop requesting positions
// xp_addr is the address of the X-Plane instance
// c is the channel to send the positions to
// wg is the wait group to signal when the function is done
func RequestPositions(ctx context.Context, xp_addr *net.UDPAddr, freq uint, c chan<- Position, feedback chan<- string) {
	// create a udp connection
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		panic(err)
	}
	defer func() {
		// stop requesting positions and close the connection
		conn.WriteToUDP(getRequest(0), xp_addr)
		conn.Close()
	}()

	// request positions
	_, err = conn.WriteToUDP(getRequest(freq), xp_addr)
	if err != nil {
		Logger.Error("Failed to request positions", "err", err)
		feedback <- "Failed to request positions"
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn.SetDeadline(time.Now().Add(1 * time.Second))
			buf := make([]byte, 1500)

			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				if err, ok := err.(net.Error); ok && err.Timeout() {
					Logger.Info("Timeout")
					feedback <- "Timeout"
					continue
				}
				Logger.Error("Failed to read from UDP", "err", err)
				feedback <- "Failed to read from UDP"
				return
			}
			buffer := bytes.NewBuffer(buf[:n])
			if string(buffer.Next(5)) != "RPOS4" {
				Logger.Warn("Invalid header", "header", buffer.String())
				feedback <- "Invalid header"
				continue
			}

			pos, err := ReadPosition(buffer)
			if err != nil {
				Logger.Warn("ReadPosition failed", "err", err)
				feedback <- "ReadPosition failed"
				continue
			}

			feedback <- ""
			c <- *pos
		}
	}
}
