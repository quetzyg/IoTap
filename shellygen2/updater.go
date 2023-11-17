package shellygen2

import "net/http"

// UpdateRequest returns a device firmware update HTTP request.
// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#shellyupdate
func (d *Device) UpdateRequest() (*http.Request, error) {
	return request(d, "Shelly.Update", map[string]string{
		"stage": "stable",
	})
}
