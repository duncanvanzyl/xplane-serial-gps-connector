package nmea

import (
	"testing"
	"time"
)

func TestGenerateGGA(t *testing.T) {
	testCases := []struct {
		name       string
		timestamp  time.Time
		lat        float64
		lon        float64
		quality    uint
		satellites uint
		hdop       float64
		alt        float64
		sep        float64
		format     formats
		expected   string
	}{
		{
			name:       "Test 0-All Zeros",
			timestamp:  time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
			lat:        0.0,
			lon:        0.0,
			quality:    0,
			satellites: 0,
			hdop:       0.0,
			alt:        0.0,
			sep:        0.0,
			format:     DEFAULTS,
			expected:   "$GPGGA,000000.000,0000.0000,N,00000.0000,E,0,0,0.0,0.00,M,0.00,M,,*5D\r\n",
		},
		{
			name:       "Test 1",
			timestamp:  time.Date(2022, time.January, 1, 12, 34, 56, 789000000, time.UTC),
			lat:        12.3456,
			lon:        98.7654,
			quality:    1,
			satellites: 10,
			hdop:       1.2,
			alt:        100.5,
			sep:        50.0,
			format:     DEFAULTS,
			expected:   "$GPGGA,123456.789,1220.7360,N,09845.9240,E,1,10,1.2,100.50,M,50.00,M,,*52\r\n",
		},
		{
			name:       "Test 2-Very Precise",
			timestamp:  time.Date(2022, time.January, 1, 12, 34, 56, 789000000, time.UTC),
			lat:        12.345678901234567,
			lon:        98.765432109876543,
			quality:    1,
			satellites: 10,
			hdop:       1.2,
			alt:        100.9876543210,
			sep:        50.0,
			format:     DEFAULTS,
			expected:   "$GPGGA,123456.789,1220.7407,N,09845.9259,E,1,10,1.2,100.99,M,50.00,M,,*59\r\n",
		},
		{
			name:       "Test 3-South and West",
			timestamp:  time.Date(2022, time.January, 1, 12, 34, 56, 789000000, time.UTC),
			lat:        -12.3456,
			lon:        -98.7654,
			quality:    1,
			satellites: 10,
			hdop:       1.2,
			alt:        100.5,
			sep:        50.0,
			format:     DEFAULTS,
			expected:   "$GPGGA,123456.789,1220.7360,S,09845.9240,W,1,10,1.2,100.50,M,50.00,M,,*5D\r\n",
		},
		{
			name:       "Test 4-South and East",
			timestamp:  time.Date(2022, time.January, 1, 12, 34, 56, 789000000, time.UTC),
			lat:        -12.3456,
			lon:        98.7654,
			quality:    1,
			satellites: 10,
			hdop:       1.2,
			alt:        100.5,
			sep:        50.0,
			format:     DEFAULTS,
			expected:   "$GPGGA,123456.789,1220.7360,S,09845.9240,E,1,10,1.2,100.50,M,50.00,M,,*4F\r\n",
		},
		{
			name:       "Test 5-North and West",
			timestamp:  time.Date(2022, time.January, 1, 12, 34, 56, 789000000, time.UTC),
			lat:        12.3456,
			lon:        -98.7654,
			quality:    1,
			satellites: 10,
			hdop:       1.2,
			alt:        100.5,
			sep:        50.0,
			format:     DEFAULTS,
			expected:   "$GPGGA,123456.789,1220.7360,N,09845.9240,W,1,10,1.2,100.50,M,50.00,M,,*40\r\n",
		},

		{
			name:       "Test 6-All Zeros-Enhanced Format",
			timestamp:  time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
			lat:        0.0,
			lon:        0.0,
			quality:    0,
			satellites: 0,
			hdop:       0.0,
			alt:        0.0,
			sep:        0.0,
			format:     ENHANCED,
			expected:   "$GPGGA,000000.000,0000.0000000,N,00000.0000000,E,0,0,0.0,0.0000,M,0.0000,M,,*5D\r\n",
		},
		{
			name:       "Test 7-Enhanced Format",
			timestamp:  time.Date(2022, time.January, 1, 12, 34, 56, 789000000, time.UTC),
			lat:        12.3456,
			lon:        98.7654,
			quality:    1,
			satellites: 10,
			hdop:       1.2,
			alt:        100.5,
			sep:        50.0,
			format:     ENHANCED,
			expected:   "$GPGGA,123456.789,1220.7360000,N,09845.9240000,E,1,10,1.2,100.5000,M,50.0000,M,,*52\r\n",
		},
		{
			name:       "Test 8-Very Precise-Enhanced Format",
			timestamp:  time.Date(2022, time.January, 1, 12, 34, 56, 789000000, time.UTC),
			lat:        12.345678901234567,
			lon:        98.765432109876543,
			quality:    1,
			satellites: 10,
			hdop:       1.2,
			alt:        100.9876543210,
			sep:        50.0,
			format:     ENHANCED,
			expected:   "$GPGGA,123456.789,1220.7407341,N,09845.9259266,E,1,10,1.2,100.9877,M,50.0000,M,,*5C\r\n",
		},
		{
			name:       "Test 9-South and West-Enhanced Format",
			timestamp:  time.Date(2022, time.January, 1, 12, 34, 56, 789000000, time.UTC),
			lat:        -12.3456,
			lon:        -98.7654,
			quality:    1,
			satellites: 10,
			hdop:       1.2,
			alt:        100.5,
			sep:        50.0,
			format:     ENHANCED,
			expected:   "$GPGGA,123456.789,1220.7360000,S,09845.9240000,W,1,10,1.2,100.5000,M,50.0000,M,,*5D\r\n",
		},
		{
			name:       "Test 10-South and East-Enhanced Format",
			timestamp:  time.Date(2022, time.January, 1, 12, 34, 56, 789000000, time.UTC),
			lat:        -12.3456,
			lon:        98.7654,
			quality:    1,
			satellites: 10,
			hdop:       1.2,
			alt:        100.5,
			sep:        50.0,
			format:     ENHANCED,
			expected:   "$GPGGA,123456.789,1220.7360000,S,09845.9240000,E,1,10,1.2,100.5000,M,50.0000,M,,*4F\r\n",
		},
		{
			name:       "Test 11-North and West-Enhanced Format",
			timestamp:  time.Date(2022, time.January, 1, 12, 34, 56, 789000000, time.UTC),
			lat:        12.3456,
			lon:        -98.7654,
			quality:    1,
			satellites: 10,
			hdop:       1.2,
			alt:        100.5,
			sep:        50.0,
			format:     ENHANCED,
			expected:   "$GPGGA,123456.789,1220.7360000,N,09845.9240000,W,1,10,1.2,100.5000,M,50.0000,M,,*40\r\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			Formats = tc.format
			result := generateGGA(tc.timestamp, tc.lat, tc.lon, tc.quality, tc.satellites, tc.hdop, tc.alt, tc.sep)
			if result != tc.expected {
				t.Errorf("Expected: %s, but got: %s", tc.expected, result)
			}
		})
	}
}
