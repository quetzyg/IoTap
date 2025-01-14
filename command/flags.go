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
func (f *StrFlag) String() string {
	return f.value
}

// Set the flag value after validating the input.
func (f *StrFlag) Set(value string) error {
	for _, v := range f.options {
		if value == v {
			f.value = value

			return nil
		}
	}

	return fmt.Errorf("expected one of: %s", strings.Join(f.options, ", "))
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
	dumpCmd       *flag.FlagSet
	dumpDriver    *StrFlag
	dumpSortField *StrFlag
	dumpFormat    *StrFlag
	dumpFile      *string

	configCmd    *flag.FlagSet
	configDriver *StrFlag
	configFile   *string

	secureCmd    *flag.FlagSet
	secureDriver *StrFlag
	secureFile   *string
	secureOff    *bool

	versionCmd    *flag.FlagSet
	versionDriver *StrFlag

	updateCmd    *flag.FlagSet
	updateDriver *StrFlag

	deployCmd    *flag.FlagSet
	deployDriver *StrFlag
	deployFile   *string

	rebootCmd    *flag.FlagSet
	rebootDriver *StrFlag
}

// setDriverFlag to a flag set. This keeps code tidy, avoiding boilerplate.
func setDriverFlag(flagSet *flag.FlagSet) *StrFlag {
	driver := NewStrFlag(device.AllDrivers, device.AllDrivers, shellygen1.Driver, shellygen2.Driver)
	flagSet.Var(driver, "driver", "Filter by device driver")

	return driver
}

// NewFlags creates a new *Flags instance.
func NewFlags() *Flags {
	flags := &Flags{}

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
	flags.dumpDriver = setDriverFlag(flags.dumpCmd)
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
	flags.dumpFile = flags.dumpCmd.String("f", "", "Output the scan results to a file")
	flags.dumpCmd.Usage = func() {
		fmt.Printf(commandUsage, Dump, os.Args[0], Dump)
		flags.dumpCmd.PrintDefaults()
	}

	// Config
	flags.configCmd = flag.NewFlagSet(Config, flag.ContinueOnError)
	flags.configDriver = setDriverFlag(flags.configCmd)
	flags.configFile = flags.configCmd.String("f", "", "Device configuration file path")
	flags.configCmd.Usage = func() {
		fmt.Printf(commandUsage, Config, os.Args[0], Config)
		flags.configCmd.PrintDefaults()
	}

	// Secure
	flags.secureCmd = flag.NewFlagSet(Secure, flag.ContinueOnError)
	flags.secureDriver = setDriverFlag(flags.secureCmd)
	flags.secureFile = flags.secureCmd.String("f", "", "Auth configuration file path (incompatible with --off)")
	flags.secureOff = flags.secureCmd.Bool("off", false, "Turn device authentication off (incompatible with -f)")
	flags.secureCmd.Usage = func() {
		fmt.Printf(commandUsage, Secure, os.Args[0], Secure)
		flags.secureCmd.PrintDefaults()
	}

	// Version
	flags.versionCmd = flag.NewFlagSet(Version, flag.ContinueOnError)
	flags.versionDriver = setDriverFlag(flags.versionCmd)
	flags.versionCmd.Usage = func() {
		fmt.Printf(commandUsage, Version, os.Args[0], Version)
		flags.versionCmd.PrintDefaults()
	}

	// Update
	flags.updateCmd = flag.NewFlagSet(Update, flag.ContinueOnError)
	flags.updateDriver = setDriverFlag(flags.updateCmd)
	flags.updateCmd.Usage = func() {
		fmt.Printf(commandUsage, Update, os.Args[0], Update)
		flags.updateCmd.PrintDefaults()
	}

	// Deploy
	flags.deployCmd = flag.NewFlagSet(Deploy, flag.ContinueOnError)
	flags.deployDriver = setDriverFlag(flags.deployCmd)
	flags.deployFile = flags.deployCmd.String("f", "", "Device deployment file path")
	flags.deployCmd.Usage = func() {
		fmt.Printf(commandUsage, Deploy, os.Args[0], Deploy)
		flags.deployCmd.PrintDefaults()
	}

	// Reboot
	flags.rebootCmd = flag.NewFlagSet(Reboot, flag.ContinueOnError)
	flags.rebootDriver = setDriverFlag(flags.rebootCmd)
	flags.rebootCmd.Usage = func() {
		fmt.Printf(commandUsage, Reboot, os.Args[0], Reboot)
		flags.rebootCmd.PrintDefaults()
	}

	return flags
}

// Usage outputs examples to the screen.
func (p *Flags) Usage() {
	flag.Usage()
}

// SortField returns the field by which the dump results should be sorted by.
func (p *Flags) SortField() string {
	return p.dumpSortField.String()
}

// DumpFormat returns the dump data format value.
func (p *Flags) DumpFormat() string {
	return p.dumpFormat.String()
}

// DumpFile returns the file path value.
func (p *Flags) DumpFile() string {
	return *p.dumpFile
}

// ConfigFile returns the configuration file path value.
func (p *Flags) ConfigFile() string {
	return *p.configFile
}

// SecureFile returns the auth configuration file path value.
func (p *Flags) SecureFile() string {
	return *p.secureFile
}

// SecureOff returns true if device authentication should be turned off, false otherwise.
func (p *Flags) SecureOff() bool {
	return *p.secureOff
}

// DeployFile returns the deployment file path value.
func (p *Flags) DeployFile() string {
	return *p.deployFile
}

// Parse the CLI arguments.
func (p *Flags) Parse(arguments []string) (*flag.FlagSet, string, error) {
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
		err = p.dumpCmd.Parse(arguments[1:])
		if err != nil {
			return p.dumpCmd, "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return p.dumpCmd, p.dumpDriver.String(), nil

	case Config:
		err = p.configCmd.Parse(arguments[1:])
		if err != nil {
			return p.configCmd, "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return p.configCmd, p.configDriver.String(), nil

	case Secure:
		err = p.secureCmd.Parse(arguments[1:])
		if err != nil {
			return p.secureCmd, "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		if p.SecureOff() && p.SecureFile() != "" {
			return p.secureCmd, "", fmt.Errorf("%w: '-f' and '--off' flags cannot be used together", ErrFlagConflict)
		}

		return p.secureCmd, p.secureDriver.String(), nil

	case Version:
		err = p.versionCmd.Parse(arguments[1:])
		if err != nil {
			return p.versionCmd, "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return p.versionCmd, p.versionDriver.String(), nil

	case Update:
		err = p.updateCmd.Parse(arguments[1:])
		if err != nil {
			return p.updateCmd, "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return p.updateCmd, p.updateDriver.String(), nil

	case Deploy:
		err = p.deployCmd.Parse(arguments[1:])
		if err != nil {
			return p.deployCmd, "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return p.deployCmd, p.deployDriver.String(), nil

	case Reboot:
		err = p.rebootCmd.Parse(arguments[1:])
		if err != nil {
			return p.rebootCmd, "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return p.rebootCmd, p.rebootDriver.String(), nil

	default:
		return nil, "", fmt.Errorf("%w: %s", ErrInvalid, arguments[0])
	}
}
