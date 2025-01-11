// Firmware to advertise a FindMy compatible device aka AirTag
// see https://github.com/biemster/FindMy for more information.
//
// To build:
// tinygo flash -target nano-rp2040 -ldflags="-X main.DerivationKey='SGVsbG8sIFdvcmxkIQ=='" .
package main

import (
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"io"
	"machine"
	"math/big"
	"time"

	"github.com/hybridgroup/go-haystack/lib/findmy"
	"golang.org/x/crypto/hkdf"
	"tinygo.org/x/bluetooth"
)

// The standard is to refresh the key pair every 15 minutes.
const rotateThreshold = 15 * time.Minute

var adapter = bluetooth.DefaultAdapter

func main() {
	// wait for USB serial to be available
	time.Sleep(2 * time.Second)

	// Get and increment index. Ideally this would be done at the end of 15
	// minutes. However writing seems to fail once advertising has started.
	derivIndex := getIndex()
	derivIndex++
	must("write new index to file", writeIndex(derivIndex))
	println("derivation secret is", DerivationSecret)
	println("derivation index is", derivIndex)

	priv, pub, err := getKeyData(derivIndex)
	must("get key data", err)
	println("private key is", base64.StdEncoding.EncodeToString(priv))
	println("public key is", base64.StdEncoding.EncodeToString(pub))

	opts := bluetooth.AdvertisementOptions{
		AdvertisementType: bluetooth.AdvertisingTypeNonConnInd,
		Interval:          bluetooth.NewDuration(1285000 * time.Microsecond), // 1285ms
		ManufacturerData:  []bluetooth.ManufacturerDataElement{findmy.NewData(pub)},
	}

	must("enable BLE stack", adapter.Enable())

	// Set the address to the first 6 bytes of the public key.
	adapter.SetRandomAddress(bluetooth.MAC{pub[5], pub[4], pub[3], pub[2], pub[1], pub[0] | 0xC0})

	println("configure advertising...")
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(opts))

	println("start advertising...")
	must("start adv", adv.Start())

	boot := time.Now()
	address, _ := adapter.Address()
	for uptime := 0; ; uptime++ {
		println("FindMy device using", address.MAC.String(), "uptime", uptime)
		time.Sleep(time.Second)
		if time.Since(boot) > rotateThreshold {
			machine.CPUReset()
		}
	}
}

const keySize = 28

var curve = elliptic.P224()

func getKeyData(i uint64) ([]byte, []byte, error) {
	secret, err := base64.StdEncoding.DecodeString(DerivationSecret)
	if err != nil {
		return nil, nil, err
	}

	info := make([]byte, 8)
	binary.LittleEndian.PutUint64(info, i)
	r := hkdf.New(sha256.New, secret, nil, info)
	for {
		priv := make([]byte, keySize)
		if _, err := io.ReadFull(r, priv); err != nil {
			return nil, nil, err
		}

		privInt := new(big.Int).SetBytes(priv)
		n := curve.Params().N
		if privInt.Sign() > 0 && privInt.Cmp(n) < 0 {
			xInt, _ := curve.ScalarBaseMult(priv)
			xBytes := xInt.Bytes()
			x := make([]byte, keySize)
			copy(x[keySize-len(xBytes):], xBytes)
			return priv, x, nil
		}
	}
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
