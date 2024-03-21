package nmea

import (
	"testing"
)

func TestToGPVTG(t *testing.T) {
	testCases := []struct {
		name     string
		heading  float64
		sog      float64
		format   formats
		expected string
	}{
		{"Zeros", 0, 0, DEFAULTS, "$GPVTG,0.000,T,0.000,M,0.000000,N,0.000000,K,D*26\r\n"},
		{"Normal", 45.123, 1, DEFAULTS, "$GPVTG,45.123,T,45.123,M,1.943845,N,3.600000,K,D*25\r\n"},
		{"Negative Heading", -45.123, 10, DEFAULTS, "$GPVTG,314.877,T,314.877,M,19.438452,N,36.000000,K,D*27\r\n"},
		{">1 Heading Rotation", 360.123, 1, DEFAULTS, "$GPVTG,0.123,T,0.123,M,1.943845,N,3.600000,K,D*25\r\n"},
		{"2 Heading Rotations", 720.123, 1, DEFAULTS, "$GPVTG,0.123,T,0.123,M,1.943845,N,3.600000,K,D*25\r\n"},
		{"Very Precise", 45.1234567890123456789, 12.34567890123456789, DEFAULTS, "$GPVTG,45.123,T,45.123,M,23.998089,N,44.444444,K,D*2E\r\n"},

		{"Zeros-Enhanced", 0, 0, ENHANCED, "$GPVTG,0.000,T,0.000,M,0.0000000,N,0.0000000,K,D*26\r\n"},
		{"Normal-Enhanced", 45.123, 1, ENHANCED, "$GPVTG,45.123,T,45.123,M,1.9438452,N,3.6000000,K,D*27\r\n"},
		{"Negative Heading-Enhanced", -45.123, 10, ENHANCED, "$GPVTG,314.877,T,314.877,M,19.4384525,N,36.0000000,K,D*22\r\n"},
		{">1 Heading Rotation-Enhanced", 360.123, 1, ENHANCED, "$GPVTG,0.123,T,0.123,M,1.9438452,N,3.6000000,K,D*27\r\n"},
		{"2 Heading Rotations-Enhanced", 720.123, 1, ENHANCED, "$GPVTG,0.123,T,0.123,M,1.9438452,N,3.6000000,K,D*27\r\n"},
		{"Very Precise-Enhanced", 45.1234567890123456789, 12.34567890123456789, ENHANCED, "$GPVTG,45.123,T,45.123,M,23.9980893,N,44.4444440,K,D*2D\r\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			Formats = tc.format
			result := ToGPVTG(tc.heading, tc.sog)
			if result != tc.expected {
				t.Errorf("Expected: %s, but got: %s", tc.expected, result)
			}
		})
	}
}
