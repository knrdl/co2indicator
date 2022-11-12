package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type LedPins struct {
	Red    uint
	Yellow uint
	Green  uint
}

type Range struct {
	min, max uint32
}

type LedRanges struct {
	Red    Range
	Yellow Range
	Green  Range
}

func (p *LedPins) Init() error {
	if err := initPin(p.Green); err != nil {
		return err
	}
	if err := initPin(p.Yellow); err != nil {
		return err
	}
	if err := initPin(p.Red); err != nil {
		return err
	}
	return nil
}

func (p *LedPins) Cleanup() error {
	if err := turnPinOff(p.Green); err != nil {
		return err
	}
	if err := turnPinOff(p.Yellow); err != nil {
		return err
	}
	if err := turnPinOff(p.Red); err != nil {
		return err
	}
	return nil
}

func initPin(id uint) error {
	_ = os.WriteFile("/sys/class/gpio/unexport", []byte(strconv.Itoa(int(id))), 0644) // try cleanup
	if err := os.WriteFile("/sys/class/gpio/export", []byte(strconv.Itoa(int(id))), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", id), []byte("out"), 0644); err != nil {
		return err
	}
	return nil
}

func turnPinOff(id uint) error {
	if err := os.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", id), []byte("0"), 0644); err != nil {
		return err
	}
	return nil
}

func turnPinOn(id uint) error {
	if err := os.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", id), []byte("1"), 0644); err != nil {
		return err
	}
	return nil
}

func isPinOn(id uint) (bool, error) {
	if data, err := os.ReadFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", id)); err == nil {
		return data[0] == []byte("1")[0], nil
	} else {
		return false, err
	}
}

func (p *LedPins) update(measurement *Measurement) error {
	co2 := measurement.Co2

	updatePin := func(pinId uint, rng Range) error {
		if co2 >= rng.min && co2 < rng.max {
			if err := turnPinOn(pinId); err != nil {
				return err
			}
		} else {
			if err := turnPinOff(pinId); err != nil {
				return err
			}
		}
		return nil
	}

	blinkPin := func(pinId uint) error {
		if isOn, err := isPinOn(pinId); err != nil {
			return err
		} else {
			if isOn {
				if err := turnPinOff(pinId); err != nil {
					return err
				}
			} else {
				if err := turnPinOn(pinId); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if err := updatePin(p.Green, LedCo2Ranges.Green); err != nil {
		return err
	}
	if err := updatePin(p.Yellow, LedCo2Ranges.Yellow); err != nil {
		return err
	}
	if err := updatePin(p.Red, LedCo2Ranges.Red); err != nil {
		return err
	}
	if co2 > LedCo2Ranges.Red.max { // value too high
		if err := blinkPin(p.Red); err != nil {
			return err
		}
	} else if time.Now().Unix()-measurement.Timestamp > 2*int64(SensorUpdateInterval/time.Second) { // sensor out of order
		_ = turnPinOff(p.Green)
		_ = turnPinOff(p.Red)
		if err := blinkPin(p.Yellow); err != nil {
			return err
		}
	}

	return nil
}
