# TinyScan

![tinyscan](../images/tinyscan.gif)

Scanner for local FindMy devices that runs on small microcontrollers that have Bluetooth and also a screen attached.

Looks for any devices nearby that are broadcasting the correct manufacturer data, and displays the MAC address and the public key for that device on the display.

## Supported hardware

The following devices currently work with the Go Haystack TinyScan firmware.

### Pimoroni Badger-2040W

https://shop.pimoroni.com/products/badger-2040-w?variant=40514062188627


```shell
tinygo flash -target badger2040-w -stack-size 8kb .
```

### Adafruit Clue


https://www.adafruit.com/clue


```shell
tinygo flash -target clue -stack-size 8kb .
```


### Adafruit PyBadge with Airlift Featherwing


https://www.adafruit.com/product/4200

https://www.adafruit.com/product/4264


```shell
tinygo flash -target pybadge -stack-size 8kb .
```

### Adafruit Pyportal


https://www.adafruit.com/product/4116


```shell
tinygo flash -target pyportal -stack-size 8kb .
```

## Your own devices

TinyScan can show the name of your own devices. Give it the advertisement keys of
each device at build time, and a beacon that uses one of those keys then shows the
name of the device with a `*` mark instead of the key.

The `haystack` tool reads the keys from the `.keys` files and builds the command
for you. Use the same device names that you used with `haystack keys`:

```shell
haystack flashscan clue blackgopher redgopher
```

Add `-onlymine` to hide every other beacon:

```shell
haystack -onlymine flashscan clue blackgopher redgopher
```

To do it by hand, put the base64 advertisement keys in the `MyDevices` flag. Each
device is a name, an `=` sign and its keys separated by commas. A `;` separates
the devices:

```shell
tinygo flash -target clue -stack-size 8kb -ldflags="-X main.MyDevices='blackgopher=KEY1,KEY2;redgopher=KEY3'" .
```

Add `-X main.OnlyMine=true` to the same `ldflags` value to hide every other beacon.

## Debugging

To show scanning errors on the TinyScan display, use the `ldflags` flag in your flash command like this:

```
tinygo flash -target badger2040-w -stack-size 8kb -ldflags="-X main.showErrors=true" .
```