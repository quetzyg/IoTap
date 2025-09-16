package shellygen2

import (
	"encoding/json/v2"
	"net"
	"net/http"

	"github.com/quetzyg/IoTap/device"
	"github.com/quetzyg/IoTap/httpclient"
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

// Unmarshaler returns a *json.Unmarshalers for decoding the JSON response
// from a network probe. On success, it hydrates a Device resource.
// This typically includes the device model, MAC address, security status,
// and firmware version. If the required fields are missing, it returns an
// error indicating an unexpected device.
func (p *Prober) Unmarshaler() *json.Unmarshalers {
	return json.UnmarshalFunc(func(data []byte, dev *Device) error {
		var v struct {
			Name     *string `json:"name"`
			Realm    *string `json:"id"`
			MAC      *string `json:"mac"`
			Model    *string `json:"model"`
			Gen      *uint8  `json:"gen"`
			Firmware *string `json:"fw_id"`
			Version  *string `json:"ver"`
			Secured  *bool   `json:"auth_en"`
		}

		err := json.Unmarshal(data, &v)
		if err != nil {
			return err
		}

		// Different Shelly generations use different JSON field names,
		// but a Gen2 device should always have these fields populated.
		if v.Realm == nil || v.MAC == nil || v.Model == nil ||
			v.Gen == nil || v.Firmware == nil || v.Version == nil ||
			v.Secured == nil {
			return device.ErrUnexpected
		}

		dev.mac, err = net.ParseMAC(device.Macify(*v.MAC))
		if err != nil {
			return err
		}

		// Default device name if unspecified
		dev.name = "N/A"
		if v.Name != nil {
			dev.name = *v.Name
		}

		dev.Realm = *v.Realm
		dev.model = *v.Model
		dev.Gen = *v.Gen
		dev.Firmware = *v.Firmware

		// Assume we're on the latest version, until the device is versioned.
		dev.VersionNext = *v.Version
		dev.Version = *v.Version

		dev.secured = *v.Secured

		return nil
	})
}
