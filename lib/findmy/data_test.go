package findmy

import (
	"testing"

	"tinygo.org/x/bluetooth"
)

func TestNewData(t *testing.T) {
	key := []byte{0xce, 0x8b, 0xad, 0x5f, 0x8a, 0x02, 0x71, 0x53, 0x8f, 0xf5, 0xaf, 0xda, 0x87, 0x49, 0x8c, 0xb0, 0x67, 0xe9, 0xa0, 0x20, 0xd6, 0xe4, 0x16, 0x78, 0x01, 0xd5, 0x5d, 0x83}
	data := NewData(key)
	if data.Data[2] != StatusBatteryFull {
		t.Errorf("expected 0x%02x, got 0x%02x", StatusBatteryFull, data.Data[2])
	}
	if data.Data[3] != 0x71 {
		t.Errorf("expected 0x71, got 0x%02x", data.Data[3])
	}
	if data.Data[4] != 0x53 {
		t.Errorf("expected 0x53, got 0x%02x", data.Data[4])
	}
}

func TestParseData(t *testing.T) {
	address := bluetooth.Address{bluetooth.MACAddress{MAC: bluetooth.MAC{0x02, 0x8a, 0x5f, 0xad, 0x8b, 0xce}}}
	startingkey := []byte{0xce, 0x8b, 0xad, 0x5f, 0x8a, 0x02, 0x71, 0x53, 0x8f, 0xf5, 0xaf, 0xda, 0x87, 0x49, 0x8c, 0xb0, 0x67, 0xe9, 0xa0, 0x20, 0xd6, 0xe4, 0x16, 0x78, 0x01, 0xd5, 0x5d, 0x83}
	data := NewData(startingkey)
	status, key, err := ParseData(address, data.Data)
	if err != nil {
		t.Fatal(err)
	}
	if status != StatusBatteryFull {
		t.Errorf("expected 0x%02x, got 0x%02x", StatusBatteryFull, status)
	}
	if !bytesEqual(key, startingkey) {
		t.Errorf("expected %v, got %v", startingkey, key)
	}
}

func TestBatteryStatus(t *testing.T) {
	tests := []struct {
		status byte
		want   string
	}{
		{StatusBatteryFull, "full"},
		{StatusBatteryMedium, "medium"},
		{StatusBatteryLow, "low"},
		{0xff, "unknown"},
	}
	for _, test := range tests {
		got := BatteryStatus(test.status)
		if got != test.want {
			t.Errorf("BatteryStatus(%d) = %q, want %q", test.status, got, test.want)
		}
	}
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, av := range a {
		if av != b[i] {
			return false
		}
	}
	return true
}
