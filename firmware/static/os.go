//go:build !tinygo

package main

import "os"

// AdvertisingKey is the public key of the device. Must be base64 encoded.
var AdvertisingKey = os.Args[1]
