package shellygen2

import (
	"encoding/json/v2"
	"fmt"
	"net/http"

	"github.com/quetzyg/IoTap/device"
)

// VersionRequest for a device version check via HTTP.
// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#shellycheckforupdate
func (d *Device) VersionRequest() (*http.Request, error) {
	return request(d, "Shelly.CheckForUpdate", nil)
}

// VersionUnmarshaler returns a *json.Unmarshalers that decodes the
// JSON response from the device's version API endpoint. It is used
// to interpret the version information so the caller can determine
// whether a newer firmware release is available.
func (d *Device) VersionUnmarshaler() *json.Unmarshalers {
	return json.UnmarshalFunc(func(data []byte, dev *Device) error {
		var v struct {
			Result *struct {
				Stable struct {
					Version string `json:"version"`
				} `json:"stable"`
			} `json:"result"`
		}

		err := json.Unmarshal(data, &v)
		if err != nil {
			return err
		}

		if v.Result != nil {
			if v.Result.Stable.Version != "" {
				dev.VersionNext = v.Result.Stable.Version
			}

			// We ignore other versions (i.e. beta)

			return nil
		}

		return device.ErrUnexpected
	})
}

// Outdated checks if the device's firmware is out of date.
func (d *Device) Outdated() bool {
	return d.Version != d.VersionNext
}

// UpdateDetails prints the device update information.
func (d *Device) UpdateDetails() string {
	if d.Outdated() {
		return fmt.Sprintf(device.UpdateDetailsFormat, d.Driver(), d.Name(), d.ip, d.Version, d.VersionNext)
	}

	return ""
}
