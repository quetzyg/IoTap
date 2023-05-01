package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Stowify/IoTune/internal/iot"
	"github.com/Stowify/IoTune/internal/iot/device/shellygen1"
	"github.com/Stowify/IoTune/internal/network"
)

const (
	defaultPath = "config.json"
	scanMode    = "scan"
	pushMode    = "push"
)

var (
	driver string
	path   string
	prober iot.Prober
	config iot.Config
	mode   string
)

const usage = `Usage:
%s [--driver DRIVER] [--config CONFIG] [--scan BOOL]

Options:
-d, --driver DRIVER	Define the IoT device driver. (default: %s)
-c, --config CONFIG	Define the configuration file. (default: %s)
-m, --mode   MODE	Define the run mode. (default: %s)

With no arguments, the tool will use the %s driver in scan mode (no config pushes).
`

func init() {
	log.SetFlags(log.LstdFlags)

	// Flag setup
	flag.StringVar(&driver, "d", shellygen1.Driver, "IoT driver name (default "+shellygen1.Driver+")")
	flag.StringVar(&driver, "driver", shellygen1.Driver, "IoT driver name (default "+shellygen1.Driver+")")

	flag.StringVar(&mode, "m", "scan", "Run mode (default scan)")
	flag.StringVar(&mode, "mode", "scan", "Run mode (default scan)")

	flag.StringVar(&path, "c", defaultPath, "Location of the config file (default "+defaultPath+")")
	flag.StringVar(&path, "config", defaultPath, "Location of the config file (default "+defaultPath+")")

	flag.Usage = func() {
		fmt.Printf(usage, os.Args[0], shellygen1.Driver, defaultPath, scanMode, shellygen1.Driver)
	}
	flag.Parse()
}

func main() {
	switch driver {
	case shellygen1.Driver:
		prober = &shellygen1.Prober{}
		config = &shellygen1.Config{}
	default:
		log.Fatalf("Unknown driver: %s", driver)
	}

	log.Printf("Loaded driver: %s\n", driver)

	log.Printf("Run mode: %s\n", mode)

	// Only load the config file if we're in push mode
	if mode == pushMode {
		f, err := os.Open(path)
		defer func(f *os.File) {
			err = f.Close()
			if err != nil {
				log.Fatalf("Config close error: %v", err)
			}
		}(f)

		if err != nil {
			log.Fatalf("Config open error: %s", err)
		}

		if err = iot.LoadConfig(f, config); err != nil {
			log.Fatalf("Config load error: %v", err)
		}
	}

	tuner := iot.NewTuner()

	log.Println("Starting IoT device scan...")
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

	if mode == pushMode {
		log.Print("Pushing configurations to IoT devices...")
		err = tuner.PushToDevices(config)
		log.Println("done!")

		var ce iot.ConfigErrors
		if errors.As(err, &ce) && !ce.Empty() {
			log.Printf("Successful device configurations: %d\n", len(devices)-len(ce))
			log.Printf("Failed device configurations: %d\n", len(ce))

			for _, e := range ce {
				log.Println(e)
			}

			return
		}

		log.Printf("All (%d) devices, successfully configured!\n", len(devices))
	}
}
