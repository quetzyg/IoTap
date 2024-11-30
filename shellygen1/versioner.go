package shellygen1

import (
	"fmt"
	"net/http"

	"github.com/Stowify/IoTune/device"
)

// VersionRequest for a device version check via HTTP.
// See: https://shelly-api-docs.shelly.cloud/gen1/#ota
func (d *Device) VersionRequest() (*http.Request, error) {
	return request(d, updatePath, nil)
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
