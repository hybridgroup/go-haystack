package main

import (
	"testing"
	"time"

	"github.com/hybridgroup/go-haystack/lib/findmy"
)

func TestSleepState(t *testing.T) {
	want := sleepState{status: findmy.StatusBatteryLow, batteryEpoch: 0xBEEF}
	got, ok := unpackSleepState(packSleepState(want))
	if !ok || got != want {
		t.Errorf("got %+v %v, want %+v true", got, ok, want)
	}

	for _, v := range []uint32{0, 0xFFFFFFFF, 0x12345678} {
		if _, ok := unpackSleepState(v); ok {
			t.Errorf("value %#x: got valid, want not valid", v)
		}
	}
}

func TestSleepTimes(t *testing.T) {
	tests := []struct {
		sleep, awake string
		wantSleep    time.Duration
		wantAwake    time.Duration
		wantOK       bool
	}{
		{"", "", 0, 0, false},
		{"", "5s", 0, 0, false},
		{"30s", "", 30 * time.Second, defaultAwakeTime, true},
		{"1m", "3s", time.Minute, 3 * time.Second, true},
		{"30s", "junk", 30 * time.Second, defaultAwakeTime, true},
		{"30s", "-1s", 30 * time.Second, defaultAwakeTime, true},
		{"junk", "2s", 0, 0, false},
		{"0s", "2s", 0, 0, false},
	}

	for _, tc := range tests {
		SleepInterval, AwakeTime = tc.sleep, tc.awake
		sleep, awake, ok := sleepTimes()
		if sleep != tc.wantSleep || awake != tc.wantAwake || ok != tc.wantOK {
			t.Errorf("SleepInterval=%q AwakeTime=%q: got %v %v %v, want %v %v %v",
				tc.sleep, tc.awake, sleep, awake, ok, tc.wantSleep, tc.wantAwake, tc.wantOK)
		}
	}
	SleepInterval, AwakeTime = "", ""
}

func TestKeyIndex(t *testing.T) {
	tests := []struct {
		elapsed, rotation time.Duration
		count, want       int
	}{
		{0, 5 * time.Minute, 3, 0},
		{4 * time.Minute, 5 * time.Minute, 3, 0},
		{5 * time.Minute, 5 * time.Minute, 3, 1},
		{14 * time.Minute, 5 * time.Minute, 3, 2},
		{15 * time.Minute, 5 * time.Minute, 3, 0},
		{time.Hour, 5 * time.Minute, 1, 0},
		{time.Hour, 0, 3, 0},
	}

	for _, tc := range tests {
		if got := keyIndex(tc.elapsed, tc.rotation, tc.count); got != tc.want {
			t.Errorf("keyIndex(%v, %v, %d): got %d, want %d", tc.elapsed, tc.rotation, tc.count, got, tc.want)
		}
	}
}

func TestBatteryEpoch(t *testing.T) {
	if got := batteryEpoch(batteryCheckInterval - time.Second); got != 0 {
		t.Errorf("before the first interval: got %d, want 0", got)
	}
	if got := batteryEpoch(batteryCheckInterval); got != 1 {
		t.Errorf("at the first interval: got %d, want 1", got)
	}
}
