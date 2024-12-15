package device

import (
	"encoding/json"
	"errors"
	"net"
	"testing"
)

func TestStrategy_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		data string
		err  error
	}{
		{
			name: "failure: invalid JSON",
			err:  &json.SyntaxError{},
		},
		{
			name: "failure: undefined strategy mode #1",
			data: "{}",
			err:  errStrategyModeUndefined,
		},
		{
			name: "failure: undefined strategy mode #2",
			data: `{"mode":""}`,
			err:  errStrategyModeUndefined,
		},
		{
			name: "failure: invalid strategy mode",
			data: `{"mode":"foo"}`,
			err:  errStrategyModeInvalid,
		},
		{
			name: "success: whitelist strategy",
			data: `{"mode":"whitelist"}`,
		},
		{
			name: "success: blacklist strategy",
			data: `{"mode":"blacklist"}`,
		},
		{
			name: "failure: blacklist strategy invalid device",
			data: `{
			  "mode": "blacklist",
			  "devices": [
				"foo"
			  ]
			}`,
			err: &net.AddrError{},
		},
		{
			name: "success: whitelist strategy invalid device",
			data: `{
			  "mode": "whitelist",
			  "devices": [
				"14:06:12:DC:7A:F0"
			  ]
			}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := json.Unmarshal([]byte(test.data), &Strategy{})

			var syntaxError *json.SyntaxError
			var addrError *net.AddrError
			switch {
			case errors.As(test.err, &syntaxError):
				var se *json.SyntaxError
				if errors.As(err, &se) {
					return
				}

			case errors.As(test.err, &addrError):
				var ae *net.AddrError
				if errors.As(err, &ae) {
					return
				}

			default:
				if errors.Is(err, test.err) {
					return
				}
			}

			t.Fatalf("expected %#v, got %#v", test.err, err)
		})
	}
}

type devStrategy struct{ mac net.HardwareAddr }

func (d *devStrategy) Driver() string { return "" }

func (d *devStrategy) IP() net.IP { return nil }

func (d *devStrategy) MAC() net.HardwareAddr { return d.mac }

func (d *devStrategy) Name() string { return "" }

func (d *devStrategy) Model() string { return "" }

func (d *devStrategy) ID() string { return "" }

func (d *devStrategy) Secured() bool { return false }

var mac = net.HardwareAddr{20, 6, 18, 220, 122, 240}

func TestStrategy_Listed(t *testing.T) {
	tests := []struct {
		name   string
		str    *Strategy
		dev    Resource
		listed bool
	}{
		{
			name: "device is not listed",
			str: &Strategy{
				mode: blacklist,
				devices: []net.HardwareAddr{
					mac,
				},
			},
			dev: &devStrategy{},
		},
		{
			name: "device is listed",
			str: &Strategy{
				mode: whitelist,
				devices: []net.HardwareAddr{
					mac,
				},
			},
			dev:    &devStrategy{mac: mac},
			listed: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			listed := test.str.Listed(test.dev)
			if listed != test.listed {
				t.Fatalf("expected %t, got %t", test.listed, listed)
			}
		})
	}
}

func TestStrategy_Excluded(t *testing.T) {
	tests := []struct {
		name     string
		str      *Strategy
		dev      Resource
		excluded bool
	}{
		{
			name: "whitelist: device is excluded",
			str: &Strategy{
				mode: whitelist,
				devices: []net.HardwareAddr{
					mac,
				},
			},
			dev:      &devStrategy{},
			excluded: true,
		},
		{
			name: "whitelist: device is included",
			str: &Strategy{
				mode: whitelist,
				devices: []net.HardwareAddr{
					mac,
				},
			},
			dev: &devStrategy{mac: mac},
		},
		{
			name: "blacklist: device is excluded",
			str: &Strategy{
				mode: blacklist,
				devices: []net.HardwareAddr{
					mac,
				},
			},
			dev:      &devStrategy{mac: mac},
			excluded: true,
		},
		{
			name: "blacklist: device is included",
			str: &Strategy{
				mode: blacklist,
				devices: []net.HardwareAddr{
					mac,
				},
			},
			dev: &devStrategy{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			excluded := test.str.Excluded(test.dev)
			if excluded != test.excluded {
				t.Fatalf("expected %t, got %t", test.excluded, excluded)
			}
		})
	}
}
