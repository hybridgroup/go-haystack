//go:build xiao_ble

package main

import "machine"

// The XIAO nRF52840 reads the battery through a 1M and 510k divider on P0.31
// (AIN7). P0.14 is open drain and must go low to connect the divider.
// See https://github.com/zmkfirmware/zmk/blob/main/app/boards/seeed/xiao_ble/xiao_ble_zmk.dts
const (
	batteryEnablePin   = machine.P0_14
	batteryEnableLevel = false
	batteryReference   = 3000
	// full-ohms / output-ohms, which is (1M + 510k) / 510k.
	batteryDividerNumerator   = 1510
	batteryDividerDenominator = 510
)

var batteryPin = machine.Pin(machine.P0_31)
