package shellygen1

import (
	"fmt"
	"net/http"

	"github.com/Stowify/IoTune/device"
)

const (
	settingsPath      = "settings"
	settingsRelayPath = "settings/relay/%d"
)

// ConfigureRequests generates a slice of *http.Requests that are to be executed in order to configure an IoT device.
func (d *Device) ConfigureRequests(cfg device.Config) ([]*http.Request, error) {
	c, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("%w: expected %q, got %q", device.ErrDriverMismatch, d.Driver(), cfg.Driver())
	}

	var requests []*http.Request

	if c.Settings != nil {
		r, err := request(d, settingsPath, c.Settings)
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)

		if c.Settings.Relay != nil {
			for i, rel := range *c.Settings.Relay {
				r, err = request(d, fmt.Sprintf(settingsRelayPath, i), rel)
				if err != nil {
					return nil, err
				}
				requests = append(requests, r)
			}
		}
	}

	return requests, nil
}
