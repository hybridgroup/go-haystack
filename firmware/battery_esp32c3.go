//go:build esp32c3 && !m5stamp_c3

package main

const (
	// The ADC uses 11 dB attenuation, which gives a full scale of about 3300
	// mV. See machine_esp32s3_adc.go in TinyGo, which states this scale.
	batteryReference = 3300

	// ADC1 is GPIO0 to GPIO4. GPIO5 is ADC2, which shares its arbiter with
	// the radio and gives noisy readings.
	batteryFirstPin = 0
	batteryLastPin  = 4
)
