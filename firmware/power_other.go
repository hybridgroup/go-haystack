// The tag is a copy of the tag on adapter_dcsupply-other.go in the bluetooth
// package, where EnableDCSupply always gives an error.
//go:build !(softdevice && (s113v7 || s132v6 || s140v6 || s140v7))

package main

// dcdcAvailable tells if the radio has the DC/DC regulator calls.
const dcdcAvailable = false
