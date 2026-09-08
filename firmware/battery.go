package main

import "github.com/hybridgroup/go-haystack/lib/findmy"

// Thresholds in millivolts for a single cell LiPo battery, which is full at
// 4200 mV and empty at about 3300 mV. A device that uses a different cell needs
// different values here.
const (
	batteryFullMillivolts   = 3900
	batteryMediumMillivolts = 3700
	batteryLowMillivolts    = 3500
)

// batteryMillivolts converts a raw ADC reading to the battery voltage. raw is
// the 16 bit value that the ADC returns, reference is the full scale of the ADC
// in millivolts, and the divider values scale the pin voltage back up to the
// battery voltage.
func batteryMillivolts(raw, reference, numerator, denominator uint32) uint16 {
	atPin := raw * reference / 0x10000
	return uint16(atPin * numerator / denominator)
}

// batteryStatus returns the FindMy status byte for a battery voltage.
func batteryStatus(millivolts uint16) byte {
	switch {
	case millivolts >= batteryFullMillivolts:
		return findmy.StatusBatteryFull
	case millivolts >= batteryMediumMillivolts:
		return findmy.StatusBatteryMedium
	case millivolts >= batteryLowMillivolts:
		return findmy.StatusBatteryLow
	default:
		return findmy.StatusBatteryCritical
	}
}
