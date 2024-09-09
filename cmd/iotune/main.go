package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	iotune "github.com/Stowify/IoTune"
	"github.com/Stowify/IoTune/device"
	"github.com/Stowify/IoTune/ip"
	"github.com/Stowify/IoTune/shellygen1"
	"github.com/Stowify/IoTune/shellygen2"
)

// Tool operating modes
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
%s NETWORK/MASK [--mode MODE] [--driver DRIVER] [--config CONFIG] [--script SCRIPT]

Options:
-m, --mode   MODE	Define the execution mode. (%s, %s, %s, %s, %s, %s) (default: %s)
-d, --driver DRIVER	Define the IoT device driver. (%s, %s, %s) (default: %s)
-c, --config CONFIG	Define the IoT device configuration file path.
-s, --script SCRIPT	Define the IoT device script file path.
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

	log.Printf("Version %s (Build time %s)\n\n", iotune.Version, iotune.BuildTime)

	// Flag setup
	flag.StringVar(&mode, "m", modeList, "Execution mode (default: "+modeList+")")
	flag.StringVar(&mode, "mode", modeList, "Execution mode (default: "+modeList+")")

	flag.StringVar(&driver, "d", device.Driver, "IoT device driver (default: "+device.Driver+")")
	flag.StringVar(&driver, "driver", device.Driver, "IoT device driver (default: "+device.Driver+")")

	flag.StringVar(&cfgPath, "c", "", "Config file path")
	flag.StringVar(&cfgPath, "config", "", "Config file path")

	flag.StringVar(&scrPath, "s", "", "Script file path")
	flag.StringVar(&scrPath, "script", "", "Script file path")

	flag.Usage = func() {
		fmt.Printf(
			usage,
			os.Args[0],

			// Execution modes
			modeList,    // 1st mode
			modeConfig,  // 2nd mode
			modeVersion, // 3rd mode
			modeUpdate,  // 4th mode
			modeScript,  // 5th mode
			modeReboot,  // 6th mode
			modeList,    // default mode

			// Device drivers
			device.Driver,     // 1st driver
			shellygen1.Driver, // 2nd driver
			shellygen2.Driver, // 3rd driver
			device.Driver,     // default driver
		)
	}
}

// loadConfig encapsulates the configuration loading logic, performing a series of checks,
// including verifying the driver, checking the file path, and error handling.
func loadConfig(driver, path string) device.Config {
	var config device.Config

	switch driver {
	case device.Driver:
		log.Fatalf("The %q driver does not support config mode", driver)
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
		log.Fatalf("The %q driver does not support script mode", driver)
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
func execScan(tuner *device.Tuner, ips []net.IP) {
	log.Println("Scanning the network for IoT devices...")
	err := tuner.Scan(ips)
	log.Println("done!")

	var ec device.Errors
	if errors.As(err, &ec) && !ec.Empty() {
		log.Println("Errors were found during the network scan:")
		log.Println(ec.Error())
	}
}

// execList is a helper function that lists the detected devices.
func execList(devices device.Collection) {
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
			"%%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%s",
			widths[0],
			widths[1],
			widths[2],
			widths[3],
			widths[4],
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
		log.Print("Applying configuration to IoT devices...")
		err := tuner.Execute(device.Configure)

		var ec device.Errors
		if errors.As(err, &ec) && !ec.Empty() {
			ec.Print(devices)

			return
		}

		log.Println("Success!")
	}
}

// execUpdate encapsulates the execution of the device.Update procedure.
func execUpdate(tuner *device.Tuner, devices device.Collection) {
	if len(devices) > 0 {
		log.Print("Updating software on IoT devices...")
		err := tuner.Execute(device.Update)

		var ec device.Errors
		if errors.As(err, &ec) && !ec.Empty() {
			ec.Print(devices)

			return
		}

		log.Println("Success!")
	}
}

// execUpdate encapsulates the execution of the device.Script procedure.
func execScript(tuner *device.Tuner, devices device.Collection) {
	if len(devices) > 0 {
		log.Print("Uploading script to IoT devices...")
		err := tuner.Execute(device.Script)

		var ec device.Errors
		if errors.As(err, &ec) && !ec.Empty() {
			ec.Print(devices)

			return
		}

		log.Println("Success!")
	}
}

// execReboot encapsulates the execution of the device.Reboot procedure.
func execReboot(tuner *device.Tuner, devices device.Collection) {
	if len(devices) > 0 {
		log.Print("Sending reboot signal to IoT devices...")
		err := tuner.Execute(device.Reboot)

		var ec device.Errors
		if errors.As(err, &ec) && !ec.Empty() {
			ec.Print(devices)

			return
		}

		log.Println("Success!")
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
	if len(os.Args) < 2 {
		log.Printf("CIDR notation value required (e.g. 192.168.0.0/24)\n\n")

		flag.Usage()

		os.Exit(1)
	}

	// Collect IP addresses for scanning
	ips, err := ip.Resolve(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to resolve IP addresses: %v", err)
	}

	// Parse from the second argument onward
	err = flag.CommandLine.Parse(os.Args[2:])
	if err != nil {
		log.Fatalf("Unable to parse arguments: %v", err)
	}

	switch mode {
	case modeList, modeConfig, modeVersion, modeUpdate, modeScript, modeReboot:
		log.Printf("Executing in %q mode\n", mode)
	default:
		log.Printf("Invalid execution mode: %s\n\n", mode)

		flag.Usage()

		os.Exit(1)
	}

	probers := resolveProber(driver)
	if len(probers) == 0 {
		log.Fatalf("Unable to resolve an IoT device prober with the %q driver", driver)
	}

	tuner := device.NewTuner(probers)

	// Avoid scanning if the config/script loading fail
	if mode == modeConfig {
		tuner.SetConfig(loadConfig(driver, cfgPath))
	}

	if mode == modeScript {
		tuner.SetScript(loadScript(driver, scrPath))
	}

	execScan(tuner, ips)

	devices := tuner.Devices()

	log.Printf("IoT devices found: %d\n", len(devices))

	switch mode {
	case modeList:
		execList(devices)
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
