package shellygen1

import (
	"net/http"
	"net/url"

	"github.com/quetzyg/IoTap/device"
)

const securePath = "settings/login"

// Secured returns true if the device requires authentication to be accessed, false otherwise.
func (d *Device) Secured() bool {
	return d.secured
}

// SetCredentials for device authentication.
func (d *Device) SetCredentials(cred *device.Credentials) {
	d.cred = cred
}

// AuthConfigRequest returns an authentication setup HTTP request.
// See: https://shelly-api-docs.shelly.cloud/gen1/#settings-login
func (d *Device) AuthConfigRequest(auth *device.AuthConfig) (*http.Request, error) {
	// A nil auth configuration disables device authentication
	if auth == nil {
		return request(d, securePath, url.Values{
			"enabled": []string{"false"},
		})
	}

	// Check if an auth config policy is set and enforce it
	if auth.Policy != nil && auth.Policy.IsExcluded(d) {
		return nil, device.ErrPolicyExcluded
	}

	return request(d, securePath, url.Values{
		"enabled":  []string{"true"},
		"username": []string{auth.Credentials.Username},
		"password": []string{auth.Credentials.Password},
	})
}
