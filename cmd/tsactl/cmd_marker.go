package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/kkettinger/go-tinysa"
	"github.com/kkettinger/tsactl/internal/util"
	"os"
	"text/tabwriter"
)

type MarkerCmd struct {
	Enable    bool         `help:"Enable marker" short:"e" group:"Marker flags:"`
	Disable   bool         `help:"Disable marker" short:"d" group:"Marker flags:"`
	Trace     *uint        `help:"Assign marker to trace" short:"t" group:"Marker flags:" placeholder:"TRACE"`
	Frequency FrequencyRel `help:"Set marker to frequency" name:"freq" short:"f" group:"Marker flags:" placeholder:"FREQ"`
	Peak      bool         `help:"Move marker to peak of assigned trace" short:"p" group:"Marker flags:"`
	Delta     MarkerDelta  `help:"Enable delta mode (off or reference marker)" group:"Marker flags:" placeholder:"<OFF|MARKER>"`
	Tracking  *bool        `help:"Enable tracking mode" name:"track" negatable:"" group:"Marker flags:"`

	Marker uint `arg:"" name:"id" help:"Marker id" optional:""`
}

func (c *MarkerCmd) Validate(ctx *kong.Context) error {
	if c.Enable && c.Disable {
		return fmt.Errorf("--enable,e and --disable,d cannot be set at the same time")
	}

	return nil
}

func (c *MarkerCmd) Run(globals *Globals, ctx *kong.Context) error {
	var ops []func(*tinysa.Device) error

	hasMarker := c.Marker != 0

	if c.Enable {
		ops = append(ops, c.EnableMarker)
	}

	if c.Disable {
		ops = append(ops, c.DisableMarker)
	}

	if c.Trace != nil {
		ops = append(ops, c.AssignTrace)
	}

	if c.Frequency.Valid {
		ops = append(ops, c.SetFrequency)
	}

	if c.Delta.Valid && c.Delta.Off {
		ops = append(ops, c.DisableDelta)
	}

	if c.Delta.Valid && c.Delta.RefMarker > 0 {
		ops = append(ops, c.EnableDelta)
	}

	if c.Peak {
		ops = append(ops, c.SetPeak)
	}

	if c.Tracking != nil {
		if *c.Tracking {
			ops = append(ops, c.EnableTracking)
		} else {
			ops = append(ops, c.DisableTracking)
		}
	}

	if !hasMarker && len(ops) > 0 {
		return fmt.Errorf("expected \"<id>\"")
	}

	d, err := initDevice(globals)
	if err != nil {
		return err
	}

	// apply operations on marker
	if hasMarker && len(ops) > 0 {
		for _, op := range ops {
			if err := op(d); err != nil {
				return err
			}
		}
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// show details about a specific marker when no flags are given
	if hasMarker {
		m, err := d.GetMarker(c.Marker)
		if err != nil {
			return err
		}
		printMarkerInfo(w, m)
		_ = w.Flush()
		return nil
	}

	// default: list details about all active marker
	fmt.Println("Active markers:")
	mAll, err := d.GetMarkerAll()
	if err != nil {
		return err
	}

	for _, m := range mAll {
		printMarkerInfo(w, m)
	}
	_ = w.Flush()

	return nil
}

func printMarkerInfo(w *tabwriter.Writer, m tinysa.Marker) {
	_, _ = fmt.Fprintf(w, "  Marker %d:\t%s\t%g\t(Index %d)\n", m.Marker, util.FormatFrequency(m.Frequency), m.Value, m.Index)
}

func (c *MarkerCmd) EnableMarker(d *tinysa.Device) error {
	fmt.Printf("enable marker #%d\n", c.Marker)
	if err := d.EnableMarker(c.Marker); err != nil {
		return fmt.Errorf("failed to enable marker #%d: %w", c.Marker, err)
	}
	return nil
}

func (c *MarkerCmd) DisableMarker(d *tinysa.Device) error {
	fmt.Printf("disable marker #%d\n", c.Marker)
	if err := d.DisableMarker(c.Marker); err != nil {
		return fmt.Errorf("failed to disable marker #%d: %w", c.Marker, err)
	}
	return nil
}

func (c *MarkerCmd) AssignTrace(d *tinysa.Device) error {
	fmt.Printf("assign marker #%d to trace #%d\n", c.Marker, *c.Trace)
	if err := d.SetMarkerTrace(c.Marker, *c.Trace); err != nil {
		return fmt.Errorf("failed to assign marker #%d to trace #%d: %w", c.Marker, *c.Trace, err)
	}
	return nil
}

func (c *MarkerCmd) SetFrequency(d *tinysa.Device) error {
	freq := c.Frequency.Value

	if c.Frequency.Relative {
		marker, err := d.GetMarker(c.Marker)
		if err != nil {
			return err
		}
		freq += int64(marker.Frequency) // #nosec G115
	}

	if freq < 0 {
		return fmt.Errorf("invalid frequency: %d", freq)
	}

	uFreq := uint64(freq)

	fmt.Printf("set marker #%d to frequency %s\n", c.Marker, util.FormatFrequency(uFreq))
	if err := d.SetMarkerFreq(c.Marker, uFreq); err != nil {
		return fmt.Errorf("failed to set marker #%d to frequency #%d: %w", c.Marker, uFreq, err)
	}
	return nil
}

func (c *MarkerCmd) EnableDelta(d *tinysa.Device) error {
	fmt.Printf("enable delta mode for marker #%d relative to marker #%d\n", c.Marker, c.Delta.RefMarker)
	if err := d.EnableMarkerDelta(c.Marker, c.Delta.RefMarker); err != nil {
		return fmt.Errorf("failed to enable delta mode for marker #%d: %w", c.Marker, err)
	}
	return nil
}

func (c *MarkerCmd) DisableDelta(d *tinysa.Device) error {
	fmt.Printf("disable delta mode for marker #%d\n", c.Marker)
	if err := d.DisableMarkerDelta(c.Marker); err != nil {
		return fmt.Errorf("failed to disable delta mode for marker #%d: %w", c.Marker, err)
	}
	return nil
}

func (c *MarkerCmd) SetPeak(d *tinysa.Device) error {
	fmt.Printf("set marker #%d to peak\n", c.Marker)
	if err := d.MoveMarkerPeak(c.Marker); err != nil {
		return fmt.Errorf("failed to set marker #%d to peak: %w", c.Marker, err)
	}
	return nil
}

func (c *MarkerCmd) EnableTracking(d *tinysa.Device) error {
	fmt.Printf("enable tracking for marker #%d\n", c.Marker)
	if err := d.EnableMarkerTracking(c.Marker); err != nil {
		return fmt.Errorf("failed to enable tracking for marker #%d: %w", c.Marker, err)
	}
	return nil
}

func (c *MarkerCmd) DisableTracking(d *tinysa.Device) error {
	fmt.Printf("disable tracking for marker #%d\n", c.Marker)
	if err := d.DisableMarkerTracking(c.Marker); err != nil {
		return fmt.Errorf("failed to disable tracking for marker #%d: %w", c.Marker, err)
	}
	return nil
}
