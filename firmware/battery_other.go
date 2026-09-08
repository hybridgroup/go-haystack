//go:build !(xiao_ble || nicenano || feather_nrf52840 || esp32s3 || (esp32c3 && !m5stamp_c3))

package main

// readBatteryMillivolts returns the battery voltage in millivolts.
//
// This board has no known battery divider, so it reports nothing and the
// advertisement keeps saying the battery is full.
func readBatteryMillivolts() (uint16, bool) {
	return 0, false
}
