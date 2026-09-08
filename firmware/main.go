// Firmware to advertise a FindMy compatible device aka AirTag
// see https://github.com/biemster/FindMy for more information.
//
// To build:
// tinygo flash -target nano-rp2040 -ldflags="-X main.AdvertisingKey='SGVsbG8sIFdvcmxkIQ=='" .
//
// For a device on a battery, see the power settings in README.md.
//
// For Linux:
// go run . SGVsbG8sIFdvcmxkIQ==
package main

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/hybridgroup/go-haystack/lib/findmy"
	"tinygo.org/x/bluetooth"
)

// Interval between advertisements. A real AirTag uses about this value when it
// is away from its owner. A longer interval uses less current, but a phone then
// finds the device less often.
const advertisingInterval = 1285 * time.Millisecond

var adapter = bluetooth.DefaultAdapter

func main() {
	// wait for USB serial to be available
	time.Sleep(2 * time.Second)

	key, err := getKeyData()
	if err != nil {
		fail("failed to get key data: " + err.Error())
	}
	println("key is", AdvertisingKey, "(", len(key), "bytes)")

	opts := bluetooth.AdvertisementOptions{
		AdvertisementType: bluetooth.AdvertisingTypeNonConnInd,
		Interval:          bluetooth.NewDuration(advertisingInterval),
		ManufacturerData:  []bluetooth.ManufacturerDataElement{findmy.NewData(key)},
	}

	must("enable BLE stack", adapter.Enable())

	// The DC/DC converter lowers the current a lot, but the board must have the
	// inductor. Set DCDC to "off" for a board that does not have it.
	if dcdcEnabled() {
		if err := adapter.EnableDCSupply(bluetooth.DCSupplyMain, true); err != nil {
			println("cannot enable DCDC:", err.Error())
		}
	}

	// The VDDH stage needs a board that is powered through VDDH, and it gains
	// little unless VDDH is much higher than VDD. It is off unless asked for.
	if dcdc0Enabled() {
		if err := adapter.EnableDCSupply(bluetooth.DCSupplyHighVoltage, true); err != nil {
			println("cannot enable DCDC0:", err.Error())
		}
	}

	// Set the address to the first 6 bytes of the public key.
	adapter.SetRandomAddress(bluetooth.MAC{key[5], key[4], key[3], key[2], key[1], key[0] | 0xC0})

	println("configure advertising...")
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(opts))

	// A lower transmit power uses less current but shortens the range, so the
	// radio keeps its default power unless TxPower asks for a level.
	if dbm, ok := txPower(); ok {
		if err := adv.SetTxPower(dbm); err != nil {
			println("cannot set tx power:", err.Error())
		}
	}

	println("start advertising...")
	must("start adv", adv.Start())

	address, _ := adapter.Address()
	println("FindMy device using", address.MAC.String())

	// The BLE stack advertises on its own from here. Park the CPU, because each
	// wake up uses current and does no useful work.
	for {
		time.Sleep(time.Hour)
	}
}

// dcdcEnabled reports if the DC/DC regulator must be turned on. The regulator
// is on unless DCDC clearly asks for it to be off, because a board without the
// DC/DC inductors browns out if it is on.
func dcdcEnabled() bool {
	switch strings.ToLower(DCDC) {
	case "off", "false", "no", "0":
		return false
	}
	return true
}

// dcdc0Enabled reports if the high voltage DC/DC regulator must be turned on.
// It is off unless DCDC0 clearly asks for it, because it only works on a board
// that feeds VDDH from the battery.
func dcdc0Enabled() bool {
	switch strings.ToLower(DCDC0) {
	case "on", "true", "yes", "1":
		return true
	}
	return false
}

// txPower returns the transmit power in dBm that TxPower asks for. It returns
// false if TxPower is empty, which keeps the default power of the radio.
func txPower() (int8, bool) {
	if TxPower == "" {
		return 0, false
	}
	dbm, err := strconv.ParseInt(TxPower, 10, 8)
	if err != nil {
		println("bad TxPower value:", TxPower)
		return 0, false
	}
	return int8(dbm), true
}

// getKeyData returns the public key data from the base64 encoded string.
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

// must calls a function and fails if an error occurs.
func must(action string, err error) {
	if err != nil {
		fail("failed to " + action + ": " + err.Error())
	}
}

// fail prints a message over and over forever.
func fail(msg string) {
	for {
		println(msg)
		time.Sleep(time.Second)
	}
}
