package main

import (
	"fmt"
	"github.com/kkettinger/go-tinysa"
	"github.com/kkettinger/tsactl/internal/util"
)

type SweepCmd struct {
	Pause        bool         `help:"Pause sweep" short:"p" group:"Sweep flags:"`
	Resume       bool         `help:"Resume sweep" short:"r" group:"Sweep flags:"`
	Mode         SweepMode    `help:"Sweep mode (${sweep_mode_opts})" short:"m" group:"Sweep flags:" placeholder:"MODE"`
	Start        FrequencyRel `help:"Start frequency" short:"s" group:"Sweep flags:" placeholder:"FREQ"`
	Stop         FrequencyRel `help:"Stop frequency" short:"e" group:"Sweep flags:" placeholder:"FREQ"`
	Span         FrequencyRel `help:"Span frequency" short:"S" group:"Sweep flags:" placeholder:"FREQ"`
	Center       FrequencyRel `help:"Center frequency" short:"C" group:"Sweep flags:" placeholder:"FREQ"`
	CenterMarker *uint        `help:"Set center frequency from marker" short:"M" group:"Sweep flags:" placeholder:"MARKER"`
	Points       *uint        `help:"Number of sweep points" short:"n" group:"Sweep flags:"`
	Time         Time         `help:"Sweep time" short:"t" group:"Sweep flags:"`
	CW           Frequency    `help:"Set continuous wave frequency" group:"Sweep flags:" placeholder:"FREQ"`
}

func (c *SweepCmd) Run(globals *Globals) error {
	var ops []func(*tinysa.Device) error

	if c.Pause {
		ops = append(ops, c.PauseSweep)
	}

	if c.Resume {
		ops = append(ops, c.ResumeSweep)
	}

	if c.Mode.Valid {
		ops = append(ops, c.SetSweepMode)
	}

	if c.Start.Valid {
		ops = append(ops, c.SetSweepStart)
	}

	if c.Stop.Valid {
		ops = append(ops, c.SetSweepStop)
	}

	if c.Span.Valid {
		ops = append(ops, c.SetSweepSpan)
	}

	if c.Center.Valid {
		ops = append(ops, c.SetSweepCenter)
	}

	if c.CenterMarker != nil {
		ops = append(ops, c.SetSweepCenterFromMarker)
	}

	if c.Points != nil {
		ops = append(ops, c.SetSweepPoints)
	}

	if c.Time.Valid {
		ops = append(ops, c.SetSweepTime)
	}

	if c.CW.Valid {
		ops = append(ops, c.SetSweepCW)
	}

	d, err := initDevice(globals)
	if err != nil {
		return err
	}

	if len(ops) > 0 {
		for _, op := range ops {
			if err := op(d); err != nil {
				return err
			}
		}
		return nil
	}

	return c.Status(d)
}

func (c *SweepCmd) Status(d *tinysa.Device) error {
	state, err := d.GetSweepStatus()
	if err != nil {
		return err
	}

	sweep, err := d.GetSweep()
	if err != nil {
		return err
	}

	fmt.Printf("Status: %s\n", state)
	if sweep.Start == sweep.Stop {
		fmt.Printf("Frequency: %s (CW)\n",
			util.FormatFrequency(sweep.Start))
	} else {
		span := sweep.Stop - sweep.Start
		center := (span / 2) + sweep.Start
		fmt.Printf("Frequency: %s to %s (%d points)\n",
			util.FormatFrequency(sweep.Start), util.FormatFrequency(sweep.Stop), sweep.Points)
		fmt.Printf("Center: %s\n", util.FormatFrequency(center))
		fmt.Printf("Span: %s\n", util.FormatFrequency(span))
	}

	return nil
}

func (c *SweepCmd) PauseSweep(d *tinysa.Device) error {
	fmt.Println("pause sweep")
	if err := d.PauseSweep(); err != nil {
		return fmt.Errorf("failed to pause sweep: %w", err)
	}
	return nil
}

func (c *SweepCmd) ResumeSweep(d *tinysa.Device) error {
	fmt.Println("resume sweep")
	if err := d.ResumeSweep(); err != nil {
		return fmt.Errorf("failed to resume sweep: %w", err)
	}
	return nil
}

func (c *SweepCmd) SetSweepMode(d *tinysa.Device) error {
	fmt.Printf("set sweep mode to %s\n", c.Mode.Mode.String())
	if err := d.SetSweepMode(c.Mode.Mode); err != nil {
		return fmt.Errorf("failed to set sweep mode to %s: %w", c.Mode.Mode.String(), err)
	}
	return nil
}

func (c *SweepCmd) SetSweepStart(d *tinysa.Device) error {
	freq := c.Start.Value

	if c.Start.Relative {
		sweep, err := d.GetSweep()
		if err != nil {
			return err
		}
		freq += int64(sweep.Start) // #nosec G115
	}

	if freq < 0 {
		return fmt.Errorf("invalid frequency: %d", freq)
	}

	uFreq := uint64(freq)

	fmt.Printf("set sweep start frequency to %s\n", util.FormatFrequency(uFreq))
	if err := d.SetSweepStart(uFreq); err != nil {
		return fmt.Errorf("failed to set sweep start frequency to %s: %w", util.FormatFrequency(uFreq), err)
	}
	return nil
}

func (c *SweepCmd) SetSweepStop(d *tinysa.Device) error {
	freq := c.Stop.Value

	if c.Stop.Relative {
		sweep, err := d.GetSweep()
		if err != nil {
			return err
		}
		freq += int64(sweep.Stop) // #nosec G115
	}

	if freq < 0 {
		return fmt.Errorf("invalid frequency: %d", freq)
	}

	uFreq := uint64(freq)

	fmt.Printf("set sweep stop frequency to %s\n", util.FormatFrequency(uFreq))
	if err := d.SetSweepStop(uFreq); err != nil {
		return fmt.Errorf("failed to set sweep stop frequency to %s: %w", util.FormatFrequency(uFreq), err)
	}
	return nil
}

func (c *SweepCmd) SetSweepSpan(d *tinysa.Device) error {
	freq := c.Span.Value

	if c.Span.Relative {
		sweep, err := d.GetSweep()
		if err != nil {
			return err
		}
		freq += int64(sweep.Stop - sweep.Start) // #nosec G115
	}

	if freq < 0 {
		return fmt.Errorf("invalid frequency: %d", freq)
	}

	uFreq := uint64(freq)

	fmt.Printf("set sweep span frequency to %s\n", util.FormatFrequency(uFreq))
	if err := d.SetSweepSpan(uFreq); err != nil {
		return fmt.Errorf("failed to set sweep span frequency to %s: %w", util.FormatFrequency(uFreq), err)
	}
	return nil
}

func (c *SweepCmd) SetSweepCenter(d *tinysa.Device) error {
	freq := c.Center.Value

	if c.Center.Relative {
		sweep, err := d.GetSweep()
		if err != nil {
			return err
		}
		freq += int64(sweep.Start + (sweep.Stop-sweep.Start)/2) // #nosec G115
	}

	if freq < 0 {
		return fmt.Errorf("invalid frequency: %d", freq)
	}

	uFreq := uint64(freq)

	fmt.Printf("set sweep center frequency to %s\n", util.FormatFrequency(uFreq))
	if err := d.SetSweepCenter(uFreq); err != nil {
		return fmt.Errorf("failed to set sweep center frequency to %s: %w", util.FormatFrequency(uFreq), err)
	}
	return nil
}

func (c *SweepCmd) SetSweepCenterFromMarker(d *tinysa.Device) error {
	marker, err := d.GetMarker(*c.CenterMarker)
	if err != nil {
		return fmt.Errorf("failed to get marker #%d: %w", *c.CenterMarker, err)
	}
	fmt.Printf("set sweep center frequency from marker #%d (%s)\n", *c.CenterMarker, util.FormatFrequency(marker.Frequency))
	if err := d.SetSweepCenter(marker.Frequency); err != nil {
		return fmt.Errorf("failed to set sweep center frequency to %s: %w", util.FormatFrequency(marker.Frequency), err)
	}
	return nil
}

func (c *SweepCmd) SetSweepPoints(d *tinysa.Device) error {
	fmt.Printf("set sweep points to %d\n", *c.Points)
	if err := d.SetSweepPoints(*c.Points); err != nil {
		return fmt.Errorf("failed to set sweep points to %d: %w", *c.Points, err)
	}
	return nil
}

func (c *SweepCmd) SetSweepTime(d *tinysa.Device) error {
	fmt.Printf("set sweep time to %s\n", util.FormatTimeDuration(c.Time.Value))
	if err := d.SetSweepTime(c.Time.Value); err != nil {
		return fmt.Errorf("failed to set sweep time to %s: %w", util.FormatTimeDuration(c.Time.Value), err)
	}
	return nil
}

func (c *SweepCmd) SetSweepCW(d *tinysa.Device) error {
	fmt.Printf("Setting sweep cw frequency to %s\n", util.FormatFrequency(c.CW.Value))
	if err := d.SetSweepContinuousWave(c.CW.Value); err != nil {
		return fmt.Errorf("failed to set sweep cw frequency to %s: %w", util.FormatFrequency(c.CW.Value), err)
	}
	return nil
}
