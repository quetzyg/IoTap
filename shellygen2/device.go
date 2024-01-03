package shellygen2

import (
	"encoding/json"
	"net"

	"github.com/Stowify/IoTune/device"
	"github.com/Stowify/IoTune/maputil"
)

const (
	// Driver name for this device implementation.
	Driver = "shelly_gen2"
)

// Device implementation for the Shelly Gen2 driver.
type Device struct {
	ip      net.IP
	mac     net.HardwareAddr
	secured bool

	Key         string
	Name        string
	Model       string
	Generation  uint8
	Firmware    string
	Version     string
	VersionNext string
	AppName     string
}

// IP address of the Device.
func (d *Device) IP() net.IP {
	return d.ip
}

// MAC address of the Device.
func (d *Device) MAC() net.HardwareAddr {
	return d.mac
}

// ID returns the Device's unique identifier.
func (d *Device) ID() string {
	return d.mac.String()
}

// Driver name of this Device implementation.
func (d *Device) Driver() string {
	return Driver
}

// Secured returns true if the device requires authentication to be accessed, false otherwise.
func (d *Device) Secured() bool {
	return d.secured
}

// UnmarshalJSON implements the Unmarshaler interface.
func (d *Device) UnmarshalJSON(data []byte) error {
	var m map[string]any

	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	// Unmarshal logic for the versioner implementation
	if maputil.KeyExists(m, "result") {
		if maputil.KeyExists(m, "result.stable.version") {
			d.VersionNext = m["result"].(map[string]any)["stable"].(map[string]any)["version"].(string)
		}

		return nil
	}

	// Device unmarshal logic
	keys := []string{
		"name",
		"id",
		"mac",
		"model",
		"gen",
		"fw_id",
		"ver",
		"app",
		"auth_en",
	}

	for _, key := range keys {
		if !maputil.KeyExists(m, key) {
			return device.ErrUnexpected
		}
	}

	// Handle a potential nil name value
	name, ok := m["name"].(string)
	if !ok {
		name = "N/A"
	}
	d.Name = name
	d.Key = m["id"].(string)

	mac := device.Macify(m["mac"].(string))
	d.mac, err = net.ParseMAC(mac)
	if err != nil {
		return err
	}

	d.Model = m["model"].(string)
	d.Generation = uint8(m["gen"].(float64))
	d.Firmware = m["fw_id"].(string)
	d.Version = m["ver"].(string)

	// Assume we're on the latest version, until we version the device.
	d.VersionNext = d.Version
	d.AppName = m["app"].(string)
	d.secured = m["auth_en"].(bool)

	return nil
}
