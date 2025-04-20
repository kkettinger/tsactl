package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/kkettinger/go-tinysa"
)

type DeviceCmd struct {
	Reset            bool  `help:"Reset device" short:"r" group:"Device flags:"`
	ResetDfu         bool  `help:"Reset device in DFU mode" group:"Device flags:"`
	GetId            bool  `help:"Get device id" name:"id" group:"Device flags:"`
	SetId            *uint `help:"Set device id" placeholder:"ID" group:"Device flags:"`
	Battery          bool  `help:"Get battery voltage (mV)" short:"b" name:"bat" group:"Device flags:"`
	BatteryOffset    bool  `help:"Get battery offset voltage (mV)" name:"bat-offset" group:"Device flags:"`
	SetBatteryOffset *uint `help:"Set battery offset voltage (mV)" name:"set-bat-offset" placeholder:"mV" group:"Device flags:"`
	Info             bool  `help:"Get firmware and hardware version" name:"info" short:"i" group:"Device flags:"`
}

func (c *DeviceCmd) Run(globals *Globals, ctx *kong.Context) error {
	var ops []func(*tinysa.Device) error

	if c.Reset {
		ops = append(ops, c.ResetDevice)
	}

	if c.ResetDfu {
		ops = append(ops, c.ResetDeviceDFU)
	}

	if c.GetId {
		ops = append(ops, c.GetDeviceId)
	}

	if c.SetId != nil {
		ops = append(ops, c.SetDeviceId)
	}

	if c.Info {
		ops = append(ops, c.GetInfo)
	}

	if c.Battery {
		ops = append(ops, c.GetBattery)
	}

	if c.BatteryOffset {
		ops = append(ops, c.GetBatteryOffset)
	}

	if c.SetBatteryOffset != nil {
		ops = append(ops, c.SetBatteryOffsetVoltage)
	}

	if len(ops) == 0 {
		_ = ctx.PrintUsage(false)
		return nil
	}

	d, err := initDevice(globals)
	if err != nil {
		return err
	}
	defer d.Close()

	for _, op := range ops {
		if err := op(d); err != nil {
			return err
		}
	}

	return nil
}

func (c *DeviceCmd) ResetDevice(d *tinysa.Device) error {
	fmt.Println("reset device")
	if err := d.Reset(false); err != nil {
		return fmt.Errorf("failed to reset device: %w", err)
	}
	return nil
}

func (c *DeviceCmd) ResetDeviceDFU(d *tinysa.Device) error {
	fmt.Printf("reset device in dfu mode")
	if err := d.Reset(true); err != nil {
		return fmt.Errorf("failed to reset device in dfu mode: %w", err)
	}
	return nil
}

func (c *DeviceCmd) GetDeviceId(d *tinysa.Device) error {
	if id, err := d.GetDeviceID(); err != nil {
		return fmt.Errorf("failed to get device id: %w", err)
	} else {
		fmt.Printf("Device id: %d\n", id)
		return nil
	}
}

func (c *DeviceCmd) SetDeviceId(d *tinysa.Device) error {
	fmt.Println("set device id to", *c.SetId)
	if err := d.SetDeviceID(*c.SetId); err != nil {
		return fmt.Errorf("failed to set device id to %d: %w", *c.SetId, err)
	}
	return nil
}

func (c *DeviceCmd) GetInfo(d *tinysa.Device) error {
	fmt.Println("Model:", d.Model())
	fmt.Println("Firmware version:", d.Version())
	fmt.Println("Hardware version:", d.HardwareVersion())
	return nil
}

func (c *DeviceCmd) GetBattery(d *tinysa.Device) error {
	if bat, err := d.GetBatteryVoltage(); err != nil {
		return fmt.Errorf("failed to get battery voltage: %w", err)
	} else {
		fmt.Printf("Battery voltage: %d mV\n", bat)
	}
	return nil
}

func (c *DeviceCmd) GetBatteryOffset(d *tinysa.Device) error {
	if offset, err := d.GetBatteryOffsetVoltage(); err != nil {
		return fmt.Errorf("failed to get battery offset voltage: %w", err)
	} else {
		fmt.Printf("Battery offset voltage: %d mV\n", offset)
	}
	return nil
}

func (c *DeviceCmd) SetBatteryOffsetVoltage(d *tinysa.Device) error {
	fmt.Println("set battery offset voltage to", *c.SetBatteryOffset)
	if err := d.SetBatteryOffsetVoltage(*c.SetBatteryOffset); err != nil {
		return fmt.Errorf("failed to set set battery voltage to %d: %w", *c.SetBatteryOffset, err)
	}
	return nil
}
