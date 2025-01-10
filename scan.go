package main

import (
	"encoding/hex"

	"errors"

	"tinygo.org/x/bluetooth"
)

const (
	// Apple, Inc.
	appleCompanyID = 0x004C

	// Offline Finding type
	findMyPayloadType = 0x12

	// Length of the payload
	findMyPayloadLength = 0x19

	// Hint byte
	findMyHint = 0x00
)

func scanDevices(verboseFlag *bool) error {
	bluetooth.DefaultAdapter.Enable()

	return bluetooth.DefaultAdapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if *verboseFlag {
			println("found device:", device.Address.String(), device.RSSI, device.LocalName())
		}

		if device.ManufacturerData() != nil && device.ManufacturerData()[0].CompanyID == appleCompanyID {
			status, key, err := parseData(device.Address, device.ManufacturerData()[0].Data)
			if err != nil {
				println("failed to parse data:", err)
				return
			}
			println(device.Address.String(), device.RSSI, hex.EncodeToString([]byte{status}), hex.EncodeToString(key))
		}
	})
}

func parseData(address bluetooth.Address, data []byte) (byte, []byte, error) {
	if len(data) < 27 {
		return 0, nil, errors.New("data is too short")
	}

	if data[0] != findMyPayloadType {
		return 0, nil, errors.New("invalid payload type")
	}

	if data[1] != findMyPayloadLength {
		return 0, nil, errors.New("invalid payload length")
	}

	if data[26] != findMyHint {
		return 0, nil, errors.New("invalid hint")
	}

	findMyStatus := data[2]
	var key [28]byte
	copy(key[6:], data[3:25])

	// turn address into key bytes
	key[0] = address.MAC[5]
	key[1] = address.MAC[4]
	key[2] = address.MAC[3]
	key[3] = address.MAC[2]
	key[4] = address.MAC[1]
	key[5] = address.MAC[0]

	return findMyStatus, key[:], nil
}
