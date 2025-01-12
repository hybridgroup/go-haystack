//go:build !darwin

package main

import (
	"encoding/hex"

	"github.com/hybridgroup/go-haystack/lib/findmy"
	"tinygo.org/x/bluetooth"
)

func scanDevices(verboseFlag *bool) error {
	bluetooth.DefaultAdapter.Enable()

	return bluetooth.DefaultAdapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if *verboseFlag {
			println("found device:", device.Address.String(), device.RSSI, device.LocalName())
		}

		if device.ManufacturerData() != nil && device.ManufacturerData()[0].CompanyID == findmy.AppleCompanyID {
			status, key, err := findmy.ParseData(device.Address.MAC, device.ManufacturerData()[0].Data)
			if err != nil {
				println("failed to parse data:", err)
				return
			}
			println(device.Address.String(), device.RSSI, hex.EncodeToString(key), "- battery", findmy.BatteryStatus(status))
		}
	})
}
