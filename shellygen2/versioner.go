package shellygen2

import (
	"fmt"
	"net/http"

	"github.com/Stowify/IoTap/device"
)

// VersionRequest for a device version check via HTTP.
// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#shellycheckforupdate
func (d *Device) VersionRequest() (*http.Request, error) {
	return request(d, "Shelly.CheckForUpdate", nil)
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
