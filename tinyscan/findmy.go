package main

import (
	"errors"

	"tinygo.org/x/bluetooth"
)

const (
	// Apple, Inc.
	appleCompanyID = 0x004C

	// Offline Finding type
	findMyPayloadType = 0x12

	// Length of the payload
	findMyPayloadLength = 0x19

	// Hint byte
	findMyHint = 0x00
)

func parseData(address bluetooth.Address, data []byte) (byte, []byte, error) {
	if len(data) < 27 {
		return 0, nil, errors.New("data is too short")
	}

	if data[0] != findMyPayloadType {
		return 0, nil, errors.New("invalid payload type")
	}

	if data[1] != findMyPayloadLength {
		return 0, nil, errors.New("invalid payload length")
	}

	if data[26] != findMyHint {
		return 0, nil, errors.New("invalid hint")
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
