package main

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"time"

	"github.com/hybridgroup/go-haystack/lib/findmy"
	"tinygo.org/x/bluetooth"
	"tinygo.org/x/tinyterm"
)

var (
	terminal *tinyterm.Terminal

	black   = color.RGBA{0, 0, 0, 255}
	adapter = bluetooth.DefaultAdapter
)

func main() {
	initTerminal()

	terminalOutput("enable interface...")

	must("enable BLE interface", adapter.Enable())
	time.Sleep(time.Second)

	terminalOutput("start scan...")

	must("start scan", adapter.Scan(scanHandler))

	for {
		time.Sleep(time.Minute)
		terminalOutput("scanning...")
	}
}

func scanHandler(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
	if device.ManufacturerData() != nil && device.ManufacturerData()[0].CompanyID == findmy.AppleCompanyID {
		status, key, err := findmy.ParseData(device.Address, device.ManufacturerData()[0].Data)
		switch {
		case err != nil:
			terminalOutput("ERROR: failed to parse data:" + err.Error())
		default:
			terminalOutput("--------------------------------")
			terminalOutput(fmt.Sprintf("%s %d (battery %s)", device.Address.String(), device.RSSI, findmy.BatteryStatus(status)))
			terminalOutput(hex.EncodeToString(key))
		}
	}
}

func must(action string, err error) {
	if err != nil {
		for {
			terminalOutput("failed to " + action + ": " + err.Error())

			time.Sleep(time.Second)
		}
	}
}

func terminalOutput(s string) {
	println(s)
	fmt.Fprintf(terminal, "\n%s", s)

	terminal.Display()
}
