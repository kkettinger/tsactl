package main

import (
	"fmt"
	"strings"
)

type MenuCmd struct {
	MenuIds []uint `arg:"" help:"Trigger menu options" name:"id"`
}

func (c *MenuCmd) Run(globals *Globals) error {
	d, err := initDevice(globals)
	if err != nil {
		return err
	}

	idsStr := make([]string, len(c.MenuIds))
	for i, v := range c.MenuIds {
		idsStr[i] = fmt.Sprintf("%d", v)
	}
	fmt.Printf("trigger menu %s\n", strings.Join(idsStr, ", "))
	if err := d.TriggerMenu(c.MenuIds); err != nil {
		return fmt.Errorf("failed to trigger menu: %w", err)
	}

	return nil
}
