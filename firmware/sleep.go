package main

import "time"

// defaultAwakeTime is how long the beacon advertises in each wake cycle when
// AwakeTime is empty.
const defaultAwakeTime = 2 * time.Second

// sleepStateMagic marks a saved state as valid. A power on reset leaves a
// value that does not have it.
const sleepStateMagic = 0xA5

// sleepState is the data that stays in an RTC register in deep sleep.
type sleepState struct {
	status       byte
	batteryEpoch uint16
}

// packSleepState puts the state in 32 bits, with the magic value at the top.
func packSleepState(s sleepState) uint32 {
	return sleepStateMagic<<24 | uint32(s.status)<<16 | uint32(s.batteryEpoch)
}

// unpackSleepState reads a value from packSleepState. It returns false if the
// value does not have the magic value.
func unpackSleepState(v uint32) (sleepState, bool) {
	if v>>24 != sleepStateMagic {
		return sleepState{}, false
	}

	return sleepState{status: byte(v >> 16), batteryEpoch: uint16(v)}, true
}

// sleepTimes returns the deep sleep time and the awake time that SleepInterval
// and AwakeTime give. It returns false if the beacon must not sleep.
func sleepTimes() (time.Duration, time.Duration, bool) {
	if SleepInterval == "" {
		return 0, 0, false
	}

	sleep, err := time.ParseDuration(SleepInterval)
	if err != nil || sleep <= 0 {
		println("bad SleepInterval value:", SleepInterval)
		return 0, 0, false
	}

	if AwakeTime == "" {
		return sleep, defaultAwakeTime, true
	}

	awake, err := time.ParseDuration(AwakeTime)
	if err != nil || awake <= 0 {
		println("bad AwakeTime value:", AwakeTime)
		return sleep, defaultAwakeTime, true
	}

	return sleep, awake, true
}

// keyIndex returns the key that the beacon uses at a time since power on. The
// beacon does not count across deep sleep, so the index comes from the time.
func keyIndex(elapsed, rotation time.Duration, count int) int {
	if rotation <= 0 || count <= 1 {
		return 0
	}

	return int(uint64(elapsed/rotation) % uint64(count))
}

// batteryEpoch returns the number of the battery check period at a time since
// power on. The battery is read again when the number changes.
func batteryEpoch(elapsed time.Duration) uint16 {
	return uint16(elapsed / batteryCheckInterval)
}
