package main

import "time"

const (
	SensorUpdateInterval = 10 * time.Second
	LedsUpdateInterval   = 1 * time.Second
)

var LedCo2Ranges = LedRanges{ // co2 values in ppm
	Green:  Range{min: 400, max: 1000},
	Yellow: Range{min: 800, max: 1400},
	Red:    Range{min: 1200, max: 2000},
}
