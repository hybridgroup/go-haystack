//go:build badger2040_w

package main

import (
	"machine"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyterm"
	"tinygo.org/x/tinyterm/displays"
)

var (
	font = &tinyfont.TomThumb
)

func initTerminal() {
	led3v3 := machine.ENABLE_3V3
	led3v3.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led3v3.High()

	display := displays.Init()

	terminal = tinyterm.NewTerminal(display)
	terminal.Configure(&tinyterm.Config{
		Font:              font,
		FontHeight:        8,
		FontOffset:        6,
		UseSoftwareScroll: true,
	})
}
