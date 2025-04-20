package util

import "testing"

func TestParseFrequency(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
	}{
		{"0.3k", 300},
		{"0.3KhZ", 300},
		{"0.7khz", 700},
		{"100M", 100_000_000},
		{"100Mhz", 100_000_000},
		{"100mhz", 100_000_000},
		{"1.5ghz", 1_500_000_000},
		{"1.5Ghz", 1_500_000_000},
		{"2.5G", 2_500_000_000},
		{"3.6Ghz", 3_600_000_000},
		{"100e6", 100_000_000},
		{"100e3", 100_000},
		{"1.4e9", 1_400_000_000},
		{"1.001ghz", 1_001_000_000},
		{"1.000000001ghz", 1_000_000_001},
		{"invalid", 0}, // Invalid input case
	}

	for _, test := range tests {
		result, err := ParseFrequency(test.input)
		if err != nil && test.expected != 0 {
			t.Errorf("ParseFrequency(%s) failed with error: %v", test.input, err)
		}
		if result != test.expected {
			t.Errorf("ParseFrequency(%s) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

func TestParseRelativeFrequency(t *testing.T) {
	tests := []struct {
		input       string
		expected    int64
		expectError bool
	}{
		{"+1000", 1000, false},
		{"-500", -500, false},
		{"+9.9mhz", 9900000, false},
		{"-1.2Khz", -1200, false},
		{"1000", 0, true},  // missing sign
		{"++100", 0, true}, // malformed
		{"+1e3", 0, true},  // scientific notation
		{"+", 0, true},     // too short
		{"", 0, true},      // empty
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseRelativeFrequency(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("error = %v, wantErr %v", err, tt.expectError)
			}
			if got != tt.expected {
				t.Errorf("got = %d, expected = %d", got, tt.expected)
			}
		})
	}
}

func TestParseTimeDuration(t *testing.T) {
	testCases := []struct {
		input    string
		expected uint64
		wantErr  bool
	}{
		// Microseconds
		{"1u", 1, false},
		{"1.5u", 1, false},
		{"100u", 100, false},
		{"1µ", 1, false},
		{"1.5µ", 1, false},
		{"100µ", 100, false},

		// Milliseconds
		{"1ms", 1000, false},
		{"1.5ms", 1500, false},
		{"9ms", 9000, false},

		// Seconds
		{"1s", 1000000, false},
		{"1.5s", 1500000, false},
		{"9s", 9000000, false},

		// Error cases
		{"1x", 0, true},
		{"1", 0, true},
		{"", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := ParseTimeDuration(tc.input)

			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error for input %s, got nil", tc.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", tc.input, err)
				return
			}

			if result != tc.expected {
				t.Errorf("For input %s: expected %d, got %d",
					tc.input, tc.expected, result)
			}
		})
	}
}
