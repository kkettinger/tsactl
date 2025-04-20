package util

import (
	"testing"
)

func TestFormatFrequency(t *testing.T) {
	tests := []struct {
		hertz    uint64
		expected string
	}{
		{12, "12 Hz"},
		{13000, "13 kHz"},
		{13001, "13.001 kHz"},
		{13500, "13.5 kHz"},
		{2000000, "2 MHz"},
		{2100000, "2.1 MHz"},
		{2100001, "2.100001 MHz"},
		{4000000000, "4 GHz"},
		{4500000000, "4.5 GHz"},
		{4500000001, "4.500000001 GHz"},
		{0, "0 Hz"},
		{999, "999 Hz"},
		{1000, "1 kHz"},
		{999999, "999.999 kHz"},
		{1000000, "1 MHz"},
		{999999999, "999.999999 MHz"},
		{1000000000, "1 GHz"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatFrequency(tt.hertz)
			if result != tt.expected {
				t.Errorf("FormatFrequency(%d) = %s, want %s", tt.hertz, result, tt.expected)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{1, "1 µs"},
		{999, "999 µs"},
		{1000, "1 ms"},
		{1500, "1.5 ms"},
		{999_999, "999.999 ms"},
		{1_000_000, "1 s"},
		{1_500_000, "1.5 s"},
	}

	for _, tt := range tests {
		result := FormatTimeDuration(tt.input)
		if result != tt.expected {
			t.Errorf("FormatTimeDuration(%d) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}
