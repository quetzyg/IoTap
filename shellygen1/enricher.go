package shellygen1

import (
	"encoding/json/v2"
	"net/http"

	"github.com/quetzyg/IoTap/device"
)

const settingsPath = "settings"

// EnrichRequest returns an HTTP request for device data enrichment.
// See: https://shelly-api-docs.shelly.cloud/gen1/#settings
func (d *Device) EnrichRequest() (*http.Request, error) {
	return request(d, settingsPath, nil)
}

// EnrichUnmarshaler returns a *json.Unmarshalers for decoding additional
// metadata from a secondary data source (i.e. endpoint).
func (d *Device) EnrichUnmarshaler() *json.Unmarshalers {
	return json.UnmarshalFunc(func(data []byte, dev *Device) error {
		var v struct {
			Device struct {
				MAC *string `json:"mac"`
			} `json:"device"`
			Name *string `json:"name"`
		}

		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}

		if v.Device.MAC == nil {
			return device.ErrUnexpected
		}

		if v.Name != nil {
			dev.name = *v.Name
		}

		return nil
	})
}
