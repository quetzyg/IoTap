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
	modeReboot        = "reboot"
)

var (
	driver string
	path   string
	conf   device.Config
	mode   string
)

const usage = `Usage:
%s [--driver DRIVER] [--config CONFIG] [--mode MODE]

Options:
-d, --driver DRIVER	Define the IoT device driver. (%s, %s) (default: %s)
-c, --config CONFIG	Define the configuration file path. (default: %s)
-m, --mode   MODE	Define the execution mode. (%s, %s, %s, %s) (default: %s)

Without arguments, the tool will run in %s mode.
`

// loadConfig is responsible for the configuration loading logic, performing a series of checks,
// including verifying the driver, checking the file path, and handling I/O operations.
func loadConfig() {
	switch driver {
	case shellygen1.Driver:
		conf = &shellygen1.Config{}
	case shellygen2.Driver:
		conf = &shellygen2.Config{}
	default:
		log.Fatalf("Unknown driver: %s", driver)
	}

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
		log.Fatalf("%s config load error: %v", driver, err)
	}

	log.Printf("Successfully loaded %q configuration from %s\n", driver, path)
}

// scan for devices on the network.
func scan(tuner *device.Tuner) {
	log.Println("Starting IoT device scan...")
	err := tuner.Scan(iotune.Address())
	log.Println("done!")

	var ec device.Errors
	if errors.As(err, &ec) && !ec.Empty() {
		log.Println("Errors were found during the scan:")

		for _, e := range ec {
			log.Printf("%v", e)
		}
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
		err := tuner.Execute(device.Configure)
		log.Println("done!")

		var e device.Errors
		if errors.As(err, &e) && !e.Empty() {
			log.Printf("Successful device configurations: %d\n", len(devices)-len(e))
			log.Printf("Failed device configurations: %d\n", len(e))

			for _, err = range e {
				log.Println(err)
			}

			return
		}

		log.Printf("All (%d) devices, successfully configured!\n", len(devices))
	}
}

// update the firmware of the available devices.
func update(tuner *device.Tuner, devices device.Collection) {
	if len(devices) > 0 {
		log.Print("Updating IoT devices...")
		err := tuner.Execute(device.Update)
		log.Println("done!")

		var e device.Errors
		if errors.As(err, &e) && !e.Empty() {
			log.Printf("Successful device updates: %d\n", len(devices)-len(e))
			log.Printf("Failed device updates: %d\n", len(e))

			for _, err = range e {
				log.Println(err)
			}

			return
		}

		log.Printf("All (%d) devices, successfully updated!\n", len(devices))
	}
}

// reboot the detected devices.
func reboot(tuner *device.Tuner, devices device.Collection) {
	if len(devices) > 0 {
		log.Print("Rebooting IoT devices...")
		err := tuner.Execute(device.Reboot)
		log.Println("done!")

		var e device.Errors
		if errors.As(err, &e) && !e.Empty() {
			log.Printf("Successful device reboots: %d\n", len(devices)-len(e))
			log.Printf("Failed device reboots: %d\n", len(e))

			for _, err = range e {
				log.Println(err)
			}

			return
		}

		log.Printf("All (%d) devices, successfully rebooted!\n", len(devices))
	}
}

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
			modeReboot, // 4th mode
			modeDump,   // default mode
			modeDump,
		)
	}
	flag.Parse()
}

func main() {
	log.Printf("Running in %q mode\n", mode)

	if mode == modeConfig {
		loadConfig()
	}

	tuner := device.NewTuner([]device.Prober{
		&shellygen1.Prober{},
		&shellygen2.Prober{},
	}, conf)

	scan(tuner)

	devices := tuner.Devices()

	log.Printf("IoT devices found: %d\n", len(devices))

	switch mode {
	case modeDump:
		dump(devices)
	case modeConfig:
		config(tuner, devices)
	case modeUpdate:
		update(tuner, devices)
	case modeReboot:
		reboot(tuner, devices)
	}
}
