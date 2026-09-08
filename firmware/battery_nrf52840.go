//go:build xiao_ble || nicenano || feather_nrf52840

package main

import "machine"

// batterySamples is the number of readings that readBatteryMillivolts averages,
// to lower the effect of noise on a high impedance divider.
const batterySamples = 4

// readBatteryMillivolts returns the battery voltage in millivolts.
//
// The board reads the battery through a resistor divider, so the value at the
// pin is scaled back up by the divider ratio.
func readBatteryMillivolts() (uint16, bool) {
	if batteryEnablePin != machine.NoPin {
		// The divider draws current all the time it is connected, so it is only
		// switched on around the reading. The pin is open drain, so releasing it
		// back to an input disconnects the divider.
		batteryEnablePin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		batteryEnablePin.Set(batteryEnableLevel)
		defer batteryEnablePin.Configure(machine.PinConfig{Mode: machine.PinInput})
	}

	machine.InitADC()
	adc := machine.ADC{Pin: batteryPin}
	adc.Configure(machine.ADCConfig{
		Reference: batteryReference,
		// The divider has a high source impedance, so the ADC needs the longest
		// sample time. See nRF52840 Product Specification v1.11 table 44,
		// Acquisition time.
		SampleTime: 40,
	})

	var total uint32
	for i := 0; i < batterySamples; i++ {
		total += uint32(adc.Get())
	}
	raw := total / batterySamples

	return batteryMillivolts(raw, batteryReference,
		batteryDividerNumerator, batteryDividerDenominator), true
}
