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
	modeList    = "list"
	modeConfig  = "config"
	modeUpdate  = "update"
	modeScript  = "script"
	modeVersion = "version"
	modeReboot  = "reboot"
)

var (
	mode    string
	driver  string
	cfgPath string
	scrPath string
)

const usage = `Usage:
%s [--mode MODE] [--driver DRIVER] [--config CONFIG] [--script SCRIPT]

Options:
-m, --mode   MODE	Define the execution mode. (%s, %s, %s, %s, %s, %s) (default: %s)
-d, --driver DRIVER	Define the IoT device driver. (%s, %s, %s) (default: %s)
-c, --config CONFIG	Define the IoT device configuration file path.
-s, --script SCRIPT	Define the IoT device script file path.

Without arguments, IoTune will execute in %q mode.
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
	flag.StringVar(&mode, "m", modeList, "Execution mode (default "+modeList+")")
	flag.StringVar(&mode, "mode", modeList, "Execution mode (default "+modeList+")")

	flag.StringVar(&driver, "d", device.Driver, "IoT driver name (default "+device.Driver+")")
	flag.StringVar(&driver, "driver", device.Driver, "IoT driver name (default "+device.Driver+")")

	flag.StringVar(&cfgPath, "c", "", "Location of the config file")
	flag.StringVar(&cfgPath, "config", "", "Location of the config file")

	flag.StringVar(&scrPath, "s", "", "Location of the script file")
	flag.StringVar(&scrPath, "script", "", "Location of the script file")

	flag.Usage = func() {
		fmt.Printf(
			usage,
			os.Args[0],
			modeList,          // 1st mode
			modeConfig,        // 2nd mode
			modeVersion,       // 3rd mode
			modeUpdate,        // 4th mode
			modeScript,        // 5th mode
			modeReboot,        // 6th mode
			modeList,          // default mode
			device.Driver,     // 1st driver
			shellygen1.Driver, // 2nd driver
			shellygen2.Driver, // 3rd driver
			device.Driver,     // default driver
			modeList,
		)
	}
	flag.Parse()
}

// loadConfig encapsulates the configuration loading logic, performing a series of checks,
// including verifying the driver, checking the file path, and error handling.
func loadConfig(driver, path string) device.Config {
	var config device.Config

	switch driver {
	case device.Driver:
		log.Fatalf("The config mode isn't supported by the %q driver", driver)
	case shellygen1.Driver:
		config = &shellygen1.Config{}
	case shellygen2.Driver:
		config = &shellygen2.Config{}
	default:
		log.Fatalf("Unknown driver: %s", driver)
	}

	if path == "" {
		log.Fatalln("The configuration file path is empty")
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

	if err = device.LoadConfig(f, config); err != nil {
		log.Fatalf("%s config load error: %v", driver, err)
	}

	log.Printf("Successfully loaded %q configuration from %s\n", driver, path)

	return config
}

// loadScript encapsulates the script loading logic, performing a series of checks,
// including verifying the driver, checking the file path, and error handling.
func loadScript(driver string, path string) *device.IoTScript {
	switch driver {
	case device.Driver, shellygen1.Driver:
		log.Fatalf("The script mode isn't supported by the %q driver", driver)
	case shellygen2.Driver:
		// All good!
	default:
		log.Fatalf("Unknown driver: %s", driver)
	}

	if path == "" {
		log.Fatalln("The script file path is empty")
	}

	script, err := device.LoadScript(path)
	if err != nil {
		log.Fatalf("%s script loading error: %v", driver, err)
	}

	log.Printf("Successfully loaded script for %q from %s\n", driver, path)

	return script
}

// execScan encapsulates the device scanning and error handling.
func execScan(tuner *device.Tuner) {
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

// execList is a helper function that lists the detected devices.
func execList(devices device.Collection, separator string) {
	if len(devices) > 0 {
		log.Println("Listing found devices:")

		// Compute the appropriate padding for each column
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

// execConfig encapsulates the execution of the device.Configure procedure.
func execConfig(tuner *device.Tuner, devices device.Collection) {
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

		log.Println("All devices were successfully configured!")
	}
}

// execUpdate encapsulates the execution of the device.Update procedure.
func execUpdate(tuner *device.Tuner, devices device.Collection) {
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

		log.Println("All devices were successfully updated!")
	}
}

// execUpdate encapsulates the execution of the device.Script procedure.
func execScript(tuner *device.Tuner, devices device.Collection) {
	if len(devices) > 0 {
		log.Print("Uploading script to IoT devices...")
		err := tuner.Execute(device.Script)
		log.Println("done!")

		var e device.Errors
		if errors.As(err, &e) && !e.Empty() {
			log.Printf("Successful script uploads: %d\n", len(devices)-len(e))
			log.Printf("Failed script uploads: %d\n", len(e))

			for _, err = range e {
				log.Println(err)
			}

			return
		}

		log.Println("The script has been successfully uploaded to all devices!")
	}
}

// execReboot encapsulates the execution of the device.Reboot procedure.
func execReboot(tuner *device.Tuner, devices device.Collection) {
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

		log.Println("All devices were successfully rebooted!")
	}
}

// execVersion encapsulates the execution of the device.Version procedure.
func execVersion(tuner *device.Tuner, devices device.Collection) {
	if len(devices) > 0 {
		log.Print("Versioning IoT devices...")
		err := tuner.Execute(device.Version)
		log.Println("done!")

		var e device.Errors
		if errors.As(err, &e) && !e.Empty() {
			log.Printf("Successfully versioned devices: %d\n", len(devices)-len(e))
			log.Printf("Failed versioned devices: %d\n", len(e))

			for _, err = range e {
				log.Println(err)
			}

			return
		}

		log.Println("All devices were successfully versioned!")

		var updatable []device.Versioner
		for _, dev := range devices {
			ver := dev.(device.Versioner)
			if ver.UpdateAvailable() {
				updatable = append(updatable, ver)
			}
		}

		if len(updatable) > 0 {
			log.Printf("A total of %d device(s) can be updated.\n", len(updatable))

			for _, dev := range updatable {
				log.Println(dev.UpdateDetails())
			}

			return
		}
	}

	log.Println("Nothing to update")
}

// resolveProber instances from a driver value.
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
	case modeList, modeConfig, modeVersion, modeUpdate, modeScript, modeReboot:
		log.Printf("Executing in %q mode\n", mode)
	default:
		log.Fatalf("Invalid execution mode: %s", mode)
	}

	probers := resolveProber(driver)
	if len(probers) == 0 {
		log.Fatalf("Unable to resolve an IoT device prober with the %q driver", driver)
	}

	tuner := device.NewTuner(probers)

	if mode == modeConfig {
		tuner.SetConfig(loadConfig(driver, cfgPath))
	}

	if mode == modeScript {
		tuner.SetScript(loadScript(driver, scrPath))
	}

	execScan(tuner)

	devices := tuner.Devices()

	log.Printf("IoT devices found: %d\n", len(devices))

	switch mode {
	case modeList:
		execList(devices, " ")
	case modeConfig:
		execConfig(tuner, devices)
	case modeVersion:
		execVersion(tuner, devices)
	case modeUpdate:
		execUpdate(tuner, devices)
	case modeScript:
		execScript(tuner, devices)
	case modeReboot:
		execReboot(tuner, devices)
	}
}
