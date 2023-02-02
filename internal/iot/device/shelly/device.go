package shelly

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/Stowify/IoTune/internal/iot"
)

const (
	driver = "shelly"

	// Endpoint paths
	probePath = "shelly"
)

// buildURL for Shelly requests.
func buildURL(ip net.IP, path string) string {
	return fmt.Sprintf("http://%s/%s", ip.String(), strings.TrimPrefix(path, "/"))
}

// ProbeRequest function implementation for the Shelly driver.
func ProbeRequest(ip net.IP) (*http.Request, iot.Device, error) {
	r, err := http.NewRequest(http.MethodGet, buildURL(ip, probePath), nil)
	if err != nil {
		return nil, nil, err
	}

	r.Header.Set(iot.ContentTypeHeader, iot.JSONMimeType)

	return r, &Device{ip: ip}, nil
}

type Device struct {
	ip net.IP

	Model        string `json:"type"`
	MAC          string `json:"mac"`
	Auth         bool   `json:"auth"`
	Firmware     string `json:"fw"`
	Discoverable bool   `json:"discoverable"`
	LongID       int    `json:"longid"`
	NumOutputs   int    `json:"num_outputs"`
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
	return driver
}
