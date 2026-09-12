package main

import "github.com/hybridgroup/go-haystack/lib/findmy"

// batteryLimits holds the thresholds in millivolts that give the status of a
// battery.
type batteryLimits struct {
	full, medium, low uint16
}

// batteryProfiles gives the thresholds for each cell type. The coin cell values
// come from the Energizer CR2032 and CR1220 datasheets, and the alkaline values
// are for two Energizer E91 cells in series.
// https://data.energizer.com/pdfs/cr2032.pdf
// https://data.energizer.com/pdfs/cr1220.pdf
// https://data.energizer.com/pdfs/e91.pdf
var batteryProfiles = map[string]batteryLimits{
	"lipo":        {3900, 3700, 3500},
	"cr2032":      {2900, 2750, 2600},
	"cr1220":      {2950, 2800, 2650},
	"aa-alkaline": {2800, 2500, 2200},
}

// defaultBatteryType is the cell that the firmware uses when BatteryType is
// empty.
const defaultBatteryType = "lipo"

// batteryMillivolts converts a raw ADC reading to the battery voltage. raw is
// the 16 bit value that the ADC returns, reference is the full scale of the ADC
// in millivolts, and the divider values scale the pin voltage back up to the
// battery voltage.
func batteryMillivolts(raw, reference, numerator, denominator uint32) uint16 {
	atPin := raw * reference / 0x10000
	return uint16(atPin * numerator / denominator)
}

// batteryStatus returns the FindMy status byte for a battery voltage.
func batteryStatus(millivolts uint16, limits batteryLimits) byte {
	switch {
	case millivolts >= limits.full:
		return findmy.StatusBatteryFull
	case millivolts >= limits.medium:
		return findmy.StatusBatteryMedium
	case millivolts >= limits.low:
		return findmy.StatusBatteryLow
	default:
		return findmy.StatusBatteryCritical
	}
}
