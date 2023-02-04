package shelly

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
	Driver = "shelly"

	// Endpoint paths
	probePath = "shelly"
)

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
	return Driver
}

// buildURL for Shelly requests.
func buildURL(ip net.IP, path string) string {
	return fmt.Sprintf("http://%s/%s", ip.String(), strings.TrimPrefix(path, "/"))
}

// Prober implementation for the Shelly driver.
type Prober struct{}

// MakeRequest function implementation for the Shelly driver.
func (p *Prober) MakeRequest(ip net.IP) (*http.Request, iot.Device, error) {
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

	var je *json.SyntaxError
	if errors.As(err, &je) {
		// We found something, but it's not outputting valid JSON
		return true
	}

	return false
}
