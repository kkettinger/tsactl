package main

import (
	"fmt"
	"github.com/kkettinger/go-tinysa"
	"os"
	"text/tabwriter"
)

type TraceCmd struct {
	Enable  bool      `help:"Enable trace" short:"e" group:"Trace flags:"`
	Disable bool      `help:"Disable trace" short:"d" group:"Trace flags:"`
	Calc    TraceCalc `help:"Enable trace calculation (${trace_calc_opts})" short:"c" placeholder:"MODE" group:"Trace flags:"`

	Trace uint `arg:"" name:"id" help:"Trace id" optional:""`
}

func (c *TraceCmd) Run(globals *Globals) error {
	var ops []func(*tinysa.Device) error

	hasTrace := c.Trace != 0

	if c.Enable {
		ops = append(ops, c.EnableTrace)
	}

	if c.Disable {
		ops = append(ops, c.DisableTrace)
	}

	if c.Calc.Valid {
		if c.Calc.Off {
			ops = append(ops, c.DisableTraceCalc)
		} else {
			ops = append(ops, c.EnableTraceCalc)
		}
	}

	if !hasTrace && len(ops) > 0 {
		return fmt.Errorf("expected \"<id>\"")
	}

	d, err := initDevice(globals)
	if err != nil {
		return err
	}

	// apply operations on trace
	if hasTrace && len(ops) > 0 {
		for _, op := range ops {
			if err := op(d); err != nil {
				return err
			}
		}
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// show details about specific trace
	if hasTrace {
		t, err := d.GetTrace(c.Trace)
		if err != nil {
			return err
		}
		printTraceInfo(w, t)
		_ = w.Flush()
		return nil
	}

	// default: list details about all active traces
	fmt.Println("Active traces:")
	tAll, err := d.GetTraceAll()
	if err != nil {
		return err
	}
	for _, t := range tAll {
		printTraceInfo(w, t)
	}
	_ = w.Flush()

	return nil
}

func printTraceInfo(w *tabwriter.Writer, t tinysa.Trace) {
	_, _ = fmt.Fprintf(w, "  Trace %d:\t%s\t%f\t%f\n", t.Trace, t.Unit, t.Scale, t.RefPos)
}

func (c *TraceCmd) EnableTrace(d *tinysa.Device) error {
	fmt.Printf("enable trace #%d\n", c.Trace)

	if err := d.EnableTrace(c.Trace); err != nil {
		return fmt.Errorf("failed to enable trace #%d: %w", c.Trace, err)
	}

	return nil
}

func (c *TraceCmd) DisableTrace(d *tinysa.Device) error {
	fmt.Printf("disable trace #%d\n", c.Trace)

	if err := d.DisableTrace(c.Trace); err != nil {
		return fmt.Errorf("failed to disable trace #%d: %w", c.Trace, err)
	}

	return nil
}

func (c *TraceCmd) DisableTraceCalc(d *tinysa.Device) error {
	fmt.Printf("disable calculations on trace #%d\n", c.Trace)

	if err := d.DisableTraceCalc(c.Trace); err != nil {
		return fmt.Errorf("failed to disable calculation for trace %d: %w", c.Trace, err)
	}

	return nil
}

func (c *TraceCmd) EnableTraceCalc(d *tinysa.Device) error {
	fmt.Printf("enable trace calculations %s for trace #%d\n", c.Calc.Mode.String(), c.Trace)

	if err := d.EnableTrace(c.Trace); err != nil {
		return fmt.Errorf("failed to enable trace #%d: %w", c.Trace, err)
	}

	if err := d.EnableTraceCalc(c.Trace, c.Calc.Mode); err != nil {
		return fmt.Errorf("failed to enable trace calculation %s for trace %d: %w", c.Calc.Mode.String(), c.Trace, err)
	}

	return nil
}
