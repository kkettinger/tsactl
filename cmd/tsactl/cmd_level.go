package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/kkettinger/go-tinysa"
)

type LevelCmd struct {
	Unit         TraceUnit `help:"Set trace unit (${trace_unit_opts})" short:"u" group:"Level flags:" placeholder:"UNIT"`
	RefLevel     *int      `help:"Set trace reference level in dBm" name:"ref" group:"Level flags:" placeholder:"REFLEVEL"`
	RefLevelAuto bool      `help:"Set trace reference level to auto" name:"ref-auto" group:"Level flags:"`
	Scale        *float64  `help:"Set trace scale" short:"s" group:"Level flags:" placeholder:"SCALE"`
	LNA          *bool     `help:"Enable low noise amplifier (LNA)" negatable:"" group:"Level flags:"`
}

func (c *LevelCmd) Validate() error {
	if c.RefLevel != nil && c.RefLevelAuto {
		return fmt.Errorf("only one of --ref or --ref-auto can be set, not both at the same time")
	}

	return nil
}

func (c *LevelCmd) Run(globals *Globals, ctx *kong.Context) error {
	var ops []func(*tinysa.Device) error

	if c.Unit.Valid {
		ops = append(ops, c.SetUnit)
	}

	if c.RefLevel != nil {
		ops = append(ops, c.SetRefLevel)
	}

	if c.RefLevelAuto {
		ops = append(ops, c.SetRefLevelAuto)
	}

	if c.Scale != nil {
		ops = append(ops, c.SetScale)
	}

	if c.LNA != nil {
		if *c.LNA {
			ops = append(ops, c.EnableLNA)
		} else {
			ops = append(ops, c.DisableLNA)
		}
	}

	if len(ops) == 0 {
		_ = ctx.PrintUsage(false)
		return nil
	}

	d, err := initDevice(globals)
	if err != nil {
		return err
	}

	for _, op := range ops {
		if err := op(d); err != nil {
			return err
		}
	}

	return nil
}

func (c *LevelCmd) SetUnit(d *tinysa.Device) error {
	fmt.Println("set trace unit to", c.Unit.Unit.String())
	if err := d.SetTraceUnit(c.Unit.Unit); err != nil {
		return fmt.Errorf("failed to set trace unit to %s: %w", c.Unit.Unit.String(), err)
	}
	return nil
}

func (c *LevelCmd) SetRefLevel(d *tinysa.Device) error {
	fmt.Println("set reference level to", *c.RefLevel)
	if err := d.SetTraceRefLevel(*c.RefLevel); err != nil {
		return fmt.Errorf("failed to set reference level to %d: %w", *c.RefLevel, err)
	}
	return nil
}

func (c *LevelCmd) SetRefLevelAuto(d *tinysa.Device) error {
	fmt.Println("set reference level to auto")
	if err := d.SetTraceRefLevelAuto(); err != nil {
		return fmt.Errorf("failed to set reference level to auto: %w", err)
	}
	return nil
}

func (c *LevelCmd) SetScale(d *tinysa.Device) error {
	fmt.Println("set display scale to", *c.Scale)
	if err := d.SetTraceScale(*c.Scale); err != nil {
		return fmt.Errorf("failed to set scale to %f: %w", *c.Scale, err)
	}
	return nil
}

func (c *LevelCmd) EnableLNA(d *tinysa.Device) error {
	fmt.Println("enable lna")
	if err := d.EnableLNA(); err != nil {
		return fmt.Errorf("failed to enable lna: %w", err)
	}
	return nil
}

func (c *LevelCmd) DisableLNA(d *tinysa.Device) error {
	fmt.Println("disable lna")
	if err := d.DisableLNA(); err != nil {
		return fmt.Errorf("failed to disable lna: %w", err)
	}
	return nil
}
