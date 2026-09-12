package findmy

import (
	"encoding/base64"
	"testing"
)

// testKey is a valid 28 byte advertisement key, base64 encoded.
const testKey = "UUksYlEEXUmIjb5h9naEnv4ZcdocNWATRE3+Xg=="

const otherKey = "TEsWX7q3MVAb8lhfzYCmMSvDRq3BFC8ltWWC4A=="

func TestParseDevices(t *testing.T) {
	first, err := base64.StdEncoding.DecodeString(testKey)
	if err != nil {
		t.Fatal(err)
	}
	second, err := base64.StdEncoding.DecodeString(otherKey)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		val  string
		want []Device
	}{
		{"empty", "", nil},
		{"one device", "demo=" + testKey, []Device{{Name: "demo", Keys: [][]byte{first}}}},
		{"two keys", "demo=" + testKey + "," + otherKey, []Device{{Name: "demo", Keys: [][]byte{first, second}}}},
		{"two devices", "demo=" + testKey + ";other=" + otherKey, []Device{
			{Name: "demo", Keys: [][]byte{first}},
			{Name: "other", Keys: [][]byte{second}},
		}},
		{"spaces", " demo = " + testKey + " ; ", []Device{{Name: "demo", Keys: [][]byte{first}}}},
		{"empty key", "demo=" + testKey + ",,", []Device{{Name: "demo", Keys: [][]byte{first}}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseDevices(tc.val)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("got %d devices, want %d", len(got), len(tc.want))
			}
			for i := range got {
				if got[i].Name != tc.want[i].Name {
					t.Errorf("device %d is %q, want %q", i, got[i].Name, tc.want[i].Name)
				}
				if len(got[i].Keys) != len(tc.want[i].Keys) {
					t.Fatalf("device %d has %d keys, want %d", i, len(got[i].Keys), len(tc.want[i].Keys))
				}
				for j := range got[i].Keys {
					if string(got[i].Keys[j]) != string(tc.want[i].Keys[j]) {
						t.Errorf("device %d key %d is wrong", i, j)
					}
				}
			}
		})
	}
}

func TestParseDevicesErrors(t *testing.T) {
	tests := []struct {
		name string
		val  string
	}{
		{"no name", "=" + testKey},
		{"no separator", testKey},
		{"no keys", "demo="},
		{"bad base64", "demo=????"},
		{"short key", "demo=SGVsbG8="},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ParseDevices(tc.val); err == nil {
				t.Error("expected an error, got none")
			}
		})
	}
}

func TestSameKey(t *testing.T) {
	key := []byte{0x0e, 0x8b, 0xad, 0x5f, 0x8a, 0x02, 0x71, 0x53, 0x8f, 0xf5, 0xaf, 0xda, 0x87, 0x49, 0x8c, 0xb0, 0x67, 0xe9, 0xa0, 0x20, 0xd6, 0xe4, 0x16, 0x78, 0x01, 0xd5, 0x5d, 0x83}

	// A key from a scan holds the bits of the random static address in byte 0.
	scanned := make([]byte, KeyLength)
	copy(scanned, key)
	scanned[0] |= 0xC0

	if !SameKey(scanned, key) {
		t.Error("a key with the address bits must match")
	}
	if !SameKey(key, key) {
		t.Error("a key must match itself")
	}

	other := make([]byte, KeyLength)
	copy(other, key)
	other[27] ^= 0x01
	if SameKey(scanned, other) {
		t.Error("a different key must not match")
	}

	if SameKey(key[:27], key) {
		t.Error("a short key must not match")
	}
}

func TestLookup(t *testing.T) {
	devices, err := ParseDevices("demo=" + testKey + ";other=" + otherKey)
	if err != nil {
		t.Fatal(err)
	}

	scanned := make([]byte, KeyLength)
	copy(scanned, devices[1].Keys[0])
	scanned[0] |= 0xC0

	name, found := Lookup(devices, scanned)
	if !found || name != "other" {
		t.Errorf("got %q %v, want \"other\" true", name, found)
	}

	unknown := make([]byte, KeyLength)
	if _, found := Lookup(devices, unknown); found {
		t.Error("an unknown key must not match")
	}
}
