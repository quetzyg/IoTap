package shellygen1

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/Stowify/IoTune/device"
	"github.com/Stowify/IoTune/maputil"
)

const (
	// Driver name of this device implementation.
	Driver = "shelly_gen1"

	// Endpoint paths
	probePath  = "settings"
	updatePath = "ota"
	rebootPath = "reboot"
)

// buildURL for Shelly Gen1 requests.
func buildURL(ip net.IP, path string) string {
	return fmt.Sprintf("http://%s/%s", ip.String(), strings.TrimPrefix(path, "/"))
}

// Device implementation for the Shelly Gen1 driver.
type Device struct {
	ip net.IP

	Model        string `json:"type"`
	Name         string `json:"name"`
	MAC          string `json:"mac"`
	AuthEnabled  bool   `json:"auth"`
	Firmware     string `json:"fw"`
	Discoverable bool   `json:"discoverable"`
	NumOutputs   uint8  `json:"num_outputs"`
}

// IP address of the Device.
func (d *Device) IP() net.IP {
	return d.ip
}

// ID returns the Device's unique identifier.
func (d *Device) ID() string {
	return d.MAC
}

// Driver name of this Device implementation.
func (d *Device) Driver() string {
	return Driver
}

// String implements the Stringer interface.
func (d *Device) String() string {
	return fmt.Sprintf(
		"%s|%s|%-15s|%-32s|%-14s|%s",
		d.Driver(),
		d.MAC,
		d.ip,
		d.Firmware,
		d.Model,
		d.Name,
	)
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
		"mqtt.id",
		"login.enabled",
		"fw",
		"discoverable",
	}

	for _, key := range keys {
		if !maputil.KeyExists(tmp, key) {
			return device.ErrUnexpected
		}
	}

	d.Model = tmp["device"].(map[string]any)["type"].(string)
	d.MAC = tmp["device"].(map[string]any)["mac"].(string)
	d.NumOutputs = uint8(tmp["device"].(map[string]any)["num_outputs"].(float64))
	d.Name = tmp["mqtt"].(map[string]any)["id"].(string)
	d.AuthEnabled = tmp["login"].(map[string]any)["enabled"].(bool)
	d.Firmware = tmp["fw"].(string)
	d.Discoverable = tmp["discoverable"].(bool)

	return nil
}
