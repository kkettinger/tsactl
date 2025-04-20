package util

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/govalues/decimal"
)

func ParseFrequency(freqStr string) (uint64, error) {
	units := map[string]decimal.Decimal{
		"hz":  decimal.MustNew(1, 0),
		"khz": decimal.MustNew(1000, 0),
		"k":   decimal.MustNew(1000, 0),
		"mhz": decimal.MustNew(1000000, 0),
		"m":   decimal.MustNew(1000000, 0),
		"ghz": decimal.MustNew(1000000000, 0),
		"g":   decimal.MustNew(1000000000, 0),
	}

	freqStr = strings.ToLower(freqStr)

	// Match frequency unit first
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)([a-z]+)$`)
	matches := re.FindStringSubmatch(freqStr)
	if matches != nil {
		num, err := decimal.Parse(matches[1])
		if err != nil {
			return 0, fmt.Errorf("invalid frequency format: %s (%w)", freqStr, err)
		}

		unit := matches[2]
		if multiplier, exists := units[unit]; exists {
			result, err := num.Mul(multiplier)
			if err != nil {
				return 0, fmt.Errorf("multiplication error: %w", err)
			}

			whole, _, ok := result.Int64(0)
			if !ok {
				return 0, fmt.Errorf("conversion error: value too large")
			}
			return uint64(whole), nil //nolint:gosec
		} else if unit != "" {
			return 0, fmt.Errorf("invalid unit '%s'", unit)
		}

		whole, _, ok := num.Int64(0)
		if !ok {
			return 0, fmt.Errorf("conversion error: value too large")
		}

		return uint64(whole), nil //nolint:gosec
	}

	// Match scientific notation
	re = regexp.MustCompile(`^(\d+(?:\.\d+)?)(?:e(\d+))?$`)
	matches = re.FindStringSubmatch(freqStr)
	if matches != nil {
		result, err := decimal.Parse(freqStr)
		if err != nil {
			return 0, fmt.Errorf("invalid frequency format: %s (%w)", freqStr, err)
		}

		signed, _, ok := result.Int64(0)
		if !ok {
			return 0, fmt.Errorf("conversion error: value too large")
		}

		return uint64(signed), nil //nolint:gosec
	}

	return 0, fmt.Errorf("invalid frequency format: %s", freqStr)
}

// TODO: refactor
func ParseRelativeFrequency(freqStr string) (int64, error) {
	if len(freqStr) < 2 {
		return 0, fmt.Errorf("invalid relative frequency: %s", freqStr)
	}

	sign := freqStr[0]
	if sign != '+' && sign != '-' {
		return 0, fmt.Errorf("relative frequency must start with '+' or '-'")
	}

	absStr := freqStr[1:]
	if strings.ContainsAny(absStr, "eE") {
		return 0, fmt.Errorf("scientific notation not supported in relative frequency: %s", freqStr)
	}

	absVal, err := ParseFrequency(absStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse frequency: %w", err)
	}

	if absVal > math.MaxInt64 {
		return 0, fmt.Errorf("value too large for int64")
	}

	val := int64(absVal)
	if sign == '-' {
		val = -val
	}

	return val, nil
}

func ParseTimeDuration(timeStr string) (uint64, error) {
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*([musnµ])s?$`)
	matches := re.FindStringSubmatch(strings.ToLower(timeStr))
	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid time format: %s", timeStr)
	}

	num, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid numeric value: %s", matches[1])
	}

	switch strings.ToLower(matches[2]) {
	case "µ": // microseconds (as µ)
		return uint64(num), nil
	case "u": // microseconds
		return uint64(num), nil
	case "m": // milliseconds
		return uint64(num * 1000), nil
	case "s": // seconds
		return uint64(num * 1000000), nil
	}

	return 0, fmt.Errorf("unknown time unit: %s", matches[2])
}
