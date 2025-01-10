package findmy

import (
	"errors"

	"tinygo.org/x/bluetooth"
)

const (
	// Apple, Inc.
	AppleCompanyID = 0x004C

	// Offline Finding type
	FindMyPayloadType = 0x12

	// Length of the payload
	FindMyPayloadLength = 0x19

	// Hint byte
	FindMyHint = 0x00
)

var (
	ErrorDataTooShort         = errors.New("data is too short")
	ErrorInvalidPayloadType   = errors.New("invalid payload type")
	ErrorInvalidPayloadLength = errors.New("invalid payload length")
	ErrorInvalidHint          = errors.New("invalid hint")
)

// ParseData parses the data from a FindMy device.
// It returns the status byte, the advertising key, and an error if any.
func ParseData(address bluetooth.Address, data []byte) (byte, []byte, error) {
	if len(data) < 27 {
		return 0, nil, ErrorDataTooShort
	}

	if data[0] != FindMyPayloadType {
		return 0, nil, ErrorInvalidPayloadType
	}

	if data[1] != FindMyPayloadLength {
		return 0, nil, ErrorInvalidPayloadLength
	}

	if data[26] != FindMyHint {
		return 0, nil, ErrorInvalidHint
	}

	findMyStatus := data[2]
	var key [28]byte
	copy(key[6:], data[3:25])

	// turn address into key bytes
	key[0] = address.MAC[5]
	key[1] = address.MAC[4]
	key[2] = address.MAC[3]
	key[3] = address.MAC[2]
	key[4] = address.MAC[1]
	key[5] = address.MAC[0]

	return findMyStatus, key[:], nil
}
