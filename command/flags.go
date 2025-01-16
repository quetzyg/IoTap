package command

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/quetzyg/IoTap/device"
	"github.com/quetzyg/IoTap/shellygen1"
	"github.com/quetzyg/IoTap/shellygen2"
)

// StrFlag is a custom flag type representing a string,
// restricted to a predefined set of options.
type StrFlag struct {
	options []string
	value   string
}

// String implements the Stringer interface.
func (sf *StrFlag) String() string {
	return sf.value
}

// Set the flag value after validating the input.
func (sf *StrFlag) Set(value string) error {
	for _, v := range sf.options {
		if value == v {
			sf.value = value

			return nil
		}
	}

	return fmt.Errorf("expected one of: %s", strings.Join(sf.options, ", "))
}

// NewStrFlag creates a new *StrFlag instance.
func NewStrFlag(def string, options ...string) *StrFlag {
	return &StrFlag{
		options: options,
		value:   def,
	}
}

// Available commands
const (
	Dump    = "dump"
	Config  = "config"
	Secure  = "secure"
	Version = "version"
	Update  = "update"
	Deploy  = "deploy"
	Reboot  = "reboot"
)

// Usage strings
const (
	usage = `Usage:
%s <IP|CIDR> <command> [flags]

Commands:
  dump    Output device scan results to STDOUT or to a file
  config  Apply configurations to multiple devices
  secure  Enable/disable device authentication mechanisms
  version Identify devices running outdated software versions
  update  Update firmware on outdated devices
  deploy  Deploy scripts to multiple devices
  reboot  Restart devices

Use %s <IP|CIDR> <command> -h for more information about the command.
`
	commandUsage = `Usage of %s:
 %s <IP|CIDR> %s [flags]

Flags:
`
)

// Flags used in the command.
type Flags struct {
	driver *StrFlag
	file   *string

	dumpCmd       *flag.FlagSet
	dumpSortField *StrFlag
	dumpFormat    *StrFlag

	configCmd *flag.FlagSet

	secureCmd *flag.FlagSet
	secureOff *bool

	versionCmd *flag.FlagSet

	updateCmd *flag.FlagSet

	deployCmd *flag.FlagSet

	rebootCmd *flag.FlagSet
}

// NewFlags creates a new *Flags instance.
func NewFlags() *Flags {
	flags := &Flags{
		driver: NewStrFlag(device.AllDrivers, device.AllDrivers, shellygen1.Driver, shellygen2.Driver),
		file:   new(string),
	}

	// Main usage
	flag.Usage = func() {
		fmt.Printf(
			usage,
			os.Args[0],
			os.Args[0],
		)
	}

	// Dump
	flags.dumpCmd = flag.NewFlagSet(Dump, flag.ContinueOnError)
	flags.dumpCmd.Var(flags.driver, "driver", "Filter by device driver")
	flags.dumpCmd.StringVar(flags.file, "f", "", "Output the scan results to a file")
	flags.dumpSortField = NewStrFlag(
		device.FieldName,
		device.FieldVendor,
		device.FieldIP,
		device.FieldMAC,
		device.FieldName,
		device.FieldModel,
		device.FieldGeneration,
	)
	flags.dumpCmd.Var(flags.dumpSortField, "sort", "Sort devices by field")
	flags.dumpFormat = NewStrFlag(device.FormatCSV, device.FormatCSV, device.FormatJSON)
	flags.dumpCmd.Var(flags.dumpFormat, "format", "Dump output format")
	flags.dumpCmd.Usage = func() {
		fmt.Printf(commandUsage, Dump, os.Args[0], Dump)
		flags.dumpCmd.PrintDefaults()
	}

	// Config
	flags.configCmd = flag.NewFlagSet(Config, flag.ContinueOnError)
	flags.configCmd.Var(flags.driver, "driver", "Filter by device driver")
	flags.configCmd.StringVar(flags.file, "f", "", "Device configuration file path")
	flags.configCmd.Usage = func() {
		fmt.Printf(commandUsage, Config, os.Args[0], Config)
		flags.configCmd.PrintDefaults()
	}

	// Secure
	flags.secureCmd = flag.NewFlagSet(Secure, flag.ContinueOnError)
	flags.secureCmd.Var(flags.driver, "driver", "Filter by device driver")
	flags.secureCmd.StringVar(flags.file, "f", "", "Auth configuration file path (incompatible with --off)")
	flags.secureOff = flags.secureCmd.Bool("off", false, "Turn device authentication off (incompatible with -f)")
	flags.secureCmd.Usage = func() {
		fmt.Printf(commandUsage, Secure, os.Args[0], Secure)
		flags.secureCmd.PrintDefaults()
	}

	// Version
	flags.versionCmd = flag.NewFlagSet(Version, flag.ContinueOnError)
	flags.versionCmd.Var(flags.driver, "driver", "Filter by device driver")
	flags.versionCmd.Usage = func() {
		fmt.Printf(commandUsage, Version, os.Args[0], Version)
		flags.versionCmd.PrintDefaults()
	}

	// Update
	flags.updateCmd = flag.NewFlagSet(Update, flag.ContinueOnError)
	flags.updateCmd.Var(flags.driver, "driver", "Filter by device driver")
	flags.updateCmd.Usage = func() {
		fmt.Printf(commandUsage, Update, os.Args[0], Update)
		flags.updateCmd.PrintDefaults()
	}

	// Deploy
	flags.deployCmd = flag.NewFlagSet(Deploy, flag.ContinueOnError)
	flags.deployCmd.Var(flags.driver, "driver", "Filter by device driver")
	flags.deployCmd.StringVar(flags.file, "f", "", "Device deployment file path")
	flags.deployCmd.Usage = func() {
		fmt.Printf(commandUsage, Deploy, os.Args[0], Deploy)
		flags.deployCmd.PrintDefaults()
	}

	// Reboot
	flags.rebootCmd = flag.NewFlagSet(Reboot, flag.ContinueOnError)
	flags.rebootCmd.Var(flags.driver, "driver", "Filter by device driver")
	flags.rebootCmd.Usage = func() {
		fmt.Printf(commandUsage, Reboot, os.Args[0], Reboot)
		flags.rebootCmd.PrintDefaults()
	}

	return flags
}

// Usage outputs examples to the screen.
func (f *Flags) Usage() {
	flag.Usage()
}

// Driver returns the driver name value.
func (f *Flags) Driver() string {
	return f.driver.String()
}

// File returns the file path value.
func (f *Flags) File() string {
	return *f.file
}

// SortField returns the field by which the dump results should be sorted by.
func (f *Flags) SortField() string {
	return f.dumpSortField.String()
}

// DumpFormat returns the dump data format value.
func (f *Flags) DumpFormat() string {
	return f.dumpFormat.String()
}

// SecureOff returns true if device authentication should be turned off, false otherwise.
func (f *Flags) SecureOff() bool {
	return *f.secureOff
}

// Parse the CLI arguments.
func (f *Flags) Parse(arguments []string) (*flag.FlagSet, string, error) {
	if len(arguments) == 0 {
		return nil, "", ErrNotFound
	}

	err := flag.CommandLine.Parse(arguments)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
	}

	// Handle commands
	switch arguments[0] {
	case Dump:
		err = f.dumpCmd.Parse(arguments[1:])
		if err != nil {
			return f.dumpCmd, "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return f.dumpCmd, f.driver.String(), nil

	case Config:
		err = f.configCmd.Parse(arguments[1:])
		if err != nil {
			return f.configCmd, "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return f.configCmd, f.driver.String(), nil

	case Secure:
		err = f.secureCmd.Parse(arguments[1:])
		if err != nil {
			return f.secureCmd, "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		if f.SecureOff() && f.File() != "" {
			return f.secureCmd, "", fmt.Errorf("%w: '-f' and '--off' flags cannot be used together", ErrFlagConflict)
		}

		return f.secureCmd, f.driver.String(), nil

	case Version:
		err = f.versionCmd.Parse(arguments[1:])
		if err != nil {
			return f.versionCmd, "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return f.versionCmd, f.driver.String(), nil

	case Update:
		err = f.updateCmd.Parse(arguments[1:])
		if err != nil {
			return f.updateCmd, "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return f.updateCmd, f.driver.String(), nil

	case Deploy:
		err = f.deployCmd.Parse(arguments[1:])
		if err != nil {
			return f.deployCmd, "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return f.deployCmd, f.driver.String(), nil

	case Reboot:
		err = f.rebootCmd.Parse(arguments[1:])
		if err != nil {
			return f.rebootCmd, "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return f.rebootCmd, f.driver.String(), nil

	default:
		return nil, "", fmt.Errorf("%w: %s", ErrInvalid, arguments[0])
	}
}
