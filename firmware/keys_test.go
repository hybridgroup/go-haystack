package main

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"
)

// TestAdvertisingKeys checks the list of keys that AdvertisingKey gives.
func TestAdvertisingKeys(t *testing.T) {
	// The build gives this value, so keep it and put it back.
	old := AdvertisingKey
	defer func() { AdvertisingKey = old }()

	tests := []struct {
		key  string
		want []string
	}{
		{"AAA", []string{"AAA"}},
		{"AAA,BBB", []string{"AAA", "BBB"}},
		{" AAA , BBB ", []string{"AAA", "BBB"}},
		{"AAA,BBB,CCC", []string{"AAA", "BBB", "CCC"}},
		// An empty part is not a key.
		{"AAA,,BBB", []string{"AAA", "BBB"}},
		{"AAA,", []string{"AAA"}},
		{"", nil},
		{",", nil},
	}

	for _, tc := range tests {
		AdvertisingKey = tc.key
		got := advertisingKeys()
		if len(got) != len(tc.want) {
			t.Errorf("AdvertisingKey=%q: got %d keys, want %d", tc.key, len(got), len(tc.want))
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("AdvertisingKey=%q: key %d is %q, want %q", tc.key, i, got[i], tc.want[i])
			}
		}
	}
}

// TestKeyData checks that only a 28 byte base64 key is usable.
func TestKeyData(t *testing.T) {
	good := base64.StdEncoding.EncodeToString(make([]byte, 28))
	short := base64.StdEncoding.EncodeToString(make([]byte, 27))

	tests := []struct {
		key    string
		wantOK bool
	}{
		{good, true},
		{short, false},
		{"", false},
		{"not base64!", false},
	}

	for _, tc := range tests {
		val, err := keyData(tc.key)
		if (err == nil) != tc.wantOK {
			t.Errorf("keyData(%q): got error %v, want ok %v", tc.key, err, tc.wantOK)
			continue
		}
		if tc.wantOK && len(val) != 28 {
			t.Errorf("keyData(%q): got %d bytes, want 28", tc.key, len(val))
		}
	}
}

// TestRotationInterval checks the interval that KeyRotation gives.
func TestRotationInterval(t *testing.T) {
	// The command line gives this value, so keep it and put it back.
	old := KeyRotation
	defer func() { KeyRotation = old }()

	tests := []struct {
		rotation string
		want     time.Duration
		wantOK   bool
	}{
		{"5m", 5 * time.Minute, true},
		{"30s", 30 * time.Second, true},
		{"1h30m", 90 * time.Minute, true},
		// No value keeps the first key for ever.
		{"", 0, false},
		{"0s", 0, false},
		{"-5m", 0, false},
		{"junk", 0, false},
		{"5", 0, false},
	}

	for _, tc := range tests {
		KeyRotation = tc.rotation
		got, ok := rotationInterval()
		if ok != tc.wantOK || got != tc.want {
			t.Errorf("KeyRotation=%q: got %v %v, want %v %v", tc.rotation, got, ok, tc.want, tc.wantOK)
		}
	}
}

// TestSleepTime checks the wait to the nearest deadline.
func TestSleepTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		deadlines []time.Time
		want      time.Duration
	}{
		{"no deadline", nil, time.Hour},
		{"zero values only", []time.Time{{}, {}}, time.Hour},
		{"one deadline", []time.Time{now.Add(5 * time.Minute)}, 5 * time.Minute},
		{"nearest of two", []time.Time{now.Add(15 * time.Minute), now.Add(5 * time.Minute)}, 5 * time.Minute},
		{"one in use", []time.Time{{}, now.Add(20 * time.Second)}, 20 * time.Second},
		// A deadline that has passed must not make a wait of zero.
		{"passed", []time.Time{now.Add(-time.Minute)}, time.Millisecond},
		// A deadline that is further away than an hour still waits an hour.
		{"far away", []time.Time{now.Add(3 * time.Hour)}, time.Hour},
	}

	for _, tc := range tests {
		if got := sleepTime(now, tc.deadlines...); got != tc.want {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestRotationText checks the message that the beacon prints at start up.
func TestRotationText(t *testing.T) {
	old := KeyRotation
	defer func() { KeyRotation = old }()

	KeyRotation = "5m"
	if got := rotationText(false); got != "off" {
		t.Errorf("rotationText(false): got %q, want %q", got, "off")
	}
	if got := rotationText(true); !strings.Contains(got, "5m") {
		t.Errorf("rotationText(true): got %q, want %q", got, "5m")
	}
}
