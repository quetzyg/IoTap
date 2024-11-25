package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Stowify/IoTune/command"
	"github.com/Stowify/IoTune/device"
	"github.com/Stowify/IoTune/ip"
	"github.com/Stowify/IoTune/meta"
	"github.com/Stowify/IoTune/shellygen1"
	"github.com/Stowify/IoTune/shellygen2"
)

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

	log.Printf("Version %s [%s] (Build time %s)\n\n", meta.Version, meta.Hash, meta.BuildTime)
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
	defer func() {
		err = f.Close()
		if err != nil {
			log.Fatalf("Config close error: %v", err)
		}
	}()

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

	var ec device.Errors
	if errors.As(err, &ec) && !ec.Empty() {
		log.Println("Errors were found during the network scan:")
		log.Println(ec.Error())

		return
	}

	log.Println("Success!")
}

// execList is a helper function that lists the detected devices.
func execList(devices device.Collection) {
	if len(devices) > 0 {
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
		log.Print("Applying configuration to devices...")
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
		log.Print("Updating software on devices...")
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
		log.Print("Uploading script to devices...")
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
		log.Print("Sending reboot signal to devices...")
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
		log.Print("Versioning devices...")
		err := tuner.Execute(device.Version)

		var ec device.Errors
		if errors.As(err, &ec) && !ec.Empty() {
			ec.Print(devices)

			return
		}

		var updatable []device.Versioner
		for _, dev := range devices {
			ver := dev.(device.Versioner)
			if ver.UpdateAvailable() {
				updatable = append(updatable, ver)
			}
		}

		if len(updatable) > 0 {
			log.Printf("Updatable devices found: %d\n", len(updatable))

			for _, dev := range updatable {
				log.Println(dev.UpdateDetails())
			}

			return
		}
	}

	log.Println("All versioned devices are up to date!")
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
	flags := command.NewFlags()

	if len(os.Args) < 2 {
		log.Printf("CIDR notation value required (e.g. 192.168.0.0/24)\n\n")

		flags.Usage()

		os.Exit(1)
	}

	// Collect IP addresses for scanning
	ips, err := ip.Resolve(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to resolve IP addresses: %v", err)
	}

	// Check if a command has been provided
	if len(os.Args) < 3 {
		log.Printf("Command expected\n\n")

		flags.Usage()

		os.Exit(1)
	}

	cmd, driver, err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Printf("%v\n\n", err)

		flags.Usage()

		os.Exit(1)
	}

	probers := resolveProber(driver)
	if len(probers) == 0 {
		log.Fatalf("Unable to resolve an IoT device prober with the %q driver", driver)
	}

	tuner := device.NewTuner(probers)

	// Avoid scanning if the config/script loading fail
	if cmd == command.Config {
		tuner.SetConfig(loadConfig(driver, flags.ConfigFile()))
	}

	if cmd == command.Script {
		tuner.SetScript(loadScript(driver, flags.ScriptFile()))
	}

	execScan(tuner, ips)

	devices := tuner.Devices()

	log.Printf("IoT devices found: %d\n", len(devices))

	switch cmd {
	case command.Dump:
		err = devices.SortBy(flags.SortField())
		if err != nil {
			log.Fatalf("Unable to sort devices: %v\n", err)
		}

		execList(devices)

	case command.Config:
		execConfig(tuner, devices)

	case command.Version:
		execVersion(tuner, devices)

	case command.Update:
		execUpdate(tuner, devices)

	case command.Script:
		execScript(tuner, devices)

	case command.Reboot:
		execReboot(tuner, devices)
	}
}
