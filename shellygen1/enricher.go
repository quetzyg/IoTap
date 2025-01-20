package shellygen1

import "net/http"

const settingsPath = "settings"

// EnrichRequest returns an HTTP request for device data enrichment.
// See: https://shelly-api-docs.shelly.cloud/gen1/#settings
func (d *Device) EnrichRequest() (*http.Request, error) {
	return request(d, settingsPath, nil)
}
