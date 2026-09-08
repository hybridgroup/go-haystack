//go:build esp32s3

package main

const (
	// The ADC uses 11 dB attenuation, which gives a full scale of about 3300
	// mV. See machine_esp32s3_adc.go in TinyGo, which states this scale.
	batteryReference = 3300

	// ADC1 is GPIO1 to GPIO10. ADC2 shares its arbiter with the radio, and
	// GPIO19 and GPIO20 are the USB data pins.
	batteryFirstPin = 1
	batteryLastPin  = 10
)
