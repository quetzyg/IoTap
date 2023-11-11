package shellygen1

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/Stowify/IoTune/internal/iot"
)

const (
	Driver = "shelly_gen1"

	// Endpoint paths
	probePath  = "shelly"
	updatePath = "ota"
)

// Device implementation for the Shelly Gen1 driver.
type Device struct {
	ip net.IP

	Model        string `json:"type"`
	MAC          string `json:"mac"`
	AuthEnabled  bool   `json:"auth"`
	Firmware     string `json:"fw"`
	Discoverable bool   `json:"discoverable"`
	LongID       bool   `json:"longid"` // true if the device identifies itself with its full MAC address, false otherwise
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
		"%s (%s) %s @ %s",
		d.Model,
		d.Firmware,
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

	keys := []string{"type", "mac", "auth", "fw", "discoverable", "longid", "num_outputs"}

	for _, key := range keys {
		if _, ok := tmp[key]; !ok {
			return iot.ErrWrongDevice
		}
	}

	d.Model = tmp["type"].(string)
	d.MAC = tmp["mac"].(string)
	d.AuthEnabled = tmp["auth"].(bool)
	d.Firmware = tmp["fw"].(string)
	d.Discoverable = tmp["discoverable"].(bool)
	d.LongID = tmp["longid"].(float64) == 0
	d.NumOutputs = uint8(tmp["num_outputs"].(float64))

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

	r.Header.Set(iot.ContentTypeHeader, iot.URLEncodedFormMimeType)

	return r, nil
}

// buildURL for Shelly Gen1 requests.
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
