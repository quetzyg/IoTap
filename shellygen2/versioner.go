package shellygen2

import (
	"fmt"
	"net/http"

	"github.com/Stowify/IoTune/device"
)

// Request for a device version check via HTTP.
// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#shellycheckforupdate
func (d *Device) Request() (*http.Request, error) {
	return request(d, "Shelly.CheckForUpdate", nil)
}

// OutOfDate checks if the device's firmware is out of date.
func (d *Device) OutOfDate() bool {
	return d.Version != d.VersionNext
}

// UpdateDetails prints the device update information.
func (d *Device) UpdateDetails() string {
	if d.OutOfDate() {
		return fmt.Sprintf(device.UpdateDetailsFormat, d.Driver(), d.Name(), d.ip, d.Version, d.VersionNext)
	}

	return ""
}
