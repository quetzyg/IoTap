package shellygen2

import (
	"bytes"
	"errors"
	"io"
	"net/http"
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
	if !reflect.DeepEqual(req1.Header, req2.Header) {
		return false
	}

	// Compare Body
	body1, _ := io.ReadAll(req1.Body)
	body2, _ := io.ReadAll(req2.Body)

	return bytes.Equal(body1, body2)
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
			name: "success: single settings",
			dev:  &Device{},
			cfg: &Config{
				BLE: &settings{
					"config": map[string]any{
						"enable": true,
					},
				},
			},
			rs: func() []*http.Request {
				r1, _ := request(&Device{}, "ble.SetConfig", map[string]any{
					"config": map[string]any{
						"enable": true,
					},
				})

				r2, _ := request(&Device{}, "Shelly.Reboot", nil)

				return []*http.Request{r1, r2}
			}(),
		},
		{
			name: "success: settings slice",
			dev:  &Device{},
			cfg: &Config{
				Input: &[]*settings{
					{
						"id": 0,
						"config": map[string]any{
							"name":   nil,
							"type":   "switch",
							"invert": true,
						},
					},
				},
			},
			rs: func() []*http.Request {
				r1, _ := request(&Device{}, "input.SetConfig", map[string]any{
					"id": 0,
					"config": map[string]any{
						"name":   nil,
						"type":   "switch",
						"invert": true,
					},
				})

				r2, _ := request(&Device{}, "Shelly.Reboot", nil)

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
