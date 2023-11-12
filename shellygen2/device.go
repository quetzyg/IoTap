package shellygen2

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	iotune "github.com/Stowify/IoTune"
	"github.com/Stowify/IoTune/device"
	"github.com/Stowify/IoTune/maputil"
)

const (
	Driver = "shelly_gen2"

	// Endpoint paths
	probePath = "shelly"
)

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

// String implements the Stringer interface.
func (d *Device) String() string {
	return fmt.Sprintf(
		"[%s] %s %s @ %s (%s) - %s",
		d.Driver(),
		d.Firmware,
		d.Model,
		d.ip,
		d.MAC,
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

// buildURL for Shelly requests.
func buildURL(ip net.IP, path string) string {
	return fmt.Sprintf("http://%s/%s", ip.String(), strings.TrimPrefix(path, "/"))
}

// Prober implementation for the Shelly Gen2 driver.
type Prober struct{}

// ProbeRequest function implementation for the Shelly Gen2 driver.
func (p *Prober) ProbeRequest(ip net.IP) (*http.Request, device.Resource, error) {
	r, err := http.NewRequest(http.MethodGet, buildURL(ip, probePath), nil)
	if err != nil {
		return nil, nil, err
	}

	r.Header.Set(iotune.ContentTypeHeader, iotune.JSONMimeType)

	return r, &Device{ip: ip}, nil
}

// UpdateRequest returns a device firmware update HTTP request.
func (d *Device) UpdateRequest() (*http.Request, error) {
	return makeRequest(d, "Shelly.Update", map[string]string{
		"stage": "stable",
	})
}
