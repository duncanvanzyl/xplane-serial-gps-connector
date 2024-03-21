package nmea

import (
	"fmt"
	"math"
)

// ToGPVTG will convert a heading, speedEast and speedSouth to a NMEA GPVTG message
func ToGPVTG(heading float64, sog float64) string {
	// Example GPVTG message:
	// $GPVTG,224.592,T,224.592,M,0.003,N,0.005,K,D*20
	// 224.592,T      True course made good over ground, in degrees
	// 224.592,M      Magnetic course made good over ground, in degrees
	// 0.003,N        Speed over ground in knots
	// 0.005,K        Speed over ground in kilometers/hour
	// D              Mode indicator: D=Diff, A=Autonomous, E=Estimated, N=Data not valid

	// heading is sometimes negative
	// limit heading to 3 decimal places
	// first heading is true (T), second is magnetic (M)
	headingS := fmt.Sprintf(Formats.hdg, math.Mod(heading+360, 360))

	// knots (N) = 1.94384 * m/s
	sogKnots := fmt.Sprintf(Formats.sog+",N", sog*1.943845249221964)
	// km/h (K) = 3.6 * m/s
	sogKmh := fmt.Sprintf(Formats.sog+",K", sog*3.6)

	// D is for Differential. A=Autonomous, D=Differential, E=Estimated, M=Manual input, N=Data not valid
	mode := "D"

	// limit heading to 3 decimal places
	// true heading (T) and magnetic heading (M) are the same
	bs := fmt.Sprintf("GPVTG,%s,T,%s,M,%s,%s,%s", headingS, headingS, sogKnots, sogKmh, mode)

	return fmt.Sprintf("$%s*%02X\r\n", bs, calculateChecksum(bs))
}
