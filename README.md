# CO₂ Indicator

## CO₂ sensor and traffic lights

### Components

* MH-Z19C Infrared CO2 Sensor: https://www.berrybase.de/mh-z19c-infrarot-co2-sensor-stecker-kabel
* LED Traffic lights: https://www.berrybase.de/led-ampel-modul-mit-3x-8mm-led-rot-gelb-gruen-5v
* Raspberry Pi

### Setup

TODO


### Cross Compile for Raspi

```shell
GOOS=linux GOARCH=arm GOARM=7  go build
scp app pi@raspberrypi:/home/pi
ssh pi@raspberrypi '/home/pi/app --server :8080 --led-pin-green=23 --led-pin-yellow=27 --led-pin-red=22'
```
