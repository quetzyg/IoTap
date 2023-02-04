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

const defaultPath = "config.json"

var (
	driver string
	path   string
	prober iot.Prober
	config iot.Config
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Flag setup
	flag.StringVar(&driver, "drv", shelly.Driver, "IoT driver name (default "+shelly.Driver+")")
	flag.StringVar(&path, "cfg", defaultPath, "Location of the config file (default "+defaultPath+")")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS]\n\n", os.Args[0])
		fmt.Printf("%s [-drv <driver>] [-cfg %s]\n", os.Args[0], defaultPath)
	}
	flag.Parse()
}

func main() {
	switch driver {
	case shelly.Driver:
		prober = &shelly.Prober{}
		config = &shelly.Config{}
	default:
		log.Fatalf("unknown driver: %s", driver)
	}

	f, err := os.Open(path)
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Fatalf("unable to close file: %v", err)
		}
	}(f)

	if err != nil {
		log.Fatal(err)
	}

	if err = iot.LoadConfig(f, config); err != nil {
		log.Fatalf("config load error: %v", err)
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
