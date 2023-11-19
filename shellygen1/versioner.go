package shellygen1

import (
	"fmt"
	"net/http"

	"github.com/Stowify/IoTune/device"
)

// VersionRequest returns a device version check HTTP request.
// See: https://shelly-api-docs.shelly.cloud/gen1/#ota
func (d *Device) VersionRequest() (*http.Request, error) {
	return request(d, updatePath, nil)
}

// UpdateAvailable checks if the device firmware can be updated.
func (d *Device) UpdateAvailable() bool {
	return d.Firmware != d.FirmwareNext
}

// UpdateDetails prints the device update information.
func (d *Device) UpdateDetails() string {
	if d.UpdateAvailable() {
		return fmt.Sprintf(device.UpdateDetailsFormat, d.Driver(), d.Name, d.ip, d.Firmware, d.FirmwareNext)
	}

	return device.UpdateDetailsNone
}
