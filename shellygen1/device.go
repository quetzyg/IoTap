package shellygen1

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	iotune "github.com/Stowify/IoTune"
	"github.com/Stowify/IoTune/device"
	"github.com/Stowify/IoTune/maputil"
)

const (
	Driver = "shelly_gen1"

	// Endpoint paths
	probePath  = "settings"
	updatePath = "ota"
)

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
		"%s, %s, %s, %s, %s, %s",
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

// UpdateRequest returns a device firmware update HTTP request.
func (d *Device) UpdateRequest() (*http.Request, error) {
	values := url.Values{}
	values.Add("update", "true")

	r, err := http.NewRequest(http.MethodPost, buildURL(d.ip, updatePath), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	r.Header.Set(iotune.ContentTypeHeader, iotune.URLEncodedFormMimeType)

	return r, nil
}

// buildURL for Shelly Gen1 requests.
func buildURL(ip net.IP, path string) string {
	return fmt.Sprintf("http://%s/%s", ip.String(), strings.TrimPrefix(path, "/"))
}

// Prober implementation for the Shelly Gen1 driver.
type Prober struct{}

// ProbeRequest function implementation for the Shelly Gen1 driver.
func (p *Prober) ProbeRequest(ip net.IP) (*http.Request, device.Resource, error) {
	r, err := http.NewRequest(http.MethodGet, buildURL(ip, probePath), nil)
	if err != nil {
		return nil, nil, err
	}

	r.Header.Set(iotune.ContentTypeHeader, iotune.JSONMimeType)

	return r, &Device{ip: ip}, nil
}
