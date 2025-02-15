package findmy

import (
	"errors"

	"tinygo.org/x/bluetooth"
)

const (
	// Apple, Inc.
	AppleCompanyID = 0x004C

	// Not yet registered
	PayloadUnregistered = 0x07

	// Registered for offline finding
	PayloadTypeRegistered = 0x12

	// Length of the payload
	PayloadLength = 0x19

	// Hint byte
	Hint = 0x00

	// Battery full
	StatusBatteryFull = 0x10

	// Battery medium
	StatusBatteryMedium = 0x40

	// Battery low
	StatusBatteryLow = 0x80

	// Battery critical
	StatusBatteryCritical = 0xC0
)

var (
	ErrorNoData               = errors.New("findmy: no data")
	ErrorDataTooShort         = errors.New("findmy: data is too short")
	ErrorUnregistered         = errors.New("findmy: unregistered device")
	ErrorInvalidPayloadType   = errors.New("findmy: invalid payload type")
	ErrorInvalidPayloadLength = errors.New("findmy: invalid payload length")
	ErrorInvalidHint          = errors.New("findmy: invalid hint")
)

// ParseData parses the data from a FindMy device.
// It returns the status byte, the advertising key, and an error if any.
func ParseData(mac bluetooth.MAC, data []byte) (byte, []byte, error) {
	if len(data) == 0 {
		return 0, nil, ErrorNoData
	}

	switch data[0] {
	case PayloadTypeRegistered:
		// registered for offline finding, so go ahead
	case PayloadUnregistered:
		return 0, nil, ErrorUnregistered
	default:
		return 0, nil, ErrorInvalidPayloadType
	}

	if len(data) < 27 {
		return 0, nil, ErrorDataTooShort
	}

	if data[1] != PayloadLength {
		return 0, nil, ErrorInvalidPayloadLength
	}

	if data[26] != Hint {
		return 0, nil, ErrorInvalidHint
	}

	findMyStatus := data[2]
	var key [28]byte
	copy(key[6:], data[3:25])

	// turn address into key bytes
	key[0] = mac[5]
	key[1] = mac[4]
	key[2] = mac[3]
	key[3] = mac[2]
	key[4] = mac[1]
	key[5] = mac[0]

	return findMyStatus, key[:], nil
}

// NewData creates the ManufacturerDataElement for the advertising data used by FindMy devices.
// See https://adamcatley.com/AirTag.html#advertising-data
func NewData(keyData []byte) bluetooth.ManufacturerDataElement {
	data := make([]byte, 0, 27)
	data = append(data, PayloadTypeRegistered, PayloadLength)
	data = append(data, StatusBatteryFull)
	data = append(data, keyData[6:]...)    // copy last 22 bytes of advertising key
	data = append(data, (keyData[0] >> 6)) // first two bits of advertising key
	data = append(data, Hint)

	return bluetooth.ManufacturerDataElement{
		CompanyID: AppleCompanyID,
		Data:      data,
	}
}

// BatteryStatus returns a string representation of the battery status.
func BatteryStatus(status byte) string {
	switch status {
	case StatusBatteryFull:
		return "full"
	case StatusBatteryMedium:
		return "medium"
	case StatusBatteryLow:
		return "low"
	case StatusBatteryCritical:
		return "critical"
	default:
		return "unknown"
	}
}
