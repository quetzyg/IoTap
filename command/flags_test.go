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

func TestFlags_SortField(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		command   string
		driver    string
		err       error
		sortField string
	}{
		{
			name:      "get default sort field value",
			args:      []string{Dump},
			command:   Dump,
			driver:    device.Driver,
			sortField: device.FieldName,
		},
		{
			name:      "get MAC sort field value",
			args:      []string{Dump, "-sort", device.FieldMAC},
			command:   Dump,
			driver:    device.Driver,
			sortField: device.FieldMAC,
		},
		{
			name:      "get default sort field value when invalid field is passed",
			args:      []string{Dump, "-sort", "foo"},
			command:   "",
			driver:    "",
			err:       ErrArgumentParse,
			sortField: device.FieldName,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flags := NewFlags()

			cmd, driver, err := flags.Parse(test.args)

			if cmd != test.command {
				t.Fatalf("Unexpected command. Got %s, expected %s", cmd, test.command)
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
		name       string
		args       []string
		command    string
		driver     string
		err        error
		dumpFormat string
	}{
		{
			name:       "get default dump format value",
			args:       []string{Dump},
			command:    Dump,
			driver:     device.Driver,
			dumpFormat: device.FormatCSV,
		},
		{
			name:       "get JSON dump format value",
			args:       []string{Dump, "-format", device.FormatJSON},
			command:    Dump,
			driver:     device.Driver,
			dumpFormat: device.FormatJSON,
		},
		{
			name:       "get default dump format value when invalid format is passed",
			args:       []string{Dump, "-format", "foo"},
			command:    "",
			driver:     "",
			err:        ErrArgumentParse,
			dumpFormat: device.FormatCSV,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flags := NewFlags()

			cmd, driver, err := flags.Parse(test.args)

			if cmd != test.command {
				t.Fatalf("Unexpected command. Got %s, expected %s", cmd, test.command)
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

func TestFlags_DumpFile(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		command  string
		driver   string
		err      error
		dumpFile string
	}{
		{
			name:     "get empty dump file path value",
			args:     []string{Dump},
			command:  Dump,
			driver:   device.Driver,
			dumpFile: "",
		},
		{
			name:     "get dump file path value",
			args:     []string{Dump, "-f", "devices.csv"},
			command:  Dump,
			driver:   device.Driver,
			dumpFile: "devices.csv",
		},
		{
			name:     "get empty dump file path value when argument is missing",
			args:     []string{Dump, "-f"},
			command:  "",
			driver:   "",
			err:      ErrArgumentParse,
			dumpFile: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flags := NewFlags()

			cmd, driver, err := flags.Parse(test.args)

			if cmd != test.command {
				t.Fatalf("Unexpected command. Got %s, expected %s", cmd, test.command)
			}

			if driver != test.driver {
				t.Fatalf("Unexpected driver. Got %s, expected %s", driver, test.driver)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}

			file := flags.DumpFile()

			if file != test.dumpFile {
				t.Fatalf("Unexpected dump file. Got %s, expected %s", file, test.dumpFile)
			}
		})
	}
}

func TestFlags_ConfigFile(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		command    string
		driver     string
		err        error
		configFile string
	}{
		{
			name:       "get empty config file path value",
			args:       []string{Config},
			command:    Config,
			driver:     device.Driver,
			configFile: "",
		},
		{
			name:       "get config file path value",
			args:       []string{Config, "-f", "config.json"},
			command:    Config,
			driver:     device.Driver,
			configFile: "config.json",
		},
		{
			name:       "get empty config file path value when argument is missing",
			args:       []string{Config, "-f"},
			command:    "",
			driver:     "",
			err:        ErrArgumentParse,
			configFile: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flags := NewFlags()

			cmd, driver, err := flags.Parse(test.args)

			if cmd != test.command {
				t.Fatalf("Unexpected command. Got %s, expected %s", cmd, test.command)
			}

			if driver != test.driver {
				t.Fatalf("Unexpected driver. Got %s, expected %s", driver, test.driver)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}

			file := flags.ConfigFile()

			if file != test.configFile {
				t.Fatalf("Unexpected config file. Got %s, expected %s", file, test.configFile)
			}
		})
	}
}

func TestFlags_DeployFile(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		command    string
		driver     string
		err        error
		deployFile string
	}{
		{
			name:       "get empty deploy file path value",
			args:       []string{Deploy},
			command:    Deploy,
			driver:     device.Driver,
			deployFile: "",
		},
		{
			name:       "get single deploy file path value",
			args:       []string{Deploy, "-f", "deployment.json"},
			command:    Deploy,
			driver:     device.Driver,
			deployFile: "deployment.json",
		},
		{
			name:       "get empty deploy file path value when argument is missing",
			args:       []string{Deploy, "-f"},
			command:    "",
			driver:     "",
			err:        ErrArgumentParse,
			deployFile: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flags := NewFlags()

			cmd, driver, err := flags.Parse(test.args)

			if cmd != test.command {
				t.Fatalf("Unexpected command. Got %s, expected %s", cmd, test.command)
			}

			if driver != test.driver {
				t.Fatalf("Unexpected driver. Got %s, expected %s", driver, test.driver)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}

			file := flags.DeployFile()

			if file != test.deployFile {
				t.Fatalf("Unexpected deploy file. Got %q, expected %q", file, test.deployFile)
			}
		})
	}
}

func TestFlags_Parse(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		command string
		driver  string
		err     error
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
			name: "failure: dump command with undefined flag",
			args: []string{Dump, "-foo"},
			err:  ErrArgumentParse,
		},
		{
			name: "failure: dump command with invalid driver flag value",
			args: []string{Dump, "-driver"},
			err:  ErrArgumentParse,
		},
		{
			name: "failure: dump command with invalid sort flag value",
			args: []string{Dump, "-sort"},
			err:  ErrArgumentParse,
		},
		{
			name: "failure: dump command with invalid format flag value",
			args: []string{Dump, "-format"},
			err:  ErrArgumentParse,
		},
		{
			name: "failure: dump command with invalid file flag value",
			args: []string{Dump, "-f"},
			err:  ErrArgumentParse,
		},
		{
			name:    "success: dump command with valid flags",
			args:    []string{Dump, "-driver", device.Driver, "-sort", device.FieldIP, "-format", device.FormatCSV, "-f", "devices.csv"},
			command: Dump,
			driver:  device.Driver,
		},
		{
			name: "success: dump command with valid + help flags",
			args: []string{Dump, "-driver", device.Driver, "-sort", device.FieldIP, "-format", device.FormatCSV, "-f", "devices.csv", "-h"},
			err:  flag.ErrHelp,
		},

		// Config
		{
			name: "failure: config command with undefined flag",
			args: []string{Config, "-foo"},
			err:  ErrArgumentParse,
		},
		{
			name: "failure: config command with invalid driver flag value",
			args: []string{Config, "-driver"},
			err:  ErrArgumentParse,
		},
		{
			name: "failure: config command with invalid file flag value",
			args: []string{Config, "-f"},
			err:  ErrArgumentParse,
		},
		{
			name:    "success: config command with valid flags",
			args:    []string{Config, "-driver", shellygen1.Driver, "-f", "config.json"},
			command: Config,
			driver:  shellygen1.Driver,
		},
		{
			name: "success: config command with valid + help flags",
			args: []string{Config, "-driver", shellygen1.Driver, "-f", "config.json", "-h"},
			err:  flag.ErrHelp,
		},

		// Version
		{
			name: "failure: version command with undefined flag",
			args: []string{Version, "-foo"},
			err:  ErrArgumentParse,
		},
		{
			name: "failure: version command with invalid driver flag value",
			args: []string{Version, "-driver"},
			err:  ErrArgumentParse,
		},
		{
			name:    "success: version command with valid flags",
			args:    []string{Version, "-driver", shellygen2.Driver},
			command: Version,
			driver:  shellygen2.Driver,
		},
		{
			name: "success: version command with valid + help flags",
			args: []string{Version, "-driver", shellygen2.Driver, "-h"},
			err:  flag.ErrHelp,
		},

		// Update
		{
			name: "failure: update command with undefined flag",
			args: []string{Update, "-foo"},
			err:  ErrArgumentParse,
		},
		{
			name: "failure: update command with invalid driver flag value",
			args: []string{Update, "-driver"},
			err:  ErrArgumentParse,
		},
		{
			name:    "success: update command with valid flags",
			args:    []string{Update, "-driver", device.Driver},
			command: Update,
			driver:  device.Driver,
		},
		{
			name: "success: update command with valid + help flags",
			args: []string{Update, "-driver", device.Driver, "-h"},
			err:  flag.ErrHelp,
		},

		// Deploy
		{
			name: "failure: deploy command with undefined flag",
			args: []string{Deploy, "-foo"},
			err:  ErrArgumentParse,
		},
		{
			name: "failure: deploy command with invalid driver flag value",
			args: []string{Deploy, "-driver"},
			err:  ErrArgumentParse,
		},
		{
			name: "failure: deploy command with invalid file flag value",
			args: []string{Deploy, "-f"},
			err:  ErrArgumentParse,
		},
		{
			name:    "success: deploy command with valid flags",
			args:    []string{Deploy, "-driver", shellygen1.Driver, "-f", "script.js"},
			command: Deploy,
			driver:  shellygen1.Driver,
		},
		{
			name: "success: deploy command with valid + help flags",
			args: []string{Deploy, "-driver", shellygen1.Driver, "-f", "script.js", "-h"},
			err:  flag.ErrHelp,
		},

		// Reboot
		{
			name: "failure: reboot command with undefined flag",
			args: []string{Reboot, "-foo"},
			err:  ErrArgumentParse,
		},
		{
			name: "failure: reboot command with invalid driver flag value",
			args: []string{Reboot, "-driver"},
			err:  ErrArgumentParse,
		},
		{
			name:    "success: reboot command with valid flags",
			args:    []string{Reboot, "-driver", shellygen2.Driver},
			command: Reboot,
			driver:  shellygen2.Driver,
		},
		{
			name: "success: reboot command with valid + help flags",
			args: []string{Reboot, "-driver", shellygen2.Driver, "-h"},
			err:  flag.ErrHelp,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd, driver, err := NewFlags().Parse(test.args)

			if cmd != test.command {
				t.Fatalf("Unexpected command. Got %s, expected %s", cmd, test.command)
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
