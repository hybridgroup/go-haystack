// Firmware to advertise a FindMy compatible device.
// aka AirTag
// see https://github.com/biemster/FindMy
//
// To build:
// tinygo flash -target nano-rp2040 -ldflags="-X main.AdvertisingKey='SGVsbG8sIFdvcmxkIQ=='" ./examples/findmy
//
// For Linux:
// go run ./examples/findmy SGVsbG8sIFdvcmxkIQ==
package main

import (
	"encoding/base64"
	"errors"
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func main() {
	time.Sleep(2 * time.Second) // wait for USB serial to be available

	key, err := getKeyData()
	if err != nil {
		panic("failed to get key data: " + err.Error())
	}
	println("key is", AdvertisingKey, "(", len(key), "bytes)")

	opts := bluetooth.AdvertisementOptions{
		AdvertisementType: bluetooth.AdvertisingTypeNonConnInd,
		Interval:          bluetooth.NewDuration(1285000 * time.Microsecond), // 1285ms
		ManufacturerData:  []bluetooth.ManufacturerDataElement{findMyData(key)},
	}

	must("enable BLE stack", adapter.Enable())
	adapter.SetRandomAddress(bluetooth.MAC{key[5], key[4], key[3], key[2], key[1], key[0] | 0xC0})
	adv := adapter.DefaultAdvertisement()

	println("advertising...")
	must("config adv", adv.Configure(opts))

	println("starting...")
	must("start adv", adv.Start())

	address, _ := adapter.Address()
	for {
		println("FindMy /", address.MAC.String())
		time.Sleep(time.Second)
	}
}

func getKeyData() ([]byte, error) {
	val, err := base64.StdEncoding.DecodeString(AdvertisingKey)
	if err != nil {
		return nil, err
	}
	if len(val) != 28 {
		return nil, errors.New("public key must be 28 bytes long")
	}

	return val, nil
}

func findMyData(keyData []byte) bluetooth.ManufacturerDataElement {
	data := make([]byte, 0, 27)
	data = append(data, 0x12, 0x19)        // Offline Finding type and length
	data = append(data, 0x10)              // state
	data = append(data, keyData[6:]...)    // copy last 22 bytes of public key
	data = append(data, (keyData[0] >> 6)) // First two bits
	data = append(data, 0x00)              // Hint (0x00)

	return bluetooth.ManufacturerDataElement{
		CompanyID: 0x004C, // Apple, Inc.
		Data:      data,
	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
