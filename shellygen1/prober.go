package shellygen1

import (
	"encoding/json/v2"
	"net"
	"net/http"

	"github.com/quetzyg/IoTap/device"
	"github.com/quetzyg/IoTap/httpclient"
)

const probePath = "shelly"

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

// Unmarshaler returns a *json.Unmarshalers for decoding the JSON response
// from a network probe. On success, it hydrates a Device resource.
// This typically includes the device model, MAC address, security status,
// and firmware version. If the required fields are missing, it returns an
// error indicating an unexpected device.
func (p *Prober) Unmarshaler() *json.Unmarshalers {
	return json.UnmarshalFunc(func(data []byte, dev *Device) error {
		var v struct {
			Model    *string `json:"type"`
			MAC      *string `json:"mac"`
			Secured  *bool   `json:"auth"`
			Firmware *string `json:"fw"`
		}

		err := json.Unmarshal(data, &v)
		if err != nil {
			return err
		}

		// Different Shelly generations use different JSON field names,
		// but a Gen1 device should always have these fields populated.
		if v.Model == nil || v.MAC == nil || v.Secured == nil || v.Firmware == nil {
			return device.ErrUnexpected
		}

		dev.mac, err = net.ParseMAC(device.Macify(*v.MAC))
		if err != nil {
			return err
		}

		dev.model = *v.Model

		// The /shelly endpoint for Gen1 devices does not provide
		// a name field, so we default to "N/A" and enrich later.
		dev.name = "N/A"

		// Assume we're on the latest version, until the device is versioned.
		dev.FirmwareNext = *v.Firmware
		dev.Firmware = *v.Firmware

		dev.secured = *v.Secured

		return nil
	})
}
