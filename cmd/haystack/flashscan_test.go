package main

import "testing"

func TestScannerLDFlags(t *testing.T) {
	tests := []struct {
		name     string
		names    []string
		keys     [][]string
		onlyMine bool
		want     string
	}{
		{"no devices", nil, nil, false, ""},
		{"one device", []string{"demo"}, [][]string{{"AAA"}}, false, "-X main.MyDevices='demo=AAA'"},
		{"more keys", []string{"demo"}, [][]string{{"AAA", "BBB"}}, false, "-X main.MyDevices='demo=AAA,BBB'"},
		{"two devices", []string{"demo", "other"}, [][]string{{"AAA"}, {"BBB", "CCC"}}, false,
			"-X main.MyDevices='demo=AAA;other=BBB,CCC'"},
		{"only mine", []string{"demo"}, [][]string{{"AAA"}}, true,
			"-X main.MyDevices='demo=AAA' -X main.OnlyMine=true"},
		{"only mine with no devices", nil, nil, true, "-X main.OnlyMine=true"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := scannerLDFlags(tc.names, tc.keys, tc.onlyMine)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
