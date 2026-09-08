//go:build tinygo

package main

// AdvertisingKey is the public key of the device. Must be base64 encoded.
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
