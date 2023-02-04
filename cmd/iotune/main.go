package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Stowify/IoTune/internal/iot"
	"github.com/Stowify/IoTune/internal/iot/device/shelly"
	"github.com/Stowify/IoTune/internal/network"
)

var (
	driver string
	prober iot.Prober
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Flag declaration
	flag.StringVar(&driver, "driver", shelly.Driver, "Specify the IoT driver. Default is "+shelly.Driver)
	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Printf("%s [-driver <driver name>] <config file>\n", os.Args[0])
	}
	flag.Parse()
}

func main() {
	switch driver {
	case shelly.Driver:
		prober = &shelly.Prober{}
	default:
		log.Fatalf("unknown driver: %s", driver)
	}

	tuner := iot.NewTuner()

	log.Printf("Starting IoT device scan using the %s driver...", driver)
	err := tuner.Scan(network.Address(), prober)
	log.Println("done!")

	var pe iot.ProbeErrors
	if errors.As(err, &pe) && !pe.Empty() {
		log.Println("Errors were found during the scan:")

		for _, e := range pe {
			log.Printf("%v", e)
		}
	}

	devices := tuner.Devices()

	log.Printf("IoT devices found: %d\n", len(devices))
}
