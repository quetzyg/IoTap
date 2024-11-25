package command

import (
	"errors"
	"testing"

	"github.com/Stowify/IoTune/device"
	"github.com/Stowify/IoTune/shellygen1"
	"github.com/Stowify/IoTune/shellygen2"
)

func TestFlags_Parse(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		command string
		driver  string
		err     error
	}{
		{
			name: "failure: nil argument",
			args: nil,
			err:  errArgumentMissing,
		},
		{
			name: "failure: empty argument",
			args: []string{""},
			err:  errCommandInvalid,
		},
		{
			name: "failure: invalid command argument",
			args: []string{"foo"},
			err:  errCommandInvalid,
		},

		// Dump
		{
			name: "failure: dump command with undefined flag",
			args: []string{Dump, "-foo"},
			err:  errArgumentParsing,
		},
		{
			name: "failure: dump command with invalid driver flag value",
			args: []string{Dump, "-driver"},
			err:  errArgumentParsing,
		},
		{
			name: "failure: dump command with invalid sort flag value",
			args: []string{Dump, "-sort"},
			err:  errArgumentParsing,
		},
		{
			name:    "success: config command with defined flags",
			args:    []string{Dump, "-driver", device.Driver, "-sort", device.FieldIP},
			command: Dump,
			driver:  device.Driver,
		},

		// Config
		{
			name: "failure: config command with undefined flag",
			args: []string{Config, "-foo"},
			err:  errArgumentParsing,
		},
		{
			name: "failure: config command with invalid driver flag value",
			args: []string{Config, "-driver"},
			err:  errArgumentParsing,
		},
		{
			name: "failure: config command with invalid file flag value",
			args: []string{Config, "-f"},
			err:  errArgumentParsing,
		},
		{
			name:    "success: config command with defined flags",
			args:    []string{Config, "-driver", shellygen1.Driver, "-f", "config.json"},
			command: Config,
			driver:  shellygen1.Driver,
		},

		// Version
		{
			name: "failure: version command with undefined flag",
			args: []string{Version, "-foo"},
			err:  errArgumentParsing,
		},
		{
			name: "failure: version command with invalid driver flag value",
			args: []string{Version, "-driver"},
			err:  errArgumentParsing,
		},
		{
			name:    "success: version command with defined flags",
			args:    []string{Version, "-driver", shellygen2.Driver},
			command: Version,
			driver:  shellygen2.Driver,
		},

		// Update
		{
			name: "failure: update command with undefined flag",
			args: []string{Update, "-foo"},
			err:  errArgumentParsing,
		},
		{
			name: "failure: update command with invalid driver flag value",
			args: []string{Update, "-driver"},
			err:  errArgumentParsing,
		},
		{
			name:    "success: update command with defined flags",
			args:    []string{Update, "-driver", device.Driver},
			command: Update,
			driver:  device.Driver,
		},

		// Script
		{
			name: "failure: script command with undefined flag",
			args: []string{Script, "-foo"},
			err:  errArgumentParsing,
		},
		{
			name: "failure: script command with invalid driver flag value",
			args: []string{Script, "-driver"},
			err:  errArgumentParsing,
		},
		{
			name: "failure: script command with invalid file flag value",
			args: []string{Script, "-f"},
			err:  errArgumentParsing,
		},
		{
			name:    "success: script command with defined flags",
			args:    []string{Script, "-driver", shellygen1.Driver, "-f", "script.js"},
			command: Script,
			driver:  shellygen1.Driver,
		},

		// Reboot
		{
			name: "failure: reboot command with undefined flag",
			args: []string{Reboot, "-foo"},
			err:  errArgumentParsing,
		},
		{
			name: "failure: reboot command with invalid driver flag value",
			args: []string{Reboot, "-driver"},
			err:  errArgumentParsing,
		},
		{
			name:    "success: reboot command with defined flags",
			args:    []string{Reboot, "-driver", shellygen2.Driver},
			command: Reboot,
			driver:  shellygen2.Driver,
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
