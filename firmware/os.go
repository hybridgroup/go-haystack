//go:build !tinygo

package main

import "os"

// AdvertisingKey holds the public keys of the device, separated by commas.
// Must be base64 encoded.
var AdvertisingKey = os.Args[1]

// KeyRotation is how long the beacon uses each key, such as "5m". It comes from
// the second argument, and an empty value keeps the first key for ever.
var KeyRotation = argument(2)

// argument returns the argument at the index, or an empty string if the command
// line does not have it.
func argument(index int) string {
	if len(os.Args) <= index {
		return ""
	}

	return os.Args[index]
}

// TxPower is the radio transmit power in dBm. An operating system does not give
// this control, so it stays empty here.
var TxPower string

// DCDC turns the DC/DC regulator on. An operating system does not give this
// control, so it is off here.
var DCDC = "off"

// DCDC0 turns the high voltage DC/DC regulator on. An operating system does not
// give this control, so it stays empty here.
var DCDC0 string

// BatteryPin is the GPIO number that reads a battery divider. An operating
// system does not give this control, so it stays empty here.
var BatteryPin string

// BatteryDivider is the ratio of the battery divider. An operating system does
// not give this control, so it stays empty here.
var BatteryDivider string
