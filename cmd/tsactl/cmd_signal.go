package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/kkettinger/go-tinysa"
)

type SignalCmd struct {
	Spur Spur `help:"Set spur removal (on, off, auto)" group:"Signal flags:"`
}

func (c *SignalCmd) Run(globals *Globals, ctx *kong.Context) error {
	var ops []func(*tinysa.Device) error

	if c.Spur.Valid {
		switch {
		case c.Spur.Auto:
			ops = append(ops, c.EnableAutoSpur)
		case c.Spur.Enable:
			ops = append(ops, c.EnableSpur)
		default:
			ops = append(ops, c.DisableSpur)
		}
	}

	if len(ops) > 0 {
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

	_ = ctx.PrintUsage(false)

	return nil
}

func (c *SignalCmd) EnableSpur(d *tinysa.Device) error {
	fmt.Println("enable spur removal")
	if err := d.EnableSpurRemoval(); err != nil {
		return fmt.Errorf("failed to enable spur removal: %w", err)
	}
	return nil
}

func (c *SignalCmd) DisableSpur(d *tinysa.Device) error {
	fmt.Println("disable spur removal")
	if err := d.DisableSpurRemoval(); err != nil {
		return fmt.Errorf("failed to disable spur removal: %w", err)
	}
	return nil
}

func (c *SignalCmd) EnableAutoSpur(d *tinysa.Device) error {
	fmt.Println("enable auto spur removal")
	if err := d.EnableAutoSpurRemoval(); err != nil {
		return fmt.Errorf("failed to enable auto spur removal: %w", err)
	}
	return nil
}
