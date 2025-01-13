package shellygen1

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/quetzyg/IoTap/device"
	"github.com/quetzyg/IoTap/httpclient"
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
	// Compare HTTP Method
	if req1.Method != req2.Method {
		return false
	}

	// Compare URL
	if req1.URL.String() != req2.URL.String() {
		return false
	}

	// Compare Headers
	return reflect.DeepEqual(req1.Header, req2.Header)
}

func TestDevice_ConfigureRequests(t *testing.T) {
	tests := []struct {
		name string
		cfg  device.Config
		rs   []*http.Request
		err  error
	}{
		{
			name: "failure: driver mismatch",
			cfg:  &config{},
			err:  device.ErrDriverMismatch,
		},
		{
			name: "failure: policy exclusion",
			cfg: &Config{
				Policy: &device.Policy{
					Mode: device.PolicyModeWhitelist,
				},
			},
			err: device.ErrPolicyExcluded,
		},
		{
			name: "success: single settings",
			cfg: &Config{
				Settings: &settings{
					"discoverable": true,
				},
			},
			rs: func() []*http.Request {
				r1 := &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme:   "http",
						Host:     "192.168.146.123",
						Path:     "settings",
						RawQuery: "discoverable=true",
					},
					Header: http.Header{},
				}

				r1.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				r2 := &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme: "http",
						Host:   "192.168.146.123",
						Path:   rebootPath,
					},
					Header: http.Header{},
				}

				r2.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				return []*http.Request{r1, r2}
			}(),
		},
		{
			name: "success: settings slice",
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
				r1 := &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme:   "http",
						Host:     "192.168.146.123",
						Path:     "settings/relay/0",
						RawQuery: "appliance_type=lock&auto_off=3&auto_on=0&btn_reverse=true&btn_type=detached&default_state=off&name=null",
					},
					Header: http.Header{},
				}

				r1.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				r2 := &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme: "http",
						Host:   "192.168.146.123",
						Path:   rebootPath,
					},
					Header: http.Header{},
				}

				r2.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				return []*http.Request{r1, r2}
			}(),
		},
	}

	shelly1 := &Device{ip: net.ParseIP("192.168.146.123")}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rs, err := shelly1.ConfigureRequests(test.cfg)

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
