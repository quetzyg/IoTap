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
	ip net.IP

	Key         string
	Name        string
	Model       string
	Generation  uint8
	MAC         net.HardwareAddr
	Firmware    string
	Version     string
	VersionNext string
	AppName     string
	AuthEnabled bool
	AuthDomain  *string
}

// IP address of the Device.
func (d *Device) IP() net.IP {
	return d.ip
}

// ID returns the Device's unique identifier.
func (d *Device) ID() string {
	return d.MAC.String()
}

// Driver name of this Device implementation.
func (d *Device) Driver() string {
	return Driver
}

// UnmarshalJSON implements the Unmarshaler interface.
func (d *Device) UnmarshalJSON(data []byte) error {
	var m map[string]any

	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	// Versioner unmarshal logic
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
		"auth_domain",
	}

	for _, key := range keys {
		if !maputil.KeyExists(m, key) {
			return device.ErrUnexpected
		}
	}

	d.Name = m["name"].(string)
	d.Key = m["id"].(string)

	mac := device.Macify(m["mac"].(string))
	d.MAC, err = net.ParseMAC(mac)
	if err != nil {
		return err
	}

	d.Model = m["model"].(string)
	d.Generation = uint8(m["gen"].(float64))
	d.Firmware = m["fw_id"].(string)
	d.Version = m["ver"].(string)

	// Assume we're on the latest version, until the device is versioned.
	d.VersionNext = d.Version
	d.AppName = m["app"].(string)
	d.AuthEnabled = m["auth_en"].(bool)

	authDomain := m["auth_domain"]
	if authDomain != nil {
		d.AuthDomain = authDomain.(*string)
	}

	return nil
}
