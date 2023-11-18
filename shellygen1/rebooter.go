package shellygen1

import (
	"net/http"
)

const rebootPath = "reboot"

// RebootRequest returns a device reboot HTTP request.
// See: https://shelly-api-docs.shelly.cloud/gen1/#reboot
func (d *Device) RebootRequest() (*http.Request, error) {
	return request(d, rebootPath, nil)
}
