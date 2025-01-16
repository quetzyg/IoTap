package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/quetzyg/IoTap/command"
	"github.com/quetzyg/IoTap/config"
	"github.com/quetzyg/IoTap/device"
	"github.com/quetzyg/IoTap/ip"
	"github.com/quetzyg/IoTap/meta"
)

func init() {
	log.SetFlags(0)
}

// banner with the CLI version information and ASCII art.
func banner() {
	fmt.Println(`8888888      88888888888`)
	fmt.Println(`  888            888`)
	fmt.Println(`  888            888`)
	fmt.Println(`  888    .d88b.  888   8888b.  88888b.`)
	fmt.Println(`  888   d88""88b 888      "88b 888 "88b`)
	fmt.Println(`  888   888  888 888  .d888888 888  888`)
	fmt.Println(`  888   Y88..88P 888  888  888 888 d88P`)
	fmt.Println(`8888888  "Y88P"  888  "Y888888 88888P"`)
	fmt.Println(`                               888`)
	fmt.Println(`                               888`)
	fmt.Println(`                               888`)

	fmt.Printf("\nVersion %s [%s] (Build time %s)\n\n", meta.Version, meta.Hash, meta.BuildTime)
}

func main() {
	banner()

	flags := command.NewFlags()

	if len(os.Args) < 2 {
		log.Printf("IP or CIDR notation value required (e.g. 192.168.146.123, 192.168.0.0/24)\n\n")

		flags.Usage()

		os.Exit(1)
	}

	// Collect IP addresses for scanning
	ips, err := ip.Resolve(os.Args[1])
	if err != nil {
		log.Printf("Unable to collect IP addresses: %v\n\n", err)

		flags.Usage()

		os.Exit(1)
	}

	cmd, driver, err := flags.Parse(os.Args[2:])
	if err != nil {
		switch {
		// User explicitly passed -h or --help
		case errors.Is(err, flag.ErrHelp):
			os.Exit(0)

		case errors.Is(err, command.ErrFlagConflict):
			log.Printf("%v\n\n", err)
			if cmd != nil {
				cmd.Usage()
			}

		case errors.Is(err, command.ErrInvalid), errors.Is(err, command.ErrNotFound):
			log.Printf("%v\n\n", err)
			flags.Usage()
		}

		os.Exit(1)
	}

	tapper := device.NewTapper(device.GetProbers(driver))

	val, err := config.LoadValues()
	switch {
	case err == nil:
		if val.Credentials != nil {
			log.Printf("Configuration successfully loaded\n\n")

			tapper.SetCredentials(val.Credentials)
		}

	case errors.Is(err, fs.ErrNotExist):
		log.Printf("No configuration to load\n\n")

	default:
		log.Printf("Unable to load configuration: %v\n\n", err)
	}

	if cmd.Name() == command.Config {
		cfg, err := device.LoadConfig(driver, flags.File())
		if err != nil {
			log.Fatalf("Unable to load device config: %v\n\n", err)
		}

		tapper.SetConfig(cfg)
	}

	if cmd.Name() == command.Secure && !flags.SecureOff() {
		auth, err := device.LoadAuthConfig(flags.File())
		if err != nil {
			log.Fatalf("Unable to load auth config: %v\n\n", err)
		}

		tapper.SetAuthConfig(auth)
	}

	if cmd.Name() == command.Deploy {
		dep, err := device.LoadDeployment(driver, flags.File())
		if err != nil {
			log.Fatalf("Unable to load deployment file: %v\n\n", err)
		}

		tapper.SetDeployment(dep)
	}

	var affected = 0

	log.Printf("Scanning %s...\n", os.Args[1])

	devices, err := tapper.Scan(ips)
	if err != nil {
		goto ErrorHandling
	}

	log.Printf("Devices found: %d\n", len(devices))

	if devices.Empty() {
		os.Exit(0)
	}

	switch cmd.Name() {
	case command.Dump:
		err = devices.SortBy(flags.SortField())
		if err != nil {
			log.Fatalf("Unable to sort results: %v\n", err)
		}

		err = device.ExecDump(devices, flags.DumpFormat(), flags.File())

	case command.Config:
		log.Print("Deploying configuration to devices...")

		affected, err = tapper.Execute(device.Configure, devices)

	case command.Secure:
		log.Print("Securing devices...")

		affected, err = tapper.Execute(device.Secure, devices)

	case command.Version:
		log.Print("Verifying device versions...")

		_, err = tapper.Execute(device.Version, devices)
		if err != nil {
			goto ErrorHandling
		}

		var outdated []device.Versioner
		for _, dev := range devices {
			ver := dev.(device.Versioner)
			if ver.Outdated() {
				outdated = append(outdated, ver)
			}
		}

		log.Printf("Outdated devices: %d\n", len(outdated))

		for _, dev := range outdated {
			log.Println(dev.UpdateDetails())
		}

	case command.Update:
		log.Print("Sending firmware update request to devices...")

		affected, err = tapper.Execute(device.Update, devices)

	case command.Deploy:
		log.Print("Deploying script(s) to devices...")

		affected, err = tapper.Execute(device.Deploy, devices)

	case command.Reboot:
		log.Print("Sending reboot request to devices...")

		affected, err = tapper.Execute(device.Reboot, devices)
	}

	if affected > 0 {
		log.Printf("Affected devices: %d\n", affected)
	}

ErrorHandling:
	if err != nil {
		log.Println("Errors found:")
		log.Println(err)

		os.Exit(1)
	}
}
