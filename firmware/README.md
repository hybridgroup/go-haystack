# Firmware

Any device supported by the TinyGo Bluetooth package can be used to create a beacon recognized by OpenHaystack.

## How to flash

A Nordic or RP2040 board:

```shell
tinygo flash -target nano-rp2040 -ldflags="-X main.AdvertisingKey='SGVsbG8sIFdvcmxkIQ=='" .
```

An ESP32-C3 or ESP32-S3 board, which uses the radio in the chip:

```shell
tinygo flash -target xiao-esp32c3 -ldflags="-X main.AdvertisingKey='SGVsbG8sIFdvcmxkIQ=='" .
```
