package shellygen2

import (
	"errors"
	"net/http"
	"testing"

	"github.com/quetzyg/IoTap/device"
)

func TestDevice_Secured(t *testing.T) {
	dev.secured = true

	if dev.Secured() != true {
		t.Fatalf("expected %t, got %t", true, dev.Secured())
	}
}

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
				uri, _ := request(&Device{}, securePath, map[string]any{
					"user":  "admin",
					"realm": "shellypro1-001122334455",
					"ha1":   nil,
				})

				return uri
			}(),
		},
		{
			name: "failure: excluded via policy",
			auth: &device.AuthConfig{
				Policy: &device.Policy{
					Mode:   device.PolicyModeBlacklist,
					Models: []string{"SPSW-201XE16EU"},
				},
			},
			err: device.ErrPolicyExcluded,
		},
		{
			name: "success: turn authentication on",
			auth: &device.AuthConfig{
				Credentials: device.Credentials{
					Username: "admin",
					Password: "secret",
				},
			},
			r: func() *http.Request {
				uri, _ := request(&Device{}, securePath, map[string]string{
					"user":  "admin",
					"realm": "shellypro1-001122334455",
					"ha1":   ha1("shellypro1-001122334455", "secret"),
				})

				return uri
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, err := (&Device{model: "SPSW-201XE16EU", Realm: "shellypro1-001122334455"}).AuthConfigRequest(test.auth)

			if err == nil && !compareHTTPRequests(r, test.r) {
				t.Fatalf("expected %#v, got %#v", test.r, r)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})

	}
}
