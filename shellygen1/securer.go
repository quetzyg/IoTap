package shellygen1

import (
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/quetzyg/IoTap/device"
	"github.com/quetzyg/IoTap/httpclient"
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
	// A nil auth configuration is considered a device authentication reset
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
		"username": []string{auth.Username},
		"password": []string{auth.Password},
	})
}

// SecureRequest decorates an HTTP request with an authorisation header.
// See: https://shelly-api-docs.shelly.cloud/gen1/#settings-login
func (d *Device) SecureRequest(r *http.Request) (*http.Request, error) {
	if d.cred == nil {
		return nil, device.ErrMissingCredentials
	}

	token := base64.StdEncoding.EncodeToString([]byte(d.cred.Username + ":" + d.cred.Password))

	r.Header.Set(httpclient.AuthorizationHeader, "Basic "+token)

	return r, nil
}
