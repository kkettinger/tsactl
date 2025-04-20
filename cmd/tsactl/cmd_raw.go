package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"strings"
)

type RawCmd struct {
	Command []string `help:"Execute raw command" arg:""`
}

func (c *RawCmd) Run(globals *Globals, ctx *kong.Context) error {
	if len(c.Command) == 0 {
		return ctx.PrintUsage(false)
	}

	d, err := initDevice(globals)
	if err != nil {
		return err
	}

	result, err := d.SendCommand(strings.Join(c.Command, " "))
	if err != nil {
		return fmt.Errorf("failed to send raw command: %w", err)
	}

	if containsBinary(result) || len(result) == 0 {
		fmt.Print(result)
	} else {
		fmt.Println(result)
	}

	return nil
}
