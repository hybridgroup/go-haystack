package main

import (
	"strconv"
	"strings"
)

// batteryConfig returns the ADC pin and the divider ratio that BatteryPin and
// BatteryDivider give. It returns false unless both values are usable.
func batteryConfig() (uint8, uint32, uint32, bool) {
	if BatteryPin == "" || BatteryDivider == "" {
		return 0, 0, 0, false
	}

	pin, err := strconv.ParseUint(BatteryPin, 10, 8)
	if err != nil {
		println("bad BatteryPin value:", BatteryPin)
		return 0, 0, 0, false
	}

	num, den, ok := parseDivider(BatteryDivider)
	if !ok {
		println("bad BatteryDivider value:", BatteryDivider)
		return 0, 0, 0, false
	}

	return uint8(pin), num, den, true
}

// parseDivider reads a divider ratio such as "1510/510". A single number is a
// ratio to 1, so "2" is the same as "2/1".
func parseDivider(s string) (uint32, uint32, bool) {
	numText, denText := s, "1"
	if i := strings.IndexByte(s, '/'); i >= 0 {
		numText, denText = s[:i], s[i+1:]
	}

	// The values stay in 16 bits so that batteryMillivolts cannot overflow.
	num, err := strconv.ParseUint(numText, 10, 16)
	if err != nil {
		return 0, 0, false
	}
	den, err := strconv.ParseUint(denText, 10, 16)
	if err != nil || den == 0 {
		return 0, 0, false
	}

	// A divider lowers the voltage, so a numerator that is less than the
	// denominator is a pair of values in the wrong order.
	if num < den {
		return 0, 0, false
	}

	return uint32(num), uint32(den), true
}
