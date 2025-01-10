package shellygen2

import (
	"net/http"

	"github.com/quetzyg/IoTap/device"
)

const securePath = "Shelly.SetAuth"

// SetCredentials for device authentication.
func (d *Device) SetCredentials(cred *device.Credentials) {
	d.cred = cred
}

// AuthConfigRequest returns an authentication setup HTTP request.
// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#shellysetauth
func (d *Device) AuthConfigRequest(auth *device.AuthConfig) (*http.Request, error) {
	// A nil auth configuration disables device authentication
	if auth == nil {
		return request(d, securePath, map[string]any{
			"user":  "admin",
			"realm": d.Realm,
			"ha1":   nil,
		})
	}

	// Check if an auth config policy is set and enforce it
	if auth.Policy != nil && auth.Policy.IsExcluded(d) {
		return nil, device.ErrPolicyExcluded
	}

	return request(d, securePath, map[string]string{
		"user":  "admin",
		"realm": d.Realm,
		"ha1":   ha1(d.Realm, auth.Credentials.Password),
	})
}
