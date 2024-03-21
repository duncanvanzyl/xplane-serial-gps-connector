package nmea

import (
	"fmt"
	"testing"
)

func TestCalculateChecksum(t *testing.T) {
	testCases := []struct {
		input    string
		expected byte
	}{
		{"PFEC,GPint,RMC05", 0x2D},
		{"GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W", 106},
		{"GPGLL,4807.038,N,01131.000,E,123519,A", 0x25},
		{"GPGSA,A,3,04,05,,09,12,,,24,,,,,2.5,1.3,2.1", 57},
		{"GPGSV,3,1,11,20,75,064,46,24,63,231,42,28,52,160,41,32,45,047,39", 0x78},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Input: %s", tc.input), func(t *testing.T) {
			result := calculateChecksum(tc.input)
			if result != tc.expected {
				t.Errorf("Expected: %d, but got: %d", tc.expected, result)
			}
		})
	}
}

func TestCalculateLL(t *testing.T) {
	testCases := []struct {
		value    float64
		format   string
		ds       [2]rune
		expected string
	}{
		{0.0, DEFAULTS.lat, LAT_RUNES, "0000.0000,N"},
		{8.123456, DEFAULTS.lat, LAT_RUNES, "0807.4074,N"},
		{89.123456, DEFAULTS.lat, LAT_RUNES, "8907.4074,N"},
		{-8.123456, DEFAULTS.lat, LAT_RUNES, "0807.4074,S"},
		{-89.123456, DEFAULTS.lat, LAT_RUNES, "8907.4074,S"},
		{0.0, DEFAULTS.lon, LON_RUNES, "00000.0000,E"},
		{8.123456, DEFAULTS.lon, LON_RUNES, "00807.4074,E"},
		{89.123456, DEFAULTS.lon, LON_RUNES, "08907.4074,E"},
		{-8.123456, DEFAULTS.lon, LON_RUNES, "00807.4074,W"},
		{-89.123456, DEFAULTS.lon, LON_RUNES, "08907.4074,W"},
		{109.123456, DEFAULTS.lon, LON_RUNES, "10907.4074,E"},
		{-109.123456, DEFAULTS.lon, LON_RUNES, "10907.4074,W"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Value: %f", tc.value), func(t *testing.T) {
			result := calculateLL(tc.value, tc.ds, tc.format)
			if result != tc.expected {
				t.Errorf("Expected: %s, but got: %s", tc.expected, result)
			}
		})
	}
}

func TestCalculateLat(t *testing.T) {
	testCases := []struct {
		lat      float64
		expected string
	}{
		{0.0, "0000.0000,N"},
		{8.123456, "0807.4074,N"},
		{89.123456, "8907.4074,N"},
		{-8.123456, "0807.4074,S"},
		{-89.123456, "8907.4074,S"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Lat: %f", tc.lat), func(t *testing.T) {
			result := calculateLat(tc.lat)
			if result != tc.expected {
				t.Errorf("Expected: %s, but got: %s", tc.expected, result)
			}
		})
	}
}

func TestCalculateLon(t *testing.T) {
	testCases := []struct {
		lon      float64
		expected string
	}{
		{0.0, "00000.0000,E"},
		{8.123456, "00807.4074,E"},
		{89.123456, "08907.4074,E"},
		{-8.123456, "00807.4074,W"},
		{-89.123456, "08907.4074,W"},
		{109.123456, "10907.4074,E"},
		{-109.123456, "10907.4074,W"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Lon: %f", tc.lon), func(t *testing.T) {
			result := calculateLon(tc.lon)
			if result != tc.expected {
				t.Errorf("Expected: %s, but got: %s", tc.expected, result)
			}
		})
	}
}
