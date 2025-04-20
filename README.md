[![Golang Test](https://github.com/kkettinger/tsactl/actions/workflows/test.yml/badge.svg)](https://github.com/kkettinger/tsactl/actions/workflows/test.yml)
[![Golang CI Lint](https://github.com/kkettinger/tsactl/actions/workflows/lint.yml/badge.svg)](https://github.com/kkettinger/tsactl/actions/workflows/lint.yml)
[![License: GPL-3.0](https://img.shields.io/badge/license-GPL--3.0-blue)](/LICENSE)

# tsactl

`tsactl` is a command line tool for controlling and interacting with the [tinySA](https://www.tinysa.org/) spectrum analyzer via its USB serial interface. It allows you to:

- Configure sweep parameters (frequency range, center, span, ...)
- Configure markers and traces
- Save screenshots in PNG format
- Save trace data to CSV (single/multiple traces)
- Trigger menu options (e.g. to enable waterfall view)
- Reset device (DFU mode for basic model)
- Load/save presets
- Execute raw commands (for when the command is not yet implemented by `tsactl`)
- And more - see the [command overview](#command-overview) for details

It [automatically detects](#device-auto-detection) connected tinySA devices and identifies the model (basic or ultra).

All platforms supported by the `go.bug.st/serial` module _should_ be working: Linux, Windows, macOS.

**Note:** Tested on Linux, Windows and Raspberry Pi with a tinySA Ultra (firmware `tinySA4_v1.4-197-gaa78ccc`). If you experience problems with your device, please [report a bug](#report-bug).


## Installation

You can download the latest `tsactl` binary for your platform from the [release page](https://github.com/kkettinger/tsactl/releases) and place it in a directory that's in your PATH.


## Device auto-detection

If no serial port is defined, the tool will iterate over all serial ports and checks if a device responds by issuing the `version` command.
Based on the `version` response, the model (basic or ultra) is identified (here `tinySA4`, meaning ultra):
```
tinySA4_v1.4-197-gaa78ccc
HW Version:V0.4.5.1
```


## Command overview

| Command         | Alias | Description                                                                      |
|-----------------|-------|----------------------------------------------------------------------------------|
| `tsactl device` | `dev` | Reset device, get device id, battery voltage, hardware and firmware version, ... |
| `tsactl level`  | `lv`  | Change trace unit, reference level, scale, ...                                   |
| `tsactl marker` | `mk`  | Enable/disable marker, assign marker to trace, set frequency, ...                |
| `tsactl menu`   |       | Trigger menu by list of ids                                                      |
| `tsactl preset` | `pr`  | Load and save presets                                                            |
| `tsactl raw`    |       | Execute raw commands                                                             |
| `tsactl save`   |       | Save screenshots as PNG, save trace data as CSV                                  |
| `tsactl signal` | `sig` | Change signal settings like spur removal                                         |
| `tsactl sweep`  | `sw`  | Show and change sweep settings                                                   |
| `tsactl trace`  | `tr`  | Enable/disable traces, trace calculations                                        |

To view all available flags for a command, run: `tsactl command --help`

### Global flags

| Flag           | Description                              | Default     | Env                     |
|----------------|------------------------------------------|-------------|-------------------------|
| `--device, -D` | Device port (e.g. /dev/ttyACM0 or COM1)  | Auto-detect | tsactl_DEVICE           |
| `--baudrate`   | Device baud rate                         | 115200      | tsactl_BAUDRATE         |
| `--debug`      | Debug output                             | False       | tsactl_DEBUG            |


## Example usage

Frequency arguments can be specified using notations like `1.5Ghz`, `1.5g`, `250k`, `250khz`, or in scientific notation such as `1.23e6`.
Most flags support relative arguments, e.g. to move the sweep center frequency by 2mhz, you can use `--center +2mhz`.

Time arguments can be written as `1.2s`, `950ms`, `750m`, `1200us` or `950u`.

### Sweep command

```sh
# Sweep between 410.5mhz and 600mhz
$ tsactl sweep --start 410.5mhz --stop 600mhz

# Check current sweep settings
$ tsactl sweep
Status: resumed
Frequency: 410.5 MHz to 600 MHz (450 points)
Center: 505.25 MHz
Span: 189.5 MHz

# Increase/decrease frequency by using relative arguments 
$ tsactl sweep --center +2mhz

# Change sweep span
$ tsactl sweep --span 20mhz
```

### Trace command

```sh
# Enable second trace and activate max hold
$ tsactl trace 2 --calc maxh

# Enable second marker, assign to trace 2, move to peak frequency and disable delta mode
$ tsactl marker 2 --trace 2 --peak --delta=off
```

### Marker command

```sh
# Set marker 2 to frequency 419mhz
$ tsactl marker 2 --freq 419mhz

# Assign marker 2 to trace 2
$ tsactl marker 2 --trace 2

# Enable/disable marker tracking
$ tsactl marker 2 --track
$ tsactl marker 2 --no-track

# Show details about active markers
$ tsactl marker
Active markers:
  Marker 1:   419 MHz      -89.4   (Index 45)
  Marker 2:   495.32 MHz   -67.9   (Index 214)
```

### Save command

```sh
# Save current screen as PNG
$ tsactl save --capture
Saved capture to SA_250415_183057.png

# Save trace 1 and 2 as CSV
$ tsactl save --trace 1,2
Saved csv to SA_250415_183132.csv

$ cat SA_250415_183132.csv
point,frequency,value_t1,value_t2
0,400000000,-90.25,-84.78
1,400445434,-89.75,-85.75
...
```

Example screenshot created with `tsactl save --capture`:

![Example screenshot](/docs/example_screenshot.png)

Example trace export created with `tsactl save --trace 1,2`: [example_trace_export.csv](/docs/example_trace_export.csv)

### Menu command

```sh
# Enable waterfall view
$ tsactl menu 6 2
```


### Signal command

```sh
# Disable spur removal
$ tsactl signal --spur off
```

### Level command
```sh
# Change trace unit
$ tsactl level --unit vpp

# Change reference level
$ tsactl level --ref -40

# Change scale
$ tsactl level --scale 10
```

### Raw command

If a specific command is missing in `tsactl`, you can execute it with the `tsactl raw` command:

```sh
$ tsactl raw sd_list
SA_250224_235406.bmp 307322
SA_250403_191532.csv 8612
SA_250407_202738.prs 1584
SA_250407_202815.prs 1584
DECT.prs 1584
WIFI.prs 1584
...

# Binary output is supported
$ tsactl raw scanraw 100mhz 120mhz 450 > scanresult.bin
```

If you are familiar with the tinySA serial commands and don't want to use those from `tsactl`, you can use a simple alias to the raw command:
```sh
$ alias tinysa='tsactl raw $@'
$ tinysa sd_list
$ tinysa load 0
$ tinysa sweep 120M 150M
$ tinysa capture > capture.bin
```


## Report bug

If you experience issues with your tinySA model or firmware version, please run `tsactl` in debug mode:
```go
tsactl <command that makes problems> --flags-that-makes-problems --debug
```

When reporting the issue on GitHub, include the debug output if possible - it helps track down the problem. Thanks!


## License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.


## Acknowledgments

- The tinySA team for creating a great spectrum analyzer
- Contributors to the Go serial library
- Contributors to the kong cli library
