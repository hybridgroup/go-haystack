//go:build feather_nrf52840

package main

import "machine"

// The Feather nRF52840 reads the battery through two 150k resistors on P0.29,
// which the board labels A6. The divider needs no enable pin.
// See the adc_vbat example in https://github.com/adafruit/Adafruit_nRF52_Arduino
const (
	batteryEnablePin   = machine.NoPin
	batteryEnableLevel = false
	batteryReference   = 3000
	// The two resistors are equal, so the battery is twice the pin.
	batteryDividerNumerator   = 2
	batteryDividerDenominator = 1
)

var batteryPin = machine.Pin(machine.A6)
