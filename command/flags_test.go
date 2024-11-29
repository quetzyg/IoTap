package command

import (
	"errors"
	"flag"
	"testing"

	"github.com/Stowify/IoTune/device"
	"github.com/Stowify/IoTune/shellygen1"
	"github.com/Stowify/IoTune/shellygen2"
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
			sortField: "name",
		},
		{
			name:      "get mac sort field value",
			args:      []string{Dump, "-sort", "mac"},
			command:   Dump,
			driver:    device.Driver,
			sortField: "mac",
		},
		{
			name:      "get default sort field value when invalid field is passed",
			args:      []string{Dump, "-sort", "foo"},
			command:   "",
			driver:    "",
			err:       ErrArgumentParse,
			sortField: "name",
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
			args:     []string{Dump, "-f", "devices.json"},
			command:  Dump,
			driver:   device.Driver,
			dumpFile: "devices.json",
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

func TestFlags_ScriptFile(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		command    string
		driver     string
		err        error
		scriptFile string
	}{
		{
			name:       "get empty script file path value",
			args:       []string{Script},
			command:    Script,
			driver:     device.Driver,
			scriptFile: "",
		},
		{
			name:       "get script file path value",
			args:       []string{Script, "-f", "script.js"},
			command:    Script,
			driver:     device.Driver,
			scriptFile: "script.js",
		},
		{
			name:       "get empty script file path value when argument is missing",
			args:       []string{Script, "-f"},
			command:    "",
			driver:     "",
			err:        ErrArgumentParse,
			scriptFile: "",
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

			file := flags.ScriptFile()

			if file != test.scriptFile {
				t.Fatalf("Unexpected script file. Got %s, expected %s", file, test.scriptFile)
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
			name: "failure: dump command with invalid file flag value",
			args: []string{Dump, "-f"},
			err:  ErrArgumentParse,
		},
		{
			name:    "success: dump command with valid flags",
			args:    []string{Dump, "-driver", device.Driver, "-sort", device.FieldIP, "-f", "devices.json"},
			command: Dump,
			driver:  device.Driver,
		},
		{
			name: "success: dump command with valid + help flags",
			args: []string{Dump, "-driver", device.Driver, "-sort", device.FieldIP, "-f", "devices.json", "-h"},
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

		// Script
		{
			name: "failure: script command with undefined flag",
			args: []string{Script, "-foo"},
			err:  ErrArgumentParse,
		},
		{
			name: "failure: script command with invalid driver flag value",
			args: []string{Script, "-driver"},
			err:  ErrArgumentParse,
		},
		{
			name: "failure: script command with invalid file flag value",
			args: []string{Script, "-f"},
			err:  ErrArgumentParse,
		},
		{
			name:    "success: script command with valid flags",
			args:    []string{Script, "-driver", shellygen1.Driver, "-f", "script.js"},
			command: Script,
			driver:  shellygen1.Driver,
		},
		{
			name: "success: script command with valid + help flags",
			args: []string{Script, "-driver", shellygen1.Driver, "-f", "script.js", "-h"},
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
