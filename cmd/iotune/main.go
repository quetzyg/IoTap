package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	iotune "github.com/Stowify/IoTune"
	"github.com/Stowify/IoTune/device"
	"github.com/Stowify/IoTune/shellygen1"
	"github.com/Stowify/IoTune/shellygen2"
)

const (
	defaultConfigPath = "config.json"
	modeDump          = "dump"
	modeConfig        = "config"
	modeUpdate        = "update"
)

var (
	driver string
	path   string
	prober device.Prober
	conf   device.Config
	mode   string
)

const usage = `Usage:
%s [--driver DRIVER] [--config CONFIG] [--mode MODE]

Options:
-d, --driver DRIVER	Define the IoT device driver. (%s, %s) (default: %s)
-c, --config CONFIG	Define the configuration file. (default: %s)
-m, --mode   MODE	Define the execution mode. (%s, %s, %s) (default: %s)

With no arguments, the tool will use the %s driver in %s mode.
`

func init() {
	log.SetFlags(log.LstdFlags)

	// Flag setup
	flag.StringVar(&driver, "d", shellygen1.Driver, "IoT driver name (default "+shellygen1.Driver+")")
	flag.StringVar(&driver, "driver", shellygen1.Driver, "IoT driver name (default "+shellygen1.Driver+")")

	flag.StringVar(&mode, "m", modeDump, "Run mode (default "+modeDump+")")
	flag.StringVar(&mode, "mode", modeDump, "Run mode (default "+modeDump+")")

	flag.StringVar(&path, "c", defaultConfigPath, "Location of the config file (default "+defaultConfigPath+")")
	flag.StringVar(&path, "config", defaultConfigPath, "Location of the config file (default "+defaultConfigPath+")")

	flag.Usage = func() {
		fmt.Printf(
			usage,
			os.Args[0],
			shellygen1.Driver, // 1st driver
			shellygen2.Driver, // 2nd driver
			shellygen1.Driver, // default driver
			defaultConfigPath,
			modeDump,   // 1st mode
			modeConfig, // 2nd mode
			modeUpdate, // 3rd mode
			modeDump,   // default mode
			shellygen1.Driver,
			modeDump,
		)
	}
	flag.Parse()
}

func main() {
	switch driver {
	case shellygen1.Driver:
		prober = &shellygen1.Prober{}
		conf = &shellygen1.Config{}
	case shellygen2.Driver:
		prober = &shellygen2.Prober{}
		conf = &shellygen2.Config{}
	default:
		log.Fatalf("Unknown driver: %s", driver)
	}

	log.Printf("Loaded driver: %s\n", driver)

	log.Printf("Run mode: %s\n", mode)

	if mode == modeConfig {
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

		if err = device.LoadConfig(f, conf); err != nil {
			log.Fatalf("Config load error: %v", err)
		}
	}

	tuner := device.NewTuner()

	log.Println("Starting IoT device scan...")
	err := tuner.Scan(iotune.Address(), prober)
	log.Println("done!")

	var pe device.ProbeErrors
	if errors.As(err, &pe) && !pe.Empty() {
		log.Println("Errors were found during the scan:")

		for _, e := range pe {
			log.Printf("%v", e)
		}
	}

	devices := tuner.Devices()

	log.Printf("IoT devices found: %d\n", len(devices))

	switch mode {
	case modeDump:
		dump(devices)
	case modeConfig:
		config(tuner, devices)
	case modeUpdate:
		update(tuner, devices)
	}
}

// dump the detected devices.
func dump(devices device.Collection) {
	if len(devices) > 0 {
		log.Println("Dumping devices:")
		for _, dev := range devices {
			log.Println(dev)
		}
	}
}

// config the detected devices.
func config(tuner *device.Tuner, devices device.Collection) {
	if len(devices) > 0 {
		log.Print("Configuring IoT devices...")
		err := tuner.ConfigureDevices(conf)
		log.Println("done!")

		var oe device.OperationErrors
		if errors.As(err, &oe) && !oe.Empty() {
			log.Printf("Successful device configurations: %d\n", len(devices)-len(oe))
			log.Printf("Failed device configurations: %d\n", len(oe))

			for _, e := range oe {
				log.Println(e)
			}

			return
		}

		log.Printf("All (%d) devices, successfully configured!\n", len(devices))
	}
}

// update the firmware of the detected devices.
func update(tuner *device.Tuner, devices device.Collection) {
	if len(devices) > 0 {
		log.Print("Updating IoT devices...")
		err := tuner.UpdateDevices()
		log.Println("done!")

		var oe device.OperationErrors
		if errors.As(err, &oe) && !oe.Empty() {
			log.Printf("Successful device updates: %d\n", len(devices)-len(oe))
			log.Printf("Failed device updates: %d\n", len(oe))

			for _, e := range oe {
				log.Println(e)
			}

			return
		}

		log.Printf("All (%d) devices, successfully updated!\n", len(devices))
	}
}
