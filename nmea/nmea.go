package nmea

import (
	"fmt"
	"math"
)

type formats struct {
	lat string
	lon string
	alt string
	sog string
	hdg string
}

const (
	// DEFAULT Presicions for NMEA messages (number of decimal places)
	DEFAULT_LAT_PRECISION = 4
	DEFAULT_LON_PRECISION = 4
	DEFAULT_ALT_PRECISION = 2
	DEFAULT_SOG_PRECISION = 6
	DEFAULT_HDG_PRECISION = 3

	// ENHANCED Presicions for NMEA messages (number of decimal places)
	ENHANCED_LAT_PRECISION = 7
	ENHANCED_LON_PRECISION = 7
	ENHANCED_ALT_PRECISION = 4
	ENHANCED_SOG_PRECISION = 7
	ENHANCED_HDG_PRECISION = 3
)

var (
	DEFAULTS = formats{
		lat: fmt.Sprintf("%%02d%%0.%df", DEFAULT_LAT_PRECISION),
		lon: fmt.Sprintf("%%03d%%0.%df", DEFAULT_LON_PRECISION),
		alt: fmt.Sprintf("%%0.%df", DEFAULT_ALT_PRECISION),
		sog: fmt.Sprintf("%%0.%df", DEFAULT_SOG_PRECISION),
		hdg: fmt.Sprintf("%%0.%df", DEFAULT_HDG_PRECISION),
	}

	ENHANCED = formats{
		lat: fmt.Sprintf("%%02d%%02.%df", ENHANCED_LAT_PRECISION),
		lon: fmt.Sprintf("%%03d%%02.%df", ENHANCED_LON_PRECISION),
		alt: fmt.Sprintf("%%0.%df", ENHANCED_ALT_PRECISION),
		sog: fmt.Sprintf("%%0.%df", ENHANCED_SOG_PRECISION),
		hdg: fmt.Sprintf("%%0.%df", ENHANCED_HDG_PRECISION),
	}

	// Formats is the formats used for the NMEA messages
	Formats = DEFAULTS
)

// calculateChecksum will calculate the checksum for a NMEA message
func calculateChecksum(s string) byte {
	cs := 0
	for _, c := range s {
		cs ^= int(c)
	}
	return byte(cs)
}

// calculateLL will convert the latitude and longitude for a NMEA message
// the format is "ddmm.mmmm" for latitude and "dddmm.mmmm" for longitude
// v is the value to convert
// ds is the directions ([N, S,] or [E, W]) positive first, negative second
// fm is the format string (e.g. "%02d%02.4f" for latitude or "%03d%02.4f" for longitude)
func calculateLL(v float64, ds []rune, fm string) string {
	// lat and lon are be supplied as floats, and are converted to degrees and minutes
	// lat and lon direction needs to be determined based on positive or negative values
	vA := math.Abs(v)
	vDegrees := math.Floor(vA)
	vMinutes := (vA - vDegrees) * 60
	vD := ds[0]
	if vA < 0 {
		vD = ds[1]
	}
	vS := fmt.Sprintf(fm+",%s", int(vDegrees), vMinutes, string(vD))
	return vS
}

// calculateLat will convert the latitude for a NMEA message
func calculateLat(lat float64) string {
	// lat needs 4 leading digits and 4 decimal places
	return calculateLL(lat, []rune{'N', 'S'}, Formats.lat)
}

// calculateLon will convert the longitude for a NMEA message
func calculateLon(lon float64) string {
	// lon needs 5 leading digits and 4 decimal places
	return calculateLL(lon, []rune{'E', 'W'}, Formats.lon)
}
