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
	Version = "version"
	Update  = "update"
	Deploy  = "deploy"
	Reboot  = "reboot"
)

// Usage strings
const (
	usage = `Usage:
%s <CIDR> <command> [flags]

Commands:
  dump    Output device scan results to STDOUT or to a file
  config  Apply configurations to multiple devices
  version Scan and check device versions
  update  Update outdated devices
  deploy  Deploy scripts to multiple devices
  reboot  Restart devices

Use %s <CIDR> <command> -h for more information about the command.
`
	commandUsage = `Usage of %s:
 %s <CIDR> %s [flags]

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
	driver := NewStrFlag(device.Driver, device.Driver, shellygen1.Driver, shellygen2.Driver)
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
	flags.dumpSortField = NewStrFlag(device.FieldName, device.FieldDriver, device.FieldIP, device.FieldMAC, device.FieldName, device.FieldModel)
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

// DeployFile returns the deployment file path value.
func (p *Flags) DeployFile() string {
	return *p.deployFile
}

// Parse the CLI arguments.
func (p *Flags) Parse(arguments []string) (string, string, error) {
	if len(arguments) == 0 {
		return "", "", ErrNotFound
	}

	err := flag.CommandLine.Parse(arguments)
	if err != nil {
		return "", "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
	}

	// Handle commands
	switch arguments[0] {
	case Dump:
		err = p.dumpCmd.Parse(arguments[1:])
		if err != nil {
			return "", "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return Dump, p.dumpDriver.String(), nil

	case Config:
		err = p.configCmd.Parse(arguments[1:])
		if err != nil {
			return "", "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return Config, p.configDriver.String(), nil

	case Version:
		err = p.versionCmd.Parse(arguments[1:])
		if err != nil {
			return "", "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return Version, p.versionDriver.String(), nil

	case Update:
		err = p.updateCmd.Parse(arguments[1:])
		if err != nil {
			return "", "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return Update, p.updateDriver.String(), nil

	case Deploy:
		err = p.deployCmd.Parse(arguments[1:])
		if err != nil {
			return "", "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return Deploy, p.deployDriver.String(), nil

	case Reboot:
		err = p.rebootCmd.Parse(arguments[1:])
		if err != nil {
			return "", "", fmt.Errorf("%w: %w", ErrArgumentParse, err)
		}

		return Reboot, p.rebootDriver.String(), nil

	default:
		return "", "", fmt.Errorf("%w: %s", ErrInvalid, arguments[0])
	}
}
