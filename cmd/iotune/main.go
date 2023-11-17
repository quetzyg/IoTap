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

func init() {
	log.SetFlags(0)

	log.Println(`8888888      88888888888`)
	log.Println(`  888            888`)
	log.Println(`  888            888`)
	log.Println(`  888    .d88b.  888  888  888 88888b.   .d88b.`)
	log.Println(`  888   d88""88b 888  888  888 888 "88b d8P  Y8b`)
	log.Println(`  888   888  888 888  888  888 888  888 88888888`)
	log.Println(`  888   Y88..88P 888  Y88b 888 888  888 Y8b.`)
	log.Println(`8888888  "Y88P"  888   "Y88888 888  888  "Y8888`)
	log.Println(``)

	log.Printf("Version %s (Build time %s)", iotune.Version, iotune.BuildTime)

	// Flag setup
	flag.StringVar(&driver, "d", device.Driver, "IoT driver name (default "+device.Driver+")")
	flag.StringVar(&driver, "driver", device.Driver, "IoT driver name (default "+device.Driver+")")

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
			device.Driver,     // default driver
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

// loadConfig is responsible for the configuration loading logic, performing a series of checks,
// including verifying the driver, checking the file path, and handling I/O operations.
func loadConfig(driver string) {
	switch driver {
	case device.Driver:
		log.Fatalln("In order to load a configuration file, a specific driver must be set")
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
func dump(devices device.Collection, separator string) {
	if len(devices) > 0 {
		log.Println("Dumping devices:")

		// Find the appropriate padding for each column
		var widths device.ColWidths
		for _, dev := range devices {
			for i, w := range dev.(device.Tabler).ColWidths() {
				if w > widths[i] {
					widths[i] = w
				}
			}
		}

		format := fmt.Sprintf(
			"%%-%ds%s%%-%ds%s%%-%ds%s%%-%ds%s%%-%ds%s%%-%ds",
			widths[0],
			separator,
			widths[1],
			separator,
			widths[2],
			separator,
			widths[3],
			separator,
			widths[4],
			separator,
			widths[5],
		)

		// Apply the computed format to each IoT device row
		for _, dev := range devices {
			log.Println(dev.(device.Tabler).Row(format))
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

// resolveProber instances according to the driver passed.
func resolveProber(driver string) []device.Prober {
	switch driver {
	case device.Driver:
		return []device.Prober{&shellygen1.Prober{}, &shellygen2.Prober{}}
	case shellygen1.Driver:
		return []device.Prober{&shellygen1.Prober{}}
	case shellygen2.Driver:
		return []device.Prober{&shellygen2.Prober{}}
	default:
		return nil
	}
}

func main() {
	switch mode {
	case modeDump, modeConfig, modeUpdate, modeReboot:
		log.Printf("Executing in %q mode\n", mode)
	default:
		log.Fatalf("Invalid run mode: %s", mode)
	}

	probers := resolveProber(driver)
	if len(probers) == 0 {
		log.Fatalf("Unable to resolve IoT device probers with driver: %s", driver)
	}

	log.Printf("IoT device probers resolved: %d\n", len(probers))

	if mode == modeConfig {
		loadConfig(driver)
	}

	tuner := device.NewTuner(probers, conf)

	scan(tuner)

	devices := tuner.Devices()

	log.Printf("IoT devices found: %d\n", len(devices))

	switch mode {
	case modeDump:
		dump(devices, " ")
	case modeConfig:
		config(tuner, devices)
	case modeUpdate:
		update(tuner, devices)
	case modeReboot:
		reboot(tuner, devices)
	}
}
