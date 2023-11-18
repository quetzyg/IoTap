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

	Model        string           `json:"type"`
	Name         string           `json:"name"`
	MAC          net.HardwareAddr `json:"mac"`
	AuthEnabled  bool             `json:"auth"`
	Firmware     string           `json:"fw"`
	Discoverable bool             `json:"discoverable"`
	NumOutputs   uint8            `json:"num_outputs"`
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
	var tmp map[string]any

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

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
		if !maputil.KeyExists(tmp, key) {
			return device.ErrUnexpected
		}
	}

	d.Model = tmp["device"].(map[string]any)["type"].(string)

	mac := device.Macify(tmp["device"].(map[string]any)["mac"].(string))
	d.MAC, err = net.ParseMAC(mac)
	if err != nil {
		return err
	}

	d.NumOutputs = uint8(tmp["device"].(map[string]any)["num_outputs"].(float64))
	d.AuthEnabled = tmp["login"].(map[string]any)["enabled"].(bool)
	d.Name = tmp["name"].(string)
	d.Firmware = tmp["fw"].(string)
	d.Discoverable = tmp["discoverable"].(bool)

	return nil
}
