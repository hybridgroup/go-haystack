package main

import "testing"

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
