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

// VersionOutdated checks if the current device firmware version is out of date.
func (d *Device) VersionOutdated() bool {
	return d.Version != d.VersionNext
}

// UpgradeDetails prints the device version upgrade path (if available).
func (d *Device) UpgradeDetails() string {
	if d.VersionOutdated() {
		return fmt.Sprintf(device.UpgradeDetailsFormat, d.Driver(), d.Name, d.ip, d.Version, d.VersionNext)
	}

	return device.UpgradeDetailsNone
}
