package shellygen2

import (
	"fmt"
	"net/http"

	"github.com/Stowify/IoTune/device"
)

// VersionRequest returns a device version check HTTP request.
// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#shellycheckforupdate
func (d *Device) VersionRequest() (*http.Request, error) {
	return request(d, "Shelly.CheckForUpdate", nil)
}

// UpdateAvailable checks if the device firmware can be updated.
func (d *Device) UpdateAvailable() bool {
	return d.Version != d.VersionNext
}

// UpdateDetails prints the device update information.
func (d *Device) UpdateDetails() string {
	if d.UpdateAvailable() {
		return fmt.Sprintf(device.UpdateDetailsFormat, d.Driver(), d.Name, d.ip, d.Version, d.VersionNext)
	}

	return ""
}
