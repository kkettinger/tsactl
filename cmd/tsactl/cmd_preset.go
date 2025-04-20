package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/kkettinger/go-tinysa"
)

type PresetCmd struct {
	Load *uint `help:"Load preset (0 = startup)" short:"l" group:"Preset flags:"`
	Save *uint `help:"Save preset (0 = startup)" short:"s" group:"Preset flags:"`
}

func (c *PresetCmd) Validate() error {
	if c.Load != nil && c.Save != nil {
		return fmt.Errorf("--load,l and --save,s cannot be called at the same time")
	}

	return nil
}

func (c *PresetCmd) Run(globals *Globals, ctx *kong.Context) error {

	var ops []func(d *tinysa.Device) error

	if c.Load != nil {
		ops = append(ops, c.LoadPreset)
	}

	if c.Save != nil {
		ops = append(ops, c.SavePreset)
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

func (c *PresetCmd) LoadPreset(d *tinysa.Device) error {
	fmt.Printf("load preset %d\n", *c.Load)
	if err := d.LoadPreset(*c.Load); err != nil {
		return fmt.Errorf("failed to load preset %d: %w", *c.Load, err)
	}
	return nil
}

func (c *PresetCmd) SavePreset(d *tinysa.Device) error {
	fmt.Printf("save preset %d\n", *c.Save)
	if err := d.SavePreset(*c.Save); err != nil {
		return fmt.Errorf("failed to save preset %d: %w", *c.Save, err)
	}
	return nil
}
