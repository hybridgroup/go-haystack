// Firmware to advertise a FindMy compatible device aka AirTag
// see https://github.com/biemster/FindMy for more information.
//
// To build:
// tinygo flash -target nano-rp2040 -ldflags="-X main.AdvertisingKey='SGVsbG8sIFdvcmxkIQ=='" .
// tinygo flash -target xiao-esp32c3 -ldflags="-X main.AdvertisingKey='SGVsbG8sIFdvcmxkIQ=='" .
//
// For a beacon that changes its key, give more keys, separated by commas, and the
// time on each key:
// -ldflags="-X main.AdvertisingKey='KEY1,KEY2' -X main.KeyRotation=5m"
//
// For a device on a battery, see the power settings in README.md.
//
// For Linux:
// go run . SGVsbG8sIFdvcmxkIQ==
// go run . KEY1,KEY2 5m
package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/hybridgroup/go-haystack/lib/findmy"
	"tinygo.org/x/bluetooth"
)

// Interval between advertisements. A real AirTag uses about this value when it
// is away from its owner. A longer interval uses less current, but a phone then
// finds the device less often.
const advertisingInterval = 1285 * time.Millisecond

// How often the firmware reads the battery voltage. A battery changes slowly,
// so a long interval keeps the CPU asleep.
const batteryCheckInterval = 15 * time.Minute

var adapter = bluetooth.DefaultAdapter

func main() {
	// wait for USB serial to be available
	time.Sleep(2 * time.Second)

	keys := advertisingKeys()
	if len(keys) == 0 {
		fail("no advertising key")
	}
	key, err := keyData(keys[0])
	if err != nil {
		fail("failed to get key data: " + err.Error())
	}

	// A beacon with one key keeps it for ever, so the interval has no use.
	rotation, rotating := rotationInterval()
	if len(keys) == 1 {
		rotating = false
	}
	println("keys are", strconv.Itoa(len(keys)), "and rotation is", rotationText(rotating))
	println("key 0 is", keys[0], "(", len(key), "bytes)")

	// This first reading must stay before adapter.Enable below. On an ESP32 it
	// starts the ADC, which must not happen while the radio runs.
	millivolts, hasBattery := readBatteryMillivolts()
	status := byte(findmy.StatusBatteryFull)
	if hasBattery {
		status = batteryStatus(millivolts)
		println("battery is", strconv.Itoa(int(millivolts)), "mV,", findmy.BatteryStatus(status))
	}

	// The payload has no space left. A non-connectable advertisement holds 31
	// bytes, and the 2 byte company ID and the 27 byte payload fill it.
	opts := bluetooth.AdvertisementOptions{
		AdvertisementType: bluetooth.AdvertisingTypeNonConnInd,
		Interval:          bluetooth.NewDuration(advertisingInterval),
		ManufacturerData:  []bluetooth.ManufacturerDataElement{findmy.NewDataWithStatus(key, status)},
	}

	must("enable BLE stack", adapter.Enable())

	// The DC/DC converter lowers the current a lot, but the board must have the
	// inductor. Set DCDC to "off" for a board that does not have it.
	if dcdcAvailable && dcdcEnabled() {
		if err := adapter.EnableDCSupply(bluetooth.DCSupplyMain, true); err != nil {
			println("cannot enable DCDC:", err.Error())
		}
	}

	// The VDDH stage needs a board that is powered through VDDH, and it gains
	// little unless VDDH is much higher than VDD. It is off unless asked for.
	if dcdcAvailable && dcdc0Enabled() {
		if err := adapter.EnableDCSupply(bluetooth.DCSupplyHighVoltage, true); err != nil {
			println("cannot enable DCDC0:", err.Error())
		}
	}

	// Set the address to the first 6 bytes of the public key. This call must
	// stay between Enable and Start, because each one needs it.
	setAddress(key)

	println("configure advertising...")
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(opts))

	// A lower transmit power uses less current but shortens the range, so the
	// radio keeps its default power unless TxPower asks for a level.
	if dbm, ok := txPower(); ok {
		if err := adv.SetTxPower(dbm); err != nil {
			println("cannot set tx power:", err.Error())
		}
	}

	println("start advertising...")
	must("start adv", adv.Start())
	advertising := true

	address, _ := adapter.Address()
	println("FindMy device using", address.MAC.String())

	// A Nordic radio advertises on its own from here, so the CPU only wakes to
	// change the key or to read the battery. A radio on the HCI path keeps a
	// poll loop running.
	var nextBattery, nextRotation time.Time
	if hasBattery {
		nextBattery = time.Now().Add(batteryCheckInterval)
	}
	if rotating {
		nextRotation = time.Now().Add(rotation)
	}

	index := 0
	for {
		time.Sleep(sleepTime(time.Now(), nextBattery, nextRotation))
		now := time.Now()

		if !nextRotation.IsZero() && !now.Before(nextRotation) {
			nextRotation = now.Add(rotation)
			index = (index + 1) % len(keys)

			// A bad key must not stop the beacon, so it keeps the key it has
			// and tries the next one later.
			next, err := keyData(keys[index])
			if err != nil {
				println("bad key", strconv.Itoa(index), err.Error())
			} else {
				key = next
				println("key", strconv.Itoa(index), "is", keys[index])
				advertising = restartAdvertising(adv, &opts, key, status, advertising)
			}
		}

		if !nextBattery.IsZero() && !now.Before(nextBattery) {
			nextBattery = now.Add(batteryCheckInterval)

			millivolts, ok := readBatteryMillivolts()
			if !ok {
				continue
			}
			newStatus := batteryStatus(millivolts)
			if newStatus == status {
				continue
			}
			status = newStatus
			println("battery is", strconv.Itoa(int(millivolts)), "mV,", findmy.BatteryStatus(status))

			advertising = restartAdvertising(adv, &opts, key, status, advertising)
		}
	}
}

// setAddress uses the first 6 bytes of the key as the BLE address, so the
// address changes with the key.
func setAddress(key []byte) {
	adapter.SetRandomAddress(bluetooth.MAC{key[5], key[4], key[3], key[2], key[1], key[0] | 0xC0})
}

// restartAdvertising puts the key and the status in the advertisement, then
// starts it again. It returns the new state of the advertisement.
//
// A failure here must not stop the device being found, so it goes on and always
// tries to advertise again.
func restartAdvertising(adv *bluetooth.Advertisement, opts *bluetooth.AdvertisementOptions, key []byte, status byte, advertising bool) bool {
	// The BLE stack refuses a new set of parameters while it advertises, so
	// stop before the payload changes.
	//
	// Stop must never run if Start failed. On the HCI path Stop waits for the
	// poll loop that Start makes, and it waits for ever if none runs.
	if advertising {
		if err := adv.Stop(); err != nil {
			println("cannot stop adv:", err.Error())
			return true
		}
	}

	setAddress(key)

	opts.ManufacturerData = []bluetooth.ManufacturerDataElement{findmy.NewDataWithStatus(key, status)}
	if err := adv.Configure(*opts); err != nil {
		println("cannot config adv:", err.Error())
	}
	if err := adv.Start(); err != nil {
		println("cannot start adv:", err.Error())
		return false
	}

	return true
}

// rotationText tells if the beacon rotates its keys.
func rotationText(rotating bool) string {
	if !rotating {
		return "off"
	}

	return KeyRotation
}

// dcdcEnabled reports if the DC/DC regulator must be turned on. The regulator
// is on unless DCDC clearly asks for it to be off, because a board without the
// DC/DC inductors browns out if it is on.
func dcdcEnabled() bool {
	switch strings.ToLower(DCDC) {
	case "off", "false", "no", "0":
		return false
	}
	return true
}

// dcdc0Enabled reports if the high voltage DC/DC regulator must be turned on.
// It is off unless DCDC0 clearly asks for it, because it only works on a board
// that feeds VDDH from the battery.
func dcdc0Enabled() bool {
	switch strings.ToLower(DCDC0) {
	case "on", "true", "yes", "1":
		return true
	}
	return false
}

// txPower returns the transmit power in dBm that TxPower asks for. It returns
// false if TxPower is empty, which keeps the default power of the radio.
func txPower() (int8, bool) {
	if TxPower == "" {
		return 0, false
	}
	dbm, err := strconv.ParseInt(TxPower, 10, 8)
	if err != nil {
		println("bad TxPower value:", TxPower)
		return 0, false
	}
	return int8(dbm), true
}

// must calls a function and fails if an error occurs.
func must(action string, err error) {
	if err != nil {
		fail("failed to " + action + ": " + err.Error())
	}
}

// fail prints a message over and over forever.
func fail(msg string) {
	for {
		println(msg)
		time.Sleep(time.Second)
	}
}
