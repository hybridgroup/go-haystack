// The tag is a copy of the tag on adapter_dcsupply-other.go in the bluetooth
// package, which is the only build that has the real EnableDCSupply call.
//go:build softdevice && (s113v7 || s132v6 || s140v6 || s140v7)

package main

import "device/nrf"

// dcdcAvailable tells if the radio has the DC/DC regulator calls.
const dcdcAvailable = true

// setLowPowerMode cancels the constant latency mode, which can be on when the firmware starts.
// See https://docs.nordicsemi.com/bundle/ps_nrf52840/page/power.html, sub power modes.
func setLowPowerMode() {
	nrf.POWER.TASKS_LOWPWR.Set(1)
}
