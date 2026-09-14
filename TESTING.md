# Testing

This page tells how to test the code on the host and how to build the code for the
microcontrollers. The CI workflow runs the same steps, so a change that passes here also
passes on CI.

## Unit tests

The unit tests run on the host, so they need no hardware:

```shell
go test ./...
```

The tests are in the `firmware`, `lib/findmy` and `cmd/haystack` packages. They cover the
key file, the key rotation, the report data, and the device list.

Also build and vet the code, as CI does:

```shell
go build ./...
go vet ./...
go test ./...
```

## On macOS

The `firmware` package advertises, which the Bluetooth package does not do on
CoreBluetooth. On macOS, use only the other packages:

```shell
go build . ./cmd/... ./lib/...
go vet . ./cmd/... ./lib/...
go test . ./cmd/... ./lib/...
```

## Microcontroller builds

The firmware and TinyScan code only compiles for a microcontroller. To check it, build one
of the targets:

```shell
tinygo build -o /dev/null -target=xiao-ble ./firmware
tinygo build -o /dev/null -stack-size 8kb -target=clue ./tinyscan
```

TinyScan needs `-stack-size 8kb` because the display driver uses more stack than the
default.

The build needs TinyGo 0.42 or later. See
[https://tinygo.org/getting-started/overview/](https://tinygo.org/getting-started/overview/)
for how to install it.

CI builds these targets:

| Directory | Targets |
| --- | --- |
| `./firmware` | `xiao-ble`, `nicenano`, `feather-nrf52840`, `xiao-esp32c3`, `xiao-esp32s3`, `m5stamp-c3`, `pico-w` |
| `./tinyscan` | `badger2040-w`, `clue`, `pybadge`, `pyportal` |

## Testing on the hardware

A build gives no proof that the beacon works, so test a new beacon with the real network:

1. Flash the device. See [How to use](./README.md#how-to-use).
2. Run `haystack scan` on the host, or use TinyScan, and look for the advertisement keys
   of the device. See [Keys in the scanners](./KEYS.md#keys-in-the-scanners).
3. Import `DEVICENAME.json` into macless-haystack and wait for the first report in the
   web UI. The device must be in range of an iPhone, and the first report can take a
   while.

A device that shows in `haystack scan` but not in the web UI usually has a key set that
macless-haystack does not ask for. See
[Changing the key set](./KEYS.md#changing-the-key-set).
