package util

import (
	"github.com/govalues/decimal"
)

func FormatFrequency(freqHz uint64) string {
	switch {
	case freqHz >= 1000000000:
		return formatDecimal(freqHz, 1000000000) + " GHz"
	case freqHz >= 1000000:
		return formatDecimal(freqHz, 1000000) + " MHz"
	case freqHz >= 1000:
		return formatDecimal(freqHz, 1000) + " kHz"
	default:
		return formatDecimal(freqHz, 1) + " Hz"
	}
}

func FormatTimeDuration(us uint64) string {
	switch {
	case us >= 1_000_000:
		return formatDecimal(us, 1_000_000) + " s"
	case us >= 1_000:
		return formatDecimal(us, 1_000) + " ms"
	default:
		return formatDecimal(us, 1) + " µs"
	}
}

func formatDecimal(value uint64, divisor uint64) string {
	dec, _ := decimal.New(int64(value), 0)   //nolint:gosec
	div, _ := decimal.New(int64(divisor), 0) //nolint:gosec
	result, _ := dec.Quo(div)
	return result.String()
}
