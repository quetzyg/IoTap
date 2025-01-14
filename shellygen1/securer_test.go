package shellygen1

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/quetzyg/IoTap/device"
	"github.com/quetzyg/IoTap/httpclient"
)

func TestDevice_SetCredentials(t *testing.T) {
	shelly1 := &Device{}

	if shelly1.cred != nil {
		t.Fatalf("expected nil, got %v", shelly1.cred)
	}

	shelly1.SetCredentials(&device.Credentials{
		Username: "admin",
		Password: "admin",
	})

	if shelly1.cred.Username != "admin" {
		t.Fatalf("expected admin, got %s", shelly1.cred.Username)
	}

	if shelly1.cred.Password != "admin" {
		t.Fatalf("expected admin, got %s", shelly1.cred.Password)
	}
}

func TestDevice_AuthConfigRequest(t *testing.T) {
	tests := []struct {
		name string
		auth *device.AuthConfig
		r    *http.Request
		err  error
	}{
		{
			name: "success: turn authentication off",
			r: func() *http.Request {
				r := &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme:   "http",
						Host:     "192.168.146.123",
						Path:     securePath,
						RawQuery: "enabled=false",
					},
					Header: http.Header{},
				}

				r.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				return r
			}(),
		},
		{
			name: "failure: excluded via policy",
			auth: &device.AuthConfig{
				Policy: &device.Policy{
					Mode:   device.PolicyModeWhitelist,
					Models: []string{"SHSW-1"},
				},
			},
			err: device.ErrPolicyExcluded,
		},
		{
			name: "success: turn authentication on",
			auth: &device.AuthConfig{
				Credentials: &device.Credentials{
					Username: "foo",
					Password: "bar",
				},
			},
			r: func() *http.Request {
				r := &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme:   "http",
						Host:     "192.168.146.123",
						Path:     securePath,
						RawQuery: "enabled=true&password=bar&username=foo",
					},
					Header: http.Header{},
				}

				r.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				return r
			}(),
		},
	}

	shelly1 := &Device{ip: net.ParseIP("192.168.146.123")}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, err := shelly1.AuthConfigRequest(test.auth)

			if err == nil {
				compareRequests(t, test.r, r)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}
