package device

import (
	"encoding/json"
	"errors"
	"net"
	"reflect"
	"testing"
)

func TestPolicy_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		policy *Policy
		err    error
	}{
		{
			name:   "failure: invalid JSON",
			policy: &Policy{},
			err:    &json.SyntaxError{},
		},
		{
			name:   "failure: undefined policy mode #1",
			data:   "{}",
			policy: &Policy{},
			err:    errPolicyModeUndefined,
		},
		{
			name:   "failure: undefined policy mode #2",
			data:   `{"mode":""}`,
			policy: &Policy{},
			err:    errPolicyModeUndefined,
		},
		{
			name:   "failure: invalid policy mode",
			data:   `{"mode":"foo"}`,
			policy: &Policy{},
			err:    errPolicyModeInvalid,
		},
		{
			name: "success: whitelist policy",
			data: `{"mode":"whitelist"}`,
			policy: &Policy{
				Mode: PolicyModeWhitelist,
			},
		},
		{
			name: "success: blacklist policy",
			data: `{"mode":"blacklist"}`,
			policy: &Policy{
				Mode: PolicyModeBlacklist,
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
			policy: &Policy{
				Mode: PolicyModeBlacklist,
			},
			err: &net.AddrError{},
		},
		{
			name: "success: whitelist policy with MAC address",
			data: `{
			  "mode": "whitelist",
			  "devices": [
				"14:06:12:DC:7A:F0"
			  ]
			}`,
			policy: &Policy{
				Mode: PolicyModeWhitelist,
				Devices: []net.HardwareAddr{
					{20, 6, 18, 220, 122, 240},
				},
			},
		},
		{
			name: "success: blacklist policy with model names",
			data: `{
			  "mode": "blacklist",
			  "models": [
				"SNSW-001X16EU",
				"SNSW-001X8EU"
			  ]
			}`,
			policy: &Policy{
				Mode: PolicyModeBlacklist,
				Models: []string{
					"SNSW-001X16EU",
					"SNSW-001X8EU",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			policy := &Policy{}

			err := json.Unmarshal([]byte(test.data), policy)

			if !reflect.DeepEqual(policy, test.policy) {
				t.Fatalf("expected %#v, got %#v", test.policy, policy)
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

func TestPolicy_Contains(t *testing.T) {
	tests := []struct {
		name      string
		str       *Policy
		dev       Resource
		contained bool
	}{
		{
			name: "device model is not contained",
			str: &Policy{
				Mode: PolicyModeWhitelist,
				Models: []string{
					"SNSW-001X16EU",
				},
			},
			dev:       &resource{},
			contained: false,
		},
		{
			name: "device model is contained",
			str: &Policy{
				Mode: PolicyModeWhitelist,
				Models: []string{
					"SNSW-001X16EU",
				},
			},
			dev:       &resource{model: "SNSW-001X16EU"},
			contained: true,
		},
		{
			name: "device MAC address is not contained",
			str: &Policy{
				Mode: PolicyModeBlacklist,
				Devices: []net.HardwareAddr{
					macAddr,
				},
			},
			dev:       &resource{},
			contained: false,
		},
		{
			name: "device MAC address is contained",
			str: &Policy{
				Mode: PolicyModeWhitelist,
				Devices: []net.HardwareAddr{
					macAddr,
				},
			},
			dev:       &resource{mac: macAddr},
			contained: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			contained := test.str.Contains(test.dev)
			if contained != test.contained {
				t.Fatalf("expected %t, got %t", test.contained, contained)
			}
		})
	}
}

func TestPolicy_IsExcluded(t *testing.T) {
	tests := []struct {
		name     string
		str      *Policy
		dev      Resource
		excluded bool
	}{
		{
			name: "whitelist: device is excluded via model name",
			str: &Policy{
				Mode: PolicyModeWhitelist,
				Models: []string{
					"SNSW-001X16EU",
				},
			},
			dev:      &resource{},
			excluded: true,
		},
		{
			name: "whitelist: device is excluded via MAC address",
			str: &Policy{
				Mode: PolicyModeWhitelist,
				Devices: []net.HardwareAddr{
					macAddr,
				},
			},
			dev:      &resource{},
			excluded: true,
		},
		{
			name: "whitelist: device is included via model name",
			str: &Policy{
				Mode: PolicyModeWhitelist,
				Models: []string{
					"SNSW-001X16EU",
				},
			},
			dev:      &resource{model: "SNSW-001X16EU"},
			excluded: false,
		},
		{
			name: "whitelist: device is included via MAC address",
			str: &Policy{
				Mode: PolicyModeWhitelist,
				Devices: []net.HardwareAddr{
					macAddr,
				},
			},
			dev:      &resource{mac: macAddr},
			excluded: false,
		},
		{
			name: "blacklist: device is excluded via model name",
			str: &Policy{
				Mode: PolicyModeBlacklist,
				Models: []string{
					"SNSW-001X16EU",
				},
			},
			dev:      &resource{model: "SNSW-001X16EU"},
			excluded: true,
		},
		{
			name: "blacklist: device is excluded via MAC address",
			str: &Policy{
				Mode: PolicyModeBlacklist,
				Devices: []net.HardwareAddr{
					macAddr,
				},
			},
			dev:      &resource{mac: macAddr},
			excluded: true,
		},
		{
			name: "blacklist: device is included via model name",
			str: &Policy{
				Mode: PolicyModeBlacklist,
				Models: []string{
					"SNSW-001X16EU",
				},
			},
			dev:      &resource{},
			excluded: false,
		},
		{
			name: "blacklist: device is included via MAC address",
			str: &Policy{
				Mode: PolicyModeBlacklist,
				Devices: []net.HardwareAddr{
					macAddr,
				},
			},
			dev:      &resource{},
			excluded: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			excluded := test.str.IsExcluded(test.dev)
			if excluded != test.excluded {
				t.Fatalf("expected %t, got %t", test.excluded, excluded)
			}
		})
	}
}
