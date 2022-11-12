package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Measurement struct {
	Co2         uint32 `json:"co2"`
	Temperature uint32 `json:"temperature"`
	Timestamp   int64  `json:"timestamp"`
}

type LatestMeasurement struct {
	value Measurement
	sync.Mutex
}

var latestMeasurement LatestMeasurement

func updateMeasurement(serial *Serial) {
	if value, err := serial.MakeMeasurement(); err == nil {
		latestMeasurement.Lock()
		latestMeasurement.value = value
		latestMeasurement.Unlock()
	} else {
		log.Println(err)
	}
}

func main() {
	deviceName := flag.String("device-name", "/dev/serial0", "serial device name of the sensor")
	serverBind := flag.String("server", "false", "start webserver and bind to address (e.g. :8080)")
	ledPinGreen := flag.Uint("led-pin-green", 0, "GPIO Pin of the green LED")
	ledPinYellow := flag.Uint("led-pin-yellow", 0, "GPIO Pin of the yellow LED")
	ledPinRed := flag.Uint("led-pin-red", 0, "GPIO Pin of the red LED")
	flag.Parse()

	serial, err := openSerialPort(*deviceName)
	if err != nil {
		panic(err)
	}

	pins := LedPins{Red: *ledPinRed, Yellow: *ledPinYellow, Green: *ledPinGreen}
	featureLedsEnabled := pins.Green != 0 && pins.Yellow != 0 && pins.Red != 0
	featureWebserverEnabled := *serverBind != "false"

	if !featureLedsEnabled && !featureWebserverEnabled {
		flag.Usage()
	} else {
		if featureLedsEnabled {
			if err := pins.Init(); err != nil {
				panic(err)
			}
			defer pins.Cleanup()

			go func(pins *LedPins) {
				ticker := time.NewTicker(LedsUpdateInterval)
				for {
					select {
					case <-ticker.C:
						latestMeasurement.Lock()
						err := pins.update(&latestMeasurement.value)
						latestMeasurement.Unlock()
						if err != nil {
							log.Println(err)
						}
					}
				}
			}(&pins)
		} else {
			log.Println("LED Output is disabled")
		}

		go func(serial *Serial) {
			ticker := time.NewTicker(SensorUpdateInterval)
			for {
				select {
				case <-ticker.C:
					updateMeasurement(serial)
				}
			}
		}(serial)
		updateMeasurement(serial)

		if featureWebserverEnabled {
			startWebserver(*serverBind)
		} else {
			log.Println("Webserver is disabled")

			// keep application running
			quitChannel := make(chan os.Signal, 1)
			signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
			<-quitChannel
		}
	}
}
