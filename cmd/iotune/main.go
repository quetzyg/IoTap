package main

import (
	"errors"
	"log"

	"github.com/Stowify/IoTune/internal/iot"
	"github.com/Stowify/IoTune/internal/iot/device/shelly"
	"github.com/Stowify/IoTune/internal/network"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	tuner := iot.NewTuner()

	log.Println("Starting IoT device scan...")
	err := tuner.Scan(network.Address(), shelly.ProbeRequest)
	log.Println("done!")

	var se iot.ScanErrors
	if errors.As(err, &se) {
		log.Println("Errors were found during the scan:")

		for _, e := range se {
			log.Printf("%v", e)
		}
	}

	devices := tuner.Devices()

	log.Printf("Known IoT devices found: %d\n", len(devices))
}
