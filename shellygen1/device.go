package shellygen1

import (
	"encoding/json"
	"net"

	"github.com/Stowify/IoTune/device"
	"github.com/Stowify/IoTune/maputil"
)

const (
	// Driver name for this device implementation.
	Driver = "shelly_gen1"
)

// Device implementation for the Shelly Gen1 driver.
type Device struct {
	ip net.IP

	Model        string
	Name         string
	MAC          net.HardwareAddr
	AuthEnabled  bool
	Firmware     string
	FirmwareNext string
	Discoverable bool
	NumOutputs   uint8
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
	if maputil.KeyExists(m, "new_version") {
		d.FirmwareNext = m["new_version"].(string)
		return nil
	}

	// Device unmarshal logic
	keys := []string{
		"device.type",
		"device.mac",
		"device.num_outputs",
		"login.enabled",
		"name",
		"fw",
		"discoverable",
	}

	for _, key := range keys {
		if !maputil.KeyExists(m, key) {
			return device.ErrUnexpected
		}
	}

	d.Model = m["device"].(map[string]any)["type"].(string)

	mac := device.Macify(m["device"].(map[string]any)["mac"].(string))
	d.MAC, err = net.ParseMAC(mac)
	if err != nil {
		return err
	}

	d.NumOutputs = uint8(m["device"].(map[string]any)["num_outputs"].(float64))
	d.AuthEnabled = m["login"].(map[string]any)["enabled"].(bool)

	// Handle a potential nil name value
	name, ok := m["name"].(string)
	if !ok {
		name = "N/A"
	}
	d.Name = name
	d.Firmware = m["fw"].(string)

	// Assume we're on the latest version, until the device is versioned.
	d.FirmwareNext = d.Firmware
	d.Discoverable = m["discoverable"].(bool)

	return nil
}
