package shellygen2

import (
	"net"
	"net/http"

	iotune "github.com/Stowify/IoTune"
	"github.com/Stowify/IoTune/device"
)

const (
	// Endpoint paths
	probePath = "shelly"
)

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