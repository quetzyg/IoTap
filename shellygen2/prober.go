package shellygen2

import (
	"net"
	"net/http"

	"github.com/Stowify/IoTune/device"
	"github.com/Stowify/IoTune/httpclient"
)

const probePath = "shelly"

// Prober implementation for the Shelly Gen2 driver.
type Prober struct{}

// Request for probing Shelly Gen2 devices on a given IP address.
func (p *Prober) Request(ip net.IP) (*http.Request, device.Resource, error) {
	r, err := http.NewRequest(http.MethodGet, buildURL(ip, probePath), nil)
	if err != nil {
		return nil, nil, err
	}

	r.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

	return r, &Device{ip: ip}, nil
}
