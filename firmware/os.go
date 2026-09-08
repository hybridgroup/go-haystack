//go:build !tinygo

package main

import "os"

// AdvertisingKey is the public key of the device. Must be base64 encoded.
var AdvertisingKey = os.Args[1]

// TxPower is the radio transmit power in dBm. An operating system does not give
// this control, so it stays empty here.
var TxPower string

// DCDC turns the DC/DC regulator on. An operating system does not give this
// control, so it is off here.
var DCDC = "off"

// DCDC0 turns the high voltage DC/DC regulator on. An operating system does not
// give this control, so it stays empty here.
var DCDC0 string
