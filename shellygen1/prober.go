package shellygen1

import (
	"net"
	"net/http"

	"github.com/Stowify/IoTune/device"
	"github.com/Stowify/IoTune/httpclient"
)

const probePath = "settings"

// Prober implementation for the Shelly Gen1 driver.
type Prober struct{}

// ProbeRequest function implementation for the Shelly Gen1 driver.
func (p *Prober) ProbeRequest(ip net.IP) (*http.Request, device.Resource, error) {
	r, err := http.NewRequest(http.MethodGet, buildURL(ip, probePath), nil)
	if err != nil {
		return nil, nil, err
	}

	r.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

	return r, &Device{ip: ip}, nil
}
