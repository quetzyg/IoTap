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

const defaultConfig = "config.json"

var (
	driver string
	config string
	prober iot.Prober
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Flag declaration
	flag.StringVar(&driver, "d", shelly.Driver, "The IoT driver name. Default is "+shelly.Driver)
	flag.StringVar(&config, "c", defaultConfig, "The configuration file path. Default is "+defaultConfig)
	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Printf("%s [-d <driver>] [-c %s]\n", os.Args[0], defaultConfig)
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

	f, err := os.Open(config)
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Fatalf("unable to close file: %v", err)
		}
	}(f)

	if err != nil {
		log.Fatal(err)
	}

	tuner := iot.NewTuner()

	log.Printf("Starting IoT device scan using the %s driver...", driver)
	err = tuner.Scan(network.Address(), prober)
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
