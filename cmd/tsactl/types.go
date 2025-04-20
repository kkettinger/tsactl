package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/kkettinger/go-tinysa"
	"github.com/kkettinger/tsactl/internal/util"
	"strconv"
	"strings"
)

type Frequency struct {
	Value uint64
	Valid bool
}

func (f *Frequency) Decode(ctx *kong.DecodeContext) error {
	var val string
	if err := ctx.Scan.PopValueInto(ctx.Value.Name, &val); err != nil {
		return err
	}

	if freq, err := util.ParseFrequency(val); err != nil {
		return fmt.Errorf("failed to parse frequency: %w", err)
	} else {
		f.Value = freq
		f.Valid = true
	}

	return nil
}

type FrequencyRel struct {
	Value    int64
	Relative bool
	Valid    bool
}

func (f *FrequencyRel) Decode(ctx *kong.DecodeContext) error {
	var val string
	if err := ctx.Scan.PopValueInto(ctx.Value.Name, &val); err != nil {
		return err
	}

	// Check if frequency is relative
	freqRel, err := util.ParseRelativeFrequency(val)
	if err == nil {
		f.Relative = true
		f.Valid = true
		f.Value = freqRel
		return nil
	}

	// Otherwise, treat as normal frequency
	if freq, err := util.ParseFrequency(val); err != nil {
		return fmt.Errorf("failed to parse frequency: %w", err)
	} else {
		f.Relative = false
		f.Valid = true
		f.Value = int64(freq) // #nosec G115
	}

	return nil
}

type Time struct {
	Value uint64
	Valid bool
}

func (f *Time) Decode(ctx *kong.DecodeContext) error {
	var val string
	if err := ctx.Scan.PopValueInto(ctx.Value.Name, &val); err != nil {
		return err
	}

	if time, err := util.ParseTimeDuration(val); err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	} else {
		f.Value = time
		f.Valid = true
	}

	return nil
}

type MarkerDelta struct {
	RefMarker uint
	Off       bool
	Valid     bool
}

func (d *MarkerDelta) Decode(ctx *kong.DecodeContext) error {
	var val string
	if err := ctx.Scan.PopValueInto(ctx.Value.Name, &val); err != nil {
		return err
	}

	val = strings.ToLower(val)
	switch val {
	case "off":
		d.Off, d.Valid = true, true
	default:
		v, err := strconv.ParseUint(val, 10, 0)
		if err != nil {
			return fmt.Errorf("failed to parse marker delta value: %w", err)
		}
		d.RefMarker = uint(v)
		d.Valid = true
	}

	return nil
}

type TraceCalc struct {
	Valid bool
	Off   bool
	Mode  tinysa.TraceCalc
}

func (o *TraceCalc) ValidOpts() []string {
	return append([]string{"off"}, tinysa.TraceCalcOptions()...)
}

func (o *TraceCalc) Decode(ctx *kong.DecodeContext) error {
	var val string
	if err := ctx.Scan.PopValueInto(ctx.Value.Name, &val); err != nil {
		return err
	}
	val = strings.ToLower(val)

	switch val {
	case "off":
		o.Valid, o.Off = true, true
	default:
		if c, ok := tinysa.TraceCalcFromString(val); ok {
			o.Valid, o.Mode = true, c
		} else {
			validOpts := strings.Join(o.ValidOpts(), ", ")
			return fmt.Errorf("invalid option '%s', must be one of: %s", val, validOpts)
		}
	}

	return nil
}

type TraceUnit struct {
	Valid bool
	Unit  tinysa.TraceUnit
}

func (o *TraceUnit) ValidOpts() []string {
	return tinysa.TraceUnitOptions()
}

func (o *TraceUnit) Decode(ctx *kong.DecodeContext) error {
	var val string
	if err := ctx.Scan.PopValueInto(ctx.Value.Name, &val); err != nil {
		return err
	}

	if c, ok := tinysa.TraceUnitFromString(val); ok {
		o.Valid, o.Unit = true, c
	} else {
		validOpts := strings.Join(o.ValidOpts(), ", ")
		return fmt.Errorf("invalid option '%s', must be one of: %s", val, validOpts)
	}

	return nil
}

type SweepMode struct {
	Valid bool
	Mode  tinysa.SweepMode
}

func (o *SweepMode) ValidOpts() []string {
	return tinysa.SweepModeOptions()
}

func (o *SweepMode) Decode(ctx *kong.DecodeContext) error {
	var val string
	if err := ctx.Scan.PopValueInto(ctx.Value.Name, &val); err != nil {
		return err
	}

	if c, ok := tinysa.SweepModeFromString(val); ok {
		o.Valid, o.Mode = true, c
	} else {
		validOpts := strings.Join(o.ValidOpts(), ", ")
		return fmt.Errorf("invalid option '%s', must be one of: %s", val, validOpts)
	}

	return nil
}

type Spur struct {
	Valid  bool
	Enable bool
	Auto   bool
}

func (o *Spur) ValidOpts() []string {
	return []string{"on", "off", "auto"}
}

func (o *Spur) Decode(ctx *kong.DecodeContext) error {
	var val string
	if err := ctx.Scan.PopValueInto(ctx.Value.Name, &val); err != nil {
		return err
	}

	val = strings.ToLower(val)
	switch val {
	case "on":
		o.Valid = true
		o.Enable = true
	case "off":
		o.Valid = true
		o.Enable = false
	case "auto":
		o.Valid = true
		o.Auto = true
	default:
		validOpts := strings.Join(o.ValidOpts(), ", ")
		return fmt.Errorf("invalid option '%s', must be one of: %s", val, validOpts)
	}

	return nil
}
