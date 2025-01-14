package shellygen2

import (
	"errors"
	"net/http"
	"testing"

	"github.com/quetzyg/IoTap/device"
)

func TestDevice_SetCredentials(t *testing.T) {
	shelly2 := &Device{}

	if shelly2.cred != nil {
		t.Fatalf("expected nil, got %v", shelly2.cred)
	}

	shelly2.SetCredentials(&device.Credentials{
		Username: "admin",
		Password: "admin",
	})

	if shelly2.cred.Username != "admin" {
		t.Fatalf("expected admin, got %s", shelly2.cred.Username)
	}

	if shelly2.cred.Password != "admin" {
		t.Fatalf("expected admin, got %s", shelly2.cred.Password)
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
				Credentials: &device.Credentials{
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

	shelly2 := &Device{model: "SPSW-201XE16EU", Realm: "shellypro1-001122334455"}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, err := shelly2.AuthConfigRequest(test.auth)

			switch {
			case err == nil:
				compareRequests(t, test.r, r)
				return

			case errors.Is(err, test.err):
				return

			default:
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}
