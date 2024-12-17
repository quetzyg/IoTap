package shellygen1

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/quetzyg/IoTap/device"
)

// config implementation for testing purposes.
type config struct{}

// Driver name of this Config implementation.
func (c *config) Driver() string {
	return "mismatch"
}

// Empty checks if the struct holding the configuration has a zero value.
func (c *config) Empty() bool {
	return true
}

func compareHTTPRequests(req1, req2 *http.Request) bool {
	// Compare URL
	if req1.URL.String() != req2.URL.String() {
		return false
	}

	// Compare HTTP Method
	if req1.Method != req2.Method {
		return false
	}

	// Compare Headers
	return reflect.DeepEqual(req1.Header, req2.Header)
}

func TestDevice_ConfigureRequests(t *testing.T) {
	tests := []struct {
		name string
		dev  *Device
		cfg  device.Config
		rs   []*http.Request
		err  error
	}{
		{
			name: "failure: driver mismatch",
			dev:  &Device{},
			cfg:  &config{},
			err:  device.ErrDriverMismatch,
		},
		{
			name: "failure: policy exclusion",
			dev:  &Device{},
			cfg: &Config{
				Policy: &device.Policy{
					Mode: device.PolicyModeWhitelist,
				},
			},
			err: device.ErrPolicyExcluded,
		},
		{
			name: "success: single setting",
			dev:  &Device{},
			cfg: &Config{
				Settings: &settings{
					"discoverable": true,
				},
			},
			rs: func() []*http.Request {
				r1, _ := request(&Device{}, "settings", url.Values{
					"discoverable": []string{"true"},
				})

				r2, _ := request(&Device{}, rebootPath, nil)

				return []*http.Request{r1, r2}
			}(),
		},
		{
			name: "success: multiple settings",
			dev:  &Device{},
			cfg: &Config{
				SettingsRelay: &[]*settings{
					{
						"name":           nil,
						"appliance_type": "lock",
						"default_state":  "off",
						"btn_type":       "detached",
						"btn_reverse":    true,
						"auto_on":        0,
						"auto_off":       3,
					},
				},
			},
			rs: func() []*http.Request {
				r1, _ := request(&Device{}, "settings/relay/0", url.Values{
					"name":           []string{"null"},
					"appliance_type": []string{"lock"},
					"default_state":  []string{"off"},
					"btn_type":       []string{"detached"},
					"btn_reverse":    []string{"true"},
					"auto_on":        []string{"0"},
					"auto_off":       []string{"3"},
				})

				r2, _ := request(&Device{}, "reboot", nil)

				return []*http.Request{r1, r2}
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rs, err := test.dev.ConfigureRequests(test.cfg)

			for i, r := range rs {
				if !compareHTTPRequests(r, test.rs[i]) {
					t.Fatalf("expected %#v, got %#v", test.rs, rs)
				}
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}
