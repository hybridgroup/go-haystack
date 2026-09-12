//go:build tinygo

package main

// AdvertisingKey holds the public keys of the device, separated by commas.
// Must be base64 encoded.
var AdvertisingKey string

// TxPower is the radio transmit power in dBm. An empty value keeps the default
// power of the radio. Set it with -ldflags "-X main.TxPower=-8".
var TxPower string

// DCDC turns the DC/DC regulator on. Set it to "off" with
// -ldflags "-X main.DCDC=off" for a board without the DC/DC inductors.
var DCDC string

// DCDC0 turns the DC/DC converter of the VDDH stage on. It needs a board that
// is powered through VDDH, and it often gains nothing from a battery. Only set
// it to "on" if a measurement shows it helps.
var DCDC0 string

// BatteryPin is the GPIO number that reads a battery divider. An empty value
// stops the reading. Set it with -ldflags "-X main.BatteryPin=2".
var BatteryPin string

// BatteryDivider is the ratio of the battery divider, such as "1510/510". A
// single number is a ratio to 1. An empty value stops the reading.
var BatteryDivider string

// BatteryType names the cell that the beacon uses, such as "cr2032". An empty
// value is a single cell LiPo. Set it with -ldflags "-X main.BatteryType=cr2032".
var BatteryType string

// BatteryThresholds gives the full, medium and low voltages in millivolts, such
// as "2900/2750/2600". It wins over BatteryType. An empty value keeps the
// values of the cell that BatteryType names.
var BatteryThresholds string

// KeyRotation is how long the beacon uses each key, such as "5m". An empty
// value keeps the first key for ever. Set it with
// -ldflags "-X main.KeyRotation=5m".
var KeyRotation string
