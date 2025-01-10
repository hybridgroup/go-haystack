//go:build clue_alpha

package main

import (
	"tinygo.org/x/tinyfont/proggy"
	"tinygo.org/x/tinyterm"
	"tinygo.org/x/tinyterm/displays"
)

var (
	font = &proggy.TinySZ8pt7b
)

func initTerminal() {
	display := displays.Init()

	terminal = tinyterm.NewTerminal(display)
	terminal.Configure(&tinyterm.Config{
		Font:              font,
		FontHeight:        10,
		FontOffset:        6,
		UseSoftwareScroll: true,
	})
}
