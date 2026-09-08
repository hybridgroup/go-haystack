//go:build nicenano

package main

import "machine"

// The nice!nano v2 has no battery divider. The battery feeds VDDH, and the
// part measures VDDH/5 on an internal ADC channel.
// See https://github.com/zmkfirmware/zmk/blob/main/app/boards/nicekeyboards/nice_nano/nice_nano_nrf52840_zmk_2_0_0.overlay
//
// The nice!nano v1.0.0 used a 2M and 806k divider on P0.04 instead, so a v1
// board reports a wrong value here.
const (
	batteryEnablePin   = machine.NoPin
	batteryEnableLevel = false
	batteryReference   = 3000
	// The internal channel measures a fifth of VDDH.
	batteryDividerNumerator   = 5
	batteryDividerDenominator = 1
)

var batteryPin = machine.ADC_VDDH.Pin
