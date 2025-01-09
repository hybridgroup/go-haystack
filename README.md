# go-haystack

![Go Haystack gopher](./images/go-haystack.png)

An easy-to-use and easy-to-setup custom FindMy network using [OpenHaystack](https://github.com/seemoo-lab/openhaystack) using [macless-haystack](https://github.com/dchristl/macless-haystack) but with tools written in Go, no Python required.

## Supported hardware

This package supports using firmware written using [TinyGo](https://tinygo.org/) and the [TinyGo Bluetooth package](https://github.com/tinygo-org/bluetooth). As a result, any of the hardware devices supported should work:

- Adafruit Bluefruit boards using nRF SoftDevice - https://github.com/tinygo-org/bluetooth?tab=readme-ov-file#adafruit-bluefruit-boards
- BBC Microbit using nRF SoftDevice - https://github.com/tinygo-org/bluetooth?tab=readme-ov-file#bbc-microbit
- Other Nordic Semi SoftDevice boards - https://github.com/tinygo-org/bluetooth?tab=readme-ov-file#flashing-the-softdevice-on-other-boards
- Boards using the NINA-FW with an ESP32 co-processor - https://github.com/tinygo-org/bluetooth?tab=readme-ov-file#esp32-nina
- Boards such as the RP2040 Pico-W using the CYW43439 co-processor - https://github.com/tinygo-org/bluetooth?tab=readme-ov-file#cyw43439-rp2040-w

You can also use the same code to run on Linux for your Raspberry Pi or other embedded system with Bluetooth available using BlueZ.

## How to install

### Apple ID

Apple-ID with 2FA enabled. Only sms/text message as second factor is supported!

### anisette-v3-server

1. Start `anisette-v3-server`

```bash
docker network create mh-network
docker run -d --restart always --name anisette -p 6969:6969 --volume anisette-v3_data:/home/Alcoholic/.config/anisette-v3 --network mh-network dadoum/anisette-v3-server
```

2. Start and set up your Macless Haystack endpoint in interactive mode:

```bash
docker run -it --restart unless-stopped --name macless-haystack -p 6176:6176 --volume mh_data:/app/endpoint/data --network mh-network christld/macless-haystack
```

###### You will be asked for your Apple-ID, password and your 2FA. If you see `serving at port 6176 over HTTP` you have all set up correctly

3. Restart macless-haystack server

```bash
docker restart macless-haystack
```

See https://github.com/dchristl/macless-haystack/blob/main/README.md#server-setup for the original instructions.

### install go-haystack

Install the go-haystack command line tool

```shell
go install github.com/hybridgroup/go-haystack
```

## How to use

1. Generate keys for a device

```shell
go-haystack -name=DEVICENAME
```

2. Flash the hardware with the target name of your device and those keys

For example:

```shell
./flash.sh nano-rp2040 DEVICENAME.keys
```

3. Upload the JSON file for that device to your running instance of `macless-haystack` using the web UI.

## Your data may appear

Eventually, if there are iPhone devices in range of your beacons, they will appear in your macless-haystack data.

Have fun, be good!
