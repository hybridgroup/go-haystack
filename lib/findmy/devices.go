package findmy

import (
	"encoding/base64"
	"errors"
	"strings"
)

var (
	ErrorNoDeviceName = errors.New("findmy: device has no name")
	ErrorNoDeviceKeys = errors.New("findmy: device has no keys")
	ErrorInvalidKey   = errors.New("findmy: advertisement key must be 28 bytes")
)

// Device is a device of the owner, with the advertisement keys that it uses.
type Device struct {
	Name string
	Keys [][]byte
}

// ParseDevices reads a list of devices such as "name=KEY1,KEY2;name2=KEY3".
// The keys are base64 and 28 bytes long.
func ParseDevices(s string) ([]Device, error) {
	var devices []Device
	for _, part := range strings.Split(s, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		name, list, found := strings.Cut(part, "=")
		name = strings.TrimSpace(name)
		if !found || name == "" {
			return nil, ErrorNoDeviceName
		}

		device := Device{Name: name}
		for _, key := range strings.Split(list, ",") {
			key = strings.TrimSpace(key)
			if key == "" {
				continue
			}

			val, err := base64.StdEncoding.DecodeString(key)
			if err != nil {
				return nil, err
			}
			if len(val) != KeyLength {
				return nil, ErrorInvalidKey
			}

			device.Keys = append(device.Keys, val)
		}

		if len(device.Keys) == 0 {
			return nil, ErrorNoDeviceKeys
		}

		devices = append(devices, device)
	}

	return devices, nil
}

// Lookup returns the name of the device that uses the key.
func Lookup(devices []Device, key []byte) (string, bool) {
	for _, device := range devices {
		for _, k := range device.Keys {
			if SameKey(key, k) {
				return device.Name, true
			}
		}
	}

	return "", false
}

// SameKey reports if the key from a scan is the key of a device. The BLE
// address keeps only 6 bits of byte 0, so the other 2 bits are not compared.
func SameKey(a, b []byte) bool {
	if len(a) != KeyLength || len(b) != KeyLength {
		return false
	}

	if a[0]&0x3F != b[0]&0x3F {
		return false
	}

	for i := 1; i < KeyLength; i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
