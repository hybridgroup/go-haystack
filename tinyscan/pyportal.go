//go:build pyportal

package main

import (
	"tinygo.org/x/drivers/ili9341"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyterm"
	"tinygo.org/x/tinyterm/displays"
)

var (
	font = &tinyfont.TomThumb
)

func initTerminal() {
	display := displays.Init()
	display.SetRotation(ili9341.Rotation270)

	terminal = tinyterm.NewTerminal(display)
	terminal.Configure(&tinyterm.Config{
		Font:              font,
		FontHeight:        8,
		FontOffset:        6,
		UseSoftwareScroll: true,
	})
}
