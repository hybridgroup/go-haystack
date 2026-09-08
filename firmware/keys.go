package main

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"
)

// advertisingKeys returns the keys that AdvertisingKey holds. A beacon with
// more than one key uses them in turn.
func advertisingKeys() []string {
	var keys []string
	for _, key := range strings.Split(AdvertisingKey, ",") {
		key = strings.TrimSpace(key)
		if key != "" {
			keys = append(keys, key)
		}
	}

	return keys
}

// keyData returns the public key data from the base64 encoded string.
func keyData(key string) ([]byte, error) {
	val, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	if len(val) != 28 {
		return nil, errors.New("public key must be 28 bytes long")
	}

	return val, nil
}

// rotationInterval returns how long the beacon uses each key. It returns false
// if KeyRotation is empty or holds a bad value, which keeps one key for ever.
func rotationInterval() (time.Duration, bool) {
	if KeyRotation == "" {
		return 0, false
	}

	d, err := time.ParseDuration(KeyRotation)
	if err != nil {
		println("bad KeyRotation value:", KeyRotation)
		return 0, false
	}
	if d <= 0 {
		println("KeyRotation must be more than zero:", KeyRotation)
		return 0, false
	}

	return d, true
}

// sleepTime returns the time to the nearest deadline. A deadline with a zero
// value is not in use. It returns an hour if no deadline is in use.
func sleepTime(now time.Time, deadlines ...time.Time) time.Duration {
	wait := time.Hour
	for _, deadline := range deadlines {
		if deadline.IsZero() {
			continue
		}
		if d := deadline.Sub(now); d < wait {
			wait = d
		}
	}

	if wait < time.Millisecond {
		return time.Millisecond
	}

	return wait
}
