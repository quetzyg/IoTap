package shellygen2

import "net/http"

// RebootRequest returns a device reboot HTTP request.
// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#shellyreboot
func (d *Device) RebootRequest() (*http.Request, error) {
	return makeRequest(d, "Shelly.Reboot", nil)
}
