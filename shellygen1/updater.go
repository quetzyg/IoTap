package shellygen1

import (
	"net/http"
	"net/url"
)

const updatePath = "ota"

// UpdateRequest returns a device firmware update HTTP request.
// See: https://shelly-api-docs.shelly.cloud/gen1/#ota
func (d *Device) UpdateRequest() (*http.Request, error) {
	return request(d, updatePath, url.Values{
		"update": []string{"true"},
	})
}
