package shellygen2

import (
	"bytes"
	"errors"
	"io"
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
				BLE: &settings{
					"config": map[string]any{
						"enable": true,
					},
				},
			},
			rs: func() []*http.Request {
				r1 := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Scheme: "http",
						Host:   "192.168.146.123",
						Path:   rpcPath,
					},
					Header: http.Header{},
					Body:   io.NopCloser(bytes.NewBufferString(`{"id":0,"src":"IoTap","method":"ble.SetConfig","params":{"config":{"enable":true}}}`)),
				}

				r1.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				r2 := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Scheme: "http",
						Host:   "192.168.146.123",
						Path:   rpcPath,
					},
					Header: http.Header{},
					Body:   io.NopCloser(bytes.NewBufferString(`{"id":0,"src":"IoTap","method":"Shelly.Reboot"}`)),
				}

				r2.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				return []*http.Request{r1, r2}
			}(),
		},
		{
			name: "success: settings slice",
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
				r1 := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Scheme: "http",
						Host:   "192.168.146.123",
						Path:   rpcPath,
					},
					Header: http.Header{},
					Body:   io.NopCloser(bytes.NewBufferString(`{"id":0,"src":"IoTap","method":"input.SetConfig","params":{"config":{"invert":true,"name":null,"type":"switch"},"id":0}}`)),
				}

				r1.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				r2 := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Scheme: "http",
						Host:   "192.168.146.123",
						Path:   rpcPath,
					},
					Header: http.Header{},
					Body:   io.NopCloser(bytes.NewBufferString(`{"id":0,"src":"IoTap","method":"Shelly.Reboot"}`)),
				}

				r2.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				return []*http.Request{r1, r2}
			}(),
		},
	}

	shelly2 := &Device{ip: net.ParseIP("192.168.146.123")}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rs, err := shelly2.ConfigureRequests(test.cfg)

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
