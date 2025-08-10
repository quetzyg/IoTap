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

func compareRequests(t *testing.T, expected, actual *http.Request) {
	if expected.Method != actual.Method {
		t.Fatalf("expected %q, got %q", expected.Method, actual.Method)
	}

	if expected.URL.String() != actual.URL.String() {
		t.Fatalf("expected %q, got %q", expected.URL.String(), actual.URL.String())
	}

	if !reflect.DeepEqual(expected.Header, actual.Header) {
		t.Fatalf("expected %#v, got %#v", expected.Header, actual.Header)
	}

	body1, _ := io.ReadAll(expected.Body)
	body2, _ := io.ReadAll(actual.Body)

	if !bytes.Equal(body1, body2) {
		t.Fatalf("expected %q, got %q", body1, body2)
	}
}

func TestDevice_ConfigureRequests(t *testing.T) {
	tests := []struct {
		cfg  device.Config
		err  error
		name string
		rs   []*http.Request
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
					Body:   io.NopCloser(bytes.NewBufferString(`{"params":{"config":{"enable":true}},"src":"IoTap","method":"ble.SetConfig","id":0}`)),
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
					Body:   io.NopCloser(bytes.NewBufferString(`{"src":"IoTap","method":"Shelly.Reboot","id":0}`)),
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
					Body:   io.NopCloser(bytes.NewBufferString(`{"params":{"config":{"invert":true,"name":null,"type":"switch"},"id":0},"src":"IoTap","method":"input.SetConfig","id":0}`)),
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
					Body:   io.NopCloser(bytes.NewBufferString(`{"src":"IoTap","method":"Shelly.Reboot","id":0}`)),
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
				compareRequests(t, test.rs[i], r)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}
