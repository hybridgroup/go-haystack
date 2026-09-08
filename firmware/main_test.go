package main

import (
	"testing"

	"github.com/hybridgroup/go-haystack/lib/findmy"
)

func TestDCDCEnabled(t *testing.T) {
	tests := []struct {
		dcdc string
		want bool
	}{
		{"", true},
		{"on", true},
		{"ON", true},
		{"true", true},
		{"1", true},
		{"off", false},
		{"OFF", false},
		{"Off", false},
		{"false", false},
		{"no", false},
		{"0", false},
	}

	for _, tc := range tests {
		DCDC = tc.dcdc
		if got := dcdcEnabled(); got != tc.want {
			t.Errorf("DCDC=%q: got %v, want %v", tc.dcdc, got, tc.want)
		}
	}
}

func TestTxPower(t *testing.T) {
	tests := []struct {
		txPower string
		want    int8
		wantOK  bool
	}{
		{"", 0, false},
		{"0", 0, true},
		{"-8", -8, true},
		{"8", 8, true},
		{"-40", -40, true},
		{"junk", 0, false},
		{"200", 0, false},
	}

	for _, tc := range tests {
		TxPower = tc.txPower
		got, ok := txPower()
		if ok != tc.wantOK || got != tc.want {
			t.Errorf("TxPower=%q: got %v %v, want %v %v", tc.txPower, got, ok, tc.want, tc.wantOK)
		}
	}
}

func TestBatteryStatus(t *testing.T) {
	tests := []struct {
		millivolts uint16
		want       byte
	}{
		{4200, findmy.StatusBatteryFull},
		{3900, findmy.StatusBatteryFull},
		{3899, findmy.StatusBatteryMedium},
		{3700, findmy.StatusBatteryMedium},
		{3699, findmy.StatusBatteryLow},
		{3500, findmy.StatusBatteryLow},
		{3499, findmy.StatusBatteryCritical},
		{0, findmy.StatusBatteryCritical},
	}

	for _, tc := range tests {
		if got := batteryStatus(tc.millivolts); got != tc.want {
			t.Errorf("%d mV: got %#x, want %#x", tc.millivolts, got, tc.want)
		}
	}
}

// TestBatteryMillivolts checks the scaling for each board that can read the
// battery. raw is what the ADC returns for a 4200 mV battery, which is a full
// single cell LiPo.
func TestBatteryMillivolts(t *testing.T) {
	tests := []struct {
		board                      string
		raw, reference, num, denom uint32
		want                       uint16
	}{
		// XIAO nRF52840: 1M and 510k divider, so the pin sees 1418 mV.
		{"xiao_ble", 30988, 3000, 1510, 510, 4198},
		// Feather nRF52840: two 150k resistors, so the pin sees 2100 mV.
		{"feather_nrf52840", 45875, 3000, 2, 1, 4198},
		// nice!nano v2: the internal channel sees VDDH/5, which is 840 mV.
		{"nicenano", 18350, 3000, 5, 1, 4195},
		// ESP32 with a 2 to 1 divider that the operator added, so the pin sees
		// 2100 mV of a 3300 mV full scale.
		{"esp32", 41705, 3300, 2, 1, 4200},
	}

	for _, tc := range tests {
		got := batteryMillivolts(tc.raw, tc.reference, tc.num, tc.denom)
		if got != tc.want {
			t.Errorf("%s: got %d mV, want %d mV", tc.board, got, tc.want)
		}
		// Every board must read a full battery as full.
		if s := batteryStatus(got); s != findmy.StatusBatteryFull {
			t.Errorf("%s: %d mV gives status %#x, want full", tc.board, got, s)
		}
	}

	// A reading of zero must not report a full battery.
	if got := batteryMillivolts(0, 3000, 1510, 510); got != 0 {
		t.Errorf("raw 0: got %d mV, want 0", got)
	}

	// The largest reading must not overflow.
	if got := batteryMillivolts(0xffff, 3600, 1510, 510); got == 0 {
		t.Error("raw 0xffff overflowed to 0")
	}
}

func TestDCDC0Enabled(t *testing.T) {
	tests := []struct {
		dcdc0 string
		want  bool
	}{
		{"", false},
		{"off", false},
		{"junk", false},
		{"on", true},
		{"ON", true},
		{"On", true},
		{"true", true},
		{"yes", true},
		{"1", true},
	}

	for _, tc := range tests {
		DCDC0 = tc.dcdc0
		if got := dcdc0Enabled(); got != tc.want {
			t.Errorf("DCDC0=%q: got %v, want %v", tc.dcdc0, got, tc.want)
		}
	}
}

// TestBatteryConfig checks the values that BatteryPin and BatteryDivider give.
func TestBatteryConfig(t *testing.T) {
	// The other tests leave these values changed, so keep and put back the
	// values that the build gave.
	oldPin, oldDivider := BatteryPin, BatteryDivider
	defer func() { BatteryPin, BatteryDivider = oldPin, oldDivider }()

	tests := []struct {
		pin, divider string
		wantPin      uint8
		wantNum      uint32
		wantDen      uint32
		wantOK       bool
	}{
		{"2", "2/1", 2, 2, 1, true},
		{"0", "1510/510", 0, 1510, 510, true},
		{"10", "2", 10, 2, 1, true},
		{"1", "1/1", 1, 1, 1, true},
		// No value stops the reading.
		{"", "2/1", 0, 0, 0, false},
		{"2", "", 0, 0, 0, false},
		{"", "", 0, 0, 0, false},
		// A denominator of zero cannot divide.
		{"2", "2/0", 0, 0, 0, false},
		// A divider lowers the voltage, so these two values are in the wrong
		// order.
		{"2", "510/1510", 0, 0, 0, false},
		{"2", "abc", 0, 0, 0, false},
		{"2", "/", 0, 0, 0, false},
		{"2", "2/x", 0, 0, 0, false},
		{"2", "-2", 0, 0, 0, false},
		{"abc", "2/1", 0, 0, 0, false},
		{"-1", "2/1", 0, 0, 0, false},
		{"300", "2/1", 0, 0, 0, false},
	}

	for _, tc := range tests {
		BatteryPin, BatteryDivider = tc.pin, tc.divider
		pin, num, den, ok := batteryConfig()
		if ok != tc.wantOK || pin != tc.wantPin || num != tc.wantNum || den != tc.wantDen {
			t.Errorf("BatteryPin=%q BatteryDivider=%q: got %d %d %d %v, want %d %d %d %v",
				tc.pin, tc.divider, pin, num, den, ok,
				tc.wantPin, tc.wantNum, tc.wantDen, tc.wantOK)
		}
	}
}
