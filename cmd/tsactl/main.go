package main

import (
	"github.com/alecthomas/kong"
	"os"
	"strings"
)

type Globals struct {
	Device   string `help:"Device serial port, e.g. /dev/ttyACM0 or COM1" short:"D" placeholder:"PORT" env:"tsactl_DEVICE"`
	Baudrate int    `help:"Device baudrate rate" default:"115200" env:"tsactl_BAUDRATE"`
	Debug    bool   `help:"Enable debug output" env:"tsactl_DEBUG"`
}

type Cli struct {
	Globals

	Version kong.VersionFlag `help:"Show tsactl version" short:"v"`

	Device DeviceCmd `help:"Access device status, ID, battery, and firmware info" cmd:"" aliases:"dev"`
	Level  LevelCmd  `help:"Set trace unit, reference level, and scale" cmd:"" aliases:"lv"`
	Marker MarkerCmd `help:"Enable marker, set frequency, and tracking" cmd:"" aliases:"mk"`
	Menu   MenuCmd   `help:"Trigger menu actions by ID" cmd:""`
	Preset PresetCmd `help:"Load or save device presets" cmd:"" aliases:"pr"`
	Raw    RawCmd    `help:"Send low-level raw commands" cmd:""`
	Save   SaveCmd   `help:"Export screen capture or trace data to file" cmd:""`
	Signal SignalCmd `help:"Configure signal processing options" cmd:"" aliases:"sig"`
	Sweep  SweepCmd  `help:"Set sweep parameters like freq range and mode" cmd:"" aliases:"sw"`
	Trace  TraceCmd  `help:"Enable traces and set calculation modes" cmd:"" aliases:"tr"`
}

var cli Cli

func main() {
	// If running without any extra arguments, default to the --help flag
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}

	// Variables
	traceCalc := TraceCalc{}
	traceCalcOpts := strings.Join(traceCalc.ValidOpts(), ", ")

	traceUnit := TraceUnit{}
	traceUnitOpts := strings.Join(traceUnit.ValidOpts(), ", ")

	sweepMode := SweepMode{}
	sweepModeOpts := strings.Join(sweepMode.ValidOpts(), ", ")

	ctx := kong.Parse(&cli,
		kong.Name("tsactl"),
		kong.Description("Command line tool for the tinySA spectrum analyzer."),
		kong.Vars{"version": getVersion()},
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"trace_calc_opts": traceCalcOpts,
			"trace_unit_opts": traceUnitOpts,
			"sweep_mode_opts": sweepModeOpts,
		})

	err := ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)
}
