package shellygen2

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
	Driver = "shelly_gen2"
)

// buildURL for Shelly requests.
func buildURL(ip net.IP, path string) string {
	return fmt.Sprintf("http://%s/%s", ip.String(), strings.TrimPrefix(path, "/"))
}

// Device implementation for the Shelly Gen2 driver.
type Device struct {
	ip net.IP

	Key         string  `json:"id"`
	Name        string  `json:"name"`
	Model       string  `json:"model"`
	Generation  uint8   `json:"gen"`
	MAC         string  `json:"mac"`
	Firmware    string  `json:"fw_id"`
	Version     string  `json:"ver"`
	AppName     string  `json:"app"`
	AuthEnabled bool    `json:"auth_en"`
	AuthDomain  *string `json:"auth_domain"`
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

// UnmarshalJSON implements the Unmarshaler interface.
func (d *Device) UnmarshalJSON(data []byte) error {
	var tmp map[string]any

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

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
		if !maputil.KeyExists(tmp, key) {
			return device.ErrUnexpected
		}
	}

	d.Name = tmp["name"].(string)
	d.Key = tmp["id"].(string)
	d.MAC = tmp["mac"].(string)
	d.Model = tmp["model"].(string)
	d.Generation = uint8(tmp["gen"].(float64))
	d.Firmware = tmp["fw_id"].(string)
	d.Version = tmp["ver"].(string)
	d.AppName = tmp["app"].(string)
	d.AuthEnabled = tmp["auth_en"].(bool)

	authDomain := tmp["auth_domain"]
	if authDomain != nil {
		d.AuthDomain = authDomain.(*string)
	}

	return nil
}
