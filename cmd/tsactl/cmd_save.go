package main

import (
	"encoding/csv"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/kkettinger/go-tinysa"
	"image/png"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	filenameCaptureDefault    = "SA_<date>_<time>.png"
	filenameTraceDefault      = "SA_<date>_<time>_<trace>.csv"
	filenameTraceMultiDefault = "SA_<date>_<time>.csv"
)

type SaveCmd struct {
	Capture bool   `help:"Save screen as PNG to file" short:"c" group:"Save flags:" `
	Trace   []uint `help:"Save trace(s) as CSV to file" short:"t" group:"Save flags:" `
	Output  string `help:"Output filepath for capture or trace" short:"o" type:"path" group:"Save flags:" placeholder:"PATH"`
}

func (c *SaveCmd) Run(globals *Globals, ctx *kong.Context) error {
	d, err := initDevice(globals)
	if err != nil {
		return err
	}

	if len(c.Trace) == 1 {
		return c.SaveSingleTrace(d)
	}

	if len(c.Trace) > 1 {
		return c.SaveMultipleTraces(d)
	}

	if c.Capture {
		return c.SaveCapture(d)
	}

	_ = ctx.PrintUsage(false)

	return nil
}

func (c *SaveCmd) SaveCapture(d *tinysa.Device) error {
	if c.Output == "" {
		c.Output = filenameCaptureDefault
	}

	// replace filename placeholders with actual values
	c.Output = replaceFilenamePlaceholdersDateTime(c.Output)

	// open file
	file, err := os.Create(c.Output)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", c.Output, err)
	}
	defer file.Close()

	// capture
	img, err := d.Capture()
	if err != nil {
		return fmt.Errorf("failed to capture screen: %w", err)
	}

	// save as PNG
	if err = png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode file '%s': %w", c.Output, err)
	}

	fmt.Printf("capture saved to %s\n", c.Output)

	return nil
}

func (c *SaveCmd) SaveSingleTrace(d *tinysa.Device) error {
	if c.Output == "" {
		c.Output = filenameTraceDefault
	}

	// replace filename placeholders with actual values
	c.Output = replaceFilenamePlaceholdersDateTime(c.Output)
	c.Output = strings.ReplaceAll(c.Output, "<trace>", fmt.Sprintf("%d", c.Trace[0]))

	// get data from device
	data, err := d.GetTraceData(c.Trace[0])
	if err != nil {
		return fmt.Errorf("failed to get trace data: %w", err)
	}

	// open file
	file, err := os.Create(c.Output)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", c.Output, err)
	}
	defer file.Close()

	// write as CSV
	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"trace", "point", "frequency", "value"}); err != nil {
		return err
	}

	for _, dp := range data {
		row := []string{
			strconv.FormatUint(uint64(dp.Trace), 10),
			strconv.FormatUint(uint64(dp.Point), 10),
			strconv.FormatUint(dp.Frequency, 10),
			strconv.FormatFloat(dp.Value, 'f', -1, 64),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	fmt.Printf("trace %d data saved to %s\n", c.Trace[0], c.Output)

	return nil
}

func (c *SaveCmd) SaveMultipleTraces(d *tinysa.Device) error {
	if c.Output == "" {
		c.Output = filenameTraceMultiDefault
	}

	// replace filename placeholders with actual values
	c.Output = replaceFilenamePlaceholdersDateTime(c.Output)

	data := make([][]tinysa.TraceData, len(c.Trace))
	for i, traceId := range c.Trace {
		d, err := d.GetTraceData(traceId)
		if err != nil {
			return fmt.Errorf("failed to get trace data: %w", err)
		}
		data[i] = d
	}

	// open file
	file, err := os.Create(c.Output)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", c.Output, err)
	}
	defer file.Close()

	// write as CSV
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// header
	header := []string{"point", "frequency"}
	for _, t := range c.Trace {
		header = append(header, fmt.Sprintf("value_t%d", t))
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	for i, dp := range data[0] {
		row := []string{
			strconv.FormatUint(uint64(dp.Point), 10),
			strconv.FormatUint(dp.Frequency, 10),
		}
		for j := range c.Trace {
			row = append(row, strconv.FormatFloat(data[j][i].Value, 'f', -1, 64))
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	fmt.Printf("traces %v saved to %s\n", c.Trace, c.Output)

	return nil
}

func replaceFilenamePlaceholdersDateTime(str string) string {
	now := time.Now()
	dateStr := now.Format("060102")
	timeStr := now.Format("150405")
	str = strings.ReplaceAll(str, "<date>", dateStr)
	str = strings.ReplaceAll(str, "<time>", timeStr)
	return str
}
