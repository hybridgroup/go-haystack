# Firmware

Any device supported by the TinyGo Bluetooth package can be used to create a beacon recognized by OpenHaystack.

## How to flash

```shell
tinygo flash -target nano-rp2040 -ldflags="-X main.PublicKey='SGVsbG8sIFdvcmxkIQ=='" .

```
