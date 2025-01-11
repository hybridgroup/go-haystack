# Firmware

Any device supported by the TinyGo Bluetooth package can be used to create a beacon recognized by OpenHaystack.

## How to flash (static key)

```shell
tinygo flash -target nano-rp2040 -ldflags="-X main.AdvertisingKey='SGVsbG8sIFdvcmxkIQ=='" ./static
```

## How to flash (dynamic key derivation)

```shell
tinygo flash -target nano-rp2040 -ldflags="-X main.DerivationKey='SGVsbG8sIFdvcmxkIQ=='" ./derivation
```
