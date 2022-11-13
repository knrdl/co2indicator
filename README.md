# CO₂ Indicator

## CO₂ sensor and traffic lights

Measure air quality and represent it with Traffic Light LEDs

## Components

* MH-Z19C Infrared CO₂ Sensor: https://www.berrybase.de/mh-z19c-infrarot-co2-sensor-stecker-kabel
* LED Traffic lights, like https://www.berrybase.de/led-ampel-modul-mit-3x-8mm-led-rot-gelb-gruen-5v
* Raspberry Pi 4
* Female-Female Jumper Wires

## Wiring

See Raspberry Pi GPIO Pin Layout: https://pinout.xyz/

| MH-Z19C (CO₂ Sensor) | Raspberry Pi Pins    |
|----------------------|----------------------|
| GND (Ground)         | 6 (Ground)           |
| Vin                  | 4 (5v Power)         |
| RX                   | 8  (GPIO 14/UART TX) |
| TX                   | 10 (GPIO 15/UART RX) |

Check CO₂ sensor device is registered as serial terminal: `$ ls /dev/serial0`

| Traffic Lights (LED Output) | Raspberry Pi Pins |
|-----------------------------|-------------------|
| GND (Ground)                | 14 (Ground)       |
| Red                         | 15 (GPIO 22)      |
| Yellow                      | 13 (GPIO 27)      |
| Green                       | 16 (GPIO 23)      |

## Deployment

The software can:

* set Traffic Light LEDs
* serve measurements via webserver

### Run binary

Download binary from releases section

```shell
./co2indicator --server :8080 --led-pin.green=23 --led-pin.yellow=27 --led-pin.red=22 --device /dev/serial0
```

Check measurements: `curl localhost:8080`

### Docker

```shell
docker run -it --rm --restart=unless-stopped \
      --device /dev/serial0 --device /sys/class/gpio --device /sys/devices/platform/soc/3f200000.gpio \
      ghcr.io/knrdl/co2indicator:edge \
      --led-pin.green=23 --led-pin.yellow=27 --led-pin.red=22
```

## Dev Setup

### Cross Compile for Raspi

```shell
cd src
GOOS=linux GOARCH=arm GOARM=7  go build
scp app pi@raspberrypi:/home/pi
ssh pi@raspberrypi '/home/pi/app --server :8080 --led-pin.green=23 --led-pin.yellow=27 --led-pin.red=22'
```
