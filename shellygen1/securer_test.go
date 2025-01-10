package shellygen1

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/quetzyg/IoTap/device"
)

func TestDevice_SetCredentials(t *testing.T) {
	if dev.cred != nil {
		t.Fatalf("expected nil, got %v", dev.cred)
	}

	dev.SetCredentials(&device.Credentials{
		Username: "admin",
		Password: "admin",
	})

	if dev.cred.Username != "admin" {
		t.Fatalf("expected admin, got %s", dev.cred.Username)
	}

	if dev.cred.Password != "admin" {
		t.Fatalf("expected admin, got %s", dev.cred.Password)
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
				uri, _ := request(&Device{}, securePath, url.Values{
					"enabled": []string{"false"},
				})

				return uri
			}(),
		},
		{
			name: "failure: excluded via policy",
			auth: &device.AuthConfig{
				Policy: &device.Policy{
					Mode:   device.PolicyModeBlacklist,
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
				uri, _ := request(&Device{}, securePath, url.Values{
					"enabled":  []string{"true"},
					"username": []string{"foo"},
					"password": []string{"bar"},
				})

				return uri
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, err := (&Device{model: "SHSW-1"}).AuthConfigRequest(test.auth)

			if err == nil && !compareHTTPRequests(r, test.r) {
				t.Fatalf("expected %#v, got %#v", test.r, r)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})

	}
}
