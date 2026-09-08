// The tag is a copy of the tag on machine_esp32xx_adc.go in TinyGo, which is
// the ADC driver that this file needs. The m5stamp-c3 has no driver.
//go:build esp32s3 || (esp32c3 && !m5stamp_c3)

package main

import "machine"

// batterySamples is the number of readings that readBatteryMillivolts averages.
// The ESP32 ADC has no sample time control, so it needs more readings.
const batterySamples = 16

// adcReady tells if machine.InitADC ran. InitADC resets the SAR peripheral and
// drives the analog bus, so it must run one time only and before the radio.
var adcReady bool

// readBatteryMillivolts returns the battery voltage in millivolts.
//
// The board has no battery divider, so BatteryPin and BatteryDivider must give
// the pin and the ratio of a divider that the operator added.
func readBatteryMillivolts() (uint16, bool) {
	pin, num, den, ok := batteryConfig()
	if !ok {
		return 0, false
	}

	if pin < batteryFirstPin || pin > batteryLastPin {
		println("BatteryPin is not an ADC1 pin:", BatteryPin)
		return 0, false
	}

	if !adcReady {
		machine.InitADC()
		adcReady = true
	}

	adc := machine.ADC{Pin: machine.Pin(pin)}
	// The driver sets the attenuation itself and does not read the config.
	if err := adc.Configure(machine.ADCConfig{}); err != nil {
		println("cannot configure battery ADC:", err.Error())
		return 0, false
	}

	var total uint32
	for i := 0; i < batterySamples; i++ {
		total += uint32(adc.Get())
	}
	raw := total / batterySamples

	return batteryMillivolts(raw, batteryReference, num, den), true
}
