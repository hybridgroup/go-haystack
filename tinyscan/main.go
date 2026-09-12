//go:build badger2040_w || clue_alpha || pybadge || pyportal

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

	showErrors string

	// MyDevices holds the devices of the owner, such as
	// "name=KEY1,KEY2;name2=KEY3". The keys are the base64 advertisement keys.
	// Set it with -ldflags "-X main.MyDevices=...".
	MyDevices string

	// OnlyMine hides every beacon that is not a device of the owner. Set it to
	// "true" with -ldflags "-X main.OnlyMine=true".
	OnlyMine string

	devices []findmy.Device
)

func main() {
	initTerminal()

	// A bad list of devices must not stop the scan.
	var err error
	devices, err = findmy.ParseDevices(MyDevices)
	if err != nil {
		terminalOutput("ERROR: failed to parse devices:" + err.Error())
	}
	terminalOutput(fmt.Sprintf("known devices: %d", len(devices)))

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
		status, key, err := findmy.ParseData(device.Address.MAC, device.ManufacturerData()[0].Data)
		switch {
		case err != nil && err == findmy.ErrorUnregistered:
			if OnlyMine != "" {
				return
			}
			terminalOutput("--------------------------------")
			terminalOutput(fmt.Sprintf("%s %d (unregistered)", device.Address.String(), device.RSSI))
			return
		case err != nil:
			if showErrors != "" {
				terminalOutput("--------------------------------")
				terminalOutput("ERROR: failed to parse data:" + err.Error())
			}
			return
		}

		name, mine := findmy.Lookup(devices, key)
		if !mine && OnlyMine != "" {
			return
		}

		terminalOutput("--------------------------------")
		terminalOutput(fmt.Sprintf("%s %d (battery %s)", device.Address.String(), device.RSSI, findmy.BatteryStatus(status)))
		if mine {
			terminalOutput("* " + name)
			return
		}
		terminalOutput(hex.EncodeToString(key))
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
