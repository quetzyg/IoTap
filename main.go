package main

import (
	"errors"
	"flag"
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
		log.Fatalf("The config command is not supported by the %q driver", driver)
	case shellygen1.Driver:
		config = &shellygen1.Config{}
	case shellygen2.Driver:
		config = &shellygen2.Config{}
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
		log.Fatalf("The script command is not supported by the %q driver", driver)
	case shellygen2.Driver:
		// All good!
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
}

// resolveProbers for a given driver.
func resolveProbers(driver string) []device.Prober {
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
		log.Printf("Unable to resolve IP addresses: %v\n\n", err)

		flags.Usage()

		os.Exit(1)
	}

	cmd, driver, err := flags.Parse(os.Args[2:])
	if err != nil {
		switch {
		// User explicitly passed -h or --help
		case errors.Is(err, flag.ErrHelp):
			os.Exit(0)

		case errors.Is(err, command.ErrInvalid), errors.Is(err, command.ErrNotFound):
			log.Printf("%v\n\n", err)
			flags.Usage()
		}

		os.Exit(1)
	}

	probers := resolveProbers(driver)
	if len(probers) == 0 {
		log.Fatalf("Unable to resolve a device prober with the %q driver", driver)
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
			log.Fatalf("Unable to sort results: %v\n", err)
		}

		err = device.ExecDump(devices, flags.DumpFormat(), flags.DumpFile())

	case command.Config:
		log.Print("Applying configuration to devices...")

		err = device.ExecConfig(tuner, devices)

	case command.Version:
		log.Print("Versioning devices...")

		ood, err := device.ExecVersion(tuner, devices)

		var ec device.Errors
		if errors.As(err, &ec) && !ec.Empty() {
			ec.Print(devices)

			return
		}

		if len(ood) > 0 {
			log.Printf("Out of date devices: %d\n", len(ood))

			for _, dev := range ood {
				log.Println(dev.UpdateDetails())
			}

			return
		}

	case command.Update:
		log.Print("Sending firmware update request to devices...")

		err = device.ExecUpdate(tuner, devices)

	case command.Script:
		log.Print("Uploading script to devices...")

		err = device.ExecScript(tuner, devices)

	case command.Reboot:
		log.Print("Sending reboot request to devices...")

		err = device.ExecReboot(tuner, devices)
	}

	// Command error handling
	var ec device.Errors
	if errors.As(err, &ec) && !ec.Empty() {
		ec.Print(devices)

		os.Exit(1)
	}

	if err != nil {
		log.Print(err)

		os.Exit(1)
	}
}
