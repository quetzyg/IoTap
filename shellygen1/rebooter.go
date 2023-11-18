package shellygen1

import (
	"net/http"

	iotune "github.com/Stowify/IoTune"
)

const rebootPath = "reboot"

// RebootRequest returns a device reboot HTTP request.
// See: https://shelly-api-docs.shelly.cloud/gen1/#reboot
func (d *Device) RebootRequest() (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, buildURL(d.ip, rebootPath), nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set(iotune.ContentTypeHeader, iotune.JSONMimeType)

	return r, nil
}
