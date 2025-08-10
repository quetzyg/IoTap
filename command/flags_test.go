package command

import (
	"errors"
	"flag"
	"testing"

	"github.com/quetzyg/IoTap/device"
	"github.com/quetzyg/IoTap/shellygen1"
	"github.com/quetzyg/IoTap/shellygen2"
)

func TestFlags_Usage(t *testing.T) {
	NewFlags().Usage()
}

func TestFlags_Driver(t *testing.T) {
	drv := NewFlags().Driver()

	if drv != device.AllDrivers {
		t.Fatalf("expected %q, got %q", device.AllDrivers, drv)
	}
}

func TestFlags_ProbeTimeout(t *testing.T) {
	timeout := NewFlags().ProbeTimeout()

	if timeout != device.ProbeTimeout {
		t.Fatalf("expected %q, got %q", device.ProbeTimeout, timeout)
	}
}

func TestFlags_SortField(t *testing.T) {
	tests := []struct {
		err       error
		name      string
		command   string
		driver    string
		sortField string
		args      []string
	}{
		{
			name:      "get default sort field value",
			args:      []string{Dump},
			command:   Dump,
			driver:    device.AllDrivers,
			sortField: device.FieldName,
		},
		{
			name:      "get MAC sort field value",
			args:      []string{Dump, "-s", device.FieldMAC},
			command:   Dump,
			driver:    device.AllDrivers,
			sortField: device.FieldMAC,
		},
		{
			name:      "get default sort field value when invalid field is passed",
			args:      []string{Dump, "-s", "foo"},
			command:   Dump,
			driver:    "",
			err:       ErrArgumentParse,
			sortField: device.FieldName,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flags := NewFlags()

			cmd, driver, err := flags.Parse(test.args)

			if cmd.Name() != test.command {
				t.Fatalf("Unexpected command. Got %s, expected %s", cmd.Name(), test.command)
			}

			if driver != test.driver {
				t.Fatalf("Unexpected driver. Got %s, expected %s", driver, test.driver)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}

			sort := flags.SortField()

			if sort != test.sortField {
				t.Fatalf("Unexpected sort field. Got %s, expected %s", sort, test.sortField)
			}
		})
	}
}

func TestFlags_DumpFormat(t *testing.T) {
	tests := []struct {
		err        error
		name       string
		command    string
		driver     string
		dumpFormat string
		args       []string
	}{
		{
			name:       "get default dump format value",
			args:       []string{Dump},
			command:    Dump,
			driver:     device.AllDrivers,
			dumpFormat: device.FormatCSV,
		},
		{
			name:       "get JSON dump format value",
			args:       []string{Dump, "-f", device.FormatJSON},
			command:    Dump,
			driver:     device.AllDrivers,
			dumpFormat: device.FormatJSON,
		},
		{
			name:       "get default dump format value when invalid format is passed",
			args:       []string{Dump, "-f", "foo"},
			command:    Dump,
			driver:     "",
			err:        ErrArgumentParse,
			dumpFormat: device.FormatCSV,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flags := NewFlags()

			cmd, driver, err := flags.Parse(test.args)

			if cmd.Name() != test.command {
				t.Fatalf("Unexpected command. Got %s, expected %s", cmd.Name(), test.command)
			}

			if driver != test.driver {
				t.Fatalf("Unexpected driver. Got %s, expected %s", driver, test.driver)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}

			format := flags.DumpFormat()

			if format != test.dumpFormat {
				t.Fatalf("Unexpected dump format. Got %s, expected %s", format, test.dumpFormat)
			}
		})
	}
}

func TestFlags_File(t *testing.T) {
	tests := []struct {
		err     error
		name    string
		command string
		driver  string
		file    string
		args    []string
	}{
		// Dump
		{
			name:    "get empty output file path value",
			args:    []string{Dump},
			command: Dump,
			driver:  device.AllDrivers,
			file:    "",
		},
		{
			name:    "get output file path value",
			args:    []string{Dump, "-o", "devices.csv"},
			command: Dump,
			driver:  device.AllDrivers,
			file:    "devices.csv",
		},
		{
			name:    "get empty output file path value when argument is missing",
			args:    []string{Dump, "-o"},
			command: Dump,
			driver:  "",
			err:     ErrArgumentParse,
			file:    "",
		},

		// Config
		{
			name:    "get empty config file path value",
			args:    []string{Config},
			command: Config,
			driver:  device.AllDrivers,
			file:    "",
		},
		{
			name:    "get config file path value",
			args:    []string{Config, "-c", "config.json"},
			command: Config,
			driver:  device.AllDrivers,
			file:    "config.json",
		},
		{
			name:    "get empty config file path value when argument is missing",
			args:    []string{Config, "-c"},
			command: Config,
			driver:  "",
			err:     ErrArgumentParse,
			file:    "",
		},

		// Deploy
		{
			name:    "get empty deploy file path value",
			args:    []string{Deploy},
			command: Deploy,
			driver:  device.AllDrivers,
			file:    "",
		},
		{
			name:    "get single deploy file path value",
			args:    []string{Deploy, "-c", "deployment.json"},
			command: Deploy,
			driver:  device.AllDrivers,
			file:    "deployment.json",
		},
		{
			name:    "get empty deploy file path value when argument is missing",
			args:    []string{Deploy, "-c"},
			command: Deploy,
			driver:  "",
			err:     ErrArgumentParse,
			file:    "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flags := NewFlags()

			cmd, driver, err := flags.Parse(test.args)

			if cmd.Name() != test.command {
				t.Fatalf("Unexpected command. Got %s, expected %s", cmd.Name(), test.command)
			}

			if driver != test.driver {
				t.Fatalf("Unexpected driver. Got %s, expected %s", driver, test.driver)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}

			file := flags.File()

			if file != test.file {
				t.Fatalf("Unexpected deploy file. Got %q, expected %q", file, test.file)
			}
		})
	}
}

func TestFlags_Parse(t *testing.T) {
	tests := []struct {
		err     error
		name    string
		command string
		driver  string
		args    []string
	}{
		{
			name: "failure: command not found #1",
			args: nil,
			err:  ErrNotFound,
		},
		{
			name: "failure: command not found #2",
			args: []string{},
			err:  ErrNotFound,
		},
		{
			name: "failure: invalid command #1",
			args: []string{""},
			err:  ErrInvalid,
		},
		{
			name: "failure: invalid command #2",
			args: []string{"foo"},
			err:  ErrInvalid,
		},

		// Dump
		{
			name:    "failure: dump command with undefined flag",
			args:    []string{Dump, "-foo"},
			command: Dump,
			err:     ErrArgumentParse,
		},
		{
			name:    "failure: dump command with missing driver flag value",
			args:    []string{Dump, "-d"},
			command: Dump,
			err:     ErrArgumentParse,
		},
		{
			name:    "failure: dump command with missing sort flag value",
			args:    []string{Dump, "-s"},
			command: Dump,
			err:     ErrArgumentParse,
		},
		{
			name:    "failure: dump command with missing format flag value",
			args:    []string{Dump, "-f"},
			command: Dump,
			err:     ErrArgumentParse,
		},
		{
			name:    "failure: dump command with missing output file flag value",
			args:    []string{Dump, "-o"},
			command: Dump,
			err:     ErrArgumentParse,
		},
		{
			name: "success: dump command with valid flags",
			args: []string{
				Dump,
				"-d",
				device.AllDrivers,
				"-s",
				device.FieldIP,
				"-f",
				device.FormatCSV,
				"-o",
				"devices.csv",
			},
			command: Dump,
			driver:  device.AllDrivers,
		},
		{
			name:    "success: dump command with help flag",
			args:    []string{Dump, "-h"},
			command: Dump,
			err:     flag.ErrHelp,
		},

		// Config
		{
			name:    "failure: config command with undefined flag",
			args:    []string{Config, "-foo"},
			command: Config,
			err:     ErrArgumentParse,
		},
		{
			name:    "failure: config command with missing driver flag value",
			args:    []string{Config, "-d"},
			command: Config,
			err:     ErrArgumentParse,
		},
		{
			name:    "failure: config command with missing config flag value",
			args:    []string{Config, "-c"},
			command: Config,
			err:     ErrArgumentParse,
		},
		{
			name:    "success: config command with valid flags",
			args:    []string{Config, "-d", shellygen1.Driver, "-c", "config.json"},
			command: Config,
			driver:  shellygen1.Driver,
		},
		{
			name:    "success: config command with help flag",
			args:    []string{Config, "-h"},
			command: Config,
			err:     flag.ErrHelp,
		},

		// Secure
		{
			name:    "failure: secure command with undefined flag",
			args:    []string{Secure, "-foo"},
			command: Secure,
			err:     ErrArgumentParse,
		},
		{
			name:    "failure: secure command with missing driver flag value",
			args:    []string{Secure, "-d"},
			command: Secure,
			err:     ErrArgumentParse,
		},
		{
			name:    "failure: secure command with missing config flag value",
			args:    []string{Secure, "-c"},
			command: Secure,
			err:     ErrArgumentParse,
		},
		{
			name:    "failure: secure command with conflicting flags",
			args:    []string{Secure, "-c", "secure.json", "--off"},
			command: Secure,
			err:     ErrFlagConflict,
		},
		{
			name:    "success: secure command with valid flags #1",
			args:    []string{Secure, "-d", shellygen1.Driver, "-c", "secure.json"},
			command: Secure,
			driver:  shellygen1.Driver,
		},
		{
			name:    "success: secure command with valid flags #2",
			args:    []string{Secure, "-d", shellygen2.Driver, "--off"},
			command: Secure,
			driver:  shellygen2.Driver,
		},
		{
			name:    "success: config command with help flag",
			args:    []string{Secure, "-h"},
			command: Secure,
			err:     flag.ErrHelp,
		},

		// Version
		{
			name:    "failure: version command with undefined flag",
			args:    []string{Version, "-foo"},
			command: Version,
			err:     ErrArgumentParse,
		},
		{
			name:    "failure: version command with missing driver flag value",
			args:    []string{Version, "-d"},
			command: Version,
			err:     ErrArgumentParse,
		},
		{
			name:    "success: version command with valid flags",
			args:    []string{Version, "-d", shellygen2.Driver},
			command: Version,
			driver:  shellygen2.Driver,
		},
		{
			name:    "success: version command with help flag",
			args:    []string{Version, "-h"},
			command: Version,
			err:     flag.ErrHelp,
		},

		// Update
		{
			name:    "failure: update command with undefined flag",
			args:    []string{Update, "-foo"},
			command: Update,
			err:     ErrArgumentParse,
		},
		{
			name:    "failure: update command with missing driver flag value",
			args:    []string{Update, "-d"},
			command: Update,
			err:     ErrArgumentParse,
		},
		{
			name:    "success: update command with valid flags",
			args:    []string{Update, "-d", device.AllDrivers},
			command: Update,
			driver:  device.AllDrivers,
		},
		{
			name:    "success: update command with help flag",
			args:    []string{Update, "-h"},
			command: Update,
			err:     flag.ErrHelp,
		},

		// Deploy
		{
			name:    "failure: deploy command with undefined flag",
			args:    []string{Deploy, "-foo"},
			command: Deploy,
			err:     ErrArgumentParse,
		},
		{
			name:    "failure: deploy command with missing driver flag value",
			args:    []string{Deploy, "-d"},
			command: Deploy,
			err:     ErrArgumentParse,
		},
		{
			name:    "failure: deploy command with missing config flag value",
			args:    []string{Deploy, "-c"},
			command: Deploy,
			err:     ErrArgumentParse,
		},
		{
			name:    "success: deploy command with valid flags",
			args:    []string{Deploy, "-d", shellygen1.Driver, "-c", "deployment.json"},
			command: Deploy,
			driver:  shellygen1.Driver,
		},
		{
			name:    "success: deploy command with help flag",
			args:    []string{Deploy, "-h"},
			command: Deploy,
			err:     flag.ErrHelp,
		},

		// Reboot
		{
			name:    "failure: reboot command with undefined flag",
			args:    []string{Reboot, "-foo"},
			command: Reboot,
			err:     ErrArgumentParse,
		},
		{
			name:    "failure: reboot command with missing driver flag value",
			args:    []string{Reboot, "-d"},
			command: Reboot,
			err:     ErrArgumentParse,
		},
		{
			name:    "success: reboot command with valid flags",
			args:    []string{Reboot, "-d", shellygen2.Driver},
			command: Reboot,
			driver:  shellygen2.Driver,
		},
		{
			name:    "success: reboot command with help flag",
			args:    []string{Reboot, "-h"},
			command: Reboot,
			err:     flag.ErrHelp,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd, driver, err := NewFlags().Parse(test.args)

			if cmd != nil && cmd.Name() != test.command {
				t.Fatalf("Unexpected command. Got %s, expected %s", cmd.Name(), test.command)
			}

			if driver != test.driver {
				t.Fatalf("Unexpected driver. Got %s, expected %s", driver, test.driver)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}
