package shellygen1

import (
	"net"
	"net/http"

	"github.com/Stowify/IoTap/device"
	"github.com/Stowify/IoTap/httpclient"
)

const probePath = "settings"

// Prober implementation for the Shelly Gen1 driver.
type Prober struct{}

// Request for probing Shelly Gen1 devices on a given IP address.
func (p *Prober) Request(ip net.IP) (*http.Request, device.Resource, error) {
	r, err := http.NewRequest(http.MethodGet, buildURL(ip, probePath), nil)
	if err != nil {
		return nil, nil, err
	}

	r.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

	return r, &Device{ip: ip}, nil
}
