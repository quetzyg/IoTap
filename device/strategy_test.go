package device

import (
	"encoding/json"
	"errors"
	"net"
	"reflect"
	"testing"
)

func TestStrategy_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		strategy *Strategy
		err      error
	}{
		{
			name:     "failure: invalid JSON",
			strategy: &Strategy{},
			err:      &json.SyntaxError{},
		},
		{
			name:     "failure: undefined strategy mode #1",
			data:     "{}",
			strategy: &Strategy{},
			err:      errStrategyModeUndefined,
		},
		{
			name:     "failure: undefined strategy mode #2",
			data:     `{"mode":""}`,
			strategy: &Strategy{},
			err:      errStrategyModeUndefined,
		},
		{
			name:     "failure: invalid strategy mode",
			data:     `{"mode":"foo"}`,
			strategy: &Strategy{},
			err:      errStrategyModeInvalid,
		},
		{
			name: "success: whitelist strategy",
			data: `{"mode":"whitelist"}`,
			strategy: &Strategy{
				mode: whitelist,
			},
		},
		{
			name: "success: blacklist strategy",
			data: `{"mode":"blacklist"}`,
			strategy: &Strategy{
				mode: blacklist,
			},
		},
		{
			name: "failure: invalid device MAC address",
			data: `{
			  "mode": "blacklist",
			  "devices": [
				"foo"
			  ]
			}`,
			strategy: &Strategy{
				mode: blacklist,
			},
			err: &net.AddrError{},
		},
		{
			name: "success: whitelist strategy with MAC address",
			data: `{
			  "mode": "whitelist",
			  "devices": [
				"14:06:12:DC:7A:F0"
			  ]
			}`,
			strategy: &Strategy{
				mode: whitelist,
				devices: []net.HardwareAddr{
					{20, 6, 18, 220, 122, 240},
				},
			},
		},
		{
			name: "success: blacklist strategy with model names",
			data: `{
			  "mode": "blacklist",
			  "models": [
				"SNSW-001X16EU",
				"SNSW-001X8EU"
			  ]
			}`,
			strategy: &Strategy{
				mode: blacklist,
				models: []string{
					"SNSW-001X16EU",
					"SNSW-001X8EU",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strategy := &Strategy{}

			err := json.Unmarshal([]byte(test.data), strategy)

			if !reflect.DeepEqual(strategy, test.strategy) {
				t.Fatalf("expected %#v, got %#v", test.strategy, strategy)
			}

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

			case errors.Is(err, test.err):
				return

			default:
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}

var macAddr = net.HardwareAddr{20, 6, 18, 220, 122, 240}

func TestStrategy_Listed(t *testing.T) {
	tests := []struct {
		name   string
		str    *Strategy
		dev    Resource
		listed bool
	}{
		{
			name: "device model is not listed",
			str: &Strategy{
				mode: whitelist,
				models: []string{
					"SNSW-001X16EU",
				},
			},
			dev:    &resource{},
			listed: false,
		},
		{
			name: "device model is listed",
			str: &Strategy{
				mode: whitelist,
				models: []string{
					"SNSW-001X16EU",
				},
			},
			dev:    &resource{model: "SNSW-001X16EU"},
			listed: true,
		},
		{
			name: "device MAC address is not listed",
			str: &Strategy{
				mode: blacklist,
				devices: []net.HardwareAddr{
					macAddr,
				},
			},
			dev:    &resource{},
			listed: false,
		},
		{
			name: "device MAC address is listed",
			str: &Strategy{
				mode: whitelist,
				devices: []net.HardwareAddr{
					macAddr,
				},
			},
			dev:    &resource{mac: macAddr},
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
			name: "whitelist: device is excluded via model name",
			str: &Strategy{
				mode: whitelist,
				models: []string{
					"SNSW-001X16EU",
				},
			},
			dev:      &resource{},
			excluded: true,
		},
		{
			name: "whitelist: device is excluded via MAC address",
			str: &Strategy{
				mode: whitelist,
				devices: []net.HardwareAddr{
					macAddr,
				},
			},
			dev:      &resource{},
			excluded: true,
		},
		{
			name: "whitelist: device is included via model name",
			str: &Strategy{
				mode: whitelist,
				models: []string{
					"SNSW-001X16EU",
				},
			},
			dev:      &resource{model: "SNSW-001X16EU"},
			excluded: false,
		},
		{
			name: "whitelist: device is included via MAC address",
			str: &Strategy{
				mode: whitelist,
				devices: []net.HardwareAddr{
					macAddr,
				},
			},
			dev:      &resource{mac: macAddr},
			excluded: false,
		},
		{
			name: "blacklist: device is excluded via model name",
			str: &Strategy{
				mode: blacklist,
				models: []string{
					"SNSW-001X16EU",
				},
			},
			dev:      &resource{model: "SNSW-001X16EU"},
			excluded: true,
		},
		{
			name: "blacklist: device is excluded via MAC address",
			str: &Strategy{
				mode: blacklist,
				devices: []net.HardwareAddr{
					macAddr,
				},
			},
			dev:      &resource{mac: macAddr},
			excluded: true,
		},
		{
			name: "blacklist: device is included via model name",
			str: &Strategy{
				mode: blacklist,
				models: []string{
					"SNSW-001X16EU",
				},
			},
			dev:      &resource{},
			excluded: false,
		},
		{
			name: "blacklist: device is included via MAC address",
			str: &Strategy{
				mode: blacklist,
				devices: []net.HardwareAddr{
					macAddr,
				},
			},
			dev:      &resource{},
			excluded: false,
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
