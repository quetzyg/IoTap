package shellygen1

import (
	"encoding/json/v2"
	"fmt"
	"net/http"

	"github.com/quetzyg/IoTap/device"
)

// VersionRequest for a device version check via HTTP.
// See: https://shelly-api-docs.shelly.cloud/gen1/#ota
func (d *Device) VersionRequest() (*http.Request, error) {
	return request(d, updatePath, nil)
}

// VersionUnmarshaler returns a *json.Unmarshalers that decodes the
// JSON response from the device's version API endpoint. It is used
// to interpret the version information so the caller can determine
// whether a newer firmware release is available.
func (d *Device) VersionUnmarshaler() *json.Unmarshalers {
	return json.UnmarshalFunc(func(data []byte, dev *Device) error {
		var v struct {
			NewVersion *string `json:"new_version"`
		}

		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}

		if v.NewVersion != nil && *v.NewVersion != "" {
			dev.FirmwareNext = *v.NewVersion
			return nil
		}

		return device.ErrUnexpected
	})
}

// Outdated checks if the device's firmware is out of date.
func (d *Device) Outdated() bool {
	return d.Firmware != d.FirmwareNext
}

// UpdateDetails prints the device update information.
func (d *Device) UpdateDetails() string {
	if d.Outdated() {
		return fmt.Sprintf(device.UpdateDetailsFormat, d.Driver(), d.Name(), d.ip, d.Firmware, d.FirmwareNext)
	}

	return ""
}
