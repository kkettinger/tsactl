package main

import (
	"github.com/kkettinger/go-tinysa"
	"log/slog"
)

func initDevice(globals *Globals) (*tinysa.Device, error) {
	var logger *slog.Logger

	// attach custom log handler when debug output is enabled
	if globals.Debug {
		logger = slog.New(&CustomHandler{})
	}

	// try to find device when no port name is given
	if globals.Device == "" {
		return tinysa.FindDevice(
			tinysa.WithBaudRate(globals.Baudrate),
			tinysa.WithLogger(logger))
	}

	return tinysa.NewDevice(globals.Device,
		tinysa.WithBaudRate(globals.Baudrate),
		tinysa.WithLogger(logger))
}
