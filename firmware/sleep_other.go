//go:build !esp32c3

package main

import "time"

// deepSleepAvailable tells if the chip has the deep sleep calls.
const deepSleepAvailable = false

func wokeFromDeepSleep() bool { return false }

func rtcTime() time.Duration { return 0 }

func loadSleepState() uint32 { return 0 }

func saveSleepState(uint32) {}

func deepSleep(time.Duration) {}
