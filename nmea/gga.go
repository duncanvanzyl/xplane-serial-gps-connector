package nmea

import (
	"fmt"
	"time"
)

// ToGPGGA will convert a latitude, longitude and altitude to a NMEA GPGGA message
func ToGPGGA(lat float64, lon float64, alt float64) string {
	// Example GPGGA message:
	// $GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*47
	// 123519       Fix taken at 12:35:19 UTC
	// 4807.038,N   Latitude 48 deg 07.038' N
	// 01131.000,E  Longitude 11 deg 31.000' E
	// 1            Fix quality: 0 = invalid
	// 08           Number of satellites being tracked
	// 0.9          Horizontal dilution of position
	// 545.4,M      Altitude, Meters, above mean sea level
	// 46.9,M       Height of geoid (mean sea level) above WGS84
	// (empty field) time in seconds since last DGPS update
	// (empty field) DGPS station ID number
	// *47          the checksum data, always begins with *

	// time is not supplied, so we will use the current time
	t := time.Now().UTC()
	ts := t.Format("150405.000")

	laS := calculateLat(lat)
	loS := calculateLon(lon)

	// altUnit set to "M" for meters
	altS := fmt.Sprintf(Formats.alt+",M", alt)

	// quality set to 8 for a simulated fix (see https://docs.novatel.com/OEM7/Content/Logs/GPGGA.htm#GPSQualityIndicators)
	quality := 8
	// numSV is the number of satellites in view and is set to 12 for a simulated fix
	numSV := 12
	// HDOP is the horizontal dilution of presicion. lower values are better. normal range is 1-2, but set to 0.5 for a simulated fix
	HDOP := 0.5

	// sep is the height of the geoid above the WGS84 ellipsoid. 0.0 is a simulated fix. We could consider calculating this later if required
	// sepUnit set to "M" for meters
	sep := 0.0
	sepS := fmt.Sprintf(Formats.alt+",M", sep)

	// diffAge and diffStation are not set, so they will default to empty strings
	diffAge := ""
	diffStation := ""

	bs := fmt.Sprintf("GPGGA,%s,%s,%s,%d,%d,%0.1f,%s,%s,%s,%s", ts, laS, loS, quality, numSV, HDOP, altS, sepS, diffAge, diffStation)

	return fmt.Sprintf("$%s*%02X\r\n", bs, calculateChecksum(bs))
}
