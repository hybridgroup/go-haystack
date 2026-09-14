# Keys

A Find My beacon is an advertisement key that goes out over Bluetooth. The key is the
name of the device, the proof that a report belongs to you, and the source of the
Bluetooth address. This page tells how the keys are made, where they are kept, and how
the beacon uses them.

## What a key is

Each key is an ECDSA key pair on the P-224 curve, which gives 28 bytes for each half.
A key has three forms, and all three are base64:

| Form | Source | Use |
| --- | --- | --- |
| Private key | the `D` value of the pair | decrypts the reports, goes in the JSON file |
| Advertisement key | the `X` value of the public point | goes out in the advertisement, and is built into the firmware |
| Hashed adv key | SHA-256 of the advertisement key | the name that macless-haystack asks the endpoint for |

An iPhone that hears the advertisement sends a report to Apple under the hashed
advertisement key. macless-haystack asks for that name and then decrypts the report with
the private key. Only you have the private key, so only you can read the location.

The hash goes into a URL path, so a hash that holds a `/` stops the reports. `haystack
keys` makes a new key when this occurs. About one half of all keys have such a hash, so
the tool tries again up to 100 times for each key.

## Generating keys

```shell
haystack keys DEVICENAME
```

This makes a set of 12 keys and writes two files. Use `-keys` for a different number and
`-v` to print the keys as they are made:

```shell
haystack -v -keys=24 keys DEVICENAME
```

The command writes the files each time, so a second run for the same name replaces the
set. The old keys are then lost, and the beacon and macless-haystack both need the new
set. See [Changing the key set](#changing-the-key-set).

## The files

`DEVICENAME.keys` holds one block for each key, in the order that the beacon uses them:

```
Private key: BASE64
Advertisement key: BASE64
Hashed adv key: BASE64
```

`haystack flash` reads only the `Advertisement key` lines from this file.

`DEVICENAME.json` is the accessory file for macless-haystack. The first private key is
`privateKey` and the others are `additionalKeys`. macless-haystack fetches reports for
all of them, so the web UI shows one device with one history.

Both file types are in `.gitignore`. Keep them private, because the private key reads the
location history of the device, and a copy of the advertisement keys lets somebody else
follow the beacon.

## Rotating keys

A beacon that always sends the same key can be followed by anybody who scans for it. To
stop this, a device gets a set of keys, and the beacon uses them one after the other. The
key and the Bluetooth address both change together, because the address is the first 6
bytes of the key.

`haystack keys` makes 12 keys, and the beacon uses each key for 5 minutes. The set lasts
one hour and then starts again.

Use `-keys` for the size of the set and `-rotate` for the time on each key:

```shell
haystack -keys=24 keys DEVICENAME
haystack -rotate=15m flash DEVICENAME xiao-ble
```

More keys give a longer time before the set repeats, but macless-haystack then asks the
endpoint for more keys at each refresh. A set of 24 keys with `-rotate=1h` covers a day.

`-keys=1` gives the behavior of the older versions, which is one key for ever. An
`-rotate` of `0s`, or an empty value, also keeps the first key for ever.

The rotation needs both parts, so a device that already has a key file with one key keeps
that one key until you make a new set.

## How the beacon uses the keys

`haystack flash` puts the advertisement keys into the firmware with the TinyGo `-ldflags`
option. The keys go in `main.AdvertisingKey` as one list separated by commas, and the
time on each key goes in `main.KeyRotation`:

```shell
tinygo flash -target xiao-ble -ldflags="-X main.AdvertisingKey='KEY1,KEY2' -X main.KeyRotation=5m" ./firmware
```

The advertisement holds only a part of the key, because a non-connectable advertisement
has no more than 31 bytes. The Bluetooth address is the first 6 bytes of the key, with
the two top bits of the first byte set, which makes a static random address. The Apple
manufacturer data then holds the last 22 bytes and the two top bits of the first byte.
A scanner puts the address and the data together to get the full 28 bytes.

The address is a part of the key, so it changes with each key. A scanner therefore cannot
follow the beacon by its address.

## Keys in the scanners

`haystack scan` shows the address and the advertisement key of every beacon in range.

TinyScan does the same on a microcontroller with a display, and it can show the name of
your own devices instead of the key. It needs the keys at build time, because it has no
file system. `haystack flashscan` reads the `.keys` files and builds the command:

```shell
haystack flashscan clue DEVICENAME OTHERDEVICE
```

Add `-onlymine` to hide every beacon that is not one of the given devices. See
[Your own devices](./tinyscan/README.md#your-own-devices).

## Changing the key set

A new set needs both parts to be brought up to date:

1. Flash the device again, so that the firmware has the new advertisement keys.
2. Import `DEVICENAME.json` into macless-haystack again, so that it asks for the new
   hashed advertisement keys.

A device that has only one of the two goes quiet in the web UI, because the beacon then
sends keys that macless-haystack does not ask for.
