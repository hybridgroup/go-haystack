//go:build pybadge

package main

import (
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyterm"
	"tinygo.org/x/tinyterm/displays"
)

var (
	font = &tinyfont.Picopixel
)

func initTerminal() {
	display := displays.Init()

	terminal = tinyterm.NewTerminal(display)
	terminal.Configure(&tinyterm.Config{
		Font:              font,
		FontHeight:        8,
		FontOffset:        4,
		UseSoftwareScroll: true,
	})
}
