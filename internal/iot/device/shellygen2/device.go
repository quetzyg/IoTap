package shellygen2

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/Stowify/IoTune/internal/iot"
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
		"%s (%s) %s %s @ %s",
		d.Model,
		d.Firmware,
		d.Name,
		d.MAC,
		d.ip,
	)
}

// UnmarshalJSON implements the Unmarshaler interface.
func (d *Device) UnmarshalJSON(data []byte) error {
	var tmp map[string]any

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	keys := []string{"id", "name", "model", "gen", "mac", "fw_id", "ver", "app", "auth_en", "auth_domain"}

	for _, key := range keys {
		if _, ok := tmp[key]; !ok {
			return iot.ErrWrongDevice
		}
	}

	d.Key = tmp["id"].(string)
	d.Name = tmp["name"].(string)
	d.Model = tmp["model"].(string)
	d.Generation = uint8(tmp["gen"].(float64))
	d.MAC = tmp["mac"].(string)
	d.Firmware = tmp["fw_id"].(string)
	d.Version = tmp["ver"].(string)
	d.AppName = tmp["app"].(string)
	d.AuthEnabled = tmp["auth_en"].(bool)
	d.AuthDomain = tmp["auth_domain"].(*string)

	return nil
}

// buildURL for Shelly requests.
func buildURL(ip net.IP, path string) string {
	return fmt.Sprintf("http://%s/%s", ip.String(), strings.TrimPrefix(path, "/"))
}

// Prober implementation for the Shelly Gen1 driver.
type Prober struct{}

// ProbeRequest function implementation for the Shelly Gen1 driver.
func (p *Prober) ProbeRequest(ip net.IP) (*http.Request, iot.Device, error) {
	r, err := http.NewRequest(http.MethodGet, buildURL(ip, probePath), nil)
	if err != nil {
		return nil, nil, err
	}

	r.Header.Set(iot.ContentTypeHeader, iot.JSONMimeType)

	return r, &Device{ip: ip}, nil
}

// IgnoreError checks if certain errors can be ignored.
func (p *Prober) IgnoreError(err error) bool {
	var ue *url.Error
	if errors.As(err, &ue) {
		// Ignore timeouts, refused connections and other classic HTTP shenanigans,
		// since (NORMALLY!) it means there's no such device at the IP address.
		return true
	}

	if errors.Is(err, iot.ErrWrongDevice) {
		// Ignore wrong devices.
		return true
	}

	var je *json.SyntaxError
	if errors.As(err, &je) {
		// We found something, but it's not outputting valid JSON
		return true
	}

	return false
}
