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
