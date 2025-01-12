//go:build darwin

package main

import (
	"encoding/hex"

	"github.com/hybridgroup/go-haystack/lib/findmy"
	"tinygo.org/x/bluetooth"
)

// unknownMAC is a MAC address w use on macOS because there is no way to obtain the actual MAC address of a device.
// see https://developer.radiusnetworks.com/2013/10/21/corebluetooth-doesnt-let-you-see-ibeacons.html
var unknownMAC = bluetooth.MAC{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

func scanDevices(verboseFlag *bool) error {
	bluetooth.DefaultAdapter.Enable()

	return bluetooth.DefaultAdapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if *verboseFlag {
			println("found device:", device.Address.String(), device.RSSI, device.LocalName())
		}

		if device.ManufacturerData() != nil && device.ManufacturerData()[0].CompanyID == findmy.AppleCompanyID {
			status, key, err := findmy.ParseData(unknownMAC, device.ManufacturerData()[0].Data)
			if err != nil {
				println("failed to parse data:", err)
				return
			}
			println(device.Address.String(), device.RSSI, hex.EncodeToString(key), "- battery", findmy.BatteryStatus(status))
		}
	})
}
